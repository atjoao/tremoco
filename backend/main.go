package main

import (
	"database/sql"
	"fmt"
	"log"
	"music/server/env"
	"music/server/functions"
	"music/server/routers"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
)

var dbConn *sql.DB = nil;

func engine() *gin.Engine {
	app := gin.New()
	
	store, err := postgres.NewStore(dbConn, []byte(env.SECRET_KEY))
	if err != nil {
		log.Panicln("Error creating session store > ", err)
	}

	app.Use(sessions.Sessions("sessions", store))

	// api routes
	api := app.Group("/api")
	{
		api.GET("/search", routers.SearchVideos)
		api.GET("/video", routers.VideoDataStream)
	}
	
	return app
}

func main() {

	var err error
	dbConn, err = sql.Open("postgres", "postgres://localdb@localhost:5432/music?sslmode=disable&")
	if err != nil {
		log.Panicln("Error connecting to postgres database > ", err)
	}

	log.Println("Connected to postgres database")
	log.Println("Executing .sql files")
	files, err := os.ReadDir("database")
	
	if err != nil {
		log.Panicln("Error reading .sql files > ", err)
	}

	for _, file := range files {
		sql, err := os.ReadFile("database/" + file.Name())
		if err != nil {
			log.Panicln("Error reading",file.Name(),".sql file > ", err)
		}

		var sqlStatement []string = strings.Split(string(sql), ";")
		for _, statement := range sqlStatement {
			fmt.Println(statement+";")
			result, err := dbConn.Exec(statement+";")
			if err != nil {
				log.Panicln("Error on",file.Name(),"command > ", err)
			}
			if result != nil {
				log.Println("Executed")
			}
		}
	}
	
	functions.ProcessAudioFiles(dbConn)

	app := engine()
	if err := app.Run(":3000"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}
