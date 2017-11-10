package db

import (
	"errors"

	"github.com/influx6/devapp/internals/users"
	"github.com/influx6/devapp/internals/users/mdb"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Delete handles receiving requests to delete a user from the database.
func Delete(ctx context.Context, log metrics.Metrics, db *mdb.UserDB, id string) error {
	log.Emit(metrics.Info("Get Existing User").With("user_id", id))

	if err := db.Delete(ctx, id); err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{"public_id": id}))
		return err
	}

	return nil
}

// GetByUsername handles receiving requests to retrieve a user from the database by it's username.
func GetByUsername(ctx context.Context, log metrics.Metrics, db *mdb.UserDB, username string) (users.User, error) {
	log.Emit(metrics.Info("Get Existing User By Username").With("username", username))

	var nu users.User
	err := db.Exec(ctx, func(col *mgo.Collection) error {
		return col.Find(bson.M{"username": username}).One(&nu)
	})

	if err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{"username": username}))
		return nu, err
	}

	return nu, nil
}

// Get handles receiving requests to retrieve a user from the database.
func Get(ctx context.Context, log metrics.Metrics, db *mdb.UserDB, id string) (users.User, error) {
	log.Emit(metrics.Info("Get Existing User").With("user_id", id))

	nu, err := db.Get(ctx, id)
	if err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{"public_id": id}))
		return nu, err
	}

	return nu, nil
}

// UserRecords defines a struct which returns the total fields and page details
// used in retrieving the records.
type UserRecords struct {
	Total           int          `json:"total"`
	Page            int          `json:"page"`
	ResponsePerPage int          `json:"responsePerPage"`
	Records         []users.User `json:"records"`
}

// GetAll handles receiving requests to retrieve all user from the database.
func GetAll(ctx context.Context, log metrics.Metrics, db *mdb.UserDB, page, responsePerPage int) (UserRecords, error) {
	log.Emit(metrics.Info("Get Existing User").WithFields(metrics.Field{
		"page":            page,
		"responsePerPage": responsePerPage,
	}))

	records, realTotalRecords, err := db.GetAllPerPage(ctx, "asc", "public_id", page, responsePerPage)
	if err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{
			"page":            page,
			"responsePerPage": responsePerPage,
		}))

		return UserRecords{}, err
	}

	return UserRecords{
		Page:            page,
		Records:         records,
		Total:           realTotalRecords,
		ResponsePerPage: responsePerPage,
	}, nil
}

// Create handles receiving requests to create a user from the server.
func Create(ctx context.Context, log metrics.Metrics, db *mdb.UserDB, nw users.NewUser) (users.User, error) {
	log.Emit(metrics.Info("Create New User"))

	newUser, err := users.New(nw)
	if err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{"username": nw.Username}))
		return users.User{}, err
	}

	if err := db.Create(ctx, newUser); err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{"username": nw.Username}))
		return users.User{}, err
	}

	return newUser, nil
}

// UpdatePassword handles receiving requests to update a user identified by it's public_id.
func UpdatePassword(ctx context.Context, log metrics.Metrics, db *mdb.UserDB, nw users.UpdateUserPassword) error {
	log.Emit(metrics.Info("Update User Password").With("user", nw.PublicID))

	if nw.PublicID == "" {
		err := errors.New("JSON UpdateUserPassword.PublicID is empty")
		log.Emit(metrics.Error(err).WithFields(metrics.Field{
			"user_id": nw.PublicID,
		}))
		return err
	}

	if nw.Password == "" {
		err := errors.New("JSON UpdateUserPassword.Password is empty")
		log.Emit(metrics.Error(err).WithFields(metrics.Field{
			"user_id": nw.PublicID,
		}))
		return err
	}

	dbuser, err := db.Get(ctx, nw.PublicID)
	if err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{
			"user_id": nw.PublicID,
		}))

		return err
	}

	if err := dbuser.ChangePassword(nw.Password); err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{
			"user_id": nw.PublicID,
		}))
		return err
	}

	if err := db.Update(ctx, dbuser.PublicID, dbuser); err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{
			"user_id": nw.PublicID,
		}))
		return err
	}

	return nil
}

// Update handles receiving requests to update a user identified by it's public_id.
func Update(ctx context.Context, log metrics.Metrics, db *mdb.UserDB, nw users.User) error {
	log.Emit(metrics.Info("Update User").With("user", nw.PublicID))

	if nw.PublicID == "" {
		err := errors.New("JSON User.PublicID is empty")
		log.Emit(metrics.Error(err).WithFields(metrics.Field{
			"user_id":  nw.PublicID,
			"username": nw.Username,
		}))

		return err
	}

	if err := db.Update(ctx, nw.PublicID, nw); err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{
			"user_id":  nw.PublicID,
			"username": nw.Username,
		}))

		return err
	}

	return nil
}
