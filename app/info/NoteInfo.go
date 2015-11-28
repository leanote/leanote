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

	ImgSrc string   `ImgSrc` // 图片, 第一张缩略图地址
	Tags   []string `Tags,omitempty`

	IsTrash bool `IsTrash` // 是否是trash, 默认是false

	IsBlog         bool   `IsBlog,omitempty`      // 是否设置成了blog 2013/12/29 新加
	UrlTitle       string `UrlTitle,omitempty`    // 博客的url标题, 为了更友好的url, 在UserId, UrlName下唯一
	IsRecommend    bool   `IsRecommend,omitempty` // 是否为推荐博客 2014/9/24新加
	IsTop          bool   `IsTop,omitempty`       // blog是否置顶
	HasSelfDefined bool   `HasSelfDefined`        // 是否已经自定义博客图片, desc, abstract

	// 2014/9/28 添加评论社交功能
	ReadNum    int `ReadNum,omitempty`    // 阅读次数 2014/9/28
	LikeNum    int `LikeNum,omitempty`    // 点赞次数 2014/9/28
	CommentNum int `CommentNum,omitempty` // 评论次数 2014/9/28

	IsMarkdown bool `IsMarkdown` // 是否是markdown笔记, 默认是false

	AttachNum int `AttachNum` // 2014/9/21, attachments num

	CreatedTime   time.Time     `CreatedTime`
	UpdatedTime   time.Time     `UpdatedTime`
	RecommendTime time.Time     `RecommendTime,omitempty` // 推荐时间
	PublicTime    time.Time     `PublicTime,omitempty`    // 发表时间, 公开为博客则设置
	UpdatedUserId bson.ObjectId `bson:"UpdatedUserId"`    // 如果共享了, 并可写, 那么可能是其它他修改了

	// 2015/1/15, 更新序号
	Usn int `Usn` // UpdateSequenceNum

	IsDeleted bool `IsDeleted` // 删除位
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

// 为了NoteController接收参数

// 更新note或content
// 肯定会传userId(谁的), NoteId
// 会传Title, Content, Tags, 一种或几种
type NoteOrContent struct {
	NotebookId string
	NoteId     string
	UserId     string
	Title      string
	Desc       string
	ImgSrc     string
	Tags       string
	Content    string
	Abstract   string
	IsNew      bool
	IsMarkdown bool
	FromUserId string // 为共享而新建
	IsBlog     bool   // 是否是blog, 更新note不需要修改, 添加note时才有可能用到, 此时需要判断notebook是否设为Blog
}
