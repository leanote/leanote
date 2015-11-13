package admin

import (
	"github.com/revel/revel"
)

// admin 首页

type Admin struct {
	AdminBaseController
}

// admin 主页
func (c Admin) Index() revel.Result {
	c.SetUserInfo()

	c.RenderArgs["title"] = "leanote"
	c.SetLocale()

	c.RenderArgs["countUser"] = userService.CountUser()
	c.RenderArgs["countNote"] = noteService.CountNote("")
	c.RenderArgs["countBlog"] = noteService.CountBlog("")

	return c.RenderTemplate("admin/index.html")
}

// 模板
func (c Admin) T(t string) revel.Result {
	c.RenderArgs["str"] = configService.GlobalStringConfigs
	c.RenderArgs["arr"] = configService.GlobalArrayConfigs
	c.RenderArgs["map"] = configService.GlobalMapConfigs
	c.RenderArgs["arrMap"] = configService.GlobalArrMapConfigs
	c.RenderArgs["version"] = configService.GetVersion()
	return c.RenderTemplate("admin/" + t + ".html")
}

func (c Admin) GetView(view string) revel.Result {
	return c.RenderTemplate("admin/" + view)
}
