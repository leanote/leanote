package blog

import (
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"html/template"
	"io/ioutil"
	//	"os"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//--------------------
// leanote 自定义主题
// 不使用revel的模板机制
// By life
//--------------------

var ts = []string{"header.html", "footer.html", "highlight.html", "comment.html", "view.html", "404.html"}
var selfTs = []string{"header.html", "footer.html", "index.html", "about_me.html"} // 用户自定义的文件列表

type BlogTpl struct {
	Template    *template.Template
	PathContent map[string]string // path => content
}

func (this *BlogTpl) Content(name string) string {
	return this.PathContent[name]
}

var BlogTplObject *BlogTpl
var CloneTemplate *template.Template

type RenderTemplateResult struct {
	Template    *template.Template
	PathContent map[string]string
	RenderArgs  map[string]interface{}

	IsPreview  bool // 是否是预览
	CurBlogTpl *BlogTpl
}

func parseTemplateError(err error) (templateName string, line int, description string) {
	description = err.Error()
	i := regexp.MustCompile(`:\d+:`).FindStringIndex(description)
	if i != nil {
		line, err = strconv.Atoi(description[i[0]+1 : i[1]-1])
		if err != nil {
		}
		templateName = description[:i[0]]
		if colon := strings.Index(templateName, ":"); colon != -1 {
			templateName = templateName[colon+1:]
		}
		templateName = strings.TrimSpace(templateName)
		description = description[i[1]+1:]
	}
	return templateName, line, description
}
func (r *RenderTemplateResult) render(req *revel.Request, resp *revel.Response, wr io.Writer) {
	err := r.Template.Execute(wr, r.RenderArgs)
	if err == nil {
		return
	}

	var templateContent []string
	templateName, line, description := parseTemplateError(err)
	var content = ""
	if templateName == "" {
		templateName = r.Template.Name()
		content = r.PathContent[templateName]
	} else {
		content = r.PathContent[templateName]
	}
	if content != "" {
		templateContent = strings.Split(content, "\n")
	}

	compileError := &revel.Error{
		Title:       "Template Execution Error",
		Path:        templateName,
		Description: description,
		Line:        line,
		SourceLines: templateContent,
	}

	// 这里, 错误!!
	// 这里应该导向到本主题的错误页面
	resp.Status = 500
	ErrorResult{r.RenderArgs, compileError, r.IsPreview, r.CurBlogTpl}.Apply(req, resp)
}

func (r *RenderTemplateResult) Apply(req *revel.Request, resp *revel.Response) {
	// Handle panics when rendering templates.
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	chunked := revel.Config.BoolDefault("results.chunked", false)

	// If it's a HEAD request, throw away the bytes.
	out := io.Writer(resp.Out)
	if req.Method == "HEAD" {
		out = ioutil.Discard
	}

	// In a prod mode, write the status, render, and hope for the best.
	// (In a dev mode, always render to a temporary buffer first to avoid having
	// error pages distorted by HTML already written)
	if chunked && !revel.DevMode {
		resp.WriteHeader(http.StatusOK, "text/html; charset=utf-8")
		r.render(req, resp, out) // 这里!!!
		return
	}

	// Render the template into a temporary buffer, to see if there was an error
	// rendering the template.  If not, then copy it into the response buffer.
	// Otherwise, template render errors may result in unpredictable HTML (and
	// would carry a 200 status code)
	var b bytes.Buffer
	r.render(req, resp, &b)
	if !chunked {
		resp.Out.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	}
	resp.WriteHeader(http.StatusOK, "text/html; charset=utf-8")
	b.WriteTo(out)
}

// 博客模板
func Init() {
	BlogTplObject = &BlogTpl{PathContent: map[string]string{}}
	BlogTplObject.Template = template.New("blog").Funcs(revel.TemplateFuncs)
	for _, path := range ts {
		fileBytes, _ := ioutil.ReadFile(revel.ViewsPath + "/Blog/" + path)
		fileStr := string(fileBytes)
		path := "blog/" + path
		//		path := path
		BlogTplObject.PathContent[path] = fileStr
		BlogTplObject.Template.New(path).Parse(fileStr) // 以blog为根
	}
	// 复制一份
	CloneTemplate, _ = BlogTplObject.Template.Clone()
}

// name = index.html, search.html, cate.html, page.html
// basePath 表未用户主题的基路径, 如/xxx/public/upload/32323232/themes/theme1, 如果没有, 则表示用自带的
// isPreview 如果是, 错误提示则显示系统的 500 错误详情信息, 供debug
//
func RenderTemplate(name string, args map[string]interface{}, basePath string, isPreview bool) revel.Result {
	var r *RenderTemplateResult
	// 传来的主题路径为空, 则用系统的
	// 都不会为空的
	if basePath == "" {
		path := "blog/" + name
		//		path := name
		t := BlogTplObject.Template.Lookup(path)
		r = &RenderTemplateResult{
			Template:    t,
			PathContent: BlogTplObject.PathContent, // 为了显示错误
			RenderArgs:  args,                      // 把args给它
		}
	} else {
		// 复制一份
		newBlogTplObject := &BlogTpl{}
		var err error
		newBlogTplObject.Template, err = CloneTemplate.Clone() // 复制一份, 为防止多用户出现问题, 因为newBlogTplObject是全局的
		if err != nil {
			return nil
		}
		newBlogTplObject.PathContent = map[string]string{}
		for k, v := range BlogTplObject.PathContent {
			newBlogTplObject.PathContent[k] = v
		}

		// 将该basePath下的所有文件提出
		files := ListDir(basePath)
		for _, t := range files {
			if !strings.Contains(t, ".html") {
				continue
			}
			fileBytes, err := ioutil.ReadFile(basePath + "/" + t)
			if err != nil {
				continue
			}
			fileStr := string(fileBytes)
			newBlogTplObject.PathContent[t] = fileStr
			newBlogTplObject.Template.New(t).Parse(fileStr)
		}

		// 如果本主题下没有, 则用系统的
		t := newBlogTplObject.Template.Lookup(name)

		if t == nil {
			path := "blog/" + name
			t = BlogTplObject.Template.Lookup(path)
		}
		r = &RenderTemplateResult{
			Template:    t,
			PathContent: newBlogTplObject.PathContent, // 为了显示错误
			RenderArgs:  args,
			CurBlogTpl:  newBlogTplObject,
			IsPreview:   isPreview,
		}
	}

	return r
}

////////////////////
//

type ErrorResult struct {
	RenderArgs map[string]interface{}
	Error      error
	IsPreview  bool
	CurBlogTpl *BlogTpl
}

// 错误显示出
func (r ErrorResult) Apply(req *revel.Request, resp *revel.Response) {
	format := req.Format
	status := resp.Status
	if status == 0 {
		status = http.StatusInternalServerError
	}

	contentType := revel.ContentTypeByFilename("xxx." + format)
	if contentType == revel.DefaultFileContentType {
		contentType = "text/plain"
	}

	// Get the error template.
	var err error
	templatePath := fmt.Sprintf("errors/%d.%s", status, format)
	err = nil
	//	tmpl, err := revel.MainTemplateLoader.Template("index.html") // 这里找到错误页面主题

	// This func shows a plaintext error message, in case the template rendering
	// doesn't work.
	showPlaintext := func(err error) {
		revel.PlaintextErrorResult{fmt.Errorf("Server Error:\n%s\n\n"+
			"Additionally, an error occurred when rendering the error page:\n%s",
			r.Error, err)}.Apply(req, resp)
	}

	// 根据是否是preview来得到404模板
	// 是, 则显示系统的错误信息, blog-500.html
	var tmpl *template.Template
	if r.IsPreview {
		tmpl = r.CurBlogTpl.Template.Lookup("blog/404.html")
	} else {
		tmpl = r.CurBlogTpl.Template.Lookup("404.html")
	}
	if tmpl == nil {
		if err == nil {
			err = fmt.Errorf("Couldn't find template %s", templatePath)
		}
		showPlaintext(err)
		return
	}

	// If it's not a revel error, wrap it in one.
	var revelError *revel.Error
	switch e := r.Error.(type) {
	case *revel.Error:
		revelError = e
	case error:
		revelError = &revel.Error{
			Title:       "Server Error",
			Description: e.Error(),
		}
	}

	if revelError == nil {
		panic("no error provided")
	}

	if r.RenderArgs == nil {
		r.RenderArgs = make(map[string]interface{})
	}
	r.RenderArgs["Error"] = revelError
	r.RenderArgs["Router"] = revel.MainRouter

	// 不是preview就不要显示错误了
	if r.IsPreview {
		var b bytes.Buffer
		out := io.Writer(resp.Out)
		//	out = ioutil.Discard
		err = tmpl.Execute(&b, r.RenderArgs)
		resp.WriteHeader(http.StatusOK, "text/html; charset=utf-8")
		b.WriteTo(out)
	}
}
