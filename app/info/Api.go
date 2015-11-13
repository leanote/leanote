package info

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//---------
// 数据结构
//---------
type NoteFile struct {
	FileId      string // 服务器端Id
	LocalFileId string // 客户端Id
	Type        string // images/png, doc, xls, 根据fileName确定
	Title       string
	HasBody     bool // 传过来的值是否要更新内容
	IsAttach    bool // 是否是附件, 不是附件就是图片
}
type ApiNote struct {
	NoteId     string
	NotebookId string
	UserId     string
	Title      string
	Desc       string
	//	ImgSrc     string
	Tags       []string
	Abstract   string
	Content    string
	IsMarkdown bool
	//	FromUserId string // 为共享而新建
	IsBlog      bool // 是否是blog, 更新note不需要修改, 添加note时才有可能用到, 此时需要判断notebook是否设为Blog
	IsTrash     bool
	IsDeleted   bool
	Usn         int
	Files       []NoteFile
	CreatedTime time.Time
	UpdatedTime time.Time
	PublicTime  time.Time
}

// 内容
type ApiNoteContent struct {
	NoteId bson.ObjectId `bson:"_id,omitempty"`
	UserId bson.ObjectId `bson:"UserId"`

	Content string `Content`

	//	CreatedTime   time.Time     `CreatedTime`
	//	UpdatedTime   time.Time     `UpdatedTime`
}

// 转换
func NoteToApiNote(note Note, files []NoteFile) ApiNote {
	apiNote := ApiNote{}
	return apiNote
}

//----------
// 用户信息
//----------

type ApiUser struct {
	UserId   string
	Username string
	Email    string
	Verified bool
	Logo     string
}

//----------
// Notebook
//----------
type ApiNotebook struct {
	NotebookId       bson.ObjectId `bson:"_id,omitempty"` // 必须要设置bson:"_id" 不然mgo不会认为是主键
	UserId           bson.ObjectId `bson:"UserId"`
	ParentNotebookId bson.ObjectId `bson:"ParentNotebookId,omitempty"` // 上级
	Seq              int           `Seq`                               // 排序
	Title            string        `Title`                             // 标题
	UrlTitle         string        `UrlTitle`                          // Url标题 2014/11.11加
	IsBlog           bool          `IsBlog,omitempty`                  // 是否是Blog 2013/12/29 新加
	CreatedTime      time.Time     `CreatedTime,omitempty`
	UpdatedTime      time.Time     `UpdatedTime,omitempty`
	Usn              int           `Usn` // UpdateSequenceNum
	IsDeleted        bool          `IsDeleted`
}

//---------
// api 返回
//---------

// 一般返回
type ApiRe struct {
	Ok  bool
	Msg string
}

func NewApiRe() ApiRe {
	return ApiRe{Ok: false}
}

// auth
type AuthOk struct {
	Ok       bool
	Token    string
	UserId   bson.ObjectId
	Email    string
	Username string
}

// 供notebook, note, tag更新的返回数据用
type ReUpdate struct {
	Ok  bool
	Msg string
	Usn int
}

func NewReUpdate() ReUpdate {
	return ReUpdate{Ok: false}
}
