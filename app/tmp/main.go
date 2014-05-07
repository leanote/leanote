// GENERATED CODE - DO NOT EDIT
package main

import (
	"flag"
	"reflect"
	"github.com/revel/revel"
	_ "github.com/leanote/leanote/app"
	controllers "github.com/leanote/leanote/app/controllers"
	_ "github.com/leanote/leanote/app/db"
	info "github.com/leanote/leanote/app/info"
	_ "github.com/leanote/leanote/app/lea/binder"
	_ "github.com/leanote/leanote/app/lea/memcache"
	_ "github.com/leanote/leanote/app/service"
	controllers1 "github.com/revel/revel/modules/static/app/controllers"
	_ "github.com/revel/revel/modules/testrunner/app"
	controllers0 "github.com/revel/revel/modules/testrunner/app/controllers"
)

var (
	runMode    *string = flag.String("runMode", "", "Run mode.")
	port       *int    = flag.Int("port", 0, "By default, read from app.conf")
	importPath *string = flag.String("importPath", "", "Go Import Path for the app.")
	srcPath    *string = flag.String("srcPath", "", "Path to the source root.")

	// So compiler won't complain if the generated code doesn't reference reflect package...
	_ = reflect.Invalid
)

func main() {
	flag.Parse()
	revel.Init(*runMode, *importPath, *srcPath)
	revel.INFO.Println("Running revel server")
	
	revel.RegisterController((*controllers.BaseController)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "E404",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers0.TestRunner)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Index",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
					46: []string{ 
						"testSuites",
					},
				},
			},
			&revel.MethodType{
				Name: "Run",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "suite", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "test", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
					69: []string{ 
						"error",
					},
				},
			},
			&revel.MethodType{
				Name: "List",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers1.Static)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Serve",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "prefix", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "filepath", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "ServeModule",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "moduleName", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "prefix", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "filepath", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.File)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "UploadImage",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "renderHtml", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UploadBlogLogo",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UploadImageJson",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "renderHtml", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Notebook)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Index",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebook", Type: reflect.TypeOf((*info.Notebook)(nil)) },
					&revel.MethodArg{Name: "i", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "name", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DeleteNotebook",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "AddNotebook",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "title", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UpdateNotebookTitle",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "title", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.NoteContentHistory)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "ListHistories",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Note)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Index",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "ListNotes",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "ListTrashNotes",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetNoteContent",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UpdateNoteOrContent",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteOrContent", Type: reflect.TypeOf((*controllers.NoteOrContent)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DeleteNote",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "userId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "isShared", Type: reflect.TypeOf((*bool)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DeleteTrash",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "MoveNote",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "CopyNote",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "CopySharedNote",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "fromUserId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "SearchNote",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "key", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "SearchNoteByTags",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "tags", Type: reflect.TypeOf((*[]string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Html2Image",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Share)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "AddShareNote",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "emails", Type: reflect.TypeOf((*[]string)(nil)) },
					&revel.MethodArg{Name: "perm", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "AddShareNotebook",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "emails", Type: reflect.TypeOf((*[]string)(nil)) },
					&revel.MethodArg{Name: "perm", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "ListShareNotes",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "userId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "GetShareNoteContent",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "sharedUserId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "ListNoteShareUserInfo",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "ListNotebookShareUserInfo",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UpdateShareNotePerm",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "perm", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "toUserId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UpdateShareNotebookPerm",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "perm", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "toUserId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DeleteShareNote",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "toUserId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DeleteShareNotebook",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "toUserId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DeleteShareNoteBySharedUser",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "fromUserId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DeleteShareNotebookBySharedUser",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "fromUserId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DeleteUserShareNoteAndNotebook",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "fromUserId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.User)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "UpdateUsername",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "username", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UpdatePwd",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "oldPwd", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "pwd", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UpdateTheme",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "theme", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "SendRegisterEmail",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "content", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "toEmail", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "ReSendActiveEmail",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UpdateEmailSendActiveEmail",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "email", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UpdateEmail",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "token", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "ActiveEmail",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "token", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "AddAccount",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "email", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "pwd", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UpdateColumnWidth",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookWidth", Type: reflect.TypeOf((*int)(nil)) },
					&revel.MethodArg{Name: "noteListWidth", Type: reflect.TypeOf((*int)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "UpdateLeftIsMin",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "leftIsMin", Type: reflect.TypeOf((*bool)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Auth)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Login",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "email", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DoLogin",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "email", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "pwd", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Logout",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Demo",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Register",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DoRegister",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "email", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "pwd", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "FindPassword",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "DoFindPassword",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "email", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "FindPassword2",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "token", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "FindPasswordUpdate",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "token", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "pwd", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Blog)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "SetNote2Blog",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "isBlog", Type: reflect.TypeOf((*bool)(nil)) },
					&revel.MethodArg{Name: "isTop", Type: reflect.TypeOf((*bool)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "SetNotebook2Blog",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "isBlog", Type: reflect.TypeOf((*bool)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Index",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "userId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "notebookId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "View",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "noteId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "SearchBlog",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "userId", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "key", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Set",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "SetUserBlogBase",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "userBlog", Type: reflect.TypeOf((*info.UserBlogBase)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "SetUserBlogComment",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "userBlog", Type: reflect.TypeOf((*info.UserBlogComment)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "SetUserBlogStyle",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "userBlog", Type: reflect.TypeOf((*info.UserBlogStyle)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "AboutMe",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "userId", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Index)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Index",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Suggestion",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "addr", Type: reflect.TypeOf((*string)(nil)) },
					&revel.MethodArg{Name: "suggestion", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Mobile)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "Index",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			&revel.MethodType{
				Name: "Logout",
				Args: []*revel.MethodArg{ 
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.RegisterController((*controllers.Oauth)(nil),
		[]*revel.MethodType{
			&revel.MethodType{
				Name: "GithubCallback",
				Args: []*revel.MethodArg{ 
					&revel.MethodArg{Name: "code", Type: reflect.TypeOf((*string)(nil)) },
				},
				RenderArgNames: map[int][]string{ 
				},
			},
			
		})
	
	revel.DefaultValidationKeys = map[string]map[int]string{ 
	}
	revel.TestSuites = []interface{}{ 
	}

	revel.Run(*port)
}
