package controllers

import (
	"github.com/revel/revel"
//	"encoding/json"
	"github.com/leanote/leanote/app/info"
	"labix.org/v2/mgo/bson"
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
func (c Notebook) AddNotebook(notebookId, title string) revel.Result {
	notebook := info.Notebook{NotebookId: bson.ObjectIdHex(notebookId), 
		Title: title,
		Seq: -1,
		UserId: c.GetObjectUserId()}
	re := notebookService.AddNotebook(notebook)
	
	if(re) {
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
func (c Notebook) SortNotebooks(notebookId2Seqs map[string]int) revel.Result {
	return c.RenderJson(notebookService.SortNotebooks(c.GetUserId(), notebookId2Seqs))
}