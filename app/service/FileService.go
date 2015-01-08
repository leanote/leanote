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

const DEFAULT_ALBUM_ID = "52d3e8ac99c37b7f0d000001"

type FileService struct {
}

// add Image
func (this *FileService) AddImage(image info.File, albumId, userId string) bool {
	image.CreatedTime = time.Now()
	if albumId != "" {
		image.AlbumId = bson.ObjectIdHex(albumId)
	} else {
		image.AlbumId = bson.ObjectIdHex(DEFAULT_ALBUM_ID)
		image.IsDefaultAlbum = true
	}
	image.UserId = bson.ObjectIdHex(userId)
	
	return db.Insert(db.Files, image)
}

// list images
// if albumId == "" get default album images
func (this *FileService) ListImagesWithPage(userId, albumId, key string, pageNumber, pageSize int) info.Page {
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, "CreatedTime", false)
	files := []info.File{}
	
	q := bson.M{"UserId": bson.ObjectIdHex(userId), "Type": ""} // life
	if albumId != "" {
		q["AlbumId"] = bson.ObjectIdHex(albumId);
	} else {
		q["IsDefaultAlbum"] = true
	}
	if key != "" {
		q["Title"] =  bson.M{"$regex": bson.RegEx{".*?" + key + ".*", "i"}}
	}
	
//	LogJ(q)
	
	count := db.Count(db.Files, q);
	
	db.Files.
		Find(q).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&files)
		
	return info.Page{Count: count, List: files}
}

func (this *FileService) UpdateImageTitle(userId, fileId, title string) bool {
	return db.UpdateByIdAndUserIdField(db.Files, fileId, userId, "Title", title)
}

// get all images names
// for upgrade
func (this *FileService) GetAllImageNamesMap(userId string) (m map[string]bool) {
	q := bson.M{"UserId": bson.ObjectIdHex(userId)}
	files := []info.File{}
	db.ListByQWithFields(db.Files, q, []string{"Name"}, &files)
	
	m = make(map[string]bool)
	if len(files) == 0 {
		return
	}
	
	for _, file := range files {
		m[file.Name] = true
	}
	return
}

// delete image
func (this *FileService) DeleteImage(userId, fileId string) (bool, string) {
	file := info.File{}
	db.GetByIdAndUserId(db.Files, fileId, userId, &file)
	
	if(file.FileId != "") {
		if db.DeleteByIdAndUserId(db.Files, fileId, userId) {
			// delete image
			// TODO
			file.Path = strings.TrimLeft(file.Path, "/")
			var err error
			if strings.HasPrefix(file.Path, "upload") {
				Log(file.Path)
				err = os.Remove(revel.BasePath + "/public/" + file.Path)
			} else {
				err = os.Remove(revel.BasePath + "/" + file.Path)
			}
			if err == nil {
				return true, ""
			}
			return false, "delete file error!"
		}
		return false, "db error"
	}
	return false, "no such item"
}

// update image title
func (this *FileService) UpdateImage(userId, fileId, title string) bool {
	return db.UpdateByIdAndUserIdField(db.Files, fileId, userId, "Title", title)
}

// 获取文件路径
// 要判断是否具有权限
// userId是否具有fileId的访问权限
func (this *FileService) GetFile(userId, fileId string) string {
	if fileId == "" {
		return ""
	}
	
	file := info.File{}
	db.Get(db.Files, fileId, &file)
	path := file.Path
	if path == "" {
		return ""
	}
	
	// 1. 判断权限
	
	// 是否是我的文件
	if userId != "" && file.UserId.Hex() == userId {
		return path
	}
	
	// 得到使用过该fileId的所有笔记NoteId
	// 这些笔记是否有public的, 若有则ok
	// 这些笔记(笔记本)是否有共享给我的, 若有则ok
	
	noteIds := noteImageService.GetNoteIds(fileId)
	if noteIds != nil && len(noteIds) > 0 {
		// 这些笔记是否有public的
		if db.Has(db.Notes, bson.M{"_id": bson.M{"$in": noteIds}, "IsBlog": true}) {
			return path
		}
		
		// 2014/12/28 修复, 如果是分享给用户组, 那就不行, 这里可以实现
		for _, noteId := range noteIds {
			note := noteService.GetNoteById(noteId.Hex())
			if shareService.HasReadPerm(note.UserId.Hex(), userId, noteId.Hex()) {
				return path;
			}
		}
		/*
		// 若有共享给我的笔记?
		// 对该笔记可读?
		if db.Has(db.ShareNotes, bson.M{"ToUserId": bson.ObjectIdHex(userId), "NoteId": bson.M{"$in": noteIds}}) {
			return path
		}
		
		// 笔记本是否共享给我?
		// 通过笔记得到笔记本
		notes := []info.Note{}
		db.ListByQWithFields(db.Notes, bson.M{"_id": bson.M{"$in": noteIds}}, []string{"NotebookId"}, &notes)	
		if notes != nil && len(notes) > 0 {
			notebookIds := make([]bson.ObjectId, len(notes))
			for i := 0; i < len(notes); i++ {
				notebookIds[i] = notes[i].NotebookId
			}
			
			if db.Has(db.ShareNotebooks, bson.M{"ToUserId": bson.ObjectIdHex(userId), "NotebookId": bson.M{"$in": notebookIds}}) {
				return path
			}
		}
		*/
	}
	
	// 可能是刚复制到owner上, 但内容又没有保存, 所以没有note->imageId的映射, 此时看是否有fromFileId
	if file.FromFileId != "" {
		fromFile := info.File{}
		db.Get2(db.Files, file.FromFileId, &fromFile)
		if fromFile.UserId.Hex() == userId {
			return fromFile.Path
		}
	}
	
	return ""
}

// 复制图片
func (this *FileService) CopyImage(userId, fileId, toUserId string) (bool, string) {
	// 是否已经复制过了
	file2 := info.File{}
	db.GetByQ(db.Files, bson.M{"UserId": bson.ObjectIdHex(toUserId), "FromFileId": bson.ObjectIdHex(fileId)}, &file2)
	if file2.FileId != "" {
		return true, file2.FileId.Hex();
	}

	// 复制之
	
	file := info.File{}
	db.GetByIdAndUserId(db.Files, fileId, userId, &file)
	
	if file.FileId == "" || file.UserId.Hex() != userId {
		return false, ""
	}
		
	_, ext := SplitFilename(file.Name)
	newFilename := NewGuid() + ext
	
	dir := "files/" + toUserId + "/images"
	filePath := dir + "/" + newFilename
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return false, ""
	}
	
	_, err = CopyFile(revel.BasePath + "/" + file.Path, revel.BasePath + "/" + filePath)
	if err != nil {
		Log(err)
		return false, ""
	}
	
	fileInfo := info.File{Name: newFilename,
		Title: file.Title,
		Path: filePath,
		Size: file.Size, 
		FromFileId: file.FileId}
	id := bson.NewObjectId();
	fileInfo.FileId = id
	fileId = id.Hex()
	Ok := this.AddImage(fileInfo, "", toUserId)
	
	if Ok {
		return Ok, id.Hex()
	}
	return false, ""
}

// 是否是我的文件
func (this *FileService) IsMyFile(userId, fileId string) bool {
	return db.Has(db.Files, bson.M{"UserId": bson.ObjectIdHex(userId), "_id": bson.ObjectIdHex(fileId)})
}
