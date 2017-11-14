package mdb_test

import (
     "encoding/json"


     "github.com/influx6/devapp/internals/tokens"

)

var tokenrecordJSON = `{


    "public_id":	"",

    "tokens":	[]string{},

    "user_id":	""

}`

var tokenrecordCreateJSON = `{


    "public_id":	"",

    "tokens":	[]string{},

    "user_id":	""

}`

var tokenrecordUpdateJSON = `{


    "public_id":	"",

    "tokens":	[]string{},

    "user_id":	""

}`

// loadJSONFor returns a new instance of a tokens.TokenRecord from the provide json content.
func loadJSONFor(content string) (tokens.TokenRecord, error) {
	var elem tokens.TokenRecord

	if err := json.Unmarshal([]byte(content), &elem); err != nil {
		return tokens.TokenRecord{}, err
	}

	return elem, nil
}

