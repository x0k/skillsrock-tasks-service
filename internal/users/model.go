package users

import "errors"

var ErrLoginIsTaken = errors.New("this login is already taken")
var ErrUserNotFound = errors.New("user not found")

type user struct {
	Login        string
	PasswordHash []byte
}

func newUser(
	login string,
	passwordHash []byte,
) *user {
	return &user{login, passwordHash}
}
