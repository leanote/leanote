package api

import (
	"github.com/leanote/leanote/app/info"
	"github.com/revel/revel"
	//	"gopkg.in/mgo.v2/bson"
	//	. "github.com/leanote/leanote/app/lea"
	//	"io/ioutil"
)

// 标签API

type ApiTag struct {
	ApiBaseContrller
}

// 获取同步的标签
// [OK]
// > afterUsn的笔记
// 返回 {ChunkHighUsn: 本下最大的usn, 借此可以知道是否还有, Notebooks: []}
func (c ApiTag) GetSyncTags(afterUsn, maxEntry int) revel.Result {
	if maxEntry == 0 {
		maxEntry = 100
	}
	tags := tagService.GeSyncTags(c.getUserId(), afterUsn, maxEntry)
	return c.RenderJson(tags)
}

// 添加Tag
// [OK]
// 不会产生冲突, 即使里面有
// 返回
/*
{
  "TagId": "551978dd99c37b9bc5000001",
  "UserId": "54a1676399c37b1c77000002",
  "Tag": "32",
  "Usn": 25,
  "Count": 1,
  "CreatedTime": "2015-03-31T00:25:01.149312407+08:00",
  "UpdatedTime": "2015-03-31T00:25:01.149312407+08:00",
  "IsDeleted": false
}
*/
func (c ApiTag) AddTag(tag string) revel.Result {
	ret := tagService.AddOrUpdateTag(c.getUserId(), tag)
	return c.RenderJson(ret)
}

// 删除标签
// [OK]
func (c ApiTag) DeleteTag(tag string, usn int) revel.Result {
	re := info.NewReUpdate()
	re.Ok, re.Msg, re.Usn = tagService.DeleteTagApi(c.getUserId(), tag, usn)
	return c.RenderJson(re)
}
