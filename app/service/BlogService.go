package service

import (
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	//	"time"
	//	"sort"
	"strconv"
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

// 得到博客统计信息
// ReadNum, LikeNum, CommentNum
func (this *BlogService) GetBlogStat(noteId string) (stat info.BlogStat) {
	note := noteService.GetBlogNote(noteId)
	stat = info.BlogStat{note.NoteId, note.ReadNum, note.LikeNum, note.CommentNum}
	return
}

// 通过id或urlTitle得到博客
func (this *BlogService) GetBlogByIdAndUrlTitle(userId string, noteIdOrUrlTitle string) (blog info.BlogItem) {
	if IsObjectId(noteIdOrUrlTitle) {
		return this.GetBlog(noteIdOrUrlTitle)
	}
	note := info.Note{}
	db.GetByQ(db.Notes, bson.M{"UserId": bson.ObjectIdHex(userId), "UrlTitle": encodeValue(noteIdOrUrlTitle),
		"IsBlog":  true,
		"IsTrash": false, "IsDeleted": false}, &note)
	return this.GetBlogItem(note)
}

// 得到某博客具体信息
func (this *BlogService) GetBlog(noteId string) (blog info.BlogItem) {
	note := noteService.GetBlogNote(noteId)
	return this.GetBlogItem(note)
}
func (this *BlogService) GetBlogItem(note info.Note) (blog info.BlogItem) {
	if note.NoteId == "" || !note.IsBlog {
		return info.BlogItem{}
	}

	// 内容
	noteContent := noteService.GetNoteContent(note.NoteId.Hex(), note.UserId.Hex())

	// 组装成blogItem
	blog = info.BlogItem{note, noteContent.Abstract, noteContent.Content, false, info.User{}}

	return
}

// 得到用户共享的notebooks
// 3/19 博客不是deleted
func (this *BlogService) ListBlogNotebooks(userId string) []info.Notebook {
	notebooks := []info.Notebook{}
	orQ := []bson.M{
		bson.M{"IsDeleted": false},
		bson.M{"IsDeleted": bson.M{"$exists": false}},
	}
	db.ListByQ(db.Notebooks, bson.M{"UserId": bson.ObjectIdHex(userId), "IsBlog": true, "$or": orQ}, &notebooks)
	return notebooks
}

// 博客列表
// userId 表示谁的blog
func (this *BlogService) ListBlogs(userId, notebookId string, page, pageSize int, sortField string, isAsc bool) (info.Page, []info.BlogItem) {
	count, notes := noteService.ListNotes(userId, notebookId, false, page, pageSize, sortField, isAsc, true)

	if notes == nil || len(notes) == 0 {
		return info.Page{}, nil
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
		var abstract string
		if noteContent, ok := noteContentsMap[note.NoteId]; ok {
			abstract = noteContent.Abstract
			content = noteContent.Content
		}
		blogs[i] = info.BlogItem{note, abstract, content, hasMore, info.User{}}
	}

	pageInfo := info.NewPage(page, pageSize, count, nil)

	return pageInfo, blogs
}

// 得到博客的标签, 那得先得到所有博客, 比较慢
/*
[
	{Tag:xxx, Count: 32}
]
*/
func (this *BlogService) GetBlogTags(userId string) []info.TagCount {
	// 得到所有博客
	tagCounts := []info.TagCount{}
	// tag不能为空
	query := bson.M{"UserId": bson.ObjectIdHex(userId), "IsBlog": true, "Tag": bson.M{"$ne": ""}}
	db.TagCounts.Find(query).Sort("-Count").All(&tagCounts)
	return tagCounts
}

// 重新计算博客的标签
// 在设置设置/取消为博客时调用
func (this *BlogService) ReCountBlogTags(userId string) bool {
	// 得到所有博客
	notes := []info.Note{}
	userIdO := bson.ObjectIdHex(userId)
	query := bson.M{"UserId": userIdO, "IsTrash": false, "IsDeleted": false, "IsBlog": true}
	db.ListByQWithFields(db.Notes, query, []string{"Tags"}, &notes)

	db.DeleteAll(db.TagCounts, bson.M{"UserId": userIdO, "IsBlog": true})
	if notes == nil || len(notes) == 0 {
		return true
	}
	// 统计所有的Tags和数目
	tagsCount := map[string]int{}
	for _, note := range notes {
		tags := note.Tags
		if tags != nil && len(tags) > 0 {
			for _, tag := range tags {
				count := tagsCount[tag]
				count++
				tagsCount[tag] = count
			}
		}
	}
	// 一个个插入
	for tag, count := range tagsCount {
		db.Insert(db.TagCounts,
			info.TagCount{UserId: userIdO, IsBlog: true, Tag: tag, Count: count})
	}
	return true
}

// 归档博客
/*
数据: 按年汇总
[
archive1,
archive2,
]
archive的数据类型是
{
Year: 2014
Posts: []
}
*/
func (this *BlogService) ListBlogsArchive(userId, notebookId string, year, month int, sortField string, isAsc bool) []info.Archive {
	//	_, notes := noteService.ListNotes(userId, notebookId, false, 1, 99999, sortField, isAsc, true);
	q := bson.M{"UserId": bson.ObjectIdHex(userId), "IsBlog": true, "IsTrash": false, "IsDeleted": false}
	if notebookId != "" {
		q["NotebookId"] = bson.ObjectIdHex(notebookId)
	}
	if year > 0 {
		now := time.Now()
		nextYear := year
		nextMonth := month
		if month == 0 {
			month = 1
			nextYear = year + 1
			nextMonth = month
		} else if month >= 12 {
			month = 12
			nextYear = year + 1
			nextMonth = 1
		} else { // month 在1-12之间
			nextMonth = month + 1
		}
		leftT := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
		rightT := time.Date(nextYear, time.Month(nextMonth), 1, 0, 0, 0, 0, now.Location())
		if sortField == "CreatedTime" || sortField == "UpdatedTime" {
			q[sortField] = bson.M{"$gte": leftT, "$lt": rightT}
		} else {
			q["PublicTime"] = bson.M{"$gte": leftT, "$lt": rightT}
		}
	}

	sorter := sortField
	if !isAsc {
		sorter = "-" + sortField
	}
	notes := []info.Note{}
	db.Notes.Find(q).Sort(sorter).All(&notes)

	if notes == nil || len(notes) == 0 {
		return nil
	}

	arcs := []info.Archive{}
	// 按年汇总
	arcsMap := map[int]info.Archive{}
	// 按月汇总
	arcsMonth := []info.ArchiveMonth{}
	var t time.Time
	var arc info.Archive
	everYear := 0
	for _, note := range notes {
		if sortField == "PublicTime" {
			t = note.PublicTime
		} else if sortField == "CreatedTime" {
			t = note.CreatedTime
		} else {
			t = note.UpdatedTime
		}
		year := t.Year()
		month := int(t.Month())
		if everYear == 0 {
			everYear = year
		}

		if everYear != year {
			yearArc := arcsMap[everYear]
			yearArc.MonthAchives = arcsMonth
			arcs = append(arcs, yearArc)
			everYear = year

			// 新的一年
			arcsMonth = []info.ArchiveMonth{}
		}

		if arcT, ok := arcsMap[year]; ok {
			arc = arcT
		} else {
			arc = info.Archive{Year: year, Posts: []*info.Post{}}
		}

		pt := this.FixNote(note)
		p := &pt
		arc.Posts = append(arc.Posts, p)
		arcsMap[year] = arc

		// month
		lm := len(arcsMonth)
		if lm == 0 || arcsMonth[lm-1].Month != month {
			arcsMonth = append(arcsMonth, info.ArchiveMonth{month, []*info.Post{p}})
		} else {
			arcsMonth[lm-1].Posts = append(arcsMonth[lm-1].Posts, p)
		}
	}
	// 最后一个
	if everYear > 0 {
		yearArc := arcsMap[everYear]
		yearArc.MonthAchives = arcsMonth
		arcs = append(arcs, yearArc)
	}

	return arcs
}

// 根据tag搜索博客
func (this *BlogService) SearchBlogByTags(tags []string, userId string, pageNumber, pageSize int, sortField string, isAsc bool) (pageInfo info.Page, blogs []info.BlogItem) {
	notes := []info.Note{}
	skipNum, sortFieldR := parsePageAndSort(pageNumber, pageSize, sortField, isAsc)

	// 不是trash的
	query := bson.M{"UserId": bson.ObjectIdHex(userId),
		"IsTrash":   false,
		"IsDeleted": false,
		"IsBlog":    true,
		"Tags":      bson.M{"$all": tags}}

	q := db.Notes.Find(query)

	// 总记录数
	count, _ := q.Count()
	if count == 0 {
		return
	}

	q.Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&notes)

	blogs = this.notes2BlogItems(notes)
	pageInfo = info.NewPage(pageNumber, pageSize, count, nil)

	return
}

