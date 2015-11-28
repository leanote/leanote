package member

import (
	"fmt"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"io/ioutil"
	"os"
	"strings"
	"time"
	//	"github.com/leanote/leanote/app/lea/blog"
)

// 博客管理

type MemberBlog struct {
	MemberBaseController
}

func (c MemberBlog) common() info.UserBlog {
	userId := c.GetUserId()
	userInfo := userService.GetUserInfo(userId)
	c.RenderArgs["userInfo"] = userInfo

	// 得到博客设置信息
	c.RenderArgs["allowCustomDomain"] = configService.GetGlobalStringConfig("allowCustomDomain")

	userBlog := blogService.GetUserBlog(userId)
	c.RenderArgs["userBlog"] = userBlog

	c.SetUserInfo()
	c.SetLocale()
	return userBlog
}

// 得到sorterField 和 isAsc
// okSorter = ['email', 'username']
func (c MemberBlog) getSorter(sorterField string, isAsc bool, okSorter []string) (string, bool) {
	sorter := ""
	c.Params.Bind(&sorter, "sorter")
	if sorter == "" {
		return sorterField, isAsc
	}

	// sorter形式 email-up, email-down
	s2 := strings.Split(sorter, "-")
	if len(s2) != 2 {
		return sorterField, isAsc
	}

	// 必须是可用的sorter
	if okSorter != nil && len(okSorter) > 0 {
		if !InArray(okSorter, s2[0]) {
			return sorterField, isAsc
		}
	}

	sorterField = strings.Title(s2[0])
	if s2[1] == "up" {
		isAsc = true
	} else {
		isAsc = false
	}
	c.RenderArgs["sorter"] = sorter
	return sorterField, isAsc
}

// 博客列表
var userPageSize = 15

func (c MemberBlog) Index(sorter, keywords string) revel.Result {
	userId := c.GetUserId()
	userInfo := userService.GetUserInfo(userId)
	c.RenderArgs["userInfo"] = userInfo

	c.RenderArgs["title"] = c.Message("Posts")
	pageNumber := c.GetPage()
	sorterField, isAsc := c.getSorter("CreatedTime", false, []string{"title", "urlTitle", "updatedTime", "publicTime", "createdTime"})
	pageInfo, blogs := blogService.ListAllBlogs(c.GetUserId(), "", keywords, false, pageNumber, userPageSize, sorterField, isAsc)
	c.RenderArgs["pageInfo"] = pageInfo
	c.RenderArgs["blogs"] = blogs
	c.RenderArgs["keywords"] = keywords

	userAndBlog := userService.GetUserAndBlog(c.GetUserId())
	c.RenderArgs["userAndBlog"] = userAndBlog

	c.common()

	return c.RenderTemplate("member/blog/list.html")
}

// 修改笔记的urlTitle
func (c MemberBlog) UpdateBlogUrlTitle(noteId, urlTitle string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Item = blogService.UpateBlogUrlTitle(c.GetUserId(), noteId, urlTitle)
	return c.RenderJson(re)
}

// 修改笔记的urlTitle
func (c MemberBlog) UpdateBlogAbstract(noteId string) revel.Result {
	c.common()
	c.RenderArgs["title"] = c.Message("Update Post Abstract")
	note := noteService.GetNoteAndContent(noteId, c.GetUserId())
	if !note.Note.IsBlog {
		return c.E404()
	}
	c.RenderArgs["note"] = note
	c.RenderArgs["noteId"] = noteId
	return c.RenderTemplate("member/blog/update_abstract.html")
}
func (c MemberBlog) DoUpdateBlogAbstract(noteId, imgSrc, desc, abstract string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.UpateBlogAbstract(c.GetUserId(), noteId, imgSrc, desc, abstract)
	return c.RenderJson(re)
}

// 基本信息设置
func (c MemberBlog) Base() revel.Result {
	c.common()
	c.RenderArgs["title"] = c.Message("Blog Base Info")
	return c.RenderTemplate("member/blog/base.html")
}
func (c MemberBlog) Comment() revel.Result {
	c.common()
	c.RenderArgs["title"] = c.Message("Comment")
	return c.RenderTemplate("member/blog/comment.html")
}

