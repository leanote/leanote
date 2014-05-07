package service

import (
	"github.com/leanote/leanote/app/info"
	"github.com/leanote/leanote/app/db"
	. "github.com/leanote/leanote/app/lea"
	"labix.org/v2/mgo/bson"
	"time"
	"strings"
	"fmt"
)

type UserService struct {
}


// 添加用户
func (this *UserService) AddUser(user info.User) bool {
	if user.UserId == "" {
		user.UserId = bson.NewObjectId()
	}
	user.CreatedTime = time.Now()
	
	if user.Email != "" {
		user.Email = strings.ToLower(user.Email)
		
		// 发送验证邮箱
		go func() {
			this.RegisterSendActiveEmail(user.UserId.Hex(), user.Email)
		}();
	}
	
	return db.Insert(db.Users, user)
}

// 通过email得到userId
func (this *UserService) GetUserId(email string) string {
	email = strings.ToLower(email)
	user := info.User{}
	db.GetByQ(db.Users, bson.M{"Email": email}, &user)
	return user.UserId.Hex()
}

// 是否存在该用户 email
func (this *UserService) IsExistsUser(email string) bool {
	if this.GetUserId(email) == "" {
		return false
	}
	return true
}

// 是否存在该用户 username
func (this *UserService) IsExistsUserByUsername(username string) bool {
	return db.Count(db.Users, bson.M{"Username": username}) >= 1
}

// 得到用户信息, userId, username, email
func (this *UserService) GetUserInfoByAny(idEmailUsername string) info.User {
	if IsObjectId(idEmailUsername) {
		return this.GetUserInfo(idEmailUsername)
	}
	
	if strings.Contains(idEmailUsername, "@") {
		return this.GetUserInfoByEmail(idEmailUsername)
	}
	
	return this.GetUserInfoByUsername(idEmailUsername)
}

// 得到用户信息 userId
func (this *UserService) GetUserInfo(userId string) info.User {
	user := info.User{}
	db.Get(db.Users, userId, &user)
	return user
}
// 得到用户信息 email
func (this *UserService) GetUserInfoByEmail(email string) info.User {
	user := info.User{}
	db.GetByQ(db.Users, bson.M{"Email": email}, &user)
	return user
}
// 得到用户信息 username
func (this *UserService) GetUserInfoByUsername(username string) info.User {
	user := info.User{}
	db.GetByQ(db.Users, bson.M{"Username": username}, &user)
	return user
}

func (this *UserService) GetUserInfoByThirdUserId(thirdUserId string) info.User {
	user := info.User{}
	db.GetByQ(db.Users, bson.M{"ThirdUserId": thirdUserId}, &user)
	return user
}
func (this *UserService) ListUserInfosByUserIds(userIds []bson.ObjectId) []info.User {
	users := []info.User{}
	db.ListByQ(db.Users, bson.M{"_id": bson.M{"$in": userIds}}, &users)
	return users
}

// 通过ids得到users, 按id的顺序组织users
func (this *UserService) GetUserInfosOrderBySeq(userIds []bson.ObjectId) []info.User {
	users := []info.User{}
	db.ListByQ(db.Users, bson.M{"_id": bson.M{"$in": userIds}}, &users);
	
	usersMap := map[bson.ObjectId]info.User{}
	for _, user := range users {
		usersMap[user.UserId] = user
	}
	
	users2 := []info.User{};
	for _, userId := range userIds {
		if user, ok := usersMap[userId]; ok {
			users2 = append(users2, user)
		}
	}
	return users2
}

// 使用email(username), pwd得到用户信息
func (this *UserService) LoginGetUserInfo(emailOrUsername, md5Pwd string) info.User {
	emailOrUsername = strings.ToLower(emailOrUsername)
	
	user := info.User{}
	if strings.Contains(emailOrUsername, "@") {
		db.GetByQ(db.Users, bson.M{"Email": emailOrUsername, "Pwd": md5Pwd}, &user)
	} else {
		db.GetByQ(db.Users, bson.M{"Username": emailOrUsername, "Pwd": md5Pwd}, &user)
	}
	
	return user
}

// 更新username
func (this *UserService) UpdateUsername(userId, username string) (bool, string) {
	if userId == "" || username == "" {
		return false, "用户已存在"
	}
	usernameRaw := username // 原先的, 可能是同一个, 但有大小写
	username = strings.ToLower(username)
	
	// 先判断是否存在
	userIdO := bson.ObjectIdHex(userId)
	if db.Has(db.Users, bson.M{"Username": username, "_id": bson.M{"$ne": userIdO}}) {
		return false, "用户已存在"
	}
		
	ok := db.UpdateByQMap(db.Users, bson.M{"_id": userIdO}, bson.M{"Username": username, "UsernameRaw": usernameRaw})
	return ok, ""
}

