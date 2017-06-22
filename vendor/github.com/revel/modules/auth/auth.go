package auth

import (
// "errors"
)

var (
	Store StorageDriver
)

// Store = gormauth.NewGormAuthDriver()

type UserAuth interface {
	// getters/setters implemented by the app-level model
	UserId() string
	Secret() string
	HashedSecret() string
	SetHashedSecret(string)

	SecretDriver
	// // implemented by secret driver
	// Authenticate() (bool, error)
}

type SecretDriver interface {
	Authenticate() (bool, error)
	HashSecret(args ...interface{}) (string, error)

	// stuff for documentation
	// UserContext is expected in these?

	// Secret expects 0 or non-0 arguments
	// When no parameter is passed, it acts as a getter
	// When one or more parameters are passed, it acts as a setter
	// A driver should specify the expected arguments and their meanings

	// Register()
	// Login()
	// Logout()
}

type StorageDriver interface {
	Save(user interface{}) error
	// Load should take a partially filled struct
	// (with values needed to look up)
	// and fills in the rest
	Load(user interface{}) error
}

// func init() {
// 	// auth.Store = gorm...
// }

// func (c App) Login(email, password string) {

// 	u := User {
// 		Email ...
// 	}

// 	good, err := auth.Authenticate(email, password)

// 	user, err := user_info.GetUserByEmail(email)
// }

// Bycrypt Authenticate() expects a single string argument of the plaintext password
// It returns true on success and false if error or password mismatch
// func Authenticate(attemptedUser UserAuth) (bool, error) {
// 	// check user in Store
// 	loadedUser, err := Store.Load(attemptedUser.UserId())
// 	if err != nil {
// 		return false, errors.New("User Not Found")
// 	}

// 	loadedUser.Authenticate(attemptedUser.Secret())

// 	// successfully authenticated
// 	return true, nil
// }