func (this *BlogService) notes2BlogItems(notes []info.Note) []info.BlogItem {
	// 得到content, 并且每个都要substring
	noteIds := make([]bson.ObjectId, len(notes))
	for i, note := range notes {
		noteIds[i] = note.NoteId
	}

	// 直接得到noteContents表的abstract
	// 这里可能是乱序的
	noteContents := noteService.ListNoteContentByNoteIds(noteIds) // 返回[info.NoteContent]
	noteContentsMap := make(map[bson.ObjectId]info.NoteContent, len(noteContents))
	for _, noteContent := range noteContents {
		noteContentsMap[noteContent.NoteId] = noteContent
	}

	// 组装成blogItem
	// 按照notes的顺序
	blogs := make([]info.BlogItem, len(noteIds))
	for i, note := range notes {
		hasMore := true
		var content, abstract string
		if noteContent, ok := noteContentsMap[note.NoteId]; ok {
			abstract = noteContent.Abstract
			content = noteContent.Content
		}
		blogs[i] = info.BlogItem{note, abstract, content, hasMore, info.User{}}
	}
	return blogs
}
func (this *BlogService) SearchBlog(key, userId string, page, pageSize int, sortField string, isAsc bool) (info.Page, []info.BlogItem) {
	count, notes := noteService.SearchNote(key, userId, page, pageSize, sortField, isAsc, true)

	if notes == nil || len(notes) == 0 {
		return info.Page{}, nil
	}

	blogs := this.notes2BlogItems(notes)
	pageInfo := info.NewPage(page, pageSize, count, nil)
	return pageInfo, blogs
}

