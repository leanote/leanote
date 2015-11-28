package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type UserService struct {
}

// 自增Usn
// 每次notebook,note添加, 修改, 删除, 都要修改
func (this *UserService) IncrUsn(userId string) int {
	user := info.User{}
	query := bson.M{"_id": bson.ObjectIdHex(userId)}
	db.GetByQWithFields(db.Users, query, []string{"Usn"}, &user)
	usn := user.Usn
	usn += 1
	Log("inc Usn")
	db.UpdateByQField(db.Users, query, "Usn", usn)
	return usn
	//	return db.Update(db.Notes, bson.M{"_id": bson.ObjectIdHex(noteId)}, bson.M{"$inc": bson.M{"ReadNum": 1}})
}

func (this *UserService) GetUsn(userId string) int {
	user := info.User{}
	query := bson.M{"_id": bson.ObjectIdHex(userId)}
	db.GetByQWithFields(db.Users, query, []string{"Usn"}, &user)
	return user.Usn
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
			emailService.RegisterSendActiveEmail(user, user.Email)
			// 发送给我 life@leanote.com
			// emailService.SendEmail("life@leanote.com", "新增用户", "{header}用户名"+user.Email+"{footer}")
		}()
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

// 得到用户名
func (this *UserService) GetUsername(userId string) string {
	user := info.User{}
	db.GetByQWithFields(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, []string{"Username"}, &user)
	return user.Username
}

// 得到用户名
func (this *UserService) GetUsernameById(userId bson.ObjectId) string {
	user := info.User{}
	db.GetByQWithFields(db.Users, bson.M{"_id": userId}, []string{"Username"}, &user)
	return user.Username
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

	// username
	return this.GetUserInfoByUsername(idEmailUsername)
}

func (this *UserService) setUserLogo(user *info.User) {
	// Logo路径问题, 有些有http: 有些没有
	if user.Logo == "" {
		user.Logo = "images/blog/default_avatar.png"
	}
	if user.Logo != "" && !strings.HasPrefix(user.Logo, "http") {
		user.Logo = strings.Trim(user.Logo, "/")
		user.Logo = configService.GetSiteUrl() + "/" + user.Logo
	}
}

// 仅得到用户
func (this *UserService) GetUser(userId string) info.User {
	user := info.User{}
	db.Get(db.Users, userId, &user)
	return user
}

// 得到用户信息 userId
func (this *UserService) GetUserInfo(userId string) info.User {
	user := info.User{}
	db.Get(db.Users, userId, &user)
	// Logo路径问题, 有些有http: 有些没有
	this.setUserLogo(&user)
	return user
}

// 得到用户信息 email
func (this *UserService) GetUserInfoByEmail(email string) info.User {
	user := info.User{}
	db.GetByQ(db.Users, bson.M{"Email": email}, &user)
	// Logo路径问题, 有些有http: 有些没有
	this.setUserLogo(&user)
	return user
}

// 得到用户信息 username
func (this *UserService) GetUserInfoByUsername(username string) info.User {
	user := info.User{}
	username = strings.ToLower(username)
	db.GetByQ(db.Users, bson.M{"Username": username}, &user)
	// Logo路径问题, 有些有http: 有些没有
	this.setUserLogo(&user)
	return user
}

func (this *UserService) GetUserInfoByThirdUserId(thirdUserId string) info.User {
	user := info.User{}
	db.GetByQ(db.Users, bson.M{"ThirdUserId": thirdUserId}, &user)
	this.setUserLogo(&user)
	return user
}
func (this *UserService) ListUserInfosByUserIds(userIds []bson.ObjectId) []info.User {
	users := []info.User{}
	db.ListByQ(db.Users, bson.M{"_id": bson.M{"$in": userIds}}, &users)
	return users
}
func (this *UserService) ListUserInfosByEmails(emails []string) []info.User {
	users := []info.User{}
	db.ListByQ(db.Users, bson.M{"Email": bson.M{"$in": emails}}, &users)
	return users
}

// 用户信息即可
func (this *UserService) MapUserInfoByUserIds(userIds []bson.ObjectId) map[bson.ObjectId]info.User {
	users := []info.User{}
	db.ListByQ(db.Users, bson.M{"_id": bson.M{"$in": userIds}}, &users)

	userMap := make(map[bson.ObjectId]info.User, len(users))
	for _, user := range users {
		this.setUserLogo(&user)
		userMap[user.UserId] = user
	}
	return userMap
}

// 用户信息和博客设置信息
func (this *UserService) MapUserInfoAndBlogInfosByUserIds(userIds []bson.ObjectId) map[bson.ObjectId]info.User {
	return this.MapUserInfoByUserIds(userIds)
}

