// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/revel/revel"


type tBaseController struct {}
var BaseController tBaseController


func (_ tBaseController) E404(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("BaseController.E404", args).Url
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (_ tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).Url
}

func (_ tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).Url
}

func (_ tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).Url
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).Url
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).Url
}


type tAuth struct {}
var Auth tAuth


func (_ tAuth) Login(
		email string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	return revel.MainRouter.Reverse("Auth.Login", args).Url
}

func (_ tAuth) DoLogin(
		email string,
		pwd string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	revel.Unbind(args, "pwd", pwd)
	return revel.MainRouter.Reverse("Auth.DoLogin", args).Url
}

func (_ tAuth) Logout(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Auth.Logout", args).Url
}

func (_ tAuth) Demo(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Auth.Demo", args).Url
}

func (_ tAuth) Register(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Auth.Register", args).Url
}

func (_ tAuth) DoRegister(
		email string,
		pwd string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	revel.Unbind(args, "pwd", pwd)
	return revel.MainRouter.Reverse("Auth.DoRegister", args).Url
}

func (_ tAuth) FindPassword(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Auth.FindPassword", args).Url
}

func (_ tAuth) DoFindPassword(
		email string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	return revel.MainRouter.Reverse("Auth.DoFindPassword", args).Url
}

func (_ tAuth) FindPassword2(
		token string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "token", token)
	return revel.MainRouter.Reverse("Auth.FindPassword2", args).Url
}

func (_ tAuth) FindPasswordUpdate(
		token string,
		pwd string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "token", token)
	revel.Unbind(args, "pwd", pwd)
	return revel.MainRouter.Reverse("Auth.FindPasswordUpdate", args).Url
}


type tBlog struct {}
var Blog tBlog


func (_ tBlog) SetNote2Blog(
		noteId string,
		isBlog bool,
		isTop bool,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	revel.Unbind(args, "isBlog", isBlog)
	revel.Unbind(args, "isTop", isTop)
	return revel.MainRouter.Reverse("Blog.SetNote2Blog", args).Url
}

func (_ tBlog) SetNotebook2Blog(
		notebookId string,
		isBlog bool,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	revel.Unbind(args, "isBlog", isBlog)
	return revel.MainRouter.Reverse("Blog.SetNotebook2Blog", args).Url
}

func (_ tBlog) Index(
		userId string,
		notebookId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "userId", userId)
	revel.Unbind(args, "notebookId", notebookId)
	return revel.MainRouter.Reverse("Blog.Index", args).Url
}

func (_ tBlog) View(
		noteId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	return revel.MainRouter.Reverse("Blog.View", args).Url
}

func (_ tBlog) SearchBlog(
		userId string,
		key string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "userId", userId)
	revel.Unbind(args, "key", key)
	return revel.MainRouter.Reverse("Blog.SearchBlog", args).Url
}

func (_ tBlog) Set(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Blog.Set", args).Url
}

func (_ tBlog) SetUserBlogBase(
		userBlog interface{},
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "userBlog", userBlog)
	return revel.MainRouter.Reverse("Blog.SetUserBlogBase", args).Url
}

func (_ tBlog) SetUserBlogComment(
		userBlog interface{},
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "userBlog", userBlog)
	return revel.MainRouter.Reverse("Blog.SetUserBlogComment", args).Url
}

func (_ tBlog) SetUserBlogStyle(
		userBlog interface{},
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "userBlog", userBlog)
	return revel.MainRouter.Reverse("Blog.SetUserBlogStyle", args).Url
}

func (_ tBlog) AboutMe(
		userId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "userId", userId)
	return revel.MainRouter.Reverse("Blog.AboutMe", args).Url
}


type tIndex struct {}
var Index tIndex


func (_ tIndex) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Index.Index", args).Url
}

func (_ tIndex) Suggestion(
		addr string,
		suggestion string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "addr", addr)
	revel.Unbind(args, "suggestion", suggestion)
	return revel.MainRouter.Reverse("Index.Suggestion", args).Url
}


type tMobile struct {}
var Mobile tMobile


func (_ tMobile) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Mobile.Index", args).Url
}

func (_ tMobile) Logout(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Mobile.Logout", args).Url
}


type tNoteContentHistory struct {}
var NoteContentHistory tNoteContentHistory


func (_ tNoteContentHistory) ListHistories(
		noteId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	return revel.MainRouter.Reverse("NoteContentHistory.ListHistories", args).Url
}


type tNote struct {}
var Note tNote


func (_ tNote) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Note.Index", args).Url
}

