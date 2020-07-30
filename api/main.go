// package main

// import (
// 	"log"
// 	"github.com/gin-gonic/gin"
// 	"github.com/jinzhu/gorm"
// 	_ "github.com/jinzhu/gorm/dialects/postgres"
// 	"thoughts-app/api/config"
// 	"thoughts-app/api/utils"
// )

// type User struct {
// 	gorm.Model
// 	UserName string `gorm:"unique;not null"`
// 	Thoughts []Thought `gorm:"foreignkey:UserID"`
// 	LikeThoughts []Thought `gorm:"foreignkey:UserID"`
// }

// type Follows struct {
// 	FollowID uint
// }

// type UserProfile struct {
// 	Sentences string
// 	User User `gorm:"foreignkey:UserID"`
// 	UserID uint
// 	Follows []User `gorm:"foreignkey:UserID"`
// 	Followers []User `gorm:"foreignkey:UserID"`
// }

// type Thought struct {
// 	gorm.Model
// 	UserID uint
// 	ParentThoughtID uint `gorm:"default:0"`
// 	IsTheme bool `gorm:"default:false"`
// 	Content string
// 	LikeUsersNumber int
// }

// func dbInit() {
// 	db, err := gorm.Open("postgres", "host=db user=app_user dbname=app_db password=postgres_password sslmode=disable port=5432")
// 	defer db.Close()
// 	if err != nil {
// 		log.Printf("Failed to open db: %s",err.Error())
// 	}
// 	db.AutoMigrate(&User{}, &UserProfile{}, &Thought{})
// 	// db.AutoMigrate(&User{})
// 	// db.AutoMigrate(&UserProfile{})
// }

// func createUser(name string) {
// 	db, err := gorm.Open("postgres", "host=db user=app_user dbname=app_db password=postgres_password sslmode=disable port=5432")
// 	defer db.Close()
// 	if err != nil {
// 		log.Printf("Failed to open db: %s",err.Error())
// 	}
// 	var user User{UserName:name}

// 	db.Create(&user)
// 	db.Create(&UserProfile{UserID:user.ID})
// }

// func createThought(userID uint, parentThoughtID uint, isTheme bool, content string) {
// 	db, err := gorm.Open("postgres", "host=db user=app_user dbname=app_db password=postgres_password sslmode=disable port=5432")
// 	defer db.Close()
// 	if err != nil {
// 		log.Printf("Failed to open db: %s",err.Error())
// 	}
// 	db.Create(&Thought{UserID: userID, ParentThoughtID: parentThoughtID, IsTheme: isTheme, Content: content, LikeUsersNumber: 0})
// }

// func editUserProfile(userID uint, sentences string) {
// 	db, err := gorm.Open("postgres", "host=db user=app_user dbname=app_db password=postgres_password sslmode=disable port=5432")
// 	defer db.Close()
// 	if err != nil {
// 		log.Printf("Failed to open db: %s",err.Error())
// 	}
// 	var userProfile UserProfile{UserID:userID}
// 	db.Model(&userProfile).Update("Sentences", sentences)
// }

// func followUser(followUserID uint, followedUserID uint) {
// 	db, err := gorm.Open("postgres", "host=db user=app_user dbname=app_db password=postgres_password sslmode=disable port=5432")
// 	defer db.Close()
// 	if err != nil {
// 		log.Printf("Failed to open db: %s",err.Error())
// 	}

// }



// func main() {
// 	utils.LoggingSettings(config.Config.LogFile)
// 	dbInit()
// 	router := gin.Default()
// 	v1 := router.Group("/v1")
// 	{
// 		v1.POST("/login", loginEndpoint)
// 		v1.POST("/thought-submit", submitEndpoint)
// 		v1.POST("/thoughts-detail/:id", detailEndpoint)
// 	}
// 	router.LoadHTMLGlob("templates/*.html")
// 	router.GET("/", func(ctx *gin.Context) {
// 		ctx.HTML(200, "index.html", gin.H{})
// 	})
// 	router.Run()
// }