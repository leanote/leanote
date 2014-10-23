package controllers

import (
	"strings"
	"time"
	"github.com/revel/revel"
//	"encoding/json"
	"gopkg.in/mgo.v2/bson"
//	. "github.com/leanote/leanote/app/lea"
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

//-----------------------------
// 前台

// 域名, 没用
func (c Blog) domain() (ok bool, userBlog info.UserBlog) {
	return 
}

// 各种地址设置
func (c Blog) setUrl(userBlog info.UserBlog, userInfo info.User) {
	// 主页 http://leanote.com/blog/life or http://blog.leanote.com/life or http:// xxxx.leanote.com or aa.com
	var indexUrl, viewUrl, searchUrl, cateUrl, aboutMeUrl, staticUrl string
	host := c.Request.Request.Host
	staticUrl = configService.GetUserUrl(strings.Split(host, ":")[0])
	// staticUrl == host, 为保证同源!!! 只有host, http://leanote.com, http://blog/leanote.com
	// life.leanote.com, lealife.com
	if userBlog.Domain != "" && configService.AllowCustomDomain() {
		// ok
		indexUrl = configService.GetUserUrl(userBlog.Domain)
		cateUrl = indexUrl + "/cate" // /xxxxx
		viewUrl = indexUrl + "/view" // /xxxxx
		searchUrl = indexUrl + "/search" // /xxxxx
		aboutMeUrl = indexUrl + "/aboutMe"
	} else if userBlog.SubDomain != "" {
		indexUrl = configService.GetUserSubUrl(userBlog.SubDomain)
		cateUrl = indexUrl + "/cate" // /xxxxx
		viewUrl = indexUrl + "/view" // /xxxxx
		searchUrl = indexUrl + "/search" // /xxxxx
		aboutMeUrl = indexUrl + "/aboutMe"
	} else {
		// ok
		blogUrl := configService.GetBlogUrl()
		userIdOrEmail := ""
		if userInfo.Username != "" {
			userIdOrEmail = userInfo.Username
		} else if userInfo.Email != "" {
			userIdOrEmail = userInfo.Email
		} else {
			userIdOrEmail = userInfo.UserId.Hex()
		}
		indexUrl = blogUrl + "/" + userIdOrEmail
		cateUrl = blogUrl + "/cate" // /notebookId
		viewUrl = blogUrl + "/view" // /xxxxx
		searchUrl = blogUrl + "/search/" + userIdOrEmail // /xxxxx
		aboutMeUrl = blogUrl + "/aboutMe/" + userIdOrEmail
	}
	
	// 分类
	// 搜索
	// 查看
	c.RenderArgs["indexUrl"] = indexUrl
	c.RenderArgs["cateUrl"] = cateUrl
	c.RenderArgs["viewUrl"] = viewUrl
	c.RenderArgs["searchUrl"] = searchUrl
	c.RenderArgs["aboutMeUrl"] = aboutMeUrl
	c.RenderArgs["staticUrl"] = staticUrl
}

// 公共
func (c Blog) blogCommon(userId string, userBlog info.UserBlog, userInfo info.User) (ok bool) {
	if userInfo.UserId == "" {
		userInfo = userService.GetUserInfoByAny(userId)
		if userInfo.UserId == "" {
			return
		}
	}
	c.RenderArgs["userInfo"] = userInfo
	
	// 分类导航
	c.RenderArgs["notebooks"] = blogService.ListBlogNotebooks(userId)
	// 最新笔记
	c.getRecentBlogs(userId)
	// 语言, url地址
	c.SetLocale();
	c.RenderArgs["isMe"] = userId == c.GetUserId()
	
	// 得到博客设置信息
	if userBlog.UserId == "" {
		userBlog = blogService.GetUserBlog(userId)
	}
	c.RenderArgs["userBlog"] = userBlog
	
	c.setUrl(userBlog, userInfo)
	
	return true
}

// 跨域判断是否是我的博客
func (c Blog) IsMe(userId string) revel.Result {
	var js = ""
	if c.GetUserId() == userId {
		js = "$('.is-me').removeClass('hide');"
	}
	return c.RenderText(js);
}

// 进入某个用户的博客
var blogPageSize = 5
var searchBlogPageSize = 30

// 分类 /cate/xxxxxxxx?notebookId=1212
func (c Blog) Cate(notebookId string) revel.Result {
	if notebookId == "" {
		return c.E404()
	}
	// 自定义域名
	hasDomain, userBlog := c.domain()
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}
	
	var notebook info.Notebook
	notebook = notebookService.GetNotebookById(notebookId)
	if !notebook.IsBlog {
		return c.E404()
	}
	if userId != "" && userId != notebook.UserId.Hex() {
		return c.E404()
	}
	userId = notebook.UserId.Hex()
		
	if !c.blogCommon(userId, userBlog, info.User{}) {
		return c.E404()
	}
	
	// 分页的话, 需要分页信息, totalPage, curPage
	page := c.GetPage()
	count, blogs := blogService.ListBlogs(userId, notebookId, page, blogPageSize, "PublicTime", false)
	
	c.RenderArgs["notebookId"] = notebookId
	c.RenderArgs["notebook"] = notebook
	c.RenderArgs["title"] = c.Message("blogClass") + " - " + notebook.Title
	c.RenderArgs["blogs"] = blogs
	c.RenderArgs["page"] = page
	c.RenderArgs["pageSize"] = blogPageSize
	c.RenderArgs["count"] = count
	
	return c.RenderTemplate("blog/index.html")
}

