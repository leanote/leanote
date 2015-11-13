package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	//	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	//	"time"
)

// 历史记录
type NoteContentHistoryService struct {
}

// 每个历史记录最大值
var maxSize = 10

// 新建一个note, 不需要添加历史记录
// 添加历史
func (this *NoteContentHistoryService) AddHistory(noteId, userId string, eachHistory info.EachHistory) {
	// 检查是否是空
	if eachHistory.Content == "" {
		return
	}

	// 先查是否存在历史记录, 没有则添加之
	history := info.NoteContentHistory{}
	db.GetByIdAndUserId(db.NoteContentHistories, noteId, userId, &history)
	if history.NoteId == "" {
		this.newHistory(noteId, userId, eachHistory)
	} else {
		// 判断是否超出 maxSize, 如果超出则pop最后一个, 再push之, 不用那么麻烦, 直接update吧, 虽然影响性能
		// TODO
		l := len(history.Histories)
		if l >= maxSize {
			// history.Histories = history.Histories[l-maxSize:] // BUG, 致使都是以前的
			history.Histories = history.Histories[:maxSize]
		}
		newHistory := []info.EachHistory{eachHistory}
		newHistory = append(newHistory, history.Histories...) // 在开头加了, 最近的在最前
		history.Histories = newHistory

		// 更新之
		db.UpdateByIdAndUserId(db.NoteContentHistories, noteId, userId, history)
	}
	return
}

// 新建历史
func (this *NoteContentHistoryService) newHistory(noteId, userId string, eachHistory info.EachHistory) {
	history := info.NoteContentHistory{NoteId: bson.ObjectIdHex(noteId),
		UserId:    bson.ObjectIdHex(userId),
		Histories: []info.EachHistory{eachHistory},
	}

	// 保存之
	db.Insert(db.NoteContentHistories, history)
}

// 列表展示
func (this *NoteContentHistoryService) ListHistories(noteId, userId string) []info.EachHistory {
	histories := info.NoteContentHistory{}
	db.GetByIdAndUserId(db.NoteContentHistories, noteId, userId, &histories)
	return histories.Histories
}