// 上一篇文章, 下一篇文章
// sorterField, baseTime是基准, sorterField=PublicTime, title
// isAsc是用户自定义的排序方式
func (this *BlogService) PreNextBlog(userId string, sorterField string, isAsc bool, noteId string, baseTime interface{}) (info.Post, info.Post) {
	userIdO := bson.ObjectIdHex(userId)

	var sortFieldT1, sortFieldT2 bson.M
	var sortFieldR1, sortFieldR2 string
	if !isAsc {
		// 降序
		/*
			------- pre
			----- now
			--- next
			--
		*/
		// 上一篇时间要比它大, 找最小的
		sortFieldT1 = bson.M{"$gte": baseTime} // 为什么要相等, 因为将notebook发布成博客, 会统一修改note的publicTime, 此时所有notes都一样
		sortFieldR1 = sorterField
		// 下一篇时间要比它小
		sortFieldT2 = bson.M{"$lte": baseTime}
		sortFieldR2 = "-" + sorterField
	} else {
		// 升序
		/*
		   --- pre
		   ----- now
		   ------- next
		   ---------
		*/
		// 上一篇要比它小, 找最大的
		sortFieldT1 = bson.M{"$lte": baseTime}
		sortFieldR1 = "-" + sorterField
		// 下一篇, 找最小的
		sortFieldT2 = bson.M{"$gte": baseTime}
		sortFieldR2 = sorterField
	}

	// 1
	// 上一篇, 比基时间要小, 但是是最后一篇, 所以是降序
	note := info.Note{}
	query := bson.M{"UserId": userIdO,
		"IsTrash":   false,
		"IsDeleted": false,
		"IsBlog":    true,
		"_id":       bson.M{"$ne": bson.ObjectIdHex(noteId)},
		sorterField: sortFieldT1,
	}
	q := db.Notes.Find(query)
	q.Sort(sortFieldR1).Limit(1).One(&note)

	// 下一篇, 比基时间要大, 但是是第一篇, 所以是升序
	if note.NoteId != "" {
		query["_id"] = bson.M{"$nin": []bson.ObjectId{bson.ObjectIdHex(noteId), note.NoteId}}
	}
	note2 := info.Note{}
	query[sorterField] = sortFieldT2
	//	Log(isAsc)
	//	LogJ(query)
	//	Log(sortFieldR2)
	q = db.Notes.Find(query)
	q.Sort(sortFieldR2).Limit(1).One(&note2)

	return this.FixNote(note), this.FixNote(note2)
}

