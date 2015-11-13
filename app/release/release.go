package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	//	"time"
)

/*
用golang exec 总是说找不到uglifyjs命令, 需要全部路径
而且node, npm要在/usr/bin下, 已建ln
*/

/*
<script src="js/jquery-1.9.0.min.js"></script>

<!-- 以后将所有的js压缩合并成一个文件 -->
<script src="js/jquery-cookie.js"></script>
<script src="js/bootstrap.js"></script>
<script type="text/javascript" src="tinymce/tinymce.js"></script>
<script src="js/common.js"></script>
<script src="js/app/note.js"></script>
<script src="js/app/tag.js"></script>
<script src="js/app/notebook.js"></script>
<script src="js/app/share.js"></script>
<script src="js/object_id.js"></script>
<script type="text/javascript" src="js/ZeroClipboard/ZeroClipboard.js"></script>
*/

//var jss = []string{"js/jquery-cookie", "js/bootstrap"}
var jss = []string{"js/jquery-cookie", "js/bootstrap",
	"js/common", "js/app/note", "js/app/tag", "js/app/notebook", "js/app/share",
	"js/object_id", "js/ZeroClipboard/ZeroClipboard"}

var base1 = "/Users/life/Documents/Go/package2/src/github.com/leanote/leanote/"
var base = "/Users/life/Documents/Go/package2/src/github.com/leanote/leanote/public/"
var cmdPath = "/usr/local/bin/uglifyjs"

func cmdError(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(os.Stderr, "The command failed to perform: %s (Command: %s, Arguments: %s)", err, "", "")
	} else {
		fmt.Println("OK")
	}
}

// filename没有扩展名
func compressJs(filename string) {
	source := base + filename + ".js"
	to := base + filename + "-min.js"
	cmd := exec.Command(cmdPath, source, "-o", to)
	_, err := cmd.CombinedOutput()
	fmt.Println(source)
	cmdError(err)
}

func combineJs() {
	// 生成一个总文件
	cmd := exec.Command("rm", base+"js/all.js")
	_, err := cmd.CombinedOutput()
	cmdError(err)

	for _, js := range jss {
		to := base + js + "-min.js"
		fmt.Println(to)
		compressJs(js)

		// 每个压缩后的文件放入之
		cmd2 := exec.Command("/bin/sh", "-c", "cat "+to+" >> "+base+"js/all.js")
		_, err := cmd2.CombinedOutput()
		cmdError(err)
		cmd2 = exec.Command("/bin/sh", "-c", "cat \n >> "+base+"js/all.js")
		_, err = cmd2.CombinedOutput()
		cmdError(err)
	}
}

// 改note-dev->note
func dev() {
	// 即替换note.js->note-min.js
	m := map[string]string{"tinymce.dev.js": "tinymce.min.js",
		"tinymce.js":                            "tinymce.min.js",
		"jquery.ztree.all-3.5.js":               "jquery.ztree.all-3.5-min.js",
		"note.js":                               "note-min.js",
		"app.js":                                "app-min.js",
		"page.js":                               "page-min.js",
		"common.js":                             "common-min.js",
		"notebook.js":                           "notebook-min.js",
		"share.js":                              "share-min.js",
		"tag.js":                                "tag-min.js",
		"jquery.slimscroll.js":                  "jquery.slimscroll-min.js",
		"jquery.contextmenu.js":                 "jquery.contextmenu-min.js",
		"editor/editor.js":                      "editor/editor-min.js",
		"/public/mdeditor/editor/scrollLink.js": "/public/mdeditor/editor/scrollLink-min.js",
		"console.log(o);":                       "",
	}
	path := base1 + "/src/views/note/note-dev.html"
	target := base1 + "/src/views/note/note.html"

	bs, _ := ioutil.ReadFile(path)
	content := string(bs)
	print(content)
	for key, value := range m {
		content = strings.Replace(content, key, value, -1)
	}

	//	var time = time.Now().Unix() % 1000

	//	content = strings.Replace(content, "-min.js", fmt.Sprintf("-min.js?r=%d", time), -1)
	//	content = strings.Replace(content, "default{{end}}.css", fmt.Sprintf("default{{end}}.css?r=%d", time), 1)
	//	content = strings.Replace(content, "writting-overwrite.css", fmt.Sprintf("writting-overwrite.css?r=%d", time), 1)

	ioutil.WriteFile(target, []byte(content), os.ModeAppend)
}

// 压缩js成一块
func tinymce() {
	//	cmdStr := "node_modules/jake/bin/cli.js minify bundle[themes:modern,plugins:table,paste,advlist,autolink,link,image,lists,charmap,hr,searchreplace,visualblocks,visualchars,code,nav,tabfocus,contextmenu,directionality,codemirror,codesyntax,textcolor,fullpage]"
	//	cmd := exec.Command("/Users/life/Documents/eclipse-workspace/go/leanote_release/tinymce-master/node_modules/jake/bin/cli.js", "minify", "bundle[themes:modern,plugins:table,paste,advlist,autolink,link,image,lists,charmap,hr,searchreplace,visualblocks,visualchars,code,nav,tabfocus,contextmenu,directionality,codemirror,codesyntax,textcolor,fullpage]")
	cmd := exec.Command("/bin/sh", "-c", "grunt minify")
	cmd.Dir = base + "/tinymce_4.1.9"

	fmt.Println("正在build tinymce")

	// 必须要先删除
	cmd2 := exec.Command("/bin/sh", "-c", "rm "+cmd.Dir+"/js/tinymce/tinymce.dev.js")
	cmd2.CombinedOutput()
	cmd2 = exec.Command("/bin/sh", "-c", "rm "+cmd.Dir+"/js/tinymce/tinymce.jquery.dev.js")
	c, _ := cmd2.CombinedOutput()
	fmt.Println(string(c))

	c, _ = cmd.CombinedOutput()
	fmt.Println(string(c))
}

func main() {
	// 压缩tinymce
	// tinymce()

	dev()

	// 其它零散的需要压缩的js
	otherJss := []string{"js/main", "js/app/page", "js/contextmenu/jquery.contextmenu",
		"js/jquery.ztree.all-3.5",
		"js/jQuery-slimScroll-1.3.0/jquery.slimscroll",
	}

	for _, js := range otherJss {
		compressJs(js)
	}

	// 先压缩后合并
	combineJs()

}
