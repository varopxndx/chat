package service

import (
	"database/sql"
	"fmt"

	"github.com/varopxndx/chat/config"

	_ "github.com/lib/pq"
)

// DB struct
type DB struct {
	db *sql.DB
}

// SetupDB creates DB struct
func SetupDB(database config.Database) *DB {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		database.Host,
		database.Port,
		database.User,
		database.Password,
		database.Name,
	)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(fmt.Sprintf("Fatal error creating DB connection: %v \n", err))
	}
	return &DB{db}
}

// Close closes DB connection
func (d DB) Close() error {
	return d.db.Close()
}
