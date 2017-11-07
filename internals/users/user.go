package users

import (
	"golang.org/x/crypto/bcrypt"

	uuid "github.com/satori/go.uuid"
)

const (
	hashComplexity = 10
	timeFormat     = "Mon Jan 2 15:04:05 -0700 MST 2006"
)

// UpdateUserPassword defines the set of data sent when updating a users password.
type UpdateUserPassword struct {
	PublicID        string `json:"public_id"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

//====================================================================================================

// NewUser holds details necessary for creating a new user.
type NewUser struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// User is a type defining the given user related fields for a given.
// @mongoapi
// @associates(@mongoapi, New, NewUser)
type User struct {
	Username  string `json:"username" bson:"username"`
	PublicID  string `json:"public_id" bson:"public_id"`
	PrivateID string `json:"private_id,omitempty" bson:"private_id"`
	Hash      string `json:"hash,omitempty" bson:"hash"`
	TwoFactor bool   `json:"two_factor_enabled,omitempty" bson:"two_factor_enabled"`
}

// New returns a new User instance based on the provided data.
func New(nw NewUser) (User, error) {
	var u User
	u.Username = nw.Username
	u.PublicID = uuid.NewV4().String()
	u.PrivateID = uuid.NewV4().String()

	u.ChangePassword(nw.Password)

	return u, nil
}

// Authenticate attempts to authenticate the giving password to the provided user.
func (u User) Authenticate(password string) error {
	pass := []byte(u.PrivateID + ":" + password)
	return bcrypt.CompareHashAndPassword([]byte(u.Hash), pass)
}

// SafeFields returns a map representing the data of the user with important
// security fields removed.
func (u User) SafeFields() map[string]interface{} {
	fields := u.Fields()

	delete(fields, "hash")
	delete(fields, "private_id")

	return fields
}

// Fields returns a map representing the data of the user.
func (u User) Fields() map[string]interface{} {
	fields := map[string]interface{}{
		"hash":               u.Hash,
		"username":           u.Username,
		"private_id":         u.PrivateID,
		"public_id":          u.PublicID,
		"two_factor_enabled": u.TwoFactor,
	}

	return fields
}

// ChangePassword uses the provided password to set the users password hash.
func (u *User) ChangePassword(password string) error {
	pass := []byte(u.PrivateID + ":" + password)
	hash, err := bcrypt.GenerateFromPassword(pass, hashComplexity)
	if err != nil {
		return err
	}

	u.Hash = string(hash)
	return nil
}
