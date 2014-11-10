package controllers

import (
	"github.com/revel/revel"
	"strings"
	//	"encoding/json"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/lea/blog"
	"gopkg.in/mgo.v2/bson"
	"fmt"
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
/*
公共
// 分类 [ok]
$.cates = [{title, cateId}]
// 单页 [ok]
$.singles = [{pageId, title}]
// 博客信息 [ok]
$.blog = {userId, desc, title, logo, openComment, disqusId}

// 公用url ok
$.indexUrl
$.cateUrl
$.searchUrl
$.postUrl
$.archiveUrl
$.singleUrl
$.themeBaseUrl

// 静态文件 [ok]
$.jQueryUrl
$.fontAsomeUrl
$.bootstrapCssUrl
$.bootstrapJsUrl
*/


func (c Blog) render(templateName string, themePath string) revel.Result {
	isPreview := false
	if c.RenderArgs["isPreview"] != nil  {
		themePath2 := c.RenderArgs["themePath"]
		if themePath2 == nil {
			return c.E404();
		}
		isPreview = true
		themePath = themePath2.(string)
		c.setPreviewUrl()
	}
	return blog.RenderTemplate(templateName, c.RenderArgs, revel.BasePath + "/" + themePath, isPreview)
}

// 404
func (c Blog) e404(themePath string) revel.Result {
	// 不知道是谁的404, 则用系统的404
	if themePath == "" {
		return c.E404();
	}
	return c.render("404.html", themePath)
}

// 二级域名或自定义域名
// life.leanote.com
// lealife.com
func (c Blog) domain() (ok bool, userBlog info.UserBlog) {
	host := c.Request.Request.Host // a.cc.com:9000
	hostArr := strings.Split(host, ".")
	if strings.Contains(host, configService.GetDefaultDomain()) {
		// 有二级域名 a.leanoe.com 3个
		if len(hostArr) > 2 {
			if userBlog = blogService.GetUserBlogBySubDomain(hostArr[0]); userBlog.UserId != "" {
				ok = true
				return
			}
		}
	} else {
		// 自定义域名
		// 把:9000去掉
		hostArr2 := strings.Split(host, ":")
		if userBlog = blogService.GetUserBlogByDomain(hostArr2[0]); userBlog.UserId != "" {
			ok = true
			return
		}
	}
	ok = false
	return
}

// 渲染模板之
func (c Blog) setPreviewUrl() {
	var indexUrl, postUrl, searchUrl, cateUrl, singleUrl, tagsUrl, archiveUrl string
	
	userId := c.GetUserId()
	themeId := c.Session["themeId"]
	theme := themeService.GetTheme(userId, themeId)
	
	siteUrl := configService.GetSiteUrl()
	blogUrl := siteUrl + "/preview" // blog.leanote.com
	userIdOrEmail := userId

	indexUrl = blogUrl + "/" + userIdOrEmail
	cateUrl = blogUrl + "/cate"                        // /notebookId
	
	postUrl = blogUrl + "/post"                        // /xxxxx
	searchUrl = blogUrl + "/search/" + userIdOrEmail   // blog.leanote.com/search/userId
	singleUrl = blogUrl + "/single"                    // blog.leanote.com/single/singleId
	archiveUrl = blogUrl + "/archives/" + userIdOrEmail // blog.leanote.com/archive/userId
	tagsUrl = blogUrl + "/tags/" + userIdOrEmail   // blog.leanote.com/archive/userId
	
	c.RenderArgs["indexUrl"] = indexUrl
	c.RenderArgs["cateUrl"] = cateUrl
	c.RenderArgs["postUrl"] = postUrl
	c.RenderArgs["searchUrl"] = searchUrl
	c.RenderArgs["singleUrl"] = singleUrl // 单页
	c.RenderArgs["archiveUrl"] = archiveUrl
	c.RenderArgs["archivesUrl"] = archiveUrl // 别名
	c.RenderArgs["tagsUrl"] = tagsUrl
	c.RenderArgs["tagPostsUrl"] =  blogUrl + "/tag/" + userIdOrEmail
	c.RenderArgs["tagUrl"] = c.RenderArgs["tagPostsUrl"]

	// themeBaseUrl 本theme的路径url, 可以加载js, css, images之类的
	c.RenderArgs["themeBaseUrl"] = "/" + theme.Path
}

// 各种地址设置
func (c Blog) setUrl(userBlog info.UserBlog, userInfo info.User) {
	// 主页 http://leanote.com/blog/life or http://blog.leanote.com/life or http:// xxxx.leanote.com or aa.com
	var indexUrl, postUrl, searchUrl, cateUrl, singleUrl, tagsUrl, archiveUrl, tagPostsUrl, staticUrl string
	host := c.Request.Request.Host
	staticUrl = configService.GetUserUrl(strings.Split(host, ":")[0])
	// staticUrl == host, 为保证同源!!! 只有host, http://leanote.com, http://blog/leanote.com
	// life.leanote.com, lealife.com
	siteUrl := configService.GetSiteUrl()
	if userBlog.Domain != "" && configService.AllowCustomDomain() {
		// ok
		indexUrl = configService.GetUserUrl(userBlog.Domain)
		cateUrl = indexUrl + "/cate"     // /xxxxx
		postUrl = indexUrl + "/post"     // /xxxxx
		searchUrl = indexUrl + "/search" // /xxxxx
		singleUrl = indexUrl + "/single"
		archiveUrl = indexUrl + "/archives"
		tagsUrl = indexUrl + "/tags"
		tagPostsUrl = indexUrl + "/tag"
	} else if userBlog.SubDomain != "" {
		indexUrl = configService.GetUserSubUrl(userBlog.SubDomain)
		cateUrl = indexUrl + "/cate"     // /xxxxx
		postUrl = indexUrl + "/post"     // /xxxxx
		searchUrl = indexUrl + "/search" // /xxxxx
		singleUrl = indexUrl + "/single"
		archiveUrl = indexUrl + "/archives"
		tagsUrl = indexUrl + "/tags"
		tagPostsUrl = indexUrl + "/tag"
	} else {
		// ok
		blogUrl := configService.GetBlogUrl() // blog.leanote.com
		userIdOrEmail := ""
		if userInfo.Username != "" {
			userIdOrEmail = userInfo.Username
		} else if userInfo.Email != "" {
			userIdOrEmail = userInfo.Email
		} else {
			userIdOrEmail = userInfo.UserId.Hex()
		}
		indexUrl = blogUrl + "/" + userIdOrEmail
		cateUrl = blogUrl + "/cate"                        // /notebookId
		postUrl = blogUrl + "/post"                        // /xxxxx
		searchUrl = blogUrl + "/search/" + userIdOrEmail   // blog.leanote.com/search/userId
		singleUrl = blogUrl + "/single"                    // blog.leanote.com/single/singleId
		archiveUrl = blogUrl + "/archives/" + userIdOrEmail // blog.leanote.com/archive/userId
		tagsUrl = blogUrl + "/tags/" + userIdOrEmail
		tagPostsUrl = blogUrl + "/tag/" + userIdOrEmail   // blog.leanote.com/archive/userId
	}

	// 分类
	// 搜索
	// 查看
	c.RenderArgs["siteUrl"] = siteUrl
	c.RenderArgs["indexUrl"] = indexUrl
	c.RenderArgs["cateUrl"] = cateUrl
	c.RenderArgs["postUrl"] = postUrl
	c.RenderArgs["searchUrl"] = searchUrl
	c.RenderArgs["singleUrl"] = singleUrl // 单页
	c.RenderArgs["archiveUrl"] = archiveUrl
	c.RenderArgs["archivesUrl"] = archiveUrl // 别名
	c.RenderArgs["tagsUrl"] = tagsUrl
	c.RenderArgs["tagPostsUrl"] = tagPostsUrl
	c.RenderArgs["tagUrl"] = c.RenderArgs["tagPostsUrl"] // 别名

	// themeBaseUrl 本theme的路径url, 可以加载js, css, images之类的
	c.RenderArgs["themeBaseUrl"] = "/" + userBlog.ThemePath

	// 其它static js
	c.RenderArgs["jQueryUrl"] = siteUrl + "/js/jquery-1.9.0.min.js"
	
	c.RenderArgs["prettifyJsUrl"] = siteUrl + "/js/google-code-prettify/prettify.js"
	c.RenderArgs["prettifyCssUrl"] = siteUrl + "/js/google-code-prettify/prettify.css"
	
	c.RenderArgs["blogCommonJsUrl"] = siteUrl + "/public/blog/js/common.js"
	
	c.RenderArgs["shareCommentCssUrl"] = siteUrl + "/public/blog/css/share_comment.css"
	c.RenderArgs["shareCommentJsUrl"] = siteUrl + "/public/blog/js/share_comment.js"
	
	c.RenderArgs["fontAwesomeUrl"] = staticUrl + "/css/font-awesome-4.2.0/css/font-awesome.css"
	
	c.RenderArgs["bootstrapCssUrl"] = siteUrl + "/css/bootstrap.css"
	c.RenderArgs["bootstrapJsUrl"] = siteUrl + "/js/bootstrap-min.js"
}

// 笔记本分类
// cates = [{title:"xxx", cateId: "xxxx"}, {}]
func (c Blog) getCates(userBlog info.UserBlog) {
	notebooks := blogService.ListBlogNotebooks(userBlog.UserId.Hex())
	notebooksMap := map[string]info.Notebook{}
	for _, each := range notebooks {
		notebooksMap[each.NotebookId.Hex()] = each
	}

	var i = 0
	cates := make([]map[string]string, len(notebooks))

	// 先要保证已有的是正确的排序
	cateIds := userBlog.CateIds
	has := map[string]bool{} // cateIds中有的
	if cateIds != nil && len(cateIds) > 0 {
		for _, cateId := range cateIds {
			if n, ok := notebooksMap[cateId]; ok {
				cates[i] = map[string]string{"Title": n.Title, "CateId": n.NotebookId.Hex()}
				i++
				has[cateId] = true
			}
		}
	}
	// 之后
	for _, n := range notebooks {
		id := n.NotebookId.Hex()
		if !has[id] {
			cates[i] = map[string]string{"Title": n.Title, "CateId": id}
			i++
		}
	}
	c.RenderArgs["cates"] = cates
}

// 单页
func (c Blog) getSingles(userId string) {
	singles := blogService.GetSingles(userId)
	/*
		if singles == nil {
			return
		}
		singles2 := make([]map[string]string, len(singles))
		for i, page := range singles {
			singles2[i] = map[string]string{"title": page["Title"], "singleId": page["SingleId"]}
		}
	*/
	c.RenderArgs["singles"] = singles
}

// $.blog = {userId, title, subTitle, desc, openComment, }
func (c Blog) setBlog(userBlog info.UserBlog, userInfo info.User) {
	blogInfo := map[string]interface{}{
		"UserId":      userBlog.UserId.Hex(),
		"Username":    userInfo.Username,
		"UserLogo":    userInfo.Logo,
		"Title":       userBlog.Title,
		"SubTitle":    userBlog.SubTitle,
		"Logo":        userBlog.Logo,
		"OpenComment": userBlog.CanComment,
		"CommentType": userBlog.CommentType, // leanote, or disqus
		"DisqusId": userBlog.DisqusId,
		"ThemeId":     userBlog.ThemeId,
		"SubDomain":   userBlog.SubDomain,
		"Domain":      userBlog.Domain,
	}
	c.RenderArgs["blogInfo"] = blogInfo
	/*
		blogInfo := info.BlogInfoCustom{
			UserId: userBlog.UserId.Hex(),
			Username: userInfo.Username,
			UserLogo: userInfo.Logo,
			Title: userBlog.Title,
			SubTitle: userBlog.SubTitle,
			Logo: userBlog.Logo,
			OpenComment: userBlog.CanComment,
			CommentType: userBlog.CommentType, // leanote, or disqus
			ThemeId: userBlog.ThemeId,
			SubDomain: userBlog.SubDomain,
			Domain: userBlog.Domain,
		}
	*/
}

func (c Blog) setPaging(pageInfo info.Page) {
	/*
		// 这些分类信息用一个变量
		paging := map[string]interface{}{ // 分页信息
			"CurPage": pageInfo.CurPage, // 当前页
			"TotalPage": pageInfo.TotalPage, // 总页数
			"PerPageSize": pageInfo.PerPageSize, // 每一页的记录数
			"Count": pageInfo.Count, // 总记录数
		}
	*/
	c.RenderArgs["paging"] = pageInfo
}

// 公共
func (c Blog) blogCommon(userId string, userBlog info.UserBlog, userInfo info.User) (ok bool, ub info.UserBlog) {
	if userInfo.UserId == "" {
		userInfo = userService.GetUserInfoByAny(userId)
		if userInfo.UserId == "" {
			return false, userBlog
		}
	}
//	c.RenderArgs["userInfo"] = userInfo

	// 最新笔记
	_, recentBlogs := blogService.ListBlogs(userId, "", 1, 5, userBlog.SortField, userBlog.IsAsc)
 	c.RenderArgs["recentPosts"] = c.fixBlogs(recentBlogs)
 	c.RenderArgs["latestPosts"] = c.RenderArgs["recentPosts"]
 	c.RenderArgs["tags"] = blogService.GetBlogTags(userId)

	// 语言, url地址
	c.SetLocale()

	// 得到博客设置信息
	if userBlog.UserId == "" {
		userBlog = blogService.GetUserBlog(userId)
	}
	c.setBlog(userBlog, userInfo)
//	c.RenderArgs["userBlog"] = userBlog

	// 分类导航
	c.getCates(userBlog)

	// 单页导航
	c.getSingles(userId)

	c.setUrl(userBlog, userInfo)

	// 当前分类Id, 全设为""
	c.RenderArgs["curCateId"] = ""
	c.RenderArgs["curSingleId"] = ""
	
	// 得到主题信息
	themeInfo := themeService.GetThemeInfo(userBlog.ThemeId.Hex(), userBlog.Style)
	LogJ(themeInfo)
	c.RenderArgs["themeInfo"] = themeInfo

	return true, userBlog
}

// 修复博客, index, cate用到
func (c Blog) fixBlog(blog info.BlogItem) info.Post {
	blog2 := info.Post{
		NoteId:      blog.NoteId.Hex(),
		Title:       blog.Title,
		CreatedTime: blog.CreatedTime,
		UpdatedTime: blog.UpdatedTime,
		PublicTime:  blog.PublicTime,
		Desc:        blog.Desc,
		Abstract:    blog.Abstract,
		Content:     blog.Content,
		Tags:        blog.Tags,
		CommentNum:  blog.CommentNum,
		ReadNum:     blog.ReadNum,
		LikeNum:     blog.LikeNum,
		IsMarkdown:  blog.IsMarkdown,
	}
	return blog2
}
func (c Blog) fixBlogs(blogs []info.BlogItem) []info.Post {
	blogs2 := make([]info.Post, len(blogs))
	for i, blog := range blogs {
		blogs2[i] = c.fixBlog(blog)
	}
	return blogs2
}

// 404
func (c Blog) E(userIdOrEmail, tag string) revel.Result {
	ok, userBlog := c.domain()
	var userId string
	if ok {
		userId = userBlog.UserId.Hex()
	}
	var userInfo info.User
	if userId != "" {
		userInfo = userService.GetUserInfoByAny(userId)
	} else {
		// blog.leanote.com/userid/tag
		userInfo = userService.GetUserInfoByAny(userIdOrEmail)
	}
	userId = userInfo.UserId.Hex()
	_, userBlog = c.blogCommon(userId, userBlog, userInfo);
	
	return c.e404(userBlog.ThemePath)
}

func (c Blog) Tags(userIdOrEmail string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath);
		}
	}()
	
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}

	var userInfo info.User
	if userId != "" {
		userInfo = userService.GetUserInfoByAny(userId)
	} else {
		// blog.leanote.com/userid/tag
		userInfo = userService.GetUserInfoByAny(userIdOrEmail)
	}
	userId = userInfo.UserId.Hex()

	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, userInfo); !ok {
		return c.e404(userBlog.ThemePath) // 404 TODO 使用用户的404
	}

	c.RenderArgs["curIsTags"] = true
	tags := blogService.GetBlogTags(userId)
	c.RenderArgs["tags"] = tags
	return c.render("tags.html", userBlog.ThemePath)
}

