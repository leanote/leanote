package service

import (
	"github.com/leanote/leanote/app/info"
	"github.com/leanote/leanote/app/db"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"regexp"
//	"time"
)

type NoteImageService struct {
}

// 通过id, userId得到noteIds
func (this *NoteImageService) GetNoteIds(imageId string) ([]bson.ObjectId) {
	noteImages := []info.NoteImage{}
	db.ListByQWithFields(db.NoteImages, bson.M{"ImageId": bson.ObjectIdHex(imageId)}, []string{"NoteId"}, &noteImages)	
	
	if noteImages != nil && len(noteImages) > 0 {
		noteIds := make([]bson.ObjectId, len(noteImages))
		cnt := len(noteImages)
		for i := 0; i < cnt; i++ {
			noteIds[i] = noteImages[i].NoteId
		}
		return noteIds
	}
	
	return nil
}

// 解析内容中的图片, 建立图片与note的关系
// <img src="/file/outputImage?fileId=12323232" />
// 图片必须是我的, 不然不添加
func (this *NoteImageService) UpdateNoteImages(userId, noteId, content string) bool {
	reg, _ := regexp.Compile("outputImage\\?fileId=([a-z0-9A-Z]{24})")
	find := reg.FindAllStringSubmatch(content, -1) // 查找所有的
	
	// 删除旧的
	db.DeleteAll(db.NoteImages, bson.M{"NoteId": bson.ObjectIdHex(noteId)})
	
	Log("--------ii--")
	
	// 添加新的
	var fileId string
	noteImage := info.NoteImage{NoteId: bson.ObjectIdHex(noteId)}
	hasAdded := make(map[string]bool)
	if find != nil && len(find) > 0 {
		for _, each := range find {
			if each != nil && len(each) == 2 {
				fileId = each[1]
				// 之前没能添加过的
				if _, ok := hasAdded[fileId]; !ok {
					Log(fileId)
					// 判断是否是我的文件
					if fileService.IsMyFile(userId, fileId) {
						Log("?????")
						Log("<><><>")
						noteImage.ImageId = bson.ObjectIdHex(fileId)
						db.Insert(db.NoteImages, noteImage)
					}
					hasAdded[fileId] = true
				}
			}
		}
	}
	
	return true
}