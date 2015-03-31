# API 列表

## 前言

### api url

所有api的url前面带/api/, 如:

`/api/user/info?userId=xxxx&token=xxxx`

除了/auth/login, /auth/register外其它的都需要另外带参数token=xxxx

### 文件目录结构
* 所有API的Controller都在app/api文件夹下
* 文件命名: Api功能Controller.go, 如ApiUserController.go
* 结构体命名为 Api功能, 如ApiUser
* API公用Controller: ApiBaseController
* init.go 注入service和定义拦截器

### 流程
用户登录后返回一个token, 以后所有的请求都携带该token. 
在init.go中的拦截器会得到token并调用sessionService判断是否登录了 

## 返回值结构
* 全部返回JSON, JSON, 除二进制文件(图片, 附件外), 如果返回其它非JSON格式的值, 肯定是出错了
* 错误信息全部返回 {Ok: false, Msg: "相应的错误信息"}
* 正确信息返回分两种:
1. 一些操作型的api, 比如updateUsername, updateLogo之类的, 成功后返回 {Ok: true, Msg:""}
2. 一些获取型api, 如getNote, 全部返回真实的返回数据, 如返回笔记:
```
{
  "NoteId": "54bdc7e305fcd13ea3000000",
  "NotebookId": "54bdc65599c37b0da9000003",
  "UserId": "54bdc65599c37b0da9000002",
  "Title": "笔记标题",
  "Desc": "",
  "Tags": null,
  "Abstract": "",
  "Content": "",
  "IsMarkdown": false,
  "IsBlog": false,
  "IsTrash": true,
  "IsDeleted": false,
  "Usn": 15,
  "Files": [],
  "CreatedTime": "2015-01-20T11:13:41.34+08:00",
  "UpdatedTime": "2015-01-20T11:13:41.34+08:00",
  "PublicTime": "0001-01-01T00:00:00Z"
}
```
* 时间类型全是返回如 "2015-01-20T11:13:41.34+08:00" 的格式, (因为golang的time转成Json就是这样, 历史原因)

-----------------

## Auth 登录与注册

### /auth/login 登录 

```
参数: email, pwd
Method: GET
返回: 
错误: {"Ok":false, "Msg":"用户名或密码有误"}
正确: 比如:
{
	"Ok":true,
	"Token":"5500830738f41138e90003232",
	"UserId":"52d26b4e99c37b609a000001",
	"Email":"leanote@leanote.com",
	"Username":"leanote"
}
```

登录成功后将使用token作为之后的请求

### /auth/logout 注销
```
参数: token
Method: GET
返回:
错误: {Ok: false, Msg: ""}
成功: {Ok: true, Msg: ""}
```

### /auth/register 注册
```
参数: email, pwd
Method: POST
返回:
错误: {Ok: false, Msg: ""}
成功: {Ok: true, Msg: ""}
```

## User 用户
### /user/info 获取用户信息
```
参数: userId
Method: GET
返回:
错误: {Ok: false, Msg: ""}
成功: type.User
```

### /user/updateUsername 修改用户名
```
参数: username(新用户名)
Method: POST
返回: 
错误: {Ok: false, Msg: ""}
成功: {Ok: true, Msg: ""}
```

### /user/updatePwd 修改密码
```
参数: oldPwd(旧密码), pwd(新密码)
Method: POST
返回:
错误: {Ok: false, Msg: ""}
成功: {Ok: true, Msg: ""}
```

### /user/updateLogo 修改头像
```
参数: file(文件)
Method: POST
返回:
错误: {Ok: false, Msg: ""}
成功: {Ok: true, Msg: ""}
```

### /user/getSyncState 获取最新同步状态
```
参数: 无
Method: POST
返回:
错误: {Ok: false, Msg: ""}
成功: {LastSyncUsn: 3232, LastSyncTime: "上次同步时间"(暂时无用)}
```

-----

## Notebook 笔记本