func (c MemberBlog) Paging() revel.Result {
	c.common()
	c.RenderArgs["title"] = c.Message("Paging")
	return c.RenderTemplate("member/blog/paging.html")
}

func (c MemberBlog) Cate() revel.Result {
	userBlog := c.common()
	c.RenderArgs["title"] = c.Message("Category")

	notebooks := blogService.ListBlogNotebooks(c.GetUserId())
	notebooksMap := map[string]info.Notebook{}
	for _, each := range notebooks {
		notebooksMap[each.NotebookId.Hex()] = each
	}

	var i = 0
	notebooks2 := make([]info.Notebook, len(notebooks))

	// 先要保证已有的是正确的排序
	cateIds := userBlog.CateIds
	has := map[string]bool{} // cateIds中有的
	if cateIds != nil && len(cateIds) > 0 {
		for _, cateId := range cateIds {
			if n, ok := notebooksMap[cateId]; ok {
				notebooks2[i] = n
				i++
				has[cateId] = true
			}
		}
	}
	// 之后
	for _, each := range notebooks {
		id := each.NotebookId.Hex()
		if !has[id] {
			notebooks2[i] = each
			i++
		}
	}
	c.RenderArgs["notebooks"] = notebooks2

	return c.RenderTemplate("member/blog/cate.html")
}

// 修改分类排序
func (c MemberBlog) UpateCateIds(cateIds []string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.UpateCateIds(c.GetUserId(), cateIds)
	return c.RenderJson(re)
}

func (c MemberBlog) UpdateCateUrlTitle(cateId, urlTitle string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Item = blogService.UpateCateUrlTitle(c.GetUserId(), cateId, urlTitle)
	return c.RenderJson(re)
}

// 保存之, 包含增加与保存
func (c MemberBlog) DoAddOrUpdateSingle(singleId, title, content string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.AddOrUpdateSingle(c.GetUserId(), singleId, title, content)
	return c.RenderJson(re)
}
func (c MemberBlog) AddOrUpdateSingle(singleId string) revel.Result {
	c.common()
	c.RenderArgs["title"] = c.Message("Add Single")
	c.RenderArgs["singleId"] = singleId
	if singleId != "" {
		c.RenderArgs["title"] = c.Message("Update Single")
		c.RenderArgs["single"] = blogService.GetSingle(singleId)
	}
	return c.RenderTemplate("member/blog/add_single.html")
}
func (c MemberBlog) SortSingles(singleIds []string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.SortSingles(c.GetUserId(), singleIds)
	return c.RenderJson(re)
}

func (c MemberBlog) DeleteSingle(singleId string) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.DeleteSingle(c.GetUserId(), singleId)
	return c.RenderJson(re)
}

// 修改页面标题
func (c MemberBlog) UpdateSingleUrlTitle(singleId, urlTitle string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Item = blogService.UpdateSingleUrlTitle(c.GetUserId(), singleId, urlTitle)
	return c.RenderJson(re)
}

func (c MemberBlog) Single() revel.Result {
	c.common()
	c.RenderArgs["title"] = c.Message("Single")
	c.RenderArgs["singles"] = blogService.GetSingles(c.GetUserId())

	return c.RenderTemplate("member/blog/single.html")
}

// 主题
func (c MemberBlog) Theme() revel.Result {
	c.common()
	activeTheme, otherThemes := themeService.GetUserThemes(c.GetUserId())
	c.RenderArgs["activeTheme"] = activeTheme
	c.RenderArgs["otherThemes"] = otherThemes

	c.RenderArgs["optionThemes"] = themeService.GetDefaultThemes()

	c.RenderArgs["title"] = c.Message("Theme")
	return c.RenderTemplate("member/blog/theme.html")
}

// 编辑主题
var baseTpls = []string{"header.html", "footer.html", "index.html", "cate.html", "search.html", "post.html", "single.html", "tags.html", "tag_posts.html", "archive.html", "share_comment.html", "404.html", "theme.json", "style.css", "blog.js"}

