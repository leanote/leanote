package service

import (
	"github.com/leanote/leanote/app/info"
	"github.com/leanote/leanote/app/db"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type NoteService struct {
}

// 通过id, userId得到note
func (this *NoteService) GetNote(noteId, userId string) (note info.Note) {
	note = info.Note{}
	db.GetByIdAndUserId(db.Notes, noteId, userId, &note)
	return
}
// 得到blog, blogService用
// 不要传userId, 因为是公开的
func (this *NoteService) GetBlogNote(noteId string) (note info.Note) {
	note = info.Note{}
	db.GetByQ(db.Notes, bson.M{"_id": bson.ObjectIdHex(noteId), "IsBlog": true, "IsTrash": false}, &note)
	return
}
// 通过id, userId得到noteContent
func (this *NoteService) GetNoteContent(noteContentId, userId string) (noteContent info.NoteContent) {
	noteContent = info.NoteContent{}
	db.GetByIdAndUserId(db.NoteContents, noteContentId, userId, &noteContent)
	return
}

// 得到笔记和内容
func (this *NoteService) GetNoteAndContent(noteId, userId string) (noteAndContent info.NoteAndContent) {
	note := this.GetNote(noteId, userId)
	noteContent := this.GetNoteContent(noteId, userId)
	return info.NoteAndContent{note, noteContent}
}

// 列出note, 排序规则, 还有分页
// CreatedTime, UpdatedTime, title 来排序
func (this *NoteService) ListNotes(userId, notebookId string,
		isTrash bool, pageNumber, pageSize int, sortField string, isAsc bool, isBlog bool) (count int, notes []info.Note) {
	notes = []info.Note{}
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)
	// 不是trash的
	query := bson.M{"UserId": bson.ObjectIdHex(userId), "IsTrash": isTrash}
	if isBlog {
		query["IsBlog"] = true
	}
	if notebookId != "" {
		query["NotebookId"] = bson.ObjectIdHex(notebookId)
	}
	
	q := db.Notes.Find(query);
	
	// 总记录数
	count, _ = q.Count()
		
	q.Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&notes)
	return
}

// 通过noteIds来查询
// ShareService调用
func (this *NoteService) ListNotesByNoteIdsWithPageSort(noteIds []bson.ObjectId, userId string, 
		pageNumber, pageSize int, sortField string, isAsc bool, isBlog bool) (notes []info.Note) {
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)
	notes = []info.Note{}
	
	// 不是trash
	db.Notes.
		Find(bson.M{"_id": bson.M{"$in": noteIds}, "IsTrash": false}).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&notes)
	return
}
// shareService调用
func (this *NoteService) ListNotesByNoteIds(noteIds []bson.ObjectId) (notes []info.Note) {
	notes = []info.Note{}
	
	db.Notes.
		Find(bson.M{"_id": bson.M{"$in": noteIds}}).
		All(&notes)
	return
}
// blog需要
func (this *NoteService) ListNoteContentsByNoteIds(noteIds []bson.ObjectId) (notes []info.NoteContent) {
	notes = []info.NoteContent{}
	
	db.NoteContents.
		Find(bson.M{"_id": bson.M{"$in": noteIds}}).
		All(&notes)
	return
}
// 只得到abstract, 不需要content
func (this *NoteService) ListNoteAbstractsByNoteIds(noteIds []bson.ObjectId) (notes []info.NoteContent) {
	notes = []info.NoteContent{}
	db.ListByQWithFields(db.NoteContents, bson.M{"_id": bson.M{"$in": noteIds}}, []string{"_id", "Abstract"}, &notes)
	return
}