### /notebook/getSyncNotebooks 得到需要同步的笔记本 
```
参数: afterUsn(int, 在此usn后的笔记本是需要同步的), maxEntry(int, 最大要同步的量)
Method: GET
返回: 
错误: {Ok: false, Msg: ""}
成功: [type.Notebook] 数组
```

### /notebook/getNotebooks 得到所有笔记本 
```
无参数
Method: GET
返回: 
错误: {Ok: false, Msg: ""}
成功: [type.Notebook] 数组
```

### /notebook/addNotebook 添加笔记本
```
参数: title(string), parentNotebookId(string, 父notebookId, 可空), seq(int) 排列
Method: POST
返回: 
错误: {Ok: false, Msg:""}
成功: type.Notebook
```

### /notebook/updateNotebook 修改笔记本
```
参数: notebookId, title, parentNotebookId, seq(int), usn(int)
Method: POST
返回: 
错误: {Ok: false, msg: ""} msg == "conflict" 表示冲突
成功: type.Notebook
```

### /notebook/deleteNotebook 删除笔记本
```
参数: notebookId, usn(int)
Method: GET
返回: 
错误: {Ok: false, msg: ""} msg == "conflict" 表示冲突
成功: {Ok: true}
```

----

## Note 笔记


### /note/getSyncNotes 获取需要同步的笔记
```
参数: afterUsn(int, 在此usn后的笔记是需要同步的), maxEntry(int, 最大要同步的量)
Method: GET
返回: 
错误: {Ok: false, Msg: ""}
成功: [type.Note] 数组, 笔记不包含Abstract和Content
```


### /note/getNotes 获得某笔记本下的笔记(无内容)
```
参数: notebookId
Method: GET
返回: 
错误: {Ok: false, Msg: ""}
成功: [type.Note] 数组, 笔记不包含Abstract和Content
```

### /note/getNoteAndContent 获得笔记与内容
```
参数: noteId
Method: GET
返回: 
错误: {Ok: false, Msg: ""}
成功: type.Note
```

### /note/getNoteContent 获得笔记内容
```
参数: noteId
Method: GET
返回: 
错误: {Ok: false, Msg: ""}
成功: type.NoteContent
```

### /note/addNote 添加笔记
```
参数: (注意首字大写)
	NotebookId string 必传
	Title string 必传
	Tags []string 可选
	Content string 必传
	Abstract string 可选, 当是markdown笔记时必须传
	IsMarkdown bool 可选
	Files []type.NoteFiles 数组  可选
Method: POST
返回: 
错误: {Ok: false, Msg: ""}
成功: type.Note, 不包含Abstract, Content
```

** 关于笔记中的图片/附件**
客户端应该添加一个"图片/附件表"来存元数据, 图片应该要缓存到本地, 附件可在需要的时候再调用相应api获取.

Content中的数据, 图片,附件在Leanote的链接是, 所以, 不管你在本地的笔记中是以什么形式来保存图片,附件的链接的,请addNote和updateNote时务必将链接修改成leanote服务器上的链接.
http://leanote.com/api/file/getImage?fileId=xx
单个附件: 
http://leanote.com/api/file/getAttach?fileId=xx 
所有附件:
http://leanote.com/api/file/getAllAttachs?noteId=xxx 
```

**注意:**
addNote时必须要把Files, 和相关的图片/附件一起传到服务器中
其中Files(文件的元数据)和其它字段以POST方式传出, 而真正数据则以http的multipart传入, 每个文件的name为"FileDatas[LocalFileId]"

图片在笔记内的链接必须为: http://leanote.com/api/file/getImage?fileId=LocalFileId或FileId
附件如果插入到了笔记内容内, 其链接必须为: http://leanote.com/api/file/getAttach?fileId=LocalFileId或FileId
其中, fileId为文件在本地的LocalFileId或服务器上的FileId

服务器端会生成FileId传给Client. Client在本地必须要建立LocalFileId与FileId的关联.

如果笔记内容传了, 且笔记内有图片, 则必须要传Files 文件元数据, 因为Server端还要对内容内的图片, 附件链接进行修改, 可能你传过去的是LocalFileId, 服务器端会将LocalFileId替换成FileId存到数据库中.

同样适用于 updateNote

http://leanote.com 不绝对, 因为用户可以自建服务, 所以在开发时需要可配置

### /note/updateNote 更新笔记
当更新了笔记某个属性时, 只要传某个属性就行, 其它不用传, 比如把笔记拉入了trash, 那么就传IsTrash: true
```
参数: (注意首字大写)
	NoteId string 必传
	Usn int 必传
	NotebookId string 可选
	Title string 可选
	Tags []string 可选
	Content string 可选
	Abstract string 可选, 当是markdown笔记时必须传
	IsMarkdown bool 可选	
	IsTrash bool 是否是trash 可选
	Files []type.NoteFiles 数组  可选
