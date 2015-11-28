package service

import (
	"encoding/json"
	"fmt"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/lea/archive"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

// 主题
type ThemeService struct {
}

var defaultStyle = "blog_default"
var elegantStyle = "blog_daqi"
var fixedStyle = "blog_left_fixed"

// admin用户的主题基路径
func (this *ThemeService) getDefaultThemeBasePath() string {
	return revel.BasePath + "/public/blog/themes"
}

// 默认主题路径
func (this *ThemeService) getDefaultThemePath(style string) string {
	if style == elegantStyle {
		return this.getDefaultThemeBasePath() + "/elegant"
	} else if style == fixedStyle {
		return this.getDefaultThemeBasePath() + "/nav_fixed"
	} else {
		return this.getDefaultThemeBasePath() + "/default"
	}
}

// blogService用
func (this *ThemeService) GetDefaultThemePath(style string) string {
	if style == elegantStyle {
		return "public/blog/themes/elegant"
	} else if style == fixedStyle {
		return "public/blog/themes/nav_fixed"
	} else {
		return "public/blog/themes/default"
	}
}

// 得到默认主题
// style是之前的值, 有3个值 blog_default, blog_daqi, blog_left_fixed
func (this *ThemeService) getDefaultTheme(style string) info.Theme {
	if style == elegantStyle {
		return info.Theme{
			IsDefault: true,
			Path:      "public/blog/themes/elegant",
			Name:      "leanote elegant",
			Author:    "leanote",
			AuthorUrl: "http://leanote.com",
			Version:   "1.0",
		}
	} else if style == fixedStyle {
		return info.Theme{
			IsDefault: true,
			Path:      "public/blog/themes/nav_fixed",
			Name:      "leanote nav fixed",
			Author:    "leanote",
			AuthorUrl: "http://leanote.com",
			Version:   "1.0",
		}
	} else { // blog default
		return info.Theme{
			IsDefault: true,
			Path:      "public/blog/themes/default",
			Name:      "leanote default",
			Author:    "leanote",
			AuthorUrl: "http://leanote.com",
			Version:   "1.0",
		}
	}
}

// 用户的主题路径设置
func (this *ThemeService) getUserThemeBasePath(userId string) string {
	return revel.BasePath + "/public/upload/" + Digest3(userId) + "/" + userId + "/themes"
}
func (this *ThemeService) getUserThemePath(userId, themeId string) string {
	return this.getUserThemeBasePath(userId) + "/" + themeId
}
func (this *ThemeService) getUserThemePath2(userId, themeId string) string {
	return "public/upload/" + Digest3(userId) + "/" + userId + "/themes/" + themeId
}

// 新建主题
// 复制默认主题到文件夹下
func (this *ThemeService) CopyDefaultTheme(userBlog info.UserBlog) (ok bool, themeId string) {
	newThemeId := bson.NewObjectId()
	themeId = newThemeId.Hex()
	userId := userBlog.UserId.Hex()
	themePath := this.getUserThemePath(userId, themeId)
	err := os.MkdirAll(themePath, 0755)
	if err != nil {
		return
	}
	// 复制默认主题
	defaultThemePath := this.getDefaultThemePath(userBlog.Style)
	err = CopyDir(defaultThemePath, themePath)
	if err != nil {
		return
	}

	// 保存到数据库中
	theme, _ := this.getThemeConfig(themePath)
	theme.ThemeId = newThemeId
	theme.Path = this.getUserThemePath2(userId, themeId)
	theme.CreatedTime = time.Now()
	theme.UpdatedTime = theme.CreatedTime
	theme.UserId = bson.ObjectIdHex(userId)

	ok = db.Insert(db.Themes, theme)
	return ok, themeId
}

// 第一次新建主题
// 设为active true
func (this *ThemeService) NewThemeForFirst(userBlog info.UserBlog) (ok bool, themeId string) {
	ok, themeId = this.CopyDefaultTheme(userBlog)
	this.ActiveTheme(userBlog.UserId.Hex(), themeId)
	// db.UpdateByQField(db.Themes, bson.M{"_id": bson.ObjectIdHex(themeId)}, "IsActive", true)
	return
}

// 新建主题, 判断是否有主题了
func (this *ThemeService) NewTheme(userId string) (ok bool, themeId string) {
	userBlog := blogService.GetUserBlog(userId)
	// 如果还没有主题, 那先复制旧的主题
	if userBlog.ThemeId == "" {
		themeService.NewThemeForFirst(userBlog)
	}
	// 再copy一个默认主题
	userBlog.Style = "defaultStyle"
	ok, themeId = this.CopyDefaultTheme(userBlog)
	return
}

// 将字符串转成Theme配置
func (this *ThemeService) parseConfig(configStr string) (theme info.Theme, err error) {
	theme = info.Theme{}
	// 除去/**/注释
	reg, _ := regexp.Compile("/\\*[\\s\\S]*?\\*/")
	configStr = reg.ReplaceAllString(configStr, "")
	// 转成map
	config := map[string]interface{}{}
	err = json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return
	}
	// 没有错, 则将Name, Version, Author, AuthorUrl
	Name := config["Name"]
	if Name != nil {
		theme.Name = Name.(string)
	}
	Version := config["Version"]
	if Version != nil {
		theme.Version = Version.(string)
	}
	Author := config["Author"]
	if Author != nil {
		theme.Author = Author.(string)
	}
	AuthorUrl := config["AuthorUrl"]
	if AuthorUrl != nil {
		theme.AuthorUrl = AuthorUrl.(string)
	}
	theme.Info = config

	return
}

