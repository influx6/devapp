User DB API
===============================

User DB API attempts to provide a simple basic documentation which details
the basic structure of the User type, and the response which will be provided
when working with this API.

The API provides a basic CRUD interface as describe below:

## Create

Create stores a given User type into the mongo db, where the DB API expects the provided type to be called
against the API function type for adding a given record User.

```go
Create(ctx context.Context, elem users.User) error
```

## Get

Get retrives a giving record of the User type from the DB API returning received result matching
the publicID value provided if found else returning an error.

```go
Get(ctx context.Context, publicID string) (users.User, error)
```

## Get All

Get retrives all records of the User type from the DB API.

```go
GetAll(ctx context.Context) ([]users.User, error)
```

## Update

Update stores a given User type into the mongo db, where the DB API expects the provided type to be called
against the API function type for updating a existing record User.

```go
Update(ctx context.Context, publicID string, elem users.User) error
```

## Delete

Delete destroys the giving record of the User type from the DB API returning an error if any occured.

```
Delete(ctx context.Context, publicID string) error
```


## Example

```go

var (
	events = metrics.New(stdout.Stdout{})

	config = mongo.Config{
		Mode:     mgo.Monotonic,
		DB:       os.Getenv("dap_MONGO_DB"),
		Host:     os.Getenv("dap_MONGO_HOST"),
		User:     os.Getenv("dap_MONGO_USER"),
		AuthDB:   os.Getenv("dap_MONGO_AUTHDB"),
		Password: os.Getenv("dap_MONGO_PASSWORD"),
	}

)

func main() {
	col := "ignitor_collection"

	ctx := context.New()
	api := mongoapi.New(testCol, events, mongo.New(config))

	elem, err := loadJSONFor(ignitorCreateJSON)
	if err != nil {
    panic(err)
	}

	defer api.Delete(ctx, elem.PublicID)

	if err := api.Create(ctx, elem); err != nil {
    panic(err)
	}

	record, err := api.Get(ctx, elem.PublicID)
	if err != nil {
    panic(err)
	}

	records, total, err := api.GetAllPerPage(ctx, "asc", "public_id", -1, -1)
	if err != nil {
    panic(err)
	}

}
```
