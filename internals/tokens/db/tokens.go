package db

import (
	"github.com/influx6/devapp/internals/tokens"
	"github.com/influx6/devapp/internals/tokens/mdb"
	"github.com/influx6/devapp/internals/users"
	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	uuid "github.com/satori/go.uuid"
)

// AddToken adds giving token into user token records.
func AddToken(ctx context.Context, log metrics.Metrics, db *mdb.TokenRecordDB, nu users.User, token string) error {
	log.Emit(metrics.Info("Add token to User TokenRecord").WithFields(metrics.Field{
		"username": nu.Username,
		"user_id":  nu.PublicID,
		"token":    token,
	}))

	// If we do not exists, create it.
	if err := db.Exec(ctx, func(col *mgo.Collection) error {
		query := bson.M{"user_id": nu.PublicID}
		log.Emit(metrics.Info("Checking existence of user token record").WithFields(metrics.Field{"username": nu.Username, "token": token, "query": query}))
		count, cerr := col.Find(query).Count()
		if cerr != nil {
			return cerr
		}

		if count == 0 {
			return mgo.ErrNotFound
		}

		return nil
	}); err != nil {
		if _, cerr := Create(ctx, log, db, nu); cerr != nil {
			return cerr
		}
	}

	if err := db.Exec(ctx, func(col *mgo.Collection) error {
		query, data := bson.M{"user_id": nu.PublicID}, bson.M{"$addToSet": bson.M{"tokens": token}}
		log.Emit(metrics.Info("Making mongodb operation").WithFields(metrics.Field{"username": nu.Username, "token": token, "query": query, "data": data}))
		return col.Update(query, data)
	}); err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{"username": nu.Username, "token": "token"}))
		return err
	}

	log.Emit(metrics.Info("Token Added to user TokenRecord").WithFields(metrics.Field{"username": nu.Username, "token": token}))
	return nil
}

// UsedToken returns true/false if giving token is already used by user in it's token records.
func UsedToken(ctx context.Context, log metrics.Metrics, db *mdb.TokenRecordDB, nu users.User, token string) (bool, error) {
	log.Emit(metrics.Info("Checked used token in User TokenRecord").WithFields(metrics.Field{
		"username": nu.Username,
		"user_id":  nu.PublicID,
		"token":    token,
	}))

	// If we do not exists, create it.
	if err := db.Exec(ctx, func(col *mgo.Collection) error {
		query := bson.M{"user_id": nu.PublicID}
		log.Emit(metrics.Info("Checking existence of user token record").WithFields(metrics.Field{"username": nu.Username, "token": token, "query": query}))
		count, cerr := col.Find(query).Count()
		if cerr != nil {
			return cerr
		}

		if count == 0 {
			return mgo.ErrNotFound
		}

		return nil
	}); err != nil {
		if _, cerr := Create(ctx, log, db, nu); cerr != nil {
			return false, cerr
		}
	}

	var used bool
	if err := db.Exec(ctx, func(col *mgo.Collection) error {
		query := bson.M{"user_id": nu.PublicID, "tokens": []string{token}}
		log.Emit(metrics.Info("Making mongodb operation").WithFields(metrics.Field{"username": nu.Username, "token": token, "query": query}))

		count, err := col.Find(query).Count()
		if err != nil {
			return err
		}

		if count > 0 {
			used = true
		}
		return nil
	}); err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{"username": nu.Username, "token": "token"}))
		return false, err
	}

	if used {
		log.Emit(metrics.Info("Token already used by user").WithFields(metrics.Field{"username": nu.Username, "token": token}))
	} else {
		log.Emit(metrics.Info("Token not yet used by user").WithFields(metrics.Field{"username": nu.Username, "token": token}))
	}

	return used, nil
}

// Create adds a new session for the specified users.
func Create(ctx context.Context, log metrics.Metrics, db *mdb.TokenRecordDB, nu users.User) (tokens.TokenRecord, error) {
	log.Emit(metrics.Info("Create User TokenRecord").WithFields(metrics.Field{
		"username": nu.Username,
		"user_id":  nu.PublicID,
	}))

	var record tokens.TokenRecord
	record.PublicID = uuid.NewV4().String()
	record.UserID = nu.PublicID

	if err := db.Create(ctx, record); err != nil {
		log.Emit(metrics.Error(err).WithFields(metrics.Field{"username": nu.Username}))
		return record, err
	}

	log.Emit(metrics.Info("User TokenRecord ready").WithFields(metrics.Field{"username": nu.Username}))

	return record, nil
}