// 标签的文章页
func (c Blog) Tag(userIdOrEmail, tag string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath);
		}
	}()
	
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}

	var userInfo info.User
	if userId != "" {
		userInfo = userService.GetUserInfoByAny(userId)
	} else {
		// blog.leanote.com/userid/tag
		userInfo = userService.GetUserInfoByAny(userIdOrEmail)
	}
	userId = userInfo.UserId.Hex()
	
	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, userInfo); !ok {
		return c.e404(userBlog.ThemePath) // 404 TODO 使用用户的404
	}
	
	if hasDomain && tag == "" {
		tag = userIdOrEmail
	}
	
	c.RenderArgs["curIsTagPosts"] = true
	c.RenderArgs["curTag"] = tag
	page := c.GetPage()
	pageInfo, blogs := blogService.SearchBlogByTags([]string{tag}, userId, page, userBlog.PerPageSize, userBlog.SortField, userBlog.IsAsc)
	c.setPaging(pageInfo)

	c.RenderArgs["posts"] = c.fixBlogs(blogs)
	tagPostsUrl := c.RenderArgs["tagPostsUrl"].(string)
	c.RenderArgs["pagingBaseUrl"] = tagPostsUrl + "/" + tag

	return c.render("tag_posts.html", userBlog.ThemePath)
}

