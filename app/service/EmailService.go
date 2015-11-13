package service

import (
	"bytes"
	"fmt"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// 发送邮件

type EmailService struct {
	tpls map[string]*template.Template
}

func NewEmailService() *EmailService {
	return &EmailService{tpls: map[string]*template.Template{}}
}

// 发送邮件
var host = ""
var emailPort = ""
var username = ""
var password = ""

func InitEmailFromDb() {
	host = configService.GetGlobalStringConfig("emailHost")
	emailPort = configService.GetGlobalStringConfig("emailPort")
	username = configService.GetGlobalStringConfig("emailUsername")
	password = configService.GetGlobalStringConfig("emailPassword")
}

func (this *EmailService) SendEmail(to, subject, body string) (ok bool, e string) {
	InitEmailFromDb()

	if host == "" || emailPort == "" || username == "" || password == "" {
		return
	}
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", username, password, hp[0])

	var content_type string

	mailtype := "html"
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + username + "<" + username + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host+":"+emailPort, auth, username, send_to, msg)

	if err != nil {
		e = fmt.Sprint(err)
		return
	}
	ok = true
	return
}

// AddUser调用
// 可以使用一个goroutine
func (this *EmailService) RegisterSendActiveEmail(userInfo info.User, email string) bool {
	token := tokenService.NewToken(userInfo.UserId.Hex(), email, info.TokenActiveEmail)
	if token == "" {
		return false
	}

	subject := configService.GetGlobalStringConfig("emailTemplateRegisterSubject")
	tpl := configService.GetGlobalStringConfig("emailTemplateRegister")

	if tpl == "" {
		return false
	}

	tokenUrl := configService.GetSiteUrl() + "/user/activeEmail?token=" + token
	// {siteUrl} {tokenUrl} {token} {tokenTimeout} {user.id} {user.email} {user.username}
	token2Value := map[string]interface{}{"siteUrl": configService.GetSiteUrl(), "tokenUrl": tokenUrl, "token": token, "tokenTimeout": strconv.Itoa(int(tokenService.GetOverHours(info.TokenActiveEmail))),
		"user": map[string]interface{}{
			"userId":   userInfo.UserId.Hex(),
			"email":    userInfo.Email,
			"username": userInfo.Username,
		},
	}

	var ok bool
	ok, _, subject, tpl = this.renderEmail(subject, tpl, token2Value)
	if !ok {
		return false
	}

	// 发送邮件
	ok, _ = this.SendEmail(email, subject, tpl)
	return ok
}

// 修改邮箱
func (this *EmailService) UpdateEmailSendActiveEmail(userInfo info.User, email string) (ok bool, msg string) {
	// 先验证该email是否被注册了
	if userService.IsExistsUser(email) {
		ok = false
		msg = "该邮箱已注册"
		return
	}

	token := tokenService.NewToken(userInfo.UserId.Hex(), email, info.TokenUpdateEmail)

	if token == "" {
		return
	}

	subject := configService.GetGlobalStringConfig("emailTemplateUpdateEmailSubject")
	tpl := configService.GetGlobalStringConfig("emailTemplateUpdateEmail")

	// 发送邮件
	tokenUrl := configService.GetSiteUrl() + "/user/updateEmail?token=" + token
	// {siteUrl} {tokenUrl} {token} {tokenTimeout} {user.userId} {user.email} {user.username}
	token2Value := map[string]interface{}{"siteUrl": configService.GetSiteUrl(), "tokenUrl": tokenUrl, "token": token, "tokenTimeout": strconv.Itoa(int(tokenService.GetOverHours(info.TokenActiveEmail))),
		"newEmail": email,
		"user": map[string]interface{}{
			"userId":   userInfo.UserId.Hex(),
			"email":    userInfo.Email,
			"username": userInfo.Username,
		},
	}

	ok, msg, subject, tpl = this.renderEmail(subject, tpl, token2Value)
	if !ok {
		return
	}

	// 发送邮件
	ok, msg = this.SendEmail(email, subject, tpl)
	return
}

func (this *EmailService) FindPwdSendEmail(token, email string) (ok bool, msg string) {
	subject := configService.GetGlobalStringConfig("emailTemplateFindPasswordSubject")
	tpl := configService.GetGlobalStringConfig("emailTemplateFindPassword")

	// 发送邮件
	tokenUrl := configService.GetSiteUrl() + "/findPassword/" + token
	// {siteUrl} {tokenUrl} {token} {tokenTimeout} {user.id} {user.email} {user.username}
	token2Value := map[string]interface{}{"siteUrl": configService.GetSiteUrl(), "tokenUrl": tokenUrl,
		"token": token, "tokenTimeout": strconv.Itoa(int(tokenService.GetOverHours(info.TokenActiveEmail)))}

	ok, msg, subject, tpl = this.renderEmail(subject, tpl, token2Value)
	if !ok {
		return
	}
	// 发送邮件
	ok, msg = this.SendEmail(email, subject, tpl)
	return
}

