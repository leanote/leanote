package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"os"
	"strings"
	"time"
)

type AttachService struct {
}

// add attach
// api调用时, 添加attach之前是没有note的
// fromApi表示是api添加的, updateNote传过来的, 此时不要incNote's usn, 因为updateNote会inc的
func (this *AttachService) AddAttach(attach info.Attach, fromApi bool) (ok bool, msg string) {
	attach.CreatedTime = time.Now()
	ok = db.Insert(db.Attachs, attach)

	note := noteService.GetNoteById(attach.NoteId.Hex())

	// api调用时, 添加attach之前是没有note的
	var userId string
	if note.NoteId != "" {
		userId = note.UserId.Hex()
	} else {
		userId = attach.UploadUserId.Hex()
	}

	if ok {
		// 更新笔记的attachs num
		this.updateNoteAttachNum(attach.NoteId, 1)
	}

	if !fromApi {
		// 增长note's usn
		noteService.IncrNoteUsn(attach.NoteId.Hex(), userId)
	}

	return
}

// 更新笔记的附件个数
// addNum 1或-1
func (this *AttachService) updateNoteAttachNum(noteId bson.ObjectId, addNum int) bool {
	num := db.Count(db.Attachs, bson.M{"NoteId": noteId})
	/*
		note := info.Note{}
		note = noteService.GetNoteById(noteId.Hex())
		note.AttachNum += addNum
		if note.AttachNum < 0 {
			note.AttachNum = 0
		}
		Log(note.AttachNum)
	*/
	return db.UpdateByQField(db.Notes, bson.M{"_id": noteId}, "AttachNum", num)
}

// list attachs
func (this *AttachService) ListAttachs(noteId, userId string) []info.Attach {
	attachs := []info.Attach{}

	// 判断是否有权限为笔记添加附件, userId为空时表示是分享笔记的附件
	if userId != "" && !shareService.HasUpdateNotePerm(noteId, userId) {
		return attachs
	}

	// 笔记是否是自己的
	note := noteService.GetNoteByIdAndUserId(noteId, userId)
	if note.NoteId == "" {
		return attachs
	}

	// TODO 这里, 优化权限控制

	db.ListByQ(db.Attachs, bson.M{"NoteId": bson.ObjectIdHex(noteId)}, &attachs)

	return attachs
}

// api调用, 通过noteIds得到note's attachs, 通过noteId归类返回
func (this *AttachService) getAttachsByNoteIds(noteIds []bson.ObjectId) map[string][]info.Attach {
	attachs := []info.Attach{}
	db.ListByQ(db.Attachs, bson.M{"NoteId": bson.M{"$in": noteIds}}, &attachs)
	noteAttchs := make(map[string][]info.Attach)
	for _, attach := range attachs {
		noteId := attach.NoteId.Hex()
		if itAttachs, ok := noteAttchs[noteId]; ok {
			noteAttchs[noteId] = append(itAttachs, attach)
		} else {
			noteAttchs[noteId] = []info.Attach{attach}
		}
	}
	return noteAttchs
}

func (this *AttachService) UpdateImageTitle(userId, fileId, title string) bool {
	return db.UpdateByIdAndUserIdField(db.Files, fileId, userId, "Title", title)
}

// Delete note to delete attas firstly
func (this *AttachService) DeleteAllAttachs(noteId, userId string) bool {
	note := noteService.GetNoteById(noteId)
	if note.UserId.Hex() == userId {
		attachs := []info.Attach{}
		db.ListByQ(db.Attachs, bson.M{"NoteId": bson.ObjectIdHex(noteId)}, &attachs)
		for _, attach := range attachs {
			attach.Path = strings.TrimLeft(attach.Path, "/")
			os.Remove(revel.BasePath + "/" + attach.Path)
		}
		return true
	}

	return false
}

// delete attach
func (this *AttachService) DeleteAttach(attachId, userId string) (bool, string) {
	attach := info.Attach{}
	db.Get(db.Attachs, attachId, &attach)

	if attach.AttachId != "" {
		// 判断是否有权限为笔记添加附件
		if !shareService.HasUpdateNotePerm(attach.NoteId.Hex(), userId) {
			return false, "No Perm"
		}

		if db.Delete(db.Attachs, bson.M{"_id": bson.ObjectIdHex(attachId)}) {
			this.updateNoteAttachNum(attach.NoteId, -1)
			attach.Path = strings.TrimLeft(attach.Path, "/")
			err := os.Remove(revel.BasePath + "/" + attach.Path)
			if err == nil {
				// userService.UpdateAttachSize(note.UserId.Hex(), -attach.Size)
				// 修改note Usn
				noteService.IncrNoteUsn(attach.NoteId.Hex(), userId)

				return true, "delete file success"
			}
			return false, "delete file error"
		}
		return false, "db error"
	}
	return false, "no such item"
}

// 获取文件路径
// 要判断是否具有权限
// userId是否具有attach的访问权限
func (this *AttachService) GetAttach(attachId, userId string) (attach info.Attach) {
	if attachId == "" {
		return
	}

	attach = info.Attach{}
	db.Get(db.Attachs, attachId, &attach)
	path := attach.Path
	if path == "" {
		return
	}

	note := noteService.GetNoteById(attach.NoteId.Hex())

	// 判断权限

	// 笔记是否是公开的
	if note.IsBlog {
		return
	}

	// 笔记是否是我的
	if note.UserId.Hex() == userId {
		return
	}

	// 我是否有权限查看或协作
	if shareService.HasReadNotePerm(attach.NoteId.Hex(), userId) {
		return
	}

	attach = info.Attach{}
	return
}

// 复制笔记时需要复制附件
// noteService调用, 权限已判断
func (this *AttachService) CopyAttachs(noteId, toNoteId, toUserId string) bool {
	attachs := []info.Attach{}
	db.ListByQ(db.Attachs, bson.M{"NoteId": bson.ObjectIdHex(noteId)}, &attachs)

	// 复制之
	toNoteIdO := bson.ObjectIdHex(toNoteId)
	for _, attach := range attachs {
		attach.AttachId = ""
		attach.NoteId = toNoteIdO

		// 文件复制一份
		_, ext := SplitFilename(attach.Name)
		newFilename := NewGuid() + ext
		dir := "files/" + toUserId + "/attachs"
		filePath := dir + "/" + newFilename
		err := os.MkdirAll(revel.BasePath+"/"+dir, 0755)
		if err != nil {
			return false
		}
		_, err = CopyFile(revel.BasePath+"/"+attach.Path, revel.BasePath+"/"+filePath)
		if err != nil {
			return false
		}
		attach.Name = newFilename
		attach.Path = filePath

		this.AddAttach(attach, false)
	}

	return true
}

// 只留下files的数据, 其它的都删除
func (this *AttachService) UpdateOrDeleteAttachApi(noteId, userId string, files []info.NoteFile) bool {
	// 现在数据库内的
	attachs := this.ListAttachs(noteId, userId)

	nowAttachs := map[string]bool{}
	if files != nil {
		for _, file := range files {
			if file.IsAttach && file.FileId != "" {
				nowAttachs[file.FileId] = true
			}
		}
	}

	for _, attach := range attachs {
		fileId := attach.AttachId.Hex()
		if !nowAttachs[fileId] {
			// 需要删除的
			// TODO 权限验证去掉
			this.DeleteAttach(fileId, userId)
		}
	}

	return false

}
