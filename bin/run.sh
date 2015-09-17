#!/bin/sh
SCRIPTPATH=$(cd "$(dirname "$0")"; pwd)

# set link

path="$SCRIPTPATH/src/github.com/leanote"
if [ ! -d "$path" ]; then
	mkdir -p "$path"
fi
rm -rf $SCRIPTPATH/src/github.com/leanote/leanote # 先删除
ln -s ../../../../ $SCRIPTPATH/src/github.com/leanote/leanote

# set GOPATH
export GOPATH=$SCRIPTPATH

# run
osName=`uname` # Darwin or Linux
osName=`tr '[A-Z]' '[a-z]' <<<"$osName"` # toLowerCase
bit=`getconf LONG_BIT` # 32, 64
if [ $bit = "64" ]
then
	bit="amd64"
else
	bit="386"
fi

script="$SCRIPTPATH/leanote-$osName-$bit"
chmod 777 $script
$script -importPath github.com/leanote/leanote