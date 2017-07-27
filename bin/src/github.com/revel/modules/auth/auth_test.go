package auth_test

import (
	"errors"
	"testing"

	"github.com/revel/modules/auth"
	"github.com/revel/modules/auth/driver/secret"
)

type User struct {
	email    string
	password string
	hashpass string

	secret.BcryptAuth // SecurityDriver for testing
}

func NewUser(email, pass string) *User {
	u := &User{
		email:    email,
		password: pass,
	}
	u.UserContext = u
	return u
}

func (self *User) UserId() string {
	return self.email
}

func (self *User) Secret() string {
	return self.password
}

func (self *User) HashedSecret() string {
	return self.hashpass
}

func (self *User) SetHashedSecret(hpass string) {
	self.hashpass = hpass
}

// func (self *User) Load() string

type TestStore struct {
	data map[string]string
}

func (self *TestStore) Save(user interface{}) error {
	u, ok := user.(*User)
	if !ok {
		return errors.New("TestStore.Save() expected arg of type User")
	}

	hPass, err := u.HashSecret(u.Secret())
	if err != nil {
		return err
	}
	self.data[u.UserId()] = hPass

	return nil
}
func (self *TestStore) Load(user interface{}) error {
	u, ok := user.(*User)
	if !ok {
		return errors.New("TestStore.Load() expected arg of type User")
	}

	hpass, ok := self.data[u.UserId()]
	if !ok {
		return errors.New("Record Not Found")
	}
	u.SetHashedSecret(hpass)
	return nil
}

func TestPasswordHash(t *testing.T) {
	auth.Store = &TestStore{
		data: make(map[string]string),
	}
	u := NewUser("demo@domain.com", "demopass")
	fail := NewUser("demo@domain.com", "")

	var err error
	u.hashpass, err = u.HashSecret(u.password)
	if err != nil {
		t.Errorf("Should have hashed password, get error: %v\n", err)
	}
	fail.hashpass, err = fail.HashSecret(fail.password)
	if err == nil {
		t.Errorf("Should have failed hashing\n")
	}
}

func TestAuthenticate(t *testing.T) {
	auth.Store = &TestStore{
		data: make(map[string]string),
	}

	// user registered a long time ago
	u := NewUser("demo@domain.com", "demopass")
	err := auth.Store.Save(u)
	if err != nil {
		t.Errorf("Should have saved user: %v", err)
	}

	// users now logging in
	pass := NewUser("demo@domain.com", "demopass")
	fail := NewUser("demo@domain.com", "invalid")

	// valid user is now trying to login
	// check user in DB
	err = auth.Store.Load(pass)
	if err != nil {
		t.Errorf("Should have loaded pass user: %v\n", err)
	}
	// check credentials
	ok, err := pass.Authenticate()
	if !ok || err != nil {
		t.Errorf("Should have authenticated user")
	}

	// invalid user is now trying to login
	err = auth.Store.Load(fail)
	if err != nil {
		t.Errorf("Should have loaded fail user")
	}
	// this should fail
	ok, err = fail.Authenticate()
	if ok || err != nil {
		t.Errorf("Should have failed to authenticate user: %v\n", err)
	}
}
