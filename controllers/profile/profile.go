package profile

//go:generate go run generate.go

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/influx6/devapp/internals/users"
	"github.com/influx6/devapp/static"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/tmplutil"
)

var (
	views = tmplutil.New().
		Add("index.layout", static.MustReadFile("layout/index.html", true)).
		Add("profile.content", MustReadFile("profile.html", true))
)

// Render renders the page for the home view.
func Render(ctx *httputil.Context) error {
	userrec, ok := ctx.Bag().Get(users.NilUser)
	if !ok {
		return errors.New("No User stored from logged in session")
	}

	user, ok := userrec.(users.User)
	if !ok {
		return errors.New("User type for User key is invalid")
	}

	tmpl, err := views.FromWith(httputil.TextContextFunctions(ctx), "index.layout", "profile.content")
	if err != nil {
		return err
	}

	ctx.SetFlash("notice", fmt.Sprintf("Welcome to your profile %+q", user.Username))
	ctx.AddHeader("Content-Type", "text/html")
	return ctx.Template(http.StatusOK, tmpl, struct {
		User users.User
	}{
		User: user,
	})
}
