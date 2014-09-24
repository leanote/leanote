package admin

import (
	"github.com/revel/revel"
//	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/info"
	"strings"
)

// admin 首页

type AdminSetting struct {
	AdminBaseController
}

// email配置
func (c AdminSetting) Email() revel.Result {
	return nil
}

// blog标签设置
func (c AdminSetting) Blog() revel.Result {
	recommendTags := configService.GetGlobalArrayConfig("recommendTags")
	newTags := configService.GetGlobalArrayConfig("newTags")
	c.RenderArgs["recommendTags"] = strings.Join(recommendTags, ",")
	c.RenderArgs["newTags"] = strings.Join(newTags, ",")
	return c.RenderTemplate("admin/setting/blog.html");
}
func (c AdminSetting) DoBlogTag(recommendTags, newTags string) revel.Result {
	re := info.NewRe()
	
	re.Ok = configService.UpdateUserArrayConfig(c.GetUserId(), "recommendTags", strings.Split(recommendTags, ","))
	re.Ok = configService.UpdateUserArrayConfig(c.GetUserId(), "newTags", strings.Split(newTags, ","))
	
	return c.RenderJson(re)
}

// demo
// blog标签设置
func (c AdminSetting) Demo() revel.Result {
	c.RenderArgs["demoUsername"] = configService.GetGlobalStringConfig("demoUsername")
	c.RenderArgs["demoPassword"] = configService.GetGlobalStringConfig("demoPassword")
	return c.RenderTemplate("admin/setting/demo.html");
}
func (c AdminSetting) DoDemo(demoUsername, demoPassword string) revel.Result {
	re := info.NewRe()
	
	userInfo := authService.Login(demoUsername, demoPassword)
	if userInfo.UserId == "" {
		re.Msg = "The User is Not Exists";
		return c.RenderJson(re)
	}
	
	re.Ok = configService.UpdateUserStringConfig(c.GetUserId(), "demoUserId", userInfo.UserId.Hex())
	re.Ok = configService.UpdateUserStringConfig(c.GetUserId(), "demoUsername", demoUsername)
	re.Ok = configService.UpdateUserStringConfig(c.GetUserId(), "demoPassword", demoPassword)
	
	return c.RenderJson(re)
}



