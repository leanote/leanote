package lea

import (
	"encoding/json"
	"regexp"
	"strconv"
)

// 验证

var rulesStr = `{
	"username": [
		{"rule": "required", "msg": "inputUsername"}, 
		{"rule": "noSpecialChars", "msg": "noSpecialChars"},
		{"rule": "minLength", "data": "4", "msg": "minLength", "msgData": "4"}
	],
	"email": [
		{"rule": "required", "msg": "inputEmail"}, 
		{"rule": "email", "msg": "errorEmail"}
	],
	"password": [
		{"rule": "required", "msg": "inputPassword"}, 
		{"rule": "password", "msg": "errorPassword"}
	],
	"subDomain": [
		{"rule": "subDomain", "msg": "errorSubDomain"}
	],
	"domain": [
		{"rule": "domain", "msg": "errorDomain"}
	],
	"perPageSize": [
		{"rule": "min", "data": "1", "msg": "errorPerPageSize"}
	],
	"sortField": [
		{"rule": "sortField", "msg": "errorSortField"}
	]
}
`
var rulesMap map[string][]map[string]string

var rules = map[string]func(string, map[string]string) (bool, string){
	"required": func(value string, rule map[string]string) (ok bool, msg string) {
		if value == "" {
			return
		}
		ok = true
		return
	},
	"minLength": func(value string, rule map[string]string) (ok bool, msg string) {
		if value == "" {
			return
		}
		data := rule["data"]
		dataI, _ := strconv.Atoi(data)
		ok = len(value) >= dataI
		return
	},
	"min": func(value string, rule map[string]string) (ok bool, msg string) {
		if value == "" {
			return
		}
		data := rule["data"]
		dataI, _ := strconv.Atoi(data)
		vI, _ := strconv.Atoi(value)
		ok = vI >= dataI
		return
	},

	"sortField": func(value string, rule map[string]string) (ok bool, msg string) {
		if value == "" {
			return
		}
		sortFields := []string{"PublicTime", "CreatedTime", "UpdatedTime", "Title"}
		ok = InArray(sortFields, value)
		return
	},

	"password": func(value string, rule map[string]string) (ok bool, msg string) {
		if value == "" {
			return
		}
		ok = len(value) >= 6
		return
	},
	"email": func(value string, rule map[string]string) (ok bool, msg string) {
		if value == "" {
			return
		}
		ok = IsEmail(value)
		return
	},
	"noSpecialChars": func(value string, rule map[string]string) (ok bool, msg string) {
		if value == "" {
			return
		}
		ok = IsUsername(value)
		return
	},
	// www.baidu.com
	//
	"domain": func(value string, rule map[string]string) (ok bool, msg string) {
		if value == "" {
			ok = true
			return // 可为空
		}
		ok2, _ := regexp.MatchString(`[^0-9a-zA-Z_\.\-]`, value)
		ok = !ok2
		if !ok {
			return
		}
		ok = true
		return
	},
	// abcd
	"subDomain": func(value string, rule map[string]string) (ok bool, msg string) {
		if value == "" {
			ok = true
			return // 可为空
		}
		if len(value) < 4 {
			ok = false
			return
		}
		ok2, _ := regexp.MatchString(`[^0-9a-zA-Z_\-]`, value)
		ok = !ok2
		return
	},
}

func InitVd() {
	json.Unmarshal([]byte(rulesStr), &rulesMap)
	LogJ(rulesMap)
}

// 验证
// Vd("username", "life")

func Vd(name, value string) (ok bool, msg string) {
	rs, _ := rulesMap[name]

	for _, rule := range rs {
		ruleFunc, _ := rules[rule["rule"]]
		if ok2, msg2 := ruleFunc(value, rule); !ok2 {
			ok = false
			if msg2 != "" {
				msg = msg2
			} else {
				msg = rule["msg"]
			}
			msgData := rule["msgData"]
			if msgData != "" {
				msg += "-" + msgData
			}
			return
		}
	}
	ok = true
	return
}

func Vds(m map[string]string) (ok bool, msg string) {
	for name, value := range m {
		ok, msg = Vd(name, value)
		if !ok {
			return
		}
	}
	ok = true
	return
}
