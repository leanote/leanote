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
// imgSrc 防止博客修改了, 但内容删除了
func (this *NoteImageService) UpdateNoteImages(userId, noteId, imgSrc, content string) bool {
	// 让主图成为内容的一员
	if imgSrc != "" {
		content = "<img src=\"" + imgSrc + "\" >" + content
	}
	reg, _ := regexp.Compile("outputImage\\?fileId=([a-z0-9A-Z]{24})")
	find := reg.FindAllStringSubmatch(content, -1) // 查找所有的
	
	// 删除旧的
	db.DeleteAll(db.NoteImages, bson.M{"NoteId": bson.ObjectIdHex(noteId)})
	
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

// 复制图片, 把note的图片都copy给我, 且修改noteContent图片路径
func (this *NoteImageService) CopyNoteImages(fromNoteId, fromUserId, newNoteId, content, toUserId string) string {
	// 得到fromNoteId的noteImages, 如果为空, 则直接返回content
	noteImages := []info.NoteImage{}
	db.ListByQWithFields(db.NoteImages, bson.M{"NoteId": bson.ObjectIdHex(fromNoteId)}, []string{"ImageId"}, &noteImages)
	
	if len(noteImages) == 0 {
		return content;
	}
	
	// <img src="/file/outputImage?fileId=12323232" />
	// 把fileId=1232替换成新的
	replaceMap := map[string]string{}
	for _, noteImage := range noteImages {
		imageId := noteImage.ImageId.Hex()
		ok, newImageId := fileService.CopyImage(fromUserId, imageId, toUserId)
		if ok {
			replaceMap[imageId] = newImageId
		}
	}
	
	if len(replaceMap) > 0 {
		// 替换之
		reg, _ := regexp.Compile("outputImage\\?fileId=([a-z0-9A-Z]{24})")
		content = reg.ReplaceAllStringFunc(content, func(each string) string {
			// each=outputImage?fileId=541bd2f599c37b4f3r000003
			fileId := each[len(each)-24:] // 得到后24位, 也即id
			if replaceFileId, ok := replaceMap[fileId]; ok {
				return "outputImage?fileId=" + replaceFileId
			}
			return each
		});
	}
	
	return content;
}
