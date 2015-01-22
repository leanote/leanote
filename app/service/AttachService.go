package service

import (
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"github.com/leanote/leanote/app/info"
	"github.com/leanote/leanote/app/db"
	"gopkg.in/mgo.v2/bson"
	"time"
	"os"
	"strings"
)

type AttachService struct {
}

// add attach
func (this *AttachService) AddAttach(attach info.Attach) (ok bool, msg string) {
	attach.CreatedTime = time.Now()
	ok = db.Insert(db.Attachs, attach)
	
	if ok {
		// 更新笔记的attachs num
		this.updateNoteAttachNum(attach.NoteId, 1)
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
	
	db.ListByQ(db.Attachs, bson.M{"NoteId": bson.ObjectIdHex(noteId)}, &attachs)
	
	return attachs
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
	
	if(attach.AttachId != "") {
		// 判断是否有权限为笔记添加附件
		if !shareService.HasUpdateNotePerm(attach.NoteId.Hex(), userId) {
			return false, "No Perm"
		}
		
		if db.Delete(db.Attachs, bson.M{"_id": bson.ObjectIdHex(attachId)}) {
			this.updateNoteAttachNum(attach.NoteId, -1)
			attach.Path = strings.TrimLeft(attach.Path, "/")
			err := os.Remove(revel.BasePath + "/" + attach.Path)
			if err == nil {
				return true, "delete file error"
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
func (this *AttachService) GetAttach(attachId, userId, token, sessionId string) (attach info.Attach) {
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
	
	//userId为空则是分享笔记的附件
	if userId == "" && sessionId != "" {
		if token != "" {
			realToken := sessionService.GetToken(sessionId)
			if token == realToken {
				Log("attach token is equal!")
				return 
			}
		} 
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
		err := os.MkdirAll(revel.BasePath + "/" + dir, 0755)
		if err != nil {
			return false
		}
		_, err = CopyFile(revel.BasePath + "/" + attach.Path, revel.BasePath + "/" + filePath)
		if err != nil {
			return false
		}
		attach.Name = newFilename
		attach.Path = filePath
		
		this.AddAttach(attach)
	}
	
	return true
}