// 读取theme.json得到值
func (this *ThemeService) getThemeConfig(themePath string) (theme info.Theme, err error) {
	theme = info.Theme{}
	configStr := GetFileStrContent(themePath + "/theme.json")
	theme, err = this.parseConfig(configStr)
	return
}

func (this *ThemeService) GetTheme(userId, themeId string) info.Theme {
	theme := info.Theme{}
	db.GetByQ(db.Themes, bson.M{"_id": bson.ObjectIdHex(themeId), "UserId": bson.ObjectIdHex(userId)}, &theme)
	return theme
}
func (this *ThemeService) GetThemeById(themeId string) info.Theme {
	theme := info.Theme{}
	db.GetByQ(db.Themes, bson.M{"_id": bson.ObjectIdHex(themeId)}, &theme)
	return theme
}

// 得到主题信息, 为了给博客用
func (this *ThemeService) GetThemeInfo(themeId, style string) map[string]interface{} {
	q := bson.M{}
	if themeId == "" {
		if style == "" {
			style = defaultStyle
		}
		q["Style"] = style
		q["IsDefault"] = true
	} else {
		q["_id"] = bson.ObjectIdHex(themeId)
	}
	theme := info.Theme{}
	db.GetByQ(db.Themes, q, &theme)
	return theme.Info
}

// 得到用户的主题
// 若用户没有主题, 则至少有一个默认主题
// 第一个返回值是当前active
func (this *ThemeService) GetUserThemes(userId string) (theme info.Theme, themes []info.Theme) {
	theme = info.Theme{}
	themes = []info.Theme{}

	//	db.ListByQ(db.Themes, bson.M{"UserId": bson.ObjectIdHex(userId)}, &themes)

	// 创建时间逆序
	query := bson.M{"UserId": bson.ObjectIdHex(userId)}
	q := db.Themes.Find(query)
	q.Sort("-CreatedTime").All(&themes)
	if len(themes) == 0 {
		userBlog := blogService.GetUserBlog(userId)
		theme = this.getDefaultTheme(userBlog.Style)
	} else {
		var has = false
		// 第一个是active的主题
		themes2 := make([]info.Theme, len(themes))
		i := 0
		for _, t := range themes {
			if t.IsActive {
				theme = t
			} else {
				has = true
				themes2[i] = t
				i++
			}
		}
		if has {
			themes = themes2
		} else {
			themes = nil
		}
	}
	return
}

// 得到默认主题供选择
func (this *ThemeService) GetDefaultThemes() (themes []info.Theme) {
	themes = []info.Theme{}
	db.ListByQ(db.Themes, bson.M{"IsDefault": true}, &themes)
	return
}

// 得到模板内容
func (this *ThemeService) GetTplContent(userId, themeId, filename string) string {
	path := this.GetThemeAbsolutePath(userId, themeId) + "/" + filename
	return GetFileStrContent(path)
}

