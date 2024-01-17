package main

import (
	"music/server/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	
	r := app.Group("/test")
	routers.TestGroup(r)

	app.Run(":3000")
}