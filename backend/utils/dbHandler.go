package utils

import (
	"database/sql"
	"log"
	"music/server/env"
)

var dbConn *sql.DB

func StartConn() *sql.DB{
	if dbConn != nil {
		return dbConn
	}

    var err error
    dbConn, err = sql.Open("postgres", env.POSTGRES_URI)
    if err != nil {
        log.Panicln("Error connecting to postgres database:", err)
    }

    err = dbConn.Ping()
    if err != nil {
        log.Panicln("Error pinging the database:", err)
    }

	return dbConn
}