// 得到主题的绝对路径
func (this *ThemeService) GetThemeAbsolutePath(userId, themeId string) string {
	theme := info.Theme{}
	db.GetByQWithFields(db.Themes, bson.M{"_id": bson.ObjectIdHex(themeId), "UserId": bson.ObjectIdHex(userId)}, []string{"Path"}, &theme)
	if theme.Path != "" {
		return revel.BasePath + "/" + theme.Path
	}
	return ""
}
func (this *ThemeService) GetThemePath(userId, themeId string) string {
	theme := info.Theme{}
	db.GetByQWithFields(db.Themes, bson.M{"_id": bson.ObjectIdHex(themeId), "UserId": bson.ObjectIdHex(userId)}, []string{"Path"}, &theme)
	if theme.Path != "" {
		return theme.Path
	}
	return ""
}

// 更新模板内容
func (this *ThemeService) UpdateTplContent(userId, themeId, filename, content string) (ok bool, msg string) {
	basePath := this.GetThemeAbsolutePath(userId, themeId)
	path := basePath + "/" + filename
	if strings.Contains(filename, ".html") {
		Log(">>")
		if ok, msg = this.ValidateTheme(basePath, filename, content); ok {
			// 模板
			if ok, msg = this.mustTpl(filename, content); ok {
				ok = PutFileStrContent(path, content)
			}
		}
		return
	} else if filename == "theme.json" {
		// 主题配置, 判断是否是正确的json
		theme, err := this.parseConfig(content)
		if err != nil {
			return false, fmt.Sprintf("%v", err)
		}
		// 正确, 更新theme信息
		ok = db.UpdateByQMap(db.Themes, bson.M{"_id": bson.ObjectIdHex(themeId), "UserId": bson.ObjectIdHex(userId)},
			bson.M{
				"Name":      theme.Name,
				"Version":   theme.Version,
				"Author":    theme.Author,
				"AuthorUrl": theme.AuthorUrl,
				"Info":      theme.Info,
			})
		if ok {
			ok = PutFileStrContent(path, content)
		}
		return
	}
	ok = PutFileStrContent(path, content)
	return
}

func (this *ThemeService) DeleteTpl(userId, themeId, filename string) (ok bool) {
	path := this.GetThemeAbsolutePath(userId, themeId) + "/" + filename
	ok = DeleteFile(path)
	return
}

// 判断是否有语法错误
func (this *ThemeService) mustTpl(filename, content string) (ok bool, msg string) {
	ok = true
	defer func() {
		if err := recover(); err != nil {
			ok = false
			Log(err)
			msg = fmt.Sprintf("%v", err)
		}
	}()
	template.Must(template.New(filename).Funcs(revel.TemplateFuncs).Parse(content))
	return
}

/////////

// 使用主题
func (this *ThemeService) ActiveTheme(userId, themeId string) (ok bool) {
	if db.Has(db.Themes, bson.M{"_id": bson.ObjectIdHex(themeId), "UserId": bson.ObjectIdHex(userId)}) {
		// 之前的设为false
		db.UpdateByQField(db.Themes, bson.M{"UserId": bson.ObjectIdHex(userId), "IsActive": true}, "IsActive", false)
		// 现在的设为true
		db.UpdateByQField(db.Themes, bson.M{"_id": bson.ObjectIdHex(themeId)}, "IsActive", true)

		// UserBlog ThemeId
		db.UpdateByQField(db.UserBlogs, bson.M{"_id": bson.ObjectIdHex(userId)}, "ThemeId", bson.ObjectIdHex(themeId))
		return true
	}
	return false
}

// 删除主题
func (this *ThemeService) DeleteTheme(userId, themeId string) (ok bool) {
	return db.Delete(db.Themes, bson.M{"_id": bson.ObjectIdHex(themeId), "UserId": bson.ObjectIdHex(userId), "IsActive": false})
}

