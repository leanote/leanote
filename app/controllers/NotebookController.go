package controllers

import (
	"encoding/json"
	"github.com/leanote/leanote/app/info"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	//	. "github.com/leanote/leanote/app/lea"
	//	"io/ioutil"
)

type Notebook struct {
	BaseController
}

func (c Notebook) Index(notebook info.Notebook, i int, name string) revel.Result {
	return c.RenderJson(notebook)
}

// 得到用户的所有笔记本
func (c Notebook) GetNotebooks() revel.Result {
	re := notebookService.GetNotebooks(c.GetUserId())
	return c.RenderJson(re)
}

func (c Notebook) DeleteNotebook(notebookId string) revel.Result {
	re, msg := notebookService.DeleteNotebook(c.GetUserId(), notebookId)
	return c.RenderJson(info.Re{Ok: re, Msg: msg})
}

// 添加notebook
func (c Notebook) AddNotebook(notebookId, title, parentNotebookId string) revel.Result {
	notebook := info.Notebook{NotebookId: bson.ObjectIdHex(notebookId),
		Title:  title,
		Seq:    -1,
		UserId: c.GetObjectUserId()}
	if parentNotebookId != "" {
		notebook.ParentNotebookId = bson.ObjectIdHex(parentNotebookId)
	}

	re, notebook := notebookService.AddNotebook(notebook)

	if re {
		return c.RenderJson(notebook)
	} else {
		return c.RenderJson(false)
	}
}

// 修改标题
func (c Notebook) UpdateNotebookTitle(notebookId, title string) revel.Result {
	return c.RenderJson(notebookService.UpdateNotebookTitle(notebookId, c.GetUserId(), title))
}

// 排序
// 无用
// func (c Notebook) SortNotebooks(notebookId2Seqs map[string]int) revel.Result {
// 	return c.RenderJson(notebookService.SortNotebooks(c.GetUserId(), notebookId2Seqs))
// }

// 调整notebooks, 可能是排序, 可能是移动到其它笔记本下
type DragNotebooksInfo struct {
	CurNotebookId    string
	ParentNotebookId string
	Siblings         []string
}

// 传过来的data是JSON.stringfy数据
func (c Notebook) DragNotebooks(data string) revel.Result {
	info := DragNotebooksInfo{}
	json.Unmarshal([]byte(data), &info)

	return c.RenderJson(notebookService.DragNotebooks(c.GetUserId(), info.CurNotebookId, info.ParentNotebookId, info.Siblings))
}

// 设置notebook <-> blog
func (c Notebook) SetNotebook2Blog(notebookId string, isBlog bool) revel.Result {
	re := notebookService.ToBlog(c.GetUserId(), notebookId, isBlog)
	return c.RenderJson(re)
}
