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
		err = dbConn.Ping()
		if err != nil {
			dbConn = nil
			log.Println("Throwing away current connection...")
		} else {
			log.Println("Reusing database connection")
			return dbConn
		}
	}

	dbConn, err = sql.Open("sqlite", "./tremoco.db")
	if err != nil {
		log.Panicln("Error creating database:", err)
	}

	log.Println("Created new a new bond with the database")

	return dbConn
}
