package controllers

import (
	"github.com/revel/revel"
//	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/lea/html2image"
	"github.com/leanote/leanote/app/info"
//	"os"
//	"github.com/leanote/leanote/app/types"
//	"io/ioutil"
//	"fmt"
)

type Note struct {
	BaseController
}

// 笔记首页, 判断是否已登录
// 已登录, 得到用户基本信息(notebook, shareNotebook), 跳转到index.html中
// 否则, 转向登录页面
func (c Note) Index() revel.Result {
	c.SetLocale()
	
	userInfo := c.GetUserInfo()
	
	userId := userInfo.UserId.Hex()
	
	// 没有登录
	if userId == "" {
		return c.Redirect("/login")
	}
	
	c.RenderArgs["openRegister"] = openRegister
	
	// 已登录了, 那么得到所有信息
	notebooks := notebookService.GetNotebooks(userId)
	shareNotebooks, sharedUserInfos := shareService.GetShareNotebooks(userId)
	
	// 还需要按时间排序(DESC)得到notes
	notes := []info.Note{}
	noteContent := info.NoteContent{}
	if len(notebooks) > 0 {
//		_, notes = noteService.ListNotes(c.GetUserId(), "", false, c.GetPage(), pageSize, defaultSortField, false, false);
		// 变成最新
		_, notes = noteService.ListNotes(c.GetUserId(), "", false, c.GetPage(), 50, defaultSortField, false, false);
		if len(notes) > 0 {
			noteContent = noteService.GetNoteContent(notes[0].NoteId.Hex(), userId)
		}
	}
	// 当然, 还需要得到第一个notes的content
	//...
	
	c.RenderArgs["userInfo"] = userInfo
	c.RenderArgs["userInfoJson"] = c.Json(userInfo)
	c.RenderArgs["notebooks"] = c.Json(notebooks)
	c.RenderArgs["shareNotebooks"] = c.Json(shareNotebooks)
	c.RenderArgs["sharedUserInfos"] = c.Json(sharedUserInfos)
	
	c.RenderArgs["notes"] = c.Json(notes)
	c.RenderArgs["noteContentJson"] = c.Json(noteContent)
	c.RenderArgs["noteContent"] = noteContent.Content
	
	c.RenderArgs["tagsJson"] = c.Json(tagService.GetTags(c.GetUserId()))
	
	if isDev, _ := revel.Config.Bool("mode.dev"); isDev {
		return c.RenderTemplate("note/note-dev.html")
	} else {
		return c.RenderTemplate("note/note.html")
	}
}

// 首页, 判断是否已登录
// 已登录, 得到用户基本信息(notebook, shareNotebook), 跳转到index.html中
// 否则, 转向登录页面
func (c Note) ListNotes(notebookId string) revel.Result {
	_, notes := noteService.ListNotes(c.GetUserId(), notebookId, false, c.GetPage(), pageSize, defaultSortField, false, false);
	return c.RenderJson(notes)
}

// 得到trash
func (c Note) ListTrashNotes() revel.Result {
	_, notes := noteService.ListNotes(c.GetUserId(), "", true, c.GetPage(), pageSize, defaultSortField, false, false);
	return c.RenderJson(notes)
}

// 得到note和内容
func (c Note) GetNoteAndContent(noteId string) revel.Result {
	return c.RenderJson(noteService.GetNoteAndContent(noteId, c.GetUserId()))
}

// 得到内容
func (c Note) GetNoteContent(noteId string) revel.Result {
	noteContent := noteService.GetNoteContent(noteId, c.GetUserId())
	return c.RenderJson(noteContent)
}

