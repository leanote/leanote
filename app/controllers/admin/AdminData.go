package admin

import (
	"archive/tar"
	"compress/gzip"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"io"
	"os"
	"time"
)

// 数据管理, 备份和恢复

type AdminData struct {
	AdminBaseController
}

func (c AdminData) Index() revel.Result {
	backups := configService.GetGlobalArrMapConfig("backups")
	// 逆序之
	backups2 := make([]map[string]string, len(backups))
	j := 0
	for i := len(backups) - 1; i >= 0; i-- {
		backups2[j] = backups[i]
		j++
	}
	c.RenderArgs["backups"] = backups2
	return c.RenderTemplate("admin/data/index.html")
}

func (c AdminData) Backup() revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = configService.Backup("")
	return c.RenderJson(re)
}

// 还原
func (c AdminData) Restore(createdTime string) revel.Result {
	re := info.Re{}
	re.Ok, re.Msg = configService.Restore(createdTime)
	return c.RenderJson(re)
}

func (c AdminData) Delete(createdTime string) revel.Result {
	re := info.Re{}
	re.Ok, re.Msg = configService.DeleteBackup(createdTime)
	return c.RenderJson(re)
}
func (c AdminData) UpdateRemark(createdTime, remark string) revel.Result {
	re := info.Re{}
	re.Ok, re.Msg = configService.UpdateBackupRemark(createdTime, remark)

	return c.RenderJson(re)
}
func (c AdminData) Download(createdTime string) revel.Result {
	backup, ok := configService.GetBackup(createdTime)
	if !ok {
		return c.RenderText("")
	}

	dbname, _ := revel.Config.String("db.dbname")
	path := backup["path"] + "/" + dbname
	allFiles := ListDir(path)

	filename := "backup_" + dbname + "_" + backup["createdTime"] + ".tar.gz"

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
	for _, file := range allFiles {
		fn := path + "/" + file
		fr, err := os.Open(fn)
		fileInfo, _ := fr.Stat()
		if err != nil {
			return c.RenderText("")
		}
		defer fr.Close()

		// 信息头
		h := new(tar.Header)
		h.Name = file
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

	return c.RenderBinary(fw, filename, revel.Attachment, time.Now()) // revel.Attachm
}
