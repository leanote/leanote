package session

import (
	"github.com/revel/revel"
//	. "github.com/leanote/leanote/app/lea"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 主要修改revel的cookie, 设置Domain
// 为了使sub domain共享cookie
// cookie.domain = leanote.com

// A signed cookie (and thus limited to 4kb in size).
// Restriction: Keys may not have a colon in them.
type Session map[string]string

const (
	SESSION_ID_KEY = "_ID"
	TIMESTAMP_KEY  = "_TS"
)

// expireAfterDuration is the time to live, in seconds, of a session cookie.
// It may be specified in config as "session.expires". Values greater than 0
// set a persistent cookie with a time to live as specified, and the value 0
// sets a session cookie.
var expireAfterDuration time.Duration
var cookieDomain = "" // life
func init() {
	// Set expireAfterDuration, default to 30 days if no value in config
	revel.OnAppStart(func() {
		var err error
		if expiresString, ok := revel.Config.String("session.expires"); !ok {
			expireAfterDuration = 30 * 24 * time.Hour
		} else if expiresString == "session" {
			expireAfterDuration = 0
		} else if expireAfterDuration, err = time.ParseDuration(expiresString); err != nil {
			panic(fmt.Errorf("session.expires invalid: %s", err))
		}
		
		cookieDomain, _ = revel.Config.String("cookie.domain")
	})
}

// Id retrieves from the cookie or creates a time-based UUID identifying this
// session.
func (s Session) Id() string {
	if sessionIdStr, ok := s[SESSION_ID_KEY]; ok {
		return sessionIdStr
	}

	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}

	s[SESSION_ID_KEY] = hex.EncodeToString(buffer)
	return s[SESSION_ID_KEY]
}

// getExpiration return a time.Time with the session's expiration date.
// If previous session has set to "session", remain it
func (s Session) getExpiration() time.Time {
	if expireAfterDuration == 0 || s[TIMESTAMP_KEY] == "session" {
		// Expire after closing browser
		return time.Time{}
	}
	return time.Now().Add(expireAfterDuration)
}

// cookie returns an http.Cookie containing the signed session.
func (s Session) cookie() *http.Cookie {
	var sessionValue string
	ts := s.getExpiration()
	s[TIMESTAMP_KEY] = getSessionExpirationCookie(ts)
	for key, value := range s {
		if strings.ContainsAny(key, ":\x00") {
			panic("Session keys may not have colons or null bytes")
		}
		if strings.Contains(value, "\x00") {
			panic("Session values may not have null bytes")
		}
		sessionValue += "\x00" + key + ":" + value + "\x00"
	}

	sessionData := url.QueryEscape(sessionValue)
	cookie := http.Cookie{
		Name:     revel.CookiePrefix + "_SESSION",
		Value:    revel.Sign(sessionData) + "-" + sessionData,
		Path:     "/",
		HttpOnly: revel.CookieHttpOnly,
		Secure:   revel.CookieSecure,
		Expires:  ts.UTC(),
	}
	
	if cookieDomain != "" {
		cookie.Domain = cookieDomain
	}
	
	return &cookie
}

// sessionTimeoutExpiredOrMissing returns a boolean of whether the session
// cookie is either not present or present but beyond its time to live; i.e.,
// whether there is not a valid session.
func sessionTimeoutExpiredOrMissing(session Session) bool {
	if exp, present := session[TIMESTAMP_KEY]; !present {
		return true
	} else if exp == "session" {
		return false
	} else if expInt, _ := strconv.Atoi(exp); int64(expInt) < time.Now().Unix() {
		return true
	}
	return false
}

// getSessionFromCookie returns a Session struct pulled from the signed
// session cookie.
func getSessionFromCookie(cookie *http.Cookie) Session {
	session := make(Session)

	// Separate the data from the signature.
	hyphen := strings.Index(cookie.Value, "-")
	if hyphen == -1 || hyphen >= len(cookie.Value)-1 {
		return session
	}
	sig, data := cookie.Value[:hyphen], cookie.Value[hyphen+1:]

	// Verify the signature.
	if !revel.Verify(data, sig) {
		revel.INFO.Println("Session cookie signature failed")
		return session
	}

	revel.ParseKeyValueCookie(data, func(key, val string) {
		session[key] = val
	})

	if sessionTimeoutExpiredOrMissing(session) {
		session = make(Session)
	}

	return session
}

// SessionFilter is a Revel Filter that retrieves and sets the session cookie.
// Within Revel, it is available as a Session attribute on Controller instances.
// The name of the Session cookie is set as CookiePrefix + "_SESSION".
func SessionFilter(c *revel.Controller, fc []revel.Filter) {
	session := restoreSession(c.Request.Request)
	// c.Session, 重新生成一个revel.Session给controller!!!
//	Log("sessoin--------")
//	LogJ(session)
	revelSession := revel.Session(session) // 强制转换 还是同一个对象, 但有个问题, 这样Session.Id()方法是用revel的了
	c.Session = revelSession
	// 生成sessionId
	c.Session.Id()
	sessionWasEmpty := len(c.Session) == 0

	// Make session vars available in templates as {{.session.xyz}}
	c.RenderArgs["session"] = c.Session

	fc[0](c, fc[1:])

	// Store the signed session if it could have changed.
	if len(c.Session) > 0 || !sessionWasEmpty {
		// 转换成lea.Session
		session = Session(c.Session)
		c.SetCookie(session.cookie())
	}
}

// restoreSession returns either the current session, retrieved from the
// session cookie, or a new session.
func restoreSession(req *http.Request) Session {
	cookie, err := req.Cookie(revel.CookiePrefix + "_SESSION")
	if err != nil {
		return make(Session)
	} else {
		return getSessionFromCookie(cookie)
	}
}

// getSessionExpirationCookie retrieves the cookie's time to live as a
// string of either the number of seconds, for a persistent cookie, or
// "session".
func getSessionExpirationCookie(t time.Time) string {
	if t.IsZero() {
		return "session"
	}
	return strconv.FormatInt(t.Unix(), 10)
}

// SetNoExpiration sets session to expire when browser session ends
func (s Session) SetNoExpiration() {
	s[TIMESTAMP_KEY] = "session"
}

// SetDefaultExpiration sets session to expire after default duration
func (s Session) SetDefaultExpiration() {
	delete(s, TIMESTAMP_KEY)
}