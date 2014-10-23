#!/bin/sh

# release leanote

SP=$(cd "$(dirname "$0")"; pwd)

cd $SP
cd ../

# tmp path to store leanote release files
tmp="/Users/life/Desktop/leanote_release"

version=x86_64.v1.0-beta

rm -rf $tmp/leanote
mkdir -p $tmp/leanote/app
mkdir -p $tmp/leanote/conf
mkdir -p $tmp/leanote/bin

# bin
cp -r ./bin/src $tmp/leanote/bin/
cp  ./bin/run.sh $tmp/leanote/bin/

# LICENSE
cp ./LICENSE $tmp/leanote/

# views
cp -r ./app/views $tmp/leanote/app
# 可不要
cp -r ./app/service $tmp/leanote/app/service
cp -r ./app/controllers $tmp/leanote/app/controllers
cp -r ./app/db $tmp/leanote/app/db
cp -r ./app/info $tmp/leanote/app/info
cp -r ./app/lea $tmp/leanote/app/lea

# conf
cp ./conf/app.conf-default $tmp/leanote/conf/app.conf
cp ./conf/routes-default $tmp/leanote/conf/routes

# others
cp -r ./messages ./public ./mongodb_backup $tmp/leanote/

# delete some files
rm -r $tmp/leanote/public/tinymce/classes
rm -r $tmp/leanote/public/upload
mkdir $tmp/leanote/public/upload
mkdir $tmp/leanote/files
rm -r $tmp/leanote/public/.codekit-cache
rm $tmp/leanote/public/.DS_Store
rm $tmp/leanote/public/config.codekit

# make link
cd $tmp/leanote/bin
mkdir ./src/github.com/leanote
ln -s ../../../../ ./src/github.com/leanote/leanote

# archieve
# << 'BLOCK
# linux
cd $SP
cd ../
cp ./bin/leanote-linux $tmp/leanote/bin/
cd $tmp
tar -cvf $tmp/leanote-linux-$version.bin.tar leanote
gzip $tmp/leanote-linux-$version.bin.tar

# mac
rm $tmp/leanote/bin/leanote-linux
cd $SP
cd ../
cp ./bin/leanote-mac $tmp/leanote/bin/
cd $tmp
tar -cvf $tmp/leanote-mac-$version.bin.tar leanote
gzip $tmp/leanote-mac-$version.bin.tar

# BLOCK'