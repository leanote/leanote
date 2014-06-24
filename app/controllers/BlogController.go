package controllers

import (
	"github.com/revel/revel"
//	"encoding/json"
	"labix.org/v2/mgo/bson"
//	. "leanote/app/lea"
	"github.com/leanote/leanote/app/info"
//	"github.com/leanote/leanote/app/types"
//	"io/ioutil"
//	"math"
//	"os"
//	"path"
)

type Blog struct {
	BaseController
}

//---------------------------
// 后台 note<->blog

// 设置/取消Blog; 置顶
func (c Blog) SetNote2Blog(noteId string, isBlog, isTop bool) revel.Result {
	if isTop {
		isBlog = true
	}
	if !isBlog {
		isTop = false
	}
	noteUpdate := bson.M{"IsBlog": isBlog, "IsTop": isTop}
	re := noteService.UpdateNote(c.GetUserId(), c.GetUserId(),
			noteId, noteUpdate)
	return c.RenderJson(re)
}

// 设置notebook <-> blog
func (c Blog) SetNotebook2Blog(notebookId string, isBlog bool) revel.Result {
	noteUpdate := bson.M{"IsBlog": isBlog}
	re := notebookService.UpdateNotebook(c.GetUserId(),
			notebookId, noteUpdate)
	return c.RenderJson(re)
}

//-----------------------------
// 前台

// 默认是admin用户的博客
// 列表
// 这里还需要得到其它博客配置信息...
// 配置信息可以放在users表中, 或添加一个user_options表(用户配置表)
var blogPageSize = 5
var searchBlogPageSize = 30
func (c Blog) Index(userId string, notebookId string) revel.Result {
	if userId == "" {
		userId = leanoteUserId
	}
	
	// userId可能是 username, email
	userInfo := userService.GetUserInfoByAny(userId)
	if userInfo.UserId == "" {
		return c.E404()
	}
	
	userId = userInfo.UserId.Hex()
	c.isMe(userId)

	c.RenderArgs["userInfo"] = userInfo
	
	// 得到博客设置信息
	userBlog := blogService.GetUserBlog(userId)
	c.RenderArgs["userBlog"] = userBlog
	
	var notebook info.Notebook
	if notebookId != "" {
		notebook = notebookService.GetNotebook(notebookId, userId)
		if !notebook.IsBlog {
			return c.E404()
		}
		
		c.RenderArgs["title"] = userBlog.Title + " - 分类: " + notebook.Title
	} else {
		c.RenderArgs["title"] = userBlog.Title
	}	
	// 分页的话, 需要分页信息, totalPage, curPage
	page := c.GetPage()
	count, blogs := blogService.ListBlogs(userId, notebookId, page, blogPageSize, "UpdatedTime", false)
	
	c.RenderArgs["blogs"] = blogs
	c.RenderArgs["page"] = page
	c.RenderArgs["pageSize"] = blogPageSize
	c.RenderArgs["count"] = count
	
	// 当前notebook
	c.RenderArgs["notebookId"] = notebookId
	c.RenderArgs["notebook"] = notebook
	
	c.RenderArgs["notebooks"] = blogService.ListBlogNotebooks(userId)
	
	
	if notebookId == "" {
		c.RenderArgs["index"] = true
	}
	
	c.getRecentBlogs(userId)
	
	return c.RenderTemplate("blog/index.html")
}

// 详情
func (c Blog) View(noteId string) revel.Result {
	blog := blogService.GetBlog(noteId)
	c.RenderArgs["blog"] = blog
	
	userInfo := userService.GetUserInfo(blog.UserId.Hex())
	c.RenderArgs["userInfo"] = userInfo
	
	c.RenderArgs["title"] = blog.Title + " - " + userInfo.Email
	
	userId := userInfo.UserId.Hex()
	c.isMe(userId)
	
	c.RenderArgs["notebooks"] = blogService.ListBlogNotebooks(userId)
	
	// 得到博客设置信息
	c.RenderArgs["userBlog"] = blogService.GetUserBlog(userId)
	
	c.getRecentBlogs(userId)
	
	return c.RenderTemplate("blog/view.html")
}

