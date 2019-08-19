package user

import (
	error2 "error"
	"testing"
)


func TestNewUserNoUsername(t *testing.T) {
	_, err := NewUser("", "d", "")
	if err != error2.ErrNoUsername{
		t.Error("Expected err to be errNoUserName")
	}
}

func TestNewUserNoPassword(t *testing.T)  {
	_, err := NewUser("user", "user@example.com", "")
	if err != error2.ErrNoPassword {
		t.Error("Expected err to be ErrNoPassword")
	}
}

type MockUserStore struct {
	findUser         *User
	findEmailUser    *User
	findUsernameUser *User
	saveUser         *User
}

func (store *MockUserStore) Find(string) (*User, error) {
	return store.findUser, nil
}

func (store *MockUserStore) FindByEmail(string) (*User, error) {
	return store.findEmailUser, nil
}

func (store *MockUserStore) FindByUsername(string) (*User, error) {
	return store.findUsernameUser, nil
}

func (store *MockUserStore) Save(user User) error {
	store.saveUser = &user
	return nil
}
func TestNewUserExistingUsername(t *testing.T)  {
	GlobalUserStore = &MockUserStore{
		findUsernameUser: &User{},
	}
	_, err := NewUser("user1", "user@example.com", "somepassword")
	if err != error2.ErrUsernameExists {
		t.Error("Expected err to be errUsernameExists")
	}
}
func TestNewUserExistingEmail(t *testing.T) {
	GlobalUserStore = &MockUserStore{
		findEmailUser: &User{},
	}

	_, err := NewUser("user", "user@example.com", "somepassword")
	if err != error2.ErrEmailExists {
		t.Error("Expected err to be errEmailExists")
	}
}