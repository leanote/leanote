package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// token
// 找回密码
// 修改密码

type TokenService struct {
}

// 生成token
func (this *TokenService) NewToken(userId string, email string, tokenType int) string {
	token := info.Token{UserId: bson.ObjectIdHex(userId), Token: NewGuidWith(email), CreatedTime: time.Now(), Email: email, Type: tokenType}

	if db.Upsert(db.Tokens, bson.M{"_id": token.UserId}, token) {
		return token.Token
	}

	return ""
}

// 删除token
func (this *TokenService) DeleteToken(userId string, tokenType int) bool {
	return db.Delete(db.Tokens, bson.M{"_id": bson.ObjectIdHex(userId), "Type": tokenType})
}

func (this *TokenService) GetOverHours(tokenType int) float64 {
	if tokenType == info.TokenPwd {
		return info.PwdOverHours
	} else if tokenType == info.TokenUpdateEmail {
		return info.UpdateEmailOverHours
	} else {
		return info.ActiveEmailOverHours
	}
}

// 验证token, 是否存在, 过时?
func (this *TokenService) VerifyToken(token string, tokenType int) (ok bool, msg string, tokenInfo info.Token) {
	overHours = this.GetOverHours(tokenType)

	ok = false
	if token == "" {
		msg = "不存在"
		return
	}

	db.GetByQ(db.Tokens, bson.M{"Token": token}, &tokenInfo)

	if tokenInfo.UserId == "" {
		msg = "不存在"
		return
	}

	// 验证是否过时
	now := time.Now()
	duration := now.Sub(tokenInfo.CreatedTime)

	if duration.Hours() > overHours {
		msg = "过期"
		return
	}

	ok = true
	return
}
