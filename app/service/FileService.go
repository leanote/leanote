package service

import (
//	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"github.com/leanote/leanote/app/info"
	"github.com/leanote/leanote/app/db"
	"gopkg.in/mgo.v2/bson"
	"time"
	"os"
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
	
	q := bson.M{"UserId": bson.ObjectIdHex(userId)}
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
			err := os.Remove(revel.BasePath + "/public/" + file.Path)
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