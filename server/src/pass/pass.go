// Package defines helpers structs to store
// credentials
package pass

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

// DB provides access to a database
type DB struct {
	Host     string
	User     string
	Password string
	Name     string // of the database
	Port     int    // default to 5432
}

// NewDB uses env variables to build DB credentials :
// DB_HOST, DB_USER, DB_PASSWORD, DB_NAME
func NewDB() (out DB, err error) {
	out.Host = os.Getenv("DB_HOST")
	if out.Host == "" {
		return DB{}, errors.New("missing env DB_HOST")
	}

	out.User = os.Getenv("DB_USER")
	if out.User == "" {
		return DB{}, errors.New("missing env DB_USER")
	}

	out.Password = os.Getenv("DB_PASSWORD")
	if out.Password == "" {
		return DB{}, errors.New("missing env DB_PASSWORD")
	}

	out.Name = os.Getenv("DB_NAME")
	if out.Name == "" {
		return DB{}, errors.New("missing env DB_NAME")
	}

	return out, nil
}

// ConnectPostgres builds a connection string and
// connect using postgres as driver name.
func (db DB) ConnectPostgres() (*sql.DB, error) {
	port := db.Port
	if port == 0 {
		port = 5432
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		db.Host, port, db.User, db.Password, db.Name)
	return sql.Open("postgres", connStr)
}
