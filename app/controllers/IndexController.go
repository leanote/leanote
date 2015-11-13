package controllers

import (
	"github.com/leanote/leanote/app/info"
	"github.com/revel/revel"
	//	. "github.com/leanote/leanote/app/lea"
)

// 首页

type Index struct {
	BaseController
}

func (c Index) Default() revel.Result {
	if configService.HomePageIsAdminsBlog() {
		blog := Blog{c.BaseController}
		return blog.Index(configService.GetAdminUsername())
	}
	return c.Index()
}

// leanote展示页, 没有登录的, 或已登录明确要进该页的
func (c Index) Index() revel.Result {
	c.SetUserInfo()
	c.RenderArgs["title"] = "leanote"
	c.RenderArgs["openRegister"] = configService.GlobalStringConfigs["openRegister"]
	c.SetLocale()

	return c.RenderTemplate("home/index.html")
}

// 建议
func (c Index) Suggestion(addr, suggestion string) revel.Result {
	re := info.NewRe()
	re.Ok = suggestionService.AddSuggestion(info.Suggestion{Addr: addr, UserId: c.GetObjectUserId(), Suggestion: suggestion})

	// 发给我
	go func() {
		emailService.SendEmail("leanote@leanote.com", "建议", "UserId: "+c.GetUserId()+" <br /> Suggestions: "+suggestion)
	}()

	return c.RenderJson(re)
}
