package info

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 只为blog, 不为note

type BlogItem struct {
	Note
	Content string // 可能是content的一部分, 截取. 点击more后就是整个信息了
	HasMore bool   // 是否是否还有
	User    User   // 用户信息
}

type UserBlogBase struct {
	Logo     string `Logo`
	Title    string `Title`    // 标题
	SubTitle string `SubTitle` // 副标题
	AboutMe  string `AboutMe`  // 关于我
}

type UserBlogComment struct {
	CanComment bool   `CanComment` // 是否可以评论
	CommentType string `CommentType` // default 或 disqus
	DisqusId   string `DisqusId`
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
	AboutMe  string        `AboutMe`  // 关于我

	CanComment bool   `CanComment` // 是否可以评论
	
	CommentType string `CommentType` // default 或 disqus
	DisqusId   string `DisqusId`

	Style string `Style` // 风格
	Css   string `Css`   // 自定义css

	SubDomain string `SubDomain` // 二级域名
	Domain    string `Domain`    // 自定义域名
}

//------------------------
// 社交功能, 点赞, 分享, 评论

// 点赞记录
type BlogLike struct {
	LikeId      bson.ObjectId `bson:"_id"`
	NoteId      bson.ObjectId `NoteId`
	UserId      bson.ObjectId `UserId`
	CreatedTime time.Time     `CreatedTime`
}

// 评论
type BlogComment struct {
	CommentId bson.ObjectId `bson:"_id"`
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
