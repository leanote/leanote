package controllers

import (
	"github.com/leanote/leanote/app/service"
	"github.com/leanote/leanote/app/info"
	"github.com/revel/revel"
	"strings"
)

// 该文件初始化所有service方法

var userService *service.UserService
var noteService *service.NoteService
var trashService *service.TrashService
var notebookService *service.NotebookService
var noteContentHistoryService *service.NoteContentHistoryService

var authService *service.AuthService
var shareService *service.ShareService
var blogService *service.BlogService
var tagService *service.TagService
var pwdService *service.PwdService
var tokenService *service.TokenService
var suggestionService *service.SuggestionService 

var albumService *service.AlbumService 
var fileService *service.FileService

var pageSize = 1000
var defaultSortField = "UpdatedTime"
var leanoteUserId = "52d26b4e99c37b609a000001"
var siteUrl = "http://leanote.com"
var openRegister = true

// 拦截器
// 不需要拦截的url
// Index 除了Note之外都不需要
var commonUrl = map[string]map[string]bool{"Index": map[string]bool{"Index": true, 
		"Login": true, 
		"DoLogin": true,
		"Logout": true,
		"Register": true,
		"DoRegister": true,
		"FindPasswword": true,
		"DoFindPassword": true,
		"FindPassword2": true,
		"FindPasswordUpdate": true,
		"Suggestion": true,
	},
	"Blog": map[string]bool{"Index": true,
		"View": true,
		"AboutMe": true,
		"SearchBlog": true,
		},
	// 用户的激活与修改邮箱都不需要登录, 通过链接地址
	"User": map[string]bool{"UpdateEmail": true,
		"ActiveEmail":true,
		},
	"Oauth": map[string]bool{"GithubCallback": true},
	"File": map[string]bool{"OutputImage": true, "OutputFile": true},
}
func needValidate(controller, method string) bool {
	// 在里面
	if v, ok := commonUrl[controller]; ok {
		// 在commonUrl里
		if _, ok2 := v[method]; ok2 {
			return false
		}
		return true
	} else {
		// controller不在这里的, 肯定要验证
		return true;
	}
}
func AuthInterceptor(c *revel.Controller) revel.Result {
	// 全部变成首字大写
	var controller = strings.Title(c.Name)
	var method = strings.Title(c.MethodName)
	
	
	// 是否需要验证?
	if !needValidate(controller, method) {
		return nil
	}
	
	
	// 验证是否已登录
	if userId, ok := c.Session["UserId"]; ok && userId != "" {
		return nil // 已登录
	}
	
	// 没有登录, 判断是否是ajax操作
	if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		re := info.NewRe()
		re.Msg = "NOTLOGIN"
		return c.RenderJson(re)
	}
	
	return c.Redirect("/login")
}

func init() {
	// interceptor
	// revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Index{}) // Index.Note自己校验
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Notebook{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Note{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Share{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &User{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &File{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &Blog{})
	revel.InterceptFunc(AuthInterceptor, revel.BEFORE, &NoteContentHistory{})
	
	// service 
	
	userService = &service.UserService{}
	noteService = &service.NoteService{}
	trashService = &service.TrashService{}
	notebookService = &service.NotebookService{}
	noteContentHistoryService = &service.NoteContentHistoryService{}
	authService = &service.AuthService{}
	shareService = &service.ShareService{}
	blogService = &service.BlogService{}
	tagService = &service.TagService{}
	pwdService = &service.PwdService{}
	tokenService = &service.TokenService{}
	suggestionService = &service.SuggestionService{}
	
	revel.OnAppStart(func() {
		leanoteUserId, _ = revel.Config.String("adminUsername")
		siteUrl, _ = revel.Config.String("site.url")
		openRegister, _ = revel.Config.Bool("register.open")
	})
}