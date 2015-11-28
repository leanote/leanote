package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"time"
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
		return false, "Leanote have been upgraded"
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
		data["UrlTitle"] = GetUrTitle(note.UserId.Hex(), note.Title, "note", noteId)
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
		data["UrlTitle"] = GetUrTitle(notebook.UserId.Hex(), notebook.Title, "notebook", notebookId)
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

// Usn设置
// 客户端 api

func (this *UpgradeService) moveTag() {
	usnI := 1
	tags := []info.Tag{}
	db.ListByQ(db.Tags, bson.M{}, &tags)
	for _, eachTag := range tags {
		tagTitles := eachTag.Tags
		now := time.Now()
		if tagTitles != nil && len(tagTitles) > 0 {
			for _, tagTitle := range tagTitles {
				noteTag := info.NoteTag{}
				noteTag.TagId = bson.NewObjectId()
				noteTag.Count = 1
				noteTag.Tag = tagTitle
				noteTag.UserId = eachTag.UserId
				noteTag.CreatedTime = now
				noteTag.UpdatedTime = now
				noteTag.Usn = usnI
				noteTag.IsDeleted = false
				db.Insert(db.NoteTags, noteTag)
				usnI++
			}
		}
	}
}

func (this *UpgradeService) setNotebookUsn() {
	usnI := 1
	notebooks := []info.Notebook{}
	db.ListByQWithFields(db.Notebooks, bson.M{}, []string{"UserId"}, &notebooks)

	for _, notebook := range notebooks {
		db.UpdateByQField(db.Notebooks, bson.M{"_id": notebook.NotebookId}, "Usn", usnI)
		usnI++
	}
}

func (this *UpgradeService) setNoteUsn() {
	usnI := 1
	notes := []info.Note{}
	db.ListByQWithFields(db.Notes, bson.M{}, []string{"UserId"}, &notes)

	for _, note := range notes {
		db.UpdateByQField(db.Notes, bson.M{"_id": note.NoteId}, "Usn", usnI)
		usnI++
	}
}

// 升级为Api, beta.4
func (this *UpgradeService) Api(userId string) (ok bool, msg string) {
	if configService.GetGlobalStringConfig("UpgradeBetaToBeta4") != "" {
		return false, "Leanote have been upgraded"
	}

	// user
	db.UpdateByQField(db.Users, bson.M{}, "Usn", 200000)

	// notebook
	db.UpdateByQField(db.Notebooks, bson.M{}, "IsDeleted", false)
	this.setNotebookUsn()

	// note
	// 1-N
	db.UpdateByQField(db.Notes, bson.M{}, "IsDeleted", false)
	this.setNoteUsn()

	// tag
	// 1-N
	/// tag, 要重新插入, 将之前的Tag表迁移到NoteTag中
	this.moveTag()

	configService.UpdateGlobalStringConfig(userId, "UpgradeBetaToBeta4", "1")

	return true, ""
}
