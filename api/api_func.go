package main

import "log"

func checkError(err error, msg string) {
	if err != nil {
		log.Printf(msg, err)
	}
}
func signIn(userName string) {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var user User{CreatedAt: time.Now().UnixNano(), UserName: userName, FollowsCount: 0, FollowersCount:0}
	err = dbmap.Insert(&user)
	checkError(err, "cannnot insert themeThought: %s")
}

func craeteThemeThought(userId int64, title string, content string) {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var thought Thought{CreatedAt:time.Now().UnixNano(), UserId: userId, Title: title, Content: content, IsTheme:true, LikesCount:0}
	err = dbmap.Insert(&thought)
	checkError(err, "cannnot insert themeThought: %s")
}

func craeteThought(userId int64, parentId int64, title string, content string) {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var thought Thought{CreatedAt:time.Now().UnixNano(), UserId: userId, ParentID: parentId, Title: title, Content: content, IsTheme:false, LikesCount:0}
	err = dbmap.Insert(&thought)
	checkError(err, "cannnot insert thought: %s")
}

func likeThought(userId int64, thoughtId int64) {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var usersLikeIT UsersLikeIT{UserId: userId, ThoughtId: thoughtId, CreatedAt: time.Now().UnixNano()}
	var thought Thought
	count, err := dbmap.SelectInt("select count(*) from usersLikeIT where user_id = $1 and where thought_id = $2", userId, thoughtId)
	checkError(err, "cannnot selectInt usersLikeIT: %s")
	if count >= 1 {
		_, err = dbmap.Exec("delete from usersLikeIT where user_id = $1 and where thought_id = $2", userId, thoughtId)
		checkError(err, "cannnot delete usersLikeIT: %s")

		err = dbmap.SelectOne(&thought, "select * from thoughts where thought_id = $1", thoughtId)
		checkError(err, "cannnot selectOne thought: %s")

		thought.LikesCount -= 1
		_, err = dbmap.Update(&thought)
		checkError(err, "cannnot update thought: %s")
	} else {
		err = dbmap.Insert(usersLikeIT)
		checkError(err, "cannnot insert usersLikeIT: %s")

		err = dbmap.SelectOne(&thought, "select * from thoughts where thought_id = $1", thoughtId)
		checkError(err, "cannnot selectOne thought: %s")

		thought.LikesCount += 1
		_, err = dbmap.Update(&thought)
		checkError(err, "cannnot update thought: %s")
	}
}

func detailUser(userId int64) User {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var user User
	err := dbmap.SelectOne(&user, "select * from users where user_id = $1", userId)
	checkError(err, "cannot selectOne user: %s")
	return user
}

func detailThought(thoughtId int64) Thought {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var thought Thought
	err := dbmap.SelectOne(&thought, "select * from thoughts where thought_id = $1", thoughtId)
	checkError(err, "cannot selectOne thought: %s")
	return thought
}

func getChildThoughs(parentId int64) []Thought {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var thoughts []Thought
	_, err := dbmap.Select(&thoughts, "select * from thoughts where parent_id = $1", parentId)
	checkError(err, "cannnot select thoughts: %s")
	return thoughts
}

func getParentThoughts() []Thought {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var thoughts []Thought
	_, err := dbmap.Select(&thoughts, "select * from thoughts where is_theme = true")
	checkError(err, "cannnot select thoughts: %s")
	return thoughts
}

func getAllThoughts() []Thought {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var thoughts []Thought
	_, err := dbmap.Select(&thoughts, "select * from thoughts order by created_at")
	checkError(err, "cannnot select thoughts: %s")
	return thoughts
}

