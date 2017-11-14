package twoauth

//go:generate go run generate.go

import (
	"errors"
	"net/http"

	"github.com/influx6/devapp/internals/users"
	"github.com/influx6/devapp/static"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/tmplutil"
)

var (
	views = tmplutil.New().
		Add("index.layout", static.MustReadFile("layout/index.html", true)).
		Add("twofactor.content", MustReadFile("twoauth.html", true)).
		Add("twofactor-qr.content", MustReadFile("twoauth-qr.html", true))
)

// Render renders the page for the two-factor authentication.
func Render(ctx *httputil.Context) error {
	userrec, ok := ctx.Bag().Get(users.NilUser)
	if !ok {
		return errors.New("No User stored from logged in session")
	}

	user, ok := userrec.(users.User)
	if !ok {
		return errors.New("User type for User key is invalid")
	}

	tmpl, err := views.FromWith(httputil.TextContextFunctions(ctx), "index.layout", "twofactor.content")
	if err != nil {
		return err
	}

	ctx.AddHeader("Content-Type", "text/html")
	return ctx.Template(http.StatusOK, tmpl, struct {
		User users.User
	}{
		User: user,
	})
}

// RenderQR renders the page for the showing twofactor authentication QR image.
func RenderQR(ctx *httputil.Context) error {
	userrec, ok := ctx.Bag().Get(users.NilUser)
	if !ok {
		return errors.New("No User stored from logged in session")
	}

	user, ok := userrec.(users.User)
	if !ok {
		return errors.New("User type for User key is invalid")
	}

	if !user.UseTwoFactor {
		return errors.New("User has not enabled twofactor")
	}

	tmpl, err := views.FromWith(httputil.TextContextFunctions(ctx), "index.layout", "twofactor-qr.content")
	if err != nil {
		return err
	}

	ctx.AddHeader("Content-Type", "text/html")
	return ctx.Template(http.StatusOK, tmpl, struct {
		User users.User
	}{
		User: user,
	})
}