func (_ tNote) ListNotes(
		notebookId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	return revel.MainRouter.Reverse("Note.ListNotes", args).Url
}

func (_ tNote) ListTrashNotes(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Note.ListTrashNotes", args).Url
}

func (_ tNote) GetNoteContent(
		noteId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	return revel.MainRouter.Reverse("Note.GetNoteContent", args).Url
}

func (_ tNote) UpdateNoteOrContent(
		noteOrContent controllers.NoteOrContent,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteOrContent", noteOrContent)
	return revel.MainRouter.Reverse("Note.UpdateNoteOrContent", args).Url
}

func (_ tNote) DeleteNote(
		noteId string,
		userId string,
		isShared bool,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	revel.Unbind(args, "userId", userId)
	revel.Unbind(args, "isShared", isShared)
	return revel.MainRouter.Reverse("Note.DeleteNote", args).Url
}

func (_ tNote) DeleteTrash(
		noteId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	return revel.MainRouter.Reverse("Note.DeleteTrash", args).Url
}

func (_ tNote) MoveNote(
		noteId string,
		notebookId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	revel.Unbind(args, "notebookId", notebookId)
	return revel.MainRouter.Reverse("Note.MoveNote", args).Url
}

func (_ tNote) CopyNote(
		noteId string,
		notebookId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	revel.Unbind(args, "notebookId", notebookId)
	return revel.MainRouter.Reverse("Note.CopyNote", args).Url
}

func (_ tNote) CopySharedNote(
		noteId string,
		notebookId string,
		fromUserId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	revel.Unbind(args, "notebookId", notebookId)
	revel.Unbind(args, "fromUserId", fromUserId)
	return revel.MainRouter.Reverse("Note.CopySharedNote", args).Url
}

func (_ tNote) SearchNote(
		key string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "key", key)
	return revel.MainRouter.Reverse("Note.SearchNote", args).Url
}

func (_ tNote) SearchNoteByTags(
		tags []string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "tags", tags)
	return revel.MainRouter.Reverse("Note.SearchNoteByTags", args).Url
}

func (_ tNote) Html2Image(
		noteId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	return revel.MainRouter.Reverse("Note.Html2Image", args).Url
}


type tShare struct {}
var Share tShare


func (_ tShare) AddShareNote(
		noteId string,
		emails []string,
		perm int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	revel.Unbind(args, "emails", emails)
	revel.Unbind(args, "perm", perm)
	return revel.MainRouter.Reverse("Share.AddShareNote", args).Url
}

func (_ tShare) AddShareNotebook(
		notebookId string,
		emails []string,
		perm int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	revel.Unbind(args, "emails", emails)
	revel.Unbind(args, "perm", perm)
	return revel.MainRouter.Reverse("Share.AddShareNotebook", args).Url
}

func (_ tShare) ListShareNotes(
		notebookId string,
		userId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	revel.Unbind(args, "userId", userId)
	return revel.MainRouter.Reverse("Share.ListShareNotes", args).Url
}

func (_ tShare) GetShareNoteContent(
		noteId string,
		sharedUserId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	revel.Unbind(args, "sharedUserId", sharedUserId)
	return revel.MainRouter.Reverse("Share.GetShareNoteContent", args).Url
}

func (_ tShare) ListNoteShareUserInfo(
		noteId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	return revel.MainRouter.Reverse("Share.ListNoteShareUserInfo", args).Url
}

func (_ tShare) ListNotebookShareUserInfo(
		notebookId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	return revel.MainRouter.Reverse("Share.ListNotebookShareUserInfo", args).Url
}

func (_ tShare) UpdateShareNotePerm(
		noteId string,
		perm int,
		toUserId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	revel.Unbind(args, "perm", perm)
	revel.Unbind(args, "toUserId", toUserId)
	return revel.MainRouter.Reverse("Share.UpdateShareNotePerm", args).Url
}

func (_ tShare) UpdateShareNotebookPerm(
		notebookId string,
		perm int,
		toUserId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	revel.Unbind(args, "perm", perm)
	revel.Unbind(args, "toUserId", toUserId)
	return revel.MainRouter.Reverse("Share.UpdateShareNotebookPerm", args).Url
}

func (_ tShare) DeleteShareNote(
		noteId string,
		toUserId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	revel.Unbind(args, "toUserId", toUserId)
	return revel.MainRouter.Reverse("Share.DeleteShareNote", args).Url
}

func (_ tShare) DeleteShareNotebook(
		notebookId string,
		toUserId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	revel.Unbind(args, "toUserId", toUserId)
	return revel.MainRouter.Reverse("Share.DeleteShareNotebook", args).Url
}