// 添加笔记
// 首先要判断Notebook是否是Blog, 是的话设为blog
// [ok]
func (this *NoteService) AddNote(note info.Note) info.Note {
	if(note.NoteId.Hex() == "") {
		noteId := bson.NewObjectId();
		note.NoteId = noteId;
	}
	note.CreatedTime = time.Now()
	note.UpdatedTime = note.CreatedTime
	note.IsTrash = false
	note.UpdatedUserId = note.UserId
	
	// 设为blog
	note.IsBlog = notebookService.IsBlog(note.NotebookId.Hex())
	
	db.Insert(db.Notes, note)
	
	// tag1
	tagService.AddTags(note.UserId.Hex(), note.Tags)
	
	return note
}

// 添加共享d笔记
func (this *NoteService) AddSharedNote(note info.Note, myUserId bson.ObjectId) info.Note {
	// 判断我是否有权限添加
	if shareService.HasUpdateNotebookPerm(note.UserId.Hex(), myUserId.Hex(), note.NotebookId.Hex()) {
		note.CreatedUserId = myUserId // 是我给共享我的人创建的
		return this.AddNote(note)
	}
	return info.Note{}
}


// 添加笔记本内容
// [ok]
func (this *NoteService) AddNoteContent(noteContent info.NoteContent) info.NoteContent {
	noteContent.CreatedTime = time.Now()
	noteContent.UpdatedTime = noteContent.CreatedTime 
	noteContent.UpdatedUserId = noteContent.UserId
	db.Insert(db.NoteContents, noteContent)
	return noteContent;
}

// 添加笔记和内容
// 这里使用 info.NoteAndContent 接收?
func (this *NoteService) AddNoteAndContent(note info.Note, noteContent info.NoteContent, myUserId bson.ObjectId) info.Note {
	if(note.NoteId.Hex() == "") {
		noteId := bson.NewObjectId()
		note.NoteId = noteId
	}
	noteContent.NoteId = note.NoteId
	if note.UserId != myUserId	{
		note = this.AddSharedNote(note, myUserId)
	} else {
		note = this.AddNote(note)
	}
	if note.NoteId != "" {
		this.AddNoteContent(noteContent)
	}
	return note
}

// 修改笔记
// [ok] TODO perm还没测
func (this *NoteService) UpdateNote(userId, updatedUserId, noteId string, needUpdate bson.M) bool {
	// updatedUserId 要修改userId的note, 此时需要判断是否有修改权限
	if userId != updatedUserId {
		if !shareService.HasUpdatePerm(userId, updatedUserId, noteId) {
			Log("NO AUTH")
			return false
		} else {
			Log("HAS AUTH -----------")
		}
	}
	
	needUpdate["UpdatedUserId"] = bson.ObjectIdHex(updatedUserId);
	needUpdate["UpdatedTime"] = time.Now();
	
	// 添加tag2
	if tags, ok := needUpdate["Tags"]; ok {
		tagService.AddTagsI(userId, tags)
	}
	
	// 是否修改了isBlog
	// 也要修改noteContents的IsBlog
	if isBlog, ok := needUpdate["IsBlog"]; ok {
		db.UpdateByIdAndUserIdMap(db.NoteContents, noteId, userId, bson.M{"IsBlog": isBlog})
	}
	
	return db.UpdateByIdAndUserIdMap(db.Notes, noteId, userId, needUpdate)
}

// 这里要判断权限, 如果userId != updatedUserId, 那么需要判断权限
// [ok] TODO perm还没测 [del]
func (this *NoteService) UpdateNoteTitle(userId, updatedUserId, noteId, title string) bool {
	// updatedUserId 要修改userId的note, 此时需要判断是否有修改权限
	if userId != updatedUserId {
		if !shareService.HasUpdatePerm(userId, updatedUserId, noteId) {
			println("NO AUTH")
			return false
		}
	}

	return db.UpdateByIdAndUserIdMap(db.Notes, noteId, userId, 
		bson.M{"UpdatedUserId": bson.ObjectIdHex(updatedUserId), "Title": title, "UpdatedTime": time.Now()})
}

