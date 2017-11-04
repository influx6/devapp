package mongoapi_test

import (
     "encoding/json"


     "github.com/influx6/devapp/internals/users"

)

var userJSON = `{


    "public_id":	"",

    "username":	"",

    "private_id":	""

}`

var userCreateJSON = `{


    "private_id":	"",

    "public_id":	"",

    "username":	""

}`

var userUpdateJSON = `{


    "username":	"",

    "private_id":	"",

    "public_id":	""

}`

// loadJSONFor returns a new instance of a users.User from the provide json content.
func loadJSONFor(content string) (users.User, error) {
	var elem users.User

	if err := json.Unmarshal([]byte(content), &elem); err != nil {
		return users.User{}, err
	}

	return elem, nil
}

