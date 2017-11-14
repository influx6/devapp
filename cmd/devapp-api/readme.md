# Devapp API
DevApp API exposes a http server which serves has a API server for the devapp project sample project.

It showcases how the project focused solely on `application/x-www-form-urlencoded` as means of data encoded format
can be used to power up a frontend for the Devapp application to provide the following:

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

## API Actions

- Create new user

```bash
> curl -i http://127.0.0.1:3000/users/new\?username\=bob\&password\=boba\&password_confirm\=boba
HTTP/1.1 200 OK
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 14 Nov 2017 12:06:31 GMT
```

- Login User

```bash
> curl -i http://127.0.0.1:3000/session/new\?username\=bob\&password\=boba
HTTP/1.1 200 OK
Authorization: Bearer NjhlMjEyNjEtMDhmMi00MzhiLWI3MmQtNzQzYjg3N2UzY2VlOjE1MGEyOTg4LTc1NjgtNGExYy1iZmJhLTBiZjZiNmE2ZWI3Yg==
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 14 Nov 2017 12:27:15 GMT
Set-Cookie: Authorization=QmVhcmVyIE5qaGxNakV5TmpFdE1EaG1NaTAwTXpoaUxXSTNNbVF0TnpRellqZzNOMlV6WTJWbE9qRTFNR0V5T1RnNExUYzFOamd0TkdFeFl5MWlabUpoTFRCaVpqWmlObUUyWldJM1lnPT0=; Path=/; Expires=Thu, 16 Nov 2017 12:14:05 GMT
```

- Get User Profile with Authorization

```bash
> curl -i -H "Authorization: Bearer NjhlMjEyNjEtMDhmMi00MzhiLWI3MmQtNzQzYjg3N2UzY2VlOjE1MGEyOTg4LTc1NjgtNGExYy1iZmJhLTBiZjZiNmE2ZWI3Yg==" http://127.0.0.1:3000/profile
HTTP/1.1 200 OK
Content-Length: 3
Content-Type: text/plain
Date: Tue, 14 Nov 2017 12:31:11 GMT

bob
```

- Enable Google Authenticator Based Two Factor Login

```bash
> curl -i -H "Authorization: Bearer NjhlMjEyNjEtMDhmMi00MzhiLWI3MmQtNzQzYjg3N2UzY2VlOjE1MGEyOTg4LTc1NjgtNGExYy1iZmJhLTBiZjZiNmE2ZWI3Yg==" http://127.0.0.1:3000/users/twofactor/enable
HTTP/1.1 200 OK
Content-Length: 166
Content-Type: text/plain
Date: Tue, 14 Nov 2017 12:32:58 GMT

otpauth://totp/devapps.inc:68e21261-08f2-438b-b72d-743b877e3cee?algorithm=SHA1&counter=0&digits=6&issuer=devapps.inc&period=30&secret=DVYYBKPGMVIVRMNWLA4KPYDNJGTKABBH
```

Now be warned, this is a sample API, normally this route `users/twofactor/enable` should be well hidden behind the backend, in fact you are better of
using the route `users/twofactor/qr` which returns a `image/png`  because the text response returned by this route is the key URL format used by
Google authenticator, you can either use the `users/twofactor/qr` and stream the png to the page or generate a QR from the URL on the frontend, because
the URL returned contains an important secret needed for this all to be secure.


- Disable Google Authenticator Based Two Factor Login

```bash
> curl -i -H "Authorization: Bearer NjhlMjEyNjEtMDhmMi00MzhiLWI3MmQtNzQzYjg3N2UzY2VlOjE1MGEyOTg4LTc1NjgtNGExYy1iZmJhLTBiZjZiNmE2ZWI3Yg==" http://127.0.0.1:3000/users/twofactor/disable?token=846371
HTTP/1.1 200 OK
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 14 Nov 2017 12:36:42 GMT
```

Since twofactor was enabled, then all calls now will require the presence of `token` for all new session/login.
Note, all provided user token keys becomes invalid after use.

- Get User profile with Two Factor enabled

```bash
> curl -i -H "Authorization: Bearer NjhlMjEyNjEtMDhmMi00MzhiLWI3MmQtNzQzYjg3N2UzY2VlOjE1MGEyOTg4LTc1NjgtNGExYy1iZmJhLTBiZjZiNmE2ZWI3Yg==" http://127.0.0.1:3000/users/twofactor/disable?token=846371
HTTP/1.1 200 OK
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 14 Nov 2017 12:36:42 GMT

bob
```

Since twofactor was enabled, then all calls now will require the presence of `token` for all new session/login.
Note, all provided user token keys becomes invalid after use.

```bash
> curl -i -H "Authorization: Bearer NjhlMjEyNjEtMDhmMi00MzhiLWI3MmQtNzQzYjg3N2UzY2VlOjE1MGEyOTg4LTc1NjgtNGExYy1iZmJhLTBiZjZiNmE2ZWI3Yg==" http://127.0.0.1:3000/profile\?token\=786382
HTTP/1.1 400 Bad Request
Content-Length: 36
Content-Type: text/plain; charset=utf-8
Date: Tue, 14 Nov 2017 14:20:20 GMT
X-Content-Type-Options: nosniff

User already used token!
```

## Note
The project development had limited time, so no test were included.
