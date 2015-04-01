# API设计

By life (life@leanote.com)

## api url

所有api的url前面带/api/, 如:

`/api/user/info?userId=xxxx&token=xxxx`

## 文件目录结构
* 所有API的Controller都在app/api文件夹下
* 文件命名: Api功能Controller.go, 如ApiUserController.go
* 结构体命名为 Api功能, 如ApiUser
* API公用Controller: ApiBaseController
* init.go 注入service和定义拦截器

## 流程
用户登录后返回一个token, 以后所有的请求都携带该token. 
在init.go中的拦截器会得到token并调用sessionService判断是否登录了 

## 返回值结构
* 全部返回JSON, JSON, 除二进制文件(图片, 附件外), 如果返回其它非JSON格式的值, 肯定是出错了
* 错误信息全部返回 {Ok: false, Msg: ""}
* 正确信息全部返回真实的返回数据, 如返回笔记:
```
{
  "NoteId": "54bdc7e305fcd13ea3000000",
  "NotebookId": "54bdc65599c37b0da9000003",
  "UserId": "54bdc65599c37b0da9000002",
  "Title": "asdfads",
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