package signup

//go:generate go run generate.go

import (
	"net/http"

	"github.com/influx6/devapp/internals/users"
	"github.com/influx6/devapp/static"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/tmplutil"
)

var (
	views = tmplutil.New().
		Add("index.layout", static.MustReadFile("layout/index.html", true)).
		Add("signup.content", MustReadFile("signup.html", true))
)

// Render renders the page for the home view.
func Render(ctx *httputil.Context) error {
	tmpl, err := views.FromWith(httputil.TextContextFunctions(ctx), "index.layout", "signup.content")
	if err != nil {
		return err
	}

	var newUser users.NewUser

	if cachedUser, ok := ctx.Bag().Get(users.NillNewUser); ok {
		newUser = cachedUser.(users.NewUser)
	}

	ctx.AddHeader("Content-Type", "text/html")
	return ctx.Template(http.StatusOK, tmpl, struct{ CachedUser users.NewUser }{CachedUser: newUser})
}
