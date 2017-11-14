package mdb_test

import (
     "encoding/json"


     "github.com/influx6/devapp/internals/sessions"

)

var sessionJSON = `{


    "expires":	nil,

    "twofactor_done":	false,

    "user_id":	"",

    "public_id":	"",

    "token":	""

}`

var sessionCreateJSON = `{


    "user_id":	"",

    "public_id":	"",

    "token":	"",

    "expires":	nil,

    "twofactor_done":	false

}`

var sessionUpdateJSON = `{


    "user_id":	"",

    "public_id":	"",

    "token":	"",

    "expires":	nil,

    "twofactor_done":	false

}`

// loadJSONFor returns a new instance of a sessions.Session from the provide json content.
func loadJSONFor(content string) (sessions.Session, error) {
	var elem sessions.Session

	if err := json.Unmarshal([]byte(content), &elem); err != nil {
		return sessions.Session{}, err
	}

	return elem, nil
}

