package controllers

import (
	"database/sql"
	"music/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func CheckHashPassword(password string, pswdhash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pswdhash), []byte(password))
	return err == nil
}

func AuthRequired(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	if userId == nil {
		ctx.AbortWithStatusJSON(403, gin.H{
			"status":  "UNAUTHORIZED",
			"message": "Login required",
		})
		return
	}
	ctx.Next()
}

func Login(ctx *gin.Context) {
	var err error

	session := sessions.Default(ctx)
	var db *sql.DB = utils.StartConn()

	const sql string = "SELECT username, passwordHash, id FROM users WHERE username = $1"
	var username string = ctx.PostForm("username")
	var password string = ctx.PostForm("password")

	if username == "" || password == "" {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": "username is empty",
		})
		return
	}

	var user utils.User
	err = db.QueryRow(sql, username).Scan(&user.Username, &user.Password, &user.Id)
	if err != nil {
		ctx.JSON(404, gin.H{
			"status":  "NOT_FOUND",
			"message": "User not found",
		})
		return
	}

	if !CheckHashPassword(password, user.Password) {
		ctx.JSON(401, gin.H{
			"status":  "WRONG_PASSWORD",
			"message": "Wrong password",
		})
		return
	}

	session.Set("userId", user.Id)
	session.Set("username", user.Username)
	if err := session.Save(); err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "OK",
		"message": "User logged in",
	})

}

func Register(ctx *gin.Context) {
	var err error
	session := sessions.Default(ctx)
	var db *sql.DB = utils.StartConn()

	const sql string = "INSERT INTO users (username, passwordHash) VALUES ($1, $2) RETURNING id;"

	var username string = ctx.PostForm("username")
	var password string = ctx.PostForm("password")
	if username == "" || password == "" {
		ctx.JSON(400, gin.H{
			"status":  "MISSED_PARAMS",
			"message": "username or password is empty",
		})
		return
	}
	var hashedPswd string = HashPassword(password)
	if hashedPswd == "" {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Internal Server Error",
		})
		return
	}

	var userId int
	err = db.QueryRow(sql, username, hashedPswd).Scan(&userId)
	if err != nil {
		ctx.JSON(409, gin.H{
			"status":  "ACCOUNT_EXISTS",
			"message": "An account with this username already exists",
		})
		return
	}

	session.Set("userId", userId)
	session.Set("username", username)
	if err := session.Save(); err != nil {
		ctx.JSON(500, gin.H{
			"status":  "SERVER_ERROR",
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "OK",
		"message": "User registred",
	})
}
