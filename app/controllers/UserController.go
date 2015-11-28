package controllers

import (
	"github.com/revel/revel"
	//	"encoding/json"
	//	"gopkg.in/mgo.v2/bson"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	//	"github.com/leanote/leanote/app/types"
	//	"io/ioutil"
	//	"fmt"
	//	"math"
	//	"os"
	//	"path"
	"strconv"
)

type User struct {
	BaseController
}

func (c User) Account(tab int) revel.Result {
	userInfo := c.GetUserInfo()
	c.RenderArgs["userInfo"] = userInfo
	c.RenderArgs["tab"] = tab
	c.SetLocale()
	return c.RenderTemplate("user/account.html")
}

// 修改用户名, 需要重置session
func (c User) UpdateUsername(username string) revel.Result {
	re := info.NewRe()
	if c.GetUserId() == configService.GetGlobalStringConfig("demoUserId") {
		re.Msg = "cannotUpdateDemo"
		return c.RenderRe(re)
	}

	if re.Ok, re.Msg = Vd("username", username); !re.Ok {
		return c.RenderRe(re)
	}

	re.Ok, re.Msg = userService.UpdateUsername(c.GetUserId(), username)
	if re.Ok {
		c.UpdateSession("Username", username)
	}
	return c.RenderRe(re)
}

// 修改密码
func (c User) UpdatePwd(oldPwd, pwd string) revel.Result {
	re := info.NewRe()
	if c.GetUserId() == configService.GetGlobalStringConfig("demoUserId") {
		re.Msg = "cannotUpdateDemo"
		return c.RenderRe(re)
	}
	if re.Ok, re.Msg = Vd("password", oldPwd); !re.Ok {
		return c.RenderRe(re)
	}
	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderRe(re)
	}
	re.Ok, re.Msg = userService.UpdatePwd(c.GetUserId(), oldPwd, pwd)
	return c.RenderRe(re)
}

// 更新主题
func (c User) UpdateTheme(theme string) revel.Result {
	re := info.NewRe()
	re.Ok = userService.UpdateTheme(c.GetUserId(), theme)
	if re.Ok {
		c.UpdateSession("Theme", theme)
	}
	return c.RenderJson(re)
}

// 发送邀请链接
func (c User) SendRegisterEmail(content, toEmail string) revel.Result {
	re := info.NewRe()
	if content == "" || !IsEmail(toEmail) {
		return c.RenderJson(re)
	}

	re.Ok = emailService.SendInviteEmail(c.GetUserInfo(), toEmail, content)
	return c.RenderJson(re)
}

//---------------------------

// 重新发送激活邮件
func (c User) ReSendActiveEmail() revel.Result {
	re := info.NewRe()
	re.Ok = emailService.RegisterSendActiveEmail(c.GetUserInfo(), c.GetEmail())
	return c.RenderJson(re)
}

// 通过点击链接
// 修改邮箱
func (c User) UpdateEmail(token string) revel.Result {
	c.SetUserInfo()

	ok, msg, email := userService.UpdateEmail(token)

	c.RenderArgs["title"] = "验证邮箱"
	c.RenderArgs["ok"] = ok
	c.RenderArgs["msg"] = msg
	c.RenderArgs["email"] = email

	// 修改session
	if ok {
		c.UpdateSession("Email", email)
	}

	return c.RenderTemplate("user/update_email.html")
}

// 注册后激活邮箱
func (c User) ActiveEmail(token string) revel.Result {
	c.SetUserInfo()

	ok, msg, email := userService.ActiveEmail(token)

	// 需要修改session
	if ok {
		c.UpdateSession("Verified", "1")
	}

	c.RenderArgs["title"] = "验证邮箱"
	c.RenderArgs["ok"] = ok
	c.RenderArgs["msg"] = msg
	c.RenderArgs["email"] = email

	return c.RenderTemplate("user/active_email.html")
}

//-----------------
// 用户偏爱
func (c User) UpdateColumnWidth(notebookWidth, noteListWidth, mdEditorWidth int) revel.Result {
	re := info.NewRe()
	re.Ok = userService.UpdateColumnWidth(c.GetUserId(), notebookWidth, noteListWidth, mdEditorWidth)
	if re.Ok {
		c.UpdateSession("NotebookWidth", strconv.Itoa(notebookWidth))
		c.UpdateSession("NoteListWidth", strconv.Itoa(noteListWidth))
		c.UpdateSession("MdEditorWidth", strconv.Itoa(mdEditorWidth))

		LogJ(c.Session)
	}
	return c.RenderJson(re)
}
func (c User) UpdateLeftIsMin(leftIsMin bool) revel.Result {
	re := info.NewRe()
	re.Ok = userService.UpdateLeftIsMin(c.GetUserId(), leftIsMin)
	if re.Ok {
		if leftIsMin {
			c.UpdateSession("LeftIsMin", "1")
		} else {
			c.UpdateSession("LeftIsMin", "0")
		}
	}
	return c.RenderJson(re)
}
