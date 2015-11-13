package controllers

import (
	"github.com/revel/revel"
	//	"encoding/json"
	//	"gopkg.in/mgo.v2/bson"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	//	"github.com/leanote/leanote/app/types"
	//	"io/ioutil"
	//	"fmt"
)

type Share struct {
	BaseController
}

// 添加共享note
func (c Share) AddShareNote(noteId string, emails []string, perm int) revel.Result {
	status := make(map[string]info.Re, len(emails))
	// 自己不能给自己添加共享
	myEmail := c.GetEmail()
	for _, email := range emails {
		if email == "" {
			continue
		}
		if myEmail != email {
			ok, msg, userId := shareService.AddShareNote(noteId, perm, c.GetUserId(), email)
			status[email] = info.Re{Ok: ok, Msg: msg, Id: userId}

		} else {
			status[email] = info.Re{Ok: false, Msg: "不能分享给自己"}
		}
	}

	return c.RenderJson(status)
}

// 添加共享notebook
func (c Share) AddShareNotebook(notebookId string, emails []string, perm int) revel.Result {
	status := make(map[string]info.Re, len(emails))
	// 自己不能给自己添加共享
	myEmail := c.GetEmail()
	for _, email := range emails {
		if email == "" {
			continue
		}
		if myEmail != email {
			ok, msg, userId := shareService.AddShareNotebook(notebookId, perm, c.GetUserId(), email)
			status[email] = info.Re{Ok: ok, Msg: msg, Id: userId}
		} else {
			status[email] = info.Re{Ok: false, Msg: "不能分享给自己"}
		}
	}

	return c.RenderJson(status)
}

// 得到notes
// userId 该userId分享给我的
func (c Share) ListShareNotes(notebookId, userId string) revel.Result {
	// 表示是默认笔记本, 不是某个特定notebook的共享
	var notes []info.ShareNoteWithPerm
	if notebookId == "" {
		notes = shareService.ListShareNotes(c.GetUserId(), userId, c.GetPage(), pageSize, defaultSortField, false)
		return c.RenderJson(notes)
	} else {
		// 有notebookId的
		return c.RenderJson(shareService.ListShareNotesByNotebookId(notebookId, c.GetUserId(), userId, c.GetPage(), pageSize, defaultSortField, false))
	}
}

// 得到内容
// sharedUserId 是谁的笔记
func (c Share) GetShareNoteContent(noteId, sharedUserId string) revel.Result {
	noteContent := shareService.GetShareNoteContent(noteId, c.GetUserId(), sharedUserId)
	return c.RenderJson(noteContent)
}

// 查看note的分享信息
// 分享给了哪些用户和权限
// ShareNotes表 userId = me, noteId = ...
// 还要查看该note的notebookId分享的信息
func (c Share) ListNoteShareUserInfo(noteId string) revel.Result {
	note := noteService.GetNote(noteId, c.GetUserId())

	noteShareUserInfos := shareService.ListNoteShareUserInfo(noteId, c.GetUserId())
	c.RenderArgs["noteOrNotebookShareUserInfos"] = noteShareUserInfos

	c.RenderArgs["noteOrNotebookShareGroupInfos"] = shareService.GetNoteShareGroups(noteId, c.GetUserId())

	c.RenderArgs["isNote"] = true
	c.RenderArgs["noteOrNotebookId"] = note.NoteId.Hex()
	c.RenderArgs["title"] = note.Title

	return c.RenderTemplate("share/note_notebook_share_user_infos.html")
}
func (c Share) ListNotebookShareUserInfo(notebookId string) revel.Result {
	notebook := notebookService.GetNotebook(notebookId, c.GetUserId())

	notebookShareUserInfos := shareService.ListNotebookShareUserInfo(notebookId, c.GetUserId())
	c.RenderArgs["noteOrNotebookShareUserInfos"] = notebookShareUserInfos

	c.RenderArgs["noteOrNotebookShareGroupInfos"] = shareService.GetNotebookShareGroups(notebookId, c.GetUserId())
	LogJ(c.RenderArgs["noteOrNotebookShareGroupInfos"])

	c.RenderArgs["isNote"] = false
	c.RenderArgs["noteOrNotebookId"] = notebook.NotebookId.Hex()
	c.RenderArgs["title"] = notebook.Title

	return c.RenderTemplate("share/note_notebook_share_user_infos.html")
}