// 搜索
func (c Blog) SearchBlog(userId, key string) revel.Result {
	c.RenderArgs["title"] = "搜索 " + key
	c.RenderArgs["key"] = key
	
	userInfo := userService.GetUserInfoByAny(userId)
	c.RenderArgs["userInfo"] = userInfo
	
	userId = userInfo.UserId.Hex()
	
	page := c.GetPage()
	_, blogs := blogService.SearchBlog(key, userId, page, searchBlogPageSize, "UpdatedTime", false)
	
	c.RenderArgs["blogs"] = blogs
	c.RenderArgs["key"] = key

	c.RenderArgs["notebooks"] = blogService.ListBlogNotebooks(userId)
	// 得到博客设置信息
	c.RenderArgs["userBlog"] = blogService.GetUserBlog(userId)
	
	c.getRecentBlogs(userId)
	
	c.isMe(userId)
	
	return c.RenderTemplate("blog/search.html")
}

// 博客设置
func (c Blog) Set() revel.Result {
	userId := c.GetUserId()
	userInfo := userService.GetUserInfo(userId)
	c.RenderArgs["userInfo"] = userInfo
	
	c.RenderArgs["notebooks"] = blogService.ListBlogNotebooks(userId)
	
	// 得到博客设置信息
	c.RenderArgs["userBlog"] = blogService.GetUserBlog(userId)
	c.RenderArgs["title"] = "博客设置"
	c.RenderArgs["isMe"] = true
	c.RenderArgs["set"] = true
	
	c.getRecentBlogs(userId)
	
	c.SetLocale();
	
	return c.RenderTemplate("blog/set.html")
}

func (c Blog) SetUserBlogBase(userBlog info.UserBlogBase) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.UpdateUserBlogBase(c.GetUserId(), userBlog)
	return c.RenderJson(re)
}
func (c Blog) SetUserBlogComment(userBlog info.UserBlogComment) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.UpdateUserBlogComment(c.GetUserId(), userBlog)
	return c.RenderJson(re)
}
func (c Blog) SetUserBlogStyle(userBlog info.UserBlogStyle) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.UpdateUserBlogStyle(c.GetUserId(), userBlog)
	return c.RenderJson(re)
}

// userId可能是其它的
func (c Blog) AboutMe(userId string) revel.Result {
	userInfo := userService.GetUserInfoByAny(userId)
	if userInfo.UserId == "" {
		return c.E404()
	}
	userId = userInfo.UserId.Hex()
	
	c.RenderArgs["userInfo"] = userInfo
	
	c.RenderArgs["notebooks"] = blogService.ListBlogNotebooks(userId)
	
	c.RenderArgs["userBlog"] = blogService.GetUserBlog(userId)
	c.RenderArgs["aboutMe"] = true
	
	c.RenderArgs["title"] = "关于我"
	
	c.isMe(userId)
	
	c.getRecentBlogs(userId)
	
	return c.RenderTemplate("blog/about_me.html")
}

// 当前的博客是否是我的
func (c Blog) isMe(userId string) {
	c.RenderArgs["isMe"] = userId == c.GetUserId()
}

// 优化, 这里不要得到count
func (c Blog) getRecentBlogs(userId string) {
	_, c.RenderArgs["recentBlogs"] = blogService.ListBlogs(userId, "", 1, 5, "UpdatedTime", false)
}

// 可以不要, 因为注册的时候已经把username设为email了
func (c Blog) setRenderUserInfo(userInfo info.User) {
	if userInfo.Username == "" {
		userInfo.Username = userInfo.Email
	}
	c.RenderArgs["userInfo"] = userInfo
}