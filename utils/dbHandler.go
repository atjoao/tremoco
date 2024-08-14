package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var dbConn *sql.DB

func StartConn() *sql.DB {
	if dbConn != nil {
		return dbConn
	}

	var err error
	connStr := os.Getenv("POSTGRES_URI")
	if connStr == "" {
		log.Panicln("POSTGRES_URI environment variable is not set")
	}

	fmt.Println("Connecting to database: ", connStr)

	dbConn, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Panicln("Error connecting to postgres database:", err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Panicln("Error pinging the database:", err)
	}

	return dbConn
}
