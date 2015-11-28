package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"time"
)

// 共享Notebook, Note服务
type ShareService struct {
}

//-----------------------------------
// 返回shareNotebooks, sharedUserInfos
// info.ShareNotebooksByUser, []info.User

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

// 谁共享给了我的Query
func (this *ShareService) getOrQ(userId string) bson.M {
	// 得到我和和我参与的组织
	groupIds := groupService.GetMineAndBelongToGroupIds(userId)

	q := bson.M{}
	if len(groupIds) > 0 {
		orQ := []bson.M{
			bson.M{"ToUserId": bson.ObjectIdHex(userId)},
			bson.M{"ToGroupId": bson.M{"$in": groupIds}},
		}
		// 不是trash的
		q["$or"] = orQ
	} else {
		q["ToUserId"] = bson.ObjectIdHex(userId)
	}
	return q
}

// 得到共享给我的笔记本和用户(谁共享给了我)
func (this *ShareService) GetShareNotebooks(userId string) (info.ShareNotebooksByUser, []info.User) {
	myUserId := userId

	// 得到共享给我的用户s信息
	// 得到我参与的组织
	q := this.getOrQ(userId)

	// 不查hasShareNotes
	// 直接查shareNotes, shareNotebooks表得到userId
	userIds1 := []bson.ObjectId{}
	db.Distinct(db.ShareNotes, q, "UserId", &userIds1)

	userIds2 := []bson.ObjectId{}
	db.Distinct(db.ShareNotebooks, q, "UserId", &userIds2) // BUG之前是userId1, 2014/12/29

	userIds := append(userIds1, userIds2...)

	userInfos := userService.GetUserInfosOrderBySeq(userIds)
	// 不要我的id
	for i, userInfo := range userInfos {
		if userInfo.UserId.Hex() == myUserId {
			userInfos = append(userInfos[:i], userInfos[i+1:]...)
			break
		}
	}

	//--------------------
	// 得到他们共享给我的notebooks

	// 这里可能会得到重复的记录
	// 权限: 笔记本分享给个人 > 笔记本分享给组织
	shareNotebooks := []info.ShareNotebook{}
	db.ShareNotebooks.Find(q).Sort("-ToUserId").All(&shareNotebooks) // 按ToUserId降序排序, 那么有ToUserId的在前面

	if len(shareNotebooks) == 0 {
		return nil, userInfos
	}

	shareNotebooksLen := len(shareNotebooks)

	// 找到了所有的notbookId, 那么找notebook表得到其详细信息
	notebookIds := []bson.ObjectId{}
	shareNotebooksMap := make(map[bson.ObjectId]info.ShareNotebook, shareNotebooksLen)
	for _, each := range shareNotebooks {
		// 之后的就不要了, 只留权限高的
		if _, ok := shareNotebooksMap[each.NotebookId]; !ok {
			// 默认的是没有notebookId的
			notebookIds = append(notebookIds, each.NotebookId)
			shareNotebooksMap[each.NotebookId] = each
		}
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
		// 我自己的, 算了
		if userId.Hex() == myUserId {
			continue
		}

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
	q := this.getOrQ(myUserId)
	q["NotebookId"] = bson.ObjectIdHex(notebookId)
	q["UserId"] = bson.ObjectIdHex(sharedUserId)
	shareNotebook := info.ShareNotebook{}
	db.GetByQ(db.ShareNotebooks,
		q,
		&shareNotebook)

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
	// 笔记的权限
	shareNotes := []info.ShareNote{}
	delete(q, "NotebookId")
	q["NoteId"] = bson.M{"$in": noteIds}
	db.ShareNotes.Find(q).Sort("-ToUserId").All(&shareNotes) // 给个的权限>给组织的权限
	notePerms := map[bson.ObjectId]int{}
	for _, each := range shareNotes {
		if _, ok := notePerms[each.NoteId]; !ok {
			notePerms[each.NoteId] = each.Perm
		}
	}
	Log("笔记权限")
	LogJ(notePerms)

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
//func (this *ShareService) getNotesPerm(noteIds []bson.ObjectId, myUserId, sharedUserId string) map[bson.ObjectId]int {
//	shareNotes := []info.ShareNote{}
//	db.ListByQ(db.ShareNotes,
//		bson.M{
//			"NoteId": bson.M{"$in": noteIds},
//			"UserId": bson.ObjectIdHex(sharedUserId),
//			"ToUserId": bson.ObjectIdHex(myUserId)}, &shareNotes)
//
//	notesPerm := make(map[bson.ObjectId]int, len(shareNotes))
//	for _, each := range shareNotes {
//		notesPerm[each.NoteId] = each.Perm
//	}
//
//	return notesPerm
//}

// 得到默认的单个的notes 共享集
// 如果真要支持排序, 这里得到所有共享的notes, 到noteService方再sort和limit
// 可以这样! 到时将零散的共享noteId放在用户基本数据中
// 这里不好排序
func (this *ShareService) ListShareNotes(myUserId, sharedUserId string,
	pageNumber, pageSize int, sortField string, isAsc bool) []info.ShareNoteWithPerm {

	skipNum, _ := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)
	shareNotes := []info.ShareNote{}

	q := this.getOrQ(myUserId)
	q["UserId"] = bson.ObjectIdHex(sharedUserId)

	db.ShareNotes.
		Find(q).
		Sort("-ToUserId"). // 给个人的权限 > 给组织的权限
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
	notesMap := map[bson.ObjectId]info.Note{}
	for _, each := range notes {
		notesMap[each.NoteId] = each
	}

	// 将shareNotes与notes结合起来
	notesWithPerm := []info.ShareNoteWithPerm{}
	hasAdded := map[bson.ObjectId]bool{} // 防止重复, 只要前面权限高的
	for _, each := range shareNotes {
		if !hasAdded[each.NoteId] {
			// 待优化
			notesWithPerm = append(notesWithPerm, info.ShareNoteWithPerm{notesMap[each.NoteId], each.Perm})
			hasAdded[each.NoteId] = true
		}
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
	return this.AddShareNotebookToUserId(notebookId, perm, userId, toUserId)
}

// 第三方注册时没有email
func (this *ShareService) AddShareNotebookToUserId(notebookId string, perm int, userId, toUserId string) (bool, string, string) {
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
	return this.AddShareNoteToUserId(noteId, perm, userId, toUserId)
}

// 第三方测试没有userId
func (this *ShareService) AddShareNoteToUserId(noteId string, perm int, userId, toUserId string) (bool, string, string) {
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

// updatedUserId是否有查看userId noteId的权限?
// userId是所有者
func (this *ShareService) HasReadPerm(userId, updatedUserId, noteId string) bool {
	q := this.getOrQ(updatedUserId) // (toUserId == "xxx" || ToGroupId in (1, 2,3))
	q["UserId"] = bson.ObjectIdHex(userId)
	q["NoteId"] = bson.ObjectIdHex(noteId)
	if !db.Has(db.ShareNotes, q) {
		// noteId的notebookId是否被共享了?
		notebookId := noteService.GetNotebookId(noteId)
		if notebookId.Hex() == "" {
			return false
		}

		delete(q, "NoteId")
		q["NotebookId"] = notebookId

		// 判断notebook是否被共享
		if !db.Has(db.ShareNotebooks, q) {
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

// updatedUserId是否有修改userId noteId的权限?
func (this *ShareService) HasUpdatePerm(userId, updatedUserId, noteId string) bool {
	q := this.getOrQ(updatedUserId) // (toUserId == "xxx" || ToGroupId in (1, 2,3))
	q["UserId"] = bson.ObjectIdHex(userId)
	q["NoteId"] = bson.ObjectIdHex(noteId)

	// note的权限
	shares := []info.ShareNote{}
	db.ShareNotes.Find(q).Sort("-ToUserId").All(&shares) // 个人 > 组织
	for _, share := range shares {
		return share.Perm == 1 // 第1个权限最大
	}

	// notebook的权限
	notebookId := noteService.GetNotebookId(noteId)
	if notebookId.Hex() == "" {
		return false
	}

	delete(q, "NoteId")
	q["NotebookId"] = notebookId
	shares2 := []info.ShareNotebook{}
	db.ShareNotebooks.Find(q).Sort("-ToUserId").All(&shares2) // 个人 > 组织
	for _, share := range shares2 {
		return share.Perm == 1 // 第1个权限最大
	}
	return false
}

// updatedUserId是否有修改userId notebookId的权限?
func (this *ShareService) HasUpdateNotebookPerm(userId, updatedUserId, notebookId string) bool {
	q := this.getOrQ(updatedUserId) // (toUserId == "xxx" || ToGroupId in (1, 2,3))
	q["UserId"] = bson.ObjectIdHex(userId)
	q["NotebookId"] = bson.ObjectIdHex(notebookId)
	shares2 := []info.ShareNotebook{}
	db.ShareNotebooks.Find(q).Sort("-ToUserId").All(&shares2) // 个人 > 组织
	for _, share := range shares2 {
		return share.Perm == 1 // 第1个权限最大
	}
	return false
}

// 共享note, notebook时使用
func (this *ShareService) AddHasShareNote(userId, toUserId string) bool {
	db.Insert(db.HasShareNotes, info.HasShareNote{UserId: bson.ObjectIdHex(userId), ToUserId: bson.ObjectIdHex(toUserId)})
	return true
}

// userId是否被共享了noteId
func (this *ShareService) HasSharedNote(noteId, myUserId string) bool {
	return db.Has(db.ShareNotes, bson.M{"ToUserId": bson.ObjectIdHex(myUserId), "NoteId": bson.ObjectIdHex(noteId)})
}

// noteId的notebook是否共享了给我
func (this *ShareService) HasSharedNotebook(noteId, myUserId, sharedUserId string) bool {
	notebookId := noteService.GetNotebookId(noteId)
	if notebookId != "" {
		return db.Has(db.ShareNotebooks, bson.M{"NotebookId": notebookId,
			"UserId":   bson.ObjectIdHex(sharedUserId),
			"ToUserId": bson.ObjectIdHex(myUserId),
		})
	}
	return false
}

// 得到共享的笔记内容
// 并返回笔记的权限!!!
func (this *ShareService) GetShareNoteContent(noteId, myUserId, sharedUserId string) (noteContent info.NoteContent) {
	noteContent = info.NoteContent{}
	// 是否单独共享了该notebook
	// 或者, 其notebook共享了我
	//	Log(this.HasSharedNote(noteId, myUserId))
	//	Log(this.HasSharedNotebook(noteId, myUserId, sharedUserId))
	//	Log(this.HasReadPerm(sharedUserId, myUserId, noteId))

	if this.HasReadPerm(sharedUserId, myUserId, noteId) {
		//	if this.HasSharedNote(noteId, myUserId) || this.HasSharedNotebook(noteId, myUserId, sharedUserId) {
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
	db.ListByQLimit(db.ShareNotes,
		bson.M{
			"NoteId":    bson.ObjectIdHex(noteId),
			"UserId":    bson.ObjectIdHex(userId),
			"ToGroupId": bson.M{"$exists": false},
		}, &shareNotes, 100)

	//	Log("<<>>>>")
	//	Log(len(shareNotes))

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
		bson.M{
			"NotebookId": bson.ObjectIdHex(notebookId),
			"UserId":     bson.ObjectIdHex(userId),
			"ToGroupId":  bson.M{"$exists": false},
		},
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

// 用户userId是否有修改noteId的权限
func (this *ShareService) HasUpdateNotePerm(noteId, userId string) bool {
	if noteId == "" || userId == "" {
		return false
	}
	note := noteService.GetNoteById(noteId)
	LogJ(note)
	if note.UserId != "" {
		noteUserId := note.UserId.Hex()
		if noteUserId != userId {
			// 是否是有权限协作的
			if this.HasUpdatePerm(noteUserId, userId, noteId) {
				return true
			} else {
				return false
			}
		} else {
			return true
		}
	} else {
		return false
	}
}

// 用户userId是否有查看noteId的权限
func (this *ShareService) HasReadNotePerm(noteId, userId string) bool {
	if noteId == "" || userId == "" {
		return false
	}
	note := noteService.GetNoteById(noteId)
	if note.UserId != "" {
		noteUserId := note.UserId.Hex()
		if noteUserId != userId {
			// 是否是有权限协作的
			if this.HasReadPerm(noteUserId, userId, noteId) {
				return true
			} else {
				return false
			}
		} else {
			return true
		}
	} else {
		return false
	}
}

//----------------
// 用户分组

// 得到笔记分享给的groups
func (this *ShareService) GetNoteShareGroups(noteId, userId string) []info.ShareNote {
	groups := groupService.GetGroupsContainOf(userId)

	// 得到有分享的分组
	shares := []info.ShareNote{}
	db.ListByQ(db.ShareNotes,
		bson.M{"NoteId": bson.ObjectIdHex(noteId), "UserId": bson.ObjectIdHex(userId), "ToGroupId": bson.M{"$exists": true}}, &shares)
	mapShares := map[bson.ObjectId]info.ShareNote{}
	for _, share := range shares {
		mapShares[share.ToGroupId] = share
	}

	// 所有的groups都有share, 但没有share的group没有shareId
	shares2 := make([]info.ShareNote, len(groups))
	for i, group := range groups {
		share, ok := mapShares[group.GroupId]
		if !ok {
			share = info.ShareNote{}
		}
		share.ToGroup = group
		shares2[i] = share
	}

	return shares2
}

// 共享笔记给分组
func (this *ShareService) AddShareNoteGroup(userId, noteId, groupId string, perm int) bool {
	if !groupService.IsExistsGroupUser(userId, groupId) {
		return false
	}

	// 先删除之
	this.DeleteShareNoteGroup(userId, noteId, groupId)

	shareNote := info.ShareNote{NoteId: bson.ObjectIdHex(noteId),
		UserId:      bson.ObjectIdHex(userId), // 冗余字段
		ToGroupId:   bson.ObjectIdHex(groupId),
		Perm:        perm,
		CreatedTime: time.Now(),
	}
	return db.Insert(db.ShareNotes, shareNote)
}

// 删除
func (this *ShareService) DeleteShareNoteGroup(userId, noteId, groupId string) bool {
	return db.Delete(db.ShareNotes, bson.M{"NoteId": bson.ObjectIdHex(noteId),
		"UserId":    bson.ObjectIdHex(userId),
		"ToGroupId": bson.ObjectIdHex(groupId),
	})
}

//-------

// 得到笔记本分享给的groups
func (this *ShareService) GetNotebookShareGroups(notebookId, userId string) []info.ShareNotebook {
	groups := groupService.GetGroupsContainOf(userId)

	// 得到有分享的分组
	shares := []info.ShareNotebook{}
	db.ListByQ(db.ShareNotebooks,
		bson.M{"NotebookId": bson.ObjectIdHex(notebookId), "UserId": bson.ObjectIdHex(userId), "ToGroupId": bson.M{"$exists": true}}, &shares)
	mapShares := map[bson.ObjectId]info.ShareNotebook{}
	for _, share := range shares {
		mapShares[share.ToGroupId] = share
	}
	LogJ(shares)

	// 所有的groups都有share, 但没有share的group没有shareId
	shares2 := make([]info.ShareNotebook, len(groups))
	for i, group := range groups {
		share, ok := mapShares[group.GroupId]
		if !ok {
			share = info.ShareNotebook{}
		}
		share.ToGroup = group
		shares2[i] = share
	}

	return shares2
}

// 共享笔记给分组
func (this *ShareService) AddShareNotebookGroup(userId, notebookId, groupId string, perm int) bool {
	if !groupService.IsExistsGroupUser(userId, groupId) {
		return false
	}

	// 先删除之
	this.DeleteShareNotebookGroup(userId, notebookId, groupId)

	shareNotebook := info.ShareNotebook{NotebookId: bson.ObjectIdHex(notebookId),
		UserId:      bson.ObjectIdHex(userId), // 冗余字段
		ToGroupId:   bson.ObjectIdHex(groupId),
		Perm:        perm,
		CreatedTime: time.Now(),
	}
	return db.Insert(db.ShareNotebooks, shareNotebook)
}

// 删除
func (this *ShareService) DeleteShareNotebookGroup(userId, notebookId, groupId string) bool {
	return db.Delete(db.ShareNotebooks, bson.M{"NotebookId": bson.ObjectIdHex(notebookId),
		"UserId":    bson.ObjectIdHex(userId),
		"ToGroupId": bson.ObjectIdHex(groupId),
	})
}

//--------------------
// 删除组时, 删除所有的
//--------------------

func (this *ShareService) DeleteAllShareNotebookGroup(groupId string) bool {
	return db.Delete(db.ShareNotebooks, bson.M{
		"ToGroupId": bson.ObjectIdHex(groupId),
	})
}
func (this *ShareService) DeleteAllShareNoteGroup(groupId string) bool {
	return db.Delete(db.ShareNotes, bson.M{
		"ToGroupId": bson.ObjectIdHex(groupId),
	})
}

//--------------------
// 删除组内用户时, 删除其分享的
//--------------------

func (this *ShareService) DeleteShareNotebookGroupWhenDeleteGroupUser(userId, groupId string) bool {
	return db.Delete(db.ShareNotebooks, bson.M{
		"UserId":    bson.ObjectIdHex(userId),
		"ToGroupId": bson.ObjectIdHex(groupId),
	})
}

func (this *ShareService) DeleteShareNoteGroupWhenDeleteGroupUser(userId, groupId string) bool {
	return db.Delete(db.ShareNotes, bson.M{
		"UserId":    bson.ObjectIdHex(userId),
		"ToGroupId": bson.ObjectIdHex(groupId),
	})
}
