package main

mport (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db, err := gorm.Open("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable port=5432")
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html")
	})
}