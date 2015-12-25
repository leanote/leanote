package service

import (
	"fmt"
	"github.com/leanote/leanote/app/db"
	"github.com/leanote/leanote/app/info"
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 配置服务
// 只是全局的, 用户的配置没有
type ConfigService struct {
	adminUserId   string
	siteUrl       string
	adminUsername string
	// 全局的
	GlobalAllConfigs    map[string]interface{}
	GlobalStringConfigs map[string]string
	GlobalArrayConfigs  map[string][]string
	GlobalMapConfigs    map[string]map[string]string
	GlobalArrMapConfigs map[string][]map[string]string
}

// appStart时 将全局的配置从数据库中得到作为全局
func (this *ConfigService) InitGlobalConfigs() bool {
	this.GlobalAllConfigs = map[string]interface{}{}
	this.GlobalStringConfigs = map[string]string{}
	this.GlobalArrayConfigs = map[string][]string{}
	this.GlobalMapConfigs = map[string]map[string]string{}
	this.GlobalArrMapConfigs = map[string][]map[string]string{}

	this.adminUsername, _ = revel.Config.String("adminUsername")
	if this.adminUsername == "" {
		this.adminUsername = "admin"
	}
	this.siteUrl, _ = revel.Config.String("site.url")

	userInfo := userService.GetUserInfoByAny(this.adminUsername)
	if userInfo.UserId == "" {
		return false
	}
	this.adminUserId = userInfo.UserId.Hex()

	configs := []info.Config{}
	db.ListByQ(db.Configs, bson.M{"UserId": userInfo.UserId}, &configs)

	for _, config := range configs {
		if config.IsArr {
			this.GlobalArrayConfigs[config.Key] = config.ValueArr
			this.GlobalAllConfigs[config.Key] = config.ValueArr
		} else if config.IsMap {
			this.GlobalMapConfigs[config.Key] = config.ValueMap
			this.GlobalAllConfigs[config.Key] = config.ValueMap
		} else if config.IsArrMap {
			this.GlobalArrMapConfigs[config.Key] = config.ValueArrMap
			this.GlobalAllConfigs[config.Key] = config.ValueArrMap
		} else {
			this.GlobalStringConfigs[config.Key] = config.ValueStr
			this.GlobalAllConfigs[config.Key] = config.ValueStr
		}
	}

	return true
}

func (this *ConfigService) GetSiteUrl() string {
	return this.siteUrl
}
func (this *ConfigService) GetAdminUsername() string {
	return this.adminUsername
}
func (this *ConfigService) GetAdminUserId() string {
	return this.adminUserId
}

// 通用方法
func (this *ConfigService) updateGlobalConfig(userId, key string, value interface{}, isArr, isMap, isArrMap bool) bool {
	// 判断是否存在
	if _, ok := this.GlobalAllConfigs[key]; !ok {
		// 需要添加
		config := info.Config{ConfigId: bson.NewObjectId(),
			UserId:      bson.ObjectIdHex(userId),
			Key:         key,
			IsArr:       isArr,
			IsMap:       isMap,
			IsArrMap:    isArrMap,
			UpdatedTime: time.Now(),
		}
		if isArr {
			v, _ := value.([]string)
			config.ValueArr = v
			this.GlobalArrayConfigs[key] = v
		} else if isMap {
			v, _ := value.(map[string]string)
			config.ValueMap = v
			this.GlobalMapConfigs[key] = v
		} else if isArrMap {
			v, _ := value.([]map[string]string)
			config.ValueArrMap = v
			this.GlobalArrMapConfigs[key] = v
		} else {
			v, _ := value.(string)
			config.ValueStr = v
			this.GlobalStringConfigs[key] = v
		}
		return db.Insert(db.Configs, config)
	} else {
		i := bson.M{"UpdatedTime": time.Now()}
		this.GlobalAllConfigs[key] = value
		if isArr {
			v, _ := value.([]string)
			i["ValueArr"] = v
			this.GlobalArrayConfigs[key] = v
		} else if isMap {
			v, _ := value.(map[string]string)
			i["ValueMap"] = v
			this.GlobalMapConfigs[key] = v
		} else if isArrMap {
			v, _ := value.([]map[string]string)
			i["ValueArrMap"] = v
			this.GlobalArrMapConfigs[key] = v
		} else {
			v, _ := value.(string)
			i["ValueStr"] = v
			this.GlobalStringConfigs[key] = v
		}
		return db.UpdateByQMap(db.Configs, bson.M{"UserId": bson.ObjectIdHex(userId), "Key": key}, i)
	}
}

// 更新用户配置
func (this *ConfigService) UpdateGlobalStringConfig(userId, key string, value string) bool {
	return this.updateGlobalConfig(userId, key, value, false, false, false)
}
func (this *ConfigService) UpdateGlobalArrayConfig(userId, key string, value []string) bool {
	return this.updateGlobalConfig(userId, key, value, true, false, false)
}
func (this *ConfigService) UpdateGlobalMapConfig(userId, key string, value map[string]string) bool {
	return this.updateGlobalConfig(userId, key, value, false, true, false)
}
func (this *ConfigService) UpdateGlobalArrMapConfig(userId, key string, value []map[string]string) bool {
	return this.updateGlobalConfig(userId, key, value, false, false, true)
}

// 获取全局配置, 博客平台使用
func (this *ConfigService) GetGlobalStringConfig(key string) string {
	return this.GlobalStringConfigs[key]
}
func (this *ConfigService) GetGlobalArrayConfig(key string) []string {
	arr := this.GlobalArrayConfigs[key]
	if arr == nil {
		return []string{}
	}
	return arr
}
func (this *ConfigService) GetGlobalMapConfig(key string) map[string]string {
	m := this.GlobalMapConfigs[key]
	if m == nil {
		return map[string]string{}
	}
	return m
}
func (this *ConfigService) GetGlobalArrMapConfig(key string) []map[string]string {
	m := this.GlobalArrMapConfigs[key]
	if m == nil {
		return []map[string]string{}
	}
	return m
}

func (this *ConfigService) IsOpenRegister() bool {
	return this.GetGlobalStringConfig("openRegister") != ""
}

//-------
// 修改共享笔记的配置
func (this *ConfigService) UpdateShareNoteConfig(registerSharedUserId string,
	registerSharedNotebookPerms, registerSharedNotePerms []int,
	registerSharedNotebookIds, registerSharedNoteIds, registerCopyNoteIds []string) (ok bool, msg string) {

	defer func() {
		if err := recover(); err != nil {
			ok = false
			msg = fmt.Sprint(err)
		}
	}()

	// 用户是否存在?
	if registerSharedUserId == "" {
		ok = true
		msg = "share userId is blank, So it share nothing to register"
		this.UpdateGlobalStringConfig(this.adminUserId, "registerSharedUserId", "")
		return
	} else {
		user := userService.GetUserInfo(registerSharedUserId)
		if user.UserId == "" {
			ok = false
			msg = "no such user: " + registerSharedUserId
			return
		} else {
			this.UpdateGlobalStringConfig(this.adminUserId, "registerSharedUserId", registerSharedUserId)
		}
	}

	notebooks := []map[string]string{}
	// 共享笔记本
	if len(registerSharedNotebookIds) > 0 {
		for i := 0; i < len(registerSharedNotebookIds); i++ {
			// 判断笔记本是否存在
			notebookId := registerSharedNotebookIds[i]
			if notebookId == "" {
				continue
			}
			notebook := notebookService.GetNotebook(notebookId, registerSharedUserId)
			if notebook.NotebookId == "" {
				ok = false
				msg = "The user has no such notebook: " + notebookId
				return
			} else {
				perm := "0"
				if registerSharedNotebookPerms[i] == 1 {
					perm = "1"
				}
				notebooks = append(notebooks, map[string]string{"notebookId": notebookId, "perm": perm})
			}
		}
	}
	this.UpdateGlobalArrMapConfig(this.adminUserId, "registerSharedNotebooks", notebooks)

	notes := []map[string]string{}
	// 共享笔记
	if len(registerSharedNoteIds) > 0 {
		for i := 0; i < len(registerSharedNoteIds); i++ {
			// 判断笔记本是否存在
			noteId := registerSharedNoteIds[i]
			if noteId == "" {
				continue
			}
			note := noteService.GetNote(noteId, registerSharedUserId)
			if note.NoteId == "" {
				ok = false
				msg = "The user has no such note: " + noteId
				return
			} else {
				perm := "0"
				if registerSharedNotePerms[i] == 1 {
					perm = "1"
				}
				notes = append(notes, map[string]string{"noteId": noteId, "perm": perm})
			}
		}
	}
	this.UpdateGlobalArrMapConfig(this.adminUserId, "registerSharedNotes", notes)

	// 复制
	noteIds := []string{}
	if len(registerCopyNoteIds) > 0 {
		for i := 0; i < len(registerCopyNoteIds); i++ {
			// 判断笔记本是否存在
			noteId := registerCopyNoteIds[i]
			if noteId == "" {
				continue
			}
			note := noteService.GetNote(noteId, registerSharedUserId)
			if note.NoteId == "" {
				ok = false
				msg = "The user has no such note: " + noteId
				return
			} else {
				noteIds = append(noteIds, noteId)
			}
		}
	}
	this.UpdateGlobalArrayConfig(this.adminUserId, "registerCopyNoteIds", noteIds)

	ok = true
	return
}

// 添加备份
func (this *ConfigService) AddBackup(path, remark string) bool {
	backups := this.GetGlobalArrMapConfig("backups") // [{}, {}]
	n := time.Now().Unix()
	nstr := fmt.Sprintf("%v", n)
	backups = append(backups, map[string]string{"createdTime": nstr, "path": path, "remark": remark})
	return this.UpdateGlobalArrMapConfig(this.adminUserId, "backups", backups)
}

func (this *ConfigService) getBackupDirname() string {
	n := time.Now()
	y, m, d := n.Date()
	return strconv.Itoa(y) + "_" + m.String() + "_" + strconv.Itoa(d) + "_" + fmt.Sprintf("%v", n.Unix())
}
func (this *ConfigService) Backup(remark string) (ok bool, msg string) {
	binPath := configService.GetGlobalStringConfig("mongodumpPath")
	config := revel.Config
	dbname, _ := config.String("db.dbname")
	host, _ := revel.Config.String("db.host")
	port, _ := revel.Config.String("db.port")
	username, _ := revel.Config.String("db.username")
	password, _ := revel.Config.String("db.password")
	// mongodump -h localhost -d leanote -o /root/mongodb_backup/leanote-9-22/ -u leanote -p nKFAkxKnWkEQy8Vv2LlM
	binPath = binPath + " -h " + host + " -d " + dbname + " --port " + port
	if username != "" {
		binPath += " -u " + username + " -p " + password
	}
	// 保存的路径
	dir := revel.BasePath + "/mongodb_backup/" + this.getBackupDirname()
	binPath += " -o " + dir
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		ok = false
		msg = fmt.Sprintf("%v", err)
		return
	}

	cmd := exec.Command("/bin/sh", "-c", binPath)
	Log(binPath)
	b, err := cmd.Output()
	if err != nil {
		msg = fmt.Sprintf("%v", err)
		ok = false
		Log("error:......")
		Log(string(b))
		return
	}
	ok = configService.AddBackup(dir, remark)
	return ok, msg
}

