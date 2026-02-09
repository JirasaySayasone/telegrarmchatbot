package db

import "database/sql"

func Connect() (*sql.DB, error) {
	return sql.Open("postgres", "postgres://user:password@localhost:8080/dbname")
}
