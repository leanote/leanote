package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	//	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"time"
	//	"strings"
)

// 用户组, 用户组用户管理

type GroupService struct {
}

// 添加分组
func (this *GroupService) AddGroup(userId, title string) (bool, info.Group) {
	group := info.Group{
		GroupId:     bson.NewObjectId(),
		UserId:      bson.ObjectIdHex(userId),
		Title:       title,
		CreatedTime: time.Now(),
	}
	return db.Insert(db.Groups, group), group
}

// 删除分组
// 判断是否有好友
func (this *GroupService) DeleteGroup(userId, groupId string) (ok bool, msg string) {
	/*
		if db.Has(db.GroupUsers, bson.M{"GroupId": bson.ObjectIdHex(groupId)}) {
			return false, "groupHasUsers"
		}
	*/
	if !this.isMyGroup(userId, groupId) {
		return false, "notMyGroup"
	}

	// 删除分组后, 需要删除所有用户分享到该组的笔记本, 笔记

	shareService.DeleteAllShareNotebookGroup(groupId)
	shareService.DeleteAllShareNoteGroup(groupId)

	db.DeleteAll(db.GroupUsers, bson.M{"GroupId": bson.ObjectIdHex(groupId)})
	return db.DeleteByIdAndUserId(db.Groups, groupId, userId), ""

	// TODO 删除分组后, 在shareNote, shareNotebook中也要删除
}

// 修改group标题
func (this *GroupService) UpdateGroupTitle(userId, groupId, title string) (ok bool) {
	return db.UpdateByIdAndUserIdField(db.Groups, groupId, userId, "Title", title)
}

// 得到用户的所有分组(包括下的所有用户)
func (this *GroupService) GetGroupsAndUsers(userId string) []info.Group {
	/*
		// 得到我的分组
		groups := []info.Group{}
		db.ListByQ(db.Groups, bson.M{"UserId": bson.ObjectIdHex(userId)}, &groups)
	*/
	// 我的分组, 及我所属的分组
	groups := this.GetGroupsContainOf(userId)

	// 得到其下的用户
	for i, group := range groups {
		group.Users = this.GetUsers(group.GroupId.Hex())
		groups[i] = group
	}
	return groups
}

// 仅仅得到所有分组
func (this *GroupService) GetGroups(userId string) []info.Group {
	// 得到分组s
	groups := []info.Group{}
	db.ListByQ(db.Groups, bson.M{"UserId": bson.ObjectIdHex(userId)}, &groups)
	return groups
}

// 得到我的和我所属组的ids
func (this *GroupService) GetMineAndBelongToGroupIds(userId string) []bson.ObjectId {
	// 所属组
	groupIds := this.GetBelongToGroupIds(userId)

	m := map[bson.ObjectId]bool{}
	for _, groupId := range groupIds {
		m[groupId] = true
	}

	// 我的组
	myGroups := this.GetGroups(userId)

	for _, group := range myGroups {
		if !m[group.GroupId] {
			groupIds = append(groupIds, group.GroupId)
		}
	}

	return groupIds
}

// 获取包含此用户的组对象数组
// 获取该用户所属组, 和我的组
func (this *GroupService) GetGroupsContainOf(userId string) []info.Group {
	// 我的组
	myGroups := this.GetGroups(userId)
	myGroupMap := map[bson.ObjectId]bool{}

	for _, group := range myGroups {
		myGroupMap[group.GroupId] = true
	}

	// 所属组
	groupIds := this.GetBelongToGroupIds(userId)

	groups := []info.Group{}
	db.ListByQ(db.Groups, bson.M{"_id": bson.M{"$in": groupIds}}, &groups)

	for _, group := range groups {
		if !myGroupMap[group.GroupId] {
			myGroups = append(myGroups, group)
		}
	}

	return myGroups
}

