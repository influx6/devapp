package signin

//go:generate go run generate.go

import (
	"net/http"

	sessionapi "github.com/influx6/devapp/internals/sessions/handler"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
)

// AuthenticateMW returns a middleware function which wraps any handler
// to validate authorization else responding with a http.StatusUnauthorized.
func AuthenticateMW(api sessionapi.SessionAPI) httputil.Middleware {
	return func(next httputil.Handler) httputil.Handler {
		return func(ctx *httputil.Context) error {
			if err := api.Authenticate(ctx); err != nil {
				ctx.Metrics().Emit(metrics.Error(err).WithMessage("Failed to authenticate request"))
				return ctx.Redirect(http.StatusTemporaryRedirect, "/")
			}

			ctx.Metrics().Emit(metrics.Info("Request authenticated"))
			return next(ctx)
		}
	}
}

// LogoutHandler returns a Handler for handling user logout.
func LogoutHandler(api sessionapi.SessionAPI, to string) httputil.Handler {
	return func(ctx *httputil.Context) error {
		if err := api.Logout(ctx); err != nil {
			ctx.Metrics().Emit(metrics.Error(err).WithMessage("Failed to end user session"))
		}

		return ctx.Redirect(http.StatusTemporaryRedirect, to)
	}
}

// IfLoggedHandler returns a Handler which will always redirect to the provided to path
// if authorization is validated, else calls the next handler.
func IfLoggedHandler(api sessionapi.SessionAPI, to string) httputil.Middleware {
	return func(next httputil.Handler) httputil.Handler {
		return func(ctx *httputil.Context) error {
			if err := api.Authenticate(ctx); err != nil {
				ctx.Metrics().Emit(metrics.Error(err).WithMessage("Failed to authenticate request"))
				return next(ctx)
			}

			ctx.Metrics().Emit(metrics.Info("Request already authenticated, redirecting..."))
			return ctx.Redirect(http.StatusTemporaryRedirect, to)
		}
	}
}

// LoginHandler returns a Handler for handling user login.
func LoginHandler(api sessionapi.SessionAPI, to string, failed string) httputil.Handler {
	return func(ctx *httputil.Context) error {
		if err := api.Login(ctx); err != nil {
			ctx.SetHeader("X-Login-Error", err.Error())
			ctx.Metrics().Emit(metrics.Error(err).WithMessage("Failed authenticate/login user, redirecting...").With("to", failed))
			return ctx.Redirect(http.StatusTemporaryRedirect, failed)
		}

		ctx.Metrics().Emit(metrics.Info("User authenticated successfully").With("to", to))
		return ctx.Redirect(http.StatusTemporaryRedirect, to)
	}
}
