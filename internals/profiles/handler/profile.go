package handler

import (
	"errors"
	"net/http"

	sessionsapi "github.com/influx6/devapp/internals/sessions/handler"
	"github.com/influx6/devapp/internals/users"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
)

// ProfileAPI implements http methods for interfacting with the profile records.
type ProfileAPI struct {
	Sessions sessionsapi.SessionAPI
}

func (p ProfileAPI) getName(ctx *httputil.Context) error {
	if userrec, ok := ctx.Bag().Get(users.NilUser); ok {
		if user, ok := userrec.(users.User); ok {
			return ctx.Blob(http.StatusOK, "text/plain", []byte(user.Username))
		}
	}

	err := errors.New("No User record instance found in context")
	ctx.Metrics().Emit(metrics.Error(err).WithMessage("Failed to retreived Session User record from context"))
	return httputil.HTTPError{
		Err:  err,
		Code: http.StatusBadRequest,
	}
}

// Get returns an appropriate response for a get request.
func (p ProfileAPI) Get(ctx *httputil.Context) error {
	return p.Sessions.AuthenticateHandlers(p.getName)(ctx)
}
