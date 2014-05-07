package controllers

import (
	"code.google.com/p/goauth2/oauth"
	"github.com/revel/revel"
	"github.com/leanote/leanote/app/lea/netutil"
	. "github.com/leanote/leanote/app/lea"
	"encoding/json"
	"fmt"
)

type Oauth struct {
	BaseController
}

var oauthCfg = &oauth.Config{
	ClientId:     "3790fbf1fc14bc6c5d85",
	ClientSecret: "e9dadfe601c7caa6df9b33db3e7539945c60dfa2",
	AuthURL:      "https://github.com/login/oauth/authorize",
	TokenURL:     "https://github.com/login/oauth/access_token",
	RedirectURL:  "http://leanote.com/oauth/githubCallback",
	Scope: "user",
}

// 用户允许后, github返回到leanote
// 通过code得到token
// https://github.com/login/oauth/authorize?access_type=&approval_prompt=&client_id=3790fbf1fc14bc6c5d85&redirect_uri=http%3A%2F%2F127.0.0.1%3A8080%2Foauth2callback&response_type=code&scope=user&state=
func (c Oauth) GithubCallback(code string) revel.Result {
	t := &oauth.Transport{Config: oauthCfg}

	// Exchange the received code for a token
	tok, err := t.Exchange(code)
	token := tok.AccessToken
	if err != nil || token == "" {
		c.RenderArgs["title"] = "error"
		return c.RenderTemplate("oauth/oauth_callback_error.html")
	}

	// 得到用户信息
	profileInfoURL := "https://api.github.com/user"
	url := fmt.Sprintf("%s?access_token=%s", profileInfoURL, token)
	content, err2 := netutil.GetContent(url)
	if err2 != nil {
		c.RenderArgs["title"] = "error"
		return c.RenderTemplate("oauth/oauth_callback_error.html")
	}
	// 转成map
	profileInfo := map[string]interface{}{}
	Log(string(content))
	err2 = json.Unmarshal(content, &profileInfo)
	if err2 != nil {
		c.RenderArgs["title"] = "error"
		return c.RenderTemplate("oauth/oauth_callback_error.html")
	}

	usernameI := profileInfo["login"]

	username, _ := usernameI.(string)
	userId := username

	// 注册
	isExists, userInfo := authService.ThirdRegister("github", userId, username)
	c.RenderArgs["isExists"] = isExists
	c.RenderArgs["userInfo"] = userInfo

	// 登录之
	c.SetSession(userInfo)

	return c.Redirect("/note")
}
