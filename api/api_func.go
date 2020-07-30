package main

import "log"

func signIn(userName string) {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var user User{CreatedAt: time.Now().UnixNano(), UserName: userName, FollowsCount: 0, FollowersCount:0}
	err = dbmap.Insert(&user)
	if err != nil {
		log.Printf("cannnot insert themeThought: %s", err.Error())
	}
}

func craeteThemeThought(userId int64, title string, content string) {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var thought Thought{CreatedAt:time.Now().UnixNano(), UserId: userId, Title: title, Content: content, IsTheme:true, LikesCount:0}
	err = dbmap.Insert(&thought)
	if err != nil {
		log.Printf("cannnot insert themeThought: %s", err.Error())
	}
}

func craeteThought(userId int64, parentId int64, title string, content string) {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var thought Thought{CreatedAt:time.Now().UnixNano(), UserId: userId, ParentID: parentId, Title: title, Content: content, IsTheme:false, LikesCount:0}
	err = dbmap.Insert(&thought)
	if err != nil {
		log.Printf("cannnot insert thought: %s", err.Error())
	}
}

func likeThought(userId int64, thoughtId int64) {
	dbmap := openDb()
	defer dbmap.Db.Close()
	var usersLikeIT UsersLikeIT{UserId: userId, ThoughtId: thoughtId, CreatedAt: time.Now().UnixNano()}
	var thought Thought
	count, err := dbmap.SelectInt("select count(*) from usersLikeIT where user_id = $1 and where thought_id = $2", userId, thoughtId)
	if err != nil {
		log.Printf("cannnot selectInt usersLikeIT: %s", err.Error())
	}
	if count >= 1 {
		_, err = dbmap.Exec("delete from usersLikeIT where user_id = $1 and where thought_id = $2", userId, thoughtId)
		if err != nil {
			log.Printf("cannnot delete usersLikeIT: %s", err.Error())
		}
		err = dbmap.SelectOne(&thought, "select * from thought where thought_id = $1", thoughtId)
		if err != nil {
			log.Printf("cannnot selectOne thought: %s", err.Error())
		}
		thought.LikesCount -= 1
		_, err = dbmap.Update(&thought)
		if err != nil {
			log.Printf("cannnot update thought: %s", err.Error())
		}
	} else {
		err = dbmap.Insert(usersLikeIT)
		if err != nil {
			log.Printf("cannnot insert usersLikeIT: %s", err.Error())
		}
		err = dbmap.SelectOne(&thought, "select * from thought where thought_id = $1", thoughtId)
		if err != nil {
			log.Printf("cannnot selectOne thought: %s", err.Error())
		}
		thought.LikesCount += 1
		_, err = dbmap.Update(&thought)
		if err != nil {
			log.Printf("cannnot update thought: %s", err.Error())
		}
	}
}


