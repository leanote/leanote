# Vendor

## Install leanote_revel cmd

基于revel 0.18

```
go install github.com/leanote/leanote/cmd/leanote_revel
leanote_revel run github.com/leanote/leanote
````

## build leanote

在当前目录生成了leanote二进制文件

```
go build -o ./leanote github.com/leanote/leanote/app/tmp
```

## 运行leanote

其中-importPath是必须的

```
./leanote -importPath=github.com/leanote/leanote -runMode=prod -port=9000
```