// 还原
func (this *ConfigService) Restore(createdTime string) (ok bool, msg string) {
	backups := this.GetGlobalArrMapConfig("backups") // [{}, {}]
	var i int
	var backup map[string]string
	for i, backup = range backups {
		if backup["createdTime"] == createdTime {
			break
		}
	}
	if i == len(backups) {
		return false, "Backup Not Found"
	}

	// 先备份当前
	ok, msg = this.Backup("Auto backup when restore from " + backup["createdTime"])
	if !ok {
		return
	}

	// mongorestore -h localhost -d leanote --directoryperdb /home/user1/gopackage/src/github.com/leanote/leanote/mongodb_backup/leanote_install_data/
	binPath := configService.GetGlobalStringConfig("mongorestorePath")
	config := revel.Config
	dbname, _ := config.String("db.dbname")
	host, _ := revel.Config.String("db.host")
	port, _ := revel.Config.String("db.port")
	username, _ := revel.Config.String("db.username")
	password, _ := revel.Config.String("db.password")
	// mongorestore -h localhost -d leanote -o /root/mongodb_backup/leanote-9-22/ -u leanote -p nKFAkxKnWkEQy8Vv2LlM
	binPath = binPath + " --drop -h " + host + " -d " + dbname + " --port " + port
	if username != "" {
		binPath += " -u " + username + " -p " + password
	}

	path := backup["path"] + "/" + dbname
	// 判断路径是否存在
	if !IsDirExists(path) {
		return false, path + " Is Not Exists"
	}

	binPath += " " + path

	cmd := exec.Command("/bin/sh", "-c", binPath)
	Log(binPath)
	b, err := cmd.Output()
	if err != nil {
		msg = fmt.Sprintf("%v", err)
		ok = false
		Log("error:......")
		Log(string(b))
		return
	}

	return true, ""
}
func (this *ConfigService) DeleteBackup(createdTime string) (bool, string) {
	backups := this.GetGlobalArrMapConfig("backups") // [{}, {}]
	var i int
	var backup map[string]string
	for i, backup = range backups {
		if backup["createdTime"] == createdTime {
			break
		}
	}
	if i == len(backups) {
		return false, "Backup Not Found"
	}

	// 删除文件夹之
	err := os.RemoveAll(backups[i]["path"])
	if err != nil {
		return false, fmt.Sprintf("%v", err)
	}

	// 删除之
	backups = append(backups[0:i], backups[i+1:]...)

	ok := this.UpdateGlobalArrMapConfig(this.adminUserId, "backups", backups)
	return ok, ""
}

