package mdb_test

import (
     "encoding/json"


     "github.com/influx6/devapp/internals/users"

)

var userJSON = `{


    "username":	"",

    "public_id":	"",

    "private_id":	"",

    "hash":	"",

    "two_factor_enabled":	false

}`

var userCreateJSON = `{


    "private_id":	"",

    "hash":	"",

    "two_factor_enabled":	false,

    "username":	"",

    "public_id":	""

}`

var userUpdateJSON = `{


    "two_factor_enabled":	false,

    "username":	"",

    "public_id":	"",

    "private_id":	"",

    "hash":	""

}`

// loadJSONFor returns a new instance of a users.User from the provide json content.
func loadJSONFor(content string) (users.User, error) {
	var elem users.User

	if err := json.Unmarshal([]byte(content), &elem); err != nil {
		return users.User{}, err
	}

	return elem, nil
}

