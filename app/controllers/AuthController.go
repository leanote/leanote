package controllers

import (
	"github.com/revel/revel"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
)

// 用户登录/注销/找回密码

type Auth struct {
	BaseController
}

//--------
// 登录
func (c Auth) Login(email string) revel.Result {
	c.RenderArgs["title"] = c.Message("login")
	c.RenderArgs["subTitle"] = c.Message("login")
	c.RenderArgs["email"] = email
	c.RenderArgs["openRegister"] = openRegister
	
	if c.Has("demo") {
		c.RenderArgs["demo"] = true
		c.RenderArgs["email"] = "demo@leanote.com"
	}
	return c.RenderTemplate("home/login.html")
}
func (c Auth) DoLogin(email, pwd string) revel.Result {
	userInfo := authService.Login(email, pwd)
	if userInfo.Email != "" {
		c.SetSession(userInfo)
		// 必须要redirect, 不然用户刷新会重复提交登录信息
//		return c.Redirect("/")
		return c.RenderJson(info.Re{Ok: true})
	}
//	return c.RenderTemplate("login.html")
	return c.RenderJson(info.Re{Ok: false, Msg: c.Message("wrongUsernameOrPassword")})
}
// 注销
func (c Auth) Logout() revel.Result {
	c.ClearSession()
	return c.Redirect("/login")
}

// 体验一下
func (c Auth) Demo() revel.Result {
	c.DoLogin("demo@leanote.com", "demo@leanote.com")
	return c.Redirect("/note")
}

//--------
// 注册
func (c Auth) Register() revel.Result {
	if !openRegister {
		return c.Redirect("/index")
	}
	
	c.RenderArgs["title"] = c.Message("register")
	c.RenderArgs["subTitle"] = c.Message("register")
	return c.RenderTemplate("home/register.html")
}
func (c Auth) DoRegister(email, pwd string) revel.Result {
	if !openRegister {
		return c.Redirect("/index")
	}
	
	re := info.NewRe();
	
	if email == "" {
		re.Msg = c.Message("inputEmail")
		return c.RenderJson(re)
	} else if !IsEmail(email) {
		re.Msg = c.Message("wrongEmail")
		return c.RenderJson(re)
	}
	
	// 密码
	if pwd == "" {
		re.Msg = c.Message("inputPassword")
		return c.RenderJson(re)
	} else if len(pwd) < 6 {
		re.Msg = c.Message("wrongPassword")
		return c.RenderJson(re)
	}
	
	// 注册
	re.Ok, re.Msg = authService.Register(email, pwd)
	
	// 注册成功, 则立即登录之
	if re.Ok {
		c.DoLogin(email, pwd)
	}
	
	return c.RenderJson(re)
}

//--------
// 找回密码
func (c Auth) FindPassword() revel.Result {
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
func (c Auth) FindPasswordUpdate(token, pwd string) revel.Result {
	re := info.NewRe();
	
	re.Ok, re.Msg = IsGoodPwd(pwd)
	if !re.Ok {
		return c.RenderJson(re)
	}

	// 修改之
	re.Ok, re.Msg = pwdService.UpdatePwd(token, pwd)
	return c.RenderJson(re)
}
