package controllers

import (
	"github.com/revel/revel"
//	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/lea/netutil"
	"github.com/leanote/leanote/app/info"
	"io/ioutil"
	"os"
	"fmt"
	"strconv"
	"strings"
)

// 首页
type File struct {
	BaseController
}

// 上传的是博客logo
// TODO logo不要设置权限, 另外的目录
func (c File) UploadBlogLogo() revel.Result {
	re := c.uploadImage("blogLogo", "");
	
	c.RenderArgs["fileUrlPath"] = re.Id
	c.RenderArgs["resultCode"] = re.Code
	c.RenderArgs["resultMsg"] = re.Msg

	return c.RenderTemplate("file/blog_logo.html")
}

// 拖拉上传, pasteImage
// noteId 是为了判断是否是协作的note, 如果是则需要复制一份到note owner中
func (c File) PasteImage(noteId string) revel.Result {
	re := c.uploadImage("pasteImage", "");
	
	userId := c.GetUserId()
	note := noteService.GetNoteById(noteId)
	if note.UserId != "" {
		noteUserId := note.UserId.Hex()
		if noteUserId != userId {
			// 是否是有权限协作的
			if shareService.HasUpdatePerm(noteUserId, userId, noteId) {
				// 复制图片之, 图片复制给noteUserId
				_, re.Id = fileService.CopyImage(userId, re.Id, noteUserId)				
			} else {
				// 怎么可能在这个笔记下paste图片呢?
				// 正常情况下不会
			}
		}
	}
	
	return c.RenderJson(re)
}

// 头像设置
func (c File) UploadAvatar() revel.Result {
	re := c.uploadImage("logo", "");
	
	c.RenderArgs["fileUrlPath"] = re.Id
	c.RenderArgs["resultCode"] = re.Code
	c.RenderArgs["resultMsg"] = re.Msg
	
	if re.Ok {
		re.Ok = userService.UpdateAvatar(c.GetUserId(), re.Id)
		if re.Ok {
			c.UpdateSession("Logo", re.Id);
		}
	}
	
	return c.RenderJson(re)
}

// leaui image plugin upload image
func (c File) UploadImageLeaui(albumId string) revel.Result {
	re := c.uploadImage("", albumId);
	return c.RenderJson(re)
}

// 上传图片, 公用方法
// upload image common func
func (c File) uploadImage(from, albumId string) (re info.Re) {
	var fileUrlPath = ""
	var fileId = ""
	var resultCode = 0 // 1表示正常
	var resultMsg = "内部错误" // 错误信息
	var Ok = false
	
	defer func() {
		re.Id = fileId // 只是id, 没有其它信息
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
	if(from == "logo" || from == "blogLogo") {
		fileUrlPath = "public/upload/" + c.GetUserId() + "/images/logo"
	} else {
		fileUrlPath = "files/" + c.GetUserId() + "/images"
	}
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
	
	var maxFileSize float64
	if(from == "logo") {
		maxFileSize = configService.GetUploadSize("uploadAvatarSize");
	} else if from == "blogLogo" {
		maxFileSize = configService.GetUploadSize("uploadBlogLogoSize");
	} else {
		maxFileSize = configService.GetUploadSize("uploadImageSize");
	}
	if maxFileSize <= 0 {
		maxFileSize = 1000
	}
	
	// > 2M?
	if(float64(len(data)) > maxFileSize * float64(1024*1024)) {
		resultCode = 0
		resultMsg = fmt.Sprintf("图片大于%vM", maxFileSize)
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
	if(from == "logo" || from == "blogLogo") {
		fileId = "public/upload/" + c.GetUserId() + "/images/logo/" + filename
	}
	
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

// 复制外网的图片, 成公共图片 放在/upload下
func (c File) CopyHttpImage(src string) revel.Result {
	re := info.NewRe()
	fileUrlPath := "upload/" + c.GetUserId() + "/images"
	dir := revel.BasePath + "/public/" +  fileUrlPath
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return c.RenderJson(re)
	}
	filesize, filename, _, ok := netutil.WriteUrl(src, dir)
	
	if !ok {
		re.Msg = "copy error"
		return c.RenderJson(re)
	}
	
	// File
	fileInfo := info.File{Name: filename,
		Title: filename,
		Path: fileUrlPath + "/" + filename,
		Size: filesize}
		
	id := bson.NewObjectId();
	fileInfo.FileId = id
	
	re.Id = id.Hex()
	re.Item = fileInfo.Path
	re.Ok = fileService.AddImage(fileInfo, "", c.GetUserId())
	
	return c.RenderJson(re)
}

//------------
// 过时 已弃用!
func (c File) UploadImage(renderHtml string) revel.Result {
	if renderHtml == "" {
		renderHtml = "file/image.html"
	}
	
	re := c.uploadImage("", "");
	
	c.RenderArgs["fileUrlPath"] = configService.GetSiteUrl() + re.Id
	c.RenderArgs["resultCode"] = re.Code
	c.RenderArgs["resultMsg"] = re.Msg

	return c.RenderTemplate(renderHtml)
}

// 已弃用
func (c File) UploadImageJson(from, noteId string) revel.Result {
	re := c.uploadImage(from, "");
	return c.RenderJson(re)
}
