package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/influx6/backoffice/models/session"
	"github.com/influx6/devapp/internals/sessions"
	"github.com/influx6/devapp/internals/sessions/db"
	"github.com/influx6/devapp/internals/sessions/mdb"
	tokensdb "github.com/influx6/devapp/internals/tokens/db"
	tokensmdb "github.com/influx6/devapp/internals/tokens/mdb"
	users "github.com/influx6/devapp/internals/users"
	userdbapi "github.com/influx6/devapp/internals/users/db"
	userdb "github.com/influx6/devapp/internals/users/mdb"

	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
)

// errors ...
var (
	ErrFailedOTPAuth        = errors.New("Failed to match otp authorization")
	ErrTwoFactorRequired    = errors.New("Session requires twofactor authentication")
	ErrTwoFactorDone        = errors.New("Session has already completed twofactor authentication")
	ErrTwoFactorNotRequired = errors.New("Session requires no twofactor authentication")
	ErrFailedLogin          = errors.New("Failed to authenticate login credentails or no AUTH credentials found")
)

// SessionAPI implements the http api for responding to session request.
type SessionAPI struct {
	DB       *mdb.SessionDB
	UserDB   *userdb.UserDB
	TokensDB *tokensmdb.TokenRecordDB
}

// TwoFactorAuthorizationCheck validates that a given user has being logged in, then checks if
// two factor authorization is enabled and if it has not be completed on the session of the user,
// if will direct to URL for two factor authentication. If user has already being authenticated with
// two factor then nothing is done.
// Must be used in conjunction with SessionAPI.Authenticate or SessionAPI.Login.
func (s SessionAPI) TwoFactorAuthorizationCheck(ctx *httputil.Context) error {
	userrec, ok := ctx.Bag().Get(users.NilUser)
	if !ok {
		return errors.New("No User stored from logged in session")
	}

	user, ok := userrec.(users.User)
	if !ok {
		return httputil.HTTPError{
			Err:  errors.New("User type for User key is invalid"),
			Code: http.StatusInternalServerError,
		}
	}

	sessionrec, ok := ctx.Bag().Get(sessions.NilSession)
	if !ok {
		return errors.New("No Session currently available")
	}

	session, ok := sessionrec.(sessions.Session)
	if !ok {
		return errors.New("No Session currently available")
	}

	if !user.UseTwoFactor {
		return nil
	}

	if user.UseTwoFactor && session.TwoFactorDone {
		return nil
	}

	return ErrTwoFactorRequired
}

// TwoFactorAuthorization receives the provided incoming authorization user token
// and validates that the provided user has the correct authorization token, which
// will be validated else fail. It expects to receive a `token` param, which contains
// user provided token.
// Must be used in conjunction with SessionAPI.Authenticate or SessionAPI.Login.
func (s SessionAPI) TwoFactorAuthorization(ctx *httputil.Context) error {
	userrec, ok := ctx.Bag().Get(users.NilUser)
	if !ok {
		return errors.New("No User stored from logged in session")
	}

	user, ok := userrec.(users.User)
	if !ok {
		return httputil.HTTPError{
			Err:  errors.New("User type for User key is invalid"),
			Code: http.StatusInternalServerError,
		}
	}

	if !user.UseTwoFactor {
		return nil
	}

	sessionrec, ok := ctx.Bag().Get(sessions.NilSession)
	if !ok {
		return errors.New("No Session currently available")
	}

	session, ok := sessionrec.(sessions.Session)
	if !ok {
		return errors.New("No Session currently available")
	}

	token, ok := ctx.Bag().GetString("token")
	token = strings.TrimSpace(token)
	if !ok || token == "" {
		ctx.SetFlash("error", "Token is not provided")
		return errors.New("No twofactor token provided")
	}

	usedToken, err := tokensdb.UsedToken(ctx, ctx.Metrics(), s.TokensDB, user, token)
	if err != nil {
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	if usedToken {
		ctx.SetFlash("error", "Token provided is already used: Try again")
		return httputil.HTTPError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("User already used token"),
		}
	}

	ctx.Metrics().Emit(metrics.Info("Recieved TwoFactor Token for Authorization").With("token", token))

	if err := user.ValidateOTP(token); err != nil {
		ctx.SetFlash("error", "Token provided is invalid: Does not match users OTP")

		// Update totp to ensure time details are properly preserved.
		if tokenerr := userdbapi.UpdateTOTP(ctx, ctx.Metrics(), s.UserDB, user); tokenerr != nil {
			ctx.Metrics().Emit(metrics.YellowAlert(tokenerr, "User DB TOTP Update Failed"))
		}

		return err
	}

	if err := tokensdb.AddToken(ctx, ctx.Metrics(), s.TokensDB, user, token); err != nil {
		return httputil.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("Failed to update user tokens TOTP Update: %+q", err),
		}
	}

	if err := userdbapi.UpdateTOTP(ctx, ctx.Metrics(), s.UserDB, user); err != nil {
		return httputil.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("User DB TOTP Update Failed: %+q", err),
		}
	}

	ctx.SetFlash("notice", "Token validated")
	session.TwoFactorDone = true

	if err := s.DB.Update(ctx, session.PublicID, session); err != nil {
		return err
	}

	ctx.Bag().Set(sessions.NilSession, session)

	return nil
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
			Err:  ErrFailedLogin,
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

	ruser, err := userdbapi.Get(ctx, ctx.Metrics(), s.UserDB, sessionUserID)
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
		return ErrFailedLogin
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

	ctx.Bag().Set(users.NilUser, ruser)
	ctx.Bag().Set(sessions.NilSession, userSession)

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

	userSession, err := db.Get(ctx, ctx.Metrics(), s.DB, user.PublicID)
	if err != nil {
		if userSession, err = db.Create(ctx, ctx.Metrics(), s.DB, httputil.FourthyEightHoursDuration, user); err != nil {
			return httputil.HTTPError{
				Err:  fmt.Errorf("Unable to create session: %+q", err),
				Code: http.StatusInternalServerError,
			}
		}
	}

	// if Session has expired, renew and reset for two factor authentication.
	if userSession.Expired() {
		userSession.TwoFactorDone = false
		userSession.Expires = time.Now().Add(httputil.FourthyEightHoursDuration)
		if err := db.Update(ctx, ctx.Metrics(), s.DB, userSession); err != nil {
			return httputil.HTTPError{
				Err:  fmt.Errorf("Unable to renew user session: %+q", err),
				Code: http.StatusInternalServerError,
			}
		}
	}

	ctx.Bag().Set(users.NilUser, user)
	ctx.Bag().Set(sessions.NilSession, userSession)

	authValue := fmt.Sprintf("Bearer %s", userSession.SessionToken())
	privateid := base64.StdEncoding.EncodeToString([]byte(authValue))

	ctx.SetHeader("Authorization", authValue)
	ctx.SetCookie(&http.Cookie{
		Name:    "Authorization",
		Value:   privateid,
		Expires: userSession.Expires,
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
