package controllers

import (
	"github.com/revel/revel"
	"strings"
	//	"encoding/json"
	"fmt"
	"github.com/leanote/leanote/app/info"
	// . "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/lea/blog"
	"gopkg.in/mgo.v2/bson"
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
	if c.ViewArgs["isPreview"] != nil {
		themePath2 := c.ViewArgs["themePath"]
		if themePath2 == nil {
			return c.E404()
		}
		isPreview = true
		themePath = themePath2.(string)
		c.setPreviewUrl()

		// 因为common的themeInfo是从UserBlog.ThemeId来取的, 所以这里要fugai下
		c.ViewArgs["themeInfo"] = c.ViewArgs["themeInfoPreview"]
	}
	return blog.RenderTemplate(templateName, c.ViewArgs, revel.BasePath+"/"+themePath, isPreview)
}

// 404
func (c Blog) e404(themePath string) revel.Result {
	// 不知道是谁的404, 则用系统的404
	if themePath == "" {
		return c.E404()
	}
	return c.render("404.html", themePath)
}

// 二级域名或自定义域名
// life.leanote.com
// lealife.com
func (c Blog) domain() (ok bool, userBlog info.UserBlog) {
	host := c.Request.Host // a.cc.com:9000
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
	userIdOrEmail := userId
	username := c.GetUsername()
	if username != "" {
		userIdOrEmail = username
	}
	themeId := c.Session["themeId"]
	theme := themeService.GetTheme(userId, themeId)

	// siteUrl := configService.GetSiteUrl()
	blogUrl := "/preview" // blog.leanote.com

	indexUrl = blogUrl + "/" + userIdOrEmail
	cateUrl = blogUrl + "/cate/" + userIdOrEmail // /notebookId

	postUrl = blogUrl + "/post/" + userIdOrEmail        // /xxxxx
	searchUrl = blogUrl + "/search/" + userIdOrEmail    // blog.leanote.com/search/userId
	singleUrl = blogUrl + "/single/" + userIdOrEmail    // blog.leanote.com/single/singleId
	archiveUrl = blogUrl + "/archives/" + userIdOrEmail // blog.leanote.com/archive/userId
	tagsUrl = blogUrl + "/tags/" + userIdOrEmail        // blog.leanote.com/archive/userId

	c.ViewArgs["indexUrl"] = indexUrl
	c.ViewArgs["cateUrl"] = cateUrl
	c.ViewArgs["postUrl"] = postUrl
	c.ViewArgs["searchUrl"] = searchUrl
	c.ViewArgs["singleUrl"] = singleUrl // 单页
	c.ViewArgs["archiveUrl"] = archiveUrl
	c.ViewArgs["archivesUrl"] = archiveUrl // 别名
	c.ViewArgs["tagsUrl"] = tagsUrl
	c.ViewArgs["tagPostsUrl"] = blogUrl + "/tag/" + userIdOrEmail
	c.ViewArgs["tagUrl"] = c.ViewArgs["tagPostsUrl"]

	// themeBaseUrl 本theme的路径url, 可以加载js, css, images之类的
	c.ViewArgs["themeBaseUrl"] = "/" + theme.Path
}

// 各种地址设置
func (c Blog) setUrl(userBlog info.UserBlog, userInfo info.User) {
	// 主页 http://leanote.com/blog/life or http://blog.leanote.com/life or http:// xxxx.leanote.com or aa.com
	// host := c.Request.Request.Host
	// var staticUrl = configService.GetUserUrl(strings.Split(host, ":")[0])
	// staticUrl == host, 为保证同源!!! 只有host, http://leanote.com, http://blog/leanote.com
	// life.leanote.com, lealife.com
	siteUrl := configService.GetSiteUrl()
	blogUrls := blogService.GetBlogUrls(&userBlog, &userInfo)
	// 分类
	// 搜索
	// 查看
	c.ViewArgs["siteUrl"] = siteUrl
	c.ViewArgs["indexUrl"] = blogUrls.IndexUrl
	c.ViewArgs["cateUrl"] = blogUrls.CateUrl
	c.ViewArgs["postUrl"] = blogUrls.PostUrl
	c.ViewArgs["searchUrl"] = blogUrls.SearchUrl
	c.ViewArgs["singleUrl"] = blogUrls.SingleUrl // 单页
	c.ViewArgs["archiveUrl"] = blogUrls.ArchiveUrl
	c.ViewArgs["archivesUrl"] = blogUrls.ArchiveUrl // 别名
	c.ViewArgs["tagsUrl"] = blogUrls.TagsUrl
	c.ViewArgs["tagPostsUrl"] = blogUrls.TagPostsUrl
	c.ViewArgs["tagUrl"] = blogUrls.TagPostsUrl // 别名

	// themeBaseUrl 本theme的路径url, 可以加载js, css, images之类的
	c.ViewArgs["themeBaseUrl"] = "/" + userBlog.ThemePath

	// 其它static js
	c.ViewArgs["jQueryUrl"] = "/js/jquery-1.9.0.min.js"

	c.ViewArgs["prettifyJsUrl"] = "/js/google-code-prettify/prettify.js"
	c.ViewArgs["prettifyCssUrl"] = "/js/google-code-prettify/prettify.css"

	c.ViewArgs["blogCommonJsUrl"] = "/public/blog/js/common.js"

	c.ViewArgs["shareCommentCssUrl"] = "/public/blog/css/share_comment.css"
	c.ViewArgs["shareCommentJsUrl"] = "/public/blog/js/share_comment.js"

	c.ViewArgs["fontAwesomeUrl"] = "/css/font-awesome-4.2.0/css/font-awesome.css"

	c.ViewArgs["bootstrapCssUrl"] = "/css/bootstrap.css"
	c.ViewArgs["bootstrapJsUrl"] = "/js/bootstrap-min.js"
}

