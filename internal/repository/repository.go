package repository

// because our DB repo is hooked into our handlers, any functions listed here can be used on the handlers
type DatabaseRepo interface {
	AllUsers() bool
}