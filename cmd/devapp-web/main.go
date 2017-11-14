package main

//go:generate go generate ./controllers/...
//go:generate go generate ./static/...

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/influx6/devapp/controllers/home"
	"github.com/influx6/devapp/controllers/profile"
	"github.com/influx6/devapp/controllers/signup"
	"github.com/influx6/devapp/controllers/twoauth"
	sessionsapi "github.com/influx6/devapp/internals/sessions/handler"
	sessionsdbapi "github.com/influx6/devapp/internals/sessions/mdb"
	tokensdbapi "github.com/influx6/devapp/internals/tokens/mdb"
	userapi "github.com/influx6/devapp/internals/users/handler"
	userdbapi "github.com/influx6/devapp/internals/users/mdb"
	"github.com/influx6/devapp/static"
	"github.com/influx6/faux/db/mongo"
	"github.com/influx6/faux/filesystem"
	"github.com/influx6/faux/filesystem/bytereaders"
	"github.com/influx6/faux/httputil"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/custom"
	mgo "gopkg.in/mgo.v2"
)

var (
	logs = metrics.New(custom.StackDisplay(os.Stdout))
	port = flag.String("port", envOrDefault("PORT", "3000"), "-port=3000 sets the port of the http server")

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
		Key:        []string{"public_id"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		DropDups:   true,
	}, mgo.Index{
		Key:        []string{"user_id"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		DropDups:   true,
	})

	tokensdb := tokensdbapi.New("user_tokens_collections", logs, mdb, mgo.Index{
		Key:        []string{"public_id"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		DropDups:   true,
	}, mgo.Index{
		Key:        []string{"user_id"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		DropDups:   true,
	})

	usersdb := userdbapi.New("user_collection", logs, mdb, mgo.Index{
		Key:        []string{"public_id"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		DropDups:   true,
	}, mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		Background: true,
		Sparse:     true,
		DropDups:   true,
	})

	staticFS := filesystem.NewSystemGroup()

	staticFS.MustRegister("/static/", filesystem.New(
		filesystem.StripPrefix("/", bytereaders.FileFromByteReader(static.FindFileReader))),
	).MustRegister("/static-controllers/home", filesystem.New(
		filesystem.StripPrefix("/", bytereaders.FileFromByteReader(home.FindFileReader))),
	).MustRegister("/static-controllers/profile", filesystem.New(
		filesystem.StripPrefix("/", bytereaders.FileFromByteReader(profile.FindFileReader))),
	).MustRegister("/static-controllers/twoauth", filesystem.New(
		filesystem.StripPrefix("/", bytereaders.FileFromByteReader(twoauth.FindFileReader))),
	)

	users := userapi.UserAPI{DB: usersdb, Tokens: tokensdb}
	sessions := sessionsapi.SessionAPI{DB: sessionsdb, UserDB: usersdb, TokensDB: tokensdb}

	m := mux.NewRouter()
	m.NotFoundHandler = httputil.HTTPFunc(httputil.NotFound)

	mw := httputil.MWi(httputil.MetricsMW(logs), httputil.LogMW)

	// static files
	m.PathPrefix("/static/").Handler(httputil.GzipServer(staticFS, true, mw))
	m.PathPrefix("/static-controllers/").Handler(httputil.GzipServer(staticFS, true, mw))

	// view routes
	m.HandleFunc("/signup", httputil.HTTPFunc(mw(signup.Render)))
	m.HandleFunc("/", httputil.HTTPFunc(mw(
		httputil.HTTPConditionFunc(
			sessions.Authenticate,
			httputil.HTTPRedirect("/profile", http.StatusTemporaryRedirect),
			home.Render,
		),
	)))

	m.HandleFunc("/profile", httputil.HTTPFunc(mw(
		httputil.HTTPConditionFunc(
			sessions.Login,
			httputil.HTTPConditionFunc(
				sessions.TwoFactorAuthorizationCheck,
				profile.Render,
				httputil.HTTPRedirect("/session/twofactor", http.StatusTemporaryRedirect),
			),
			httputil.HTTPRedirect("/", http.StatusTemporaryRedirect),
		),
	)))

	m.HandleFunc("/users/new", httputil.HTTPFunc(mw(
		httputil.HTTPConditionFunc(
			users.CreateUserFromURLEncoded,
			httputil.HTTPRedirect("/", http.StatusTemporaryRedirect),
			signup.Render,
		),
	)))

	// login-logout
	m.HandleFunc("/session/twofactor", httputil.HTTPFunc(mw(
		httputil.HTTPConditionFunc(
			sessions.Authenticate,
			httputil.HTTPConditionFunc(
				sessions.TwoFactorAuthorizationCheck,
				httputil.HTTPRedirect("/profile", http.StatusTemporaryRedirect),
				twoauth.Render,
			),
			httputil.HTTPRedirect("/", http.StatusTemporaryRedirect),
		),
	)))

	m.HandleFunc("/session/twofactor/finish", httputil.HTTPFunc(mw(
		httputil.HTTPConditionFunc(
			sessions.Authenticate,
			httputil.HTTPConditionFunc(
				sessions.TwoFactorAuthorization,
				httputil.HTTPRedirect("/profile", http.StatusTemporaryRedirect),
				twoauth.Render,
			),
			httputil.HTTPRedirect("/", http.StatusTemporaryRedirect),
		),
	)))

	m.HandleFunc("/session/new", httputil.HTTPFunc(mw(
		httputil.HTTPConditionFunc(
			sessions.Login,
			httputil.HTTPConditionFunc(
				sessions.TwoFactorAuthorizationCheck,
				httputil.HTTPRedirect("/profile", http.StatusTemporaryRedirect),
				httputil.HTTPRedirect("/session/twofactor", http.StatusTemporaryRedirect),
			),
			httputil.HTTPRedirect("/", http.StatusTemporaryRedirect),
		),
	)))

	m.HandleFunc("/session/destroy", httputil.HTTPFunc(mw(
		httputil.HTTPConditionFunc(
			sessions.Logout,
			httputil.HTTPRedirect("/", http.StatusTemporaryRedirect),
			httputil.HTTPRedirect("/", http.StatusTemporaryRedirect),
		),
	)))

	m.HandleFunc("/users/twofactor", httputil.HTTPFunc(mw(
		httputil.HTTPConditionFunc(
			sessions.Login,
			httputil.HTTPConditionFunc(
				twoauth.RenderQR,
				users.SetTwoFactorAsSeen,
				httputil.HTTPRedirect("/profile", http.StatusTemporaryRedirect),
			),
			httputil.HTTPRedirect("/", http.StatusTemporaryRedirect),
		),
	)))

	m.HandleFunc("/users/twofactor/enable", httputil.HTTPFunc(mw(
		httputil.HTTPConditionFunc(
			sessions.Authenticate,
			httputil.HTTPConditionFunc(
				users.EnableTwoFactor,
				httputil.HTTPRedirect("/users/twofactor", http.StatusTemporaryRedirect),
				httputil.HTTPRedirect("/", http.StatusTemporaryRedirect),
			),
			httputil.BadRequest,
		),
	)))

	m.HandleFunc("/users/twofactor/disable", httputil.HTTPFunc(mw(
		httputil.HTTPConditionFunc(
			sessions.Authenticate,
			httputil.HTTPConditionFunc(
				users.DisableTwoFactor,
				httputil.HTTPRedirect("/profile", http.StatusTemporaryRedirect),
				httputil.HTTPRedirect("/", http.StatusTemporaryRedirect),
			),
			httputil.BadRequest,
		),
	)))

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
