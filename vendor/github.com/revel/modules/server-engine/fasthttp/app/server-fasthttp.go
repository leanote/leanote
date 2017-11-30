package fasthttp

import (
	"bytes"
	"fmt"
	"github.com/revel/revel"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type FastHTTPServer struct {
	Server           *fasthttp.Server
	ServerInit       *revel.EngineInit
	MaxMultipartSize int64
}

var serverLog = revel.AppLog

func init() {
	revel.RegisterServerEngine("fasthttp", func() revel.ServerEngine {
		return &FastHTTPServer{}
	})
	revel.RegisterModuleInit(func(m *revel.Module) {
		serverLog = m.Log
	})
}

func (f *FastHTTPServer) Init(init *revel.EngineInit) {
	f.MaxMultipartSize = int64(revel.Config.IntDefault("server.request.max.multipart.filesize", 32)) << 20 /* 32 MB */
	fastHttpContextStack = revel.NewStackLock(revel.Config.IntDefault("server.context.stack", 100),
		revel.Config.IntDefault("server.context.maxstack", 200),
		func() interface{} { return NewFastHttpContext(f) })
	fastHttpMultipartFormStack = revel.NewStackLock(revel.Config.IntDefault("server.form.stack", 100),
		revel.Config.IntDefault("server.form.maxstack", 200),
		func() interface{} { return &FastHttpMultipartForm{} })

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		f.RequestHandler(ctx)
	}
	f.ServerInit = init
	f.Server = &fasthttp.Server{
		ReadTimeout:  time.Duration(revel.Config.IntDefault("http.timeout.read", 0)) * time.Second,
		WriteTimeout: time.Duration(revel.Config.IntDefault("http.timeout.write", 0)) * time.Second,
		Handler:      requestHandler,
	}

}

// Handler is assigned in the Init
func (f *FastHTTPServer) Start() {
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("\nListening on fasthttp %s...\n", f.ServerInit.Address)
	}()
	if revel.HTTPSsl {
		if f.ServerInit.Network != "tcp" {
			// This limitation is just to reduce complexity, since it is standard
			// to terminate SSL upstream when using unix domain sockets.
			serverLog.Fatal("SSL is only supported for TCP sockets. Specify a port to listen on.")
		}
		serverLog.Fatal("Failed to listen https:", "error",
			f.Server.ListenAndServeTLS(f.ServerInit.Address, revel.HTTPSslCert, revel.HTTPSslKey))
	} else {
		listener, err := net.Listen(f.ServerInit.Network, f.ServerInit.Address)
		if err != nil {
			serverLog.Fatal("Failed to listen http:", "error", err)
		}
		serverLog.Info("Listening fasthttp ", f.ServerInit.Network, f.ServerInit.Address)
		// revel.ERROR.Fatalln("Failed to serve:", f.Server.ListenAndServe(f.ServerInit.Address))
		serverLog.Fatal("Server exited", "error", f.Server.Serve(listener))
		println("***ENDING ***")
	}

}
func (f *FastHTTPServer) RequestHandler(ctx *fasthttp.RequestCtx) {
	// TODO this
	//if maxRequestSize := int64(revel.Config.IntDefault("http.maxrequestsize", 0)); maxRequestSize > 0 {
	//   buffer := &bytes.Buffer{}
	//   err := ctx.Request.ReadLimitBody(buffer,maxRequestSize)
	//   if err!=nil {
	//       // Send the error back to the client
	//       ctx.SetStatusCode(http.StatusRequestEntityTooLarge)
	//       return
	//   }
	//}
	context := fastHttpContextStack.Pop().(*FastHttpContext)
	defer func() {
		fastHttpContextStack.Push(context)
	}()
	context.SetContext(ctx)
	f.ServerInit.Callback(context)
}

func (f *FastHTTPServer) Event(event int, args interface{}) {

	switch event {
	case revel.ENGINE_BEFORE_INITIALIZED:
	case revel.ENGINE_STARTED:

	}

}
func (f *FastHTTPServer) Name() string {
	return "fasthttp"
}
func (f *FastHTTPServer) Engine() interface{} {
	return f
}
func (g *FastHTTPServer) Stats() map[string]interface{} {
	return map[string]interface{}{
		"FastHTTP Engine Context": fastHttpContextStack.String(),
		"FastHTTP Engine Forms":   fastHttpMultipartFormStack.String(),
	}
}