// 归档
func (c Blog) Archives(userIdOrEmail string, cateId string, year, month int) (re revel.Result) {
	notebookId := cateId
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			re = c.e404(userBlog.ThemePath);
		}
	}()
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}

	// 用户id为空, 转至博客平台
	if userIdOrEmail == "" {
		userIdOrEmail = configService.GetAdminUsername()
	}
	var userInfo info.User
	if userId != "" {
		userInfo = userService.GetUserInfoByAny(userId)
	} else {
		userInfo = userService.GetUserInfoByAny(userIdOrEmail)
	}
	userId = userInfo.UserId.Hex()

	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, userInfo); !ok {
		return c.e404(userBlog.ThemePath) // 404 TODO 使用用户的404
	}

	arcs := blogService.ListBlogsArchive(userId, notebookId, year, month, "PublicTime", false)
	c.RenderArgs["archives"] = arcs

	c.RenderArgs["curIsArchive"] = true
	if notebookId != "" {
		notebook := notebookService.GetNotebookById(notebookId)
		c.RenderArgs["curCateTitle"] = notebook.Title
		c.RenderArgs["curCateId"] = notebookId
	}
	c.RenderArgs["curYear"] = year
	c.RenderArgs["curMonth"] = month
	

	return c.render("archive.html", userBlog.ThemePath)
}

