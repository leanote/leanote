package service

import (
	"gopkg.in/mgo.v2/bson"
//	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
//	"github.com/revel/revel"
	"strings"
	. "github.com/leanote/leanote/app/lea"
	"fmt"
	"strconv"
	"errors"
)

// 登录与权限

type AuthService struct {
}

// 使用bcrypt认证或者Md5认证
func (this *AuthService) Login(emailOrUsername, pwd string) (info.User, error) {
	emailOrUsername = strings.Trim(emailOrUsername, " ")
  //	pwd = strings.Trim(pwd, " ")
	userInfo := userService.GetUserInfoByName(emailOrUsername)
  passwd := userInfo.Pwd
  if len(passwd) == 32 && Md5(pwd) != passwd {
  	return userInfo,  errors.New("wrong username or password")
  } 
  if len(passwd) > 32 {
    hex := []byte(passwd)
	  if !CompareHash(hex, pwd) {
		  return userInfo,  errors.New("wrong username or password")
	  } 
  } 
	return userInfo, nil
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
func (this *AuthService) Register(email, pwd, fromUserId string) (bool, string) {
	// 用户是否已存在
	if userService.IsExistsUser(email) {
		return false, "userHasBeenRegistered-" + email
	}
  digest, err := GenerateHash(pwd)
	if err != nil {
		return false,"GenerateHash error"
	}
  passwd := string(digest)
	user := info.User{UserId: bson.NewObjectId(), Email: email, Username: email, Pwd: passwd}
	if fromUserId != "" && IsObjectId(fromUserId) {
		user.FromUserId = bson.ObjectIdHex(fromUserId)
	}
	LogJ(user)
	return this.register(user)
}

func (this *AuthService) register(user info.User) (bool, string) {
	if userService.AddUser(user) {
		// 添加笔记本, 生活, 学习, 工作
		userId := user.UserId.Hex();
		notebook := info.Notebook{
			Seq: -1,
			UserId: user.UserId}
		title2Id := map[string]bson.ObjectId{"life": bson.NewObjectId(), "study": bson.NewObjectId(), "work": bson.NewObjectId()}
		for title, objectId := range title2Id {
			notebook.Title = title
			notebook.NotebookId = objectId
			notebook.UserId = user.UserId
			notebookService.AddNotebook(notebook);
		}
		
		// 添加leanote -> 该用户的共享
		registerSharedUserId := configService.GetGlobalStringConfig("registerSharedUserId")
		if(registerSharedUserId != "") {
			registerSharedNotebooks := configService.GetGlobalArrMapConfig("registerSharedNotebooks")
			registerSharedNotes := configService.GetGlobalArrMapConfig("registerSharedNotes")
			registerCopyNoteIds := configService.GetGlobalArrayConfig("registerCopyNoteIds")
			
			// 添加共享笔记本
			for _, notebook := range registerSharedNotebooks {
				perm, _ := strconv.Atoi(notebook["perm"])
				shareService.AddShareNotebookToUserId(notebook["notebookId"], perm, registerSharedUserId, userId);
			}
			
			// 添加共享笔记
			for _, note := range registerSharedNotes {
				perm, _ := strconv.Atoi(note["perm"])
				shareService.AddShareNoteToUserId(note["noteId"], perm, registerSharedUserId, userId);
			}
			
			// 复制笔记
			for _, noteId := range registerCopyNoteIds {
				note := noteService.CopySharedNote(noteId, title2Id["life"].Hex(), registerSharedUserId, user.UserId.Hex());
//				Log(noteId)
//				Log("Copy")
//				LogJ(note)
				noteUpdate := bson.M{"IsBlog": false} // 不要是博客
				noteService.UpdateNote(user.UserId.Hex(), note.NoteId.Hex(), noteUpdate, -1)
			}
		}
		
		//---------------
		// 添加一条userBlog
		blogService.UpdateUserBlog(info.UserBlog{UserId: user.UserId, 
			Title: user.Username + " 's Blog", 
			SubTitle: "Love Leanote!",
			AboutMe: "Hello, I am (^_^)",
			CanComment: true,
			})
		// 添加一个单页面
		blogService.AddOrUpdateSingle(user.UserId.Hex(), "", "About Me", "Hello, I am (^_^)")
	}
	
	return true, ""
}

//--------------
// 第三方注册

// 第三方得到用户名, 可能需要多次判断
func (this *AuthService) getUsername(thirdType, thirdUsername string) (username string) {
	username = thirdType + "-" + thirdUsername
	i := 1
	for ;; {
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
		Username: username, 
		ThirdUserId: thirdUserId,
		ThirdUsername: thirdUsername,
		}
	_, _ = this.register(userInfo)
	return
}
