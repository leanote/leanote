#!/bin/sh

# release leanote

# 当前路径
SP=$(cd "$(dirname "$0")"; pwd)

# tmp path to store leanote release files
tmp="/Users/life/Desktop/leanote_release"

# version
V="v1.4.2"

##=================================
# 1. 先build 成 3个平台, 2种bit = 6种
##=================================

# cd /Documents/Go/package2/src/github.com/leanote/leanote/bin
# GOOS=darwin GOARCH=amd64 go build -o leanote-darwin-amd64 ../app/tmp

cd $SP

# $1 = darwin, linux
# $2 = amd64
function build()
{
	echo build-$1-$2
	if [ $1 = "linux" -o $1 = "darwin" ]
	then
		suffix=""
	else
		suffix=".exe"
	fi
	
	GOOS=$1 GOARCH=$2 go build -o leanote-$1-$2$suffix ../app/tmp
}

build "linux" "386";
build "linux" "amd64";
build "linux" "arm";

build "windows" "386";
build "windows" "amd64";

# build "darwin" "386";
build "darwin" "amd64";


##======================
# 2. release目录准备工作
##======================
rm -rf $tmp/leanote
mkdir -p $tmp/leanote/app
mkdir -p $tmp/leanote/conf
mkdir -p $tmp/leanote/bin

##==================
# 3. 复制
##==================

cd $SP
cd ../

# bin
cp -r ./bin/src $tmp/leanote/bin/
cp  ./bin/run.sh $tmp/leanote/bin/
cp  ./bin/run.bat $tmp/leanote/bin/

# views
cp -r ./app/views $tmp/leanote/app
# 可不要, 源码
#cp -r ./app/service $tmp/leanote/app/service
#cp -r ./app/controllers $tmp/leanote/app/controllers
#cp -r ./app/db $tmp/leanote/app/db
#cp -r ./app/info $tmp/leanote/app/info
#cp -r ./app/lea $tmp/leanote/app/lea

# conf
cp ./conf/app.conf $tmp/leanote/conf/app.conf
cp ./conf/routes $tmp/leanote/conf/routes
# 处理app.conf, 还原配置
cat $tmp/leanote/conf/app.conf | sed 's/db.dbname=leanote.*#/db.dbname=leanote #/' > $tmp/leanote/conf/app.conf2 # 不能直接覆盖
rm $tmp/leanote/conf/app.conf
mv $tmp/leanote/conf/app.conf2 $tmp/leanote/conf/app.conf

# others
cp -r ./messages ./public ./mongodb_backup $tmp/leanote/

# delete some files
rm -r $tmp/leanote/public/tinymce/classes
rm -r $tmp/leanote/public/upload
mkdir $tmp/leanote/public/upload
rm -r $tmp/leanote/public/.codekit-cache
rm $tmp/leanote/public/.DS_Store
rm $tmp/leanote/public/config.codekit

# make link
# cd $tmp/leanote/bin
# ln -s ../../../../ ./src/github.com/leanote/leanote

# archieve
# << 'BLOCK

##===========
# 4. 打包
##===========
# $1 = linux
# $2 = 386, amd64

# 创建一个$V的目录存放之
rm -rf $tmp/$V
mkdir $tmp/$V

function tarRelease()
{
	echo tar-$1-$2
	cd $SP
	cd ../
	rm $tmp/leanote/bin/leanote-* # 删除之前的bin文件
	rm $tmp/leanote/bin/run* # 删除之前的run.sh 或 run.bat
	
	if [ $1 = "linux" -o $1 = "darwin" ]
	then
		suffix=""
		if [ $2 = "arm" ]
		then
			cp ./bin/run-arm.sh $tmp/leanote/bin/run.sh
		else
			cp ./bin/run-$1-$2.sh $tmp/leanote/bin/run.sh
		fi
	else
		cp ./bin/run.bat $tmp/leanote/bin/
		suffix=".exe"
	fi
	
	cp ./bin/leanote-$1-$2$suffix $tmp/leanote/bin/
	cd $tmp
	tar -cf $tmp/$V/leanote-$1-$2-$V.bin.tar leanote
	gzip $tmp/$V/leanote-$1-$2-$V.bin.tar
}

tarRelease "linux" "386";
tarRelease "linux" "amd64";
tarRelease "linux" "arm";

tarRelease "windows" "386";
tarRelease "windows" "amd64";

# tarRelease "darwin" "386";
tarRelease "darwin" "amd64";

# BLOCK'
