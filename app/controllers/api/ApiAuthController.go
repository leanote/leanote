package api

import (
	"github.com/revel/revel"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
//	"strconv"
)

// 用户登录后生成一个token, 将这个token保存到session中
// 以后每次的请求必须带这个token, 并从session中获取userId

// 用户登录/注销/找回密码

type ApiAuth struct {
	ApiBaseContrller
}

// 为了demo和register
func (c ApiAuth) doLogin(email, pwd string) revel.Result {
	sessionId := c.Session.Id()
	var msg = ""
	
	userInfo := authService.Login(email, pwd)
	if userInfo.Email != "" {
		c.SetSession(userInfo)
		sessionService.ClearLoginTimes(sessionId)
		return c.RenderJson(info.Re{Ok: true})
	} else {
		// 登录错误, 则错误次数++
		msg = "wrongUsernameOrPassword"
	}
	
	return c.RenderJson(info.Re{Ok: false, Item: sessionService.LoginTimesIsOver(sessionId) , Msg: c.Message(msg)})
}
func (c ApiAuth) Login(email, pwd string, captcha string) revel.Result {
	sessionId := c.Session.Id()
	var msg = ""
	
	// > 5次需要验证码, 直到登录成功
	if sessionService.LoginTimesIsOver(sessionId) && sessionService.GetCaptcha(sessionId) != captcha {
		msg = "captchaError"
	} else {
		userInfo := authService.Login(email, pwd)
		if userInfo.Email != "" {
			c.SetSession(userInfo)
			sessionService.ClearLoginTimes(sessionId)
			return c.RenderJson(info.Re{Ok: true})
		} else {
			// 登录错误, 则错误次数++
			msg = "wrongUsernameOrPassword"
			sessionService.IncrLoginTimes(sessionId)
		}
	}
	
	return c.RenderJson(info.Re{Ok: false, Item: sessionService.LoginTimesIsOver(sessionId) , Msg: c.Message(msg)})
}

// 注销
func (c ApiAuth) Logout() revel.Result {
	sessionId := c.Session.Id()
	sessionService.Clear(sessionId)
	c.ClearSession()
	re := info.NewRe()
	re.Ok = true
	return c.RenderJson(re)
}

// 体验一下
func (c ApiAuth) Demo() revel.Result {
	c.doLogin(configService.GetGlobalStringConfig("demoUsername"), configService.GetGlobalStringConfig("demoPassword"))
	return c.Redirect("/note")
}

//--------
// 注册
func (c ApiAuth) Register(from string) revel.Result {
	if !configService.IsOpenRegister() {
		return c.Redirect("/index")
	}
	c.SetLocale()
	c.RenderArgs["from"] = from
	
	c.RenderArgs["title"] = c.Message("register")
	c.RenderArgs["subTitle"] = c.Message("register")
	return c.RenderTemplate("home/register.html")
}
func (c ApiAuth) DoRegister(email, pwd string) revel.Result {
	if !configService.IsOpenRegister() {
		return c.Redirect("/index")
	}
	
	re := info.NewRe();
	
	if re.Ok, re.Msg = Vd("email", email); !re.Ok {
		return c.RenderRe(re);
	}
	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderRe(re);
	}
	
	// 注册
	re.Ok, re.Msg = authService.Register(email, pwd)
	
	// 注册成功, 则立即登录之
	if re.Ok {
		c.doLogin(email, pwd)
	}
	
	return c.RenderRe(re)
}

//--------
// 找回密码
func (c ApiAuth) FindPassword() revel.Result {
	c.RenderArgs["title"] = c.Message("findPassword")
	c.RenderArgs["subTitle"] = c.Message("findPassword")
	return c.RenderTemplate("home/find_password.html")
}
func (c ApiAuth) DoFindPassword(email string) revel.Result {
	pwdService.FindPwd(email)
	re := info.NewRe()
	re.Ok = true
	return c.RenderJson(re)
}
// 点击链接后, 先验证之
func (c ApiAuth) FindPassword2(token string) revel.Result {
	c.RenderArgs["title"] = c.Message("findPassword")
	c.RenderArgs["subTitle"] = c.Message("findPassword")
	if token == "" {
		return c.RenderTemplate("find_password2_timeout.html")
	}
	ok, _, findPwd := tokenService.VerifyToken(token, info.TokenPwd);
	if !ok {
		return c.RenderTemplate("home/find_password2_timeout.html")
	}
	c.RenderArgs["findPwd"] = findPwd
	
	c.RenderArgs["title"] = c.Message("updatePassword")
	c.RenderArgs["subTitle"] = c.Message("updatePassword")
	
	return c.RenderTemplate("home/find_password2.html")
}
// 找回密码修改密码
func (c ApiAuth) FindPasswordUpdate(token, pwd string) revel.Result {
	re := info.NewRe();

	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderRe(re);
	}

	// 修改之
	re.Ok, re.Msg = pwdService.UpdatePwd(token, pwd)
	return c.RenderRe(re)
}
