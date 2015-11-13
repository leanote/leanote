package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"time"
	//	"strings"
)

// Session存储到mongodb中
type SessionService struct {
}

func (this *SessionService) Update(sessionId, key string, value interface{}) bool {
	return db.UpdateByQMap(db.Sessions, bson.M{"SessionId": sessionId},
		bson.M{key: value, "UpdatedTime": time.Now()})
}

// 注销时清空session
func (this *SessionService) Clear(sessionId string) bool {
	return db.Delete(db.Sessions, bson.M{"SessionId": sessionId})
}
func (this *SessionService) Get(sessionId string) info.Session {
	session := info.Session{}
	db.GetByQ(db.Sessions, bson.M{"SessionId": sessionId}, &session)

	// 如果没有session, 那么插入一条之
	if session.Id == "" {
		session.Id = bson.NewObjectId()
		session.SessionId = sessionId
		session.CreatedTime = time.Now()
		session.UpdatedTime = session.CreatedTime
		db.Insert(db.Sessions, session)
	}

	return session
}

//------------------
// 错误次数处理

// 登录错误时间是否已超过了
func (this *SessionService) LoginTimesIsOver(sessionId string) bool {
	session := this.Get(sessionId)
	return session.LoginTimes > 5
}

// 登录成功后清空错误次数
func (this *SessionService) ClearLoginTimes(sessionId string) bool {
	return this.Update(sessionId, "LoginTimes", 0)
}

// 增加错误次数
func (this *SessionService) IncrLoginTimes(sessionId string) bool {
	session := this.Get(sessionId)
	return this.Update(sessionId, "LoginTimes", session.LoginTimes+1)
}

//----------
// 验证码
func (this *SessionService) GetCaptcha(sessionId string) string {
	session := this.Get(sessionId)
	return session.Captcha
}
func (this *SessionService) SetCaptcha(sessionId, captcha string) bool {
	this.Get(sessionId)
	Log(sessionId)
	Log(captcha)
	ok := this.Update(sessionId, "Captcha", captcha)
	Log(ok)
	return ok
}

//-----------
// API
func (this *SessionService) GetUserId(sessionId string) string {
	session := this.Get(sessionId)
	// 更新updateTime, 避免过期
	db.UpdateByQMap(db.Sessions, bson.M{"SessionId": sessionId},
		bson.M{"UpdatedTime": time.Now()})
	return session.UserId
}

// 登录成功后设置userId
func (this *SessionService) SetUserId(sessionId, userId string) bool {
	this.Get(sessionId)
	ok := this.Update(sessionId, "UserId", userId)
	return ok
}
