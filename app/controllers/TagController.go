package controllers

import (
	"github.com/revel/revel"
	//	"encoding/json"
	//	"gopkg.in/mgo.v2/bson"
	//	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/info"
	//	"os/exec"
)

type Tag struct {
	BaseController
}

// 更新Tag
func (c Tag) UpdateTag(tag string) revel.Result {
	ret := info.NewRe()
	ret.Ok = true
	ret.Item = tagService.AddOrUpdateTag(c.GetUserId(), tag)
	return c.RenderJson(ret)
}

// 删除标签
func (c Tag) DeleteTag(tag string) revel.Result {
	ret := info.Re{}
	ret.Ok = true
	ret.Item = tagService.DeleteTag(c.GetUserId(), tag)
	return c.RenderJson(ret)
}
