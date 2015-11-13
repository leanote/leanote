package info

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 配置, 每一个配置一行记录
type Config struct {
	ConfigId    bson.ObjectId       `bson:"_id"`
	UserId      bson.ObjectId       `UserId`
	Key         string              `Key`
	ValueStr    string              `ValueStr,omitempty`    // "1"
	ValueArr    []string            `ValueArr,omitempty`    // ["1","b","c"]
	ValueMap    map[string]string   `ValueMap,omitempty`    // {"a":"bb", "CC":"xx"}
	ValueArrMap []map[string]string `ValueArrMap,omitempty` // [{"a":"B"}, {}, {}]
	IsArr       bool                `IsArr`                 // 是否是数组
	IsMap       bool                `IsMap`                 // 是否是Map
	IsArrMap    bool                `IsArrMap`              // 是否是数组Map

	// StringConfigs map[string]string   `StringConfigs` // key => value
	// ArrayConfigs  map[string][]string `ArrayConfigs`  // key => []value

	UpdatedTime time.Time `UpdatedTime`
}