// 修改笔记本内容
// [ok] TODO perm未测
func (this *NoteService) UpdateNoteContent(userId, updatedUserId, noteId, content, abstract string) bool {
	// updatedUserId 要修改userId的note, 此时需要判断是否有修改权限
	if userId != updatedUserId {
		if !shareService.HasUpdatePerm(userId, updatedUserId, noteId) {
			Log("NO AUTH")
			return false
		}
	}
	
	if db.UpdateByIdAndUserIdMap(db.NoteContents, noteId, userId, 
		bson.M{"UpdatedUserId": bson.ObjectIdHex(updatedUserId), 
		"Content": content, 
		"Abstract": abstract, 
		"UpdatedTime": time.Now()}) {
		
		// 这里, 添加历史记录
		noteContentHistoryService.AddHistory(noteId, userId, info.EachHistory{UpdatedUserId: bson.ObjectIdHex(updatedUserId), 
			Content: content,
			UpdatedTime: time.Now(),
		})
		return true
	}
	return false
}

// 更新tags
// [ok] [del]
func (this *NoteService) UpdateTags(noteId string, userId string, tags []string) bool {
	return db.UpdateByIdAndUserIdField(db.Notes, noteId, userId, "Tags", tags)
}

// 移动note
// trash, 正常的都可以用
// 1. 要检查下notebookId是否是自己的
// 2. 要判断之前是否是blog, 如果不是, 那么notebook是否是blog?
func (this *NoteService) MoveNote(noteId, notebookId, userId string) info.Note {
	if notebookService.IsMyNotebook(notebookId, userId) {
		re := db.UpdateByIdAndUserId(db.Notes, noteId, userId, 
			bson.M{"$set": bson.M{"IsTrash": false, 
				"NotebookId": bson.ObjectIdHex(notebookId)}})
				
		if re {
			// 更新blog状态
			this.updateToNotebookBlog(noteId, notebookId, userId)
		}
		
		return this.GetNote(noteId, userId);
	}
	return info.Note{}
}

// 如果自己的blog状态是true, 不用改变, 
// 否则, 如果notebookId的blog是true, 则改为true之
// 返回blog状态
func (this *NoteService) updateToNotebookBlog(noteId, notebookId, userId string) bool {
	if this.IsBlog(noteId) {
		return true
	}
	if notebookService.IsBlog(notebookId) {
		db.UpdateByIdAndUserId(db.Notes, noteId, userId, 
			bson.M{"$set": bson.M{"IsBlog": true}})
		return true
	}
	return false
}
// 判断是否是blog
func (this *NoteService) IsBlog(noteId string) bool {
	note := info.Note{}
	db.GetByQWithFields(db.Notes, bson.M{"_id": bson.ObjectIdHex(noteId)}, []string{"IsBlog"}, &note);
	return note.IsBlog
}

// 复制note
// 正常的可以用
// 先查, 再新建
// 要检查下notebookId是否是自己的
func (this *NoteService) CopyNote(noteId, notebookId, userId string) info.Note {
	if notebookService.IsMyNotebook(notebookId, userId) {
		note := this.GetNote(noteId, userId)
		noteContent := this.GetNoteContent(noteId, userId)
		
		// 重新生成noteId
		note.NoteId = bson.NewObjectId();
		note.NotebookId = bson.ObjectIdHex(notebookId)
		
		noteContent.NoteId = note.NoteId
		this.AddNoteAndContent(note, noteContent, note.UserId);
		
		// 更新blog状态
		isBlog := this.updateToNotebookBlog(note.NoteId.Hex(), notebookId, userId)
		
		note.IsBlog = isBlog
		return note
	}
	
	return info.Note{}
}

