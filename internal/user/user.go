package user

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("user not found")

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
