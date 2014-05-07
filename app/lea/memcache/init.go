package memcache

import (
	"github.com/robfig/gomemcache/memcache"
)

var client *memcache.Client

func init() {
	// client = memcache.New("localhost:11211")	
}