func (this *ConfigService) UpdateBackupRemark(createdTime, remark string) (bool, string) {
	backups := this.GetGlobalArrMapConfig("backups") // [{}, {}]
	var i int
	var backup map[string]string
	for i, backup = range backups {
		if backup["createdTime"] == createdTime {
			break
		}
	}
	if i == len(backups) {
		return false, "Backup Not Found"
	}
	backup["remark"] = remark

	ok := this.UpdateGlobalArrMapConfig(this.adminUserId, "backups", backups)
	return ok, ""
}

// 得到备份
func (this *ConfigService) GetBackup(createdTime string) (map[string]string, bool) {
	backups := this.GetGlobalArrMapConfig("backups") // [{}, {}]
	var i int
	var backup map[string]string
	for i, backup = range backups {
		if backup["createdTime"] == createdTime {
			break
		}
	}
	if i == len(backups) {
		return map[string]string{}, false
	}
	return backup, true
}

//--------------
// sub domain
var defaultDomain string
var schema = "http://"
var port string

func init() {
	revel.OnAppStart(func() {
		/*
			不用配置的, 因为最终通过命令可以改, 而且有的使用nginx代理
			port  = strconv.Itoa(revel.HttpPort)
			if port != "80" {
				port = ":" + port
			} else {
				port = "";
			}
		*/

		siteUrl, _ := revel.Config.String("site.url") // 已包含:9000, http, 去掉成 leanote.com
		if strings.HasPrefix(siteUrl, "http://") {
			defaultDomain = siteUrl[len("http://"):]
		} else if strings.HasPrefix(siteUrl, "https://") {
			defaultDomain = siteUrl[len("https://"):]
			schema = "https://"
		}

		// port localhost:9000
		ports := strings.Split(defaultDomain, ":")
		if len(ports) == 2 {
			port = ports[1]
		}
		if port == "80" {
			port = ""
		} else {
			port = ":" + port
		}
	})
}

