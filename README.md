

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
[中文](https://github.com/leanote/leanote#1-介绍)


