package memcache

import (
	"github.com/robfig/gomemcache/memcache"
	"encoding/json"
	"strconv"
)

var client *memcache.Client

// onAppStart后调用
func InitMemcache() {
	client = memcache.New("localhost:11211")	
}

//------------
// map

func SetMap(key string, value map[string]string, expiration int32) {
	// 把value转成byte
	bytes, _ := json.Marshal(value)
	if expiration == -1 {
		expiration = 30 * 24 * 60 * 60 // 30天
	}
	client.Set(&memcache.Item{Key: key, Value: bytes, Expiration: expiration})
}

func GetMap(key string) map[string]string {
	item, err := client.Get(key)
	if err != nil {
		return nil
	}
	
	m := map[string]string{}
	json.Unmarshal(item.Value, &m)	
	return m
}

//------------
// string
func GetString(key string) string {
	item, err := client.Get(key)
	if err != nil {
		return ""
	}
	return string(item.Value)
}
func SetString(key string, value string, expiration int32) {
	if expiration == -1 {
		expiration = 30 * 24 * 60 * 60 // 30天
	}
	client.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: expiration})
}

//-------------------------
// int, 是通过转成string来存的

func GetInt(key string) int {
	str := GetString(key)
	i, _ := strconv.Atoi(str)
	return i
}
func SetInt(key string, value int, expiration int32) {
	str := strconv.Itoa(value)
	SetString(key, str, expiration)
}
