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
	
	re.Ok = configService.UpdateGlobalArrayConfig(c.GetUserId(), "recommendTags", strings.Split(recommendTags, ","))
	re.Ok = configService.UpdateGlobalArrayConfig(c.GetUserId(), "newTags", strings.Split(newTags, ","))
	
	return c.RenderJson(re)
}

// 共享设置
func (c AdminSetting) ShareNote(registerSharedUserId string, 
	registerSharedNotebookPerms, registerSharedNotePerms []int, 
	registerSharedNotebookIds, registerSharedNoteIds, registerCopyNoteIds []string) revel.Result {
	
	re := info.NewRe()
	re.Ok, re.Msg = configService.UpdateShareNoteConfig(registerSharedUserId, registerSharedNotebookPerms, registerSharedNotePerms, registerSharedNotebookIds, registerSharedNoteIds, registerCopyNoteIds);
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
	
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "demoUserId", userInfo.UserId.Hex())
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "demoUsername", demoUsername)
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "demoPassword", demoPassword)
	
	return c.RenderJson(re)
}

// ToImage
// 长微博的bin路径phantomJs
func (c AdminSetting) ToImage() revel.Result {
	c.RenderArgs["toImageBinPath"] = configService.GetGlobalStringConfig("toImageBinPath")
	return c.RenderTemplate("admin/setting/toImage.html");
}
func (c AdminSetting) DoToImage(toImageBinPath string) revel.Result {
	re := info.NewRe()
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "toImageBinPath", toImageBinPath)
	return c.RenderJson(re)
}

// SubDomain
func (c AdminSetting) SubDomain() revel.Result {
	c.RenderArgs["str"] = configService.GlobalStringConfigs
	c.RenderArgs["arr"] = configService.GlobalArrayConfigs
	
	c.RenderArgs["noteSubDomain"] = configService.GetGlobalStringConfig("noteSubDomain")
	c.RenderArgs["blogSubDomain"] = configService.GetGlobalStringConfig("blogSubDomain")
	c.RenderArgs["leaSubDomain"] = configService.GetGlobalStringConfig("leaSubDomain")
	
	return c.RenderTemplate("admin/setting/subDomain.html");
}
func (c AdminSetting) DoSubDomain(noteSubDomain, blogSubDomain, leaSubDomain, blackSubDomains, allowCustomDomain, blackCustomDomains string) revel.Result {
	re := info.NewRe()
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "noteSubDomain", noteSubDomain)
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "blogSubDomain", blogSubDomain)
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "leaSubDomain", leaSubDomain)
	
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "allowCustomDomain", allowCustomDomain)
	re.Ok = configService.UpdateGlobalArrayConfig(c.GetUserId(), "blackSubDomains", strings.Split(blackSubDomains, ","))
	re.Ok = configService.UpdateGlobalArrayConfig(c.GetUserId(), "blackCustomDomains", strings.Split(blackCustomDomains, ","))
	
	return c.RenderJson(re)
}

func (c AdminSetting) OpenRegister(openRegister string) revel.Result {
	re := info.NewRe()
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "openRegister", openRegister)
	return c.RenderJson(re)
}

func (c AdminSetting) Mongodb(mongodumpPath, mongorestorePath string) revel.Result {
	re := info.NewRe()
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "mongodumpPath", mongodumpPath)
	re.Ok = configService.UpdateGlobalStringConfig(c.GetUserId(), "mongorestorePath", mongorestorePath)

	return c.RenderJson(re)
}