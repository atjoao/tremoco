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

func engine() *gin.Engine {
	app := gin.New()
	
	db, err := sql.Open("postgres", "postgres://localdb@localhost:5432/music?sslmode=disable")
	if err != nil {
		log.Panicln("Error connecting to postgres database > ", err)
	} 

	store, err := postgres.NewStore(db, []byte(env.SECRET_KEY))
	if err != nil {
		log.Panicln("Error creating session store > ", err)
	}

	app.Use(sessions.Sessions("sessions", store))

	video := app.Group("/api")
	{
		video.GET("/search", routers.SearchVideos)
		video.GET("/video", routers.VideoDataStream)
	}
	
	return app
}

func main() {
	
	app := engine()
	if err := app.Run(":3000"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}
