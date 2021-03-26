# http://revel.github.io/manual/tool.html
SCRIPTPATH=$(dirname "$PWD")
echo $SCRIPTPATH;
cd $SCRIPTPATH;

revel run -a .