func (this *ConfigService) GetSchema() string {
	return schema
}

// 默认
func (this *ConfigService) GetDefaultDomain() string {
	return defaultDomain
}

// 包含http://
func (this *ConfigService) GetDefaultUrl() string {
	return schema + defaultDomain
}

// note
func (this *ConfigService) GetNoteDomain() string {
	subDomain := this.GetGlobalStringConfig("noteSubDomain")
	if subDomain != "" {
		return subDomain + port
	}
	return this.GetDefaultDomain() + "/note"
}
func (this *ConfigService) GetNoteUrl() string {
	return schema + this.GetNoteDomain()
}

// blog
func (this *ConfigService) GetBlogDomain() string {
	subDomain := this.GetGlobalStringConfig("blogSubDomain")
	if subDomain != "" {
		return subDomain + port
	}
	return this.GetDefaultDomain() + "/blog"
}
func (this *ConfigService) GetBlogUrl() string {
	return schema + this.GetBlogDomain()
}

// lea
func (this *ConfigService) GetLeaDomain() string {
	subDomain := this.GetGlobalStringConfig("leaSubDomain")
	if subDomain != "" {
		return subDomain + port
	}
	return this.GetDefaultDomain() + "/lea"
}
func (this *ConfigService) GetLeaUrl() string {
	return schema + this.GetLeaDomain()
}