// 显示分类的最近博客, json
func (c Blog) ListCateLatest(notebookId, callback string) revel.Result {
	if notebookId == "" {
		return c.E404()
	}
	// 自定义域名
	hasDomain, userBlog := c.domain()
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}
	
	var notebook info.Notebook
	notebook = notebookService.GetNotebookById(notebookId)
	if !notebook.IsBlog {
		return c.E404()
	}
	if userId != "" && userId != notebook.UserId.Hex() {
		return c.E404()
	}
	userId = notebook.UserId.Hex()
		
	if !c.blogCommon(userId, userBlog, info.User{}) {
		return c.E404()
	}
	
	// 分页的话, 需要分页信息, totalPage, curPage
	page := 1
	_, blogs := blogService.ListBlogs(userId, notebookId, page, 5, "PublicTime", false)
	re := info.NewRe()
	re.Ok = true
	re.List = blogs
	return c.RenderJsonP(callback, re)
}

func (c Blog) Index(userIdOrEmail string) revel.Result {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}

	// 用户id为空, 转至博客平台
	if userIdOrEmail == "" {
		userIdOrEmail = leanoteUserId;
	}
	var userInfo info.User
	if userId != "" {
		userInfo = userService.GetUserInfoByAny(userId)
	} else {
		userInfo = userService.GetUserInfoByAny(userIdOrEmail)
	}
	userId = userInfo.UserId.Hex()

	if !c.blogCommon(userId, userBlog, userInfo) {
		return c.E404()
	}
	
	// 分页的话, 需要分页信息, totalPage, curPage
	page := c.GetPage()
	count, blogs := blogService.ListBlogs(userId, "", page, blogPageSize, "PublicTime", false)
	
	c.RenderArgs["blogs"] = blogs
	c.RenderArgs["page"] = page
	c.RenderArgs["pageSize"] = blogPageSize
	c.RenderArgs["count"] = count
	c.RenderArgs["index"] = true
	c.RenderArgs["notebookId"] = ""
	c.RenderArgs["title"] = userBlog.Title
	
	return c.RenderTemplate("blog/index.html")
}

// 详情
func (c Blog) View(noteId string) revel.Result {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}
	
	blog := blogService.GetBlog(noteId)
	userInfo := userService.GetUserInfo(blog.UserId.Hex())
	if userId != "" && userInfo.UserId.Hex() != userId {
		return c.E404()
	}
	c.RenderArgs["blog"] = blog
	c.RenderArgs["userInfo"] = userInfo
	c.RenderArgs["title"] = blog.Title + " - " + userInfo.Username
	
	userId = userInfo.UserId.Hex()
	c.blogCommon(userId, userBlog, info.User{})
	
	// 得到访问者id
	visitUserId := c.GetUserId()
	if(visitUserId != "") {
		visitUserInfo := userService.GetUserInfo(visitUserId)
		c.RenderArgs["visitUserInfoJson"] = c.Json(visitUserInfo)
		c.RenderArgs["visitUserInfo"] = visitUserInfo
	} else {
		c.RenderArgs["visitUserInfoJson"] = "{}";
	}
	
	return c.RenderTemplate("blog/view.html")
}

// 搜索
func (c Blog) Search(userIdOrEmail, key string) revel.Result {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}
	
	c.RenderArgs["title"] = c.Message("search") + " - " + key
	c.RenderArgs["key"] = key
	
	var userInfo info.User
	if userId != "" {
		userInfo = userService.GetUserInfoByAny(userId)
	} else {
		userInfo = userService.GetUserInfoByAny(userIdOrEmail)
	}
	c.RenderArgs["userInfo"] = userInfo
	userId = userInfo.UserId.Hex()
	c.blogCommon(userId, userBlog, userInfo)
	
	page := c.GetPage()
	_, blogs := blogService.SearchBlog(key, userId, page, searchBlogPageSize, "PublicTime", false)
	
	c.RenderArgs["blogs"] = blogs
	c.RenderArgs["key"] = key

	return c.RenderTemplate("blog/search.html")
}