// 进入某个用户的博客
// TODO 可以配置, 排序和数目
var blogPageSize = 5
var searchBlogPageSize = 30

// 分类 /cate/xxxxxxxx?notebookId=1212
func (c Blog) Cate(notebookId string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath);
		}
	}()
	
	if notebookId == "" {
		return c.e404(userBlog.ThemePath);
	}
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}

	var notebook info.Notebook
	notebook = notebookService.GetNotebookById(notebookId)
	if !notebook.IsBlog {
		return c.e404("")
	}
	if userId != "" && userId != notebook.UserId.Hex() {
		return c.e404("")
	}
	userId = notebook.UserId.Hex()

	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, info.User{}); !ok {
		return c.e404(userBlog.ThemePath)
	}

	// 分页的话, 需要分页信息, totalPage, curPage
	page := c.GetPage()
	pageInfo, blogs := blogService.ListBlogs(userId, notebookId, page, userBlog.PerPageSize, userBlog.SortField, userBlog.IsAsc)
	blogs2 := c.fixBlogs(blogs)
	c.RenderArgs["posts"] = blogs2

	c.setPaging(pageInfo)

	c.RenderArgs["curCateTitle"] = notebook.Title
	c.RenderArgs["curCateId"] = notebookId
	cateUrl := c.RenderArgs["cateUrl"].(string)
	c.RenderArgs["pagingBaseUrl"] = cateUrl + "/" + notebookId
	c.RenderArgs["curIsCate"] = true

	return c.render("cate.html", userBlog.ThemePath)
}

