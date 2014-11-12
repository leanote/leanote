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
/*
<li>Migrate "About me" to single(a single post)</li>
<li>Add some default themes to administrator</li>
<li>Generate "UrlTitle" for all notes. "UrlTitle" is a friendly url for post</li>
<li>Generate "UrlTitle" for all notebooks</li>
<li>Generate "UrlTitle" for all singles</li>
*/
func (this *UpgradeService) UpgradeBetaToBeta2(userId string) (ok bool, msg string) {
	if configService.GetGlobalStringConfig("UpgradeBetaToBeta2") != "" {
		return false, "已升级"
	}
	
	// 1. aboutMe -> page
	userBlogs := []info.UserBlog{}
	db.ListByQ(db.UserBlogs, bson.M{}, &userBlogs)
	
	for _, userBlog := range userBlogs {
		blogService.AddOrUpdateSingle(userBlog.UserId.Hex(), "", "About Me", userBlog.AboutMe)
	}
	
	
	// 2. 默认主题, 给admin用户
	themeService.UpgradeThemeBeta2()
	
	// 3. UrlTitles
	
	// 3.1 note
	notes := []info.Note{}
	db.ListByQ(db.Notes, bson.M{}, &notes)
	for _, note := range notes {
		data := bson.M{}	
		noteId := note.NoteId.Hex()
		// PublicTime, RecommendTime = UpdatedTime
		if note.IsBlog && note.PublicTime.Year() < 100 {
			data["PublicTime"] = note.UpdatedTime
			data["RecommendTime"] = note.UpdatedTime
			Log("Time " + noteId)
		}
		data["UrlTitle"] = GetUrTitle(note.UserId.Hex(), note.Title, "note")
		db.UpdateByIdAndUserIdMap2(db.Notes, note.NoteId, note.UserId, data)
		Log(noteId)
	}
	
	// 3.2
	Log("notebook")
	notebooks := []info.Notebook{}
	db.ListByQ(db.Notebooks, bson.M{}, &notebooks)
	for _, notebook := range notebooks {
		notebookId := notebook.NotebookId.Hex()
		data := bson.M{}
		data["UrlTitle"] = GetUrTitle(notebook.UserId.Hex(), notebook.Title, "notebook")
		db.UpdateByIdAndUserIdMap2(db.Notebooks, notebook.NotebookId, notebook.UserId, data)
		Log(notebookId)
	}
	
	// 3.3 single
	/*
	singles := []info.BlogSingle{}
	db.ListByQ(db.BlogSingles, bson.M{}, &singles)
	for _, single := range singles {
		singleId := single.SingleId.Hex()
		blogService.UpdateSingleUrlTitle(single.UserId.Hex(), singleId, single.Title)
		Log(singleId)
	}
	*/
	
	// 删除索引
	db.ShareNotes.DropIndex("UserId", "ToUserId", "NoteId")
	
	ok = true
	msg = "success"
	configService.UpdateGlobalStringConfig(userId, "UpgradeBetaToBeta2", "1")
	
	return
}
