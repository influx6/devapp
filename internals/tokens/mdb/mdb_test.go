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

	"github.com/influx6/devapp/internals/tokens/mdb"
)

var (
	events = metrics.New(custom.StackDisplay(os.Stdout))

	config = mongo.Config{
		Mode:     mgo.Monotonic,
		DB:       os.Getenv("tokens_MONGO_DB"),
		Host:     os.Getenv("tokens_MONGO_HOST"),
		User:     os.Getenv("tokens_MONGO_USER"),
		AuthDB:   os.Getenv("tokens_MONGO_AUTHDB"),
		Password: os.Getenv("tokens_MONGO_PASSWORD"),
	}

	testCol = "tokenrecord_test_collection"
)

// TestGetTokenRecord validates the retrieval of a TokenRecord
// record from a mongodb.
func TestGetTokenRecord(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(tokenrecordCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for TokenRecord record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for TokenRecord record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for TokenRecord into db: %+q.", err)
	}
	tests.Passed("Successfully added record for TokenRecord into db.")

	_, err = api.Get(ctx, elem.PublicID)
	if err != nil {
		tests.Failed("Successfully retrieved stored record for TokenRecord from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved stored record for TokenRecord from db.")
}

// TestGetAllTokenRecord validates the retrieval of all TokenRecord
// record from a mongodb.
func TestGetAllTokenRecord(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(tokenrecordCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for TokenRecord record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for TokenRecord record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for TokenRecord into db: %+q.", err)
	}
	tests.Passed("Successfully added record for TokenRecord into db.")

	records, _, err := api.GetAllPerPage(ctx, "asc", "public_id", -1, -1)
	if err != nil {
		tests.Failed("Successfully retrieved all records for TokenRecord from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved all records for TokenRecord from db.")

	if len(records) == 0 {
		tests.Failed("Successfully retrieved atleast 1 record for TokenRecord from db.")
	}
	tests.Passed("Successfully retrieved atleast 1 record for TokenRecord from db.")
}

// TestGetAllTokenRecordOrderBy validates the retrieval of all TokenRecord
// record from a mongodb.
func TestGetAllTokenRecordByOrder(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(tokenrecordCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for TokenRecord record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for TokenRecord record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for TokenRecord into db: %+q.", err)
	}
	tests.Passed("Successfully added record for TokenRecord into db.")

	records, err := api.GetAllByOrder(ctx, "asc", "public_id")
	if err != nil {
		tests.Failed("Successfully retrieved all records for TokenRecord from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved all records for TokenRecord from db.")

	if len(records) == 0 {
		tests.Failed("Successfully retrieved atleast 1 record for TokenRecord from db.")
	}
	tests.Passed("Successfully retrieved atleast 1 record for TokenRecord from db.")
}

// TestTokenRecordCreate validates the creation of a TokenRecord
// record with a mongodb.
func TestTokenRecordCreate(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(tokenrecordCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for TokenRecord record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for TokenRecord record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for TokenRecord into db: %+q.", err)
	}
	tests.Passed("Successfully added record for TokenRecord into db.")
}

// TestTokenRecordUpdate validates the update of a TokenRecord
// record with a mongodb.
func TestTokenRecordUpdate(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(tokenrecordCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for TokenRecord record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for TokenRecord record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for TokenRecord into db: %+q.", err)
	}
	tests.Passed("Successfully added record for TokenRecord into db.")

	//TODO: Update something.

	if err := api.Update(ctx, elem.PublicID, elem); err != nil {
		tests.Failed("Successfully updated record for TokenRecord into db: %+q.", err)
	}
	tests.Passed("Successfully updated record for TokenRecord into db.")
}

// TestTokenRecordDelete validates the removal of a TokenRecord
// record from a mongodb.
func TestTokenRecordDelete(t *testing.T) {
	api := mdb.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(tokenrecordCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for TokenRecord record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for TokenRecord record")

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for TokenRecord into db: %+q.", err)
	}
	tests.Passed("Successfully added record for TokenRecord into db.")

	if err := api.Delete(ctx, elem.PublicID); err != nil {
		tests.Failed("Successfully removed record for TokenRecord into db: %+q.", err)
	}
	tests.Passed("Successfully removed record for TokenRecord into db.")
}
