package lea

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2/bson"
	"io"
	math_rand "math/rand"
	"regexp"
	"strings"
	"time"
)

// 字符串

// md5
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 3位数的转换, 为了用bson.id -> 3位数
func Digest3(str string) string {
	var b rune = 0
	for _, k := range str {
		b += k
	}
	return fmt.Sprintf("%d", b%1000)
}
func Digest2(str string) string {
	var b rune = 0
	for _, k := range str {
		b += k
	}
	return fmt.Sprintf("%d", b%100)
}

// Guid
func NewGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

// 后面加个str生成之, 更有保障, 确保唯一
func NewGuidWith(str string) string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString([]byte(string(b) + str)))
}

// 以byte来截取
func SubstringByte(str string, start int) string {
	return substr(str, start, len(str)-start, false)
}

func Substring(str string, start int) string {
	return substr(str, start, len(str)-start, true)
}

func Substr(str string, start, length int) string {
	return substr(str, start, length, true)
}

func substr(str string, start, length int, isRune bool) string {
	rs := []rune(str)
	rs2 := []byte(str)

	rl := len(rs)
	if !isRune {
		rl = len(rs2)
	}
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	if isRune {
		return string(rs[start:end])
	}
	return string(rs2[start:end])
}

// 截取含html的字符串
/*
<div>
	abcdef
	<p>xxxxxxxxxxx</p>
	<img src="" />
	<div>
		abcdfef
		<p>aksdfjsld
		dfsadf
		--------------- 这里截断
		</p>
		<p>xxxxxxxxxxxxxxx</p>
		abcdef
	</div>
</div>
那么只要在最后的html中添加</p></div></div>补全结束即可
思路:
1. 找到截断的位置
2. 去掉成对的标签
3. 去掉不要结束的标签
4. 补全标签
*/
func ReplaceAll(oldStr, pattern, newStr string) string {
	p, _ := regexp.Compile(pattern)
	s := p.ReplaceAll([]byte(oldStr), []byte(newStr))
	return string(s)
}

// 获取纯文本
func SubStringHTMLToRaw(param string, length int) string {
	if param == "" {
		return param
	}
	n := 0
	var temp rune // 中文问题, 用rune来解决
	rStr := []rune(param)
	lenStr := len(rStr)
	isCode := false

	resultRune := make([]rune, length)
	// s := ""
	for i := 0; i < lenStr; i++ {
		temp = rStr[i]
		if temp == '<' {
			isCode = true
			continue
		} else if temp == '>' {
			isCode = false
			resultRune[n] = ' '

			n++
			if n >= length {
				break
			}
			continue
		}
		if !isCode {
			resultRune[n] = temp
			// s += string(temp)
			n++
			if n >= length {
				break
			}
		}
	}
	result := string(resultRune[0:n])
	return strings.Trim(result, " ")
}

// 自带方法补全html
func fixHtml(result string) string {
	// 取出所有标签
	tempResult := ReplaceAll(result, "(>)[^<>]*(<?)", "$1$2") // 把标签中间的所有内容都去掉了

	// 2. 去掉不需要结素标记的HTML标记
	tempResult = ReplaceAll(tempResult, "</?(embed|AREA|BASE|BASEFONT|BR|COL|COLGROUP|DD|DT|FRAME|HEAD|HR|IMG|INPUT|ISINDEX|LI|LINK|META|OPTION|PARAM|area|wbr|br|col|colgroup|dd|dt|frame|hr|img|input|isindex|link|meta|param)[^<>]*/?>", "")

	// 把<div class=xxx的class=xxx去掉
	tempResult = ReplaceAll(tempResult, "<(/?[a-zA-Z]+)[^<>]*>", "<$1>")

	// 3 只能用正则,+stack来去有结束的
	// golang的正则暂不支持back reference, 以后可以用它来去掉重复的标签
	p, _ := regexp.Compile("<(/?[a-zA-Z]+)[^<>]*>") // 得到所有的<div>, </div>...
	strs := p.FindAllString(tempResult, -1)

	//	fmt.Println(strs)
	stack := make([]string, len(strs))
	stackP := -1
	for _, each := range strs {
		if stackP >= 0 {
			// 匹配
			if stack[stackP][1:] == each[2:] {
				stackP--
				continue
			}
		}
		stackP++
		stack[stackP] = each
	}
	// 补全tag
	if stackP != -1 {
		fmt.Println(stack[0 : stackP+1])

		for _, each := range stack[0 : stackP+1] {
			if each[1] != '/' {
				result += "</" + each[1:]
			}
		}
	}

	return result
}