func (c Blog) Index(userIdOrEmail string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath);
		}
	}()
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}

	// 用户id为空, 转至博客平台
	if userIdOrEmail == "" {
		userIdOrEmail = configService.GetAdminUsername()
	}
	var userInfo info.User
	if userId != "" {
		userInfo = userService.GetUserInfoByAny(userId)
	} else {
		userInfo = userService.GetUserInfoByAny(userIdOrEmail)
	}
	userId = userInfo.UserId.Hex()

	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, userInfo); !ok {
		return c.e404(userBlog.ThemePath) // 404 TODO 使用用户的404
	}

	// 分页的话, 需要分页信息, totalPage, curPage
	page := c.GetPage()
	pageInfo, blogs := blogService.ListBlogs(userId, "", page, userBlog.PerPageSize, userBlog.SortField, userBlog.IsAsc)
	blogs2 := c.fixBlogs(blogs)
	c.RenderArgs["posts"] = blogs2

	c.setPaging(pageInfo)
	c.RenderArgs["pagingBaseUrl"] = c.RenderArgs["indexUrl"]

	c.RenderArgs["curIsIndex"] = true

	return c.render("index.html", userBlog.ThemePath)
}

// 详情
// 为了上一篇, 下一篇
func (c Blog) fixNote(note info.Note) map[string]interface{} {
	if note.NoteId == "" {
		return nil
	}
	return map[string]interface{}{
		"NoteId":      note.NoteId.Hex(),
		"Title":       note.Title,
		"CreatedTime": note.CreatedTime,
		"UpdatedTime": note.UpdatedTime,
		"PublicTime":  note.PublicTime,
		"Desc":        note.Desc,
		"Tags":        note.Tags,
		"CommentNum":  note.CommentNum,
		"ReadNum":     note.ReadNum,
		"LikeNum":     note.LikeNum,
		"IsMarkdown":  note.IsMarkdown,
	}
}

