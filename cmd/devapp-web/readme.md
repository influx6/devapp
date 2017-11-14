
# Devapp API
DevApp API exposes a http server which serves a frontend server for the devapp project sample project.
It showcases how the project focused solely on `application/x-www-form-urlencoded` as means of data encoded format,
can be used to power create a application to provide the following:

- User Management and Authentication
- TwoFactor Authorization for User
- TwoFactor Managed Routes

## Run

To startup the project first ensure the following environment variables are set, has the project relies on MongoDB.

```bash
export DEVAPP_MONGO_USER=db
export DEVAPP_MONGO_DB=devapps
export DEVAPP_MONGO_AUTHDB=devapps
export DEVAPP_MONGO_PASSWORD=somepassword
export DEVAPP_MONGO_HOST=some-app.mlab.com:49535
```

The binary will loaded above environment variables to be able to work as needed, then:

```bash
> go run ./main.go
```
A server running on `http://127.0.0.1:3000` will be up for the application, just ensure port `3000` is not in used,
else run

```bash
> go run ./main.go -p 4050
```

## Note
The project development had limited time, so no test were included.
