package mongoapi_test

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

	"github.com/influx6/devapp/internals/users/mongoapi"
)

var (
	events = metrics.New(custom.StackDisplay(os.Stdout))

	config = mongo.Config{
		Mode:     mgo.Monotonic,
		DB:       os.Getenv("users_MONGO_DB"),
		Host:     os.Getenv("users_MONGO_HOST"),
		User:     os.Getenv("users_MONGO_USER"),
		AuthDB:   os.Getenv("users_MONGO_AUTHDB"),
		Password: os.Getenv("users_MONGO_PASSWORD"),
	}

	testCol = "user_test_collection"
)

// TestGetUser validates the retrieval of a User
// record from a mongodb.
func TestGetUser(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(userCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for User record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for User record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for User into db: %+q.", err)
	}
	tests.Passed("Successfully added record for User into db.")

	_, err = api.Get(ctx, elem.PublicID)
	if err != nil {
		tests.Failed("Successfully retrieved stored record for User from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved stored record for User from db.")
}

// TestGetAllUser validates the retrieval of all User
// record from a mongodb.
func TestGetAllUser(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(userCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for User record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for User record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for User into db: %+q.", err)
	}
	tests.Passed("Successfully added record for User into db.")

	records, _, err := api.GetAllPerPage(ctx, "asc", "public_id", -1, -1)
	if err != nil {
		tests.Failed("Successfully retrieved all records for User from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved all records for User from db.")

	if len(records) == 0 {
		tests.Failed("Successfully retrieved atleast 1 record for User from db.")
	}
	tests.Passed("Successfully retrieved atleast 1 record for User from db.")
}

// TestGetAllUserOrderBy validates the retrieval of all User
// record from a mongodb.
func TestGetAllUserByOrder(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(userCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for User record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for User record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for User into db: %+q.", err)
	}
	tests.Passed("Successfully added record for User into db.")

	records, err := api.GetAllByOrder(ctx, "asc", "public_id")
	if err != nil {
		tests.Failed("Successfully retrieved all records for User from db: %+q.", err)
	}
	tests.Passed("Successfully retrieved all records for User from db.")

	if len(records) == 0 {
		tests.Failed("Successfully retrieved atleast 1 record for User from db.")
	}
	tests.Passed("Successfully retrieved atleast 1 record for User from db.")
}

// TestUserCreate validates the creation of a User
// record with a mongodb.
func TestUserCreate(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(userCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for User record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for User record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for User into db: %+q.", err)
	}
	tests.Passed("Successfully added record for User into db.")
}

// TestUserUpdate validates the update of a User
// record with a mongodb.
func TestUserUpdate(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(userCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for User record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for User record")

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for User into db: %+q.", err)
	}
	tests.Passed("Successfully added record for User into db.")

	//TODO: Update something.

	if err := api.Update(ctx, elem.PublicID, elem); err != nil {
		tests.Failed("Successfully updated record for User into db: %+q.", err)
	}
	tests.Passed("Successfully updated record for User into db.")
}

// TestUserDelete validates the removal of a User
// record from a mongodb.
func TestUserDelete(t *testing.T) {
	api := mongoapi.New(testCol, events, mongo.New(config))

	ctx := context.WithTimeout(context.NewValueBag(), 10*time.Second)

	elem, err := loadJSONFor(userCreateJSON)
	if err != nil {
		tests.Failed("Successfully loaded JSON for User record: %+q.", err)
	}
	tests.Passed("Successfully loaded JSON for User record")

	if err := api.Create(ctx, elem); err != nil {
		tests.Failed("Successfully added record for User into db: %+q.", err)
	}
	tests.Passed("Successfully added record for User into db.")

	if err := api.Delete(ctx, elem.PublicID); err != nil {
		tests.Failed("Successfully removed record for User into db: %+q.", err)
	}
	tests.Passed("Successfully removed record for User into db.")
}
