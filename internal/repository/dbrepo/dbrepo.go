package dbrepo

import (
	"database/sql"

	"github.com/andkolbe/go-chat/internal/config"
	"github.com/andkolbe/go-chat/internal/repository"
)

// only works with Postgres
type postgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB // database connection pool. *sql.DB is for postgres
}

// lets us pass our connection pool and app config and return a repository
// because we return a pointer to postgresDPRepo, it will connect to postgres
// we can set up other repos for other dbs in the future
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo {
		App: a,
		DB: conn,
	}
}