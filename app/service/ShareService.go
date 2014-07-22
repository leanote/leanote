package service

import (
	"leanote/app/db"
	"leanote/app/info"
	//	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"time"
)

// 共享Notebook, Note服务
type ShareService struct {
}

//-----------------
// 总体来说, 这个方法比较麻烦, 速度未知. 以后按以下方案来缓存用户基础数据

// 以后建个用户的基本数据表, 放所有notebook, sharedNotebook的缓存!!
// 每更新一次则启动一个goroutine异步更新
// 共享只能共享本notebook下的, 如果其子也要共享, 必须设置其子!!!
// 那么, 父, 子在shareNotebooks表中都会有记录

// 得到用户的所有*被*共享的Notebook
// 1 得到别人共享给我的所有notebooks
// 2 按parent进行层次化
// 3 每个层按seq进行排序
// 4 按用户分组
// [ok]
func (this *ShareService) GetShareNotebooks(userId string) (info.ShareNotebooksByUser, []info.User) {
	//-------------
	// 查询HasShareNote表得到所有其它用户信息
	hasShareNotes := []info.HasShareNote{}
	db.ListByQ(db.HasShareNotes, bson.M{"ToUserId": bson.ObjectIdHex(userId)}, &hasShareNotes)

	userIds := make([]bson.ObjectId, len(hasShareNotes))
	for i, each := range hasShareNotes {
		userIds[i] = each.UserId
	}
	userInfos := userService.GetUserInfosOrderBySeq(userIds)

	//--------------------
	// 得到他们共享给我的notebooks

	shareNotebooks := []info.ShareNotebook{}
	db.ShareNotebooks.Find(bson.M{"ToUserId": bson.ObjectIdHex(userId)}).All(&shareNotebooks)

	if len(shareNotebooks) == 0 {
		return nil, userInfos
	}

	shareNotebooksLen := len(shareNotebooks)

	// 找到了所有的notbookId, 那么找notebook表得到其详细信息
	notebookIds := make([]bson.ObjectId, shareNotebooksLen)
	shareNotebooksMap := make(map[bson.ObjectId]info.ShareNotebook, shareNotebooksLen)
	for i, each := range shareNotebooks {
		// 默认的是没有notebookId的
		notebookIds[i] = each.NotebookId
		shareNotebooksMap[each.NotebookId] = each
	}

	// 1, 2
	subNotebooks := notebookService.GetNotebooksByNotebookIds(notebookIds)
	// 填充其它信息变成SubShareNotebooks
	subShareNotebooks := this.parseToSubShareNotebooks(&subNotebooks, &shareNotebooksMap)

	// 3 按用户进行分组成ShareNotebooksByUser
	shareNotebooksByUsersMap := map[bson.ObjectId][]info.ShareNotebooks{}
	// 先建立userId => []
	for _, eachSub := range subShareNotebooks {
		userId := eachSub.Notebook.UserId
		if _, ok := shareNotebooksByUsersMap[userId]; ok {
			shareNotebooksByUsersMap[userId] = append(shareNotebooksByUsersMap[userId], eachSub)
		} else {
			shareNotebooksByUsersMap[userId] = []info.ShareNotebooks{eachSub}
		}
	}
	shareNotebooksByUser := info.ShareNotebooksByUser{}
	for userId, eachShareNotebooks := range shareNotebooksByUsersMap {
		// 4, 按用户排序
		shareNotebooksByUser[userId.Hex()] = sortSubShareNotebooks(eachShareNotebooks)
	}

	return shareNotebooksByUser, userInfos
}

// 排序
func sortSubShareNotebooks(eachNotebooks info.SubShareNotebooks) info.SubShareNotebooks {
	// 遍历子, 则子往上进行排序
	for _, eachNotebook := range eachNotebooks {
		if eachNotebook.Subs != nil && len(eachNotebook.Subs) > 0 {
			eachNotebook.Subs = sortSubShareNotebooks(eachNotebook.Subs)
		}
	}

	// 子排完了, 本层排
	sort.Sort(&eachNotebooks)
	return eachNotebooks
}

