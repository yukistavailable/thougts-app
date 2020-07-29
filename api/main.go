package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	Follows []int
	Followers []int
}

type Thought struct {
	gorm.Model
	UserID uint
	ParentThoughtID uint `gorm:"default:0"`
	IsTheme bool `gorm:"default:false"`
	Content string
	LikeUsers []int
}



func main() {
	db, err := gorm.Open("postgres", "host=db user=app_user dbname=app_db password=postgres_password sslmode=disable port=5432")
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{})
	})
	router.Run()
}