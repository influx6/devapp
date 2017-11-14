// Package mdb provides a auto-generated package which contains a mongo CRUD API for the specific TokenRecord struct in package tokens.
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

	"github.com/influx6/devapp/internals/tokens"
)

// TokenRecordFields defines an interface which exposes method to return a map of all
// attributes associated with the defined structure as decided by the structure.
type TokenRecordFields interface {
	Fields() map[string]interface{}
}

// TokenRecordBSON defines an interface which exposes method to return a bson.M type
// which contains all related fields for the giving  object.
type TokenRecordBSON interface {
	BSON() bson.M
}

// TokenRecordBSONConsumer defines an interface which accepts a map of data which will be consumed
// into the giving implementing structure as decided by the structure.
type TokenRecordBSONConsumer interface {
	BSONConsume(bson.M) error
}

// TokenRecordConsumer defines an interface which accepts a map of data which will be consumed
// into the giving implementing structure as decided by the structure.
type TokenRecordConsumer interface {
	Consume(map[string]interface{}) error
}

// Mongod defines a interface which exposes a method for retrieving a
// mongo.Database and mongo.Session.
type Mongod interface {
	New() (*mgo.Database, *mgo.Session, error)
}

// TokenRecordDB defines a structure which provide DB CRUD operations
// using mongo as the underline db.
type TokenRecordDB struct {
	col             string
	db              Mongod
	metrics         metrics.Metrics
	ensuredIndex    bool
	incompleteIndex bool
	indexes         []mgo.Index
}

// New returns a new instance of TokenRecordDB.
func New(col string, m metrics.Metrics, mo Mongod, indexes ...mgo.Index) *TokenRecordDB {
	return &TokenRecordDB{
		db:      mo,
		col:     col,
		metrics: m,
		indexes: indexes,
	}
}

// ensureIndex attempts to ensure all provided indexes into the specific collection.
func (mdb *TokenRecordDB) ensureIndex() error {
	m := metrics.NewTrace("TokenRecordDB.ensureIndex")
	defer mdb.metrics.Emit(metrics.Info("TokenRecordDB.ensureIndex").WithTrace(m.End()))

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
func (mdb *TokenRecordDB) Count(ctx context.Context) (int, error) {
	m := metrics.NewTrace("TokenRecordDB.Count")
	defer mdb.metrics.Emit(metrics.Info("TokenRecordDB.Count").WithTrace(m.End()))

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
// on the given tokens.TokenRecord struct.
func (mdb *TokenRecordDB) Delete(ctx context.Context, publicID string) error {
	m := metrics.NewTrace("TokenRecordDB.Delete")
	defer mdb.metrics.Emit(metrics.Info("TokenRecordDB.Delete").With("publicID", publicID).WithTrace(m.End()))

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
// tokens.TokenRecord.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given TokenRecord struct.
func (mdb *TokenRecordDB) Create(ctx context.Context, elem tokens.TokenRecord) error {
	m := metrics.NewTrace("TokenRecordDB.Create")
	defer mdb.metrics.Emit(metrics.Info("TokenRecordDB.Create").With("publicID", elem.PublicID).WithTrace(m.End()))

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

	if fields, ok := interface{}(elem).(TokenRecordBSON); ok {
		if err := database.C(mdb.col).Insert(fields.BSON()); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to create TokenRecord record").
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

	if fields, ok := interface{}(elem).(TokenRecordFields); ok {
		if err := database.C(mdb.col).Insert(bson.M(fields.Fields())); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to create TokenRecord record").
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

		"tokens": elem.Tokens,

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
		mdb.metrics.Emit(metrics.Errorf("Failed to create TokenRecord record").
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

// GetAllPerPage retrieves all records from the db and returns a slice of tokens.TokenRecord type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given TokenRecord struct.
func (mdb *TokenRecordDB) GetAllPerPage(ctx context.Context, order string, orderBy string, page int, responsePerPage int) ([]tokens.TokenRecord, int, error) {
	m := metrics.NewTrace("TokenRecordDB.GetAll")
	defer mdb.metrics.Emit(metrics.Info("TokenRecordDB.GetAll").WithTrace(m.End()))

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

	var items []tokens.TokenRecord

	if err := database.C(mdb.col).Find(query).Skip(indexToStart).Limit(totalWanted).Sort(orderBy).All(&items); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve all records of TokenRecord type from db").
			With("collection", mdb.col).
			With("query", query).
			With("error", err.Error()))

		return nil, -1, err
	}

	return items, totalRecords, nil

}

// GetAllByOrder retrieves all records from the db and returns a slice of tokens.TokenRecord type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given TokenRecord struct.
func (mdb *TokenRecordDB) GetAllByOrder(ctx context.Context, order, orderBy string) ([]tokens.TokenRecord, error) {
	m := metrics.NewTrace("TokenRecordDB.GetAllByOrder")
	defer mdb.metrics.Emit(metrics.Info("TokenRecordDB.GetAllByOrder").WithTrace(m.End()))

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

	var items []tokens.TokenRecord

	if err := database.C(mdb.col).Find(query).Sort(orderBy).All(&items); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve all records of TokenRecord type from db").
			With("collection", mdb.col).
			With("query", query).
			With("error", err.Error()))

		return nil, err
	}

	return items, nil

}

// GetByField retrieves a record from the db using the provided field key and value
// returns the tokens.TokenRecord type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given TokenRecord struct.
func (mdb *TokenRecordDB) GetByField(ctx context.Context, key string, value interface{}) (tokens.TokenRecord, error) {
	m := metrics.NewTrace("TokenRecordDB.GetByField")
	defer mdb.metrics.Emit(metrics.Info("TokenRecordDB.GetByField").With(key, value).WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").With(key, value).
			With("collection", mdb.col).
			With("error", err.Error()))

		return tokens.TokenRecord{}, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))

		return tokens.TokenRecord{}, err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session").
			With(key, value).
			With("collection", mdb.col).
			With("error", err.Error()))

		return tokens.TokenRecord{}, err
	}

	defer session.Close()

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").With(key, value).
			With("collection", mdb.col).
			With("error", err.Error()))

		return tokens.TokenRecord{}, err
	}

	query := bson.M{key: value}

	var item tokens.TokenRecord

	if err := database.C(mdb.col).Find(query).One(&item); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve all records of TokenRecord type from db").
			With("query", query).
			With("collection", mdb.col).
			With("error", err.Error()))

		return tokens.TokenRecord{}, err
	}

	return item, nil

}

