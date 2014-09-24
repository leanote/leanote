package service

import (
	"github.com/leanote/leanote/app/info"
//	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/db"
	"gopkg.in/mgo.v2/bson"
	"github.com/revel/revel"
	"time"
)

// 配置服务
type ConfigService struct {
	// 全局的
	GlobalStringConfigs map[string]string
	GlobalArrayConfigs map[string][]string
	
	// 两种配置, 用户自己的
	UserStringConfigs map[string]string
	UserArrayConfigs map[string][]string
	
	// 合并之后的
	AllStringConfigs map[string]string
	AllArrayConfigs map[string][]string
}

var adminUserId = ""

// appStart时 将全局的配置从数据库中得到作为全局
func (this *ConfigService) InitGlobalConfigs() bool {
	this.GlobalStringConfigs = map[string]string{}
	this.GlobalArrayConfigs = map[string][]string{}
	
	this.UserStringConfigs = map[string]string{}
	this.UserArrayConfigs = map[string][]string{}
	
	this.AllStringConfigs = map[string]string{}
	this.AllArrayConfigs = map[string][]string{}
	
	adminUsername, _ := revel.Config.String("adminUsername")
	if adminUsername == "" {
		adminUsername = "admin"
	}
	
	userInfo := userService.GetUserInfoByAny(adminUsername)
	if userInfo.UserId == "" {
		return false
	}
	adminUserId = userInfo.UserId.Hex()
	
	configs := info.Config{}
	db.Get2(db.Configs, userInfo.UserId, &configs)
	
	if configs.UserId == "" {
		db.Insert(db.Configs, info.Config{UserId: userInfo.UserId, StringConfigs: map[string]string{}, ArrayConfigs: map[string][]string{}})
	}
	
	this.GlobalStringConfigs = configs.StringConfigs;
	this.GlobalArrayConfigs = configs.ArrayConfigs;
	
	// 复制到所有配置上
	for key, value := range this.GlobalStringConfigs {
		this.AllStringConfigs[key] = value
	}
	for key, value := range this.GlobalArrayConfigs {
		this.AllArrayConfigs[key] = value
	}
	
	return true
}

// 用户登录后获取用户自定义的配置, 并将所有的配置都用上
func (this *ConfigService) InitUserConfigs(userId string) bool {
	configs := info.Config{}
	db.Get(db.Configs, userId, &configs)
	
	if configs.UserId == "" {
		db.Insert(db.Configs, info.Config{UserId: bson.ObjectIdHex(userId), StringConfigs: map[string]string{}, ArrayConfigs: map[string][]string{}})
	}
	
	this.UserStringConfigs = configs.StringConfigs;
	this.UserArrayConfigs = configs.ArrayConfigs;
	
	// 合并配置
	for key, value := range this.UserStringConfigs {
		this.AllStringConfigs[key] = value
	}
	for key, value := range this.UserArrayConfigs {
		this.AllArrayConfigs[key] = value
	}
	
	return true
}

// 获取配置
func (this *ConfigService) GetStringConfig(key string) string {
	return this.AllStringConfigs[key]
}
func (this *ConfigService) GetArrayConfig(key string) []string {
	arr := this.AllArrayConfigs[key]
	if arr == nil {
		return []string{}
	}
	return arr
}

// 更新用户配置
func (this *ConfigService) UpdateUserStringConfig(userId, key string, value string) bool {
	this.UserStringConfigs[key] = value
	this.AllStringConfigs[key] = value
	if userId == adminUserId {
		this.GlobalStringConfigs[key] = value
	}
	
	// 保存到数据库中
	return db.UpdateByQMap(db.Configs, bson.M{"_id": bson.ObjectIdHex(userId)}, 
	bson.M{"StringConfigs": this.UserStringConfigs, "UpdatedTime": time.Now()})
}
func (this *ConfigService) UpdateUserArrayConfig(userId, key string, value []string) bool {
	this.UserArrayConfigs[key] = value
	this.AllArrayConfigs[key] = value
	if userId == adminUserId {
		this.GlobalArrayConfigs[key] = value
	}
	
	// 保存到数据库中
	return db.UpdateByQMap(db.Configs, bson.M{"_id": bson.ObjectIdHex(userId)}, 
	bson.M{"ArrayConfigs": this.UserArrayConfigs, "UpdatedTime": time.Now()})
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