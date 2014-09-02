package service

import (
	"github.com/leanote/leanote/app/info"
	"github.com/leanote/leanote/app/db"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
//	"time"
)

/*
每添加,更新note时, 都将tag添加到tags表中
*/
type TagService struct {
}

func (this *TagService) GetTags(userId string) []string {
	tag := info.Tag{}
	db.Get(db.Tags, userId, &tag)
	LogJ(tag)
	return tag.Tags
}

func (this *TagService) AddTagsI(userId string, tags interface{}) bool {
	if ts, ok2 := tags.([]string); ok2 {
		return this.AddTags(userId, ts)
	}
	return false
}
func (this *TagService) AddTags(userId string, tags []string) bool {
	for _, tag := range tags {
		if !db.Upsert(db.Tags, 
			bson.M{"_id": bson.ObjectIdHex(userId)}, 
			bson.M{"$addToSet": bson.M{"Tags": tag}}) {
			return false
		}
	}
	return true
}