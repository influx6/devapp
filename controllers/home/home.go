package home

import "github.com/influx6/faux/httputil"

// Render renders the page for the home view.
func Render(ctx *httputil.Context) error {
	return ctx.HTMLBlob(200, MustReadFileByte("home.html", true))
}