// Get retrieves a record from the db using the publicID and returns the tokens.TokenRecord type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given TokenRecord struct.
func (mdb *TokenRecordDB) Get(ctx context.Context, publicID string) (tokens.TokenRecord, error) {
	m := metrics.NewTrace("TokenRecordDB.Get")
	defer mdb.metrics.Emit(metrics.Info("TokenRecordDB.Get").With("publicID", publicID).WithTrace(m.End()))

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return tokens.TokenRecord{}, err
	}

	if err := mdb.ensureIndex(); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to apply index").
			With("collection", mdb.col).
			With("error", err.Error()))
		return tokens.TokenRecord{}, err
	}

	database, session, err := mdb.db.New()
	if err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to create session").
			With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return tokens.TokenRecord{}, err
	}

	defer session.Close()

	if context.IsExpired(ctx) {
		err := fmt.Errorf("Context has expired")
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve record").With("publicID", publicID).
			With("collection", mdb.col).
			With("error", err.Error()))
		return tokens.TokenRecord{}, err
	}

	query := bson.M{"public_id": publicID}

	var item tokens.TokenRecord

	if err := database.C(mdb.col).Find(query).One(&item); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to retrieve all records of TokenRecord type from db").
			With("query", query).
			With("collection", mdb.col).
			With("error", err.Error()))

		return tokens.TokenRecord{}, err
	}

	return item, nil

}

// Update uses a record from the db using the publicID and returns the tokens.TokenRecord type.
// Records using this DB must have a public id value, expressed either by a bson or json tag
// on the given TokenRecord struct.
func (mdb *TokenRecordDB) Update(ctx context.Context, publicID string, elem tokens.TokenRecord) error {
	m := metrics.NewTrace("TokenRecordDB.Update")
	defer mdb.metrics.Emit(metrics.Info("TokenRecordDB.Update").With("publicID", publicID).WithTrace(m.End()))

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

	query := bson.M{"public_id": publicID}

	if fields, ok := interface{}(elem).(TokenRecordBSON); ok {
		if err := database.C(mdb.col).Update(query, fields.BSON()); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to update TokenRecord record").
				With("collection", mdb.col).
				With("public_id", publicID).
				With("query", query).
				With("error", err.Error()))

			return err
		}

		mdb.metrics.Emit(metrics.Info("Update record").
			With("query", query).
			With("collection", mdb.col).
			With("public_id", publicID).
			With("data", fields.BSON()))

		return nil
	}

	if fields, ok := interface{}(elem).(TokenRecordFields); ok {
		if err := database.C(mdb.col).Update(query, fields.Fields()); err != nil {
			mdb.metrics.Emit(metrics.Errorf("Failed to update TokenRecord record").
				With("query", query).
				With("public_id", publicID).
				With("collection", mdb.col).
				With("error", err.Error()))
			return err
		}

		mdb.metrics.Emit(metrics.Info("Create record").
			With("collection", mdb.col).
			With("query", query).
			With("data", fields.Fields()).
			With("public_id", publicID))

		return nil
	}

	queryData := bson.M(map[string]interface{}{

		"public_id": elem.PublicID,

		"tokens": elem.Tokens,

		"user_id": elem.UserID,
	})

	if err := database.C(mdb.col).Update(query, queryData); err != nil {
		mdb.metrics.Emit(metrics.Errorf("Failed to update TokenRecord record").
			With("collection", mdb.col).
			With("query", query).
			With("data", queryData).
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
func (mdb *TokenRecordDB) Exec(ctx context.Context, fx func(col *mgo.Collection) error) error {
	m := metrics.NewTrace("TokenRecordDB.Exec")
	defer mdb.metrics.Emit(metrics.Info("TokenRecordDB.Exec").WithTrace(m.End()))

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
