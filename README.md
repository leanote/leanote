[中文](https://github.com/leanote/leanote#1-介绍)

## 1. Introduction

Leanote, not just a notebook!
![leanote.png](leanote.png "")

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

More information about how to install leanote please see:
* [leanote binary distribution installation tutorial](https://github.com/leanote/leanote/wiki/leanote-binary-distribution-installation-tutorial)
* [leanote develop distribution installation tutorial](https://github.com/leanote/leanote/wiki/leanote-develop-distribution-installation-tutorial)

### 3.1. Download leanote

Leanote V1.0.2-beta has been released. Binaries:

* Linux: [leanote-linux-x86_64.v1.0-beta.2.bin.tar.gz](https://github.com/leanote/leanote/releases/download/1.0-beta/leanote-linux-x86_64.v1.0-beta.2.bin.tar.gz)
* MacOS X: [leanote-mac-x86_64.v1.0-beta.2.bin.tar.gz](https://github.com/leanote/leanote/releases/download/1.0-beta/leanote-mac-x86_64.v1.0-beta.bin.2.tar.gz)

Or you can clone [Leanote bin repository](https://github.com/leanote/leanote-bin) (Recommend)

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
user2 username: admin, password: abc123 (administrator)
user3 username: demo, password: demo@leanote.com (this user is for demo)
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

``app.secret`` **required** **important**
The secret key used for cryptographic operations (revel.Sign).

FOR SECURITY, YOU MUST CHANGE IT!!

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

## Docs
* [leanote binary distribution installation tutorial](https://github.com/leanote/leanote/wiki/leanote-binary-distribution-installation-tutorial)
* [leanote develop distribution installation tutorial](https://github.com/leanote/leanote/wiki/leanote-develop-distribution-installation-tutorial)
* [leanote blog theme api](https://github.com/leanote/leanote/wiki/leanote-blog-theme-api_en)

## Discussion
* [leanote bbs](http://bbs.leanote.com)
* [leanote google group](https://groups.google.com/forum/#!forum/leanote)
* QQ Group: 158716820

-----------------------------------------------------------------------

## 1. 介绍

Leanote, 不只是笔记!
![leanote.png](leanote.png "")

**特性**

* 知识管理: 通过leanote来管理知识, leanote有易操作的界面, 包含两款编辑器tinymce和markdown. 在leanote, 你可以尽情享受写作.
* 分享: 你也可以通过分享知识给好友, 让好友拥有你的知识.
* 协作: 在分享的同时也可以与好友一起协作知识.
* 博客: leanote也可以作为你的博客, 将知识公开成博客, 让leanote把你的知识传播的更远!

## 2. 为什么我们要创建leanote?
说实话, 我们曾是evernote的忠实粉丝, 但是我们也发现evernote的不足:
* evernote的编辑器不能满足我们的需求, 不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放.
* 我们是markdown的爱好者, 可是evernote竟然没有.
* 我们也想将知识公开, 所以我们有自己的博客, 如wordpress, 但为什么这两者不能合二为一呢?
* 还有...

## 3.安装leanote
leanote是一款私有云笔记, 你可以下载它安装在自己的服务器上, 当然也可以在 http://leanote.com 上注册.

这里详细整理了leanote二进版和leanote开发版的安装教程, 请移步至:
* [leanote二进制详细安装教程](https://github.com/leanote/leanote/wiki/leanote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)
* [leanote开发版详细安装教程](https://github.com/leanote/leanote/wiki/leanote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)

### 3.1. 下载leanote

Leanote V1.0-beta.2 已发布, 二进制文件(暂时没有windows版的):

* Linux: [leanote-linux-x86_64.v1.0-beta.2.bin.tar.gz](https://github.com/leanote/leanote/releases/download/1.0-beta/leanote-linux-x86_64.v1.0-beta.2.bin.tar.gz)
* MacOS X: [leanote-mac-x86_64.v1.0-beta.2.bin.tar.gz](https://github.com/leanote/leanote/releases/download/1.0-beta/leanote-mac-x86_64.v1.0-beta.2.bin.tar.gz)

或者直接检出[Leanote bin repository](https://github.com/leanote/leanote-bin) (推荐, 因为为最新版本)

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
user2 username: admin, password: abc123 (管理员, 重要!)
user3 username: demo@leanote.com, password: demo@leanote.com (为体验使用)
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

``app.secret`` **重要**
请随意修改一个, app的密钥, 不能使用默认的, 不然会有安全问题

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

## 相关文档
* [leanote二进制版详细安装教程](https://github.com/leanote/leanote/wiki/leanote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)
* [leanote开发版详细安装教程](https://github.com/leanote/leanote/wiki/leanote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)
* [Leanote source leanote源码导读](https://github.com/leanote/leanote/wiki/Leanote-source-leanote源码导读)
* [leanote blog theme api(中文版)](https://github.com/leanote/leanote/wiki/leanote-blog-theme-api)
* [How to develop leanote 如何开发leanote](https://github.com/leanote/leanote/wiki/How-to-develop-leanote-如何开发leanote)


## 讨论
* [leanote 社区](http://bbs.leanote.com)
* QQ群: 158716820
* [leanote google group](https://groups.google.com/forum/#!forum/leanote)
