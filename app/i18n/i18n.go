package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// convert revel msg to js msg

var msgBasePath = "/Users/life/Documents/Go/package2/src/github.com/leanote/leanote/messages/"
var targetBasePath = "/Users/life/Documents/Go/package2/src/github.com/leanote/leanote/public/js/i18n/"

func parse(filename string) {
	file, err := os.Open(msgBasePath + filename)
	reader := bufio.NewReader(file)
	msg := map[string]string{}
	if err != nil {
		fmt.Println(err)
		return
	}
	for true {
		line, _, err := reader.ReadLine()

		if err != nil {
			break
		}

		if len(line) == 0 {
			continue
		}
		// 对每一行进行处理
		if line[0] == '#' || line[1] == '#' {
			continue
		}
		lineStr := string(line)

		// 找到第一个=位置
		pos := strings.Index(lineStr, "=")

		if pos < 0 {
			continue
		}

		key := string(line[0:pos])
		value := string(line[pos+1:])

		//		fmt.Println(lineStr)
		//		fmt.Println(value)

		msg[key] = value
	}

	// JSON
	b, _ := json.Marshal(msg)
	str := string(b)
	fmt.Println(str)

	targetName := targetBasePath + filename + ".js"
	file2, err2 := os.OpenFile(targetName, os.O_RDWR|os.O_CREATE, 0644)
	if err2 != nil {
		file2, err2 = os.Create(targetName)
	}
	file2.WriteString("var MSG = " + str + ";" + `
function getMsg(key, data) {
	var msg = MSG[key]
	if(msg) {
		if(data) {
			if(!isArray(data)) {
				data = [data];
			}
			for(var i = 0; i < data.length; ++i) {
				msg = msg.replace("%s", data[i]);
			}
		}
		return msg;
	}
	return key;
}`)
}

// 生成js的i18n文件
func main() {
	parse("msg.en")
	parse("msg.zh")
	parse("msg.fr")
	parse("msg.pt")
	parse("blog.zh")
	parse("blog.en")
	parse("blog.fr")
	parse("blog.pt")
}
