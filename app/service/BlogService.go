package service

import (
	"github.com/leanote/leanote/app/info"
	"github.com/leanote/leanote/app/db"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
//	"time"
//	"sort"
	"strings"
	"time"
)

// blog
/*
note, notebook都可设为blog
关键是, 怎么得到blog列表? 还要分页

??? 不用新建, 直接使用notes表, 添加IsBlog字段. 新建表 blogs {NoteId, UserId, CreatedTime, IsTop(置顶)}, NoteId, UserId 为unique!!

// 设置一个note为blog
添加到blogs中

// 设置/取消notebook为blog
创建一个note时, 如果其notebookId已设为blog, 那么添加该note到blog中.
设置一个notebook为blog时, 将其下所有的note添加到blogs里 -> 更新其IsBlog为true
取消一个notebook不为blog时, 删除其下的所有note -> 更新其IsBlog为false

*/
type BlogService struct {
}

// 得到某博客具体信息
func (this *BlogService) GetBlog(noteId string) (blog info.BlogItem) {
	note := noteService.GetBlogNote(noteId)
	
	if note.NoteId == "" || !note.IsBlog {
		return
	}
	
	// 内容
	noteContent := noteService.GetNoteContent(note.NoteId.Hex(), note.UserId.Hex())
	
	// 组装成blogItem
	blog = info.BlogItem{note, noteContent.Content, false, info.User{}}	
	
	return
}

// 得到用户共享的notebooks
func (this *BlogService) ListBlogNotebooks(userId string) []info.Notebook {
	notebooks := []info.Notebook{}
	db.ListByQ(db.Notebooks, bson.M{"UserId": bson.ObjectIdHex(userId), "IsBlog": true}, &notebooks)
	return notebooks
}

// 博客列表
// userId 表示谁的blog
func (this *BlogService) ListBlogs(userId, notebookId string, page, pageSize int, sortField string, isAsc bool) (int, []info.BlogItem) {
	count, notes := noteService.ListNotes(userId, notebookId, false, page, pageSize, sortField, isAsc, true);
	
	if(notes == nil || len(notes) == 0) {
		return 0, nil
	}
	
	// 得到content, 并且每个都要substring
	noteIds := make([]bson.ObjectId, len(notes))
	for i, note := range notes {
		noteIds[i] = note.NoteId
	}
	
	// 直接得到noteContents表的abstract
	// 这里可能是乱序的
	noteContents := noteService.ListNoteAbstractsByNoteIds(noteIds) // 返回[info.NoteContent]
	noteContentsMap := make(map[bson.ObjectId]info.NoteContent, len(noteContents))
	for _, noteContent := range noteContents {
		noteContentsMap[noteContent.NoteId] = noteContent
	}
	
	// 组装成blogItem
	// 按照notes的顺序
	blogs := make([]info.BlogItem, len(noteIds))
	for i, note := range notes {
		hasMore := true
		var content string
		if noteContent, ok := noteContentsMap[note.NoteId]; ok {
			content = noteContent.Abstract
		}
		blogs[i] = info.BlogItem{note, content, hasMore, info.User{}}
	}
	return count, blogs
}

func (this *BlogService) SearchBlog(key, userId string, page, pageSize int, sortField string, isAsc bool) (int, []info.BlogItem) {
	count, notes := noteService.SearchNote(key, userId, page, pageSize, sortField, isAsc, true);
	
	if(notes == nil || len(notes) == 0) {
		return 0, nil
	}
	
	// 得到content, 并且每个都要substring
	noteIds := make([]bson.ObjectId, len(notes))
	for i, note := range notes {
		noteIds[i] = note.NoteId
	}
	
	// 直接得到noteContents表的abstract
	// 这里可能是乱序的
	noteContents := noteService.ListNoteAbstractsByNoteIds(noteIds) // 返回[info.NoteContent]
	noteContentsMap := make(map[bson.ObjectId]info.NoteContent, len(noteContents))
	for _, noteContent := range noteContents {
		noteContentsMap[noteContent.NoteId] = noteContent
	}
	
	// 组装成blogItem
	// 按照notes的顺序
	blogs := make([]info.BlogItem, len(noteIds))
	for i, note := range notes {
		hasMore := true
		var content string
		if noteContent, ok := noteContentsMap[note.NoteId]; ok {
			content = noteContent.Abstract
		}
		blogs[i] = info.BlogItem{note, content, hasMore, info.User{}}
	}
	return count, blogs
}