func (c MemberBlog) UpdateTheme(themeId string, isNew int) revel.Result {
	// 查看用户是否有该theme, 若没有则复制default之
	// 得到主题的文件列表
	userBlog := blogService.GetUserBlog(c.GetUserId())
	if themeId == "" {
		_, themeId = themeService.NewThemeForFirst(userBlog)
		return c.Redirect("/member/blog/updateTheme?themeId=" + themeId)
	}

	c.common()
	c.RenderArgs["title"] = c.Message("Update Theme")
	c.RenderArgs["isNew"] = isNew

	// 先复制之
	c.RenderArgs["themeId"] = themeId

	// 得到脚本目录
	userId := c.GetUserId()

	theme := themeService.GetTheme(userId, themeId)
	if theme.ThemeId == "" {
		return c.E404()
	}
	c.RenderArgs["theme"] = theme

	path := revel.BasePath + "/" + theme.Path

	tpls := ListDir(path)
	myTpls := make([]string, len(baseTpls))
	tplMap := map[string]bool{}
	for i, t := range baseTpls {
		myTpls[i] = t
		tplMap[t] = true
	}
	// 得到没有的tpls
	for _, t := range tpls {
		if t == "images" {
			continue
		}
		if !tplMap[t] {
			myTpls = append(myTpls, t)
		}
	}

	c.RenderArgs["myTpls"] = myTpls

	return c.RenderTemplate("member/blog/update_theme.html")
}

// 得到文件内容
func (c MemberBlog) GetTplContent(themeId string, filename string) revel.Result {
	re := info.NewRe()
	re.Ok = true
	re.Item = themeService.GetTplContent(c.GetUserId(), themeId, filename)

	return c.RenderJson(re)
}
func (c MemberBlog) UpdateTplContent(themeId, filename, content string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = themeService.UpdateTplContent(c.GetUserId(), themeId, filename, content)
	return c.RenderRe(re)
}

func (c MemberBlog) DeleteTpl(themeId, filename string) revel.Result {
	re := info.NewRe()
	re.Ok = themeService.DeleteTpl(c.GetUserId(), themeId, filename)
	return c.RenderJson(re)
}

func (c MemberBlog) ListThemeImages(themeId string) revel.Result {
	re := info.NewRe()
	userId := c.GetUserId()
	path := themeService.GetThemeAbsolutePath(userId, themeId) + "/images"
	os.MkdirAll(path, 0755)
	images := ListDir(path)
	re.List = images
	re.Ok = true
	return c.RenderJson(re)
}

func (c MemberBlog) DeleteThemeImage(themeId, filename string) revel.Result {
	re := info.NewRe()
	path := themeService.GetThemeAbsolutePath(c.GetUserId(), themeId) + "/images/" + filename
	re.Ok = DeleteFile(path)
	return c.RenderJson(re)
}

// 上传主题图片
func (c MemberBlog) UploadThemeImage(themeId string) revel.Result {
	re := c.uploadImage(themeId)
	c.RenderArgs["fileUrlPath"] = re.Id
	c.RenderArgs["resultCode"] = re.Code
	c.RenderArgs["resultMsg"] = re.Msg
	return c.RenderTemplate("file/blog_logo.html")
}
func (c MemberBlog) uploadImage(themeId string) (re info.Re) {
	var fileId = ""
	var resultCode = 0     // 1表示正常
	var resultMsg = "内部错误" // 错误信息
	var Ok = false

	defer func() {
		re.Id = fileId // 只是id, 没有其它信息
		re.Code = resultCode
		re.Msg = resultMsg
		re.Ok = Ok
	}()

	file, handel, err := c.Request.FormFile("file")
	if err != nil {
		return re
	}
	defer file.Close()
	// 生成上传路径
	dir := themeService.GetThemeAbsolutePath(c.GetUserId(), themeId) + "/images"
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return re
	}
	// 生成新的文件名
	filename := handel.Filename

	var ext string

	_, ext = SplitFilename(filename)
	if ext != ".gif" && ext != ".jpg" && ext != ".png" && ext != ".bmp" && ext != ".jpeg" {
		resultMsg = "不是图片"
		return re
	}

	filename = filename
	data, err := ioutil.ReadAll(file)
	if err != nil {
		LogJ(err)
		return re
	}

	// > 2M?
	if len(data) > 5*1024*1024 {
		resultCode = 0
		resultMsg = "图片大于2M"
		return re
	}

	toPath := dir + "/" + filename
	err = ioutil.WriteFile(toPath, data, 0777)
	if err != nil {
		LogJ(err)
		return re
	}
	TransToGif(toPath, 0, true)
	resultCode = 1
	resultMsg = "上传成功!"

	return re
}

