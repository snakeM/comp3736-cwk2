package main

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

const (
	// Database connection constants
	HOST     = "sc20mra-comp3736.database.windows.net"
	DATABASE = "sc20mra-comp3736-cwk2"
	USER     = "augerm"
	PASSWORD = "jos3chuf7YISH"
	PORT     = 1433
)

// Connect to the database; return a database object
func initDatabaseConnection() (*sql.DB, error) {
	var connectionString string = fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%d;database=%s;", 
		HOST, 
		USER, 
		PASSWORD, 
		PORT, 
		DATABASE)

	db, err := sql.Open("sqlserver", connectionString)
	return db, err
}

