package info

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 举报
type Report struct {
	ReportId bson.ObjectId `bson:"_id"`
	NoteId   bson.ObjectId `NoteId`

	UserId bson.ObjectId `UserId` // UserId回复ToUserId
	Reason string        `Reason` // 评论内容

	CommentId bson.ObjectId `CommendId,omitempty` // 对某条评论进行回复

	CreatedTime time.Time `CreatedTime`
}
