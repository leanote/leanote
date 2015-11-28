package member

import (
	"github.com/leanote/leanote/app/info"
	"github.com/revel/revel"
)

// 分组管理
type MemberGroup struct {
	MemberBaseController
}

// 首页, 显示所有分组和用户
func (c MemberGroup) Index() revel.Result {
	c.SetUserInfo()
	c.SetLocale()
	c.RenderArgs["title"] = c.Message("My Group")
	c.RenderArgs["groups"] = groupService.GetGroupsAndUsers(c.GetUserId())
	return c.RenderTemplate("member/group/index.html")
}

// 添加分组
func (c MemberGroup) AddGroup(title string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Item = groupService.AddGroup(c.GetUserId(), title)
	return c.RenderJson(re)
}

func (c MemberGroup) UpdateGroupTitle(groupId, title string) revel.Result {
	re := info.NewRe()
	re.Ok = groupService.UpdateGroupTitle(c.GetUserId(), groupId, title)
	return c.RenderJson(re)
}

func (c MemberGroup) DeleteGroup(groupId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = groupService.DeleteGroup(c.GetUserId(), groupId)
	return c.RenderRe(re)
}

// 添加用户
func (c MemberGroup) AddUser(groupId, email string) revel.Result {
	re := info.NewRe()
	userInfo := userService.GetUserInfoByAny(email)
	if userInfo.UserId == "" {
		re.Msg = "userNotExists"
	} else {
		re.Ok, re.Msg = groupService.AddUser(c.GetUserId(), groupId, userInfo.UserId.Hex())
		re.Item = userInfo
	}
	return c.RenderRe(re)
}

func (c MemberGroup) DeleteUser(groupId, userId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = groupService.DeleteUser(c.GetUserId(), groupId, userId)
	return c.RenderRe(re)
}