func (c Blog) Post(noteId string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			Log(err)
			re = c.e404(userBlog.ThemePath);
		}
	}()
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}

	blogInfo := blogService.GetBlog(noteId)
	userInfo := userService.GetUserInfo(blogInfo.UserId.Hex())
	if userId != "" && userInfo.UserId.Hex() != userId {
		panic("error")
	}

	userId = userInfo.UserId.Hex()
	_, userBlog = c.blogCommon(userId, userBlog, info.User{})

	c.RenderArgs["post"] = c.fixBlog(blogInfo)
	// c.RenderArgs["userInfo"] = userInfo
	c.RenderArgs["curIsPost"] = true

	// 得到访问者id
	/*
	visitUserId := c.GetUserId()
	if visitUserId != "" {
		visitUserInfo := userService.GetUserInfo(visitUserId)
		c.RenderArgs["visitUserInfoJson"] = c.Json(visitUserInfo)
		c.RenderArgs["visitUserInfo"] = visitUserInfo
	} else {
		c.RenderArgs["visitUserInfoJson"] = "{}"
	}
	*/

	// 上一篇, 下一篇
	var baseTime interface{}
	if userBlog.SortField == "PublicTime" {
		baseTime = blogInfo.PublicTime
	} else if userBlog.SortField == "CreatedTime" {
		baseTime = blogInfo.CreatedTime
	} else if userBlog.SortField == "UpdatedTime" {
		baseTime = blogInfo.UpdatedTime
	} else {
		baseTime = blogInfo.Title
	}

	preNote, nextNote := blogService.PreNextBlog(userId, userBlog.SortField, userBlog.IsAsc, baseTime)
	c.RenderArgs["prePost"] = c.fixNote(preNote)
	c.RenderArgs["nextPost"] = c.fixNote(nextNote)
	Log("---------")
	return c.render("post.html", userBlog.ThemePath)
}