// 发送邀请链接
func (this *EmailService) SendInviteEmail(userInfo info.User, email, content string) bool {
	subject := configService.GetGlobalStringConfig("emailTemplateInviteSubject")
	tpl := configService.GetGlobalStringConfig("emailTemplateInvite")

	token2Value := map[string]interface{}{"siteUrl": configService.GetSiteUrl(),
		"registerUrl": configService.GetSiteUrl() + "/register?from=" + userInfo.Username,
		"content":     content,
		"user": map[string]interface{}{
			"username": userInfo.Username,
			"email":    userInfo.Email,
		},
	}
	var ok bool
	ok, _, subject, tpl = this.renderEmail(subject, tpl, token2Value)
	if !ok {
		return false
	}
	// 发送邮件
	ok, _ = this.SendEmail(email, subject, tpl)
	return ok
}

// 发送评论
func (this *EmailService) SendCommentEmail(note info.Note, comment info.BlogComment, userId, content string) bool {
	subject := configService.GetGlobalStringConfig("emailTemplateCommentSubject")
	tpl := configService.GetGlobalStringConfig("emailTemplateComment")

	// title := "评论提醒"

	/*
		toUserId := note.UserId.Hex()
		// title := "评论提醒"

		// 表示回复回复的内容, 那么发送给之前回复的
		if comment.CommentId != "" {
			toUserId = comment.UserId.Hex()
		}
		toUserInfo := userService.GetUserInfo(toUserId)
		sendUserInfo := userService.GetUserInfo(userId)

		subject := note.Title + " 收到 " + sendUserInfo.Username + " 的评论";
		if comment.CommentId != "" {
			subject = "您在 " + note.Title + " 发表的评论收到 " + sendUserInfo.Username;
			if userId == note.UserId.Hex() {
				subject += "(作者)";
			}
			subject += " 的评论";
		}
	*/

	toUserId := note.UserId.Hex()
	// 表示回复回复的内容, 那么发送给之前回复的
	if comment.CommentId != "" {
		toUserId = comment.UserId.Hex()
	}
	toUserInfo := userService.GetUserInfo(toUserId) // 被评论者
	sendUserInfo := userService.GetUserInfo(userId) // 评论者

	// {siteUrl} {blogUrl}
	// {blog.id} {blog.title} {blog.url}
	// {commentUser.userId} {commentUser.username} {commentUser.email}
	// {commentedUser.userId} {commentedUser.username} {commentedUser.email}
	token2Value := map[string]interface{}{"siteUrl": configService.GetSiteUrl(), "blogUrl": configService.GetBlogUrl(),
		"blog": map[string]string{
			"id":    note.NoteId.Hex(),
			"title": note.Title,
			"url":   configService.GetBlogUrl() + "/view/" + note.NoteId.Hex(),
		},
		"commentContent": content,
		// 评论者信息
		"commentUser": map[string]interface{}{"userId": sendUserInfo.UserId.Hex(),
			"username":     sendUserInfo.Username,
			"email":        sendUserInfo.Email,
			"isBlogAuthor": userId == note.UserId.Hex(),
		},
		// 被评论者信息
		"commentedUser": map[string]interface{}{"userId": toUserId,
			"username":     toUserInfo.Username,
			"email":        toUserInfo.Email,
			"isBlogAuthor": toUserId == note.UserId.Hex(),
		},
	}

	ok := false
	ok, _, subject, tpl = this.renderEmail(subject, tpl, token2Value)
	if !ok {
		return false
	}

	// 发送邮件
	ok, _ = this.SendEmail(toUserInfo.Email, subject, tpl)
	return ok
}

// 验证模板是否正确
func (this *EmailService) ValidTpl(str string) (ok bool, msg string) {
	defer func() {
		if err := recover(); err != nil {
			ok = false
			msg = fmt.Sprint(err)
		}
	}()
	header := configService.GetGlobalStringConfig("emailTemplateHeader")
	footer := configService.GetGlobalStringConfig("emailTemplateFooter")
	str = strings.Replace(str, "{{header}}", header, -1)
	str = strings.Replace(str, "{{footer}}", footer, -1)
	_, err := template.New("tpl name").Parse(str)
	if err != nil {
		msg = fmt.Sprint(err)
		return
	}
	ok = true
	return
}

// ok, msg, subject, tpl
func (this *EmailService) getTpl(str string) (ok bool, msg string, tpl *template.Template) {
	defer func() {
		if err := recover(); err != nil {
			ok = false
			msg = fmt.Sprint(err)
		}
	}()

	var err error
	var has bool

	if tpl, has = this.tpls[str]; !has {
		tpl, err = template.New("tpl name").Parse(str)
		if err != nil {
			msg = fmt.Sprint(err)
			return
		}
		this.tpls[str] = tpl
	}
	ok = true
	return
}

