package memcache

import (
	"github.com/robfig/gomemcache/memcache"
	"encoding/json"
)

func Set(key string, value map[string]string, expiration int32) {
	// 把value转成byte
	bytes, _ := json.Marshal(value)
	if expiration == -1 {
		expiration = 30 * 24 * 60 * 60 // 30天
	}
	client.Set(&memcache.Item{Key: key, Value: bytes, Expiration: expiration})
}

func Get(key string) map[string]string {
	item, err := client.Get(key)
	if err != nil {
		return nil
	}
	
	m := map[string]string{}
	json.Unmarshal(item.Value, &m)	
	return m
}