// This file contains the repository implementation layer.
package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type Repository struct {
	Db *sql.DB
}

type NewRepositoryOptions struct {
	Dsn string
}

func NewRepository(opts NewRepositoryOptions) *Repository {
	db, err := sql.Open("postgres", opts.Dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(1 * time.Hour)

	return &Repository{
		Db: db,
	}
}
