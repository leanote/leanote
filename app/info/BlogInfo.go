package info

import (
	"labix.org/v2/mgo/bson"
)

// 只为blog, 不为note

type BlogItem struct {
	Note
	Content string // 可能是content的一部分, 截取. 点击more后就是整个信息了
	HasMore bool   // 是否是否还有
}

type UserBlogBase struct {
	Logo       string        `Logo`
	Title      string        `Title`      // 标题
	SubTitle   string        `SubTitle`   // 副标题
	AboutMe    string        `AboutMe` // 关于我
}

type UserBlogComment struct {
	CanComment bool          `CanComment` // 是否可以评论
	DisqusId   string        `DisqusId`
}

type UserBlogStyle struct {
	Style      string        `Style`   // 风格
}

// 每个用户一份博客设置信息
type UserBlog struct {
	UserId     bson.ObjectId `bson:"_id"` // 谁的
	Logo       string        `Logo`
	Title      string        `Title`      // 标题
	SubTitle   string        `SubTitle`   // 副标题
	AboutMe    string        `AboutMe` // 关于我
	
	CanComment bool          `CanComment` // 是否可以评论
	DisqusId   string        `DisqusId`
	
	Style      string        `Style`   // 风格
}