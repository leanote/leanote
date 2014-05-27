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

Please see [How-to-develop-leanote](https://github.com/leanote/leanote/wiki/How-to-develop-leanote)

## Discussion
[leanote google group](https://groups.google.com/forum/#!forum/leanote)