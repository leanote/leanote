package session

import (
	"github.com/robfig/revel"
	"github.com/leanote/leanote/app/lea/memcache"
//	. "leanote/app/lea"
)

// 使用filter
// 很巧妙就使用了memcache来处理session
// revel的session(cookie)只存sessionId, 其它信息存在memcache中

func SessionFilter(c *revel.Controller, fc []revel.Filter) {
	sessionId := c.Session.Id()
	
	// 从memcache中得到cache, 赋给session
	cache := revel.Session(memcache.Get(sessionId))
	if cache == nil {
		cache = revel.Session{}
		cache.Id()
	}
	c.Session = cache

	fc[0](c, fc[1:])
	
	// 再把session保存之
	memcache.Set(sessionId, c.Session, -1)	
	
	// 只留下sessionId
	c.Session = revel.Session{revel.SESSION_ID_KEY: sessionId}
}