//
// 使用主题
func (c MemberBlog) ActiveTheme(themeId string) revel.Result {
	re := info.NewRe()
	re.Ok = themeService.ActiveTheme(c.GetUserId(), themeId)
	return c.RenderJson(re)
}

// 删除主题
func (c MemberBlog) DeleteTheme(themeId string) revel.Result {
	re := info.NewRe()
	re.Ok = themeService.DeleteTheme(c.GetUserId(), themeId)
	return c.RenderJson(re)
}

// 管理员公开主题
func (c MemberBlog) PublicTheme(themeId string) revel.Result {
	re := info.NewRe()
	re.Ok = themeService.PublicTheme(c.GetUserId(), themeId)
	return c.RenderJson(re)
}

// 导出
func (c MemberBlog) ExportTheme(themeId string) revel.Result {
	re := info.NewRe()
	var path string
	re.Ok, path = themeService.ExportTheme(c.GetUserId(), themeId)
	if !re.Ok {
		return c.RenderText("error...")
	}
	fw, err := os.Open(path)
	if err != nil {
		return c.RenderText("error")
	}
	return c.RenderBinary(fw, GetFilename(path), revel.Attachment, time.Now()) // revel.Attachment
}

// 导入主题
func (c MemberBlog) ImportTheme() revel.Result {
	re := info.NewRe()

	file, handel, err := c.Request.FormFile("file")
	if err != nil {
		re.Msg = fmt.Sprintf("%v", err)
		return c.RenderJson(re)
	}

	defer file.Close()
	// 生成上传路径
	userId := c.GetUserId()
	dir := revel.BasePath + "/public/upload/" + userId + "/tmp"
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		re.Msg = fmt.Sprintf("%v", err)
		return c.RenderJson(re)
	}
	// 生成新的文件名
	filename := handel.Filename

	var ext string
	_, ext = SplitFilename(filename)
	if ext != ".zip" {
		re.Msg = "Please upload zip file"
		return c.RenderJson(re)
	}

	filename = filename
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return c.RenderJson(re)
	}

	// > 10M?
	if len(data) > 10*1024*1024 {
		re.Msg = "File is big than 10M"
		return c.RenderJson(re)
	}

	toPath := dir + "/" + filename
	err = ioutil.WriteFile(toPath, data, 0777)
	if err != nil {
		re.Msg = fmt.Sprintf("%v", err)
		return c.RenderJson(re)
	}

	// 上传好后, 增加之
	re.Ok, re.Msg = themeService.ImportTheme(c.GetUserId(), toPath)
	return c.RenderRe(re)
}

// 安装
func (c MemberBlog) InstallTheme(themeId string) revel.Result {
	re := info.NewRe()
	re.Ok = themeService.InstallTheme(c.GetUserId(), themeId)
	return c.RenderJson(re)
}

// 新建主题
func (c MemberBlog) NewTheme() revel.Result {
	_, themeId := themeService.NewTheme(c.GetUserId())
	return c.Redirect("/member/blog/updateTheme?isNew=1&themeId=" + themeId)
}

//-----------
//
func (c MemberBlog) SetUserBlogBase(userBlog info.UserBlogBase) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.UpdateUserBlogBase(c.GetUserId(), userBlog)
	return c.RenderJson(re)
}
func (c MemberBlog) SetUserBlogComment(userBlog info.UserBlogComment) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.UpdateUserBlogComment(c.GetUserId(), userBlog)
	return c.RenderJson(re)
}
func (c MemberBlog) SetUserBlogStyle(userBlog info.UserBlogStyle) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.UpdateUserBlogStyle(c.GetUserId(), userBlog)
	return c.RenderJson(re)
}

func (c MemberBlog) SetUserBlogPaging(perPageSize int, sortField string, isAsc bool) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = blogService.UpdateUserBlogPaging(c.GetUserId(), perPageSize, sortField, isAsc)
	return c.RenderRe(re)
}
