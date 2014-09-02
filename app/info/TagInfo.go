package info

import (
	"gopkg.in/mgo.v2/bson"
)

// 这里主要是为了统计每个tag的note数目
type TagNote struct {
	TagId   bson.ObjectId `bson:"_id,omitempty"` // 必须要设置bson:"_id" 不然mgo不会认为是主键
	UserId  bson.ObjectId `bson:"UserId"`
	Tag   string        `Title`   // 标题
	NoteNum int           `NoteNum` // note数目
}

// 每个用户一条记录, 存储用户的所有tags
type Tag struct {
	UserId  bson.ObjectId `bson:"_id"`
	Tags []string `Tags`
}