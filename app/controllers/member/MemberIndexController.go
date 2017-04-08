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
	c.ViewArgs["title"] = c.Message("Leanote Member Center")

	c.ViewArgs["countNote"] = noteService.CountNote(c.GetUserId())
	c.ViewArgs["countBlog"] = noteService.CountBlog(c.GetUserId())

	c.SetLocale()

	return c.RenderTemplate("member/index.html")
}

// 模板
func (c MemberIndex) T(t string) revel.Result {
	c.ViewArgs["str"] = configService.GlobalStringConfigs
	c.ViewArgs["arr"] = configService.GlobalArrayConfigs
	c.ViewArgs["map"] = configService.GlobalMapConfigs
	c.ViewArgs["arrMap"] = configService.GlobalArrMapConfigs
	return c.RenderTemplate("admin/" + t + ".html")
}

func (c MemberIndex) GetView(view string) revel.Result {
	return c.RenderTemplate("admin/" + view)
}
