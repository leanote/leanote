package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	"gopkg.in/mgo.v2/bson"
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
	note := noteService.GetNote(noteId, userId)
	// 如果是垃圾, 则彻底删除
	if note.IsTrash {
		return this.DeleteTrash(noteId, userId)
	}

	// 首先删除其共享
	if shareService.DeleteShareNoteAll(noteId, userId) {
		// 更新note isTrash = true
		if db.UpdateByIdAndUserId(db.Notes, noteId, userId, bson.M{"$set": bson.M{"IsTrash": true, "Usn": userService.IncrUsn(userId)}}) {
			// recount notebooks' notes number
			notebookIdO := noteService.GetNotebookId(noteId)
			notebookId := notebookIdO.Hex()
			notebookService.ReCountNotebookNumberNotes(notebookId)
			return true
		}
	}

	return false
}

// 删除别人共享给我的笔记
// 先判断我是否有权限, 笔记是否是我创建的
func (this *TrashService) DeleteSharedNote(noteId, myUserId string) bool {
	note := noteService.GetNoteById(noteId)
	userId := note.UserId.Hex()
	if shareService.HasUpdatePerm(userId, myUserId, noteId) && note.CreatedUserId.Hex() == myUserId {
		return db.UpdateByIdAndUserId(db.Notes, noteId, userId, bson.M{"$set": bson.M{"IsTrash": true, "Usn": userService.IncrUsn(userId)}})
	}
	return false
}

// recover
func (this *TrashService) recoverNote(noteId, notebookId, userId string) bool {
	re := db.UpdateByIdAndUserId(db.Notes, noteId, userId,
		bson.M{"$set": bson.M{"IsTrash": false,
			"Usn":        userService.IncrUsn(userId),
			"NotebookId": bson.ObjectIdHex(notebookId)}})
	return re
}

// 删除trash
func (this *TrashService) DeleteTrash(noteId, userId string) bool {
	note := noteService.GetNote(noteId, userId)
	if note.NoteId == "" {
		return false
	}
	// delete note's attachs
	ok := attachService.DeleteAllAttachs(noteId, userId)

	// 设置删除位
	ok = db.UpdateByIdAndUserIdMap(db.Notes, noteId, userId,
		bson.M{"IsDeleted": true,
			"Usn": userService.IncrUsn(userId)})
	// delete note
	//	ok = db.DeleteByIdAndUserId(db.Notes, noteId, userId)

	// delete content
	ok = db.DeleteByIdAndUserId(db.NoteContents, noteId, userId)

	// 删除content history
	ok = db.DeleteByIdAndUserId(db.NoteContentHistories, noteId, userId)

	// 重新统计tag's count
	// TODO 这里会改变tag's Usn
	tagService.reCountTagCount(userId, note.Tags)

	return ok
}

func (this *TrashService) DeleteTrashApi(noteId, userId string, usn int) (bool, string, int) {
	note := noteService.GetNote(noteId, userId)

	if note.NoteId == "" || note.IsDeleted {
		return false, "notExists", 0
	}

	if note.Usn != usn {
		return false, "conflict", 0
	}

	// delete note's attachs
	ok := attachService.DeleteAllAttachs(noteId, userId)

	// 设置删除位
	afterUsn := userService.IncrUsn(userId)
	ok = db.UpdateByIdAndUserIdMap(db.Notes, noteId, userId,
		bson.M{"IsDeleted": true,
			"Usn": afterUsn})

	// delete content
	ok = db.DeleteByIdAndUserId(db.NoteContents, noteId, userId)

	// 删除content history
	ok = db.DeleteByIdAndUserId(db.NoteContentHistories, noteId, userId)

	// 一个BUG, iOS删除直接调用这个API, 结果没有重新recount
	// recount notebooks' notes number
	notebookService.ReCountNotebookNumberNotes(note.NotebookId.Hex())

	return ok, "", afterUsn
}

// 列出note, 排序规则, 还有分页
// CreatedTime, UpdatedTime, title 来排序
func (this *TrashService) ListNotes(userId string,
	pageNumber, pageSize int, sortField string, isAsc bool) (notes []info.Note) {
	_, notes = noteService.ListNotes(userId, "", true, pageNumber, pageSize, sortField, isAsc, false)
	return
}
