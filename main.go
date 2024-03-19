package main

import (
	"fmt"
	"log"
	"music/server/controllers"
	"music/server/env"
	"music/server/functions"
	"music/server/utils"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
)

func engine() *gin.Engine {
	app := gin.New()
	app.LoadHTMLGlob("templates/*")
	
	store, err := postgres.NewStore(utils.StartConn(), []byte(env.SECRET_KEY))
	if err != nil {
		log.Panicln("Error creating session store > ", err)
	}

	app.Use(sessions.Sessions("sessions", store))

	// api routes
	api := app.Group("/api")
	api.Use(controllers.AuthRequired)
	{
		api.GET("/search", controllers.SearchVideos)
		api.GET("/video", controllers.VideoDataStream)
		api.GET("/stream/:audioId", controllers.StreamAudio)
		api.GET("/cover/:audioId", controllers.GetAudioCover)
	}

	auth := app.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	// public routes
	app.Static("/assets", "./assets")

	app.GET("/", func(ctx *gin.Context) {
		const loadTemplate string = "dash.tmpl"
		if sessions.Default(ctx).Get("userId") != nil {
			ctx.HTML(200, loadTemplate, gin.H{
				"Title": "Inicio",
				"loadTemplate": loadTemplate,
			})
		} else {
			ctx.Redirect(302, "/login")
		}
	})

	app.GET("/login", func(ctx *gin.Context) {
		const loadTemplate string = "login.tmpl"
		ctx.HTML(200, loadTemplate, gin.H{
			"Title": "Login",
			"loadTemplate": loadTemplate,
		})
	})

	return app
}

func main() {
	var err error
	dbConn := utils.StartConn()
	if dbConn != nil {
		log.Println("Connected to postgres database")
	}

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
			fmt.Println(statement)
			result, err := dbConn.Exec(statement)
			if err != nil {
				log.Panicln("Error on",file.Name(),"command > ", err)
			}
			if result != nil {
				log.Println("Executed")
			}
		}
	}
	
	functions.ProcessAudioFiles()

	app := engine()
	if err := app.Run(":3000"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}