// 得到分组, shareService用
func (this *GroupService) GetGroup(userId, groupId string) info.Group {
	// 得到分组s
	group := info.Group{}
	db.GetByIdAndUserId(db.Groups, groupId, userId, &group)
	return group
}

// 得到某分组下的用户
func (this *GroupService) GetUsers(groupId string) []info.User {
	// 得到UserIds
	groupUsers := []info.GroupUser{}
	db.ListByQWithFields(db.GroupUsers, bson.M{"GroupId": bson.ObjectIdHex(groupId)}, []string{"UserId"}, &groupUsers)
	if len(groupUsers) == 0 {
		return nil
	}
	userIds := make([]bson.ObjectId, len(groupUsers))
	for i, each := range groupUsers {
		userIds[i] = each.UserId
	}
	// 得到userInfos
	return userService.ListUserInfosByUserIds(userIds)
}

// 得到我所属的所有分组ids
func (this *GroupService) GetBelongToGroupIds(userId string) []bson.ObjectId {
	// 得到UserIds
	groupUsers := []info.GroupUser{}
	db.ListByQWithFields(db.GroupUsers, bson.M{"UserId": bson.ObjectIdHex(userId)}, []string{"GroupId"}, &groupUsers)
	if len(groupUsers) == 0 {
		return nil
	}
	groupIds := make([]bson.ObjectId, len(groupUsers))
	for i, each := range groupUsers {
		groupIds[i] = each.GroupId
	}
	return groupIds
}

func (this *GroupService) isMyGroup(ownUserId, groupId string) (ok bool) {
	return db.Has(db.Groups, bson.M{"_id": bson.ObjectIdHex(groupId), "UserId": bson.ObjectIdHex(ownUserId)})
}

// 判断组中是否包含指定用户
func (this *GroupService) IsExistsGroupUser(userId, groupId string) (ok bool) {
	// 如果我拥有这个组, 那也行
	if this.isMyGroup(userId, groupId) {
		return true
	}
	return db.Has(db.GroupUsers, bson.M{"UserId": bson.ObjectIdHex(userId), "GroupId": bson.ObjectIdHex(groupId)})
}

// 为group添加用户
// 用户是否已存在?
func (this *GroupService) AddUser(ownUserId, groupId, userId string) (ok bool, msg string) {
	// groupId是否是ownUserId的?
	/*
		if !this.IsExistsGroupUser(ownUserId, groupId) {
			return false, "forbiddenNotMyGroup"
		}
	*/
	if !this.isMyGroup(ownUserId, groupId) {
		return false, "forbiddenNotMyGroup"
	}

	// 是否已存在
	if db.Has(db.GroupUsers, bson.M{"GroupId": bson.ObjectIdHex(groupId), "UserId": bson.ObjectIdHex(userId)}) {
		return false, "userExistsInGroup"
	}

	return db.Insert(db.GroupUsers, info.GroupUser{
		GroupUserId: bson.NewObjectId(),
		GroupId:     bson.ObjectIdHex(groupId),
		UserId:      bson.ObjectIdHex(userId),
		CreatedTime: time.Now(),
	}), ""
}

// 删除用户
func (this *GroupService) DeleteUser(ownUserId, groupId, userId string) (ok bool, msg string) {
	// groupId是否是ownUserId的?
	/*
		if !this.IsExistsGroupUser(ownUserId, groupId) {
			return false, "forbiddenNotMyGroup"
		}
	*/
	if !this.isMyGroup(ownUserId, groupId) {
		return false, "forbiddenNotMyGroup"
	}

	// 删除该用户分享到本组的笔记本, 笔记
	shareService.DeleteShareNotebookGroupWhenDeleteGroupUser(userId, groupId)
	shareService.DeleteShareNoteGroupWhenDeleteGroupUser(userId, groupId)

	return db.Delete(db.GroupUsers, bson.M{"GroupId": bson.ObjectIdHex(groupId), "UserId": bson.ObjectIdHex(userId)}), ""
}