// 复制别人的共享笔记给我
// TODO 判断是否共享了给我
func (this *NoteService) CopySharedNote(noteId, notebookId, fromUserId, myUserId string) info.Note {
	if notebookService.IsMyNotebook(notebookId, myUserId) {
		note := this.GetNote(noteId, fromUserId)
		if note.NoteId == "" {
			return info.Note{}
		}
		noteContent := this.GetNoteContent(noteId, fromUserId)
		
		// 重新生成noteId
		note.NoteId = bson.NewObjectId();
		note.NotebookId = bson.ObjectIdHex(notebookId)
		note.UserId = bson.ObjectIdHex(myUserId)
		note.IsTop = false
		note.IsBlog = false // 别人的可能是blog
		
		// content
		noteContent.NoteId = note.NoteId
		noteContent.UserId = note.UserId
		
		// 添加之
		note = this.AddNoteAndContent(note, noteContent, note.UserId);
		
		// 更新blog状态
		isBlog := this.updateToNotebookBlog(note.NoteId.Hex(), notebookId, myUserId)
		
		note.IsBlog = isBlog
		return note
	}
	
	return info.Note{}
}

// 通过noteId得到notebookId
// shareService call
// [ok]
func (this *NoteService) GetNotebookId(noteId string) bson.ObjectId {
	note := &info.Note{}
	db.Get(db.Notes, noteId, note)
	return note.NotebookId
}

//------------------
// 搜索Note
func (this *NoteService) SearchNote(key, userId string, pageNumber, pageSize int, sortField string, isAsc, isBlog bool) (count int, notes []info.Note) {
	notes = []info.Note{}
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)
	
	// 不是trash的
	query := bson.M{"UserId": bson.ObjectIdHex(userId), 
		"IsTrash": false, 
		"Title": bson.M{"$regex": bson.RegEx{".*?" + key + ".*", "i"}},
	}
	if isBlog {
		query["IsBlog"] = true
	}
	q := db.Notes.Find(query);
	
	// 总记录数
	count, _ = q.Count()
	
	q.Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&notes)
		
	// 如果 < pageSize 那么搜索content, 且id不在这些id之间的
	if len(notes) < pageSize {
		notes = this.searchNoteFromContent(notes, userId, key, pageSize, sortFieldR, isBlog)
	}
	return
}

// 搜索noteContents, 补集pageSize个
func (this *NoteService) searchNoteFromContent(notes []info.Note, userId, key string, pageSize int, sortField string, isBlog bool) []info.Note {
	var remain = pageSize - len(notes)
	noteIds := make([]bson.ObjectId, len(notes))
	for i, note := range notes {
		noteIds[i] = note.NoteId
	}
	LogJ(noteIds)
	noteContents := []info.NoteContent{}
	query := bson.M{"_id": bson.M{"$nin": noteIds}, "UserId": bson.ObjectIdHex(userId), "Content": bson.M{"$regex": bson.RegEx{".*?" + key + ".*", "i"}}}
	if isBlog {
		query["IsBlog"] = true
	}
	db.NoteContents.
		Find(query).
		Sort(sortField).
		Limit(remain).
		Select(bson.M{"_id": true}).
		All(&noteContents)
	var lenContent = len(noteContents)
	if(lenContent == 0) {
		return notes
	}
	
	// 收集ids
	noteIds2 := make([]bson.ObjectId, lenContent)
	for i, content := range noteContents {
		noteIds2[i] = content.NoteId
	}
	
//	Log(" content search ")
//	Log(lenContent)
	
	// 得到notes
	notes2 := this.ListNotesByNoteIds(noteIds2)
	
	// 合并之
	notes = append(notes, notes2...)
	return notes
}

//----------------
// tag搜索
func (this *NoteService) SearchNoteByTags(tags []string, userId string, pageNumber, pageSize int, sortField string, isAsc bool) (count int, notes []info.Note) {
	notes = []info.Note{}
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)
	
	// 不是trash的
	query := bson.M{"UserId": bson.ObjectIdHex(userId), 
		"IsTrash": false, 
		"Tags": bson.M{"$all": tags}}
	
	q := db.Notes.Find(query);
	
	// 总记录数
	count, _ = q.Count()
	
	Log(count)
		
	q.Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&notes)
	return
}