//-------
// p
// 平台 lea+
// 博客列表
func (this *BlogService) ListAllBlogs(tag string, keywords string, isRecommend bool, page, pageSize int, sorterField string, isAsc bool) (info.Page, []info.BlogItem) {
	pageInfo := info.Page{CurPage: page}
	notes := []info.Note{}
	
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, sorterField, isAsc)
	
	// 不是trash的
	query := bson.M{"IsTrash": false, "IsBlog": true, "Title": bson.M{"$ne":"欢迎来到leanote!"}}
	if tag != "" {
		query["Tags"] = bson.M{"$in": []string{tag}}
	}
	// 不是demo的博客
	demoUserId := configService.GetGlobalStringConfig("demoUserId")
	if demoUserId != "" {
		query["UserId"] = bson.M{"$ne": bson.ObjectIdHex(demoUserId)}
	}
	
	if isRecommend {
		query["IsRecommend"] = isRecommend
	}
	if keywords != "" {
		query["Title"] = bson.M{"$regex": bson.RegEx{".*?" + keywords + ".*", "i"}}
	}
	q := db.Notes.Find(query);
	
	// 总记录数
	count, _ := q.Count()
	
	q.Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&notes)
	
	if(notes == nil || len(notes) == 0) {
		return pageInfo, nil
	}
	
	// 得到content, 并且每个都要substring
	noteIds := make([]bson.ObjectId, len(notes))
	userIds := make([]bson.ObjectId, len(notes))
	for i, note := range notes {
		noteIds[i] = note.NoteId
		userIds[i] = note.UserId
	}
	
	// 可以不要的
	// 直接得到noteContents表的abstract
	// 这里可能是乱序的
	/*
	noteContents := noteService.ListNoteAbstractsByNoteIds(noteIds) // 返回[info.NoteContent]
	noteContentsMap := make(map[bson.ObjectId]info.NoteContent, len(noteContents))
	for _, noteContent := range noteContents {
		noteContentsMap[noteContent.NoteId] = noteContent
	}
	*/
	
	// 得到用户信息
	userMap := userService.MapUserInfoAndBlogInfosByUserIds(userIds)
	
	// 组装成blogItem
	// 按照notes的顺序
	blogs := make([]info.BlogItem, len(noteIds))
	for i, note := range notes {
		hasMore := true
		var content string
		/*
		if noteContent, ok := noteContentsMap[note.NoteId]; ok {
			content = noteContent.Abstract
		}
		*/
		blogs[i] = info.BlogItem{note, content, hasMore, userMap[note.UserId]}
	}
	pageInfo = info.NewPage(page, pageSize, count, nil)
	
	return pageInfo, blogs
}

//------------------------
// 博客设置
func (this *BlogService) fixUserBlog(userBlog *info.UserBlog) {
	/*
	if userBlog.Title == "" {
		userInfo := userService.GetUserInfo(userBlog)
		userBlog.Title = userInfo.Username + " 's Blog"
	}
	*/
	
	// Logo路径问题, 有些有http: 有些没有
	Log(userBlog.Logo)
	if userBlog.Logo != "" && !strings.HasPrefix(userBlog.Logo, "http") {
		userBlog.Logo = strings.Trim(userBlog.Logo, "/")
		userBlog.Logo = siteUrl + "/" + userBlog.Logo
	}
}
func (this *BlogService) GetUserBlog(userId string) info.UserBlog {
	userBlog := info.UserBlog{}
	db.Get(db.UserBlogs, userId, &userBlog)
	this.fixUserBlog(&userBlog)
	return userBlog
}