// 公开主题, 只有管理员才有权限, 之前没公开的变成公开
func (this *ThemeService) PublicTheme(userId, themeId string) (ok bool) {
	// 是否是管理员?
	userInfo := userService.GetUserInfo(userId)
	if userInfo.Username == configService.GetAdminUsername() {
		theme := this.GetThemeById(themeId)
		return db.UpdateByQField(db.Themes, bson.M{"UserId": bson.ObjectIdHex(userId), "_id": bson.ObjectIdHex(themeId)}, "IsDefault", !theme.IsDefault)
	}
	return false
}

// 导出主题
func (this *ThemeService) ExportTheme(userId, themeId string) (ok bool, path string) {
	theme := this.GetThemeById(themeId)
	// 打包
	// 验证路径, 别把整个项目打包了
	Log(theme.Path)
	if theme.Path == "" ||
		(!strings.HasPrefix(theme.Path, "public/upload") &&
			!strings.HasPrefix(theme.Path, "public/blog/themes")) ||
		strings.Contains(theme.Path, "..") {
		return
	}

	sourcePath := revel.BasePath + "/" + theme.Path
	targetPath := revel.BasePath + "/public/upload/" + userId + "/tmp"
	err := os.MkdirAll(targetPath, 0755)
	if err != nil {
		Log(err)
		return
	}
	targetName := targetPath + "/" + theme.Name + ".zip"
	Log(sourcePath)
	Log(targetName)
	ok = archive.Zip(sourcePath, targetName)
	if !ok {
		return
	}

	return true, targetName
}

// 导入主题
// path == /llllllll/..../public/upload/.../aa.zip, 绝对路径
func (this *ThemeService) ImportTheme(userId, path string) (ok bool, msg string) {
	themeIdO := bson.NewObjectId()
	themeId := themeIdO.Hex()
	targetPath := this.getUserThemePath(userId, themeId) // revel.BasePath + "/public/upload/" + userId + "/themes/" + themeId

	err := os.MkdirAll(targetPath, 0755)
	if err != nil {
		msg = "error"
		return
	}
	if ok, msg = archive.Unzip(path, targetPath); !ok {
		DeleteFile(targetPath)
		Log("oh no")
		return
	}

	// 主题验证
	if ok, msg = this.ValidateTheme(targetPath, "", ""); !ok {
		DeleteFile(targetPath)
		return
	}
	// 解压成功, 那么新建之
	// 保存到数据库中
	theme, _ := this.getThemeConfig(targetPath)
	if theme.Name == "" {
		ok = false
		DeleteFile(targetPath)
		msg = "解析错误"
		return
	}
	theme.ThemeId = themeIdO
	theme.Path = this.getUserThemePath2(userId, themeId)
	theme.CreatedTime = time.Now()
	theme.UpdatedTime = theme.CreatedTime
	theme.UserId = bson.ObjectIdHex(userId)

	ok = db.Insert(db.Themes, theme)
	if !ok {
		DeleteFile(targetPath)
	}
	DeleteFile(path)
	return
}

// 升级用
// public/

func (this *ThemeService) UpgradeThemeBeta2() (ok bool) {
	adminUserId := configService.GetAdminUserId()
	this.upgradeThemeBeta2(adminUserId, defaultStyle, true)
	this.upgradeThemeBeta2(adminUserId, elegantStyle, false)
	this.upgradeThemeBeta2(adminUserId, fixedStyle, false)
	return true
}
func (this *ThemeService) upgradeThemeBeta2(userId, style string, isActive bool) (ok bool) {
	// 解压成功, 那么新建之
	// 保存到数据库中
	targetPath := this.GetDefaultThemePath(style)
	theme, _ := this.getThemeConfig(revel.BasePath + "/" + targetPath)
	if theme.Name == "" {
		ok = false
		return
	}
	themeIdO := bson.NewObjectId()
	theme.ThemeId = themeIdO
	theme.Path = targetPath // public
	theme.CreatedTime = time.Now()
	theme.UpdatedTime = theme.CreatedTime
	theme.UserId = bson.ObjectIdHex(userId)
	theme.IsActive = isActive
	theme.IsDefault = true
	theme.Style = style
	ok = db.Insert(db.Themes, theme)
	return ok
}

