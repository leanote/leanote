package controllers

import (
	"github.com/revel/revel"
//	"encoding/json"
	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/info"
	"io/ioutil"
	"os"
)

// 首页
type File struct {
	BaseController
}

// 上传图片 editor
func (c File) UploadImage(renderHtml string) revel.Result {
	if renderHtml == "" {
		renderHtml = "file/image.html"
	}
	
	re := c.uploadImage();
	
	c.RenderArgs["fileUrlPath"] = "http://leanote.com" + re.Id
	c.RenderArgs["resultCode"] = re.Code
	c.RenderArgs["resultMsg"] = re.Msg

	return c.RenderTemplate(renderHtml)
}

// 上传的是博客logo
func (c File) UploadBlogLogo() revel.Result {
	return c.UploadImage("file/blog_logo.html");
}

// 拖拉上传
func (c File) UploadImageJson(renderHtml string) revel.Result {
	re := c.uploadImage();
	re.Id = "http://leanote.com" + re.Id
//	re.Id = re.Id
	return c.RenderJson(re)
}

// 上传图片, 公用方法
func (c File) uploadImage() (re info.Re) {
	var fileUrlPath = ""
	var resultCode = 0 // 1表示正常
	var resultMsg = "内部错误" // 错误信息
	
	defer func() {
		re.Id = fileUrlPath
		re.Code = resultCode
		re.Msg = resultMsg
	}()
	
	file, handel, err := c.Request.FormFile("file")
	if err != nil {
		return re
	}
	defer file.Close()
	// 生成上传路径
	fileUrlPath = "/upload/" + c.GetUserId() + "/images"
	dir := revel.BasePath + "/public/" + fileUrlPath
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		Log(err)
		return re
	}
	// 生成新的文件名
	filename := handel.Filename
	_, ext := SplitFilename(filename)
	if(ext != ".gif" && ext != ".jpg" && ext != ".png" && ext != ".bmp" && ext != ".jpeg") {
		resultMsg = "不是图片"
		return re
	}
	filename = NewGuid() + ext
	data, err := ioutil.ReadAll(file)
	if err != nil {
		LogJ(err)
		return re
	}
	
	// > 2M?
	if(len(data) > 5 * 1024 * 1024) {
		resultCode = 0
		resultMsg = "图片大于2M"
		return re
	}
	
	toPath := dir + "/" + filename;
	err = ioutil.WriteFile(toPath, data, 0777)
	if err != nil {
		LogJ(err)
		return re
	}
	// 改变成gif图片
	_, toPathGif := TransToGif(toPath, 0, true)
	
	fileUrlPath += "/" + GetFilename(toPathGif)
	resultCode = 1
	resultMsg = "上传成功!"
	
	return re
}