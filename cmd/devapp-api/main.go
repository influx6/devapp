package main

//go:generate go generate ./controllers/...
//go:generate go generate ./static/...

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	profilesapi "github.com/influx6/devapp/internals/profiles/handler"
	sessionsapi "github.com/influx6/devapp/internals/sessions/handler"
	sessionsdbapi "github.com/influx6/devapp/internals/sessions/mdb"
	tokensdbapi "github.com/influx6/devapp/internals/tokens/mdb"
	userapi "github.com/influx6/devapp/internals/users/handler"
	userdbapi "github.com/influx6/devapp/internals/users/mdb"
	"github.com/influx6/faux/db/mongo"
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

	var profiles profilesapi.ProfileAPI
	users := userapi.UserAPI{DB: usersdb, Tokens: tokensdb}
	sessions := sessionsapi.SessionAPI{DB: sessionsdb, UserDB: usersdb, TokensDB: tokensdb}

	m := mux.NewRouter()
	m.NotFoundHandler = httputil.HTTPFunc(httputil.NotFound)

	mw := httputil.MWi(httputil.MetricsMW(logs), httputil.LogMW)

	m.HandleFunc("/profile", httputil.HTTPFunc(mw(
		httputil.HTTPConditionErrorFunc(
			sessions.Authenticate,
			httputil.HTTPConditionErrorFunc(
				sessions.TwoFactorAuthorization,
				profiles.Get,
				httputil.BadRequestWithError,
			),
			httputil.BadRequestWithError,
		),
	)))

	m.HandleFunc("/users/new", httputil.HTTPFunc(mw(
		httputil.HTTPConditionErrorFunc(
			users.CreateUserFromURLEncoded,
			httputil.OKRequest,
			httputil.BadRequestWithError,
		),
	)))

	m.HandleFunc("/session/new", httputil.HTTPFunc(mw(
		httputil.HTTPConditionErrorFunc(
			sessions.Login,
			httputil.HTTPConditionErrorFunc(
				sessions.TwoFactorAuthorization,
				httputil.OKRequest,
				httputil.BadRequestWithError,
			),
			httputil.BadRequestWithError,
		),
	)))

	m.HandleFunc("/session/destroy", httputil.HTTPFunc(mw(
		httputil.HTTPConditionErrorFunc(
			sessions.Logout,
			httputil.OKRequest,
			httputil.BadRequestWithError,
		),
	)))

	m.HandleFunc("/users/twofactor/enable", httputil.HTTPFunc(mw(
		httputil.HTTPConditionErrorFunc(
			sessions.Authenticate,
			httputil.HTTPConditionErrorFunc(
				users.EnableTwoFactor,
				httputil.Then(users.UserTwoFactorQTURL, users.SetTwoFactorAsSeen),
				httputil.BadRequestWithError,
			),
			httputil.BadRequestWithError,
		),
	)))

	m.HandleFunc("/users/twofactor/qr", httputil.HTTPFunc(mw(
		httputil.HTTPConditionErrorFunc(
			sessions.Authenticate,
			httputil.Then(users.UserTwoFactorQRImage, users.SetTwoFactorAsSeen),
			httputil.BadRequestWithError,
		),
	)))

	m.HandleFunc("/users/twofactor/disable", httputil.HTTPFunc(mw(
		httputil.HTTPConditionErrorFunc(
			sessions.Authenticate,
			httputil.HTTPConditionErrorFunc(
				sessions.TwoFactorAuthorization,
				users.DisableTwoFactor,
				httputil.BadRequestWithError,
			),
			httputil.BadRequestWithError,
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