// 返回info.UserAndBlog
func (this *UserService) MapUserAndBlogByUserIds(userIds []bson.ObjectId) map[string]info.UserAndBlog {
	users := []info.User{}
	db.ListByQ(db.Users, bson.M{"_id": bson.M{"$in": userIds}}, &users)

	userBlogs := []info.UserBlog{}
	db.ListByQ(db.UserBlogs, bson.M{"_id": bson.M{"$in": userIds}}, &userBlogs)

	userBlogMap := make(map[bson.ObjectId]info.UserBlog, len(userBlogs))
	for _, user := range userBlogs {
		userBlogMap[user.UserId] = user
	}

	userAndBlogMap := make(map[string]info.UserAndBlog, len(users))

	for _, user := range users {
		this.setUserLogo(&user)

		userBlog, ok := userBlogMap[user.UserId]
		if !ok {
			continue
		}

		userAndBlogMap[user.UserId.Hex()] = info.UserAndBlog{
			UserId:    user.UserId,
			Username:  user.Username,
			Email:     user.Email,
			Logo:      user.Logo,
			BlogTitle: userBlog.Title,
			BlogLogo:  userBlog.Logo,
			BlogUrl:   blogService.GetUserBlogUrl(&userBlog, user.Username),
		}
	}
	return userAndBlogMap
}

// 得到用户信息+博客主页
func (this *UserService) GetUserAndBlogUrl(userId string) info.UserAndBlogUrl {
	user := this.GetUserInfo(userId)
	userBlog := blogService.GetUserBlog(userId)

	blogUrls := blogService.GetBlogUrls(&userBlog, &user)

	return info.UserAndBlogUrl{
		User:    user,
		BlogUrl: blogUrls.IndexUrl,
		PostUrl: blogUrls.PostUrl,
	}
}

// 得到userAndBlog公开信息
func (this *UserService) GetUserAndBlog(userId string) info.UserAndBlog {
	user := this.GetUserInfo(userId)
	userBlog := blogService.GetUserBlog(userId)
	return info.UserAndBlog{
		UserId:    user.UserId,
		Username:  user.Username,
		Email:     user.Email,
		Logo:      user.Logo,
		BlogTitle: userBlog.Title,
		BlogLogo:  userBlog.Logo,
		BlogUrl:   blogService.GetUserBlogUrl(&userBlog, user.Username),
		BlogUrls:  blogService.GetBlogUrls(&userBlog, &user),
	}
}

// 通过ids得到users, 按id的顺序组织users
func (this *UserService) GetUserInfosOrderBySeq(userIds []bson.ObjectId) []info.User {
	users := []info.User{}
	db.ListByQ(db.Users, bson.M{"_id": bson.M{"$in": userIds}}, &users)

	usersMap := map[bson.ObjectId]info.User{}
	for _, user := range users {
		usersMap[user.UserId] = user
	}

	hasAppend := map[bson.ObjectId]bool{} // 为了防止userIds有重复的
	users2 := []info.User{}
	for _, userId := range userIds {
		if user, ok := usersMap[userId]; ok && !hasAppend[userId] {
			hasAppend[userId] = true
			users2 = append(users2, user)
		}
	}
	return users2
}

// 使用email(username), 得到用户信息
func (this *UserService) GetUserInfoByName(emailOrUsername string) info.User {
	emailOrUsername = strings.ToLower(emailOrUsername)

	user := info.User{}
	if strings.Contains(emailOrUsername, "@") {
		db.GetByQ(db.Users, bson.M{"Email": emailOrUsername}, &user)
	} else {
		db.GetByQ(db.Users, bson.M{"Username": emailOrUsername}, &user)
	}
	this.setUserLogo(&user)
	return user
}

// 更新username
func (this *UserService) UpdateUsername(userId, username string) (bool, string) {
	if userId == "" || username == "" || username == "admin" { // admin用户是内置的, 不能设置
		return false, "usernameIsExisted"
	}
	usernameRaw := username // 原先的, 可能是同一个, 但有大小写
	username = strings.ToLower(username)

	// 先判断是否存在
	userIdO := bson.ObjectIdHex(userId)
	if db.Has(db.Users, bson.M{"Username": username, "_id": bson.M{"$ne": userIdO}}) {
		return false, "usernameIsExisted"
	}

	ok := db.UpdateByQMap(db.Users, bson.M{"_id": userIdO}, bson.M{"Username": username, "UsernameRaw": usernameRaw})
	return ok, ""
}

// 修改头像
func (this *UserService) UpdateAvatar(userId, avatarPath string) bool {
	userIdO := bson.ObjectIdHex(userId)
	return db.UpdateByQField(db.Users, bson.M{"_id": userIdO}, "Logo", avatarPath)
}