// 更新note或content
// 肯定会传userId(谁的), NoteId
// 会传Title, Content, Tags, 一种或几种
type NoteOrContent struct {
	NotebookId string
	NoteId string
	UserId string
	Title string
	Desc string
	ImgSrc string
	Tags []string
	Content string
	Abstract string
	IsNew bool
	IsMarkdown bool
	FromUserId string // 为共享而新建
	IsBlog bool // 是否是blog, 更新note不需要修改, 添加note时才有可能用到, 此时需要判断notebook是否设为Blog
}
// 这里不能用json, 要用post
func (c Note) UpdateNoteOrContent(noteOrContent NoteOrContent) revel.Result {
	// 新添加note
	if noteOrContent.IsNew {
		userId := c.GetObjectUserId();
		myUserId := userId
		// 为共享新建?
		if noteOrContent.FromUserId != "" {
			userId = bson.ObjectIdHex(noteOrContent.FromUserId)
		}
		
		note := info.Note{UserId: userId, 
			NoteId: bson.ObjectIdHex(noteOrContent.NoteId), 
			NotebookId: bson.ObjectIdHex(noteOrContent.NotebookId), 
			Title: noteOrContent.Title, 
			Tags: noteOrContent.Tags,
			Desc: noteOrContent.Desc,
			ImgSrc: noteOrContent.ImgSrc,
			IsBlog: noteOrContent.IsBlog,
			IsMarkdown: noteOrContent.IsMarkdown,
		};
		noteContent := info.NoteContent{NoteId: note.NoteId, 
			UserId: userId, 
			IsBlog: note.IsBlog,
			Content: noteOrContent.Content, 
			Abstract: noteOrContent.Abstract};
		
		note = noteService.AddNoteAndContent(note, noteContent, myUserId)
		return c.RenderJson(note)
	}
	
	noteUpdate := bson.M{}
	needUpdateNote := false
	
	// Desc前台传来
	if c.Has("Desc") {
		needUpdateNote = true
		noteUpdate["Desc"] = noteOrContent.Desc;
	}
	if c.Has("ImgSrc") {
		needUpdateNote = true
		noteUpdate["ImgSrc"] = noteOrContent.ImgSrc;
	}
	if c.Has("Title") {
		needUpdateNote = true
		noteUpdate["Title"] = noteOrContent.Title;
	}
	
	if c.Has("Tags[]") {
		needUpdateNote = true
		noteUpdate["Tags"] = noteOrContent.Tags;
	}
	
	if needUpdateNote { 
		noteService.UpdateNote(noteOrContent.UserId, c.GetUserId(), 
			noteOrContent.NoteId, noteUpdate)
	}
	
	//-------------
	
	if c.Has("Content") {
		noteService.UpdateNoteContent(noteOrContent.UserId, c.GetUserId(), 
			noteOrContent.NoteId, noteOrContent.Content, noteOrContent.Abstract)
	}
	
	return c.RenderJson(true)
}

// 删除note/ 删除别人共享给我的笔记
// userId 是note.UserId
func (c Note) DeleteNote(noteId, userId string, isShared bool) revel.Result {
	if(!isShared) {
		return c.RenderJson(trashService.DeleteNote(noteId, c.GetUserId()));
	}
	
	return c.RenderJson(trashService.DeleteSharedNote(noteId, userId, c.GetUserId()));
}
// 删除trash
func (c Note) DeleteTrash(noteId string) revel.Result {
	return c.RenderJson(trashService.DeleteTrash(noteId, c.GetUserId()));
}
// 移动note
func (c Note) MoveNote(noteId, notebookId string) revel.Result {
	return c.RenderJson(noteService.MoveNote(noteId, notebookId, c.GetUserId()));
}
// 复制note
func (c Note) CopyNote(noteId, notebookId string) revel.Result {
	return c.RenderJson(noteService.CopyNote(noteId, notebookId, c.GetUserId()));
}
// 复制别人共享的笔记给我
func (c Note) CopySharedNote(noteId, notebookId, fromUserId string) revel.Result {
	return c.RenderJson(noteService.CopySharedNote(noteId, notebookId, fromUserId, c.GetUserId()));
}

//------------
// search
// 通过title搜索
func (c Note) SearchNote(key string) revel.Result {
	_, blogs := noteService.SearchNote(key, c.GetUserId(), c.GetPage(), pageSize, "UpdatedTime", false, false)
	return c.RenderJson(blogs)
}
// 通过tags搜索
func (c Note) SearchNoteByTags(tags []string) revel.Result {
	_, blogs := noteService.SearchNoteByTags(tags, c.GetUserId(), c.GetPage(), pageSize, "UpdatedTime", false)
	return c.RenderJson(blogs)
}

//-----------------
// html2image
func (c Note) Html2Image(noteId string) revel.Result {
	re := info.NewRe()
	userId := c.GetUserId()
	note := noteService.GetNote(noteId, userId)
	if note.NoteId == "" {
		return c.RenderJson(re)
	}
	content := noteService.GetNoteContent(noteId, userId)
	
	// path 判断是否需要重新生成之
	fileUrlPath := "/upload/" + userId + "/images/weibo"
	dir := revel.BasePath + "/public/" + fileUrlPath
	if !ClearDir(dir) {
		return c.RenderJson(re)
	}
	
	filename := note.NoteId.Hex() + ".png";
	path := dir + "/" + filename
	
	// 生成之
	html2image.ToImage(userId, c.GetUsername(), noteId, note.Title, content.Content, path)
	
	re.Ok = true
	re.Id = fileUrlPath + "/" + filename
	return c.RenderJson(re)
}