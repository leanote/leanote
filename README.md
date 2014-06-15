[中文](https://github.com/leanote/leanote#1-介绍)

## 1. Introduction

Leanote, your own cloud note.

**Some Features**

* Knowledge: Manage your knowledge in leanote. leanote contains the tinymce editor and a markdown editor, just enjoy yourself writing.
* Share: Share your knowledge with your friends in leanote. You can invite your friends to join your notepad in the cloud so you can share knowledge.
* Cooperation: Collaborate with friends to improve your skills.
* Blog: Publish your knowledge and make leanote your blog.

## 2. Why we created leanote
To be honest, our inspiration comes from Evernote. We use Evernote to manage our knowledge everyday. But we find that:
* Evernote's editor can't meet our needs, it does not have document navigation, it does not render code properly (as a programmer, syntax highlighted code rendering is a basic need), it cannot resize images and so forth
* We like markdown, but Evernote does not support it.
* We want to share our knowledge, so all of us have our blogs (e.g. on Wordpress) and our Evernote accounts, but why can not those two be one!
* ......

## 3. How to install leanote

### 3.1. Download leanote

Leanote V0.3 has been released. Binaries:

* Linux: [leanote-linux-v0.3.bin.tar.gz](https://github.com/leanote/leanote/releases/download/0.3/leanote-linux-v0.3.bin.tar.gz)
* MacOS X: [leanote-mac-v0.3.bin.tar.gz](https://github.com/leanote/leanote/releases/download/0.3/leanote-mac-v0.3.bin.tar.gz)

### 3.2. Install MongoDB

Leanote is written in go using [revel](https://revel.github.io/) and [MongoDB](https://www.mongodb.org). Thus, you need to first install MongoDB.

For more tips please have a look at [our wiki](https://github.com/leanote/leanote/wiki/Install-Mongodb)

### 3.3. Import initial MongoDB data

The mongodb data is in `[PATH_TO_LEANOTE]/mongodb_backup/leanote_install_data`

```
$> mongorestore -h localhost -d leanote --directoryperdb PATH_TO_LEANOTE/mongodb_backup/leanote_install_data
```

The initial database contains two users:

```
user1 username: leanote, password: abc123
user2 username: admin, password: abc123
```

### 3.4. Configuration

Modify `[PATH_TO_LEANOTE]/conf/app.conf`. Available configuration options are:

``mongodb``  **required**

```Shell
db.host=localhost
db.port=27017
db.dbname=leanote
db.username=
db.password=
```

``http.port``

Default is 80

``site.url``

Default is `http://localhost`, you must edit this when hosting leanote anywhere else. This is used when uploading images.

``email``

For password recovery mails

``adminUsername``

Default is `admin`. The landing page is the admin user's blog.

For more infomation please see `app/app.conf` and the [revel manuals](https://revel.github.io/)

### 3.5. Run leanote

```
$> cd PATH_TO_LEANOTE/bin
$> sudo sh run.sh
```

## 4. How to develop leanote

Please see [How-to-develop-leanote](https://github.com/leanote/leanote/wiki/How-to-develop-leanote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91leanote)


## 5. Contributors
Thank you to all the [contributors](https://github.com/leanote/leanote/graphs/contributors) on
this project. Your help is much appreciated.

## 6. Contributing

Please fork this repository and contribute back using [pull requests](https://github.com/leanote/leanote/pulls).

## Discussion
[leanote google group](https://groups.google.com/forum/#!forum/leanote)

-----------------------------------------------------------------------

## 1. 介绍

Leanote, 你的私人云笔记!

**特性**

* 知识管理: 通过leanote来管理知识, leanote有易操作的界面, 包含两款编辑器tinymce和markdown. 在leanote, 你可以尽情享受写作.
* 分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.
* 协作: 在分享的同时也可以与好友一起协作知识.
* 博客: leanote也可以作为你的博客, 将知识公开成博客, 让leanote把你的知识传播的更远!

## 2. 为什么我们要创建leanote?
说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:
* leanote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.
* 我们是markdown的爱好者, 可是evernote竟然没有.
* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?
* 还有...

## 3.安装leanote
leanote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 http://leanote.com 上注册.

### 3.1. 下载leanote

Leanote V0.3 已发布, 二进制文件(暂时没有windows版的):

* Linux: [leanote-linux-v0.3.bin.tar.gz](https://github.com/leanote/leanote/releases/download/0.3/leanote-linux-v0.3.bin.tar.gz)
* MacOS X: [leanote-mac-v0.3.bin.tar.gz](https://github.com/leanote/leanote/releases/download/0.3/leanote-mac-v0.3.bin.tar.gz)

### 3.2. 安装 MongodbDB

Leanote是由golang(使用[revel](https://revel.github.io/)框架 和 [MongoDB](https://www.mongodb.org)数据库), 你需要先安装Mongodb.

安装MongodbDB, 导入数据更多细节请查看: [wiki](https://github.com/leanote/leanote/wiki/Install-Mongodb)

### 3.3. 导入初始数据

MongodbDB初始数据在 `[PATH_TO_LEANOTE]/mongodb_backup/leanote_install_data`

```
$> mongorestore -h localhost -d leanote --directoryperdb PATH_TO_LEANOTE/mongodb_backup/leanote_install_data
```

初始数据包含两个用户:

```
user1 username: leanote, password: abc123
user2 username: admin, password: abc123
```

### 3.4. 配置

修改 `[PATH_TO_LEANOTE]/conf/app.conf`. 有以下选项:

``mongodb``  **必须配置!**

```Shell
db.host=localhost
db.port=27017
db.dbname=leanote
db.username=
db.password=
```

``http.port``

默认为 80

``site.url``

默认是 `http://localhost`, 这会在上传图片后的图片路径中用户, 还有发邮件, 找回密码验证邮箱时用到.

``email``

找回密码和验证邮箱时使用

``adminUsername``

默认是 `admin`. 首页即为该用户的博客

更多配置请查看 `app/app.conf` 和 [revel 手册](https://revel.github.io/)

### 3.5. 运行leanote

```
$> cd PATH_TO_LEANOTE/bin
$> sudo sh run.sh
```

## 4. 如何对leanote进行二次开发

请查看 [How-to-develop-leanote](https://github.com/leanote/leanote/wiki/How-to-develop-leanote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91leanote)

## 5. 贡献者
多谢 [贡献者](https://github.com/leanote/leanote/graphs/contributors) 的贡献, leanote因有你们而更完美!

## 6. 加入我们

欢迎提交[pull requests](https://github.com/leanote/leanote/pulls) 到leanote.

leanote还有很多问题, 如果你喜欢它, 欢迎加入我们一起完善leanote.

## 讨论
[leanote google group](https://groups.google.com/forum/#!forum/leanote)