#!/bin/sh
SCRIPTPATH=$(cd "$(dirname "$0")"; pwd)

# set link
path="$SCRIPTPATH/src/github.com/leanote"
if [ ! -d "$path" ]; then 
	mkdir -p "$path" 
fi 
ln -s ../../../../ $SCRIPTPATH/src/github.com/leanote/leanote

# set GOPATH
export GOPATH=$GOPATH:$SCRIPTPATH

# run
osName=`uname`
if [ $osName == "Darwin" ]; then
	"$SCRIPTPATH/leanote-mac" -importPath github.com/leanote/leanote
else
	"$SCRIPTPATH/leanote-linux" -importPath github.com/leanote/leanote
fi