// 安装主题
// 得到该主题路径
func (this *ThemeService) InstallTheme(userId, themeId string) (ok bool) {
	theme := this.GetThemeById(themeId)
	// 不是默认主题, 即不是admin用户的主题, 不能乱安装
	if !theme.IsDefault {
		return false
	}

	// 用户之前是否有主题?
	userBlog := blogService.GetUserBlog(userId)
	if userBlog.ThemeId == "" {
		this.NewThemeForFirst(userBlog)
	}

	// 生成新主题
	newThemeId := bson.NewObjectId()
	themeId = newThemeId.Hex()
	themePath := this.getUserThemePath(userId, themeId)
	err := os.MkdirAll(themePath, 0755)
	if err != nil {
		return
	}
	// 复制默认主题
	sourceThemePath := revel.BasePath + "/" + theme.Path
	err = CopyDir(sourceThemePath, themePath)
	if err != nil {
		return
	}

	// 保存到数据库中
	theme, _ = this.getThemeConfig(themePath)
	theme.ThemeId = newThemeId
	theme.Path = this.getUserThemePath2(userId, themeId)
	theme.CreatedTime = time.Now()
	theme.UpdatedTime = theme.CreatedTime
	theme.UserId = bson.ObjectIdHex(userId)

	ok = db.Insert(db.Themes, theme)

	// 激活之
	this.ActiveTheme(userId, themeId)

	return ok
}

// 验证主题是否全法, 存在循环引用?
// filename, newContent 表示在修改模板时要判断模板修改时是否有错误
func (this *ThemeService) ValidateTheme(path string, filename, newContent string) (ok bool, msg string) {
	Log("theme Path")
	Log(path)
	// 建立一个有向图
	// 将该path下的所有文件提出, 得到文件的引用情况
	files := ListDir(path)
	LogJ(files)
	size := len(files)
	if size > 100 {
		ok = false
		msg = "tooManyFiles"
		return
	}
	/*
		111111111
		111000000
	*/
	vector := make([][]int, size)
	for i := 0; i < size; i++ {
		vector[i] = make([]int, size)
	}
	fileIndexMap := map[string]int{}   // fileName => index
	fileContent := map[string]string{} // fileName => content
	index := 0
	// 得到文件内容, 和建立索引, 每个文件都有一个index, 对应数组位置
	for _, t := range files {
		if !strings.Contains(t, ".html") {
			continue
		}
		if t != filename {
			fileBytes, err := ioutil.ReadFile(path + "/" + t)
			if err != nil {
				continue
			}

			fileIndexMap[t] = index
			// html内容
			fileStr := string(fileBytes)
			fileContent[t] = fileStr
		} else {
			fileIndexMap[t] = index
			fileContent[t] = newContent
		}
		index++
	}
	// 分析文件内容, 建立有向图
	reg, _ := regexp.Compile("{{ *template \"(.+?\\.html)\".*}}")
	for filename, content := range fileContent {
		thisIndex := fileIndexMap[filename]
		finds := reg.FindAllStringSubmatch(content, -1) // 子匹配
		LogJ(finds)
		//		Log(content)
		if finds != nil && len(finds) > 0 {
			for _, includes := range finds {
				include := includes[1]
				includeIndex, has := fileIndexMap[include]
				Log(includeIndex)
				Log("??")
				Log(has)
				if has {
					vector[thisIndex][includeIndex] = 1
				}
			}
		}
	}

	LogJ(vector)
	LogJ(fileIndexMap)
	// 建立图后, 判断是否有环
	if this.hasRound(vector, index) {
		ok = false
		msg = "themeValidHasRoundInclude"
	} else {
		ok = true
	}
	return
}

// 检测有向图是否有环, DFS
func (this *ThemeService) hasRound(vector [][]int, size int) (ok bool) {
	for i := 0; i < size; i++ {
		visited := make([]int, size)
		if this.hasRoundEach(vector, i, size, visited) {
			return true
		}
	}
	return false
}

// 从每个节点出发, 判断是否有环
func (this *ThemeService) hasRoundEach(vector [][]int, index int, size int, visited []int) (ok bool) {
	if visited[index] > 0 {
		return true
	}
	visited[index] = 1
	// 遍历它的孩子
	for i := 0; i < size; i++ {
		if vector[index][i] > 0 {
			return this.hasRoundEach(vector, i, size, visited)
		}
	}
	visited[index] = 0
	return false
}