// 将普通的notebooks添加perm及shareNotebook信息
func (this *ShareService) parseToSubShareNotebooks(subNotebooks *info.SubNotebooks, shareNotebooksMap *map[bson.ObjectId]info.ShareNotebook) info.SubShareNotebooks {
	subShareNotebooks := info.SubShareNotebooks{}
	for _, each := range *subNotebooks {
		shareNotebooks := info.ShareNotebooks{}
		shareNotebooks.Notebook = each.Notebook                              // 基本信息有了
		shareNotebooks.ShareNotebook = (*shareNotebooksMap)[each.NotebookId] // perm有了

		// 公用的, 单独赋值
		shareNotebooks.Seq = shareNotebooks.ShareNotebook.Seq
		shareNotebooks.NotebookId = shareNotebooks.ShareNotebook.NotebookId

		// 还有其子, 递归解析之
		if each.Subs != nil && len(each.Subs) > 0 {
			shareNotebooks.Subs = this.parseToSubShareNotebooks(&each.Subs, shareNotebooksMap)
		}
		subShareNotebooks = append(subShareNotebooks, shareNotebooks)
	}

	return subShareNotebooks
}

//-------------

// 得到共享笔记本下的notes
func (this *ShareService) ListShareNotesByNotebookId(notebookId, myUserId, sharedUserId string,
	page, pageSize int, sortField string, isAsc bool) []info.ShareNoteWithPerm {
	// 1 首先判断是否真的sharedUserId 共享了 notebookId 给 myUserId
	shareNotebook := info.ShareNotebook{}
	db.GetByQ(db.ShareNotebooks, bson.M{"NotebookId": bson.ObjectIdHex(notebookId),
		"UserId": bson.ObjectIdHex(sharedUserId), "ToUserId": bson.ObjectIdHex(myUserId)}, &shareNotebook)

	if shareNotebook.NotebookId == "" {
		return nil
	}

	perm := shareNotebook.Perm

	// 2 得到该notebook下分页显示所有的notes
	_, notes := noteService.ListNotes(sharedUserId, notebookId, false, page, pageSize, sortField, isAsc, false)

	// 3 添加权限信息
	// 3.1 如果该notebook自己有其它权限信息, 比如1, 那么覆盖notebook的权限信息
	noteIds := make([]bson.ObjectId, len(notes))
	for i, note := range notes {
		noteIds[i] = note.NoteId
	}
	notePerms := this.getNotesPerm(noteIds, myUserId, sharedUserId)

	// 3.2 组合
	notesWithPerm := make([]info.ShareNoteWithPerm, len(notes))
	for i, each := range notes {
		thisPerm := perm
		if selfPerm, ok := notePerms[each.NoteId]; ok {
			thisPerm = selfPerm
		}

		notesWithPerm[i] = info.ShareNoteWithPerm{each, thisPerm}
	}
	return notesWithPerm
}

// 得到note的perm信息
func (this *ShareService) getNotesPerm(noteIds []bson.ObjectId, myUserId, sharedUserId string) map[bson.ObjectId]int {
	shareNotes := []info.ShareNote{}
	db.ListByQ(db.ShareNotes, bson.M{"NoteId": bson.M{"$in": noteIds}, "UserId": bson.ObjectIdHex(sharedUserId), "ToUserId": bson.ObjectIdHex(myUserId)}, &shareNotes)

	notesPerm := make(map[bson.ObjectId]int, len(shareNotes))
	for _, each := range shareNotes {
		notesPerm[each.NoteId] = each.Perm
	}

	return notesPerm
}

// 得到默认的单个的notes 共享集
// 如果真要支持排序, 这里得到所有共享的notes, 到noteService方再sort和limit
// 可以这样! 到时将零散的共享noteId放在用户基本数据中
// 这里不好排序
func (this *ShareService) ListShareNotes(myUserId, sharedUserId string,
	pageNumber, pageSize int, sortField string, isAsc bool) []info.ShareNoteWithPerm {

	skipNum, _ := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)
	shareNotes := []info.ShareNote{}

	db.ShareNotes.
		Find(bson.M{"UserId": bson.ObjectIdHex(sharedUserId), "ToUserId": bson.ObjectIdHex(myUserId)}).
		//		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&shareNotes)

	if len(shareNotes) == 0 {
		return nil
	}

	noteIds := make([]bson.ObjectId, len(shareNotes))
	for i, each := range shareNotes {
		noteIds[i] = each.NoteId
	}
	notes := noteService.ListNotesByNoteIds(noteIds)
	notesMap := make(map[bson.ObjectId]info.Note, len(notes))
	for _, each := range notes {
		notesMap[each.NoteId] = each
	}

	// 将shareNotes与notes结合起来
	notesWithPerm := make([]info.ShareNoteWithPerm, len(shareNotes))
	for i, each := range shareNotes {
		notesWithPerm[i] = info.ShareNoteWithPerm{notesMap[each.NoteId], each.Perm}
	}
	return notesWithPerm
}

