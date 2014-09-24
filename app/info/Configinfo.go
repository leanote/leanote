package info

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 配置
// 用户配置高于全局配置
type Config struct {
	UserId        bson.ObjectId       `bson:"_id"`
	StringConfigs map[string]string   `StringConfigs` // key => value
	ArrayConfigs  map[string][]string `ArrayConfigs`  // key => []value
	UpdatedTime   time.Time           `UpdatedTime`
}
