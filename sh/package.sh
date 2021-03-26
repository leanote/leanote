# http://revel.github.io/manual/tool.html
SCRIPTPATH=$(dirname "$PWD")
echo $SCRIPTPATH;
cd $SCRIPTPATH;
revel package --run-mode=prod --target-path=sh/leanote.tar.gz -a .