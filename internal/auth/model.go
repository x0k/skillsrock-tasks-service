package auth

import "errors"

var ErrLoginIsTaken = errors.New("this login is already taken")
var ErrUserNotFound = errors.New("user not found")
var ErrPasswordsMismatch = errors.New("passwords mismatch")

type User struct {
	Login        string
	PasswordHash []byte
}

func NewUser(
	login string,
	passwordHash []byte,
) *User {
	return &User{login, passwordHash}
}