//----------------------
// 已经登录了的用户修改密码
func (this *UserService) UpdatePwd(userId, oldPwd, pwd string) (bool, string) {
	userInfo := this.GetUserInfo(userId)
	if !ComparePwd(oldPwd, userInfo.Pwd) {
		return false, "oldPasswordError"
	}

	passwd := GenPwd(pwd)
	if passwd == "" {
		return false, "GenerateHash error"
	}

	ok := db.UpdateByQField(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, "Pwd", passwd)
	return ok, ""
}

// 管理员重置密码
func (this *UserService) ResetPwd(adminUserId, userId, pwd string) (ok bool, msg string) {
	if configService.GetAdminUserId() != adminUserId {
		return
	}

	passwd := GenPwd(pwd)
	if passwd == "" {
		return false, "GenerateHash error"
	}
	ok = db.UpdateByQField(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, "Pwd", passwd)
	return
}

// 修改主题
func (this *UserService) UpdateTheme(userId, theme string) bool {
	ok := db.UpdateByQField(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, "Theme", theme)
	return ok
}

// 帐户类型设置
func (this *UserService) UpdateAccount(userId, accountType string, accountStartTime, accountEndTime time.Time,
	maxImageNum, maxImageSize, maxAttachNum, maxAttachSize, maxPerAttachSize int) bool {
	return db.UpdateByQI(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, info.UserAccount{
		AccountType:      accountType,
		AccountStartTime: accountStartTime,
		AccountEndTime:   accountEndTime,
		MaxImageNum:      maxImageNum,
		MaxImageSize:     maxImageSize,
		MaxAttachNum:     maxAttachNum,
		MaxAttachSize:    maxAttachSize,
		MaxPerAttachSize: maxPerAttachSize,
	})
}

//---------------
// 修改email

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
		email = strings.ToLower(tokenInfo.Email)
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

//------------
// 偏好设置

// 宽度
func (this *UserService) UpdateColumnWidth(userId string, notebookWidth, noteListWidth, mdEditorWidth int) bool {
	return db.UpdateByQMap(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)},
		bson.M{"NotebookWidth": notebookWidth, "NoteListWidth": noteListWidth, "MdEditorWidth": mdEditorWidth})
}

// 左侧是否隐藏
func (this *UserService) UpdateLeftIsMin(userId string, leftIsMin bool) bool {
	return db.UpdateByQMap(db.Users, bson.M{"_id": bson.ObjectIdHex(userId)}, bson.M{"LeftIsMin": leftIsMin})
}

//-------------
// user admin
func (this *UserService) ListUsers(pageNumber, pageSize int, sortField string, isAsc bool, email string) (page info.Page, users []info.User) {
	users = []info.User{}
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)
	query := bson.M{}
	if email != "" {
		orQ := []bson.M{
			bson.M{"Email": bson.M{"$regex": bson.RegEx{".*?" + email + ".*", "i"}}},
			bson.M{"Username": bson.M{"$regex": bson.RegEx{".*?" + email + ".*", "i"}}},
		}
		query["$or"] = orQ
	}
	q := db.Users.Find(query)
	// 总记录数
	count, _ := q.Count()
	// 列表
	q.Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&users)
	page = info.NewPage(pageNumber, pageSize, count, nil)
	return
}

func (this *UserService) GetAllUserByFilter(userFilterEmail, userFilterWhiteList, userFilterBlackList string, verified bool) []info.User {
	query := bson.M{}

	if verified {
		query["Verified"] = true
	}

	orQ := []bson.M{}
	if userFilterEmail != "" {
		orQ = append(orQ, bson.M{"Email": bson.M{"$regex": bson.RegEx{".*?" + userFilterEmail + ".*", "i"}}},
			bson.M{"Username": bson.M{"$regex": bson.RegEx{".*?" + userFilterEmail + ".*", "i"}}},
		)
	}
	if userFilterWhiteList != "" {
		userFilterWhiteList = strings.Replace(userFilterWhiteList, "\r", "", -1)
		emails := strings.Split(userFilterWhiteList, "\n")
		orQ = append(orQ, bson.M{"Email": bson.M{"$in": emails}})
	}
	if len(orQ) > 0 {
		query["$or"] = orQ
	}

	emailQ := bson.M{}
	if userFilterBlackList != "" {
		userFilterWhiteList = strings.Replace(userFilterBlackList, "\r", "", -1)
		bEmails := strings.Split(userFilterBlackList, "\n")
		emailQ["$nin"] = bEmails
		query["Email"] = emailQ
	}

	LogJ(query)
	users := []info.User{}
	q := db.Users.Find(query)
	q.All(&users)
	Log(len(users))

	return users
}

// 统计
func (this *UserService) CountUser() int {
	return db.Count(db.Users, bson.M{})
}
