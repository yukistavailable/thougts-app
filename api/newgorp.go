package main

import (
	"database/sql"
	"log"
	"strconv"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/gorp.v2"
	"thoughts-app/api/config"
	"thoughts-app/api/utils"
	"thoughts-app/api/funcs"
	// "fmt"
)

type jsonUser struct {
	Name string `json:"name"`
}

type jsonThought struct {
	userId int64 `json:"user_id"`
	parentId int64 `json:"parent_id"`
	title string `json:title`
	content string `json:content`
}


func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "host=db user=app_user dbname=app_db password=postgres_password sslmode=disable port=5432")
	if err != nil {
		log.Printf("Failed to open db: %s",err.Error())
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(funcs.User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(funcs.Thought{}, "thoughts").SetKeys(true, "Id")
	dbmap.AddTableWithName(funcs.UsersThoughtIT{}, "usersThoughtIT").SetKeys(true, "Id")
	dbmap.AddTableWithName(funcs.UsersLikeIT{}, "usersLikeIT").SetKeys(true, "Id")
	dbmap.AddTableWithName(funcs.FollowIT{}, "followIT").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Printf("Failed to create table: %s",err.Error())
	}
	return dbmap
}


func main() {
	utils.LoggingSettings(config.Config.LogFile)
	log.Println("Run Server")
	dbmap := initDb()
	router := gin.Default()
	// router.LoadHTMLGlob("templates/*.html")
	v1 := router.Group("/v1")
	{
		v1.POST("/login", func(ctx *gin.Context) {
			var userName jsonUser
			err := ctx.BindJSON(&userName)
			funcs.CheckError(err, "cannot bind json user : %s")
			userId := funcs.SignIn(userName.Name, dbmap)
			ctx.JSON(200, gin.H{"userId":userId})
		})
		v1.GET("/allUsers", func(ctx *gin.Context) {
			users := funcs.AllUsers(dbmap)
			ctx.JSON(200, gin.H{
				"users": users,
			})
		})
		v1.POST("/thought-submit", func(ctx *gin.Context) {
			// userId, err := strconv.ParseInt(ctx.PostForm("userId"), 10, 64)
			// funcs.CheckError(err, "cannnot strconv.Atoi: %s")
			// parentId, err := strconv.ParseInt(ctx.PostForm("parentId"), 10, 64)
			// funcs.CheckError(err, "cannnot strconv.Atoi: $s")
			// title := ctx.PostForm("title")
			// content := ctx.PostForm("content")
			var thought jsonThought
			err := ctx.BindJSON(&thought)
			funcs.CheckError(err, "cannot bind json thought: %s")
			thoughtId := funcs.CreateThought(thought.userId, thought.parentId, thought.title, thought.content, dbmap)
			ctx.JSON(200, gin.H{"thoughtId":thoughtId})
		})
		v1.POST("/theme-thought-submit", func(ctx *gin.Context) {
			var thought jsonThought
			err := ctx.BindJSON(&thought)
			funcs.CheckError(err, "cannnot bind json themethought: %s")
			thoughtId := funcs.CreateThemeThought(thought.userId, thought.title, thought.content, dbmap)
			ctx.JSON(200, gin.H{"thoughtId":thoughtId})
		})
		v1.GET("/allThoughts", func(ctx *gin.Context) {
			thoughts := funcs.GetAllThoughts(dbmap)
			ctx.JSON(200, gin.H{"tohughts":thoughts})
		})
		v1.GET("/thoughts-detail/:id", func(ctx *gin.Context) {
			thoughtId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
			funcs.CheckError(err, "cannnot strocnv.Atoi: %s")
			thought := funcs.DetailThought(thoughtId, dbmap)
			ctx.JSON(200, gin.H{"thought":thought})
		})

	}
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run()
}