//-------
// p
// 平台 lea+
// 博客列表
func (this *BlogService) ListAllBlogs(userId, tag string, keywords string, isRecommend bool, page, pageSize int, sorterField string, isAsc bool) (info.Page, []info.BlogItem) {
	pageInfo := info.Page{CurPage: page}
	notes := []info.Note{}

	skipNum, sortFieldR := parsePageAndSort(page, pageSize, sorterField, isAsc)

	// 不是trash的
	query := bson.M{"IsTrash": false, "IsDeleted": false, "IsBlog": true, "Title": bson.M{"$ne": "欢迎来到leanote!"}}
	if tag != "" {
		query["Tags"] = bson.M{"$in": []string{tag}}
	}
	if userId != "" {
		query["UserId"] = bson.ObjectIdHex(userId)
	}
	// 不是demo的博客
	demoUserId := configService.GetGlobalStringConfig("demoUserId")
	if userId == "" && demoUserId != "" {
		query["UserId"] = bson.M{"$ne": bson.ObjectIdHex(demoUserId)}
	}

	if isRecommend {
		query["IsRecommend"] = isRecommend
	}
	if keywords != "" {
		query["Title"] = bson.M{"$regex": bson.RegEx{".*?" + keywords + ".*", "i"}}
	}
	q := db.Notes.Find(query)

	// 总记录数
	count, _ := q.Count()

	q.Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&notes)

	if notes == nil || len(notes) == 0 {
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
		if len(note.Tags) == 1 && note.Tags[0] == "" {
			note.Tags = nil
		}
		blogs[i] = info.BlogItem{note, "", content, hasMore, userMap[note.UserId]}
	}
	pageInfo = info.NewPage(page, pageSize, count, nil)

	return pageInfo, blogs
}

