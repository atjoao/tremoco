package utils

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var dbConn *sql.DB

func StartConn() *sql.DB {
	var err error
	if dbConn != nil {
		return dbConn
	}

	dbConn, err = sql.Open("sqlite", "./tremoco.db")
	if err != nil {
		log.Panicln("Error creating database:", err)
	}

	return dbConn
}
