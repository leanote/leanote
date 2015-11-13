package api

import (
	"github.com/revel/revel"
	//	"encoding/json"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"time"
	//	"github.com/leanote/leanote/app/types"
	"io/ioutil"
	//	"fmt"
	//	"math"
	"os"

	//	"path"
	//	"strconv"
)

type ApiUser struct {
	ApiBaseContrller
}

// 获取用户信息
// [OK]
func (c ApiUser) Info() revel.Result {
	re := info.NewApiRe()

	userInfo := c.getUserInfo()
	if userInfo.UserId == "" {
		return c.RenderJson(re)
	}
	apiUser := info.ApiUser{
		UserId:   userInfo.UserId.Hex(),
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Logo:     userInfo.Logo,
		Verified: userInfo.Verified,
	}
	return c.RenderJson(apiUser)
}

// 修改用户名
// [OK]
func (c ApiUser) UpdateUsername(username string) revel.Result {
	re := info.NewApiRe()
	if c.GetUsername() == "demo" {
		re.Msg = "cannotUpdateDemo"
		return c.RenderJson(re)
	}

	if re.Ok, re.Msg = Vd("username", username); !re.Ok {
		return c.RenderJson(re)
	}

	re.Ok, re.Msg = userService.UpdateUsername(c.getUserId(), username)
	return c.RenderJson(re)
}

// 修改密码
// [OK]
func (c ApiUser) UpdatePwd(oldPwd, pwd string) revel.Result {
	re := info.NewApiRe()
	if c.GetUsername() == "demo" {
		re.Msg = "cannotUpdateDemo"
		return c.RenderJson(re)
	}
	if re.Ok, re.Msg = Vd("password", oldPwd); !re.Ok {
		return c.RenderJson(re)
	}
	if re.Ok, re.Msg = Vd("password", pwd); !re.Ok {
		return c.RenderJson(re)
	}
	re.Ok, re.Msg = userService.UpdatePwd(c.getUserId(), oldPwd, pwd)
	return c.RenderJson(re)
}

// 获得同步状态
// [OK]
func (c ApiUser) GetSyncState() revel.Result {
	ret := bson.M{"LastSyncUsn": userService.GetUsn(c.getUserId()), "LastSyncTime": time.Now().Unix()}
	return c.RenderJson(ret)
}

// 头像设置
// 参数file=文件
// 成功返回{Logo: url} 头像新url
// [OK]
func (c ApiUser) UpdateLogo() revel.Result {
	ok, msg, url := c.uploadImage()

	if ok {
		ok = userService.UpdateAvatar(c.getUserId(), url)
		return c.RenderJson(map[string]string{"Logo": url})
	} else {
		re := info.NewApiRe()
		re.Msg = msg
		return c.RenderJson(re)
	}
}

// 上传图片
func (c ApiUser) uploadImage() (ok bool, msg, url string) {
	var fileUrlPath = ""
	ok = false

	file, handel, err := c.Request.FormFile("file")
	if err != nil {
		return
	}
	defer file.Close()
	// 生成上传路径
	fileUrlPath = "public/upload/" + c.getUserId() + "/images/logo"

	dir := revel.BasePath + "/" + fileUrlPath
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}
	// 生成新的文件名
	filename := handel.Filename

	var ext string

	_, ext = SplitFilename(filename)
	if ext != ".gif" && ext != ".jpg" && ext != ".png" && ext != ".bmp" && ext != ".jpeg" {
		msg = "notImage"
		return
	}

	filename = NewGuid() + ext
	data, err := ioutil.ReadAll(file)
	if err != nil {
		LogJ(err)
		return
	}

	// > 5M?
	if len(data) > 5*1024*1024 {
		msg = "fileIsTooLarge"
		return
	}

	toPath := dir + "/" + filename
	err = ioutil.WriteFile(toPath, data, 0777)
	if err != nil {
		LogJ(err)
		return
	}

	ok = true
	url = configService.GetSiteUrl() + "/" + fileUrlPath + "/" + filename
	return
}
