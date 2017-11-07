// Package mdb provides a auto-generated package which contains a mongo CRUD API for the specific Session struct in package sessions.
//
//
package mdb

import (
	"fmt"
	"strings"

	mgo "gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"

	"github.com/influx6/faux/context"

	"github.com/influx6/faux/metrics"

	"github.com/influx6/devapp/internals/sessions"
)

// SessionFields defines an interface which exposes method to return a map of all
// attributes associated with the defined structure as decided by the structure.
type SessionFields interface {
	Fields() map[string]interface{}
}

// SessionBSON defines an interface which exposes method to return a bson.M type
// which contains all related fields for the giving  object.
type SessionBSON interface {
	BSON() bson.M
}

// SessionBSONConsumer defines an interface which accepts a map of data which will be consumed
// into the giving implementing structure as decided by the structure.
type SessionBSONConsumer interface {
	BSONConsume(bson.M) error
}

// SessionConsumer defines an interface which accepts a map of data which will be consumed
// into the giving implementing structure as decided by the structure.
type SessionConsumer interface {
	Consume(map[string]interface{}) error
}

// Mongod defines a interface which exposes a method for retrieving a
// mongo.Database and mongo.Session.
type Mongod interface {
	New() (*mgo.Database, *mgo.Session, error)
}

// SessionDB defines a structure which provide DB CRUD operations
// using mongo as the underline db.
type SessionDB struct {
	col             string
	db              Mongod
	metrics         metrics.Metrics
	ensuredIndex    bool
	incompleteIndex bool
	indexes         []mgo.Index
}

// New returns a new instance of SessionDB.
func New(col string, m metrics.Metrics, mo Mongod, indexes ...mgo.Index) *SessionDB {
	return &SessionDB{
		db:      mo,
		col:     col,
		metrics: m,
		indexes: indexes,
	}
}

// ensureIndex attempts to ensure all provided indexes into the specific collection.
func (mdb *SessionDB) ensureIndex() error {
	m := metrics.NewTrace("SessionDB.ensureIndex")
	defer mdb.metrics.Emit(metrics.Info("SessionDB.ensureIndex").WithTrace(m.End()))

	if mdb.ensuredIndex {
		return nil
	}

	if len(mdb.indexes) == 0 {
		return nil
	}

	// If we had an error before index was complete, then skip, we cant not
	// stop all ops because of failed index.
	if !mdb.ensuredIndex && mdb.incompleteIndex {
		return nil
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session for index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	collection := database.C(mdb.col)

	for _, index := range mdb.indexes {
		if err := collection.EnsureIndex(index); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to ensure session index").
				With("collection", mdb.col).
				With("index", index).
				With("error", err.Error()))

			mdb.incompleteIndex = true
			return err
		}

		mdb.metrics.Emit(metrics.Info("Succeeded in ensuring collection index").
			With("collection", mdb.col).
			With("index", index))
	}

	mdb.ensuredIndex = true

	mdb.metrics.Emit(metrics.Info("Finished adding index").
		With("collection", mdb.col))

	return nil
}

// Count attempts to return the total number of record from the db.
func (mdb *SessionDB) Count(ctx context.Context) (int, error) {
	m := metrics.NewTrace("SessionDB.Count")
	defer mdb.metrics.Emit(metrics.Info("SessionDB.Count").WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")

		mdb.metrics.Emit(metrics.Errorf("Failed to get record count").
			With("collection", mdb.col).
			With("error", err.Error()))
		return -1, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))

		return -1, err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to get record count").
			With("collection", mdb.col).
			With("error", err.Error()))

		return -1, err
	}

	defer session.Close()

	query := bson.M{}

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to get record count").
			With("collection", mdb.col).
			With("error", err.Error()))

		return -1, err
	}

	total, err := database.C(mdb.col).Find(query).Count()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to get record count").
			With("collection", mdb.col).
			With("query", query).
			With("error", err.Error()))

		return -1, err
	}

	mdb.metrics.Emit(metrics.Info("Deleted record").
		With("collection", mdb.col).
		With("query", query))

	return total, err
}

