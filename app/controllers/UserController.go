package controllers

import (
	"github.com/revel/revel"
//	"encoding/json"
//	"gopkg.in/mgo.v2/bson"
	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/info"
//	"github.com/leanote/leanote/app/types"
//	"io/ioutil"
	"fmt"
//	"math"
//	"os"
//	"path"
	"strconv"
)

type User struct {
	BaseController
}

// 修改用户名, 需要重置session
func (c User) UpdateUsername(username string) revel.Result {
	re := info.NewRe();
	// 判断是否满足最基本的, 4位, 不含特殊字符, 大小写无关. email大小写无关
	if len(username) < 4 {
		re.Ok = false
		re.Msg = "至少4位"
		return c.RenderJson(re);
	}
	if !IsUsername(username) {
		re.Ok = false
		re.Msg = "不能包含特殊字符"
		return c.RenderJson(re);
	}
	
	re.Ok, re.Msg = userService.UpdateUsername(c.GetUserId(), username)
	if(re.Ok) {
		c.UpdateSession("Username", username)
	}
	return c.RenderJson(re);
}

// 修改密码
func (c User) UpdatePwd(oldPwd, pwd string) revel.Result {
	re := info.NewRe();
	if oldPwd == "" {
		re.Msg = "旧密码错误"
		return c.RenderJson(re);
	}
	
	re.Ok, re.Msg = IsGoodPwd(pwd)
	if !re.Ok {
		return c.RenderJson(re);
	}
	
	re.Ok, re.Msg = userService.UpdatePwd(c.GetUserId(), oldPwd, pwd)
	return c.RenderJson(re);
}

// 更新主题
func (c User) UpdateTheme(theme string) revel.Result {
	re := info.NewRe();
	re.Ok = userService.UpdateTheme(c.GetUserId(), theme)
	if re.Ok {
		c.UpdateSession("Theme", theme)
	}
	return c.RenderJson(re);
}

// 发送邀请链接
func (c User) SendRegisterEmail(content, toEmail string) revel.Result {
	re := info.NewRe()
	if content == "" || !IsEmail(toEmail) {
		return c.RenderJson(re);
	}
	
	// 发送邮件
	var userInfo = c.GetUserInfo();
	siteUrl, _ := revel.Config.String("site.url")
	url := siteUrl + "/register?from=" + userInfo.Username
	body := fmt.Sprintf("点击链接注册leanote: <a href='%v'>%v</a>. ", url, url);
	body = content + "<br />" + body
	re.Ok = SendEmail(toEmail, userInfo.Username + "邀请您注册leanote", "邀请注册", body)
	
	return c.RenderJson(re);
}

//---------------------------

// 重新发送激活邮件
func (c User) ReSendActiveEmail() revel.Result {
	re := info.NewRe()
	re.Ok = userService.RegisterSendActiveEmail(c.GetUserId(), c.GetEmail())
	return c.RenderJson(re)
}

// 修改Email发送激活邮箱
func (c User) UpdateEmailSendActiveEmail(email string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = userService.UpdateEmailSendActiveEmail(c.GetUserId(), email)
	return c.RenderJson(re)
}

// 通过点击链接
// 修改邮箱
func (c User) UpdateEmail(token string) revel.Result {
	c.SetUserInfo();
	
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
	c.SetUserInfo();
	
	ok, msg, email := userService.ActiveEmail(token)
	
	// 需要修改session
	if ok {
		c.UpdateSession("Verified", "1");
	}
	
	c.RenderArgs["title"] = "验证邮箱"
	c.RenderArgs["ok"] = ok
	c.RenderArgs["msg"] = msg
	c.RenderArgs["email"] = email
	
	return c.RenderTemplate("user/active_email.html")
}

//------------
// 第三方账号添加leanote账号
func (c User) AddAccount(email, pwd string) revel.Result {
	re := info.NewRe()
		
	if email == "" {
		re.Msg = "请输入邮箱"
		return c.RenderJson(re)
	} else if !IsEmail(email) {
		re.Msg = "请输入正确的邮箱"
		return c.RenderJson(re)
	}
	
	// 密码
	if pwd == "" {
		re.Msg = "请输入密码"
		return c.RenderJson(re)
	} else if len(pwd) < 6 {
		re.Msg = "密码长度至少6位"
		return c.RenderJson(re)
	}
	
	re.Ok, re.Msg = userService.ThirdAddUser(c.GetUserId(), email, pwd)
	
	if re.Ok {
		c.UpdateSession("Email", email);
	}
	
	return c.RenderJson(re)
}

//-----------------
// 用户偏爱
func (c User) UpdateColumnWidth(notebookWidth, noteListWidth int) revel.Result {
	re := info.NewRe()
	re.Ok = userService.UpdateColumnWidth(c.GetUserId(), notebookWidth, noteListWidth)
	if re.Ok {
		c.UpdateSession("NotebookWidth", strconv.Itoa(notebookWidth));
		c.UpdateSession("NoteListWidth", strconv.Itoa(noteListWidth));
	}
	return c.RenderJson(re)
}
func (c User) UpdateLeftIsMin(leftIsMin bool) revel.Result {
	re := info.NewRe()
	re.Ok = userService.UpdateLeftIsMin(c.GetUserId(), leftIsMin)
	if re.Ok {
		if leftIsMin {
			c.UpdateSession("LeftIsMin", "1");
		} else {
			c.UpdateSession("LeftIsMin", "0");
		}
	}
	return c.RenderJson(re)
}