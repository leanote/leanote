package secret

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"github.com/revel/modules/auth"
)

// example implementation of a Revel auth security driver
// This driver should be embedded into your app-level User model
// It expects your User model to have `Password` and `HashedPassword` string fields
//
// Your User model also needs to set itself as the UserContext for the BcryptAuth driver
//
// func NewUser(email, pass string) *User {
// 	u := &User{
// 		email:    email,
// 		password: pass,
// 	}
// 	u.UserContext = u
// }
//
type BcryptAuth struct {
	UserContext auth.UserAuth
}

// Bcrypt Secret() returns the hashed version of the password.
// It expects an argument of type string, which is the plain text password
func (self *BcryptAuth) HashSecret(args ...interface{}) (string, error) {
	if auth.Store == nil {
		return "", errors.New("Auth module StorageDriver not set")
	}
	argLen := len(args)
	if argLen == 0 {
		// we are getting
		return string(self.UserContext.HashedSecret()), nil
	}

	if argLen == 1 {
		// we are setting
		password, ok := args[0].(string)
		if !ok {
			return "", errors.New("Wrong argument type provided, expected plaintext password as string")
		}
		hPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		self.UserContext.SetHashedSecret(string(hPass))
		return self.UserContext.HashedSecret(), nil
	}

	// bad argument count
	return "", errors.New("Too many arguments provided, expected one")
}

// Bycrypt Authenticate() expects a single string argument of the plaintext password
// It returns true on success and false if error or password mismatch
func (self *BcryptAuth) Authenticate() (bool, error) {
	// check password
	err := bcrypt.CompareHashAndPassword([]byte(self.UserContext.HashedSecret()), []byte(self.UserContext.Secret()))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// successfully authenticated
	return true, nil
}
