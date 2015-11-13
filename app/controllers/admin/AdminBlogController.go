package admin

import (
	"github.com/revel/revel"
	//	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/info"
)

// admin 首页

type AdminBlog struct {
	AdminBaseController
}

// admin 主页
func (c AdminBlog) Index(sorter, keywords string) revel.Result {
	pageNumber := c.GetPage()
	sorterField, isAsc := c.getSorter("CreatedTime", false, []string{"title", "userId", "isRecommed", "createdTime"})
	pageInfo, blogs := blogService.ListAllBlogs("", "", keywords, false, pageNumber, userPageSize, sorterField, isAsc)
	c.RenderArgs["pageInfo"] = pageInfo
	c.RenderArgs["blogs"] = blogs
	c.RenderArgs["keywords"] = keywords
	return c.RenderTemplate("admin/blog/list.html")
}

func (c AdminBlog) SetRecommend(noteId string, recommend bool) revel.Result {
	re := info.NewRe()
	re.Ok = blogService.SetRecommend(noteId, recommend)
	return c.RenderJson(re)
}
