package handler

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/influx6/devapp/internals/users"
	"github.com/influx6/devapp/internals/users/db"
	"github.com/influx6/devapp/internals/users/mdb"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
	"github.com/sec51/twofactor"
)

// UserAPI exposes controller methods for user http apis
type UserAPI struct {
	DB *mdb.UserDB
}

// EnableTwoFactor enables two factor autnetication and directs user to run to
// UserQRPage to retrieve user qr code for authentication.
// This handle must be used in conguction with Session.Authenticate
// has a user must be logged in and providing it's authentication token in header.
func (u UserAPI) EnableTwoFactor(ctx *httputil.Context) error {
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

	totp, err := twofactor.NewTOTP(user.PublicID, users.TwofactorOrg, crypto.SHA1, 6)
	if err != nil {
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	totpbytes, err := totp.ToBytes()
	if err != nil {
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	user.TOTP = base64.StdEncoding.EncodeToString(totpbytes)
	user.UseTwoFactor = true

	if err := u.DB.Update(ctx, user.PublicID, user); err != nil {
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	ctx.Bag().Set(users.NilUser, user)

	return nil
}

// SetTwoFactorAsSeen sets the two factor state of the user has already seen.
func (u UserAPI) SetTwoFactorAsSeen(ctx *httputil.Context) error {
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

	user.SeenTwoFactor = true

	if err := u.DB.Update(ctx, user.PublicID, user); err != nil {
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	ctx.Bag().Set(users.NilUser, user)

	return nil
}

// DisableTwoFactor disables two factor autnetication for a authenticated user.
// This handle must be used in conguction with Session.Authenticate
// has a user must be logged in and providing it's authentication token in header.
func (u UserAPI) DisableTwoFactor(ctx *httputil.Context) error {
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

	user.TOTP = ""
	user.UseTwoFactor = false
	user.SeenTwoFactor = false

	if err := u.DB.Update(ctx, user.PublicID, user); err != nil {
		return httputil.HTTPError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	ctx.Bag().Set(users.NilUser, user)

	return nil
}

// UserTwoFactorQRImage returns a png/image response which is written to
// as response for the user qr for google authenticator.
// This handle must be used in conguction with Session.Authenticate
// has a user must be logged in.
func (u UserAPI) UserTwoFactorQRImage(ctx *httputil.Context) error {
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

	qr, err := user.TwoFactorQR()
	if err != nil {
		return err
	}

	return ctx.Stream(http.StatusOK, "image/png", bytes.NewBuffer(qr))
}

// UserTwoFactorQTURL returns a text/plain response which is the URL
// for the given user twofactor key URL for use with google authenticator.
func (u UserAPI) UserTwoFactorQTURL(ctx *httputil.Context) error {
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

	qr, err := user.TwoFactorURL()
	if err != nil {
		return err
	}

	return ctx.Stream(http.StatusOK, "text/plain", bytes.NewBufferString(qr))
}

// UpdateUserTOTP attempts to update the totp field for the user.
func (u UserAPI) UpdateUserTOTP(ctx *httputil.Context) error {
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

	if err := db.UpdateTOTP(ctx, ctx.Metrics(), u.DB, user); err != nil {
		ctx.Metrics().Emit(metrics.Error(err).WithMessage("Failed to update user totp").With("user", user.PublicID))
		return err
	}

	return nil
}

// CreateUserFromURLEncoded provides the http api function for creting a user.
func (u UserAPI) CreateUserFromURLEncoded(ctx *httputil.Context) error {
	var ok bool
	var username, password, passwordConfirm string

	defer func() {
		ctx.Bag().Set(users.NillNewUser, users.NewUser{Username: username, Password: password, PasswordConfirm: passwordConfirm})
	}()

	username, ok = ctx.Bag().GetString("username")
	username = strings.TrimSpace(username)
	if !ok || username == "" {
		ctx.SetFlash("error", "User username not provided!")
		return errors.New("Must provided username value")
	}

	password, ok = ctx.Bag().GetString("password")
	password = strings.TrimSpace(password)
	if !ok || password == "" {
		ctx.SetFlash("error", "User password not provided!")
		return errors.New("Must provided password value")
	}

	passwordConfirm, ok = ctx.Bag().GetString("password_confirm")
	passwordConfirm = strings.TrimSpace(passwordConfirm)
	if !ok || passwordConfirm == "" {
		ctx.SetFlash("error", "User password confirm not provided!")
		return errors.New("Must provided password confirm value")
	}

	if passwordConfirm != password {
		ctx.SetFlash("error", "Password and  PasswordConfirm do not match!")
		return errors.New("Password and  PasswordConfirm do not match")
	}

	var newuser users.NewUser
	newuser.Username = username
	newuser.Password = password
	newuser.PasswordConfirm = passwordConfirm

	ctx.Metrics().Emit(metrics.Info("Creating new user (Content: application/x-www-form-urlencoded)").With("user", newuser))

	_, err := db.Create(ctx, ctx.Metrics(), u.DB, newuser)
	if err != nil {
		ctx.SetFlash("error", "User was not created successfully!")
		ctx.Metrics().Emit(metrics.Error(err).WithMessage("Creating new user(Content: application/x-www-form-urlencoded)").With("user", newuser))
		return err
	}

	ctx.SetFlash("success", "User created successfully!")
	ctx.Metrics().Emit(metrics.Info("Created new user (Content: application/x-www-form-urlencoded)").With("user", newuser))

	return nil
}
