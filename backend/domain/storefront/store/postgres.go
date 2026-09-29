package store

import "database/sql"

type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository constructs a PostgreSQL TODO repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
