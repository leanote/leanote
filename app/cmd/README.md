全部代码来自https://github.com/revel/cmd

因为要改parse2, 所以改只要一点点代码
harness/
    build.go 只要gensource, 其它的先return
main.go 改动很小
build.go 改动很小
parser2/
    source_processors.go 改了 fsWalk 过滤掉 public, files, build 等文件夹