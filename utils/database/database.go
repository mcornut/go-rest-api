package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Connect func
func Connect(user, password, dbname, host, port string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	return sql.Open("postgres", connStr)
}
