package admin

import (
	"github.com/revel/revel"
//	. "github.com/leanote/leanote/app/lea"
)

// admin 首页

type AdminUser struct {
	AdminBaseController
}

// admin 主页
var userPageSize = 10
func (c AdminUser) Index(sorter, keywords string) revel.Result {
	pageNumber := c.GetPage()
	sorterField, isAsc := c.getSorter("CreatedTime", false, []string{"email", "username", "verified", "createdTime"});
	pageInfo, users := userService.ListUsers(pageNumber, userPageSize, sorterField, isAsc, keywords);
	c.RenderArgs["pageInfo"] = pageInfo
	c.RenderArgs["users"] = users
	c.RenderArgs["keywords"] = keywords
	return c.RenderTemplate("admin/user/list.html");
}

func (c AdminUser) Add() revel.Result {
	return c.RenderTemplate("admin/user/add.html");
}
