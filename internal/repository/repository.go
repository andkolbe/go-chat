package repository

import "github.com/andkolbe/go-chat/internal/models"

// because our DB repo is hooked into our handlers, any functions listed here can be used on the handlers
type DatabaseRepo interface {
	GetUserByID(id int) (models.User, error)
	AddUser(user models.User) error
	UpdateUser(user models.User) error
	Authenticate(username, testPassword string) (int, string, error)
}