func (this *ShareService) notes2NotesWithPerm(notes []info.Note) {

}

// 添加一个notebook共享
// [ok]
func (this *ShareService) AddShareNotebook1(shareNotebook info.ShareNotebook) bool {
	// 添加一条记录说明两者存在关系
	this.AddHasShareNote(shareNotebook.UserId.Hex(), shareNotebook.ToUserId.Hex())

	shareNotebook.CreatedTime = time.Now()
	return db.Insert(db.ShareNotebooks, shareNotebook)
}

// 添加共享笔记本
func (this *ShareService) AddShareNotebook(notebookId string, perm int, userId, email string) (bool, string, string) {
	// 通过email得到被共享的userId
	toUserId := userService.GetUserId(email)
	if toUserId == "" {
		return false, "无该用户", ""
	}

	// 添加一条记录说明两者存在关系
	this.AddHasShareNote(userId, toUserId)

	// 先删除之
	db.Delete(db.ShareNotebooks, bson.M{"NotebookId": bson.ObjectIdHex(notebookId),
		"UserId":   bson.ObjectIdHex(userId),
		"ToUserId": bson.ObjectIdHex(toUserId),
	})

	shareNotebook := info.ShareNotebook{NotebookId: bson.ObjectIdHex(notebookId),
		UserId:      bson.ObjectIdHex(userId),
		ToUserId:    bson.ObjectIdHex(toUserId),
		Perm:        perm,
		CreatedTime: time.Now(),
	}
	return db.Insert(db.ShareNotebooks, shareNotebook), "", toUserId
}

// 添加一个note共享
// [ok]
/*
func (this *ShareService) AddShareNote(shareNote info.ShareNote) bool {
	shareNote.CreatedTime = time.Now()
	return db.Insert(db.ShareNotes, shareNote)
}
*/
func (this *ShareService) AddShareNote(noteId string, perm int, userId, email string) (bool, string, string) {
	// 通过email得到被共享的userId
	toUserId := userService.GetUserId(email)
	if toUserId == "" {
		return false, "无该用户", ""
	}

	// 添加一条记录说明两者存在关系
	this.AddHasShareNote(userId, toUserId)

	// 先删除之
	db.Delete(db.ShareNotes, bson.M{"NoteId": bson.ObjectIdHex(noteId),
		"UserId":   bson.ObjectIdHex(userId),
		"ToUserId": bson.ObjectIdHex(toUserId),
	})

	shareNote := info.ShareNote{NoteId: bson.ObjectIdHex(noteId),
		UserId:      bson.ObjectIdHex(userId),
		ToUserId:    bson.ObjectIdHex(toUserId),
		Perm:        perm,
		CreatedTime: time.Now(),
	}
	return db.Insert(db.ShareNotes, shareNote), "", toUserId
}