// 通过subject, body和值得到内容
func (this *EmailService) renderEmail(subject, body string, values map[string]interface{}) (ok bool, msg string, o string, b string) {
	ok = false
	msg = ""
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			ok = false
			msg = fmt.Sprint(err) // 这里的err其实就是panic传入的内容，
		}
	}()

	var tpl *template.Template

	values["siteUrl"] = configService.GetSiteUrl()

	// subject
	if subject != "" {
		ok, msg, tpl = this.getTpl(subject)
		if !ok {
			return
		}
		var buffer bytes.Buffer
		err := tpl.Execute(&buffer, values)
		if err != nil {
			msg = fmt.Sprint(err)
			return
		}
		o = buffer.String()
	} else {
		o = ""
	}

	// content
	header := configService.GetGlobalStringConfig("emailTemplateHeader")
	footer := configService.GetGlobalStringConfig("emailTemplateFooter")
	body = strings.Replace(body, "{{header}}", header, -1)
	body = strings.Replace(body, "{{footer}}", footer, -1)
	values["subject"] = o
	ok, msg, tpl = this.getTpl(body)
	if !ok {
		return
	}
	var buffer2 bytes.Buffer
	err := tpl.Execute(&buffer2, values)
	if err != nil {
		msg = fmt.Sprint(err)
		return
	}
	b = buffer2.String()

	return
}

// 发送email给用户
// 需要记录
func (this *EmailService) SendEmailToUsers(users []info.User, subject, body string) (ok bool, msg string) {
	if users == nil || len(users) == 0 {
		msg = "no users"
		return
	}

	// 尝试renderHtml
	ok, msg, _, _ = this.renderEmail(subject, body, map[string]interface{}{})
	if !ok {
		Log(msg)
		return
	}

	go func() {
		for _, user := range users {
			LogJ(user)
			m := map[string]interface{}{}
			m["userId"] = user.UserId.Hex()
			m["username"] = user.Username
			m["email"] = user.Email
			ok2, msg2, subject2, body2 := this.renderEmail(subject, body, m)
			ok = ok2
			msg = msg2
			if ok2 {
				sendOk, msg := this.SendEmail(user.Email, subject2, body2)
				this.AddEmailLog(user.Email, subject, body, sendOk, msg) // 把模板记录下
				// 记录到Email Log
				if sendOk {
					// Log("ok " + user.Email)
				} else {
					// Log("no " + user.Email)
				}
			} else {
				// Log(msg);
			}
		}
	}()

	return
}

func (this *EmailService) SendEmailToEmails(emails []string, subject, body string) (ok bool, msg string) {
	if emails == nil || len(emails) == 0 {
		msg = "no emails"
		return
	}

	// 尝试renderHtml
	ok, msg, _, _ = this.renderEmail(subject, body, map[string]interface{}{})
	if !ok {
		Log(msg)
		return
	}

	//	go func() {
	for _, email := range emails {
		if email == "" {
			continue
		}
		m := map[string]interface{}{}
		m["email"] = email
		ok, msg, subject, body = this.renderEmail(subject, body, m)
		if ok {
			sendOk, msg := this.SendEmail(email, subject, body)
			this.AddEmailLog(email, subject, body, sendOk, msg)
			// 记录到Email Log
			if sendOk {
				Log("ok " + email)
			} else {
				Log("no " + email)
			}
		} else {
			Log(msg)
		}
	}
	//	}()

	return
}

// 添加邮件日志
func (this *EmailService) AddEmailLog(email, subject, body string, ok bool, msg string) {
	log := info.EmailLog{LogId: bson.NewObjectId(), Email: email, Subject: subject, Body: body, Ok: ok, Msg: msg, CreatedTime: time.Now()}
	db.Insert(db.EmailLogs, log)
}

// 展示邮件日志

func (this *EmailService) DeleteEmails(ids []string) bool {
	idsO := make([]bson.ObjectId, len(ids))
	for i, id := range ids {
		idsO[i] = bson.ObjectIdHex(id)
	}
	db.DeleteAll(db.EmailLogs, bson.M{"_id": bson.M{"$in": idsO}})

	return true
}
func (this *EmailService) ListEmailLogs(pageNumber, pageSize int, sortField string, isAsc bool, email string) (page info.Page, emailLogs []info.EmailLog) {
	emailLogs = []info.EmailLog{}
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)
	query := bson.M{}
	if email != "" {
		query["Email"] = bson.M{"$regex": bson.RegEx{".*?" + email + ".*", "i"}}
	}
	q := db.EmailLogs.Find(query)
	// 总记录数
	count, _ := q.Count()
	// 列表
	q.Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&emailLogs)
	page = info.NewPage(pageNumber, pageSize, count, nil)
	return
}