// 笔记本分类
// cates = [{title:"xxx", cateId: "xxxx"}, {}]
func (c Blog) getCateUrlTitle(n *info.Notebook) string {
	if n.UrlTitle != "" {
		return n.UrlTitle
	}
	return n.NotebookId.Hex()
}
func (c Blog) getCates(userBlog info.UserBlog) {
	notebooks := blogService.ListBlogNotebooks(userBlog.UserId.Hex())
	notebooksMap := map[string]info.Notebook{}
	for _, each := range notebooks {
		notebooksMap[each.NotebookId.Hex()] = each
	}

	var i = 0
	cates := make([]*info.Cate, len(notebooks))

	// 先要保证已有的是正确的排序
	cateIds := userBlog.CateIds
	has := map[string]bool{} // cateIds中有的
	cateMap := map[string]*info.Cate{}
	if cateIds != nil && len(cateIds) > 0 {
		for _, cateId := range cateIds {
			if n, ok := notebooksMap[cateId]; ok {
				parentNotebookId := ""
				if n.ParentNotebookId != "" {
					parentNotebookId = n.ParentNotebookId.Hex()
				}
				cates[i] = &info.Cate{Title: n.Title, UrlTitle: c.getCateUrlTitle(&n), CateId: n.NotebookId.Hex(), ParentCateId: parentNotebookId}
				cateMap[cates[i].CateId] = cates[i]
				i++
				has[cateId] = true
			}
		}
	}

	// 之后添加没有排序的
	for _, n := range notebooks {
		id := n.NotebookId.Hex()
		if !has[id] {
			parentNotebookId := ""
			if n.ParentNotebookId != "" {
				parentNotebookId = n.ParentNotebookId.Hex()
			}
			cates[i] = &info.Cate{Title: n.Title, UrlTitle: c.getCateUrlTitle(&n), CateId: id, ParentCateId: parentNotebookId}
			cateMap[cates[i].CateId] = cates[i]
			i++
		}
	}

	//	LogJ(">>")
	//	LogJ(cates)

	// 建立层级
	hasParent := map[string]bool{} // 有父的cate
	for _, cate := range cates {
		parentCateId := cate.ParentCateId
		if parentCateId != "" {
			if parentCate, ok := cateMap[parentCateId]; ok {
				//				Log("________")
				//				LogJ(parentCate)
				//				LogJ(cate)
				if parentCate.Children == nil {
					parentCate.Children = []*info.Cate{cate}
				} else {
					parentCate.Children = append(parentCate.Children, cate)
				}
				hasParent[cate.CateId] = true
			}
		}
	}

	// 得到没有父的cate, 作为第一级cate
	catesTree := []*info.Cate{}
	for _, cate := range cates {
		if !hasParent[cate.CateId] {
			catesTree = append(catesTree, cate)
		}
	}

	c.ViewArgs["cates"] = cates
	c.ViewArgs["catesTree"] = catesTree
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
	c.ViewArgs["singles"] = singles
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
		"DisqusId":    userBlog.DisqusId,
		"ThemeId":     userBlog.ThemeId,
		"SubDomain":   userBlog.SubDomain,
		"Domain":      userBlog.Domain,
	}
	c.ViewArgs["blogInfo"] = blogInfo
}

