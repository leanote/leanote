package controllers

import (
	"github.com/revel/revel"
	"leanote/app/info"
	. "leanote/app/lea"
)

// 首页

type Index struct {
	BaseController
}

// leanote展示页, 没有登录的, 或已登录明确要进该页的
func (c Index) Index() revel.Result {
	c.SetUserInfo()
	c.RenderArgs["title"] = "leanote"
	c.RenderArgs["openRegister"] = openRegister
	c.SetLocale()

	return c.RenderTemplate("home/index.html")
}

// 建议
func (c Index) Suggestion(addr, suggestion string) revel.Result {
	re := info.NewRe()
	re.Ok = suggestionService.AddSuggestion(info.Suggestion{Addr: addr, UserId: c.GetObjectUserId(), Suggestion: suggestion})

	// 发给我
	go func() {
		SendToLeanote("建议", "建议", "UserId: "+c.GetUserId()+" <br /> Suggestions: "+suggestion)
	}()

	return c.RenderJson(re)
}
