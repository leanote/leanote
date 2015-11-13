package info

import (
	"gopkg.in/mgo.v2/bson"
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

	Logo string `Logo` // 9-24
	// 主题
	Theme string `Theme`

	// 用户配置
	NotebookWidth int  `NotebookWidth` // 笔记本宽度
	NoteListWidth int  `NoteListWidth` // 笔记列表宽度
	MdEditorWidth int  `MdEditorWidth` // markdown 左侧编辑器宽度
	LeftIsMin     bool `LeftIsMin`     // 左侧是否是隐藏的, 默认是打开的

	// 这里 第三方登录
	ThirdUserId   string `ThirdUserId`   // 用户Id, 在第三方中唯一可识别
	ThirdUsername string `ThirdUsername` // 第三方中username, 为了显示
	ThirdType     int    `ThirdType`     // 第三方类型

	// 用户的帐户类型

	ImageNum   int           `bson:"ImageNum" json:"-"`   // 图片数量
	ImageSize  int           `bson:"ImageSize" json:"-"`  // 图片大小
	AttachNum  int           `bson:"AttachNum" json:"-"`  // 附件数量
	AttachSize int           `bson:"AttachSize" json:"-"` // 附件大小
	FromUserId bson.ObjectId `FromUserId,omitempty`       // 邀请的用户

	AccountType      string    `bson:"AccountType" json:"-"`      // normal(为空), premium
	AccountStartTime time.Time `bson:"AccountStartTime" json:"-"` // 开始日期
	AccountEndTime   time.Time `bson:"AccountEndTime" json:"-"`   // 结束日期
	// 阈值
	MaxImageNum      int `bson:"MaxImageNums" json:"-"`     // 图片数量
	MaxImageSize     int `bson:"MaxImageSize" json:"-"`     // 图片大小
	MaxAttachNum     int `bson:"MaxAttachNum" json:"-"`     // 图片数量
	MaxAttachSize    int `bson:"MaxAttachSize" json:"-"`    // 图片大小
	MaxPerAttachSize int `bson:"MaxPerAttachSize" json:"-"` // 单个附件大小

	// 2015/1/15, 更新序号
	Usn            int       `Usn`                   // UpdateSequenceNum , 全局的
	FullSyncBefore time.Time `bson:"FullSyncBefore"` // 需要全量同步的时间, 如果 > 客户端的LastSyncTime, 则需要全量更新
}

type UserAccount struct {
	AccountType      string    `bson:"AccountType" json:"-"`      // normal(为空), premium
	AccountStartTime time.Time `bson:"AccountStartTime" json:"-"` // 开始日期
	AccountEndTime   time.Time `bson:"AccountEndTime" json:"-"`   // 结束日期
	// 阈值
	MaxImageNum      int `bson:"MaxImageNums" json:"-"`     // 图片数量
	MaxImageSize     int `bson:"MaxImageSize" json:"-"`     // 图片大小
	MaxAttachNum     int `bson:"MaxAttachNum" json:"-"`     // 图片数量
	MaxAttachSize    int `bson:"MaxAttachSize" json:"-"`    // 图片大小
	MaxPerAttachSize int `bson:"MaxPerAttachSize" json:"-"` // 单个附件大小
}

// note主页需要
type UserAndBlogUrl struct {
	User
	BlogUrl string `BlogUrl`
	PostUrl string `PostUrl`
}

// 用户与博客信息结合, 公开
type UserAndBlog struct {
	UserId    bson.ObjectId `bson:"_id,omitempty"` // 必须要设置bson:"_id" 不然mgo不会认为是主键
	Email     string        `Email`                // 全是小写
	Username  string        `Username`             // 不区分大小写, 全是小写
	Logo      string        `Logo`
	BlogTitle string        `BlogTitle` // 博客标题
	BlogLogo  string        `BlogLogo`  // 博客Logo
	BlogUrl   string        `BlogUrl`   // 博客链接, 主页

	BlogUrls // 各个页面
}