func (c Blog) setPaging(pageInfo info.Page) {
	c.ViewArgs["paging"] = pageInfo
}

// 公共
func (c Blog) blogCommon(userId string, userBlog info.UserBlog, userInfo info.User) (ok bool, ub info.UserBlog) {
	if userInfo.UserId == "" {
		userInfo = userService.GetUserInfoByAny(userId)
		if userInfo.UserId == "" {
			return false, userBlog
		}
	}
	//	c.ViewArgs["userInfo"] = userInfo

	// 最新笔记
	_, recentBlogs := blogService.ListBlogs(userId, "", 1, 5, userBlog.SortField, userBlog.IsAsc)
	c.ViewArgs["recentPosts"] = blogService.FixBlogs(recentBlogs)
	c.ViewArgs["latestPosts"] = c.ViewArgs["recentPosts"]
	c.ViewArgs["tags"] = blogService.GetBlogTags(userId)

	// 语言, url地址
	c.SetLocale()

	// 得到博客设置信息
	if userBlog.UserId == "" {
		userBlog = blogService.GetUserBlog(userId)
	}
	c.setBlog(userBlog, userInfo)
	//	c.ViewArgs["userBlog"] = userBlog

	// 分类导航
	c.getCates(userBlog)

	// 单页导航
	c.getSingles(userId)

	c.setUrl(userBlog, userInfo)

	// 当前分类Id, 全设为""
	c.ViewArgs["curCateId"] = ""
	c.ViewArgs["curSingleId"] = ""

	// 得到主题信息
	themeInfo := themeService.GetThemeInfo(userBlog.ThemeId.Hex(), userBlog.Style)
	c.ViewArgs["themeInfo"] = themeInfo

	//	Log(">>")
	//	Log(userBlog.Style)
	//	Log(userBlog.ThemeId.Hex())

	return true, userBlog
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
	_, userBlog = c.blogCommon(userId, userBlog, userInfo)

	return c.e404(userBlog.ThemePath)
}

func (c Blog) Tags(userIdOrEmail string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath)
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

	c.ViewArgs["curIsTags"] = true
	tags := blogService.GetBlogTags(userId)
	c.ViewArgs["tags"] = tags
	return c.render("tags.html", userBlog.ThemePath)
}

// 标签的文章页
func (c Blog) Tag(userIdOrEmail, tag string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath)
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

	c.ViewArgs["curIsTagPosts"] = true
	c.ViewArgs["curTag"] = tag
	page := c.GetPage()
	pageInfo, blogs := blogService.SearchBlogByTags([]string{tag}, userId, page, userBlog.PerPageSize, userBlog.SortField, userBlog.IsAsc)
	c.setPaging(pageInfo)

	c.ViewArgs["posts"] = blogService.FixBlogs(blogs)
	tagPostsUrl := c.ViewArgs["tagPostsUrl"].(string)
	c.ViewArgs["pagingBaseUrl"] = tagPostsUrl + "/" + tag

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
			re = c.e404(userBlog.ThemePath)
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
	c.ViewArgs["archives"] = arcs

	c.ViewArgs["curIsArchive"] = true
	if notebookId != "" {
		notebook := notebookService.GetNotebookById(notebookId)
		c.ViewArgs["curCateTitle"] = notebook.Title
		c.ViewArgs["curCateId"] = notebookId
	}
	c.ViewArgs["curYear"] = year
	c.ViewArgs["curMonth"] = month

	return c.render("archive.html", userBlog.ThemePath)
}

// 进入某个用户的博客
var blogPageSize = 5
var searchBlogPageSize = 30

