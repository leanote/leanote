package session

import (
	"github.com/revel/revel"
	"github.com/leanote/leanote/app/lea/memcache"
	. "github.com/leanote/leanote/app/lea"
)

// 使用filter
// 很巧妙就使用了memcache来处理session
// revel的session(cookie)只存sessionId, 其它信息存在memcache中

func MSessionFilter(c *revel.Controller, fc []revel.Filter) {
	sessionId := c.Session.Id()
	
	// 从memcache中得到cache, 赋给session
	cache := revel.Session(memcache.GetMap(sessionId))
	
	Log("memcache")
	LogJ(cache)
	if cache == nil {
		cache = revel.Session{}
		cache.Id()
	}
	c.Session = cache
	
	// Make session vars available in templates as {{.session.xyz}}
	c.RenderArgs["session"] = c.Session

	fc[0](c, fc[1:])
	
	// 再把session保存之
	LogJ(c.Session)
	memcache.SetMap(sessionId, c.Session, -1)	

	// 只留下sessionId
	c.Session = revel.Session{revel.SESSION_ID_KEY: sessionId}
}