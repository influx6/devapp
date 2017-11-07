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
	"github.com/influx6/devapp/controllers/profile"
	"github.com/influx6/devapp/controllers/signin"
	"github.com/influx6/devapp/controllers/signup"
	sessionsapi "github.com/influx6/devapp/internals/sessions/handler"
	sessionsdbapi "github.com/influx6/devapp/internals/sessions/mdb"
	userapi "github.com/influx6/devapp/internals/users/handler"
	userdbapi "github.com/influx6/devapp/internals/users/mdb"
	"github.com/influx6/devapp/static"
	"github.com/influx6/faux/db/mongo"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/custom"
	mgo "gopkg.in/mgo.v2"
)

var (
	logs = metrics.New(custom.StackDisplay(os.Stdout))
	port = flag.String("port", envOrDefault("PORT", "3000"), "-port=3000 sets the port of the http server")

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

	dbconf = mongo.Config{
		Mode:     mgo.Monotonic,
		DB:       os.Getenv("DEVAPP_MONGO_DB"),
		User:     os.Getenv("DEVAPP_MONGO_USER"),
		Host:     os.Getenv("DEVAPP_MONGO_HOST"),
		AuthDB:   os.Getenv("DEVAPP_MONGO_AUTHDB"),
		Password: os.Getenv("DEVAPP_MONGO_PASSWORD"),
	}
)

func main() {
	flag.Parse()

	logs.Emit(metrics.Info("Using db config").With("config", dbconf))

	mdb := mongo.New(dbconf)
	sessionsdb := sessionsdbapi.New("session_collection", logs, mdb, mgo.Index{
		Key:        []string{"public_id", "user_id"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		DropDups:   true,
	})

	usersdb := userdbapi.New("user_collection", logs, mdb, mgo.Index{
		Key:        []string{"public_id", "username"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		DropDups:   true,
	})

	users := userapi.UserAPI{DB: usersdb}
	sessions := sessionsapi.SessionAPI{DB: sessionsdb, UserDB: usersdb}

	mw := httputil.MWi(httputil.MetricsMW(logs), httputil.LogMW)
	authmw := httputil.MWi(mw, signin.AuthenticateMW(sessions))
	prefixmw := httputil.MWi(mw, httputil.StripPrefixMW("/static/"))
	confirmLoginMW := httputil.MWi(mw, signin.IfLoggedHandler(sessions, "/profile"))

	m := mux.NewRouter()
	m.NotFoundHandler = httputil.HTTPFunc(httputil.NotFound)

	// static files
	m.PathPrefix("/static/").Handler(httputil.GzipServer(staticFS, true, prefixmw))

	// view routes
	m.HandleFunc("/", httputil.HTTPFunc(confirmLoginMW(home.Render)))
	m.HandleFunc("/signup", httputil.HTTPFunc(mw(signup.Render)))
	m.HandleFunc("/profile", httputil.HTTPFunc(authmw(profile.Render)))

	m.HandleFunc("/users/new", httputil.HTTPFunc(mw(signup.SignupHandler(users, "/"))))

	// login-logout
	m.HandleFunc("/session/new", httputil.HTTPFunc(confirmLoginMW(signin.LoginHandler(sessions, "/profile", "/"))))
	m.HandleFunc("/session/destroy", httputil.HTTPFunc(mw(signin.LogoutHandler(sessions, "/profile"))))

	// api routes
	m.HandleFunc("/api/users/new", httputil.HTTPFunc(mw(users.CreateUserFromURLEncoded)))
	m.HandleFunc("/api/sessions/login", httputil.HTTPFunc(mw(sessions.Login)))
	m.HandleFunc("/api/sessions/logout", httputil.HTTPFunc(mw(sessions.Logout)))

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
