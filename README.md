## Introduction

leanote, it's a cloud note. You can create your own cloud note in leanote. 

## Features
* Knowledge: manage your knowledge in leanote. Leanote contains tinymce editor and markdown editor, just enjoy yourself and leanote.
* Share: share your knowledge to your friends in leanote. Well, you are not alone, you can invite your friends to join your cloud note and share your knowledge each other.
* Cooperation: collaborate with friends to improve your knowledge.
* Blog: public your knowledge and leanote be your blog.

## Why we create leanote
Well, to be honest, our inspiration comes from evernote, and we use evenote to manage our knowledge everyday. But we find that:
* Evernote's editor can't meet our needs, it hasn't document navigation, can't put our codes(as a programmer, put codes is the basic needs), can't resize images...)
* We like markdown, but evernote don't support it.
* We want to public our knowledge, so we have our blog(such as wordpress) and evernote, but why can't be one!
* ......

## How to use it
leanote build with golang(revel) and mongodb. so you must install mongodb at first.

### Install mongodb
go http://www.mongodb.org to download and install it.


### Export initial mongodb data
the mongodb data in path_to_leante/mongodb_backup/leanote_install_data
```
mongorestore -h localhost -d leanote --directoryperdb path_to_leante/mongodb_backup/leanote_install_data
```
the initial data contains two users
```
user1 username: leanote, password: abc123
user2 username: admin, password: abc123
```

### Configuration
rename path_to_leante/conf/app-default.conf to path_to_leante/conf/app.conf

``mongodb``  **required**

the mongodb info

```Shell
db.host=localhost
db.port=27017
db.dbname=leanote # required
db.username=root
db.password=root123
# or you can set the mongodb url for more complex needs
# mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
db.url=mongodb://root:root123@localhost:27017/leanote
```
``email``

for find password

``http.port``

default is 80

``site.url``

default is http://localhost, you must config it when your domain isn't it, it is used when upload images.

``adminUsername``

default is admin. The index site is the adminUsername's blog

for more infomation please see app/app.conf and revel manuals http://revel.github.io

### Run leanote
```
cd path_to_leanote/bin
sh run.sh
```

## How to develop leanote

leanote is a app based on revel(http://revel.github.io), so if you want to develop leanote as you want, you must be familar with revel.

### install golang

install golang and set GOPATH

### install revel
* go get github.com/revel/revel
* go get github.com/revel/cmd/revel

### get leanote
* go get github.com/leanote/leanote/app

### build/run leanote via revel

now you can modify leanote source and build/run with revel

```
revel run github.com/leanote/leanote
```

Welcome contribute your code to leanote! Thanks.

## License

leanote is released as open-source software under the GNU General Public License (GPL), version 2 or later. 
