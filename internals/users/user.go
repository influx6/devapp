package users

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	"github.com/influx6/twofactor"
	uuid "github.com/satori/go.uuid"
)

// consts ...
const (
	hashComplexity = 10
	timeFormat     = "Mon Jan 2 15:04:05 -0700 MST 2006"
	TwofactorOrg   = "devapps.inc"
)

var (
	// NilUser defines a nil type of user value useful for retrieving a User from
	// a httputil.Context.ValueBag
	NilUser = ((*User)(nil))
	// NillNewUser defines a nil version of NewUser struct used as unique key in httputil.Context.ValueBag.
	NillNewUser = ((*NewUser)(nil))
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
	TOTP          string `json:"totp" bson:"totp"`
	Username      string `json:"username" bson:"username"`
	PublicID      string `json:"public_id" bson:"public_id"`
	PrivateID     string `json:"private_id,omitempty" bson:"private_id"`
	Hash          string `json:"hash,omitempty" bson:"hash"`
	UseTwoFactor  bool   `json:"use_twofactor" bson:"use_twofactor"`
	SeenTwoFactor bool   `json:"seen_twofactor" bson:"seen_twofactor"`
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

// ValidateOTP validates provided OTP code from google authenticator.
func (u *User) ValidateOTP(userCode string) error {
	data, err := base64.StdEncoding.DecodeString(u.TOTP)
	if err != nil {
		return err
	}

	totp, err := twofactor.TOTPFromBytes(data, TwofactorOrg)
	if err != nil {
		return err
	}

	if verr := totp.Validate(userCode); verr != nil {
		return verr
	}

	totpbytes, err := totp.ToBytes()
	if err != nil {
		return err
	}

	u.TOTP = base64.StdEncoding.EncodeToString(totpbytes)
	return nil
}

// QR returns the User.TwoFactorQR data as string instead of bytes.
func (u User) QR() (string, error) {
	twqr, err := u.TwoFactorQR()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(twqr), nil
}

// TwoFactorQR returns QR code associated with given
func (u *User) TwoFactorQR() ([]byte, error) {
	totp, err := u.twoFactorObject()
	if err != nil {
		return nil, err
	}

	qr, err := totp.QR()
	if err != nil {
		return nil, err
	}

	totpbytes, err := totp.ToBytes()
	if err != nil {
		return nil, err
	}

	u.TOTP = base64.StdEncoding.EncodeToString(totpbytes)

	return qr, nil
}

// TwoFactorURL returns key url for the given user twofactor object.
func (u *User) TwoFactorURL() (string, error) {
	totp, err := u.twoFactorObject()
	if err != nil {
		return "", err
	}

	qtURL, err := totp.URL()
	if err != nil {
		return "", err
	}

	totpbytes, err := totp.ToBytes()
	if err != nil {
		return qtURL, err
	}

	u.TOTP = base64.StdEncoding.EncodeToString(totpbytes)

	return qtURL, nil
}

func (u User) twoFactorObject() (*twofactor.Totp, error) {
	data, err := base64.StdEncoding.DecodeString(u.TOTP)
	if err != nil {
		return nil, err
	}

	totp, err := twofactor.TOTPFromBytes(data, TwofactorOrg)
	if err != nil {
		return nil, err
	}

	return totp, nil
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
		"hash":           u.Hash,
		"username":       u.Username,
		"private_id":     u.PrivateID,
		"public_id":      u.PublicID,
		"totp":           u.TOTP,
		"use_twofactor":  u.UseTwoFactor,
		"seen_twofactor": u.SeenTwoFactor,
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
