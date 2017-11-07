package handler

import (
	"errors"

	"github.com/influx6/devapp/internals/users"
	"github.com/influx6/devapp/internals/users/db"
	"github.com/influx6/devapp/internals/users/mdb"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
)

// UserAPI exposes controller methods for user http apis
type UserAPI struct {
	DB *mdb.UserDB
}

// CreateUserFromURLEncoded provides the http api function for creting a user.
func (u UserAPI) CreateUserFromURLEncoded(ctx *httputil.Context) error {
	username, ok := ctx.Bag().GetString("username")
	if !ok {
		return errors.New("Must provided username value")
	}

	password, ok := ctx.Bag().GetString("password")
	if !ok {
		return errors.New("Must provided password value")
	}

	passwordConfirm, ok := ctx.Bag().GetString("password_confirm")
	if !ok {
		return errors.New("Must provided password confirm value")
	}

	var newuser users.NewUser
	newuser.Username = username
	newuser.Password = password
	newuser.PasswordConfirm = passwordConfirm

	ctx.Metrics().Emit(metrics.Info("Creating new user (Content: application/x-www-form-urlencoded)").With("user", newuser))

	_, err := db.Create(ctx, ctx.Metrics(), u.DB, newuser)
	if err != nil {
		ctx.Metrics().Emit(metrics.Error(err).WithMessage("Creating new user(Content: application/x-www-form-urlencoded)").With("user", newuser))
		return err
	}

	// ctx.Status(http.StatusOK)
	ctx.Metrics().Emit(metrics.Info("Created new user (Content: application/x-www-form-urlencoded)").With("user", newuser))

	return nil
}
