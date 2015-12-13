package controllers

import (
	"encoding/json"
	"github.com/leanote/leanote/app/info"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	//	. "github.com/leanote/leanote/app/lea"
	//	"io/ioutil"
	//	"fmt"
	"bytes"
	"math"
	"strconv"
	"strings"
)

// 公用Controller, 其它Controller继承它
type BaseController struct {
	*revel.Controller
}

func (c BaseController) GetUserId() string {
	if userId, ok := c.Session["UserId"]; ok {
		return userId
	}
	return ""
}

// 是否已登录
func (c BaseController) HasLogined() bool {
	return c.GetUserId() != ""
}

func (c BaseController) GetObjectUserId() bson.ObjectId {
	userId := c.GetUserId()
	if userId != "" {
		return bson.ObjectIdHex(userId)
	}
	return ""
}

func (c BaseController) GetEmail() string {
	if email, ok := c.Session["Email"]; ok {
		return email
	}
	return ""
}

func (c BaseController) GetUsername() string {
	if email, ok := c.Session["Username"]; ok {
		return email
	}
	return ""
}

// 得到用户信息
func (c BaseController) GetUserInfo() info.User {
	if userId, ok := c.Session["UserId"]; ok && userId != "" {
		return userService.GetUserInfo(userId)
		/*
			notebookWidth, _ := strconv.Atoi(c.Session["NotebookWidth"])
			noteListWidth, _ := strconv.Atoi(c.Session["NoteListWidth"])
			mdEditorWidth, _ := strconv.Atoi(c.Session["MdEditorWidth"])
			LogJ(c.Session)
			user := info.User{UserId: bson.ObjectIdHex(userId),
				Email: c.Session["Email"],
				Logo: c.Session["Logo"],
				Username: c.Session["Username"],
				UsernameRaw: c.Session["UsernameRaw"],
				Theme: c.Session["Theme"],
				NotebookWidth: notebookWidth,
				NoteListWidth: noteListWidth,
				MdEditorWidth: mdEditorWidth,
				}
			if c.Session["Verified"] == "1" {
				user.Verified = true
			}
			if c.Session["LeftIsMin"] == "1" {
				user.LeftIsMin = true
			}
			return user
		*/
	}
	return info.User{}
}

func (c BaseController) GetUserAndBlogUrl() info.UserAndBlogUrl {
	if userId, ok := c.Session["UserId"]; ok && userId != "" {
		return userService.GetUserAndBlogUrl(userId)
	}
	return info.UserAndBlogUrl{}
}

// 这里的session都是cookie中的, 与数据库session无关
func (c BaseController) GetSession(key string) string {
	v, ok := c.Session[key]
	if !ok {
		v = ""
	}
	return v
}
func (c BaseController) SetSession(userInfo info.User) {
	if userInfo.UserId.Hex() != "" {
		c.Session["UserId"] = userInfo.UserId.Hex()
		c.Session["Email"] = userInfo.Email
		c.Session["Username"] = userInfo.Username
		c.Session["UsernameRaw"] = userInfo.UsernameRaw
		c.Session["Theme"] = userInfo.Theme
		c.Session["Logo"] = userInfo.Logo

		c.Session["NotebookWidth"] = strconv.Itoa(userInfo.NotebookWidth)
		c.Session["NoteListWidth"] = strconv.Itoa(userInfo.NoteListWidth)

		if userInfo.Verified {
			c.Session["Verified"] = "1"
		} else {
			c.Session["Verified"] = "0"
		}

		if userInfo.LeftIsMin {
			c.Session["LeftIsMin"] = "1"
		} else {
			c.Session["LeftIsMin"] = "0"
		}
	}
}

func (c BaseController) ClearSession() {
	delete(c.Session, "UserId")
	delete(c.Session, "Email")
	delete(c.Session, "Username")
	delete(c.Session, "theme")
}

// 修改session
func (c BaseController) UpdateSession(key, value string) {
	c.Session[key] = value
}

// 返回json
func (c BaseController) Json(i interface{}) string {
	//	b, _ := json.MarshalIndent(i, "", "	")
	b, _ := json.Marshal(i)
	return string(b)
}

// 得到第几页
func (c BaseController) GetPage() int {
	page := 0
	c.Params.Bind(&page, "page")
	if page == 0 {
		return 1
	}
	return page
}

// 判断是否含有某参数
func (c BaseController) Has(key string) bool {
	if _, ok := c.Params.Values[key]; ok {
		return true
	}
	return false
}

/*
func (c Blog) GetPage(page, count int, list interface{}) info.Page {
	return info.Page{page, int(math.Ceil(float64(count)/float64(page))), list}
}
*/

func (c BaseController) GetTotalPage(page, count int) int {
	return int(math.Ceil(float64(count) / float64(page)))
}

//-------------
func (c BaseController) E404() revel.Result {
	c.RenderArgs["title"] = "404"
	return c.NotFound("", nil)
}

// 设置本地
func (c BaseController) SetLocale() string {
	locale := string(c.Request.Locale) // zh-CN
	lang := locale
	if strings.Contains(locale, "-") {
		pos := strings.Index(locale, "-")
		lang = locale[0:pos]
	}
	if lang != "zh" && lang != "en" && lang != "fr" && lang != "pt" {
		lang = "en"
	}
	c.RenderArgs["locale"] = lang
	c.RenderArgs["siteUrl"] = configService.GetSiteUrl()

	c.RenderArgs["blogUrl"] = configService.GetBlogUrl()
	c.RenderArgs["leaUrl"] = configService.GetLeaUrl()
	c.RenderArgs["noteUrl"] = configService.GetNoteUrl()

	return lang
}

// 设置userInfo
func (c BaseController) SetUserInfo() info.User {
	userInfo := c.GetUserInfo()
	c.RenderArgs["userInfo"] = userInfo
	if userInfo.Username == configService.GetAdminUsername() {
		c.RenderArgs["isAdmin"] = true
	}
	return userInfo
}

// life
// 返回解析后的字符串, 只是为了解析模板得到字符串
func (c BaseController) RenderTemplateStr(templatePath string) string {
	// Get the Template.
	// 返回 GoTemplate{tmpl, loader}
	template, err := revel.MainTemplateLoader.Template(templatePath)
	if err != nil {
	}

	tpl := &revel.RenderTemplateResult{
		Template:   template,
		RenderArgs: c.RenderArgs, // 把args给它
	}

	var buffer bytes.Buffer
	tpl.Template.Render(&buffer, c.RenderArgs)
	return buffer.String()
}

// json, result
// 为了msg
// msg-v1-v2-v3
func (c BaseController) RenderRe(re info.Re) revel.Result {
	oldMsg := re.Msg
	if re.Msg != "" {
		if strings.Contains(re.Msg, "-") {
			msgAndValues := strings.Split(re.Msg, "-")
			if len(msgAndValues) == 2 {
				re.Msg = c.Message(msgAndValues[0], msgAndValues[1])
			} else {
				others := msgAndValues[0:]
				a := make([]interface{}, len(others))
				for i, v := range others {
					a[i] = v
				}
				re.Msg = c.Message(msgAndValues[0], a...)
			}
		} else {
			re.Msg = c.Message(re.Msg)
		}
	}
	if strings.HasPrefix(re.Msg, "???") {
		re.Msg = oldMsg
	}
	return c.RenderJson(re)
}
