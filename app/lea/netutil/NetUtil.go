package netutil

import (
	"os"
	"strings"
	//	"path/filepath"
	. "github.com/leanote/leanote/app/lea"
	"io/ioutil"
	"net"
	"net/http"
)

// net的util

// toPath 文件保存的目录
// 默认是/tmp
// 返回文件的完整目录
func WriteUrl(url string, toPath string) (length int64, newFilename, path string, ok bool) {
	if url == "" {
		return
	}
	content, err := GetContent(url)
	if err != nil {
		return
	}

	length = int64(len(content))

	// a.html?a=a11&xxx
	url = trimQueryParams(url)
	_, ext := SplitFilename(url)
	if toPath == "" {
		toPath = "/tmp"
	}
	//	dir := filepath.Dir(toPath)
	newFilename = NewGuid() + ext
	fullPath := toPath + "/" + newFilename

	// 写到文件中
	file, err := os.Create(fullPath)
	defer file.Close()
	if err != nil {
		return
	}
	file.Write(content)

	path = fullPath
	ok = true
	return
}

// 得到内容
func GetContent(url string) (content []byte, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	Log(err)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	} else {
	}
	if resp == nil || resp.Body == nil || err != nil || resp.StatusCode != http.StatusOK {
		return
	}

	var buf []byte
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		Log(err)
		return
	}

	content = buf
	err = nil
	return
}

// 将url ?, #后面的字符串去掉
func trimQueryParams(url string) string {
	pos := strings.Index(url, "?")
	if pos != -1 {
		url = Substr(url, 0, pos)
	}

	pos = strings.Index(url, "#")
	if pos != -1 {
		url = Substr(url, 0, pos)
	}

	pos = strings.Index(url, "!")
	if pos != -1 {
		url = Substr(url, 0, pos)
	}
	return url
}

// 通过domain得到ip
func GetIpFromDomain(domain string) string {
	ip, _ := net.LookupIP(domain)
	if ip != nil && len(ip) > 0 {
		return ip[0].String()
	}
	return ""
}