// 修改之
func (this *BlogService) UpdateUserBlog(userBlog info.UserBlog) bool {
	return db.Upsert(db.UserBlogs, bson.M{"_id": userBlog.UserId}, userBlog)
}
// 修改之UserBlogBase
func (this *BlogService) UpdateUserBlogBase(userId string, userBlog info.UserBlogBase) bool {
	ok := db.UpdateByQMap(db.UserBlogs, bson.M{"_id": bson.ObjectIdHex(userId)}, userBlog)
	return ok
}
func (this *BlogService) UpdateUserBlogComment(userId string, userBlog info.UserBlogComment) bool {
	return db.UpdateByQMap(db.UserBlogs, bson.M{"_id": bson.ObjectIdHex(userId)}, userBlog)
}
func (this *BlogService) UpdateUserBlogStyle(userId string, userBlog info.UserBlogStyle) bool {
	return db.UpdateByQMap(db.UserBlogs, bson.M{"_id": bson.ObjectIdHex(userId)}, userBlog)
}


//---------------------
// 后台管理

// 推荐博客
func (this *BlogService) SetRecommend(noteId string, isRecommend bool) bool {
	data := bson.M{"IsRecommend": isRecommend}
	if isRecommend {
		data["RecommendTime"] = time.Now()
	}
	return db.UpdateByQMap(db.Notes, bson.M{"_id": bson.ObjectIdHex(noteId), "IsBlog": true}, data)
}

//----------------------
// 博客社交, 评论

// 返回所有liked用户, bool是否还有
func (this *BlogService) ListLikedUsers(noteId string, isAll bool) ([]info.User, bool) {
	// 默认前5
	pageSize := 5
	skipNum, sortFieldR := parsePageAndSort(1, pageSize, "CreatedTime", false)
		
	likes := []info.BlogLike{}
	query := bson.M{"NoteId": bson.ObjectIdHex(noteId)}
	q := db.BlogLikes.Find(query);
	
	// 总记录数
	count, _ := q.Count()
	if count == 0 {
		return nil, false
	}
	
	if isAll {
		q.Sort(sortFieldR).Skip(skipNum).Limit(pageSize).All(&likes)
	} else {
		q.Sort(sortFieldR).All(&likes)
	}
	
	// 得到所有userIds
	userIds := make([]bson.ObjectId, len(likes))
	for i, like := range likes {
		userIds[i] = like.UserId
	}
	// 得到用户信息
	userMap := userService.MapUserInfoAndBlogInfosByUserIds(userIds)
	
	users := make([]info.User, len(likes));
	for i, like := range likes {
		users[i] = userMap[like.UserId]
	}
	
	return users, count > pageSize
}

func (this *BlogService) IsILikeIt(noteId, userId string) bool {
	if userId == "" {
		return false
	}
	if db.Has(db.BlogLikes, bson.M{"NoteId": bson.ObjectIdHex(noteId), "UserId": bson.ObjectIdHex(userId)}) {
		return true
	}
	return false
}

// 阅读次数统计+1
func (this *BlogService) IncReadNum(noteId string) bool {
	note := noteService.GetNoteById(noteId)
	if note.IsBlog {
		return db.Update(db.Notes, bson.M{"_id": bson.ObjectIdHex(noteId)}, bson.M{"$inc": bson.M{"ReadNum": 1}})
	}
	return false
}

// 点赞
// retun ok , isLike
func (this *BlogService) LikeBlog(noteId, userId string) (ok bool, isLike bool) {
	ok = false
	isLike = false
	if noteId == "" || userId == "" {
		return 
	}
	// 判断是否点过赞, 如果点过那么取消点赞
	note := noteService.GetNoteById(noteId)
	if !note.IsBlog /*|| note.UserId.Hex() == userId */{
		return 
	}
	
	noteIdO := bson.ObjectIdHex(noteId)
	userIdO := bson.ObjectIdHex(userId)
	var n int
	if !db.Has(db.BlogLikes, bson.M{"NoteId": noteIdO, "UserId": userIdO}) {
		n = 1
		// 添加之
		db.Insert(db.BlogLikes, info.BlogLike{LikeId: bson.NewObjectId(), NoteId: noteIdO, UserId: userIdO, CreatedTime: time.Now()})
		isLike = true
	} else {
		// 已点过, 那么删除之
		n = -1
		db.Delete(db.BlogLikes, bson.M{"NoteId": noteIdO, "UserId": userIdO})
		isLike = false
	}
	ok = db.Update(db.Notes, bson.M{"_id": noteIdO}, bson.M{"$inc": bson.M{"LikeNum": n}})
	
	return
}

