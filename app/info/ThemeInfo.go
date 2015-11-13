package info

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 主题, 每个用户有多个主题, 这里面有主题的配置信息
// 模板, css, js, images, 都在路径Path下
type Theme struct {
	ThemeId   bson.ObjectId          `bson:"_id,omitempty"` // 必须要设置bson:"_id" 不然mgo不会认为是主键
	UserId    bson.ObjectId          `UserId`
	Name      string                 `Name`
	Version   string                 `Version`
	Author    string                 `Author`
	AuthorUrl string                 `AuthorUrl`
	Path      string                 `Path`     // 文件夹路径, public/upload/54d7620d99c37b030600002c/themes/54d867c799c37b533e000001
	Info      map[string]interface{} `Info`     // 所有信息
	IsActive  bool                   `IsActive` // 是否在用

	IsDefault bool   `IsDefault`       // leanote默认主题, 如果用户修改了默认主题, 则先copy之. 也是admin用户的主题
	Style     string `Style,omitempty` // 之前的, 只有default的用户才有blog_default, blog_daqi, blog_left_fixed

	CreatedTime time.Time `CreatedTime`
	UpdatedTime time.Time `UpdatedTime`
}