func (_ tShare) DeleteShareNoteBySharedUser(
		noteId string,
		fromUserId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "noteId", noteId)
	revel.Unbind(args, "fromUserId", fromUserId)
	return revel.MainRouter.Reverse("Share.DeleteShareNoteBySharedUser", args).Url
}

func (_ tShare) DeleteShareNotebookBySharedUser(
		notebookId string,
		fromUserId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	revel.Unbind(args, "fromUserId", fromUserId)
	return revel.MainRouter.Reverse("Share.DeleteShareNotebookBySharedUser", args).Url
}

func (_ tShare) DeleteUserShareNoteAndNotebook(
		fromUserId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "fromUserId", fromUserId)
	return revel.MainRouter.Reverse("Share.DeleteUserShareNoteAndNotebook", args).Url
}


type tUser struct {}
var User tUser


func (_ tUser) UpdateUsername(
		username string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	return revel.MainRouter.Reverse("User.UpdateUsername", args).Url
}

func (_ tUser) UpdatePwd(
		oldPwd string,
		pwd string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "oldPwd", oldPwd)
	revel.Unbind(args, "pwd", pwd)
	return revel.MainRouter.Reverse("User.UpdatePwd", args).Url
}

func (_ tUser) UpdateTheme(
		theme string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "theme", theme)
	return revel.MainRouter.Reverse("User.UpdateTheme", args).Url
}

func (_ tUser) SendRegisterEmail(
		content string,
		toEmail string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "content", content)
	revel.Unbind(args, "toEmail", toEmail)
	return revel.MainRouter.Reverse("User.SendRegisterEmail", args).Url
}

func (_ tUser) ReSendActiveEmail(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("User.ReSendActiveEmail", args).Url
}

func (_ tUser) UpdateEmailSendActiveEmail(
		email string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	return revel.MainRouter.Reverse("User.UpdateEmailSendActiveEmail", args).Url
}

func (_ tUser) UpdateEmail(
		token string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "token", token)
	return revel.MainRouter.Reverse("User.UpdateEmail", args).Url
}

func (_ tUser) ActiveEmail(
		token string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "token", token)
	return revel.MainRouter.Reverse("User.ActiveEmail", args).Url
}

func (_ tUser) AddAccount(
		email string,
		pwd string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	revel.Unbind(args, "pwd", pwd)
	return revel.MainRouter.Reverse("User.AddAccount", args).Url
}

func (_ tUser) UpdateColumnWidth(
		notebookWidth int,
		noteListWidth int,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookWidth", notebookWidth)
	revel.Unbind(args, "noteListWidth", noteListWidth)
	return revel.MainRouter.Reverse("User.UpdateColumnWidth", args).Url
}

func (_ tUser) UpdateLeftIsMin(
		leftIsMin bool,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "leftIsMin", leftIsMin)
	return revel.MainRouter.Reverse("User.UpdateLeftIsMin", args).Url
}


type tFile struct {}
var File tFile


func (_ tFile) UploadImage(
		renderHtml string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "renderHtml", renderHtml)
	return revel.MainRouter.Reverse("File.UploadImage", args).Url
}

func (_ tFile) UploadBlogLogo(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("File.UploadBlogLogo", args).Url
}

func (_ tFile) UploadImageJson(
		renderHtml string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "renderHtml", renderHtml)
	return revel.MainRouter.Reverse("File.UploadImageJson", args).Url
}


type tNotebook struct {}
var Notebook tNotebook


func (_ tNotebook) Index(
		notebook interface{},
		i int,
		name string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebook", notebook)
	revel.Unbind(args, "i", i)
	revel.Unbind(args, "name", name)
	return revel.MainRouter.Reverse("Notebook.Index", args).Url
}

func (_ tNotebook) DeleteNotebook(
		notebookId string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	return revel.MainRouter.Reverse("Notebook.DeleteNotebook", args).Url
}

func (_ tNotebook) AddNotebook(
		notebookId string,
		title string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	revel.Unbind(args, "title", title)
	return revel.MainRouter.Reverse("Notebook.AddNotebook", args).Url
}

func (_ tNotebook) UpdateNotebookTitle(
		notebookId string,
		title string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "notebookId", notebookId)
	revel.Unbind(args, "title", title)
	return revel.MainRouter.Reverse("Notebook.UpdateNotebookTitle", args).Url
}


type tOauth struct {}
var Oauth tOauth


func (_ tOauth) GithubCallback(
		code string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "code", code)
	return revel.MainRouter.Reverse("Oauth.GithubCallback", args).Url
}


