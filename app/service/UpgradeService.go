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
		db.UpdateByIdAndUserIdMap2(db.Notes, note.NoteId, note.UserId, bson.M{"PublicTime": note.UpdatedTime, "RecommendTime": note.UpdatedTime})
		Log(note.NoteId.Hex())
	}
	
	return true
}