func (this *ConfigService) GetUserUrl(domain string) string {
	return schema + domain + port
}
func (this *ConfigService) GetUserSubUrl(subDomain string) string {
	return schema + subDomain + "." + this.GetDefaultDomain()
}

// 是否允许自定义域名
func (this *ConfigService) AllowCustomDomain() bool {
	return configService.GetGlobalStringConfig("allowCustomDomain") != ""
}

// 是否是好的自定义域名
func (this *ConfigService) IsGoodCustomDomain(domain string) bool {
	blacks := this.GetGlobalArrayConfig("blackCustomDomains")
	for _, black := range blacks {
		if strings.Contains(domain, black) {
			return false
		}
	}
	return true
}
func (this *ConfigService) IsGoodSubDomain(domain string) bool {
	blacks := this.GetGlobalArrayConfig("blackSubDomains")
	LogJ(blacks)
	for _, black := range blacks {
		if domain == black {
			return false
		}
	}
	return true
}

// 上传大小
func (this *ConfigService) GetUploadSize(key string) float64 {
	f, _ := strconv.ParseFloat(this.GetGlobalStringConfig(key), 64)
	return f
}
func (this *ConfigService) GetInt64(key string) int64 {
	f, _ := strconv.ParseInt(this.GetGlobalStringConfig(key), 10, 64)
	return f
}
func (this *ConfigService) GetInt32(key string) int32 {
	f, _ := strconv.ParseInt(this.GetGlobalStringConfig(key), 10, 32)
	return int32(f)
}
func (this *ConfigService) GetUploadSizeLimit() map[string]float64 {
	return map[string]float64{
		"uploadImageSize":    this.GetUploadSize("uploadImageSize"),
		"uploadBlogLogoSize": this.GetUploadSize("uploadBlogLogoSize"),
		"uploadAttachSize":   this.GetUploadSize("uploadAttachSize"),
		"uploadAvatarSize":   this.GetUploadSize("uploadAvatarSize"),
	}
}

// 为用户得到全局的配置
// NoteController调用
func (this *ConfigService) GetGlobalConfigForUser() map[string]interface{} {
	uploadSizeConfigs := this.GetUploadSizeLimit()
	config := map[string]interface{}{}
	for k, v := range uploadSizeConfigs {
		config[k] = v
	}
	return config
}

// 主页是否是管理员的博客页
func (this *ConfigService) HomePageIsAdminsBlog() bool {
	return this.GetGlobalStringConfig("homePage") == ""
}

func (this *ConfigService) GetVersion() string {
	return "1.4.2"
}
