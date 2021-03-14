package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // postgres
)

type (
	// IDatabase is an interface that has database related functions
	IDatabase interface {
		Close() error
		Ping() error
	}

	// Database is a database instance
	Database struct {
		Context  context.Context
		Database *sql.DB
	}

	// Configuration is a database configuration
	Configuration struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		Schema   string
		Driver   string
	}
)

// New returns an instance to the database.
func New(ctx context.Context, config *Configuration) (*Database, error) {
	dbConfig := fmt.Sprintf(
		"host=%s port=%d user=%s password='%s' dbname=%s search_path=%s sslmode=require",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
		config.Schema,
	)

	db, err := sql.Open(config.Driver, dbConfig)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Database{
		Context:  ctx,
		Database: db,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.Database.Close()
}

// Ping is used to ping database connection
func (d *Database) Ping() error {
	return d.Database.Ping()
}
