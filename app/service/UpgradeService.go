package service

import (
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/db"
	"gopkg.in/mgo.v2/bson"
//	"time"
)


type UpgradeService struct {
}

// 添加了PublicTime, RecommendTime
func (this *UpgradeService) UpgradeBlog() bool {
	notes := []info.Note{}
	db.ListByQ(db.Notes, bson.M{"IsBlog": true}, &notes)
	
	// PublicTime, RecommendTime = UpdatedTime
	for _, note := range notes {
		if note.IsBlog && note.PublicTime.Year() < 100 {
			db.UpdateByIdAndUserIdMap2(db.Notes, note.NoteId, note.UserId, bson.M{"PublicTime": note.UpdatedTime, "RecommendTime": note.UpdatedTime})
			Log(note.NoteId.Hex())
		}
	}
	
	return true
}

// 11-5自定义博客升级, 将aboutMe移至pages
func (this *UpgradeService) UpgradeBetaToSelfBlog(userId string) (ok bool, msg string) {
	if configService.GetGlobalStringConfig("UpgradeBetaToSelfBlog") != "" {
		return false, "已升级"
	}
	
	// 1. aboutMe -> page
	userBlogs := []info.UserBlog{}
	db.ListByQ(db.UserBlogs, bson.M{}, &userBlogs)
	
	for _, userBlog := range userBlogs {
		blogService.AddOrUpdateSingle(userBlog.UserId.Hex(), "", "About Me", userBlog.AboutMe)
	}
	
	configService.UpdateGlobalStringConfig(userId, "UpgradeBetaToSelfBlog", "ok")
	
	ok = true
	msg = "success"
	
	// 2. 默认主题, 给admin用户
	return
}
