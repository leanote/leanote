package controllers

import (
	"github.com/revel/revel"
//	"github.com/leanote/leanote/app/info"
)

// 首页
type Mobile struct {
	BaseController
}

// leanote展示页, 没有登录的, 或已登录明确要进该页的
func (c Mobile) Index() revel.Result {
	c.SetLocale()
	
	userInfo := c.GetUserInfo()
	userId := userInfo.UserId.Hex()
	
	// 没有登录
	if userId == "" {
		return c.RenderTemplate("mobile/login.html")
	}
	
	/*
	// 已登录了, 那么得到所有信息
	notebooks := notebookService.GetNotebooks(userId)
	shareNotebooks, sharedUserInfos := shareService.GetShareNotebooks(userId)
	
	c.RenderArgs["userInfo"] = userInfo
	c.RenderArgs["userInfoJson"] = c.Json(userInfo)
	c.RenderArgs["notebooks"] = c.Json(notebooks)
	c.RenderArgs["shareNotebooks"] = c.Json(shareNotebooks)
	c.RenderArgs["sharedUserInfos"] = c.Json(sharedUserInfos)
	c.RenderArgs["tagsJson"] = c.Json(tagService.GetTags(c.GetUserId()))
	*/
	
	return c.RenderTemplate("mobile/angular.html");
}

func (c Mobile) Logout() revel.Result {
	c.ClearSession()
	return c.RenderTemplate("mobile/login.html");
}