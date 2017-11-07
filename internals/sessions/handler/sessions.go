package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/influx6/backoffice/models/session"
	"github.com/influx6/devapp/internals/sessions/db"
	"github.com/influx6/devapp/internals/sessions/mdb"
	userdbapi "github.com/influx6/devapp/internals/users/db"
	userdb "github.com/influx6/devapp/internals/users/mdb"

	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
)

// SessionAPI implements the http api for responding to session request.
type SessionAPI struct {
	DB     *mdb.SessionDB
	UserDB *userdb.UserDB
}

// Authenticate handles relogin request for a previously authenticated/logged in user, returning
// appropriate token for session.
// HTTP Method: GET
// Header:
// 		{
// 			"Authorization":"Bearer <TOKEN>",
// 		}
//
// 		WHERE: <TOKEN> = <USERID>:<SESSIONTOKEN>
//
func (s SessionAPI) Authenticate(ctx *httputil.Context) error {
	authorization, err := s.GetAuthroization(ctx)
	if err != nil {
		return httputil.HTTPError{
			Code: http.StatusBadRequest,
			Err:  errors.New("Has no 'Authorization' header"),
		}
	}

	ctx.Metrics().Emit(metrics.Info("Retreived Authorization Value").With("authorization", authorization))

	authtype, token, err := httputil.ParseAuthorization(authorization)
	if err != nil {
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	if authtype != "Bearer" {
		err = errors.New("Only `Bearer` Authorization supported")
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	// Retrieve Authorization UserID and Token.
	sessionUserID, sessionToken, err := session.ParseToken(token)
	if err != nil {
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	ctx.Metrics().Emit(metrics.Info("Retreived Session User record").With("user", sessionUserID))

	_, err = userdbapi.Get(ctx, ctx.Metrics(), s.UserDB, sessionUserID)
	if err != nil {
		ctx.Metrics().Emit(metrics.Error(err).WithMessage("Failed to retreived Session User record").With("user", sessionUserID))
		return httputil.HTTPError{
			Err:  fmt.Errorf("User not found for %+q: %+q", sessionUserID, err),
			Code: http.StatusBadRequest,
		}
	}

	userSession, err := db.Get(ctx, ctx.Metrics(), s.DB, sessionUserID)
	if err != nil {
		ctx.Metrics().Emit(metrics.Error(err).WithMessage("Failed to retreived Session record").With("user", sessionUserID))
		return err
	}

	ctx.Metrics().Emit(metrics.Info("Retreived Session record").
		With("authorization", authorization).
		With("session", userSession.Fields()))

	if !userSession.ValidateToken(sessionToken) {
		err := errors.New("Invalid user session's token")
		ctx.Metrics().Emit(metrics.Error(err).WithMessage("Failed to validate token").
			With("user", sessionUserID).
			With("token", sessionToken).
			With("session", userSession))
		return err
	}

	if userSession.Expired() {
		err := errors.New("User session has expired")
		ctx.Metrics().Emit(metrics.Error(err).With("user", sessionUserID))

		if derr := db.Delete(ctx, ctx.Metrics(), s.DB, userSession.PublicID); derr != nil {
			ctx.Metrics().Emit(metrics.Error(derr).WithMessage("Failed to remove expired session record").With("user", sessionUserID))
			return derr
		}

		return err
	}

	ctx.Metrics().Emit(metrics.Info("User session is valid and authenticated").With("user", sessionUserID))

	return nil
}

// Login handles login request and authenticates user, returning
// appropriate token for session.
// HTTP Method: GET
// Optional Header:
// 		{
// 			"Authorization":"Bearer <TOKEN>",
// 		}
//
// 		WHERE: <TOKEN> = <USERID>:<SESSIONTOKEN>
//
func (s SessionAPI) Login(ctx *httputil.Context) error {
	if err := s.Authenticate(ctx); err == nil {
		return nil
	}

	username, ok := ctx.Bag().GetString("username")
	if !ok {
		return errors.New("username not provided")
	}

	password, ok := ctx.Bag().GetString("password")
	if !ok {
		return errors.New("password not provided")
	}

	user, err := userdbapi.GetByUsername(ctx, ctx.Metrics(), s.UserDB, username)
	if err != nil {
		return err
	}

	if err = user.Authenticate(password); err != nil {
		return httputil.HTTPError{
			Err:  fmt.Errorf("Invalid Credentials: %+q", err),
			Code: http.StatusUnauthorized,
		}
	}

	newSession, err := db.Create(ctx, ctx.Metrics(), s.DB, httputil.FourthyEightHoursDuration, user)
	if err != nil {
		return httputil.HTTPError{
			Err:  fmt.Errorf("Unable to create session: %+q", err),
			Code: http.StatusInternalServerError,
		}
	}

	authValue := fmt.Sprintf("Bearer %s", newSession.SessionToken())
	privateid := base64.StdEncoding.EncodeToString([]byte(authValue))

	ctx.SetHeader("Authorization", authValue)
	ctx.SetCookie(&http.Cookie{
		Name:    "Authorization",
		Value:   privateid,
		Expires: newSession.Expires,
		Path:    "/",
	})

	ctx.Metrics().Emit(metrics.Info("Created Session").With("private_id", privateid).With("user_id", user.PublicID))

	return nil
}

// Logout handles logout request of a user, deleting
// session for user in db.
// HTTP Method: GET
// Header:
// 		{
// 			"Authorization":"Bearer <TOKEN>",
// 		}
//
// 		WHERE: <TOKEN> = <USERID>:<SESSIONTOKEN>
//
func (s SessionAPI) Logout(ctx *httputil.Context) error {
	authorization, err := s.GetAuthroization(ctx)
	if err != nil {
		return httputil.HTTPError{
			Code: http.StatusBadRequest,
			Err:  errors.New("Has no 'Authorization' header"),
		}
	}

	authtype, token, err := httputil.ParseAuthorization(authorization)
	if err != nil {
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	if authtype != "Bearer" {
		err = errors.New("Only `Bearer` Authorization supported")
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	// Retrieve Authorization UserID and Token.
	sessionUserID, sessionToken, err := session.ParseToken(token)
	if err != nil {
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	_, err = userdbapi.Get(ctx, ctx.Metrics(), s.UserDB, sessionUserID)
	if err != nil {
		return httputil.HTTPError{
			Err:  fmt.Errorf("User not found for %+q: %+q", sessionUserID, err),
			Code: http.StatusBadRequest,
		}
	}

	userSession, err := db.Get(ctx, ctx.Metrics(), s.DB, sessionUserID)
	if err != nil {
		return err
	}

	if !userSession.ValidateToken(sessionToken) {
		err = errors.New("Invalid user session's token")
		return err
	}

	if err = db.Delete(ctx, ctx.Metrics(), s.DB, userSession.UserID); err != nil {
		return err
	}

	return nil
}

// GetAuthroization returns authorization value for giving request.
func (s SessionAPI) GetAuthroization(ctx *httputil.Context) (string, error) {
	if ctx.HasHeader("Authorization", "") {
		return ctx.GetHeader("Authorization"), nil
	}

	for _, cookie := range ctx.Cookies() {
		if strings.ToLower(cookie.Name) == "authorization" {
			val, err := base64.StdEncoding.DecodeString(cookie.Value)
			return string(val), err
		}
	}

	return "", errors.New("no valid authorization found")
}
