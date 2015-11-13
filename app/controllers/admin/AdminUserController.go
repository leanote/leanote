package admin

import (
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	//	"time"
	"github.com/leanote/leanote/app/info"
)

// admin 首页

type AdminUser struct {
	AdminBaseController
}

// admin 主页
var userPageSize = 10

func (c AdminUser) Index(sorter, keywords string, pageSize int) revel.Result {
	pageNumber := c.GetPage()
	if userPageSize == 0 {
		pageSize = userPageSize
	}
	sorterField, isAsc := c.getSorter("CreatedTime", false, []string{"email", "username", "verified", "createdTime", "accountType"})
	pageInfo, users := userService.ListUsers(pageNumber, pageSize, sorterField, isAsc, keywords)
	c.RenderArgs["pageInfo"] = pageInfo
	c.RenderArgs["users"] = users
	c.RenderArgs["keywords"] = keywords
	return c.RenderTemplate("admin/user/list.html")
}

func (c AdminUser) Add() revel.Result {
	return c.RenderTemplate("admin/user/add.html")
}

// 添加
func (c AdminUser) Register(email, pwd string) revel.Result {
	re := info.NewRe()

	if re.Ok, re.Msg = Vd("email", email); !re.Ok {
		return c.RenderRe(re)
	}
	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderRe(re)
	}

	// 注册
	re.Ok, re.Msg = authService.Register(email, pwd, "")

	return c.RenderRe(re)
}

// 修改帐户
func (c AdminUser) ResetPwd(userId string) revel.Result {
	userInfo := userService.GetUserInfo(userId)
	c.RenderArgs["userInfo"] = userInfo
	return c.RenderTemplate("admin/user/reset_pwd.html")
}

func (c AdminUser) DoResetPwd(userId, pwd string) revel.Result {
	re := info.NewRe()
	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderRe(re)
	}
	re.Ok, re.Msg = userService.ResetPwd(c.GetUserId(), userId, pwd)
	return c.RenderRe(re)
}
