package app

import (
	"github.com/revel/revel"
	. "github.com/leanote/leanote/app/lea"
	_ "github.com/leanote/leanote/app/lea/binder"
	"reflect"
	"fmt"
	"html/template"
	"math"
	"strconv"
	"time"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		// AuthFilter,						// Invoke the action.
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		
//		session.SessionFilter,         // leanote memcache session life
		
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}
	
	revel.TemplateFuncs["raw"] = func(str string) template.HTML {
		return template.HTML(str)
	}
	revel.TemplateFuncs["add"] = func(i int) template.HTML {
		i = i + 1;
		return template.HTML(fmt.Sprintf("%v", i))
	}
	revel.TemplateFuncs["concat"] = func(s1, s2 string) template.HTML {
		return template.HTML(s1 + s2)
	}
	revel.TemplateFuncs["datetime"] = func(t time.Time) template.HTML {
		return template.HTML(t.Format("2006-01-02 15:04:05"))
	}
	
	// interface是否有该字段
	revel.TemplateFuncs["has"] = func(i interface{}, key string) bool {
		t := reflect.TypeOf(i) 
		_, ok := t.FieldByName(key)
		return ok
	}

	// tags
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
	
	// pagination
	revel.TemplateFuncs["page"] = func(userId, notebookId string, page, pageSize, count int) template.HTML {
		if count == 0 {
			return "";
		}
		totalPage := int(math.Ceil(float64(count)/float64(pageSize)))
		
		preClass := ""
		prePage := page - 1
		if prePage == 0 {
			prePage = 1
		}
		nextClass := ""
		nextPage := page + 1
		var preUrl, nextUrl string
		
		urlBase := "/blog/" + userId
		if notebookId != "" {
			urlBase += "/" + notebookId
		}
		
		preUrl = urlBase + "?page="  + strconv.Itoa(prePage)
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
	
	// init Email
	revel.OnAppStart(func() {
		InitEmail()
	})
}