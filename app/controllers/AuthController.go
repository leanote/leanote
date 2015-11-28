package controllers

import (
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"strings"
	//	"strconv"
)

// 用户登录/注销/找回密码

type Auth struct {
	BaseController
}

//--------
// 登录
func (c Auth) Login(email, from string) revel.Result {
	c.RenderArgs["title"] = c.Message("login")
	c.RenderArgs["subTitle"] = c.Message("login")
	c.RenderArgs["email"] = email
	c.RenderArgs["from"] = from
	c.RenderArgs["openRegister"] = configService.IsOpenRegister()

	sessionId := c.Session.Id()
	if sessionService.LoginTimesIsOver(sessionId) {
		c.RenderArgs["needCaptcha"] = true
	}

	c.SetLocale()

	if c.Has("demo") {
		c.RenderArgs["demo"] = true
		c.RenderArgs["email"] = "demo@leanote.com"
	}
	return c.RenderTemplate("home/login.html")
}

// 为了demo和register
func (c Auth) doLogin(email, pwd string) revel.Result {
	sessionId := c.Session.Id()
	var msg = ""

	userInfo, err := authService.Login(email, pwd)
	if err != nil {
		// 登录错误, 则错误次数++
		msg = "wrongUsernameOrPassword"
	} else {
		c.SetSession(userInfo)
		sessionService.ClearLoginTimes(sessionId)
		return c.RenderJson(info.Re{Ok: true})
	}

	return c.RenderJson(info.Re{Ok: false, Item: sessionService.LoginTimesIsOver(sessionId), Msg: c.Message(msg)})
}
func (c Auth) DoLogin(email, pwd string, captcha string) revel.Result {
	sessionId := c.Session.Id()
	var msg = ""

	// > 5次需要验证码, 直到登录成功
	if sessionService.LoginTimesIsOver(sessionId) && sessionService.GetCaptcha(sessionId) != captcha {
		msg = "captchaError"
	} else {
		userInfo, err := authService.Login(email, pwd)
		if err != nil {
			// 登录错误, 则错误次数++
			msg = "wrongUsernameOrPassword"
			sessionService.IncrLoginTimes(sessionId)
		} else {
			c.SetSession(userInfo)
			sessionService.ClearLoginTimes(sessionId)
			return c.RenderJson(info.Re{Ok: true})
		}
	}

	return c.RenderJson(info.Re{Ok: false, Item: sessionService.LoginTimesIsOver(sessionId), Msg: c.Message(msg)})
}

// 注销
func (c Auth) Logout() revel.Result {
	sessionId := c.Session.Id()
	sessionService.Clear(sessionId)
	c.ClearSession()
	return c.Redirect("/login")
}

// 体验一下
func (c Auth) Demo() revel.Result {
	email := configService.GetGlobalStringConfig("demoUsername")
	pwd := configService.GetGlobalStringConfig("demoPassword")

	userInfo, err := authService.Login(email, pwd)
	if err != nil {
		return c.RenderJson(info.Re{Ok: false})
	} else {
		c.SetSession(userInfo)
		return c.Redirect("/note")
	}
	return nil
}

//--------
// 注册
func (c Auth) Register(from, iu string) revel.Result {
	if !configService.IsOpenRegister() {
		return c.Redirect("/index")
	}
	c.SetLocale()
	c.RenderArgs["from"] = from
	c.RenderArgs["iu"] = iu

	c.RenderArgs["title"] = c.Message("register")
	c.RenderArgs["subTitle"] = c.Message("register")
	return c.RenderTemplate("home/register.html")
}
func (c Auth) DoRegister(email, pwd, iu string) revel.Result {
	if !configService.IsOpenRegister() {
		return c.Redirect("/index")
	}

	re := info.NewRe()

	if re.Ok, re.Msg = Vd("email", email); !re.Ok {
		return c.RenderRe(re)
	}
	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderRe(re)
	}

	email = strings.ToLower(email)

	// 注册
	re.Ok, re.Msg = authService.Register(email, pwd, iu)

	// 注册成功, 则立即登录之
	if re.Ok {
		c.doLogin(email, pwd)
	}

	return c.RenderRe(re)
}

//--------
// 找回密码
func (c Auth) FindPassword() revel.Result {
	c.SetLocale()
	c.RenderArgs["title"] = c.Message("findPassword")
	c.RenderArgs["subTitle"] = c.Message("findPassword")
	return c.RenderTemplate("home/find_password.html")
}
func (c Auth) DoFindPassword(email string) revel.Result {
	pwdService.FindPwd(email)
	re := info.NewRe()
	re.Ok = true
	return c.RenderJson(re)
}

// 点击链接后, 先验证之
func (c Auth) FindPassword2(token string) revel.Result {
	c.SetLocale()
	c.RenderArgs["title"] = c.Message("findPassword")
	c.RenderArgs["subTitle"] = c.Message("findPassword")
	if token == "" {
		return c.RenderTemplate("find_password2_timeout.html")
	}
	ok, _, findPwd := tokenService.VerifyToken(token, info.TokenPwd)
	if !ok {
		return c.RenderTemplate("home/find_password2_timeout.html")
	}
	c.RenderArgs["findPwd"] = findPwd

	c.RenderArgs["title"] = c.Message("updatePassword")
	c.RenderArgs["subTitle"] = c.Message("updatePassword")

	return c.RenderTemplate("home/find_password2.html")
}

// 找回密码修改密码
func (c Auth) FindPasswordUpdate(token, pwd string) revel.Result {
	re := info.NewRe()

	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderRe(re)
	}

	// 修改之
	re.Ok, re.Msg = pwdService.UpdatePwd(token, pwd)
	return c.RenderRe(re)
}
