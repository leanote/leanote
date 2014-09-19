package controllers

import (
	"github.com/revel/revel"
//	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/info"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 首页
type File struct {
	BaseController
}

// 上传图片 editor
// 过时
func (c File) UploadImage(renderHtml string) revel.Result {
	if renderHtml == "" {
		renderHtml = "file/image.html"
	}
	
	re := c.uploadImage("", "");
	
	c.RenderArgs["fileUrlPath"] = siteUrl + re.Id
	c.RenderArgs["resultCode"] = re.Code
	c.RenderArgs["resultMsg"] = re.Msg

	return c.RenderTemplate(renderHtml)
}

// 上传的是博客logo
// TODO logo不要设置权限, 另外的目录
func (c File) UploadBlogLogo() revel.Result {
	return c.UploadImage("file/blog_logo.html");
}

// 拖拉上传, pasteImage
func (c File) UploadImageJson(renderHtml, from string) revel.Result {
	re := c.uploadImage(from, "");
	return c.RenderJson(re)
}

// leaui image plugin
func (c File) UploadImageLeaui(albumId string) revel.Result {
	re := c.uploadImage("", albumId);
	return c.RenderJson(re)
}

// 上传图片, 公用方法
func (c File) uploadImage(from, albumId string) (re info.Re) {
	var fileUrlPath = ""
	var fileId = ""
	var resultCode = 0 // 1表示正常
	var resultMsg = "内部错误" // 错误信息
	var Ok = false
	
	defer func() {
		re.Id = fileId
		re.Code = resultCode
		re.Msg = resultMsg
		re.Ok = Ok
	}()
	
	file, handel, err := c.Request.FormFile("file")
	if err != nil {
		return re
	}
	defer file.Close()
	// 生成上传路径
	fileUrlPath = "files/" + c.GetUserId() + "/images"
	dir := revel.BasePath + "/" +  fileUrlPath
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return re
	}
	// 生成新的文件名
	filename := handel.Filename
	
	var ext string;
	if from == "pasteImage" {
		ext = ".png"; // TODO 可能不是png类型
	} else {
		_, ext = SplitFilename(filename)
		if(ext != ".gif" && ext != ".jpg" && ext != ".png" && ext != ".bmp" && ext != ".jpeg") {
			resultMsg = "不是图片"
			return re
		}
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
	filename = GetFilename(toPathGif)
	filesize := GetFilesize(toPathGif)
	fileUrlPath += "/" + filename
	resultCode = 1
	resultMsg = "上传成功!"
	
	// File
	fileInfo := info.File{Name: filename,
		Title: handel.Filename,
		Path: fileUrlPath,
		Size: filesize}
		
	id := bson.NewObjectId();
	fileInfo.FileId = id
	fileId = id.Hex()
	Ok = fileService.AddImage(fileInfo, albumId, c.GetUserId())
	
	fileInfo.Path = ""; // 不要返回
	re.Item = fileInfo
	
	return re
}

// get all images by userId with page
func (c File) GetImages(albumId, key string, page int) revel.Result {
	re := fileService.ListImagesWithPage(c.GetUserId(), albumId, key, page, 12)
	return c.RenderJson(re)
}

func (c File) UpdateImageTitle(fileId, title string) revel.Result {
	re := info.NewRe()
	re.Ok = fileService.UpdateImageTitle(c.GetUserId(), fileId, title)
	return c.RenderJson(re)
}

func (c File) DeleteImage(fileId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = fileService.DeleteImage(c.GetUserId(), fileId)
	return c.RenderJson(re)
}

// update image uploader to leaui image, 
// scan all user's images and insert into db
func (c File) UpgradeLeauiImage() revel.Result {
	re := info.NewRe()
	
	if ok, _ := revel.Config.Bool("upgradeLeauiImage"); !ok {
		re.Msg = "Not allowed"
		return c.RenderJson(re)
	}
	
	uploadPath := revel.BasePath + "/public/upload";
	userIds := ListDir(uploadPath)
	if userIds == nil {
		re.Msg = "no user"
		return c.RenderJson(re)
	}
	
	msg := "";
	
	for _, userId := range userIds {
		dirPath := uploadPath + "/" + userId +  "/images"
		images := ListDir(dirPath)
		if images == nil {
			msg += userId + " no images "
			continue;
		}
		
		hadImages := fileService.GetAllImageNamesMap(userId)
		
		i := 0
		for _, filename := range images {
			if _, ok := hadImages[filename]; !ok {
				fileUrlPath := "/upload/" + userId + "/images/" + filename
				fileInfo := info.File{Name: filename,
					Title: filename,
					Path: fileUrlPath,
					Size: GetFilesize(dirPath + "/" + filename)}
				fileService.AddImage(fileInfo, "", userId)
				i++
			}
		}
		msg += userId + ": " + strconv.Itoa(len(images)) + " -- " + strconv.Itoa(i) + " images "
	}
	
	re.Msg = msg
	return c.RenderJson(re)
}

//-----------

// 输出image
// 权限判断
func (c File) OutputImage(noteId, fileId string) revel.Result {
	path := fileService.GetFile(c.GetUserId(), fileId); // 得到路径
	if path == "" {
		return c.RenderText("")
	}
	fn := revel.BasePath + "/" +  strings.TrimLeft(path, "/")
    file, _ := os.Open(fn)
    return c.RenderFile(file, revel.Inline) // revel.Attachment
}

// 协作时复制图片到owner
func (c File) CopyImage(userId, fileId, toUserId string) revel.Result {
	re := info.NewRe()
	
	re.Ok, re.Id = fileService.CopyImage(userId, fileId, toUserId)
	
	return c.RenderJson(re)
}
