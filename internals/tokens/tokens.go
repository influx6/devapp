package tokens

// TokenRecord defines a struct which stores all used twofactor token by user.
// @mongoapi
type TokenRecord struct {
	PublicID string   `bson:"public_id"`
	Tokens   []string `bson:"tokens"`
	UserID   string   `bson:"user_id"`
}
