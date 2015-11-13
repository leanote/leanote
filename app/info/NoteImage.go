package info

import (
	"gopkg.in/mgo.v2/bson"
)

// 笔记内部图片
type NoteImage struct {
	NoteImageId bson.ObjectId `bson:"_id,omitempty"` // 必须要设置bson:"_id" 不然mgo不会认为是主键
	NoteId      bson.ObjectId `bson:"NoteId"`        // 笔记
	ImageId     bson.ObjectId `bson:"ImageId"`       // 图片fileId
}
