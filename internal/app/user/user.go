package user

import (
	"crypto/md5"
	error2 "error"
	"fmt"
	"generateId"
	"github.com/giantswarm/go.crypto/bcrypt"
)

type User struct {
	ID             string
	Email          string
	HashedPassword string
	Username       string
}
func (user *User) AvatarURL() string {
	return fmt.Sprintf(
		"//www.gravatar.com/avatar/%x",
		md5.Sum([]byte(user.Email)),
	)
}

func (user *User) ImagesRoute() string {
	return "/user/" + user.ID
}

const (
	hashCost       = 10
	passwordLength = 6
	userIDLength   = 16
)

func NewUser(username, email, password string) (User, error) {
	user := User{
		Email: email,
		Username: username,
	}
	if username == ""{
		return user, error2.ErrNoUsername
	}
	if email == "" {
		return user, error2.ErrNoEmail
	}

	if password == "" {
		return user, error2.ErrNoPassword
	}

	if len(password) < passwordLength {
		return user, error2.ErrPasswordTooShort
	}

	//Check if the username exists
	existingUser, err := GlobalUserStore.FindByUsername(username)

	if err != nil{
		return user, err
	}
	if existingUser != nil{
		return user, error2.ErrEmailExists
	}

	// Check if the email exists
	existingUser, err = GlobalUserStore.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, error2.ErrEmailExists
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	user.HashedPassword = string(hashedPassword)
	user.ID = generateId.GenerateID("usr", userIDLength)
	return user, err
}

func FindUser(username, password string) (*User, error) {
	out := &User{
		Username: username,
	}

	existingUser, err := GlobalUserStore.FindByUsername(username)
	if err != nil{
		return out, err
	}

	if existingUser == nil{
		return out, error2.ErrCredentialsIncorrect
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(existingUser.HashedPassword),
		[]byte(password),) != nil{
			return out, error2.ErrCredentialsIncorrect
	}

	return existingUser, nil
}

func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error)  {
	out := *user
	out.Email = email

	//Check if email exists
	existingUser, err := GlobalUserStore.FindByEmail(email)
	if err != nil {
		return out, err
	}
	if existingUser != nil && existingUser.ID != user.ID {
		return out, error2.ErrEmailExists
	}

	//Update email address
	user.Email = email

	//Check password
	if currentPassword == ""{
		return out, error2.ErrPasswordNotEmpty
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.HashedPassword),
		[]byte(currentPassword)) != nil{
			return out, error2.ErrPasswordIncorrect
	}

	if currentPassword == "" {
		return out, error2.ErrNoPassword
	}

	if len(newPassword) < passwordLength{
		return out, error2.ErrPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	user.HashedPassword = string(hashedPassword)

	return out, err
}