// 博客设置
func (c Blog) Set() revel.Result {
	userId := c.GetUserId()
	userInfo := userService.GetUserInfo(userId)
	c.RenderArgs["userInfo"] = userInfo
	
	// 得到博客设置信息
	c.RenderArgs["title"] = c.Message("blogSet")
	c.RenderArgs["isMe"] = true
	c.RenderArgs["set"] = true
	
	c.RenderArgs["allowCustomDomain"] = configService.GetGlobalStringConfig("allowCustomDomain")
	
	userBlog := blogService.GetUserBlog(userId)
	c.blogCommon(userId, userBlog, info.User{})
	
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
func (c Blog) AboutMe(userIdOrEmail string) revel.Result {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}
	
	var userInfo info.User
	if userId != "" {
		userInfo = userService.GetUserInfoByAny(userId)
	} else {
		userInfo = userService.GetUserInfoByAny(userIdOrEmail)
	}
	
	if userInfo.UserId == "" {
		return c.E404()
	}
	userId = userInfo.UserId.Hex()
	
	c.RenderArgs["userInfo"] = userInfo
	c.RenderArgs["aboutMe"] = true

	c.RenderArgs["title"] = c.Message("aboutMe")
	c.blogCommon(userId, userBlog, info.User{})
	
	return c.RenderTemplate("blog/about_me.html")
}

// 优化, 这里不要得到count
func (c Blog) getRecentBlogs(userId string) {
	_, c.RenderArgs["recentBlogs"] = blogService.ListBlogs(userId, "", 1, 5, "PublicTime", false)
}

// 可以不要, 因为注册的时候已经把username设为email了
func (c Blog) setRenderUserInfo(userInfo info.User) {
	if userInfo.Username == "" {
		userInfo.Username = userInfo.Email
	}
	c.RenderArgs["userInfo"] = userInfo
}

//---------------------------
// 后台 note<->blog

// 设置/取消Blog; 置顶
func (c Blog) SetNote2Blog(noteId string, isBlog, isTop bool) revel.Result {
	noteUpdate := bson.M{}
	if isTop {
		isBlog = true
	}
	if !isBlog {
		isTop = false
	}
	noteUpdate["IsBlog"] = isBlog
	noteUpdate["IsTop"] = isTop
	if isBlog {
		noteUpdate["PublicTime"] = time.Now()
	}
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

//----------------
// 社交, 点赞, 评论

// 我是否点过赞?
// 所有点赞的用户列表
// 各个评论中是否我也点过赞?
func (c Blog) GetLike(noteId string) revel.Result {
	userId := c.GetUserId()
	
	// 我也点过?
	isILikeIt := blogService.IsILikeIt(noteId, userId)
	// 点赞用户列表
	likedUsers, hasMoreLikedUser := blogService.ListLikedUsers(noteId, false)

	result := map[string]interface{}{}
	result["isILikeIt"] = isILikeIt
	result["likedUsers"] = likedUsers
	result["hasMoreLikedUser"] = hasMoreLikedUser
	
	return c.RenderJson(result)
}
func (c Blog) GetLikeAndComments(noteId string) revel.Result {
	userId := c.GetUserId()
	
	// 我也点过?
	isILikeIt := blogService.IsILikeIt(noteId, userId)
	// 点赞用户列表
	likedUsers, hasMoreLikedUser := blogService.ListLikedUsers(noteId, false)
	// 评论
	page := c.GetPage()
	pageInfo, comments, commentUserInfo := blogService.ListComments(userId, noteId, page, 15)
	
	result := map[string]interface{}{}
	result["isILikeIt"] = isILikeIt
	result["likedUsers"] = likedUsers
	result["hasMoreLikedUser"] = hasMoreLikedUser
	result["pageInfo"] = pageInfo
	result["comments"] = comments
	result["commentUserInfo"] = commentUserInfo
	
	return c.RenderJson(result)
}

func (c Blog) IncReadNum(noteId string) revel.Result {
	blogService.IncReadNum(noteId)
	return nil
}
// 点赞
func (c Blog) LikeBlog(noteId string) revel.Result {
	userId := c.GetUserId()
	re := info.NewRe()
	re.Ok, re.Item = blogService.LikeBlog(noteId, userId)
	
	return c.RenderJson(re)
}
func (c Blog) ListLikes(noteId string) revel.Result {
	return nil
}

func (c Blog) ListComments(noteId string) revel.Result {
	// 评论
	userId := c.GetUserId()
	page := c.GetPage()
	pageInfo, comments, commentUserInfo := blogService.ListComments(userId, noteId, page, 15)
	
	result := map[string]interface{}{}
	result["pageInfo"] = pageInfo
	result["comments"] = comments
	result["commentUserInfo"] = commentUserInfo
	
	return c.RenderJson(result)
}
func (c Blog) DeleteComment(noteId, commentId string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.DeleteComment(noteId, commentId, c.GetUserId())
	return c.RenderJson(re)
}
func (c Blog) Comment(noteId, content, toCommentId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Item = blogService.Comment(noteId, toCommentId, c.GetUserId(), content);
	return c.RenderJson(re)
}
func (c Blog) LikeComment(commentId string) revel.Result {
	re := info.NewRe()
	ok, isILikeIt, num := blogService.LikeComment(commentId, c.GetUserId())
	re.Ok = ok
	re.Item = bson.M{"IsILikeIt": isILikeIt, "Num": num}
	return c.RenderJson(re)
}

func (c Blog) Report(noteId, commentId, reason string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.Report(noteId, commentId, reason, c.GetUserId());
	return c.RenderJson(re)
}