// Delete attempts to remove the record from the db using the provided publicID.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given sessions.Session struct.
func (mdb *SessionDB) Delete(ctx context.Context, publicID string) error {
	m := metrics.NewTrace("SessionDB.Delete")
	defer mdb.metrics.Emit(metrics.Info("SessionDB.Delete").With("publicID", publicID).WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to delete record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to delete record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	query := bson.M{
		"publicID": publicID,
	}

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to delete record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := database.C(mdb.col).Remove(query); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to delete record").
			With("collection", mdb.col).
			With("query", query).
			With("publicID", publicID).With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(metrics.Info("Deleted record").
		With("collection", mdb.col).
		With("query", query).
		With("publicID", publicID))

	return nil
}

// Create attempts to add the record into the db using the provided instance of the
// sessions.Session.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Session struct.
func (mdb *SessionDB) Create(ctx context.Context, elem sessions.Session) error {
	m := metrics.NewTrace("SessionDB.Create")
	defer mdb.metrics.Emit(metrics.Info("SessionDB.Create").With("publicID", elem.PublicID).WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to create record").With("publicID", elem.PublicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session").
			With("publicID", elem.PublicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to create record").With("publicID", elem.PublicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if fields, ok := interface{}(elem).(SessionBSON); ok {
		if err := database.C(mdb.col).Insert(fields.BSON()); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to create Session record").
				With("collection", mdb.col).
				With("elem", elem).
				With("error", err.Error()))
			return err
		}

		mdb.metrics.Emit(metrics.Info("Create record").
			With("collection", mdb.col).
			With("elem", elem))

		return nil
	}

	if fields, ok := interface{}(elem).(SessionFields); ok {
		if err := database.C(mdb.col).Insert(bson.M(fields.Fields())); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to create Session record").
				With("collection", mdb.col).
				With("elem", elem).
				With("error", err.Error()))
			return err
		}

		mdb.metrics.Emit(metrics.Info("Create record").
			With("collection", mdb.col).
			With("elem", elem))

		return nil
	}

	query := bson.M(map[string]interface{}{

		"public_id": elem.PublicID,

		"token": elem.Token,

		"user_id": elem.UserID,
	})

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to create record").With("publicID", elem.PublicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := database.C(mdb.col).Insert(query); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create Session record").
			With("collection", mdb.col).
			With("query", query).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(metrics.Info("Create record").
		With("collection", mdb.col).
		With("query", query))

	return nil
}

// GetAllPerPage retrieves all records from the db and returns a slice of sessions.Session type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Session struct.
func (mdb *SessionDB) GetAllPerPage(ctx context.Context, order string, orderBy string, page int, responsePerPage int) ([]sessions.Session, int, error) {
	m := metrics.NewTrace("SessionDB.GetAll")
	defer mdb.metrics.Emit(metrics.Info("SessionDB.GetAll").WithTrace(m.End()))

	switch strings.ToLower(order) {
	case "dsc", "desc":
		orderBy = "-" + orderBy
	}

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, -1, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, -1, err
	}

	if page <= 0 && responsePerPage <= 0 {
		records, err := mdb.GetAllByOrder(ctx, order, orderBy)
		return records, len(records), err
	}

	// Get total number of records.
	totalRecords, err := mdb.Count(ctx)
	if err != nil {
		return nil, -1, err
	}

	var totalWanted, indexToStart int

	if page <= 1 && responsePerPage > 0 {
		totalWanted = responsePerPage
		indexToStart = 0
	} else {
		totalWanted = responsePerPage * page
		indexToStart = totalWanted / 2

		if page > 1 {
			indexToStart++
		}
	}

	mdb.metrics.Emit(metrics.Info("DB:Query:GetAllPerPage").WithFields(metrics.Field{
		"starting_index":       indexToStart,
		"total_records_wanted": totalWanted,
		"order":                order,
		"orderBy":              orderBy,
		"page":                 page,
		"responsePerPage":      responsePerPage,
	}))

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, -1, err
	}

	defer session.Close()

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, -1, err
	}

	query := bson.M{}

	var items []sessions.Session

	if err := database.C(mdb.col).Find(query).Skip(indexToStart).Limit(totalWanted).Sort(orderBy).All(&items); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve all records of Session type from db").
			With("collection", mdb.col).
			With("query", query).
			With("error", err.Error()))

		return nil, -1, err
	}

	return items, totalRecords, nil

}

// GetAllByOrder retrieves all records from the db and returns a slice of sessions.Session type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Session struct.
func (mdb *SessionDB) GetAllByOrder(ctx context.Context, order, orderBy string) ([]sessions.Session, error) {
	m := metrics.NewTrace("SessionDB.GetAllByOrder")
	defer mdb.metrics.Emit(metrics.Info("SessionDB.GetAllByOrder").WithTrace(m.End()))

	switch strings.ToLower(order) {
	case "dsc", "desc":
		orderBy = "-" + orderBy
	}

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")

		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	defer session.Close()

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")

		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").
			With("collection", mdb.col).
			With("error", err.Error()))
		return nil, err
	}

	query := bson.M{}

	var items []sessions.Session

	if err := database.C(mdb.col).Find(query).Sort(orderBy).All(&items); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve all records of Session type from db").
			With("collection", mdb.col).
			With("query", query).
			With("error", err.Error()))

		return nil, err
	}

	return items, nil

}

