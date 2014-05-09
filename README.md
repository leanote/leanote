## Introduction

Leanote, a cloud note. You can create your own cloud note by leanote. 

## Features
* Knowledge: manage your knowledge in leanote. leanote contains tinymce editor and markdown editor, just enjoy yourself in writting.
* Share: share your knowledge to your friends in leanote. Well, you are not alone, you can invite your friends to join your cloud note and share your knowledge each other.
* Cooperation: collaborate with friends to improve your knowledge.
* Blog: public your knowledge and leanote be your blog.

## Why we create leanote
To be honest, our inspiration comes from evernote, and we use evenote to manage our knowledge everyday. But we find that:
* Evernote's editor can't meet our needs, it hasn't document navigation, can't put our codes(as a programmer, put codes is the basic needs), can't resize images...)
* We like markdown, but evernote don't support it.
* We want to public our knowledge, so we have our blog(such as wordpress) and evernote, but why can't be the one!
* ......

## How to use it

## Download leanote

Leanote v0.1 has released.

* linux: https://github.com/leanote/leanote/releases/download/0.1/leanote-linux-v0.1.bin.tar.gz
* mac ox: https://github.com/leanote/leanote/releases/download/0.1/leanote-mac-v0.1.tar.gz

### Install mongodb

Leanote build with golang(revel) and mongodb. so you must install mongodb at first.

For more tips please go https://github.com/leanote/leanote/wiki/mongodb-in-leanote

Go http://www.mongodb.org to download and install it.

### Export initial mongodb data

The mongodb data is in path_to_leante/mongodb_backup/leanote_install_data

```
$> mongorestore -h localhost -d leanote --directoryperdb path_to_leante/mongodb_backup/leanote_install_data
```

The initial data contains two users:

```
user1 username: leanote, password: abc123
user2 username: admin, password: abc123
```

### Configuration

Modify path_to_leanote/conf/app.conf, the configuration options contains:

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

Default is http://localhost, you must config it when your domain isn't it, it is used when upload images.

``email``

for find password

``adminUsername``

Default is admin. The index site is the adminUsername's blog

For more infomation please see app/app.conf and revel manuals http://revel.github.io

### Run leanote

```
$> cd path_to_leanote/bin
$> sudo sh run.sh
```

## How to develop leanote

For more tips please go https://github.com/leanote/leanote/wiki/How-to-develop-leanote

Leanote is a app based on revel(http://revel.github.io), so if you want to develop leanote as you want, you must be familar with revel.

### Install golang

Install golang and set GOPATH

### Install revel
```
$> go get github.com/revel/revel
$> go get github.com/revel/cmd/revel
```

### Get leanote

```
$> go get github.com/leanote/leanote/app
```

### Build/Run leanote via revel

```
$> cp conf/app.conf-default to conf/app.conf
$> cp conf/routes-default to conf/routes
```

Now you can modify leanote source and build/run with revel

```
$> revel run github.com/leanote/leanote
```

Welcome to join with us and contribute your code to leanote! Thanks.