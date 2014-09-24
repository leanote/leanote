package controllers

import (
	"github.com/revel/revel"
//	"encoding/json"
	. "github.com/leanote/leanote/app/lea"
//	"github.com/leanote/leanote/app/types"
//	"io/ioutil"
//	"math"
//	"os"
//	"path"
)

// lea++博客平台
type Lea struct {
	BaseController
}

// 进入某个用户的博客
func (c Lea) Index(tag, keywords string) revel.Result {
	c.RenderArgs["nav"] = "recommend"
	return c.p(tag, keywords, true)
}

func (c Lea) Latest(tag, keywords string) revel.Result {
	c.RenderArgs["nav"] = "latest"
	return c.p(tag, keywords, false);
}

func (c Lea) p(tag, keywords string, recommend bool) revel.Result {
	var tags = []string{}
	if recommend {
		tags = configService.GetGlobalArrayConfig("recommendTags")
	} else {
		tags = configService.GetGlobalArrayConfig("newTags")
	}
	// 如果不在所在的tag就不能搜索
	if !InArray(tags, tag) {
		tag = ""
	}
	c.RenderArgs["tag"] = tag
	
	page := c.GetPage()
	pageInfo, blogs := blogService.ListAllBlogs(tag, keywords, recommend, page, 10, "UpdatedTime", false)
	
	c.RenderArgs["pageInfo"] = pageInfo
	c.RenderArgs["blogs"] = blogs
	c.RenderArgs["tags"] = tags
	c.RenderArgs["keywords"] = keywords
	
	return c.RenderTemplate("lea/index.html");
}