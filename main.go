package main

//go:generate go generate ./controllers/home/...
//go:generate go generate ./controllers/signup/...
//go:generate go generate ./static/...

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/influx6/devapp/controllers/home"
	"github.com/influx6/devapp/static"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/custom"
)

var (
	logs     = metrics.New(custom.FlatDisplay(os.Stdout))
	port     = flag.String("port", envOrDefault("PORT", "3000"), "-port=3000 sets the port of the http server")
	staticFS = httputil.VirtualFileSystem{
		GetFileFunc: func(path string) (*httputil.VirtualFile, error) {
			path = strings.TrimPrefix(path, "/")
			reader, dataSize, err := static.FindFile(path, false)
			if err != nil {
				return nil, err
			}

			bureader, ok := reader.(*bytes.Reader)
			if !ok {
				return nil, errors.New("Expected bytes.Reader type")
			}

			return httputil.NewVirtualFile(bureader, path, dataSize, time.Now()), nil
		},
	}
)

func main() {
	flag.Parse()

	mw := httputil.MWi(httputil.MetricsMW(logs), httputil.LogMW)
	prefixmw := httputil.MWi(mw, httputil.StripPrefixMW("/static/"))

	m := mux.NewRouter()
	m.NotFoundHandler = httputil.HTTPFunc(httputil.NotFound)

	m.HandleFunc("/", httputil.HTTPFunc(mw(home.Render)))
	m.PathPrefix("/static/").Handler(httputil.GzipServer(staticFS, true, prefixmw))

	server, err := httputil.Listen(false, fmt.Sprintf(":%s", *port), m)
	if err != nil {
		log.Fatalf("Failed to start server: %+q", err)
		return
	}

	server.Wait()
}

func envOrDefault(name string, def string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}

	return def
}