//------------------------
// 博客设置
func (this *BlogService) fixUserBlog(userBlog *info.UserBlog) {
	// Logo路径问题, 有些有http: 有些没有
	if userBlog.Logo != "" && !strings.HasPrefix(userBlog.Logo, "http") {
		userBlog.Logo = strings.Trim(userBlog.Logo, "/")
		userBlog.Logo = configService.GetSiteUrl() + "/" + userBlog.Logo
	}

	if userBlog.SortField == "" {
		userBlog.SortField = "PublicTime"
	}
	if userBlog.PerPageSize <= 0 {
		userBlog.PerPageSize = 10
	}

	// themePath
	if userBlog.Style == "" {
		userBlog.Style = defaultStyle
	}
	if userBlog.ThemeId == "" {
		userBlog.ThemePath = themeService.GetDefaultThemePath(userBlog.Style)
	} else {
		userBlog.ThemePath = themeService.GetThemePath(userBlog.UserId.Hex(), userBlog.ThemeId.Hex())
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

// 分页与排序
func (this *BlogService) UpdateUserBlogPaging(userId string, perPageSize int, sortField string, isAsc bool) (ok bool, msg string) {
	if ok, msg = Vds(map[string]string{"perPageSize": strconv.Itoa(perPageSize), "sortField": sortField}); !ok {
		return
	}
	ok = db.UpdateByQMap(db.UserBlogs, bson.M{"_id": bson.ObjectIdHex(userId)},
		bson.M{"PerPageSize": perPageSize, "SortField": sortField, "IsAsc": isAsc})
	return
}

func (this *BlogService) GetUserBlogBySubDomain(subDomain string) info.UserBlog {
	blogUser := info.UserBlog{}
	db.GetByQ(db.UserBlogs, bson.M{"SubDomain": subDomain}, &blogUser)
	this.fixUserBlog(&blogUser)
	return blogUser
}
func (this *BlogService) GetUserBlogByDomain(domain string) info.UserBlog {
	blogUser := info.UserBlog{}
	db.GetByQ(db.UserBlogs, bson.M{"Domain": domain}, &blogUser)
	this.fixUserBlog(&blogUser)
	return blogUser
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
func (this *BlogService) ListLikedUsers(noteId string, isAll bool) ([]info.UserAndBlog, bool) {
	// 默认前5
	pageSize := 5
	skipNum, sortFieldR := parsePageAndSort(1, pageSize, "CreatedTime", false)

	likes := []info.BlogLike{}
	query := bson.M{"NoteId": bson.ObjectIdHex(noteId)}
	q := db.BlogLikes.Find(query)

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
	userMap := userService.MapUserAndBlogByUserIds(userIds)

	users := make([]info.UserAndBlog, len(likes))
	for i, like := range likes {
		users[i] = userMap[like.UserId.Hex()]
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
	if !note.IsBlog /*|| note.UserId.Hex() == userId */ {
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
		NoteId:      bson.ObjectIdHex(noteId),
		UserId:      bson.ObjectIdHex(userId),
		Content:     content,
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
			this.sendEmail(note, comment2, userId, content)
		}()
	}

	return ok, comment
}

// 发送email
func (this *BlogService) sendEmail(note info.Note, comment info.BlogComment, userId, content string) {
	emailService.SendCommentEmail(note, comment, userId, content)
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

	if userId == configService.GetAdminUserId() || note.UserId.Hex() == userId || comment.UserId.Hex() == userId {
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
			bson.M{"$pull": bson.M{"LikeUserIds": userId}})
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
func (this *BlogService) ListComments(userId, noteId string, page, pageSize int) (info.Page, []info.BlogCommentPublic, map[string]info.UserAndBlog) {
	pageInfo := info.Page{CurPage: page}

	comments2 := []info.BlogComment{}

	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "CreatedTime", false)

	query := bson.M{"NoteId": bson.ObjectIdHex(noteId)}
	q := db.BlogComments.Find(query)

	// 总记录数
	count, _ := q.Count()
	q.Sort(sortFieldR).Skip(skipNum).Limit(pageSize).All(&comments2)

	if len(comments2) == 0 {
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

	note := noteService.GetNoteById(noteId)

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
	userMap := userService.MapUserAndBlogByUserIds(userIds)
	pageInfo = info.NewPage(page, pageSize, count, nil)

	return pageInfo, comments, userMap
}

// 举报
func (this *BlogService) Report(noteId, commentId, reason, userId string) bool {
	note := noteService.GetNoteById(noteId)
	if !note.IsBlog {
		return false
	}

	report := info.Report{ReportId: bson.NewObjectId(),
		NoteId:      bson.ObjectIdHex(noteId),
		UserId:      bson.ObjectIdHex(userId),
		Reason:      reason,
		CreatedTime: time.Now(),
	}
	if commentId != "" {
		report.CommentId = bson.ObjectIdHex(commentId)
	}
	return db.Insert(db.Reports, report)
}

//---------------
// 分类排序

// CateIds
func (this *BlogService) UpateCateIds(userId string, cateIds []string) bool {
	return db.UpdateByQField(db.UserBlogs, bson.M{"_id": bson.ObjectIdHex(userId)}, "CateIds", cateIds)
}

// 修改笔记本urlTitle
func (this *BlogService) UpateCateUrlTitle(userId string, cateId, urlTitle string) (ok bool, url string) {
	url = urlTitle
	/*
		// 先清空
		ok = db.UpdateByIdAndUserIdMap(db.Notebooks, cateId, userId, bson.M{
			"UrlTitle": "",
		})
	*/
	url = GetUrTitle(userId, urlTitle, "notebook", cateId)
	ok = db.UpdateByIdAndUserIdMap(db.Notebooks, cateId, userId, bson.M{
		"UrlTitle": url,
	})
	// 返回给前端的是decode
	url = decodeValue(url)
	return
}

// 修改笔记urlTitle
func (this *BlogService) UpateBlogUrlTitle(userId string, noteId, urlTitle string) (ok bool, url string) {
	url = urlTitle
	// 先清空
	ok = db.UpdateByIdAndUserIdMap(db.Notes, noteId, userId, bson.M{
		"UrlTitle": "",
	})
	url = GetUrTitle(userId, urlTitle, "note", noteId)
	ok = db.UpdateByIdAndUserIdMap(db.Notes, noteId, userId, bson.M{
		"UrlTitle": url,
	})
	// 返回给前端的是decode
	url = decodeValue(url)
	return
}

// 修改博客的图片, 描述, 摘要
func (this *BlogService) UpateBlogAbstract(userId string, noteId, imgSrc, desc, abstract string) (ok bool) {
	ok = db.UpdateByIdAndUserIdMap(db.Notes, noteId, userId, bson.M{
		"ImgSrc":         imgSrc,
		"Desc":           desc,
		"HasSelfDefined": true,
	})
	ok = db.UpdateByIdAndUserIdMap(db.NoteContents, noteId, userId, bson.M{
		"Abstract": abstract,
	})
	return ok
}

// 单页
func (this *BlogService) GetSingles(userId string) []map[string]string {
	userBlog := this.GetUserBlog(userId)
	singles := userBlog.Singles
	LogJ(singles)
	return singles
}
func (this *BlogService) GetSingle(singleId string) info.BlogSingle {
	page := info.BlogSingle{}
	db.Get(db.BlogSingles, singleId, &page)
	return page
}
func (this *BlogService) GetSingleByUserIdAndUrlTitle(userId, singleIdOrUrlTitle string) info.BlogSingle {
	page := info.BlogSingle{}
	if IsObjectId(singleIdOrUrlTitle) {
		db.Get(db.BlogSingles, singleIdOrUrlTitle, &page)
	} else {
		db.GetByQ(db.BlogSingles, bson.M{"UserId": bson.ObjectIdHex(userId), "UrlTitle": encodeValue(singleIdOrUrlTitle)}, &page)
	}
	return page
}

func (this *BlogService) updateBlogSingles(userId string, isDelete bool, isAdd bool, singleId, title, urlTitle string) (ok bool) {
	userBlog := this.GetUserBlog(userId)
	singles := userBlog.Singles
	if singles == nil {
		singles = []map[string]string{}
	}
	if isDelete || !isAdd { // 删除或更新, 需要找到相应的
		i := 0
		for _, p := range singles {
			if p["SingleId"] == singleId {
				break
			}
			i++
		}
		// 得到i
		// 找不到
		if i == len(singles) {
		} else {
			// 找到了, 如果是删除, 则删除
			if isDelete {
				singles = append(singles[:i], singles[i+1:]...)
			} else {
				// 是更新
				if title != "" {
					singles[i]["Title"] = title
				}
				if urlTitle != "" {
					singles[i]["UrlTitle"] = urlTitle
				}
			}
		}
	} else {
		// 是添加, 直接添加到最后
		singles = append(singles, map[string]string{"SingleId": singleId, "Title": title, "UrlTitle": urlTitle})
	}
	return db.UpdateByQField(db.UserBlogs, bson.M{"_id": bson.ObjectIdHex(userId)}, "Singles", singles)
}

// 删除页面
func (this *BlogService) DeleteSingle(userId, singleId string) (ok bool) {
	ok = db.DeleteByIdAndUserId(db.BlogSingles, singleId, userId)
	if ok {
		// 还要修改UserBlog中的Singles
		this.updateBlogSingles(userId, true, false, singleId, "", "")
	}
	return
}

// 修改urlTitle
func (this *BlogService) UpdateSingleUrlTitle(userId, singleId, urlTitle string) (ok bool, url string) {
	url = urlTitle
	/*
		// 先清空
		ok = db.UpdateByIdAndUserIdMap(db.BlogSingles, singleId, userId, bson.M{
			"UrlTitle": "",
		})
	*/
	url = GetUrTitle(userId, urlTitle, "single", singleId)
	ok = db.UpdateByIdAndUserIdMap(db.BlogSingles, singleId, userId, bson.M{
		"UrlTitle": url,
	})
	if ok {
		// 还要修改UserBlog中的Singles
		this.updateBlogSingles(userId, false, false, singleId, "", url)
	}
	// 返回给前端的是decode
	url = decodeValue(url)
	return
}

// 更新或添加
func (this *BlogService) AddOrUpdateSingle(userId, singleId, title, content string) (ok bool) {
	ok = false
	if singleId != "" {
		ok = db.UpdateByIdAndUserIdMap(db.BlogSingles, singleId, userId, bson.M{
			"Title":       title,
			"Content":     content,
			"UpdatedTime": time.Now(),
		})
		if ok {
			// 还要修改UserBlog中的Singles
			this.updateBlogSingles(userId, false, false, singleId, title, "")
		}
		return
	}
	// 添加
	page := info.BlogSingle{
		SingleId:    bson.NewObjectId(),
		UserId:      bson.ObjectIdHex(userId),
		Title:       title,
		Content:     content,
		UrlTitle:    GetUrTitle(userId, title, "single", singleId),
		CreatedTime: time.Now(),
	}
	page.UpdatedTime = page.CreatedTime
	ok = db.Insert(db.BlogSingles, page)

	// 还要修改UserBlog中的Singles
	this.updateBlogSingles(userId, false, true, page.SingleId.Hex(), title, page.UrlTitle)

	return
}

// 重新排序
func (this *BlogService) SortSingles(userId string, singleIds []string) (ok bool) {
	if singleIds == nil || len(singleIds) == 0 {
		return
	}
	userBlog := this.GetUserBlog(userId)
	singles := userBlog.Singles
	if singles == nil || len(singles) == 0 {
		return
	}

	singlesMap := map[string]map[string]string{}
	for _, page := range singles {
		singlesMap[page["SingleId"]] = page
	}

	singles2 := make([]map[string]string, len(singles))
	for i, singleId := range singleIds {
		singles2[i] = singlesMap[singleId]
	}

	return db.UpdateByQField(db.UserBlogs, bson.M{"_id": bson.ObjectIdHex(userId)}, "Singles", singles2)
}

// 得到用户的博客url
func (this *BlogService) GetUserBlogUrl(userBlog *info.UserBlog, username string) string {
	if userBlog != nil {
		if userBlog.Domain != "" && configService.AllowCustomDomain() {
			return configService.GetUserUrl(userBlog.Domain)
		} else if userBlog.SubDomain != "" {
			return configService.GetUserSubUrl(userBlog.SubDomain)
		}
		if username == "" {
			username = userBlog.UserId.Hex()
		}
	}
	return configService.GetBlogUrl() + "/" + username
}

// 得到所有url
func (this *BlogService) GetBlogUrls(userBlog *info.UserBlog, userInfo *info.User) info.BlogUrls {
	var indexUrl, postUrl, searchUrl, cateUrl, singleUrl, tagsUrl, archiveUrl, tagPostsUrl string
	if userBlog.Domain != "" && configService.AllowCustomDomain() { // http://demo.com
		// ok
		indexUrl = configService.GetUserUrl(userBlog.Domain)
		cateUrl = indexUrl + "/cate"     // /xxxxx
		postUrl = indexUrl + "/post"     // /xxxxx
		searchUrl = indexUrl + "/search" // /xxxxx
		singleUrl = indexUrl + "/single"
		archiveUrl = indexUrl + "/archives"
		tagsUrl = indexUrl + "/tags"
		tagPostsUrl = indexUrl + "/tag"
	} else if userBlog.SubDomain != "" { // demo.leanote.com
		indexUrl = configService.GetUserSubUrl(userBlog.SubDomain)
		cateUrl = indexUrl + "/cate"     // /xxxxx
		postUrl = indexUrl + "/post"     // /xxxxx
		searchUrl = indexUrl + "/search" // /xxxxx
		singleUrl = indexUrl + "/single"
		archiveUrl = indexUrl + "/archives"
		tagsUrl = indexUrl + "/tags"
		tagPostsUrl = indexUrl + "/tag"
	} else {
		// ok
		blogUrl := configService.GetBlogUrl() // blog.leanote.com
		userIdOrEmail := ""
		if userInfo.Username != "" {
			userIdOrEmail = userInfo.Username
		} else if userInfo.Email != "" {
			userIdOrEmail = userInfo.Email
		} else {
			userIdOrEmail = userInfo.UserId.Hex()
		}
		indexUrl = blogUrl + "/" + userIdOrEmail
		cateUrl = blogUrl + "/cate/" + userIdOrEmail        // /username/notebookId
		postUrl = blogUrl + "/post/" + userIdOrEmail        // /username/xxxxx
		searchUrl = blogUrl + "/search/" + userIdOrEmail    // blog.leanote.com/search/username
		singleUrl = blogUrl + "/single/" + userIdOrEmail    // blog.leanote.com/single/username/singleId
		archiveUrl = blogUrl + "/archives/" + userIdOrEmail // blog.leanote.com/archive/username
		tagsUrl = blogUrl + "/tags/" + userIdOrEmail
		tagPostsUrl = blogUrl + "/tag/" + userIdOrEmail // blog.leanote.com/archive/username
	}

	return info.BlogUrls{
		IndexUrl:    indexUrl,
		CateUrl:     cateUrl,
		SearchUrl:   searchUrl,
		SingleUrl:   singleUrl,
		PostUrl:     postUrl,
		ArchiveUrl:  archiveUrl,
		TagsUrl:     tagsUrl,
		TagPostsUrl: tagPostsUrl,
	}
}

// 转成post
func (this *BlogService) FixBlogs(blogs []info.BlogItem) []info.Post {
	blogs2 := make([]info.Post, len(blogs))
	for i, blog := range blogs {
		blogs2[i] = this.FixBlog(blog)
	}
	return blogs2
}
func (this *BlogService) FixBlog(blog info.BlogItem) info.Post {
	urlTitle := blog.UrlTitle
	if urlTitle == "" {
		urlTitle = blog.NoteId.Hex()
	}
	blog2 := info.Post{
		NoteId:      blog.NoteId.Hex(),
		Title:       blog.Title,
		UrlTitle:    urlTitle,
		ImgSrc:      blog.ImgSrc,
		CreatedTime: blog.CreatedTime,
		UpdatedTime: blog.UpdatedTime,
		PublicTime:  blog.PublicTime,
		Desc:        blog.Desc,
		Abstract:    blog.Abstract,
		Content:     blog.Content,
		Tags:        blog.Tags,
		CommentNum:  blog.CommentNum,
		ReadNum:     blog.ReadNum,
		LikeNum:     blog.LikeNum,
		IsMarkdown:  blog.IsMarkdown,
	}
	if blog2.Tags != nil && len(blog2.Tags) > 0 && blog2.Tags[0] != "" {
	} else {
		blog2.Tags = nil
	}
	return blog2
}

func (this *BlogService) FixNote(note info.Note) info.Post {
	if note.NoteId == "" {
		return info.Post{}
	}
	urlTitle := note.UrlTitle
	if urlTitle == "" {
		urlTitle = note.NoteId.Hex()
	}
	return info.Post{
		NoteId:      note.NoteId.Hex(),
		Title:       note.Title,
		ImgSrc:      note.ImgSrc,
		UrlTitle:    urlTitle,
		CreatedTime: note.CreatedTime,
		UpdatedTime: note.UpdatedTime,
		PublicTime:  note.PublicTime,
		Desc:        note.Desc,
		Tags:        note.Tags,
		CommentNum:  note.CommentNum,
		ReadNum:     note.ReadNum,
		LikeNum:     note.LikeNum,
		IsMarkdown:  note.IsMarkdown,
	}
}
