package db

import (
	"time"

	"github.com/influx6/devapp/internals/sessions"
	"github.com/influx6/devapp/internals/sessions/mdb"
	"github.com/influx6/devapp/internals/users"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Create adds a new session for the specified users.
func Create(ctx context.Context, log metrics.Metrics, db *mdb.SessionDB, expiration time.Duration, nu users.User) (*sessions.Session, error) {
	log.Emit(metrics.Info("Create New Session").WithFields(metrics.Field{
		"username": nu.Username,
		"user_id":  nu.PublicID,
	}))

	currentTime := time.Now()

	// Attempt to retrieve session from db if we still have an outstanding non-expired sessions.
	newSession, err := db.Get(ctx, nu.PublicID)
	if err == nil {
		log.Emit(metrics.Info("Found user session").With("session", newSession).With("expired", currentTime.After(newSession.Expires)).
			WithFields(metrics.Field{"username": nu.Username, "user_id": nu.PublicID}))

		// We have an existing session and the time of expiring is still counting, simly return
		if !newSession.Expires.IsZero() && currentTime.Before(newSession.Expires) {
			return &newSession, nil
		}

		// 	If we still have active session, then simply kick it and return safely.
		if newSession.Expires.IsZero() || currentTime.After(newSession.Expires) {

			// Delete this sessions
			if err := db.Delete(ctx, newSession.UserID); err != nil {
				log.Emit(metrics.Error(err).WithMessage("Failed to delete old session").
					WithFields(metrics.Field{"username": nu.Username, "user_id": nu.PublicID}))
				return nil, err
			}
		}
	}

	// Create new session and store session into db.
	newSession = sessions.New(nu.PublicID, time.Now().Add(expiration))
	if err := db.Create(ctx, newSession); err != nil {
		log.Emit(metrics.Error(err).
			WithMessage("Failed to save new session").
			WithFields(metrics.Field{"username": nu.Username, "user_id": nu.PublicID}))
		return nil, err
	}

	return &newSession, nil
}

// Get retrieves the session associated with the giving User.
func Get(ctx context.Context, log metrics.Metrics, db *mdb.SessionDB, userID string) (sessions.Session, error) {
	log.Emit(metrics.Info("Get Existing Session").WithFields(metrics.Field{
		"user_id": userID,
	}))

	var existingSession sessions.Session

	var res map[string]interface{}
	query := bson.M{"user_id": userID}
	// Attempt to retrieve session from db if we still have an outstanding non-expired sessions.
	if err := db.Exec(ctx, func(col *mgo.Collection) error {
		return col.Find(query).One(&res)
	}); err != nil {
		log.Emit(metrics.Error(err).
			WithMessage("Failed to retrieve session from db: %+q").
			WithFields(metrics.Field{"user_id": userID, "query": query}))
		return existingSession, err
	}

	if err := existingSession.Consume(res); err != nil {
		log.Emit(metrics.Error(err).
			WithMessage("Failed to unmarshall session data to Session Object: %+q").
			WithFields(metrics.Field{"user_id": userID, "query": query}))
		return existingSession, err
	}

	return existingSession, nil
}

// Delete removes an existing session from the db for a specified users.
func Delete(ctx context.Context, log metrics.Metrics, db *mdb.SessionDB, userID string) error {
	log.Emit(metrics.Info("Delete Existing Session").WithFields(metrics.Field{
		"user_id": userID,
	}))

	// Delete this sessions
	if err := db.Exec(ctx, func(col *mgo.Collection) error {
		return col.Remove(bson.M{"user_id": userID})
	}); err != nil {
		log.Emit(metrics.Error(err).WithMessage("Failed to remove session from db").
			WithFields(metrics.Field{"user_id": userID}))
		return err
	}

	log.Emit(metrics.Info("Removed user session from db").
		WithFields(metrics.Field{"user_id": userID}))

	return nil
}
