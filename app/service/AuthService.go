package service

import (
	"gopkg.in/mgo.v2/bson"
	//	"github.com/leanote/leanote/app/db"
	"fmt"
	"github.com/revel/revel"
	"leanote/app/info"
	. "leanote/app/lea"
)

// 登录与权限

type AuthService struct {
}

// pwd已md5了
func (this *AuthService) Login(emailOrUsername, pwd string) info.User {
	return userService.LoginGetUserInfo(emailOrUsername, Md5(pwd))
}

// 注册
/*
注册 leanote@leanote.com userId = "5368c1aa99c37b029d000001"
添加 在博客上添加一篇欢迎note, note1 5368c1b919807a6f95000000

将nk1(只读), nk2(可写) 分享给该用户
将note1 复制到用户的生活nk上
*/
// 1. 添加用户
// 2. 将leanote共享给我
// [ok]
func (this *AuthService) Register(email, pwd string) (bool, string) {
	// 用户是否已存在
	if userService.IsExistsUser(email) {
		return false, email + " 已被注册"
	}
	user := info.User{UserId: bson.NewObjectId(), Email: email, Username: email, Pwd: Md5(pwd)}
	return this.register(user)
}

func (this *AuthService) register(user info.User) (bool, string) {
	if userService.AddUser(user) {
		// 添加笔记本, 生活, 学习, 工作
		notebook := info.Notebook{
			Seq:    -1,
			UserId: user.UserId}
		title2Id := map[string]bson.ObjectId{"life": bson.NewObjectId(), "study": bson.NewObjectId(), "work": bson.NewObjectId()}
		for title, objectId := range title2Id {
			notebook.Title = title
			notebook.NotebookId = objectId
			notebook.UserId = user.UserId
			notebookService.AddNotebook(notebook)
		}

		email := user.Email

		// 添加leanote -> 该用户的共享
		leanoteUserId, _ := revel.Config.String("register.sharedUserId")    // "5368c1aa99c37b029d000001";
		nk1, _ := revel.Config.String("register.sharedUserShareNotebookId") // 5368c1aa99c37b029d000002" // leanote
		welcomeNoteId, _ := revel.Config.String("register.welcomeNoteId")   // "5368c1b919807a6f95000000" // 欢迎来到leanote

		if leanoteUserId != "" && nk1 != "" && welcomeNoteId != "" {
			shareService.AddShareNotebook(nk1, 0, leanoteUserId, email)
			shareService.AddShareNote(welcomeNoteId, 0, leanoteUserId, email)

			// 将welcome copy给我
			note := noteService.CopySharedNote(welcomeNoteId, title2Id["life"].Hex(), leanoteUserId, user.UserId.Hex())

			// 公开为博客
			noteUpdate := bson.M{"IsBlog": true}
			noteService.UpdateNote(user.UserId.Hex(), user.UserId.Hex(), note.NoteId.Hex(), noteUpdate)
		}

		//---------------
		// 添加一条userBlog
		blogService.UpdateUserBlog(info.UserBlog{UserId: user.UserId,
			Title:    user.Username + " 's Blog",
			SubTitle: "love leanote!",
			AboutMe:  "Hello, I am (^_^)",
		})
	}

	return true, ""
}

//--------------
// 第三方注册

// 第三方得到用户名, 可能需要多次判断
func (this *AuthService) getUsername(thirdType, thirdUsername string) (username string) {
	username = thirdType + "-" + thirdUsername
	i := 1
	for {
		if !userService.IsExistsUserByUsername(username) {
			return
		}
		username = fmt.Sprintf("%v%v", username, i)
	}
}

func (this *AuthService) ThirdRegister(thirdType, thirdUserId, thirdUsername string) (exists bool, userInfo info.User) {
	userInfo = userService.GetUserInfoByThirdUserId(thirdUserId)
	if userInfo.UserId != "" {
		exists = true
		return
	}

	username := this.getUsername(thirdType, thirdUsername)
	userInfo = info.User{UserId: bson.NewObjectId(),
		Username:      username,
		ThirdUserId:   thirdUserId,
		ThirdUsername: thirdUsername,
	}
	_, _ = this.register(userInfo)
	return
}