// updatedUserId是否有修改userId noteId的权限?
func (this *ShareService) HasUpdatePerm(userId, updatedUserId, noteId string) bool {
	// 1. noteId是否被共享了?
	// 得到该note share的信息
	/*
		UserId                 bson.ObjectId `bson:"UserId"`
		ToUserId               bson.ObjectId `bson:"ToUserId"`
		NoteId                 bson.ObjectId `bson:"NoteId"`
		Perm                   int           `bson:"Perm"` // 权限, 0只读, 1可修改
	*/
	if !db.Has(db.ShareNotes,
		bson.M{"UserId": bson.ObjectIdHex(userId), "ToUserId": bson.ObjectIdHex(updatedUserId), "NoteId": bson.ObjectIdHex(noteId), "Perm": 1}) {
		// noteId的notebookId是否被共享了?
		notebookId := noteService.GetNotebookId(noteId)
		if notebookId.Hex() == "" {
			return false
		}

		// 判断notebook是否被共享
		if !db.Has(db.ShareNotebooks,
			bson.M{"UserId": bson.ObjectIdHex(userId), "ToUserId": bson.ObjectIdHex(updatedUserId), "NotebookId": notebookId, "Perm": 1}) {
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

// updatedUserId是否有修改userId notebookId的权限?
func (this *ShareService) HasUpdateNotebookPerm(userId, updatedUserId, notebookId string) bool {
	// 判断notebook是否被共享
	if !db.Has(db.ShareNotebooks,
		bson.M{"UserId": bson.ObjectIdHex(userId), "ToUserId": bson.ObjectIdHex(updatedUserId), "NotebookId": bson.ObjectIdHex(notebookId), "Perm": 1}) {
		return false
	} else {
		return true
	}
}

// 共享note, notebook时使用
func (this *ShareService) AddHasShareNote(userId, toUserId string) bool {
	db.Insert(db.HasShareNotes, info.HasShareNote{UserId: bson.ObjectIdHex(userId), ToUserId: bson.ObjectIdHex(toUserId)})
	return true
}

// userId是否被共享了noteId
func (this *ShareService) hasSharedNote(noteId, myUserId string) bool {
	return db.Has(db.ShareNotes, bson.M{"ToUserId": bson.ObjectIdHex(myUserId), "NoteId": bson.ObjectIdHex(noteId)})
}

// noteId的notebook是否共享了给我
func (this *ShareService) hasSharedNotebook(noteId, myUserId, sharedUserId string) bool {
	note := noteService.GetNote(noteId, sharedUserId)
	if note.NoteId != "" {
		return db.Has(db.ShareNotebooks, bson.M{"NotebookId": note.NotebookId,
			"UserId":   bson.ObjectIdHex(sharedUserId),
			"ToUserId": bson.ObjectIdHex(myUserId),
		})
	}
	return false
}

// 得到共享的笔记内容
// 首先要判断这个note是否我被共享了
func (this *ShareService) GetShareNoteContent(noteId, myUserId, sharedUserId string) (noteContent info.NoteContent) {
	noteContent = info.NoteContent{}
	// 是否单独共享了该notebook
	// 或者, 其notebook共享了我
	if this.hasSharedNote(noteId, myUserId) || this.hasSharedNotebook(noteId, myUserId, sharedUserId) {
		db.Get(db.NoteContents, noteId, &noteContent)
	} else {
	}
	return
}

// 查看note的分享信息
// 分享给了哪些用户和权限
// ShareNotes表 userId = me, noteId = ...
// 还要查看该note的notebookId分享的信息
func (this *ShareService) ListNoteShareUserInfo(noteId, userId string) []info.ShareUserInfo {
	// 得到shareNote信息, 得到所有的ToUserId
	shareNotes := []info.ShareNote{}
	db.ListByQLimit(db.ShareNotes, bson.M{"NoteId": bson.ObjectIdHex(noteId), "UserId": bson.ObjectIdHex(userId)}, &shareNotes, 100)

	if len(shareNotes) == 0 {
		return nil
	}

	shareNotesMap := make(map[bson.ObjectId]info.ShareNote, len(shareNotes))
	for _, each := range shareNotes {
		shareNotesMap[each.ToUserId] = each
	}

	toUserIds := make([]bson.ObjectId, len(shareNotes))
	for i, eachShareNote := range shareNotes {
		toUserIds[i] = eachShareNote.ToUserId
	}

	note := noteService.GetNote(noteId, userId)
	if note.NoteId == "" {
		return nil
	}

	// 查看其notebook的shareNotebooks信息
	shareNotebooks := []info.ShareNotebook{}
	db.ListByQ(db.ShareNotebooks,
		bson.M{"NotebookId": note.NotebookId, "UserId": bson.ObjectIdHex(userId), "ToUserId": bson.M{"$in": toUserIds}},
		&shareNotebooks)
	shareNotebooksMap := make(map[bson.ObjectId]info.ShareNotebook, len(shareNotebooks))
	for _, each := range shareNotebooks {
		shareNotebooksMap[each.ToUserId] = each
	}

	// 得到用户信息
	userInfos := userService.ListUserInfosByUserIds(toUserIds)

	if len(userInfos) == 0 {
		return nil
	}

	shareUserInfos := make([]info.ShareUserInfo, len(userInfos))

	for i, userInfo := range userInfos {
		_, hasNotebook := shareNotebooksMap[userInfo.UserId]
		shareUserInfos[i] = info.ShareUserInfo{ToUserId: userInfo.UserId,
			Email:             userInfo.Email,
			Perm:              shareNotesMap[userInfo.UserId].Perm,
			NotebookHasShared: hasNotebook,
		}
	}

	return shareUserInfos
}

// 得到notebook的share信息
// TODO 这里必须要分页, 最多取100个用户; 限制
func (this *ShareService) ListNotebookShareUserInfo(notebookId, userId string) []info.ShareUserInfo {
	// notebook的shareNotebooks信息
	shareNotebooks := []info.ShareNotebook{}

	db.ListByQLimit(db.ShareNotebooks,
		bson.M{"NotebookId": bson.ObjectIdHex(notebookId), "UserId": bson.ObjectIdHex(userId)},
		&shareNotebooks, 100)

	if len(shareNotebooks) == 0 {
		return nil
	}

	// 得到用户信息
	toUserIds := make([]bson.ObjectId, len(shareNotebooks))
	for i, each := range shareNotebooks {
		toUserIds[i] = each.ToUserId
	}
	userInfos := userService.ListUserInfosByUserIds(toUserIds)

	if len(userInfos) == 0 {
		return nil
	}

	shareNotebooksMap := make(map[bson.ObjectId]info.ShareNotebook, len(shareNotebooks))
	for _, each := range shareNotebooks {
		shareNotebooksMap[each.ToUserId] = each
	}

	shareUserInfos := make([]info.ShareUserInfo, len(userInfos))
	for i, userInfo := range userInfos {
		shareUserInfos[i] = info.ShareUserInfo{ToUserId: userInfo.UserId,
			Email: userInfo.Email,
			Perm:  shareNotebooksMap[userInfo.UserId].Perm,
		}
	}

	return shareUserInfos
}

//----------------
// 改变note share权限
func (this *ShareService) UpdateShareNotePerm(noteId string, perm int, userId, toUserId string) bool {
	return db.UpdateByQField(db.ShareNotes,
		bson.M{"NoteId": bson.ObjectIdHex(noteId), "UserId": bson.ObjectIdHex(userId), "ToUserId": bson.ObjectIdHex(toUserId)},
		"Perm",
		perm,
	)
}

func (this *ShareService) UpdateShareNotebookPerm(notebookId string, perm int, userId, toUserId string) bool {
	return db.UpdateByQField(db.ShareNotebooks,
		bson.M{"NotebookId": bson.ObjectIdHex(notebookId), "UserId": bson.ObjectIdHex(userId), "ToUserId": bson.ObjectIdHex(toUserId)},
		"Perm",
		perm,
	)
}

//---------------
// 删除share note
func (this *ShareService) DeleteShareNote(noteId string, userId, toUserId string) bool {
	return db.DeleteAll(db.ShareNotes,
		bson.M{"NoteId": bson.ObjectIdHex(noteId), "UserId": bson.ObjectIdHex(userId), "ToUserId": bson.ObjectIdHex(toUserId)})
}

// 删除笔记时要删除该noteId的所有...
func (this *ShareService) DeleteShareNoteAll(noteId string, userId string) bool {
	return db.DeleteAll(db.ShareNotes,
		bson.M{"NoteId": bson.ObjectIdHex(noteId), "UserId": bson.ObjectIdHex(userId)})
}

// 删除share notebook
func (this *ShareService) DeleteShareNotebook(notebookId string, userId, toUserId string) bool {
	return db.DeleteAll(db.ShareNotebooks,
		bson.M{"NotebookId": bson.ObjectIdHex(notebookId), "UserId": bson.ObjectIdHex(userId), "ToUserId": bson.ObjectIdHex(toUserId)})
}

// 删除userId分享给toUserId的所有
func (this *ShareService) DeleteUserShareNoteAndNotebook(userId, toUserId string) bool {
	query := bson.M{"UserId": bson.ObjectIdHex(userId), "ToUserId": bson.ObjectIdHex(toUserId)}
	db.DeleteAll(db.ShareNotebooks, query)
	db.DeleteAll(db.ShareNotes, query)
	db.DeleteAll(db.HasShareNotes, query)

	return true
}
