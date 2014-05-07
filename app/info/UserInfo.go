package info

import (
	"labix.org/v2/mgo/bson"
	"time"
)

// 第三方类型
const (
	ThirdGithub = iota
	ThirdQQ
)

type User struct {
	UserId      bson.ObjectId `bson:"_id,omitempty"` // 必须要设置bson:"_id" 不然mgo不会认为是主键
	Email       string        `Email`                // 全是小写
	Verified    bool          `Verified`             // Email是否已验证过?
	Username    string        `Username`             // 不区分大小写, 全是小写
	UsernameRaw string        `UsernameRaw`          // 可能有大小写
	Pwd         string        `bson:"Pwd" json:"-"`
	CreatedTime time.Time     `CreatedTime`

	// 主题
	Theme string `Theme`

	// 用户配置
	NotebookWidth int  `NotebookWidth` // 笔记本宽度
	NoteListWidth int  `NoteListWidth` // 笔记列表宽度
	LeftIsMin     bool `LeftIsMin`     // 左侧是否是隐藏的, 默认是打开的

	// 这里 第三方登录
	ThirdUserId   string `ThirdUserId`   // 用户Id, 在第三方中唯一可识别
	ThirdUsername string `ThirdUsername` // 第三方中username, 为了显示
	ThirdType     int    `ThirdType`     // 第三方类型
}
