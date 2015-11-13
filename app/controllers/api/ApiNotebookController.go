package api

import (
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	//	"io/ioutil"
)

// 笔记本API

type ApiNotebook struct {
	ApiBaseContrller
}

// 从Notebook -> ApiNotebook
func (c ApiNotebook) fixNotebooks(notebooks []info.Notebook) []info.ApiNotebook {
	if notebooks == nil {
		return nil
	}
	apiNotebooks := make([]info.ApiNotebook, len(notebooks))
	for i, notebook := range notebooks {
		apiNotebooks[i] = c.fixNotebook(&notebook)
	}
	return apiNotebooks
}
func (c ApiNotebook) fixNotebook(notebook *info.Notebook) info.ApiNotebook {
	if notebook == nil {
		return info.ApiNotebook{}
	}
	return info.ApiNotebook{
		NotebookId:       notebook.NotebookId,
		UserId:           notebook.UserId,
		ParentNotebookId: notebook.ParentNotebookId,
		Seq:              notebook.Seq,
		Title:            notebook.Title,
		UrlTitle:         notebook.UrlTitle,
		IsBlog:           notebook.IsBlog,
		CreatedTime:      notebook.CreatedTime,
		UpdatedTime:      notebook.UpdatedTime,
		Usn:              notebook.Usn,
		IsDeleted:        notebook.IsDeleted,
	}
}

// 获取同步的笔记本
// [OK]
// > afterUsn的笔记
// 返回 {ChunkHighUsn: 本下最大的usn, 借此可以知道是否还有, Notebooks: []}
func (c ApiNotebook) GetSyncNotebooks(afterUsn, maxEntry int) revel.Result {
	if maxEntry == 0 {
		maxEntry = 100
	}
	notebooks := notebookService.GeSyncNotebooks(c.getUserId(), afterUsn, maxEntry)
	return c.RenderJson(c.fixNotebooks(notebooks))
}

// 得到用户的所有笔记本
// [OK]
// info.SubNotebooks
func (c ApiNotebook) GetNotebooks() revel.Result {
	notebooks := notebookService.GeSyncNotebooks(c.getUserId(), 0, 99999)
	return c.RenderJson(c.fixNotebooks(notebooks))
}

// 添加notebook
// [OK]
func (c ApiNotebook) AddNotebook(title, parentNotebookId string, seq int) revel.Result {
	notebook := info.Notebook{NotebookId: bson.NewObjectId(),
		Title:  title,
		Seq:    seq,
		UserId: bson.ObjectIdHex(c.getUserId())}
	if parentNotebookId != "" && bson.IsObjectIdHex(parentNotebookId) {
		notebook.ParentNotebookId = bson.ObjectIdHex(parentNotebookId)
	}
	re := info.NewRe()
	re.Ok, notebook = notebookService.AddNotebook(notebook)
	if !re.Ok {
		return c.RenderJson(re)
	}
	return c.RenderJson(c.fixNotebook(&notebook))
}

// 修改笔记
// [OK]
func (c ApiNotebook) UpdateNotebook(notebookId, title, parentNotebookId string, seq, usn int) revel.Result {
	re := info.NewApiRe()

	ok, msg, notebook := notebookService.UpdateNotebookApi(c.getUserId(), notebookId, title, parentNotebookId, seq, usn)
	if !ok {
		re.Ok = false
		re.Msg = msg
		return c.RenderJson(re)
	}
	LogJ(notebook)
	return c.RenderJson(c.fixNotebook(&notebook))
}

// 删除笔记本
// [OK]
func (c ApiNotebook) DeleteNotebook(notebookId string, usn int) revel.Result {
	re := info.NewApiRe()
	re.Ok, re.Msg = notebookService.DeleteNotebookForce(c.getUserId(), notebookId, usn)
	return c.RenderJson(re)
}