//------------
// 改变share note 权限
func (c Share) UpdateShareNotePerm(noteId string, perm int, toUserId string) revel.Result {
	return c.RenderJson(shareService.UpdateShareNotePerm(noteId, perm, c.GetUserId(), toUserId))
}

// 改变share notebook 权限
func (c Share) UpdateShareNotebookPerm(notebookId string, perm int, toUserId string) revel.Result {
	return c.RenderJson(shareService.UpdateShareNotebookPerm(notebookId, perm, c.GetUserId(), toUserId))
}

//---------------
// 删除share note
func (c Share) DeleteShareNote(noteId string, toUserId string) revel.Result {
	return c.RenderJson(shareService.DeleteShareNote(noteId, c.GetUserId(), toUserId))
}

// 删除share notebook
func (c Share) DeleteShareNotebook(notebookId string, toUserId string) revel.Result {
	return c.RenderJson(shareService.DeleteShareNotebook(notebookId, c.GetUserId(), toUserId))
}

// 删除share note, 被共享方删除
func (c Share) DeleteShareNoteBySharedUser(noteId string, fromUserId string) revel.Result {
	return c.RenderJson(shareService.DeleteShareNote(noteId, fromUserId, c.GetUserId()))
}

// 删除share notebook, 被共享方删除
func (c Share) DeleteShareNotebookBySharedUser(notebookId string, fromUserId string) revel.Result {
	return c.RenderJson(shareService.DeleteShareNotebook(notebookId, fromUserId, c.GetUserId()))
}

// 删除fromUserId分享给我的所有note, notebook
func (c Share) DeleteUserShareNoteAndNotebook(fromUserId string) revel.Result {
	return c.RenderJson(shareService.DeleteUserShareNoteAndNotebook(fromUserId, c.GetUserId()))
}

//-------------
// 用户组

// 将笔记分享给分组
func (c Share) AddShareNoteGroup(noteId, groupId string, perm int) revel.Result {
	re := info.NewRe()
	re.Ok = shareService.AddShareNoteGroup(c.GetUserId(), noteId, groupId, perm)
	return c.RenderJson(re)
}

// 删除
func (c Share) DeleteShareNoteGroup(noteId, groupId string) revel.Result {
	re := info.NewRe()
	re.Ok = shareService.DeleteShareNoteGroup(c.GetUserId(), noteId, groupId)
	return c.RenderJson(re)
}

// 更新, 也是一样, 先删后加
func (c Share) UpdateShareNoteGroupPerm(noteId, groupId string, perm int) revel.Result {
	return c.AddShareNoteGroup(noteId, groupId, perm)
}

//------

// 将笔记分享给分组
func (c Share) AddShareNotebookGroup(notebookId, groupId string, perm int) revel.Result {
	re := info.NewRe()
	re.Ok = shareService.AddShareNotebookGroup(c.GetUserId(), notebookId, groupId, perm)
	return c.RenderJson(re)
}

// 删除
func (c Share) DeleteShareNotebookGroup(notebookId, groupId string) revel.Result {
	re := info.NewRe()
	re.Ok = shareService.DeleteShareNotebookGroup(c.GetUserId(), notebookId, groupId)
	return c.RenderJson(re)
}

// 更新, 也是一样, 先删后加
func (c Share) UpdateShareNotebookGroupPerm(notebookId, groupId string, perm int) revel.Result {
	return c.AddShareNotebookGroup(notebookId, groupId, perm)
}
