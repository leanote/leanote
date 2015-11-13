package app

import (
	"encoding/json"
	"fmt"
	"github.com/leanote/leanote/app/controllers"
	"github.com/leanote/leanote/app/controllers/admin"
	"github.com/leanote/leanote/app/controllers/api"
	"github.com/leanote/leanote/app/controllers/member"
	"github.com/leanote/leanote/app/db"
	. "github.com/leanote/leanote/app/lea"
	_ "github.com/leanote/leanote/app/lea/binder"
	"github.com/leanote/leanote/app/lea/route"
	"github.com/leanote/leanote/app/service"
	"github.com/revel/revel"
	"html/template"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter, // Recover from panics and display an error page instead.
		route.RouterFilter,
		// revel.RouterFilter,            // Use the routing table to select the right Action
		// AuthFilter,						// Invoke the action.
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.

		// 使用SessionFilter标准版从cookie中得到sessionID, 然后通过MssessionFilter从Memcache中得到
		// session, 之后MSessionFilter将session只存sessionID然后返回给SessionFilter返回到web
		// session.SessionFilter,         // leanote session
		// session.MSessionFilter,         // leanote memcache session

		revel.FlashFilter,       // Restore and write the flash cookie.
		revel.ValidationFilter,  // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,        // Resolve the requested language
		revel.InterceptorFilter, // Run interceptors around the action.
		revel.CompressFilter,    // Compress the result.
		revel.ActionInvoker,     // Invoke the action.
	}

	revel.TemplateFuncs["raw"] = func(str string) template.HTML {
		return template.HTML(str)
	}
	revel.TemplateFuncs["trim"] = func(str string) string {
		str = strings.Trim(str, " ")
		str = strings.Trim(str, " ")

		str = strings.Trim(str, "\n")
		str = strings.Trim(str, "&nbsp;")

		// 以下两个空格不一样
		str = strings.Trim(str, " ")
		str = strings.Trim(str, " ")
		return str
	}
	revel.TemplateFuncs["add"] = func(i int) string {
		i = i + 1
		return fmt.Sprintf("%v", i)
	}
	revel.TemplateFuncs["sub"] = func(i int) int {
		i = i - 1
		return i
	}
	// 增加或减少
	revel.TemplateFuncs["incr"] = func(n, i int) int {
		n = n + i
		return n
	}
	revel.TemplateFuncs["join"] = func(arr []string) template.HTML {
		if arr == nil {
			return template.HTML("")
		}
		return template.HTML(strings.Join(arr, ","))
	}
	revel.TemplateFuncs["concat"] = func(s1, s2 string) template.HTML {
		return template.HTML(s1 + s2)
	}
	revel.TemplateFuncs["concatStr"] = func(strs ...string) string {
		str := ""
		for _, s := range strs {
			str += s
		}
		return str
	}
	revel.TemplateFuncs["decodeUrlValue"] = func(i string) string {
		v, _ := url.ParseQuery("a=" + i)
		return v.Get("a")
	}
	revel.TemplateFuncs["json"] = func(i interface{}) string {
		b, _ := json.Marshal(i)
		return string(b)
	}
	revel.TemplateFuncs["jsonJs"] = func(i interface{}) template.JS {
		b, _ := json.Marshal(i)
		return template.JS(string(b))
	}
	revel.TemplateFuncs["datetime"] = func(t time.Time) template.HTML {
		return template.HTML(t.Format("2006-01-02 15:04:05"))
	}
	revel.TemplateFuncs["dateFormat"] = func(t time.Time, format string) template.HTML {
		return template.HTML(t.Format(format))
	}
	revel.TemplateFuncs["unixDatetime"] = func(unixSec string) template.HTML {
		sec, _ := strconv.Atoi(unixSec)
		t := time.Unix(int64(sec), 0)
		return template.HTML(t.Format("2006-01-02 15:04:05"))
	}

	// interface是否有该字段
	revel.TemplateFuncs["has"] = func(i interface{}, key string) bool {
		t := reflect.TypeOf(i)
		_, ok := t.FieldByName(key)
		return ok
	}

	// tags
	// 2014/12/30 标签添加链接
	revel.TemplateFuncs["blogTags"] = func(renderArgs map[string]interface{}, tags []string) template.HTML {
		if tags == nil || len(tags) == 0 {
			return ""
		}
		locale, _ := renderArgs[revel.CurrentLocaleRenderArg].(string)
		tagStr := ""
		lenTags := len(tags)

		tagPostUrl, _ := renderArgs["tagPostsUrl"].(string)

		for i, tag := range tags {
			str := revel.Message(locale, tag)
			var classes = "label"
			if strings.HasPrefix(str, "???") {
				str = tag
			}
			if InArray([]string{"red", "blue", "yellow", "green"}, tag) {
				classes += " label-" + tag
			} else {
				classes += " label-default"
			}

			classes += " label-post"
			var url = tagPostUrl + "/" + url.QueryEscape(tag)
			tagStr += "<a class=\"" + classes + "\" href=\"" + url + "\">" + str + "</a>"
			if i != lenTags-1 {
				tagStr += " "
			}
		}
		return template.HTML(tagStr)
	}

	revel.TemplateFuncs["blogTagsForExport"] = func(renderArgs map[string]interface{}, tags []string) template.HTML {
		if tags == nil || len(tags) == 0 {
			return ""
		}
		tagStr := ""
		lenTags := len(tags)

		for i, tag := range tags {
			str := tag
			var classes = "label"
			if InArray([]string{"red", "blue", "yellow", "green"}, tag) {
				classes += " label-" + tag
			} else {
				classes += " label-default"
			}

			classes += " label-post"
			tagStr += "<span class=\"" + classes + "\" >" + str + "</span>"
			if i != lenTags-1 {
				tagStr += " "
			}
		}
		return template.HTML(tagStr)
	}

	// 不用revel的msg
	revel.TemplateFuncs["leaMsg"] = func(renderArgs map[string]interface{}, key string) template.HTML {
		locale, _ := renderArgs[revel.CurrentLocaleRenderArg].(string)
		str := revel.Message(locale, key)
		if strings.HasPrefix(str, "???") {
			str = key
		}
		return template.HTML(str)
	}

	// lea++
	revel.TemplateFuncs["blogTagsLea"] = func(renderArgs map[string]interface{}, tags []string, typeStr string) template.HTML {
		if tags == nil || len(tags) == 0 {
			return ""
		}
		locale, _ := renderArgs[revel.CurrentLocaleRenderArg].(string)
		tagStr := ""
		lenTags := len(tags)

		tagPostUrl := "http://lea.leanote.com/"
		if typeStr == "recommend" {
			tagPostUrl += "?tag="
		} else if typeStr == "latest" {
			tagPostUrl += "latest?tag="
		} else {
			tagPostUrl += "subscription?tag="
		}

		for i, tag := range tags {
			str := revel.Message(locale, tag)
			var classes = "label"
			if strings.HasPrefix(str, "???") {
				str = tag
			}
			if InArray([]string{"red", "blue", "yellow", "green"}, tag) {
				classes += " label-" + tag
			} else {
				classes += " label-default"
			}
			classes += " label-post"
			var url = tagPostUrl + url.QueryEscape(tag)
			tagStr += "<a class=\"" + classes + "\" href=\"" + url + "\">" + str + "</a>"
			if i != lenTags-1 {
				tagStr += " "
			}
		}
		return template.HTML(tagStr)
	}

	/*
		revel.TemplateFuncs["blogTags"] = func(tags []string) template.HTML {
			if tags == nil || len(tags) == 0 {
				return ""
			}
			// TODO 这里判断语言, 从语言包中拿
			tagMap := map[string]string{"red": "红色", "yellow": "黄色", "blue": "蓝色", "green": "绿色"}
			tagStr := ""
			lenTags := len(tags)
			for i, tag := range tags {
				if text, ok := tagMap[tag]; ok {
					tagStr += text
				} else {
					tagStr += tag
				}
				if i != lenTags - 1 {
					tagStr += ","
				}
			}
			return template.HTML(tagStr)
		}
	*/
	revel.TemplateFuncs["li"] = func(a string) string {
		return ""
	}
	// str连接
	revel.TemplateFuncs["urlConcat"] = func(url string, v ...interface{}) string {
		html := ""
		for i := 0; i < len(v); i = i + 2 {
			item := v[i]
			if i+1 == len(v) {
				break
			}
			value := v[i+1]
			if item != nil && value != nil {
				keyStr, _ := item.(string)
				valueStr, err := value.(string)
				if !err {
					valueInt, _ := value.(int)
					valueStr = strconv.Itoa(valueInt)
				}
				if keyStr != "" && valueStr != "" {
					s := keyStr + "=" + valueStr
					if html != "" {
						html += "&" + s
					} else {
						html += s
					}
				}
			}
		}

		if html != "" {
			if strings.Index(url, "?") >= 0 {
				return url + "&" + html
			} else {
				return url + "?" + html
			}
		}
		return url
	}

	revel.TemplateFuncs["urlCond"] = func(url string, sorterI, keyords interface{}) template.HTML {
		return ""
	}

	// http://stackoverflow.com/questions/14226416/go-lang-templates-always-quotes-a-string-and-removes-comments
	revel.TemplateFuncs["rawMsg"] = func(renderArgs map[string]interface{}, message string, args ...interface{}) template.JS {
		str, ok := renderArgs[revel.CurrentLocaleRenderArg].(string)
		if !ok {
			return ""
		}
		return template.JS(revel.Message(str, message, args...))
	}

	// 为后台管理sorter th使用
	// 必须要返回HTMLAttr, 返回html, golang 会执行安全检查返回ZgotmplZ
	// sorterI 可能是nil, 所以用interfalce{}来接收
	/*
		data-url="/adminUser/index"
		data-sorter="email"
		class="th-sortable {{if eq .sorter "email-up"}}th-sort-up{{else}}{{if eq .sorter "email-down"}}th-sort-down{{end}}{{end}}"
	*/
	revel.TemplateFuncs["sorterTh"] = func(url, sorterField string, sorterI interface{}) template.HTMLAttr {
		sorter := ""
		if sorterI != nil {
			sorter, _ = sorterI.(string)
		}
		html := "data-url=\"" + url + "\" data-sorter=\"" + sorterField + "\""
		html += " class=\"th-sortable "
		if sorter == sorterField+"-up" {
			html += "th-sort-up\""
		} else if sorter == sorterField+"-down" {
			html += "th-sort-down"
		}
		html += "\""
		return template.HTMLAttr(html)
	}

	// pagination
	revel.TemplateFuncs["page"] = func(urlBase string, page, pageSize, count int) template.HTML {
		if count == 0 {
			return ""
		}
		totalPage := int(math.Ceil(float64(count) / float64(pageSize)))

		preClass := ""
		prePage := page - 1
		if prePage == 0 {
			prePage = 1
		}
		nextClass := ""
		nextPage := page + 1
		var preUrl, nextUrl string

		preUrl = urlBase + "?page=" + strconv.Itoa(prePage)
		nextUrl = urlBase + "?page=" + strconv.Itoa(nextPage)

		// 没有上一页了
		if page == 1 {
			preClass = "disabled"
			preUrl = "#"
		}
		// 没有下一页了
		if totalPage <= page {
			nextClass = "disabled"
			nextUrl = "#"
		}
		return template.HTML("<li class='" + preClass + "'><a href='" + preUrl + "'>Previous</a></li> <li  class='" + nextClass + "'><a href='" + nextUrl + "'>Next</a></li>")
	}
	// life
	// https://groups.google.com/forum/#!topic/golang-nuts/OEdSDgEC7js
	// http://play.golang.org/p/snygrVpQva
	// http://grokbase.com/t/gg/golang-nuts/142a6dhfh3/go-nuts-text-template-using-comparison-operators-eq-gt-etc-on-non-existent-variable-causes-the-template-to-stop-outputting-but-with-no-error-correct-behaviour
	/*
		revel.TemplateFuncs["gt"] = func(a1, a2 interface{}) bool {
			switch a1.(type) {
			case string:
				switch a2.(type) {
				case string:
					return reflect.ValueOf(a1).String() > reflect.ValueOf(a2).String()
				}
			case int, int8, int16, int32, int64:
				switch a2.(type) {
				case int, int8, int16, int32, int64:
					return reflect.ValueOf(a1).Int() > reflect.ValueOf(a2).Int()
				}
			case uint, uint8, uint16, uint32, uint64:
				switch a2.(type) {
				case uint, uint8, uint16, uint32, uint64:
					return reflect.ValueOf(a1).Uint() > reflect.ValueOf(a2).Uint()
				}
			case float32, float64:
				switch a2.(type) {
				case float32, float64:
					return reflect.ValueOf(a1).Float() > reflect.ValueOf(a2).Float()
				}
			}
			return false
		}
	*/

	/*
		{{range $i := N 1 10}}
	        <div>{{$i}}</div>
	    {{end}}
	*/
	revel.TemplateFuncs["N"] = func(start, end int) (stream chan int) {
		stream = make(chan int)
		go func() {
			for i := start; i <= end; i++ {
				stream <- i
			}
			close(stream)
		}()
		return
	}

	// init Email
	revel.OnAppStart(func() {
		// 数据库
		db.Init("", "")
		// email配置
		InitEmail()
		InitVd()
		// memcache.InitMemcache() // session服务
		// 其它service
		service.InitService()
		controllers.InitService()
		admin.InitService()
		member.InitService()
		service.ConfigS.InitGlobalConfigs()
		api.InitService()
	})
}