func (c Blog) Single(singleId string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath);
		}
	}()
	
	userId := ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}

	single := blogService.GetSingle(singleId)
	userInfo := userService.GetUserInfo(single.UserId.Hex())
	if userId != "" && userInfo.UserId.Hex() != userId {
		return c.e404(userBlog.ThemePath)
	}

	userId = userInfo.UserId.Hex()
	_, userBlog = c.blogCommon(userId, userBlog, info.User{})

	c.RenderArgs["single"] = map[string]interface{}{
		"SingleId":    single.SingleId.Hex(),
		"Title":       single.Title,
		"Content":     single.Content,
		"CreatedTime": single.CreatedTime,
		"UpdatedTime": single.UpdatedTime,
	}
	c.RenderArgs["curSingleId"] = singleId
	c.RenderArgs["curIsSingle"] = true

	return c.render("single.html", userBlog.ThemePath)
}

// 搜索
func (c Blog) Search(userIdOrEmail, keywords string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath);
		}
	}()
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
//	c.RenderArgs["userInfo"] = userInfo
	userId = userInfo.UserId.Hex()
	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, userInfo); !ok {
		return c.e404(userBlog.ThemePath)
	}

	page := c.GetPage()
	pageInfo, blogs := blogService.SearchBlog(keywords, userId, page, userBlog.PerPageSize, userBlog.SortField, userBlog.IsAsc)
	c.setPaging(pageInfo)

	c.RenderArgs["posts"] = c.fixBlogs(blogs)
	c.RenderArgs["keywords"] = keywords
	searchUrl, _ := c.RenderArgs["searchUrl"].(string)
	c.RenderArgs["pagingBaseUrl"] = searchUrl + "?keywords=" + keywords
	c.RenderArgs["curIsSearch"] = true

	return c.render("search.html", userBlog.ThemePath)
}

// 可以不要, 因为注册的时候已经把username设为email了
func (c Blog) setRenderUserInfo(userInfo info.User) {
	if userInfo.Username == "" {
		userInfo.Username = userInfo.Email
	}
	c.RenderArgs["userInfo"] = userInfo
}

//----------------
// 社交, 点赞, 评论

// 得到博客统计信息
func (c Blog) GetPostStat(noteId string) revel.Result {
	re := info.NewRe()
	re.Ok = true
	statInfo := blogService.GetBlogStat(noteId)
	re.Item = statInfo
	return c.RenderJson(re);
}

