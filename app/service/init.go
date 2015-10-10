package service

import (
	"regexp"
	"strings"
	"net/url"
	"strconv"
	"gopkg.in/mgo.v2"
	"github.com/leanote/leanote/app/db"
	"gopkg.in/mgo.v2/bson"
)

// init service, for share service bettween services
// 初始化, 实例service
// 为了共享service

var notebookService, NotebookS *NotebookService
var noteService, NoteS *NoteService
var noteContentHistoryService, NoteContentHistoryS *NoteContentHistoryService
var trashService, TrashS *TrashService
var shareService, ShareS *ShareService
var userService, UserS *UserService
var groupService, GroupS *GroupService
var tagService, TagS *TagService
var blogService, BlogS *BlogService
var tokenService, TokenS *TokenService
var noteImageService, NoteImageS *NoteImageService
var fileService, FileS *FileService
var albumService, AlbumS *AlbumService
var attachService, AttachS *AttachService
var configService, ConfigS *ConfigService
var PwdS *PwdService
var SuggestionS *SuggestionService
var emailService, EmailS *EmailService
var AuthS *AuthService
var UpgradeS *UpgradeService
var SessionS, sessionService *SessionService
var ThemeS, themeService *ThemeService

// onAppStart调用
func InitService() {
	NotebookS = &NotebookService{}
	NoteS = &NoteService{}
	NoteContentHistoryS = &NoteContentHistoryService{}
	TrashS = &TrashService{}
	ShareS = &ShareService{}
	UserS = &UserService{}
	GroupS = &GroupService{}
	TagS = &TagService{}
	BlogS = &BlogService{}
	TokenS = &TokenService{}
	NoteImageS = &NoteImageService{}
	FileS = &FileService{}
	AlbumS = &AlbumService{}
	AttachS = &AttachService{}
	ConfigS = &ConfigService{}
	PwdS = &PwdService{}
	SuggestionS = &SuggestionService{}
	AuthS = &AuthService{}
	EmailS = NewEmailService()
	UpgradeS = &UpgradeService{}
	SessionS = &SessionService{}
	ThemeS = &ThemeService{}
	
	notebookService = NotebookS
	noteService = NoteS
	noteContentHistoryService = NoteContentHistoryS
	trashService = TrashS
	shareService = ShareS
	userService = UserS
	groupService = GroupS
	tagService = TagS
	blogService = BlogS
	tokenService = TokenS
	noteImageService = NoteImageS
	fileService = FileS
	albumService = AlbumS
	attachService = AttachS
	configService = ConfigS
	emailService = EmailS
	sessionService = SessionS
	themeService = ThemeS
}

//----------------
// service 公用方法

// 将name=val的val进行encoding
func decodeValue(val string) string {
	v, _ := url.ParseQuery("a=" + val)
	return v.Get("a")
}

func encodeValue(val string) string {
	if val == "" {
		return val
	}
	v := url.Values{}
	v.Set("", val)
	return v.Encode()[1:]
}

// 添加笔记时通过title得到urlTitle
func fixUrlTitle(urlTitle string) string {
	if urlTitle != "" {
		// 把特殊字段给替换掉
		//		str := `life "%&()+,/:;<>=?@\|`
		reg, _ := regexp.Compile("/|#|\\$|!|\\^|\\*|'| |\"|%|&|\\(|\\)|\\+|\\,|/|:|;|<|>|=|\\?|@|\\||\\\\")
		urlTitle = reg.ReplaceAllString(urlTitle, "-")
		urlTitle = strings.Trim(urlTitle, "-") // 左右单独的-去掉
		// 把空格替换成-
		//		urlTitle = strings.Replace(urlTitle, " ", "-", -1)
		for strings.Index(urlTitle, "--") >= 0 { // 防止出现连续的--
			urlTitle = strings.Replace(urlTitle, "--", "-", -1)
		}
		return encodeValue(urlTitle)
	}
	return urlTitle
}

func getUniqueUrlTitle(userId string, urlTitle string, types string, padding int) string {
	urlTitle2 := urlTitle

	// 判断urlTitle是不是过长, 过长则截断, 300
	// 不然生成index有问题
	// it will not index a single field with more than 1024 bytes.
	// If you're indexing a field that is 2.5MB, it's not really indexing it, it's being skipped.
	if len(urlTitle2) > 320 {
		urlTitle2 = urlTitle2[:300] // 为什么要少些, 因为怕无限循环, 因为把padding截了
	}

	if padding > 1 {
		urlTitle2 = urlTitle + "-" + strconv.Itoa(padding)
	}
	userIdO := bson.ObjectIdHex(userId)

	var collection *mgo.Collection
	if types == "note" {
		collection = db.Notes
	} else if types == "notebook" {
		collection = db.Notebooks
	} else if types == "single" {
		collection = db.BlogSingles
	}
	for db.Has(collection, bson.M{"UserId": userIdO, "UrlTitle": urlTitle2}) { // 用户下唯一
		padding++
		urlTitle2 = urlTitle + "-" + strconv.Itoa(padding)
	}

	return urlTitle2
}

// types == note,notebook,single
func GetUrTitle(userId string, title string, types string) string {
	urlTitle := strings.Trim(title, " ")
	if urlTitle == "" {
		urlTitle = "Untitled-" + userId
	}
	urlTitle = fixUrlTitle(urlTitle)
	return getUniqueUrlTitle(userId, urlTitle, types, 1)
}
