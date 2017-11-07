package profile

//go:generate go run generate.go

import (
	"net/http"

	"github.com/influx6/devapp/static"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/tmplutil"
)

var (
	views = tmplutil.New().
		Add("index.layout", static.MustReadFile("templates/index.html", true)).
		Add("profile.content", MustReadFile("profile.html", true))
)

// Render renders the page for the home view.
func Render(ctx *httputil.Context) error {
	tmpl, err := views.From("index.layout", "profile.content")
	if err != nil {
		return err
	}

	ctx.AddHeader("Content-Type", "text/html")
	return ctx.Template(http.StatusOK, tmpl, struct{}{})
}