// 分类 /cate/xxxxxxxx?notebookId=1212
func (c Blog) Cate(userIdOrEmail string, notebookId string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			re = c.e404(userBlog.ThemePath)
		}
	}()

	userId, userInfo := c.userIdOrEmail(hasDomain, userBlog, userIdOrEmail)
	notebookId2 := notebookId
	var notebook info.Notebook
	if userId == "" { // 证明没有userIdOrEmail, 只有singleId, 那么直接查
		notebook = notebookService.GetNotebookById(notebookId)
		userId = notebook.UserId.Hex()
	} else {
		notebook = notebookService.GetNotebookByUserIdAndUrlTitle(userId, notebookId)
		notebookId2 = notebook.NotebookId.Hex()
	}
	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, userInfo); !ok {
		return c.e404(userBlog.ThemePath) // 404 TODO 使用用户的404
	}
	if !notebook.IsBlog {
		panic("")
	}

	// 分页的话, 需要分页信息, totalPage, curPage
	page := c.GetPage()
	pageInfo, blogs := blogService.ListBlogs(userId, notebookId2, page, userBlog.PerPageSize, userBlog.SortField, userBlog.IsAsc)
	blogs2 := blogService.FixBlogs(blogs)
	c.ViewArgs["posts"] = blogs2

	c.setPaging(pageInfo)

	c.ViewArgs["curCateTitle"] = notebook.Title
	c.ViewArgs["curCateId"] = notebookId2
	cateUrl := c.ViewArgs["cateUrl"].(string)
	c.ViewArgs["pagingBaseUrl"] = cateUrl + "/" + notebookId
	c.ViewArgs["curIsCate"] = true

	return c.render("cate.html", userBlog.ThemePath)
}

func (c Blog) userIdOrEmail(hasDomain bool, userBlog info.UserBlog, userIdOrEmail string) (userId string, userInfo info.User) {
	userId = ""
	if hasDomain {
		userId = userBlog.UserId.Hex()
	}
	if userId != "" {
		userInfo = userService.GetUserInfoByAny(userId)
	} else {
		if userIdOrEmail != "" {
			userInfo = userService.GetUserInfoByAny(userIdOrEmail)
		} else {
			return
		}
	}
	userId = userInfo.UserId.Hex()
	return
}

func (c Blog) Index(userIdOrEmail string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath)
		}
	}()
	// 用户id为空, 则是admin用户的博客
	if userIdOrEmail == "" {
		userIdOrEmail = configService.GetAdminUsername()
	}
	userId, userInfo := c.userIdOrEmail(hasDomain, userBlog, userIdOrEmail)
	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, userInfo); !ok {
		return c.e404(userBlog.ThemePath) // 404 TODO 使用用户的404
	}

	// 分页的话, 需要分页信息, totalPage, curPage
	page := c.GetPage()
	pageInfo, blogs := blogService.ListBlogs(userId, "", page, userBlog.PerPageSize, userBlog.SortField, userBlog.IsAsc)
	blogs2 := blogService.FixBlogs(blogs)
	c.ViewArgs["posts"] = blogs2

	c.setPaging(pageInfo)
	c.ViewArgs["pagingBaseUrl"] = c.ViewArgs["indexUrl"]

	c.ViewArgs["curIsIndex"] = true

	return c.render("index.html", userBlog.ThemePath)
}

func (c Blog) Post(userIdOrEmail, noteId string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			// Log(err)
			re = c.e404(userBlog.ThemePath)
		}
	}()

	userId, userInfo := c.userIdOrEmail(hasDomain, userBlog, userIdOrEmail)
	var blogInfo info.BlogItem
	if userId == "" { // 证明没有userIdOrEmail, 只有singleId, 那么直接查
		blogInfo = blogService.GetBlog(noteId)
		userId = blogInfo.UserId.Hex()
	} else {
		blogInfo = blogService.GetBlogByIdAndUrlTitle(userId, noteId)
	}
	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, userInfo); !ok {
		return c.e404(userBlog.ThemePath) // 404 TODO 使用用户的404
	}
	if blogInfo.NoteId == "" {
		return c.e404(userBlog.ThemePath) // 404 TODO 使用用户的404
	}

	post := blogService.FixBlog(blogInfo)
	c.ViewArgs["post"] = post
	// c.ViewArgs["userInfo"] = userInfo
	c.ViewArgs["curIsPost"] = true

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

	prePost, nextPost := blogService.PreNextBlog(userId, userBlog.SortField, userBlog.IsAsc, post.NoteId, baseTime)
	if prePost.NoteId != "" {
		c.ViewArgs["prePost"] = prePost
	}
	if nextPost.NoteId != "" {
		c.ViewArgs["nextPost"] = nextPost
	}
	return c.render("post.html", userBlog.ThemePath)
}