type (
	FastHttpContext struct {
		Request  *FastHttpRequest
		Response *FastHttpResponse
	}

	FastHttpRequest struct {
		toQuery         bool
		url             *url.URL
		query           url.Values
		Original        *fasthttp.RequestCtx
		FormParsed      bool
		form            url.Values
		MultiFormParsed bool
		// WebSocket       *websocket.Conn
		ParsedForm *FastHttpMultipartForm
		header     *FastHttpHeader
		Engine     *FastHTTPServer
	}

	FastHttpResponse struct {
		Original *fasthttp.RequestCtx
		header   *FastHttpHeader
		Writer   io.Writer
		Engine   *FastHTTPServer
	}
	FastHttpMultipartForm struct {
		Form *multipart.Form
	}
	FastHttpHeader struct {
		Source     interface{}
		isResponse bool
	}
	FastHttpCookie []byte
)

var (
	fastHttpContextStack       *revel.SimpleLockStack
	fastHttpMultipartFormStack *revel.SimpleLockStack
)

func NewFastHttpContext(instance *FastHTTPServer) *FastHttpContext {
	if instance == nil {
		instance = &FastHTTPServer{MaxMultipartSize: 32 << 20}
	}
	c := &FastHttpContext{
		Request: &FastHttpRequest{header: &FastHttpHeader{isResponse: false},
			Engine: instance},
		Response: &FastHttpResponse{header: &FastHttpHeader{isResponse: true},
			Engine: instance},
	}
	c.Response.header.Source = c.Response
	c.Request.header.Source = c.Request
	return c
}
func (c *FastHttpContext) GetRequest() revel.ServerRequest {
	return c.Request
}
func (c *FastHttpContext) GetResponse() revel.ServerResponse {
	return c.Response
}
func (c *FastHttpContext) SetContext(context *fasthttp.RequestCtx) {
	c.Response.SetContext(context)
	c.Request.SetContext(context)
}
func (c *FastHttpContext) Destroy() {
	c.Response.Destroy()
	c.Request.Destroy()
}

func (r *FastHttpRequest) Get(key int) (value interface{}, err error) {
	switch key {
	case revel.HTTP_SERVER_HEADER:
		value = r.GetHeader()
	case revel.HTTP_MULTIPART_FORM:
		value, err = r.GetMultipartForm()
	case revel.HTTP_QUERY:
		value = r.GetQuery()
	case revel.HTTP_FORM:
		value, err = r.GetForm()
	case revel.HTTP_REQUEST_URI:
		value = string(r.Original.RequestURI())
	case revel.HTTP_REMOTE_ADDR:
		value = r.Original.RemoteAddr().String()
	case revel.HTTP_METHOD:
		value = string(r.Original.Method())
	case revel.HTTP_PATH:
		value = string(r.Original.Path())
	case revel.HTTP_HOST:
		value = string(r.Original.Request.Host())
	case revel.HTTP_URL:
		if r.url == nil {
			r.url, _ = url.Parse(string(r.Original.Request.URI().FullURI()))
		}
		value = r.url
	case revel.HTTP_BODY:
		value = bytes.NewBuffer(r.Original.Request.Body())
	default:
		err = revel.ENGINE_UNKNOWN_GET
	}

	return
}
func (r *FastHttpRequest) Set(key int, value interface{}) bool {
	return false
}

func (r *FastHttpRequest) GetQuery() url.Values {
	if !r.toQuery {
		// Attempt to convert to query
		r.query = url.Values{}
		r.Original.QueryArgs().VisitAll(func(key, value []byte) {
			r.query.Set(string(key), string(value))
		})
		r.toQuery = true
	}
	return r.query
}
func (r *FastHttpRequest) GetForm() (url.Values, error) {
	if !r.FormParsed {
		r.form = url.Values{}
		r.Original.PostArgs().VisitAll(func(key, value []byte) {
			println("Set value", string(key), string(value))
			r.query.Set(string(key), string(value))
		})
		r.FormParsed = true
	}
	return r.form, nil
}
func (r *FastHttpRequest) GetMultipartForm() (revel.ServerMultipartForm, error) {
	if !r.MultiFormParsed {
		// TODO Limit size r.Engine.MaxMultipartSize
		form, err := r.Original.MultipartForm()
		if err != nil {
			return nil, err
		}

		r.ParsedForm = fastHttpMultipartFormStack.Pop().(*FastHttpMultipartForm)
		r.ParsedForm.Form = form
	}

	return r.ParsedForm, nil
}
func (r *FastHttpRequest) GetHeader() revel.ServerHeader {
	return r.header
}
func (r *FastHttpRequest) GetRaw() interface{} {
	return r.Original
}
func (r *FastHttpRequest) SetContext(req *fasthttp.RequestCtx) {
	r.Original = req

}
func (r *FastHttpRequest) Destroy() {
	r.Original = nil
	r.FormParsed = false
	r.MultiFormParsed = false
	r.ParsedForm = nil
	r.toQuery = false

}

