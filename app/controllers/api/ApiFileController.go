package api

import (
	"github.com/revel/revel"
	//	"encoding/json"
	//	. "github.com/leanote/leanote/app/lea"
	//	"gopkg.in/mgo.v2/bson"
	//	"github.com/leanote/leanote/app/lea/netutil"
	//	"github.com/leanote/leanote/app/info"
	//	"io/ioutil"
	"os"
	//	"strconv"
	"archive/tar"
	"compress/gzip"
	"io"
	"strings"
	"time"
)

// 文件操作, 图片, 头像上传, 输出

type ApiFile struct {
	ApiBaseContrller
}

/*
// 协作时复制图片到owner
func (c ApiFile) CopyImage(userId, fileId, toUserId string) revel.Result {
	re := info.NewRe()

	re.Ok, re.Id = fileService.CopyImage(userId, fileId, toUserId)

	return c.RenderJson(re)
}

// get all images by userId with page
func (c ApiFile) GetImages(albumId, key string, page int) revel.Result {
	imagesPage := fileService.ListImagesWithPage(c.getUserId(), albumId, key, page, 12)
	re := info.NewRe()
	re.Ok = true
	re.Item = imagesPage
	return c.RenderJson(re)
}

func (c ApiFile) UpdateImageTitle(fileId, title string) revel.Result {
	re := info.NewRe()
	re.Ok = fileService.UpdateImageTitle(c.getUserId(), fileId, title)
	return c.RenderJson(re)
}

func (c ApiFile) DeleteImage(fileId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = fileService.DeleteImage(c.getUserId(), fileId)
	return c.RenderJson(re)
}

*/

//-----------

// 输出image
// [OK]
func (c ApiFile) GetImage(fileId string) revel.Result {
	path := fileService.GetFile(c.getUserId(), fileId) // 得到路径
	if path == "" {
		return c.RenderText("")
	}
	fn := revel.BasePath + "/" + strings.TrimLeft(path, "/")
	file, _ := os.Open(fn)
	return c.RenderFile(file, revel.Inline) // revel.Attachment
}

// 下载附件
// [OK]
func (c ApiFile) GetAttach(fileId string) revel.Result {
	attach := attachService.GetAttach(fileId, c.getUserId()) // 得到路径
	path := attach.Path
	if path == "" {
		return c.RenderText("No Such File")
	}
	fn := revel.BasePath + "/" + strings.TrimLeft(path, "/")
	file, _ := os.Open(fn)
	return c.RenderBinary(file, attach.Title, revel.Attachment, time.Now()) // revel.Attachment
}

// 下载所有附件
// [OK]
func (c ApiFile) GetAllAttachs(noteId string) revel.Result {
	note := noteService.GetNoteById(noteId)
	if note.NoteId == "" {
		return c.RenderText("")
	}
	// 得到文件列表
	attachs := attachService.ListAttachs(noteId, c.getUserId())
	if attachs == nil || len(attachs) == 0 {
		return c.RenderText("")
	}

	/*
		dir := revel.BasePath + "/files/tmp"
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return c.RenderText("")
		}
	*/

	filename := note.Title + ".tar.gz"
	if note.Title == "" {
		filename = "all.tar.gz"
	}

	// file write
	fw, err := os.Create(revel.BasePath + "/files/" + filename)
	if err != nil {
		return c.RenderText("")
	}
	// defer fw.Close() // 不需要关闭, 还要读取给用户下载

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 遍历文件列表
	for _, attach := range attachs {
		fn := revel.BasePath + "/" + strings.TrimLeft(attach.Path, "/")
		fr, err := os.Open(fn)
		fileInfo, _ := fr.Stat()
		if err != nil {
			return c.RenderText("")
		}
		defer fr.Close()

		// 信息头
		h := new(tar.Header)
		h.Name = attach.Title
		h.Size = fileInfo.Size()
		h.Mode = int64(fileInfo.Mode())
		h.ModTime = fileInfo.ModTime()

		// 写信息头
		err = tw.WriteHeader(h)
		if err != nil {
			panic(err)
		}

		// 写文件
		_, err = io.Copy(tw, fr)
		if err != nil {
			panic(err)
		}
	} // for

	//    tw.Close()
	//    gw.Close()
	//    fw.Close()
	//    file, _ := os.Open(dir + "/" + filename)
	// fw.Seek(0, 0)
	return c.RenderBinary(fw, filename, revel.Attachment, time.Now()) // revel.Attachment
}
