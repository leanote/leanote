# Vendor

脱离 revel 每次升级不向下兼容 的控制!! 基于revel 0.18

## Install leanote_revel cmd

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