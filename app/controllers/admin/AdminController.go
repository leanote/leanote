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

	c.ViewArgs["title"] = "leanote"
	c.SetLocale()

	c.ViewArgs["countUser"] = userService.CountUser()
	c.ViewArgs["countNote"] = noteService.CountNote("")
	c.ViewArgs["countBlog"] = noteService.CountBlog("")

	return c.RenderTemplate("admin/index.html")
}

// 模板
func (c Admin) T(t string) revel.Result {
	c.ViewArgs["str"] = configService.GlobalStringConfigs
	c.ViewArgs["arr"] = configService.GlobalArrayConfigs
	c.ViewArgs["map"] = configService.GlobalMapConfigs
	c.ViewArgs["arrMap"] = configService.GlobalArrMapConfigs
	c.ViewArgs["version"] = configService.GetVersion()
	return c.RenderTemplate("admin/" + t + ".html")
}

func (c Admin) GetView(view string) revel.Result {
	return c.RenderTemplate("admin/" + view)
}
