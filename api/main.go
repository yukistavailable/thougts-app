package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"thoughts-app/api/config"
	"thoughts-app/api/utils"
)

type User struct {
	gorm.Model
	UserName string `gorm:"unique;not null"`
	Thoughts []Thought `gorm:"foreignkey:UserID"`
}

type UserProfile struct {
	Sentences string
	User User `gorm:"foreignkey:UserID"`
	UserID uint
	Follows []uint
	Followers []uint
}

type Thought struct {
	gorm.Model
	UserID uint
	ParentThoughtID uint `gorm:"default:0"`
	IsTheme bool `gorm:"default:false"`
	Content string
	LikeUsers []uint
}

func dbInit() {
	db, err := gorm.Open("postgres", "host=db user=app_user dbname=app_db password=postgres_password sslmode=disable port=5432")
	defer db.Close()
	if err != nil {
		log.Println(err.Error())
	}
	db.AutoMigrate(&User{}, &UserProfile{}, &Thought{})
}


func main() {
	utils.LoggingSettings(config.Config.LogFile)
	dbInit()
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})
	router.Run()
}