// GetByField retrieves a record from the db using the provided field key and value
// returns the sessions.Session type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Session struct.
func (mdb *SessionDB) GetByField(ctx context.Context, key string, value interface{}) (sessions.Session, error) {
	m := metrics.NewTrace("SessionDB.GetByField")
	defer mdb.metrics.Emit(metrics.Info("SessionDB.GetByField").With(key, value).WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").With(key, value).
			With("collection", mdb.col).
			With("error", err.Error()))

		return sessions.Session{}, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))

		return sessions.Session{}, err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session").
			With(key, value).
			With("collection", mdb.col).
			With("error", err.Error()))

		return sessions.Session{}, err
	}

	defer session.Close()

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").With(key, value).
			With("collection", mdb.col).
			With("error", err.Error()))

		return sessions.Session{}, err
	}

	query := bson.M{key: value}

	var item sessions.Session

	if err := database.C(mdb.col).Find(query).One(&item); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve all records of Session type from db").
			With("query", query).
			With("collection", mdb.col).
			With("error", err.Error()))

		return sessions.Session{}, err
	}

	return item, nil

}

// Get retrieves a record from the db using the publicID and returns the sessions.Session type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Session struct.
func (mdb *SessionDB) Get(ctx context.Context, publicID string) (sessions.Session, error) {
	m := metrics.NewTrace("SessionDB.Get")
	defer mdb.metrics.Emit(metrics.Info("SessionDB.Get").With("publicID", publicID).WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return sessions.Session{}, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return sessions.Session{}, err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session").
			With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return sessions.Session{}, err
	}

	defer session.Close()

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return sessions.Session{}, err
	}

	query := bson.M{"public_id": publicID}

	var item sessions.Session

	if err := database.C(mdb.col).Find(query).One(&item); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record of Session type from db").
			With("query", query).
			With("collection", mdb.col).
			With("error", err.Error()))

		return sessions.Session{}, err
	}

	return item, nil

}

// Update uses a record from the db using the publicID and returns the sessions.Session type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given Session struct.
func (mdb *SessionDB) Update(ctx context.Context, publicID string, elem sessions.Session) error {
	m := metrics.NewTrace("SessionDB.Update")
	defer mdb.metrics.Emit(metrics.Info("SessionDB.Update").With("publicID", publicID).WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to finish, context has expired").
			With("collection", mdb.col).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session").
			With("publicID", publicID).
			With("collection", mdb.col).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to finish, context has expired").
			With("collection", mdb.col).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	if fields, ok := interface{}(elem).(SessionBSON); ok {
		query := fields.BSON()

		if err := database.C(mdb.col).Insert(query); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to update Session record").
				With("collection", mdb.col).
				With("public_id", publicID).
				With("query", query).
				With("error", err.Error()))

			return err
		}

		mdb.metrics.Emit(metrics.Info("Update record").
			With("collection", mdb.col).
			With("public_id", publicID).
			With("query", query).
			With("error", err.Error()))

		return nil
	}

	if fields, ok := interface{}(elem).(SessionFields); ok {
		query := bson.M(fields.Fields())

		if err := database.C(mdb.col).Insert(query); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to update Session record").
				With("query", query).
				With("public_id", publicID).
				With("collection", mdb.col).
				With("error", err.Error()))
			return err
		}

		mdb.metrics.Emit(metrics.Info("Create record").
			With("collection", mdb.col).
			With("query", query).
			With("public_id", publicID).
			With("error", err.Error()))

		return nil
	}

	query := bson.M{"publicID": publicID}
	queryData := bson.M(map[string]interface{}{

		"public_id": elem.PublicID,

		"token": elem.Token,

		"user_id": elem.UserID,
	})

	if err := database.C(mdb.col).Update(query, queryData); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to update Session record").
			With("collection", mdb.col).
			With("query", query).
			With("public_id", publicID).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(metrics.Info("Update record").
		With("collection", mdb.col).
		With("public_id", publicID).
		With("query", query))

	return nil
}

// Exec provides a function which allows the execution of a custom function against the collection.
func (mdb *SessionDB) Exec(ctx context.Context, fx func(col *mgo.Collection) error) error {
	m := metrics.NewTrace("SessionDB.Exec")
	defer mdb.metrics.Emit(metrics.Info("SessionDB.Exec").WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to execute operation").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	defer session.Close()

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to finish, context has expired").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	if err := fx(database.C(mdb.col)); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to execute operation").
			With("collection", mdb.col).
			With("error", err.Error()))
		return err
	}

	mdb.metrics.Emit(metrics.Info("Operation executed").
		With("collection", mdb.col))

	return nil
}
