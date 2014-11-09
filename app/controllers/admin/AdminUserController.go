package admin

import (
	"github.com/revel/revel"
	. "github.com/leanote/leanote/app/lea"
	"time"
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
	sorterField, isAsc := c.getSorter("CreatedTime", false, []string{"email", "username", "verified", "createdTime", "accountType"});
	pageInfo, users := userService.ListUsers(pageNumber, pageSize, sorterField, isAsc, keywords);
	c.RenderArgs["pageInfo"] = pageInfo
	c.RenderArgs["users"] = users
	c.RenderArgs["keywords"] = keywords
	return c.RenderTemplate("admin/user/list.html");
}

func (c AdminUser) Add() revel.Result {
	return c.RenderTemplate("admin/user/add.html");
}

// 添加
func (c AdminUser) Register(email, pwd string) revel.Result {
	re := info.NewRe();
	
	if re.Ok, re.Msg = Vd("email", email); !re.Ok {
		return c.RenderRe(re);
	}
	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderRe(re);
	}
	
	// 注册
	re.Ok, re.Msg = authService.Register(email, pwd)
	
	return c.RenderRe(re)
}

// 修改帐户
func (c AdminUser) UpdateAccount(userId string) revel.Result {
	userInfo := userService.GetUserInfo(userId)
	c.RenderArgs["userInfo"] = userInfo
	return c.RenderTemplate("admin/user/update_account.html");
}

func (c AdminUser) DoUpdateAccount(userId, accountType string, accountStartTime, accountEndTime string, maxImageNum, maxImageSize, maxAttachNum, maxAttachSize, maxPerAttachSize int) revel.Result {
	re := info.NewRe();
	s, _ := time.Parse("2006-01-02 15:04:02", accountStartTime)
	e, _ := time.Parse("2006-01-02 15:04:02", accountEndTime)
	re.Ok = userService.UpdateAccount(userId, accountType, s, e, maxImageNum, maxImageSize, maxAttachNum, maxAttachSize, maxPerAttachSize )
	return c.RenderRe(re)
}