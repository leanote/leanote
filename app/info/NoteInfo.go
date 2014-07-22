package info

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 只存笔记基本信息
// 内容不存放
type Note struct {
	NoteId        bson.ObjectId `bson:"_id,omitempty"`           // 必须要设置bson:"_id" 不然mgo不会认为是主键
	UserId        bson.ObjectId `bson:"UserId"`                  // 谁的
	CreatedUserId bson.ObjectId `bson:"CreatedUserId,omitempty"` // 谁创建的(UserId != CreatedUserId, 是因为共享). 只是共享才有, 默认为空, 不存 必须要加omitempty
	NotebookId    bson.ObjectId `bson:"NotebookId"`
	Title         string        `Title` // 标题
	Desc          string        `Desc`  // 描述, 非html

	ImgSrc  string   `ImgSrc` // 图片, 第一张缩略图地址
	Tags    []string `Tags,omitempty`
	IsTrash bool     `IsTrash` // 是否是trash, 默认是false

	IsBlog bool `IsBlog,omitempty` // 是否设置成了blog 2013/12/29 新加
	IsTop  bool `IsTop,omitempty`  // blog是否置顶

	IsMarkdown bool `IsMarkdown` // 是否是markdown笔记, 默认是false

	CreatedTime   time.Time     `CreatedTime`
	UpdatedTime   time.Time     `UpdatedTime`
	UpdatedUserId bson.ObjectId `bson:"UpdatedUserId"` // 如果共享了, 并可写, 那么可能是其它他修改了
}

// 内容
type NoteContent struct {
	NoteId bson.ObjectId `bson:"_id,omitempty"`
	UserId bson.ObjectId `bson:"UserId"`

	IsBlog bool `IsBlog,omitempty` // 为了搜索博客

	Content  string `Content`
	Abstract string `Abstract` // 摘要, 有html标签, 比content短, 在博客展示需要, 不放在notes表中

	CreatedTime   time.Time     `CreatedTime`
	UpdatedTime   time.Time     `UpdatedTime`
	UpdatedUserId bson.ObjectId `bson:"UpdatedUserId"` // 如果共享了, 并可写, 那么可能是其它他修改了
}

// 基本信息和内容在一起
type NoteAndContent struct {
	Note
	NoteContent
}

// 历史记录
// 每一个历史记录对象
type EachHistory struct {
	UpdatedUserId bson.ObjectId `UpdatedUserId`
	UpdatedTime   time.Time     `UpdatedTime`
	Content       string        `Content`
}
type NoteContentHistory struct {
	NoteId    bson.ObjectId `bson:"_id,omitempty"`
	UserId    bson.ObjectId `bson:"UserId"` // 所属者
	Histories []EachHistory `Histories`
}
