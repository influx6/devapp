package main

//go:generate go generate ./controllers/home/...
//go:generate go generate ./controllers/signup/...
//go:generate go generate ./static/...

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/influx6/devapp/controllers/home"
	"github.com/influx6/faux/httputil"
)

var port = flag.String("port", envOrDefault("PORT", "3000"), "-port=3000 sets the port of the http server")

func main() {
	flag.Parse()

	m := mux.NewRouter()
	m.HandleFunc("/", httputil.HTTPFunc(home.Render))

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
