package users

// NewUser holds details necessary for creating a new user.
type NewUser struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// User holds given data about user credentials.
// @mongoapi
// @associates(@mongoapi, New, NewUser)
type User struct {
	Username  string `json:"username" bson:"username"`
	PrivateID string `json:"private_id" bson:"private_id"`
	PublicID  string `json:"public_id" bson:"public_id"`
}

// Fields returns the values for given struct as map.
func (u User) Fields() map[string]interface{} {
	return map[string]interface{}{
		"username":   u.Username,
		"public_id":  u.PublicID,
		"private_id": u.PrivateID,
	}
}
