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
)

type User struct {
	CreatedAt int64 `db:"created_at"`
	Id int64 `db:"user_id, primarykey, autoincrement"`
	UserName string `db:"user_name,size:128"`
	Profile string `db:"profile,size:512"`
	FollowsCount int `db:follows_count, default:0`
	FollowersCount int `db:followers_count, default:0`
}

type Thought struct {
	Id int64 `db:"thought_id, primarykey, autoincrement"`
	CreatedAt int64 `db:"created_at"`
	UserId int64 `db:"user_id"`
	ParentId int64 `db:"parent_id"`
	Title string `db:"title, size:128`
	Content string `db:"content, size:2048`
	IsTheme bool `db:"is_theme,default:false"`
	LikesCount int `db:"likes_count, default:0"` 
}

type UsersThoughtIT struct {
	Id int64 `db:"id, primarykey, autoincrement"`
	CreatedAt int64 `db:"created_at"`
	UserId int64 `db:"user_id"`
	ThoughtId int64 `db:"thought_id"`
}

type UsersLikeIT struct {
	Id int64 `db:"id, primarykey, autoincrement"`
	CreatedAt int64 `db:"created_at"`
	UserId int64 `db:"user_id"`
	ThoughtId int64 `db:"thought_id"`
}

type FollowIT struct {
	Id int64 `db:"id, primarykey, autoincrement"`
	CreatedAt int64 `db:"created_at"`
	FollowUserId int64 `db:"follow_user_id"`
	FollowedUserId int64 `db:"followed_user_id"`
}



func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "host=db user=app_user dbname=app_db password=postgres_password sslmode=disable port=5432")
	if err != nil {
		log.Printf("Failed to open db: %s",err.Error())
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Thought{}, "thoughts").SetKeys(true, "Id")
	dbmap.AddTableWithName(UsersThoughtIT{}, "usersThoughtIT").SetKeys(true, "Id")
	dbmap.AddTableWithName(UsersLikeIT{}, "usersLikeIT").SetKeys(true, "Id")
	dbmap.AddTableWithName(FollowIT{}, "followIT").SetKeys(true, "Id")
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
	dbmap.Db.Close()
	router := gin.Default()
	// router.LoadHTMLGlob("templates/*.html")
	v1 := router.Group("/v1")
	{
		v1.POST("/login", func(ctx *gin.Context) {
			userName := ctx.PostForm("userName")
			userId := funcs.SignIn(userName)
			ctx.JSON(200, gin.H{"userId":userId})
		})
		v1.GET("/allUsers", func(ctx *gin.Context) {
			users := funcs.AllUsers()
			ctx.JSON(200, gin.H{
				"users": users,
			})
		})
		v1.POST("/thought-submit", func(ctx *gin.Context) {
			userId, err := strconv.ParseInt(ctx.PostForm("userId"), 10, 64)
			funcs.CheckError(err, "cannnot strconv.Atoi: %s")
			parentId, err := strconv.ParseInt(ctx.PostForm("parentId"), 10, 64)
			funcs.CheckError(err, "cannnot strconv.Atoi: $s")
			title := ctx.PostForm("title")
			content := ctx.PostForm("content")
			thoughtId := funcs.CreateThought(userId, parentId, title, content)
			ctx.JSON(200, gin.H{"thoughtId":thoughtId})
		})
		v1.GET("/allThoughts", func(ctx *gin.Context) {
			thoughts := funcs.GetAllThoughts()
			ctx.JSON(200, gin.H{"tohughts":thoughts})
		})
		v1.GET("/thoughts-detail/:id", func(ctx *gin.Context) {
			thoughtId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
			funcs.CheckError(err, "cannnot strocnv.Atoi: %s")
			thought := funcs.DetailThought(thoughtId)
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