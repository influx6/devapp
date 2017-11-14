package mdb_test

import (
     "encoding/json"


     "github.com/influx6/devapp/internals/users"

)

var userJSON = `{


    "totp":	"",

    "username":	"",

    "public_id":	"",

    "private_id":	"",

    "hash":	"",

    "use_twofactor":	false

}`

var userCreateJSON = `{


    "totp":	"",

    "username":	"",

    "public_id":	"",

    "private_id":	"",

    "hash":	"",

    "use_twofactor":	false

}`

var userUpdateJSON = `{


    "totp":	"",

    "username":	"",

    "public_id":	"",

    "private_id":	"",

    "hash":	"",

    "use_twofactor":	false

}`

// loadJSONFor returns a new instance of a users.User from the provide json content.
func loadJSONFor(content string) (users.User, error) {
	var elem users.User

	if err := json.Unmarshal([]byte(content), &elem); err != nil {
		return users.User{}, err
	}

	return elem, nil
}

