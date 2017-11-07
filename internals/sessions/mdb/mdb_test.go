package mdb_test

import (
	"os"
	"time"

	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/influx6/faux/tests"

	"github.com/influx6/faux/metrics"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/db/mongo"

	"github.com/influx6/faux/metrics/custom"

	"github.com/influx6/devapp/internals/sessions/mdb"
)

var (
	events = metrics.New(custom.StackDisplay(os.Stdout))

	config = mongo.Config{
		Mode:     mgo.Monotonic,
		DB:       os.Getenv("sessions_MONGO_DB"),
		Host:     os.Getenv("sessions_MONGO_HOST"),
		User:     os.Getenv("sessions_MONGO_USER"),
		AuthDB:   os.Getenv("sessions_MONGO_AUTHDB"),
		Password: os.Getenv("sessions_MONGO_PASSWORD"),
	}

	testCol = "session_test_collection"
)

// TestGetSession validates the retrieval of a Session
// record from a mongodb.
func TestGetSession(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(sessionCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Session record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Session record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Session into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Session into db.")

	_, err = api.Get(ctx, elem.PublicID)
	if err != nil {
		tests.Failed("Successfully retrieved stored record for Session from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved stored record for Session from db.")
}

// TestGetAllSession validates the retrieval of all Session
// record from a mongodb.
func TestGetAllSession(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(sessionCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Session record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Session record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Session into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Session into db.")

	records, _, err := api.GetAllPerPage(ctx, "asc", "public_id", -1, -1)
	if err != nil {
		tests.Failed("Successfully retrieved all records for Session from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved all records for Session from db.")

	if len(records) == 0 {
		tests.Failed("Successfully retrieved atleast 1 record for Session from db.")
	}
	tests.Passed("Successfully retrieved atleast 1 record for Session from db.")
}

// TestGetAllSessionOrderBy validates the retrieval of all Session
// record from a mongodb.
func TestGetAllSessionByOrder(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(sessionCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Session record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Session record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Session into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Session into db.")

	records, err := api.GetAllByOrder(ctx, "asc", "public_id")
	if err != nil {
		tests.Failed("Successfully retrieved all records for Session from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved all records for Session from db.")

	if len(records) == 0 {
		tests.Failed("Successfully retrieved atleast 1 record for Session from db.")
	}
	tests.Passed("Successfully retrieved atleast 1 record for Session from db.")
}

// TestSessionCreate validates the creation of a Session
// record with a mongodb.
func TestSessionCreate(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(sessionCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Session record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Session record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Session into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Session into db.")
}

// TestSessionUpdate validates the update of a Session
// record with a mongodb.
func TestSessionUpdate(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(sessionCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Session record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Session record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Session into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Session into db.")

	//TODO: Update something.

	if err := api.Update(ctx, elem.PublicID, elem); err != nil {
		tests.Failed("Successfully updated record for Session into db: %+q.", err)
	}
	tests.Passed("Successfully updated record for Session into db.")
}

// TestSessionDelete validates the removal of a Session
// record from a mongodb.
func TestSessionDelete(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(sessionCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for Session record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for Session record")

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for Session into db: %+q.", err)
	}
	tests.Passed("Successfully added record for Session into db.")

	if err := api.Delete(ctx, elem.PublicID); err != nil {
		tests.Failed("Successfully removed record for Session into db: %+q.", err)
	}
	tests.Passed("Successfully removed record for Session into db.")
}
