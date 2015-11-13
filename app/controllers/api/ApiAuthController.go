package api

import (
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	//	"strconv"
)

// 用户登录后生成一个token, 将这个token保存到session中
// 以后每次的请求必须带这个token, 并从session中获取userId

// 用户登录/注销/找回密码

type ApiAuth struct {
	ApiBaseContrller
}

// 登录
// [ok]
// 成功返回 {Ok: true, Item: token }
// 失败返回 {Ok: false, Msg: ""}
func (c ApiAuth) Login(email, pwd string) revel.Result {
	var msg = ""

	userInfo, err := authService.Login(email, pwd)
	if err == nil {
		token := bson.NewObjectId().Hex()
		sessionService.SetUserId(token, userInfo.UserId.Hex())
		return c.RenderJson(info.AuthOk{Ok: true, Token: token, UserId: userInfo.UserId, Email: userInfo.Email, Username: userInfo.Username})
	} else {
		// 登录错误, 则错误次数++
		msg = "wrongUsernameOrPassword"
	}
	return c.RenderJson(info.ApiRe{Ok: false, Msg: c.Message(msg)})
}

// 注销
// [Ok]
func (c ApiAuth) Logout() revel.Result {
	token := c.getToken()
	sessionService.Clear(token)
	re := info.NewApiRe()
	re.Ok = true
	return c.RenderJson(re)
}

// 注册
// [Ok]
// 成功后并不返回用户ID, 需要用户重新登录
func (c ApiAuth) Register(email, pwd string) revel.Result {
	re := info.NewApiRe()
	if !configService.IsOpenRegister() {
		re.Msg = "notOpenRegister" // 未开放注册
		return c.RenderJson(re)
	}

	if re.Ok, re.Msg = Vd("email", email); !re.Ok {
		return c.RenderJson(re)
	}
	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderJson(re)
	}

	// 注册
	re.Ok, re.Msg = authService.Register(email, pwd, "")
	return c.RenderJson(re)
}