// jsonP
// 我是否点过赞? 得到我的信息
// 所有点赞的用户列表
// 各个评论中是否我也点过赞?
func (c Blog) GetLikes(noteId string, callback string) revel.Result {
	userId := c.GetUserId()
	result := map[string]interface{}{}
	isILikeIt := false
	if userId != ""  {
		isILikeIt = blogService.IsILikeIt(noteId, userId)
		result["visitUserInfo"] = userService.GetUserAndBlog(userId)
	}
	// 点赞用户列表
	likedUsers, hasMoreLikedUser := blogService.ListLikedUsers(noteId, false)

	re := info.NewRe()
	re.Ok = true
	result["isILikeIt"] = isILikeIt
	result["likedUsers"] = likedUsers
	result["hasMoreLikedUser"] = hasMoreLikedUser

	re.Item = result
	return c.RenderJsonP(callback, re)
}
func (c Blog) GetLikesAndComments(noteId, callback string) revel.Result {
	userId := c.GetUserId()
	result := map[string]interface{}{}

	// 我也点过?
	isILikeIt := false
	if userId != ""  {
		isILikeIt = blogService.IsILikeIt(noteId, userId)
		result["visitUserInfo"] = userService.GetUserAndBlog(userId)
	}
	
	// 点赞用户列表
	likedUsers, hasMoreLikedUser := blogService.ListLikedUsers(noteId, false)
	// 评论
	page := c.GetPage()
	pageInfo, comments, commentUserInfo := blogService.ListComments(userId, noteId, page, 15)

	re := info.NewRe()
	re.Ok = true
	result["isILikeIt"] = isILikeIt
	result["likedUsers"] = likedUsers
	result["hasMoreLikedUser"] = hasMoreLikedUser
	result["pageInfo"] = pageInfo
	result["comments"] = comments
	result["commentUserInfo"] = commentUserInfo
	re.Item = result
	return c.RenderJsonP(callback, re)
}

func (c Blog) IncReadNum(noteId string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.IncReadNum(noteId)
	return c.RenderJson(re)
}

// 点赞, 要用jsonp
func (c Blog) LikePost(noteId string, callback string) revel.Result {
	re := info.NewRe()
	userId := c.GetUserId()
	re.Ok, re.Item = blogService.LikeBlog(noteId, userId)
	return c.RenderJsonP(callback, re)
}
func (c Blog) GetComments(noteId string) revel.Result {
	// 评论
	userId := c.GetUserId()
	page := c.GetPage()
	pageInfo, comments, commentUserInfo := blogService.ListComments(userId, noteId, page, 15)
	re := info.NewRe()
	re.Ok = true
	result := map[string]interface{}{}
	result["pageInfo"] = pageInfo
	result["comments"] = comments
	result["commentUserInfo"] = commentUserInfo
	re.Item = result

	return c.RenderJson(re)
}
// jsonp
func (c Blog) DeleteComment(noteId, commentId string, callback string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.DeleteComment(noteId, commentId, c.GetUserId())
	return c.RenderJsonP(callback, re)
}
// jsonp
func (c Blog) CommentPost(noteId, content, toCommentId string, callback string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Item = blogService.Comment(noteId, toCommentId, c.GetUserId(), content)
	return c.RenderJsonP(callback, re)
}
// jsonp
func (c Blog) LikeComment(commentId string, callback string) revel.Result {
	re := info.NewRe()
	ok, isILikeIt, num := blogService.LikeComment(commentId, c.GetUserId())
	re.Ok = ok
	re.Item = bson.M{"IsILikeIt": isILikeIt, "Num": num}
	return c.RenderJsonP(callback, re)
}

// 显示分类的最近博客, jsonp
func (c Blog) ListCateLatest(notebookId, callback string) revel.Result {
	if notebookId == "" {
		return c.e404("")
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
		return c.e404(userBlog.ThemePath)
	}
	if userId != "" && userId != notebook.UserId.Hex() {
		return c.e404(userBlog.ThemePath)
	}
	userId = notebook.UserId.Hex()

	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, info.User{}); !ok {
		return c.e404(userBlog.ThemePath)
	}

	// 分页的话, 需要分页信息, totalPage, curPage
	page := 1
	_, blogs := blogService.ListBlogs(userId, notebookId, page, 5, userBlog.SortField, userBlog.IsAsc)
	re := info.NewRe()
	re.Ok = true
	re.List = blogs
	return c.RenderJsonP(callback, re)
}
