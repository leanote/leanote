package info

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 在数据库中每个
// 修改字段必须要在NotebookService中修改ParseAndSortNotebooks(没有匿名字段), 以后重构
type Notebook struct {
	NotebookId       bson.ObjectId `bson:"_id,omitempty"` // 必须要设置bson:"_id" 不然mgo不会认为是主键
	UserId           bson.ObjectId `bson:"UserId"`
	ParentNotebookId bson.ObjectId `bson:"ParentNotebookId,omitempty"` // 上级
	Seq              int           `Seq`                               // 排序
	Title            string        `Title`                             // 标题
	UrlTitle         string        `UrlTitle`                          // Url标题 2014/11.11加
	NumberNotes      int           `NumberNotes`                       // 笔记数
	IsTrash          bool          `IsTrash,omitempty`                 // 是否是trash, 默认是false
	IsBlog           bool          `IsBlog,omitempty`                  // 是否是Blog 2013/12/29 新加
	CreatedTime      time.Time     `CreatedTime,omitempty`
	UpdatedTime      time.Time     `UpdatedTime,omitempty`

	// 2015/1/15, 更新序号
	Usn       int  `Usn` // UpdateSequenceNum
	IsDeleted bool `IsDeleted`
}

// 仅仅是为了返回前台
type SubNotebooks []*Notebooks // 存地址, 为了生成tree
type Notebooks struct {
	Notebook
	Subs SubNotebooks // 子notebook 在数据库中是没有的
}

// SubNotebook sort
func (this SubNotebooks) Len() int {
	return len(this)
}
func (this SubNotebooks) Less(i, j int) bool {
	return (*this[i]).Seq < (*this[j]).Seq
}
func (this SubNotebooks) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

/*
修改方案, 因为要共享notebook的问题, 所以还是每个notebook一条记录
{
	notebookId,
	title,
	seq,
	parentNoteBookId, // 上级
	userId
}

得到所有该用户的notebook, 然后组装成tree返回之
更新顺序
添加notebook
更新notebook
删除notebook
*/
