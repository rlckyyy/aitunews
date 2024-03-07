package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

const (
	RoleUser    = "user"
	RoleTeacher = "teacher"
	RoleAdmin   = "admin"
)

type User struct {
	ID             int
	Name           string
	Email          string
	Role           string
	HashedPassword []byte
}

type News struct {
	ID       int
	Title    string
	Content  string
	Created  time.Time
	Category string
}

type Comments struct {
	ID     int
	UserId int
	NewsId int
	Text   string
}