// 评论
// 在noteId博客下userId 给toUserId评论content
// commentId可为空(针对某条评论评论)
func (this *BlogService) Comment(noteId, toCommentId, userId, content string) (bool, info.BlogComment) {
	var comment info.BlogComment
	if content == "" {
		return false, comment
	}
	
	note := noteService.GetNoteById(noteId)
	if !note.IsBlog {
		return false, comment
	}

	comment = info.BlogComment{CommentId: bson.NewObjectId(), 
		NoteId: bson.ObjectIdHex(noteId), 
		UserId: bson.ObjectIdHex(userId),
		Content: content,
		CreatedTime: time.Now(),
	}
	var comment2 = info.BlogComment{}
	if toCommentId != "" {
		comment2 = info.BlogComment{}
		db.Get(db.BlogComments, toCommentId, &comment2)
		if comment2.CommentId != "" {
			comment.ToCommentId = comment2.CommentId
			comment.ToUserId = comment2.UserId
		}
	} else {
		// comment.ToUserId = note.UserId
	}
	ok := db.Insert(db.BlogComments, comment)
	if ok {
		// 评论+1
		db.Update(db.Notes, bson.M{"_id": bson.ObjectIdHex(noteId)}, bson.M{"$inc": bson.M{"CommentNum": 1}})
	}
	
	if userId != note.UserId.Hex() || toCommentId != "" {
		go func() {
			this.sendEmail(note, comment2, userId, content);
		}()
	}
	
	return ok, comment
}

// 发送email
func (this *BlogService) sendEmail(note info.Note, comment info.BlogComment, userId, content string) {
	emailService.SendCommentEmail(note, comment, userId, content);
	/*
	toUserId := note.UserId.Hex()
	// title := "评论提醒"
	
	// 表示回复回复的内容, 那么发送给之前回复的
	if comment.CommentId != "" {
		toUserId = comment.UserId.Hex()
	}
	toUserInfo := userService.GetUserInfo(toUserId)
	sendUserInfo := userService.GetUserInfo(userId)
	
	subject := note.Title + " 收到 " + sendUserInfo.Username + " 的评论";
	if comment.CommentId != "" {
		subject = "您在 " + note.Title + " 发表的评论收到 " + sendUserInfo.Username;
		if userId == note.UserId.Hex() {
			subject += "(作者)";
		}
		subject += " 的评论";
	}

	body := "{header}<b>评论内容</b>: <br /><blockquote>" + content + "</blockquote>";
	href := "http://"+ configService.GetBlogDomain() + "/view/" + note.NoteId.Hex()
	body += "<br /><b>博客链接</b>: <a href='" + href + "'>" + href + "</a>{footer}";
	
	emailService.SendEmail(toUserInfo.Email, subject, body)
	*/
}

// 作者(或管理员)可以删除所有评论
// 自己可以删除评论
func (this *BlogService) DeleteComment(noteId, commentId, userId string) bool {
	note := noteService.GetNoteById(noteId)
	if !note.IsBlog {
		return false
	}
	
	comment := info.BlogComment{}
	db.Get(db.BlogComments, commentId, &comment)
	
	if comment.CommentId == "" {
		return false
	}
	
	if userId == adminUserId || note.UserId.Hex() == userId || comment.UserId.Hex() == userId {
		 if db.Delete(db.BlogComments, bson.M{"_id": bson.ObjectIdHex(commentId)}) {
			// 评论-1
			db.Update(db.Notes, bson.M{"_id": bson.ObjectIdHex(noteId)}, bson.M{"$inc": bson.M{"CommentNum": -1}})
			return true
		 }
	}
		
	return false
}

