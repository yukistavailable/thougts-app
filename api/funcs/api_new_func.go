package funcs

import (
	"database/sql"
	"log"
	"time"
	"gopkg.in/gorp.v2"
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

func openDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "host=db user=app_user dbname=app_db password=postgres_password sslmode=disable port=5432")
	if err != nil {
		log.Printf("Failed to open db: %s",err.Error())
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	return dbmap
}


func CheckError(err error, msg string) {
	if err != nil {
		log.Printf(msg, err)
	}
}
func SignIn(userName string, dbmap *gorp.DbMap) int64 {
	// dbmap := openDb()
	// defer dbmap.Db.Close()
	user := User{CreatedAt: time.Now().UnixNano(), UserName: userName, FollowsCount: 0, FollowersCount:0}
	err := dbmap.Insert(&user)
	CheckError(err, "cannnot insert user: %s")
	return user.Id
}

func CraeteThemeThought(userId int64, title string, content string, dbmap *gorp.DbMap) {
	// dbmap := openDb()
	// defer dbmap.Db.Close()
	thought := Thought{CreatedAt:time.Now().UnixNano(), UserId: userId, Title: title, Content: content, IsTheme:true, LikesCount:0}
	err := dbmap.Insert(&thought)
	CheckError(err, "cannnot insert themeThought: %s")
}

func CreateThought(userId int64, parentId int64, title string, content string, dbmap *gorp.DbMap) int64 {
	// dbmap := openDb()
	// defer dbmap.Db.Close()
	thought := Thought{CreatedAt:time.Now().UnixNano(), UserId: userId, ParentId: parentId, Title: title, Content: content, IsTheme:false, LikesCount:0}
	err := dbmap.Insert(&thought)
	CheckError(err, "cannnot insert thought: %s")
	return thought.Id
}

func LikeThought(userId int64, thoughtId int64) {
	dbmap := openDb()
	defer dbmap.Db.Close()
	usersLikeIT := UsersLikeIT{UserId: userId, ThoughtId: thoughtId, CreatedAt: time.Now().UnixNano()}
	var thought Thought
	count, err := dbmap.SelectInt("select count(*) from usersLikeIT where user_id = $1 and where thought_id = $2", userId, thoughtId)
	CheckError(err, "cannnot selectInt usersLikeIT: %s")
	if count >= 1 {
		_, err = dbmap.Exec("delete from usersLikeIT where user_id = $1 and where thought_id = $2", userId, thoughtId)
		CheckError(err, "cannnot delete usersLikeIT: %s")

		err = dbmap.SelectOne(&thought, "select * from thoughts where thought_id = $1", thoughtId)
		CheckError(err, "cannnot selectOne thought: %s")

		thought.LikesCount--
		_, err = dbmap.Update(&thought)
		CheckError(err, "cannnot update thought: %s")
	} else {
		err = dbmap.Insert(usersLikeIT)
		CheckError(err, "cannnot insert usersLikeIT: %s")

		err = dbmap.SelectOne(&thought, "select * from thoughts where thought_id = $1", thoughtId)
		CheckError(err, "cannnot selectOne thought: %s")

		thought.LikesCount++
		_, err = dbmap.Update(&thought)
		CheckError(err, "cannnot update thought: %s")
	}
}

func DetailUser(userId int64, dbmap *gorp.DbMap) User {
	// dbmap := openDb()
	// defer dbmap.Db.Close()
	var user User
	err := dbmap.SelectOne(&user, "select * from users where user_id = $1", userId)
	CheckError(err, "cannot selectOne user: %s")
	return user
}

func AllUsers(dbmap *gorp.DbMap) []User {
	// dbmap := openDb()
	// defer dbmap.Db.Close()
	var users []User
	_ , err := dbmap.Select(&users, "select * from users order by created_at")
	CheckError(err, "cannnot select users")
	return users
}

func DetailThought(thoughtId int64, dbmap *gorp.DbMap) Thought {
	// dbmap := openDb()
	// defer dbmap.Db.Close()
	var thought Thought
	err := dbmap.SelectOne(&thought, "select * from thoughts where thought_id = $1", thoughtId)
	CheckError(err, "cannot selectOne thought: %s")
	return thought
}

func GetChildThoughs(parentId int64) []Thought {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var thoughts []Thought
	_, err := dbmap.Select(&thoughts, "select * from thoughts where parent_id = $1", parentId)
	CheckError(err, "cannnot select thoughts: %s")
	return thoughts
}

func GetParentThoughts() []Thought {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var thoughts []Thought
	_, err := dbmap.Select(&thoughts, "select * from thoughts where is_theme = true")
	CheckError(err, "cannnot select thoughts: %s")
	return thoughts
}

func GetAllThoughts(dbmap *gorp.DbMap) []Thought {
	// dbmap := openDb()
	// defer dbmap.Db.Close()
	var thoughts []Thought
	_, err := dbmap.Select(&thoughts, "select * from thoughts order by created_at")
	CheckError(err, "cannnot select thoughts: %s")
	return thoughts
}

