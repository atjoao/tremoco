package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"music/controllers"
	"music/functions"
	"music/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*
//go:embed database/*
//go:embed assets/*
var embedded embed.FS

func engine() *gin.Engine {
	app := gin.Default()

	templ := template.Must(template.New("").ParseFS(embedded, "templates/*.tmpl"))
	app.SetHTMLTemplate(templ)

	store, err := postgres.NewStore(utils.StartConn(), []byte(os.Getenv("SESSION_KEY")))
	if err != nil {
		log.Panicln("Error creating session store > ", err)
	}

	app.Use(sessions.Sessions("sessions", store))

	// html render
	html := app.Group("/html")
	html.Use(controllers.AuthRequired)
	{
		html.GET("/sidebar", controllers.Sidebar)
	}

	// api routes
	api := app.Group("/api")
	api.Use(controllers.AuthRequired)
	{
		api.GET("/search", controllers.SearchVideos)
		api.GET("/video", controllers.VideoDataStream)
		api.GET("/stream/:audioId", controllers.StreamAudio)
		api.GET("/cover/:audioId", controllers.GetAudioCover)

		api.GET("/playlists", controllers.GetUserPlaylists)
		api.POST("/playlist/create", controllers.CreatePlaylist)
		api.GET("/playlist/get/:audioId", controllers.GetPlaylistsMusic)
		api.POST("/playlist/change", controllers.ChangePlaylist)
		api.GET("/playlist/:playlistId", controllers.GetPlaylist)
		api.DELETE("/playlist/delete/:playlistId", controllers.DeletePlaylist)
		api.POST("/playlist/edit/:playlistId", controllers.EditPlaylist)

		api.GET("/proxy", controllers.ProxyContent)

	}

	auth := app.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	// public routes

	assets, err := fs.Sub(embedded, "assets")
	if err != nil {
		panic(err)
	}

	app.StaticFS("/assets", http.FS(assets))

	app.GET("/", func(ctx *gin.Context) {
		const loadTemplate string = "dash.tmpl"
		if sessions.Default(ctx).Get("userId") != nil {
			ctx.HTML(200, loadTemplate, gin.H{
				"Title":        "Home",
				"loadTemplate": loadTemplate,
				"username":     sessions.Default(ctx).Get("username"),
			})
		} else {
			ctx.Redirect(302, "/login")
		}
	})

	app.GET("/login", func(ctx *gin.Context) {
		if sessions.Default(ctx).Get("userId") != nil {
			ctx.Redirect(302, "/")
		}
		const loadTemplate string = "login.tmpl"
		ctx.HTML(200, loadTemplate, gin.H{
			"Title":        "Login",
			"loadTemplate": loadTemplate,
		})
	})

	app.GET("/register", func(ctx *gin.Context) {
		if sessions.Default(ctx).Get("userId") != nil {
			ctx.Redirect(302, "/")
		}
		const loadTemplate string = "register.tmpl"
		ctx.HTML(200, loadTemplate, gin.H{
			"Title":        "Register",
			"loadTemplate": loadTemplate,
		})
	})

	app.GET("/logout", func(ctx *gin.Context) {
		if sessions.Default(ctx).Get("userId") != nil {
			ctx.Redirect(302, "/login")
		}

		sessions.Default(ctx).Clear()
		sessions.Default(ctx).Save()

		ctx.Redirect(302, "/login")
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
	files, err := embedded.ReadDir("database")

	if err != nil {
		log.Panicln("Error reading .sql files > ", err)
	}

	for _, file := range files {
		sql, err := embedded.ReadFile("database/" + file.Name())
		if err != nil {
			log.Panicln("Error reading", file.Name(), ".sql file > ", err)
		}

		var sqlStatement []string = strings.Split(string(sql), ";")
		for _, statement := range sqlStatement {
			fmt.Println(statement)
			_, err := dbConn.Exec(statement)
			if err != nil {
				log.Panicln("Error on", file.Name(), "command > ", err)
			}
		}
	}

	functions.ProcessAudioFiles()
	functions.RemoveMusicFromDb()

	app := engine()
	if err := app.Run(":3000"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}