Method: POST
返回: 
错误: {Ok: false, msg: ''} msg == 'conflict' 表示冲突
成功: type.Note, 不包含Abstract和Content
```

### /note/deleteTrash 彻底删除笔记
```
参数: noteId, usn
Method: GET
返回: 
错误: {Ok: false, msg: ''} msg == 'conflict' 表示冲突
成功: type.UpdateRet
```

-------

## Tag 标签

### /tag/getSyncTags 获取需要同步的标签
```
参数: afterUsn(int, 在此usn后的标签是需要同步的), maxEntry(int, 最大要同步的量)
Method: GET
返回: 
错误: {Ok: false, Msg: ""}
成功: [type.Tag] 数组
```

### /tag/addTag 添加标签
```
参数: tag(string)
Method: POST
返回: 
错误: {Ok: false, Msg: ""}
成功: type.Tag
```

### /tag/deleteTag 删除标签 
```
参数: tag(string)
Method: POST
返回: 
错误: {Ok: false, Msg: ""}
成功: type.UpdateRet
```

### File 文件(获取图片, 附件)

### /file/getImage 获取图片
```
参数: fileId
Method: GET
返回: 
错误: 非二进制文件数据
成功: 二进制文件
```

### /file/getAttach 获取附件
```
参数: fileId
Method: GET
返回: 
错误: 非二进制文件数据
成功: 二进制文件
```

### /file/getAllAttachs 获取所有附件
```
参数: noteId
Method: GET
返回: 
错误: 非二进制文件数据
成功: 二进制文件
```

--------

## 数据类型

### type.User 用户信息

```
User {
	UserId  string
	Username string
	Email string
	Verified bool
	Logo string
}
```

### type.Notebook 笔记本

```
Notebook {
	NotebookId        
	UserId           
	ParentNotebookId // 上级
	Seq              int // 排序
	Title            string 
	IsBlog           bool  
	CreatedTime      time.Time   
	UpdatedTime      time.Time 
	
	// 更新序号
	Usn int  // UpdateSequenceNum 
}
```

### type.Note 笔记
```
Note {
	NoteId     string
	NotebookId string
	UserId     string
	Title      string
	Tags       []string
	Content    string
	IsMarkdown bool
	IsBlog     bool 
	IsTrash bool
	Files []NoteFile // 图片, 附件
	CreatedTime time.Time
	UpdatedTime time.Time
	PublicTime time.Time
	
	// 更新序号
	Usn int
}
```

### type.NoteContent 笔记内容
```
NoteContent {
	NoteId string
	UserId string
	Content string
}
```

### type.NoteFile 笔记文件(图片,附件)
```
NoteFile {
	FileId string // 服务器端Id
	LocalFileId string // 客户端Id
	Type string // images/png, doc, xls, 根据fileName确定
	Title string
	HasBody bool // 传过来的值是否要更新内容, 如果有true, 则必须传文件
	IsAttach bool // 是否是附件, 不是附件就是图片
}
```

### type.Tag 标签
```
Tag {
	TagId string
	UserId string
	Tag string 
	CreatedTime
	UpdatedTime
	IsDeleted bool  // 删除位
	// 更新序号
	Usn 
}
```

### type.UpdateRe 更新后返回的值, 包含Usn
```
ReUpdate {
	Ok bool
	Msg string
	
	// 更新序号
	Usn int
}
```
