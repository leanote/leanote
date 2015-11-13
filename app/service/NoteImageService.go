package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	//	"time"
)

type NoteImageService struct {
}

// 通过id, userId得到noteIds
func (this *NoteImageService) GetNoteIds(imageId string) []bson.ObjectId {
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

// TODO 这个web可以用, 但api会传来, 不用用了
// 解析内容中的图片, 建立图片与note的关系
// <img src="/file/outputImage?fileId=12323232" />
// 图片必须是我的, 不然不添加
// imgSrc 防止博客修改了, 但内容删除了
func (this *NoteImageService) UpdateNoteImages(userId, noteId, imgSrc, content string) bool {
	// 让主图成为内容的一员
	if imgSrc != "" {
		content = "<img src=\"" + imgSrc + "\" >" + content
	}
	// life 添加getImage
	reg, _ := regexp.Compile("(outputImage|getImage)\\?fileId=([a-z0-9A-Z]{24})")
	find := reg.FindAllStringSubmatch(content, -1) // 查找所有的

	// 删除旧的
	db.DeleteAll(db.NoteImages, bson.M{"NoteId": bson.ObjectIdHex(noteId)})

	// 添加新的
	var fileId string
	noteImage := info.NoteImage{NoteId: bson.ObjectIdHex(noteId)}
	hasAdded := make(map[string]bool)
	if find != nil && len(find) > 0 {
		for _, each := range find {
			if each != nil && len(each) == 3 {
				fileId = each[2] // 现在有两个子表达式了
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
	/* 弃用之
	// 得到fromNoteId的noteImages, 如果为空, 则直接返回content
	noteImages := []info.NoteImage{}
	db.ListByQWithFields(db.NoteImages, bson.M{"NoteId": bson.ObjectIdHex(fromNoteId)}, []string{"ImageId"}, &noteImages)
	if len(noteImages) == 0 {
		return content;
	}
	for _, noteImage := range noteImages {
		imageId := noteImage.ImageId.Hex()
		ok, newImageId := fileService.CopyImage(fromUserId, imageId, toUserId)
		if ok {
			replaceMap[imageId] = newImageId
		}
	}
	*/

	// 因为很多图片上传就会删除, 所以直接从内容中查看图片id进行复制

	// <img src="/file/outputImage?fileId=12323232" />
	// 把fileId=1232替换成新的
	replaceMap := map[string]string{}

	reg, _ := regexp.Compile("(outputImage|getImage)\\?fileId=([a-z0-9A-Z]{24})")
	content = reg.ReplaceAllStringFunc(content, func(each string) string {
		// each = outputImage?fileId=541bd2f599c37b4f3r000003
		// each = getImage?fileId=541bd2f599c37b4f3r000003

		fileId := each[len(each)-24:] // 得到后24位, 也即id

		if _, ok := replaceMap[fileId]; !ok {
			if bson.IsObjectIdHex(fileId) {
				ok2, newImageId := fileService.CopyImage(fromUserId, fileId, toUserId)
				if ok2 {
					replaceMap[fileId] = newImageId
				} else {
					replaceMap[fileId] = ""
				}
			} else {
				replaceMap[fileId] = ""
			}
		}

		replaceFileId := replaceMap[fileId]
		if replaceFileId != "" {
			if each[0] == 'o' {
				return "outputImage?fileId=" + replaceFileId
			}
			return "getImage?fileId=" + replaceFileId
		}
		return each
	})

	return content
}

//
func (this *NoteImageService) getImagesByNoteIds(noteIds []bson.ObjectId) map[string][]info.File {
	noteNoteImages := []info.NoteImage{}
	db.ListByQ(db.NoteImages, bson.M{"NoteId": bson.M{"$in": noteIds}}, &noteNoteImages)

	// 得到imageId, 再去files表查所有的Files
	imageIds := []bson.ObjectId{}

	// 图片1 => N notes
	imageIdNotes := map[string][]string{} // imageId => [noteId1, noteId2, ...]
	for _, noteImage := range noteNoteImages {
		imageId := noteImage.ImageId
		imageIds = append(imageIds, imageId)

		imageIdHex := imageId.Hex()
		noteId := noteImage.NoteId.Hex()
		if notes, ok := imageIdNotes[imageIdHex]; ok {
			imageIdNotes[imageIdHex] = append(notes, noteId)
		} else {
			imageIdNotes[imageIdHex] = []string{noteId}
		}
	}

	// 得到所有files
	files := []info.File{}
	db.ListByQ(db.Files, bson.M{"_id": bson.M{"$in": imageIds}}, &files)

	// 建立note->file关联
	noteImages := make(map[string][]info.File)
	for _, file := range files {
		fileIdHex := file.FileId.Hex() // == imageId
		// 这个fileIdHex有哪些notes呢?
		if notes, ok := imageIdNotes[fileIdHex]; ok {
			for _, noteId := range notes {
				if files, ok2 := noteImages[noteId]; ok2 {
					noteImages[noteId] = append(files, file)
				} else {
					noteImages[noteId] = []info.File{file}
				}
			}
		}
	}
	return noteImages
}
