package info

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 只为blog, 不为note

type BlogItem struct {
	Note
	Abstract string
	Content  string // 可能是content的一部分, 截取. 点击more后就是整个信息了
	HasMore  bool   // 是否是否还有
	User     User   // 用户信息
}

type UserBlogBase struct {
	Logo     string `Logo`
	Title    string `Title`    // 标题
	SubTitle string `SubTitle` // 副标题
	//	AboutMe  string `AboutMe`  // 关于我
}

type UserBlogComment struct {
	CanComment  bool   `CanComment`  // 是否可以评论
	CommentType string `CommentType` // default 或 disqus
	DisqusId    string `DisqusId`
}

type UserBlogStyle struct {
	Style string `Style` // 风格
	Css   string `Css`   // 自定义css
}

// 每个用户一份博客设置信息
type UserBlog struct {
	UserId   bson.ObjectId `bson:"_id"` // 谁的
	Logo     string        `Logo`
	Title    string        `Title`    // 标题
	SubTitle string        `SubTitle` // 副标题
	AboutMe  string        `AboutMe`  // 关于我, 弃用

	CanComment bool `CanComment` // 是否可以评论

	CommentType string `CommentType` // default 或 disqus
	DisqusId    string `DisqusId`

	Style string `Style` // 风格
	Css   string `Css`   // 自定义css

	ThemeId   bson.ObjectId `ThemeId,omitempty`         // 主题Id
	ThemePath string        `bson:"ThemePath" json:"-"` // 不存值, 从Theme中获取, 相对路径 public/

	CateIds []string            `CateIds,omitempty` // 分类Id, 排序好的
	Singles []map[string]string `Singles,omitempty` // 单页, 排序好的, map包含: ["Title"], ["SingleId"]

	PerPageSize int    `PerPageSize,omitempty`
	SortField   string `SortField`       // 排序字段
	IsAsc       bool   `IsAsc,omitempty` // 排序类型, 降序, 升序, 默认是false, 表示降序

	SubDomain string `SubDomain` // 二级域名
	Domain    string `Domain`    // 自定义域名

}

// 博客统计信息
type BlogStat struct {
	NoteId     bson.ObjectId `bson:"_id,omitempty"`
	ReadNum    int           `ReadNum,omitempty`    // 阅读次数 2014/9/28
	LikeNum    int           `LikeNum,omitempty`    // 点赞次数 2014/9/28
	CommentNum int           `CommentNum,omitempty` // 评论次数 2014/9/28
}

// 单页
type BlogSingle struct {
	SingleId    bson.ObjectId `bson:"_id,omitempty"`
	UserId      bson.ObjectId `UserId`
	Title       string        `Title`
	UrlTitle    string        `UrlTitle` // 2014/11/11
	Content     string        `Content`
	UpdatedTime time.Time     `UpdatedTime`
	CreatedTime time.Time     `CreatedTime`
}

//------------------------
// 社交功能, 点赞, 分享, 评论

// 点赞记录
type BlogLike struct {
	LikeId      bson.ObjectId `bson:"_id,omitempty"`
	NoteId      bson.ObjectId `NoteId`
	UserId      bson.ObjectId `UserId`
	CreatedTime time.Time     `CreatedTime`
}

// 评论
type BlogComment struct {
	CommentId bson.ObjectId `bson:"_id,omitempty"`
	NoteId    bson.ObjectId `NoteId`

	UserId  bson.ObjectId `UserId`  // UserId回复ToUserId
	Content string        `Content` // 评论内容

	ToCommentId bson.ObjectId `ToCommendId,omitempty` // 对某条评论进行回复
	ToUserId    bson.ObjectId `ToUserId,omitempty`    // 为空表示直接评论, 不回空表示回复某人

	LikeNum     int      `LikeNum`     // 点赞次数, 评论也可以点赞
	LikeUserIds []string `LikeUserIds` // 点赞的用户ids

	CreatedTime time.Time `CreatedTime`
}

type BlogCommentPublic struct {
	BlogComment
	IsILikeIt bool
}

type BlogUrls struct {
	IndexUrl    string
	CateUrl     string
	SearchUrl   string
	SingleUrl   string
	PostUrl     string
	ArchiveUrl  string
	TagsUrl     string
	TagPostsUrl string
}