//----------------------
// 已经登录了的用户修改密码
func (this *UserService) UpdatePwd(userId, oldPwd, pwd string) (bool, string) {
	userInfo := this.GetUserInfo(userId)
	if userInfo.Pwd != Md5(oldPwd) {
		return false, "旧密码错误"
	}
	ok := db.UpdateByQField(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, "Pwd", Md5(pwd))
	return ok, ""
}

// 修改主题
func (this *UserService) UpdateTheme(userId, theme string) (bool) {
	ok := db.UpdateByQField(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, "Theme", theme)
	return ok
}

//---------------
// 修改email

// 发送激活邮件

// AddUser调用
// 可以使用一个goroutine
func (this *UserService) RegisterSendActiveEmail(userId string, email string) bool {
	token := tokenService.NewToken(userId, email, info.TokenActiveEmail)
	
	if token == "" {
		return false
	}
	
	// 发送邮件
	url := "http://leanote.com/user/activeEmail?token=" + token
	body := fmt.Sprintf("请点击链接验证邮箱: <a href='%v'>%v</a>. %v小时后过期.", url, url, tokenService.GetOverHours(info.TokenActiveEmail));
	if !SendEmail(email, "leanote-验证邮箱", "验证邮箱", body) {
		return false
	}
	
	// 发送给我 life@leanote.com
	SendEmail("life@leanote.com", "新增用户", "新增用户", "用户名" + email);
	
	return true
}

// 修改邮箱
func (this *UserService) UpdateEmailSendActiveEmail(userId, email string) (ok bool, msg string) {
	// 先验证该email是否被注册了
	if userService.IsExistsUser(email) {
		ok = false
		msg = "该邮箱已注册"
		return
	}

	token := tokenService.NewToken(userId, email, info.TokenUpdateEmail)
	
	if token == "" {
		return
	}
	
	// 发送邮件
	url := "http://115.28.133.226/user/updateEmail?token=" + token
	body := "邮箱验证后您的登录邮箱为: <b>" + email + "</b><br />";
	body += fmt.Sprintf("请点击链接验证邮箱: <a href='%v'>%v</a>. %v小时后过期.", url, url, tokenService.GetOverHours(info.TokenUpdateEmail));
	if !SendEmail(email, "leanote-验证邮箱", "验证邮箱", body) {
		msg = "发送失败"
		return 
	}
	ok = true
	return
}

// 注册后验证邮箱
func (this *UserService) ActiveEmail(token string) (ok bool, msg, email string) {
	tokenInfo := info.Token{}
	if ok, msg, tokenInfo = tokenService.VerifyToken(token, info.TokenActiveEmail); ok {
		// 修改之后的邮箱
		email = tokenInfo.Email
		userInfo := this.GetUserInfoByEmail(email)
		if userInfo.UserId == "" {
			ok = false
			msg = "不存在该用户"
			return 
		}
		
		// 修改之, 并将verified = true
		ok = db.UpdateByQMap(db.Users, bson.M{"_id": userInfo.UserId}, bson.M{"Verified": true})
		return
	}
	
	ok = false
	msg = "该链接已过期"
	return
}

// 修改邮箱
// 在此之前, 验证token是否过期
// 验证email是否有人注册了
func (this *UserService) UpdateEmail(token string) (ok bool, msg, email string) {
	tokenInfo := info.Token{}
	if ok, msg, tokenInfo = tokenService.VerifyToken(token, info.TokenUpdateEmail); ok {
		// 修改之后的邮箱
		email = tokenInfo.Email
		// 先验证该email是否被注册了
		if userService.IsExistsUser(email) {
			ok = false
			msg = "该邮箱已注册"
			return
		}
		
		// 修改之, 并将verified = true
		ok = db.UpdateByQMap(db.Users, bson.M{"_id": tokenInfo.UserId}, bson.M{"Email": email, "Verified": true})
		return
	}
	
	ok = false
	msg = "该链接已过期"
	return
}

//---------
// 第三方添加账号
func (this *UserService) ThirdAddUser(userId, email, pwd string) (ok bool, msg string) {
	// 判断email是否存在
	if this.IsExistsUser(email) {
		msg = "该用户已存在"
		return
	}
	
	ok = db.UpdateByQMap(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, bson.M{"Email": email, "Pwd": Md5(pwd)})
	return
}


//------------
// 偏好设置

// 宽度
func (this *UserService)UpdateColumnWidth(userId string, notebookWidth, noteListWidth int) bool {
	return db.UpdateByQMap(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, bson.M{"NotebookWidth": notebookWidth, "NoteListWidth": noteListWidth})
}
// 左侧是否隐藏
func  (this *UserService)UpdateLeftIsMin(userId string, leftIsMin bool) bool {
	return db.UpdateByQMap(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, bson.M{"LeftIsMin": leftIsMin})
}
