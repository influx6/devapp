package signup

//go:generate go run generate.go

import (
	"net/http"

	userapi "github.com/influx6/devapp/internals/users/handler"
	"github.com/influx6/devapp/static"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/tmplutil"
)

var (
	views = tmplutil.New().
		Add("index.layout", static.MustReadFile("templates/index.html", true)).
		Add("signup.content", MustReadFile("signup.html", true))
)

// Render renders the page for the home view.
func Render(ctx *httputil.Context) error {
	tmpl, err := views.From("index.layout", "signup.content")
	if err != nil {
		return err
	}

	ctx.AddHeader("Content-Type", "text/html")
	return ctx.Template(http.StatusOK, tmpl, struct{}{})
}

// SingupHandler returns a middleware function which wraps any handler
// to validate authorization else responding with a http.StatusUnauthorized.
func SignupHandler(api userapi.UserAPI, to string) httputil.Handler {
	return func(ctx *httputil.Context) error {
		if err := api.CreateUserFromURLEncoded(ctx); err != nil {
			ctx.Metrics().Emit(metrics.Error(err).WithMessage("Failed to create user"))
			return err
		}

		ctx.Metrics().Emit(metrics.Info("Created new user succesfully"))
		return ctx.Redirect(http.StatusTemporaryRedirect, to)
	}
}
