SCRIPTPATH=$(cd "$(dirname "$0")"; pwd)
cd $SCRIPTPATH
go run . build -v ../../ ./tmptmp
rm -rf ./tmptmp