func (r *FastHttpResponse) Get(key int) (value interface{}, err error) {
	switch key {
	case revel.HTTP_SERVER_HEADER:
		value = r.Header()
	case revel.HTTP_STREAM_WRITER:
		value = r
	case revel.HTTP_WRITER:
		value = r.Writer
	default:
		err = revel.ENGINE_UNKNOWN_GET
	}
	return
}

func (r *FastHttpResponse) Set(key int, value interface{}) (set bool) {
	switch key {
	case revel.ENGINE_RESPONSE_STATUS:
		r.Header().SetStatus(value.(int))
		set = true
	case revel.HTTP_WRITER:
		r.SetWriter(value.(io.Writer))
		set = true
	}
	return
}

func (r *FastHttpResponse) GetWriter() io.Writer {
	return r.Writer
}
func (r *FastHttpResponse) Header() revel.ServerHeader {
	return r.header
}
func (r *FastHttpResponse) GetRaw() interface{} {
	return r.Original
}
func (r *FastHttpResponse) WriteStream(name string, contentlen int64, modtime time.Time, reader io.Reader) error {

	// do a simple io.Copy, we do it directly into the writer which may be configured to be a compressed
	// writer
	ius := r.Original.Request.Header.Peek("If-Unmodified-Since")
	if t, err := http.ParseTime(string(ius)); ius != nil && err == nil && !modtime.IsZero() {
		// The Date-Modified header truncates sub-second precision, so
		// use mtime < t+1s instead of mtime <= t to check for unmodified.
		if modtime.Before(t.Add(1 * time.Second)) {
			h := r.Original.Response.Header
			h.Del("Content-Type")
			h.Del("Content-Length")
			if h.Peek("Etag") != nil {
				h.Del("Last-Modified")
			}
			h.SetStatusCode(http.StatusNotModified)
			return nil
		}
	}

	if contentlen != -1 {
		r.Original.Response.Header.Set("Content-Length", strconv.FormatInt(contentlen, 10))
	}
	if _, err := io.Copy(r.Writer, reader); err != nil {
		r.Original.Response.Header.SetStatusCode(http.StatusInternalServerError)
		return err
	} else {
		r.Original.Response.Header.SetStatusCode(http.StatusOK)
	}

	return nil
}
func (r *FastHttpResponse) Destroy() {
	if c, ok := r.Writer.(io.Closer); ok {
		c.Close()
	}
	r.Original = nil
	r.Writer = nil

}
func (r *FastHttpResponse) SetContext(w *fasthttp.RequestCtx) {
	r.Original = w
	r.Writer = w.Response.BodyWriter()
}
func (r *FastHttpResponse) SetWriter(writer io.Writer) {
	r.Writer = writer
}
func (r *FastHttpHeader) SetCookie(cookie string) {
	if r.isResponse {
		r.Source.(*FastHttpResponse).Original.Response.Header.Add("Set-Cookie", cookie)
	}
}
func (r *FastHttpHeader) GetCookie(key string) (value revel.ServerCookie, err error) {
	if !r.isResponse {
		var cookie []byte
		if cookie = r.Source.(*FastHttpRequest).Original.Request.Header.Cookie(key); cookie != nil {
			value = FastHttpCookie(cookie)
		} else {
			err = http.ErrNoCookie
		}

	}
	return
}
func (r *FastHttpHeader) Set(key string, value string) {
	if r.isResponse {
		r.Source.(*FastHttpResponse).Original.Response.Header.Set(key, value)
	}
}
func (r *FastHttpHeader) Add(key string, value string) {
	if r.isResponse {
		r.Source.(*FastHttpResponse).Original.Response.Header.Add(key, value)
	}
}
func (r *FastHttpHeader) Del(key string) {
	if r.isResponse {
		r.Source.(*FastHttpResponse).Original.Response.Header.Del(key)
	}
}
func (r *FastHttpHeader) Get(key string) (value []string) {
	if !r.isResponse {
		value = strings.Split(string(r.Source.(*FastHttpRequest).Original.Request.Header.Peek(key)), ",")
	} else {
		value = strings.Split(string(r.Source.(*FastHttpResponse).Original.Response.Header.Peek(key)), ",")
	}
	return
}
func (r *FastHttpHeader) SetStatus(statusCode int) {
	if r.isResponse {
		r.Source.(*FastHttpResponse).Original.Response.SetStatusCode(statusCode)
	}
}
func (r FastHttpCookie) GetValue() string {
	return string(r)
}
func (f *FastHttpMultipartForm) GetFiles() map[string][]*multipart.FileHeader {
	return f.Form.File
}
func (f *FastHttpMultipartForm) GetValues() url.Values {
	return url.Values(f.Form.Value)
}
func (f *FastHttpMultipartForm) RemoveAll() error {
	return f.Form.RemoveAll()
}