// 获取摘要, HTML
func SubStringHTML(param string, length int, end string) string {
	if param == "" {
		return param
	}
	result := ""

	rStr := []rune(param)
	lenStr := len(rStr)

	if lenStr <= length {
		result = param
	} else {
		// 1
		n := 0
		var temp rune   // 中文问题, 用rune来解决
		isCode := false //是不是HTML代码
		isHTML := false //是不是HTML特殊字符,如&nbsp;
		var i = 0
		for ; i < lenStr; i++ {
			temp = rStr[i]
			if temp == '<' {
				isCode = true
			} else if temp == '&' {
				isHTML = true
			} else if temp == '>' && isCode {
				// n = n - 1
				isCode = false
			} else if temp == ';' && isHTML {
				isHTML = false
			}
			if !isCode && !isHTML {
				n = n + 1
			}
			// 每一次都相加, 速度非常慢!, 重新分配内存, 7倍的差距
			// result += string(temp)
			if n >= length {
				break
			}
		}

		result = string(rStr[0:i])

		if end != "" {
			result += end
		}
	}

	// 使用goquery来取出html, 为了补全html
	htmlReader := bytes.NewBufferString(result)
	dom, err1 := goquery.NewDocumentFromReader(htmlReader)
	if err1 == nil {
		html, _ := dom.Html()
		html = strings.Replace(html, "<html><head></head><body>", "", 1)
		html = strings.Replace(html, "</body></html>", "", 1)

		// TODO 把style="float: left"去掉
		return html

		// 如果有错误, 则使用自己的方法补全, 有风险
	} else {
		return fixHtml(result)
	}
}

// 是否是合格的密码
func IsGoodPwd(pwd string) (bool, string) {
	if pwd == "" {
		return false, "密码不能为空"
	}
	if len(pwd) < 6 {
		return false, "密码至少6位"
	}
	return true, ""
}

// 是否是email
func IsEmail(email string) bool {
	if email == "" {
		return false
	}
	ok, _ := regexp.MatchString(`^([a-zA-Z0-9]+[_|\_|\.|\-]?)*[_a-z\-A-Z0-9]+@([a-zA-Z0-9]+[_|\_|\.|\-]?)*[a-zA-Z0-9\-]+\.[0-9a-zA-Z]{2,6}$`, email)
	return ok
}

// 是否只包含数字, 字母 -, _
func IsUsername(username string) bool {
	if username == "" {
		return false
	}
	ok, _ := regexp.MatchString(`[^0-9a-zA-Z_\-]`, username)
	return !ok
}

// 是否是ObjectId
func IsObjectId(id string) (ok bool) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		// 证明有错误发生
		if err := recover(); err != nil {
			ok = false
		} else {
			ok = true
		}
	}()
	bson.ObjectIdHex(id)
	return
}

// 随机密码
// num 几位
func RandomPwd(num int) string {
	chars := make([]byte, 62)
	j := 0
	for i := 48; i <= 57; i++ {
		chars[j] = byte(i)
		j++
	}
	for i := 65; i <= 90; i++ {
		chars[j] = byte(i)
		j++
	}
	for i := 97; i <= 122; i++ {
		chars[j] = byte(i)
		j++
	}
	j--

	str := ""
	math_rand.Seed(time.Now().UnixNano())
	for i := 0; i < num; i++ {
		x := math_rand.Intn(j)
		str += string(chars[x])
	}

	return str
}

func InArray(arr []string, str string) bool {
	if arr == nil {
		return false
	}
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// 将名称的特殊字符去掉
func FixFilename(filename string) string {
	if filename != "" {
		// 把特殊字段给替换掉
		//		str := `life "%&()+,/:;<>=?@\|`
		// . == \\.
		// $ === \\$
		reg, _ := regexp.Compile("\\.|/|#|\\$|!|\\^|\\*|'| |\"|%|&|\\(|\\)|\\+|\\,|/|:|;|<|>|=|\\?|@|\\||\\\\")
		filename = reg.ReplaceAllString(filename, "-")
		filename = strings.Trim(filename, "-") // 左右单独的-去掉
		// 把空格替换成-
		//		filename = strings.Replace(filename, " ", "-", -1)
		for strings.Index(filename, "--") >= 0 { // 防止出现连续的--
			filename = strings.Replace(filename, "--", "-", -1)
		}
		return filename
	}
	return filename
}

// 是否是合法的时间
// 不是, 0001-01-01T00:00:00Z, 且比今天小
func IsValidTime(t time.Time) bool {
	if t.Year() > 20 {
		now := time.Now()
		Log("------")
		Log(t)
		Log(now)
		if t.Before(now) {
			return true
		}
	}
	return false
}

// url传过来的时间没有时区信息, 转到本地时间
func ToLocalTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
}

// 修复传过来的时间, 如果比今天大, 则设为现在
func FixUrlTime(t time.Time) time.Time {
	localTime := ToLocalTime(t)
	if IsValidTime(localTime) {
		return localTime
	}
	return time.Now()
}

// 得到用户的随机文件路径 3位/userId/2位
func GetRandomFilePath(userId, uuid string) string {
	if uuid == "" {
		uuid = NewGuid()
	}
	return Digest3(userId) + "/" + userId + "/" + Digest2(uuid)
}
