package main

import (
	"database/sql"
	"log"
	"music/server/env"
	"music/server/routers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
)

var Db *sql.DB = nil;

func engine() *gin.Engine {
	app := gin.New()
	
	store, err := postgres.NewStore(Db, []byte(env.SECRET_KEY))
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
	Db, err = sql.Open("postgres", "postgres://localdb@localhost:5432/music?sslmode=disable")
	if err != nil {
		log.Panicln("Error connecting to postgres database > ", err)
	}
	
	log.Println("Connected to postgres database")
	app := engine()
	if err := app.Run(":3000"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}