// 点赞/取消赞
func (this *BlogService) LikeComment(commentId, userId string) (ok bool, isILike bool, num int) {
	ok = false
	isILike = false
	num = 0
	comment := info.BlogComment{}
	
	db.Get(db.BlogComments, commentId, &comment)
	
	var n int
	if comment.LikeUserIds != nil && len(comment.LikeUserIds) > 0 && InArray(comment.LikeUserIds, userId) {
		n = -1
		// 从点赞名单删除
		db.Update(db.BlogComments, bson.M{"_id": bson.ObjectIdHex(commentId)}, 
			bson.M{"$pull":  bson.M{"LikeUserIds": userId}})
		isILike = false
	} else {
		n = 1
		// 添加之
		db.Update(db.BlogComments, bson.M{"_id": bson.ObjectIdHex(commentId)}, 
			bson.M{"$push": bson.M{"LikeUserIds": userId}})
		isILike = true
	}
	
	if comment.LikeUserIds == nil {
		num = 0
	} else {
		num = len(comment.LikeUserIds) + n
	}
	
	ok = db.Update(db.BlogComments, bson.M{"_id": bson.ObjectIdHex(commentId)}, 
			bson.M{"$set": bson.M{"LikeNum": num}})
			
	return
}

// 评论列表
// userId主要是显示userId是否点过某评论的赞
// 还要获取用户信息
func (this *BlogService) ListComments(userId, noteId string, page, pageSize int) (info.Page, []info.BlogCommentPublic, map[string]info.User) {
	pageInfo := info.Page{CurPage: page}
	
	comments2 := []info.BlogComment{}
	
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "CreatedTime", false)
		
	query := bson.M{"NoteId": bson.ObjectIdHex(noteId)}
	q := db.BlogComments.Find(query);
	
	// 总记录数
	count, _ := q.Count()
	q.Sort(sortFieldR).Skip(skipNum).Limit(pageSize).All(&comments2)
	
	if(len(comments2) == 0) {
		return pageInfo, nil, nil
	}
	
	comments := make([]info.BlogCommentPublic, len(comments2))
	// 我是否点过赞呢?
	for i, comment := range comments2 {
		comments[i].BlogComment = comment
		if comment.LikeNum > 0 && comment.LikeUserIds != nil && len(comment.LikeUserIds) > 0 && InArray(comment.LikeUserIds, userId) {
			comments[i].IsILikeIt = true
		}
	}
	
	note := noteService.GetNoteById(noteId);
	
	// 得到用户信息
	userIdsMap := map[bson.ObjectId]bool{note.UserId: true}
	for _, comment := range comments {
		userIdsMap[comment.UserId] = true
		if comment.ToUserId != "" { // 可能为空
			userIdsMap[comment.ToUserId] = true
		}
	}
	userIds := make([]bson.ObjectId, len(userIdsMap))
	i := 0
	for userId, _ := range userIdsMap {
		userIds[i] = userId
		i++
	}
	
	// 得到用户信息
	userMap := userService.MapUserInfoByUserIds(userIds)
	userMap2 := make(map[string]info.User, len(userMap))
	for userId, v := range userMap {
		userMap2[userId.Hex()] = v
	}
	
	pageInfo = info.NewPage(page, pageSize, count, nil)
	
	return pageInfo, comments, userMap2
}

// 举报
func (this *BlogService) Report(noteId, commentId, reason, userId string) (bool) {
	note := noteService.GetNoteById(noteId)
	if !note.IsBlog {
		return false
	}

	report := info.Report{ReportId: bson.NewObjectId(), 
		NoteId: bson.ObjectIdHex(noteId), 
		UserId: bson.ObjectIdHex(userId),
		Reason: reason,
		CreatedTime: time.Now(),
	}
	if commentId != "" {
		report.CommentId = bson.ObjectIdHex(commentId)
	}
	return db.Insert(db.Reports, report)
}