func (c Blog) Single(userIdOrEmail, singleId string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath)
		}
	}()

	userId, userInfo := c.userIdOrEmail(hasDomain, userBlog, userIdOrEmail)
	var single info.BlogSingle
	if userId == "" { // 证明没有userIdOrEmail, 只有singleId, 那么直接查
		single = blogService.GetSingle(singleId)
		userId = single.UserId.Hex()
	} else {
		single = blogService.GetSingleByUserIdAndUrlTitle(userId, singleId)
	}
	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, userInfo); !ok {
		return c.e404(userBlog.ThemePath) // 404 TODO 使用用户的404
	}
	if single.SingleId == "" {
		panic("")
	}

	c.ViewArgs["single"] = map[string]interface{}{
		"SingleId":    single.SingleId.Hex(),
		"Title":       single.Title,
		"UrlTitle":    single.UrlTitle,
		"Content":     single.Content,
		"CreatedTime": single.CreatedTime,
		"UpdatedTime": single.UpdatedTime,
	}
	c.ViewArgs["curSingleId"] = single.SingleId.Hex()
	c.ViewArgs["curIsSingle"] = true

	return c.render("single.html", userBlog.ThemePath)
}

// 搜索
func (c Blog) Search(userIdOrEmail, keywords string) (re revel.Result) {
	// 自定义域名
	hasDomain, userBlog := c.domain()
	defer func() {
		if err := recover(); err != nil {
			re = c.e404(userBlog.ThemePath)
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
	//	c.ViewArgs["userInfo"] = userInfo
	userId = userInfo.UserId.Hex()
	var ok = false
	if ok, userBlog = c.blogCommon(userId, userBlog, userInfo); !ok {
		return c.e404(userBlog.ThemePath)
	}

	page := c.GetPage()
	pageInfo, blogs := blogService.SearchBlog(keywords, userId, page, userBlog.PerPageSize, userBlog.SortField, userBlog.IsAsc)
	c.setPaging(pageInfo)

	c.ViewArgs["posts"] = blogService.FixBlogs(blogs)
	c.ViewArgs["keywords"] = keywords
	searchUrl, _ := c.ViewArgs["searchUrl"].(string)
	c.ViewArgs["pagingBaseUrl"] = searchUrl + "?keywords=" + keywords
	c.ViewArgs["curIsSearch"] = true

	return c.render("search.html", userBlog.ThemePath)
}

// 可以不要, 因为注册的时候已经把username设为email了
func (c Blog) setRenderUserInfo(userInfo info.User) {
	if userInfo.Username == "" {
		userInfo.Username = userInfo.Email
	}
	c.ViewArgs["userInfo"] = userInfo
}

//----------------
// 社交, 点赞, 评论

// 得到博客统计信息
func (c Blog) GetPostStat(noteId string) revel.Result {
	re := info.NewRe()
	re.Ok = true
	statInfo := blogService.GetBlogStat(noteId)
	re.Item = statInfo
	return c.RenderJSON(re)
}

// jsonP
// 我是否点过赞? 得到我的信息
// 所有点赞的用户列表
// 各个评论中是否我也点过赞?
func (c Blog) GetLikes(noteId string, callback string) revel.Result {
	userId := c.GetUserId()
	result := map[string]interface{}{}
	isILikeIt := false
	if userId != "" {
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
	return c.RenderJSONP(callback, re)
}
func (c Blog) GetLikesAndComments(noteId, callback string) revel.Result {
	userId := c.GetUserId()
	result := map[string]interface{}{}

	// 我也点过?
	isILikeIt := false
	if userId != "" {
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
	return c.RenderJSONP(callback, re)
}

func (c Blog) IncReadNum(noteId string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.IncReadNum(noteId)
	return c.RenderJSON(re)
}

// 点赞, 要用jsonp
func (c Blog) LikePost(noteId string, callback string) revel.Result {
	re := info.NewRe()
	userId := c.GetUserId()
	re.Ok, re.Item = blogService.LikeBlog(noteId, userId)
	return c.RenderJSONP(callback, re)
}
func (c Blog) GetComments(noteId string, callback string) revel.Result {
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

	if callback != "" {
		return c.RenderJSONP(callback, result)
	}

	return c.RenderJSON(re)
}

// jsonp
func (c Blog) DeleteComment(noteId, commentId string, callback string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.DeleteComment(noteId, commentId, c.GetUserId())
	return c.RenderJSONP(callback, re)
}

// jsonp
func (c Blog) CommentPost(noteId, content, toCommentId string, callback string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Item = blogService.Comment(noteId, toCommentId, c.GetUserId(), content)
	return c.RenderJSONP(callback, re)
}

// jsonp
func (c Blog) LikeComment(commentId string, callback string) revel.Result {
	re := info.NewRe()
	ok, isILikeIt, num := blogService.LikeComment(commentId, c.GetUserId())
	re.Ok = ok
	re.Item = bson.M{"IsILikeIt": isILikeIt, "Num": num}
	return c.RenderJSONP(callback, re)
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
	return c.RenderJSONP(callback, re)
}
