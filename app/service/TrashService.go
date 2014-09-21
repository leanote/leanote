package service

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
)

// 回收站
// 可以移到noteSerice中

// 不能删除notebook, 如果其下有notes!
// 这样回收站里只有note

// 删除笔记后(或删除笔记本)后入回收站
// 把note, notebook设个标记即可!
// 已经在trash里的notebook, note不能是共享!, 所以要删除共享

type TrashService struct {
}

//---------------------

// 删除note
// 应该放在回收站里
// 有trashService
func (this *TrashService) DeleteNote(noteId, userId string) bool {
	// 首先删除其共享
	if shareService.DeleteShareNoteAll(noteId, userId) {
		// 更新note isTrash = true
		return db.UpdateByIdAndUserId(db.Notes, noteId, userId, bson.M{"$set": bson.M{"IsTrash": true}})
	}
	return false
}

// 删除别人共享给我的笔记
// 先判断我是否有权限, 笔记是否是我创建的
func (this *TrashService) DeleteSharedNote(noteId, userId, myUserId string) bool {
	note := noteService.GetNote(noteId, userId)
	if shareService.HasUpdatePerm(userId, myUserId, noteId) && note.CreatedUserId.Hex() == myUserId {
		return db.UpdateByIdAndUserId(db.Notes, noteId, userId, bson.M{"$set": bson.M{"IsTrash": true}})
	}
	return false
}

// recover
func (this *TrashService) recoverNote(noteId, notebookId, userId string) bool {
	re := db.UpdateByIdAndUserId(db.Notes, noteId, userId, 
		bson.M{"$set": bson.M{"IsTrash": false, 
			"NotebookId": bson.ObjectIdHex(notebookId)}})
	return re;
}

// 删除trash
func (this *TrashService) DeleteTrash(noteId, userId string) bool {
	// delete note's attachs
	ok := attachService.DeleteAllAttachs(noteId, userId)
	
	// delete note
	ok = db.DeleteByIdAndUserId(db.Notes, noteId, userId)
	// delete content
	ok = db.DeleteByIdAndUserId(db.NoteContents, noteId, userId)
	
	return ok
}

// 列出note, 排序规则, 还有分页
// CreatedTime, UpdatedTime, title 来排序
func (this *TrashService) ListNotes(userId string, 
		pageNumber, pageSize int, sortField string, isAsc bool) (notes []info.Note) {
	_, notes = noteService.ListNotes(userId, "", true, pageNumber, pageSize, sortField, isAsc, false)
	return
}