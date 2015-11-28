package member

import (
	"github.com/revel/revel"
)

// admin 首页

type MemberIndex struct {
	MemberBaseController
}

// admin 主页
func (c MemberIndex) Index() revel.Result {
	c.SetUserInfo()
	c.RenderArgs["title"] = c.Message("Leanote Member Center")

	c.RenderArgs["countNote"] = noteService.CountNote(c.GetUserId())
	c.RenderArgs["countBlog"] = noteService.CountBlog(c.GetUserId())

	c.SetLocale()

	return c.RenderTemplate("member/index.html")
}

// 模板
func (c MemberIndex) T(t string) revel.Result {
	c.RenderArgs["str"] = configService.GlobalStringConfigs
	c.RenderArgs["arr"] = configService.GlobalArrayConfigs
	c.RenderArgs["map"] = configService.GlobalMapConfigs
	c.RenderArgs["arrMap"] = configService.GlobalArrMapConfigs
	return c.RenderTemplate("admin/" + t + ".html")
}

func (c MemberIndex) GetView(view string) revel.Result {
	return c.RenderTemplate("admin/" + view)
}
