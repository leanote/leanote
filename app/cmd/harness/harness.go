// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

// Package harness for a Revel Framework.
//
// It has a following responsibilities:
// 1. Parse the user program, generating a main.go file that registers
//    controller classes and starts the user's server.
// 2. Build and run the user program.  Show compile errors.
// 3. Monitor the user source and re-build / restart the program when necessary.
//
// Source files are generated in the app/tmp directory.
package harness

import (
	"crypto/tls"
	"fmt"
	"time"
	"go/build"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/revel/cmd/model"
	"github.com/revel/cmd/utils"
	"github.com/revel/cmd/watcher"
	"html/template"
	"io/ioutil"
	"sync"
	"encoding/json"
)

var (
	doNotWatch = []string{"tmp", "views", "routes"}

	lastRequestHadError int32
)

// Harness reverse proxies requests to the application server.
// It builds / runs / rebuilds / restarts the server when code is changed.
type Harness struct {
	app        *App                   // The application
	useProxy   bool                   // True if proxy is in use
	serverHost string                 // The proxy server host
	port       int                    // The proxy serber port
	proxy      *httputil.ReverseProxy // The proxy
	watcher    *watcher.Watcher       // The file watched
	mutex      *sync.Mutex            // A mutex to prevent concurrent updates
	paths      *model.RevelContainer  // The Revel container
	config     *model.CommandConfig   // The configuration
	runMode    string                 // The runmode the harness is running in
	isError    bool                   // True if harness is in error state
	ranOnce    bool                   // True app compiled once
}

func (h *Harness) renderError(iw http.ResponseWriter, ir *http.Request, err error) {
	// Render error here
	// Grab the template from three places
	// 1) Application/views/errors
	// 2) revel_home/views/errors
	// 3) views/errors
	if err == nil {
		utils.Logger.Panic("Caller passed in a nil error")
	}
	templateSet := template.New("__root__")
	seekViewOnPath := func(view string) (path string) {
		path = filepath.Join(h.paths.ViewsPath, "errors", view)
		if !utils.Exists(path) {
			path = filepath.Join(h.paths.RevelPath, "templates", "errors", view)
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			utils.Logger.Error("Unable to read template file", path)
		}
		_, err = templateSet.New("errors/" + view).Parse(string(data))
		if err != nil {
			utils.Logger.Error("Unable to parse template file", path)
		}
		return
	}
	target := []string{seekViewOnPath("500.html"), seekViewOnPath("500-dev.html")}
	if !utils.Exists(target[0]) {
		fmt.Fprintf(iw, "Target template not found not found %s<br />\n", target[0])
		fmt.Fprintf(iw, "An error ocurred %s", err.Error())
		return
	}
	var revelError *utils.SourceError
	switch e := err.(type) {
	case *utils.SourceError:
		revelError = e
	case error:
		revelError = &utils.SourceError{
			Title:       "Server Error",
			Description: e.Error(),
		}
	}

	if revelError == nil {
		panic("no error provided")
	}
	viewArgs := map[string]interface{}{}
	viewArgs["RunMode"] = h.paths.RunMode
	viewArgs["DevMode"] = h.paths.DevMode
	viewArgs["Error"] = revelError

	// Render the template from the file
	err = templateSet.ExecuteTemplate(iw, "errors/500.html", viewArgs)
	if err != nil {
		utils.Logger.Error("Failed to execute", "error", err)
	}
}

// ServeHTTP handles all requests.
// It checks for changes to app, rebuilds if necessary, and forwards the request.
func (h *Harness) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Don't rebuild the app for favicon requests.
	if lastRequestHadError > 0 && r.URL.Path == "/favicon.ico" {
		return
	}

	// Flush any change events and rebuild app if necessary.
	// Render an error page if the rebuild / restart failed.
	err := h.watcher.Notify()
	if err != nil {
		// In a thread safe manner update the flag so that a request for
		// /favicon.ico does not trigger a rebuild
		atomic.CompareAndSwapInt32(&lastRequestHadError, 0, 1)
		h.renderError(w, r, err)
		return
	}

	// In a thread safe manner update the flag so that a request for
	// /favicon.ico is allowed
	atomic.CompareAndSwapInt32(&lastRequestHadError, 1, 0)

	// Reverse proxy the request.
	// (Need special code for websockets, courtesy of bradfitz)
	if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		h.proxyWebsocket(w, r, h.serverHost)
	} else {
		h.proxy.ServeHTTP(w, r)
	}
}

// NewHarness method returns a reverse proxy that forwards requests
// to the given port.
func NewHarness(c *model.CommandConfig, paths *model.RevelContainer, runMode string, noProxy bool) *Harness {
	// Get a template loader to render errors.
	// Prefer the app's views/errors directory, and fall back to the stock error pages.
	//revel.MainTemplateLoader = revel.NewTemplateLoader(
	//	[]string{filepath.Join(revel.RevelPath, "templates")})
	//if err := revel.MainTemplateLoader.Refresh(); err != nil {
	//	revel.RevelLog.Error("Template loader error", "error", err)
	//}

	addr := paths.HTTPAddr
	port := paths.Config.IntDefault("harness.port", 0)
	scheme := "http"

	if paths.HTTPSsl {
		scheme = "https"
	}

	// If the server is running on the wildcard address, use "localhost"
	if addr == "" {
		utils.Logger.Warn("No http.addr specified in the app.conf listening on localhost interface only. " +
			"This will not allow external access to your application")
		addr = "localhost"
	}

	if port == 0 {
		port = getFreePort()
	}

	serverURL, _ := url.ParseRequestURI(fmt.Sprintf(scheme+"://%s:%d", addr, port))

	serverHarness := &Harness{
		port:       port,
		serverHost: serverURL.String()[len(scheme+"://"):],
		proxy:      httputil.NewSingleHostReverseProxy(serverURL),
		mutex:      &sync.Mutex{},
		paths:      paths,
		useProxy:   !noProxy,
		config:     c,
		runMode:    runMode,
	}

	if paths.HTTPSsl {
		serverHarness.proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return serverHarness
}

// Refresh method rebuilds the Revel application and run it on the given port.
// called by the watcher
func (h *Harness) Refresh() (err *utils.SourceError) {
	t  := time.Now();
	fmt.Println("Changed detected, recompiling")
	err = h.refresh()
	if err!=nil && !h.ranOnce && h.useProxy {
		addr := fmt.Sprintf("%s:%d", h.paths.HTTPAddr, h.paths.HTTPPort)

		fmt.Printf("\nError compiling code, to view error details see proxy running on http://%s\n\n",addr)
	}

	h.ranOnce = true
	fmt.Printf("\nTime to recompile %s\n",time.Now().Sub(t).String())
	return
}

func (h *Harness) refresh() (err *utils.SourceError) {
	// Allow only one thread to rebuild the process
	// If multiple requests to rebuild are queued only the last one is executed on
	// So before a build is started we wait for a second to determine if
	// more requests for a build are triggered.
	// Once no more requests are triggered the build will be processed
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.app != nil {
		h.app.Kill()
	}

	utils.Logger.Info("Rebuild Called")
	var newErr error
	h.app, newErr = Build(h.config, h.paths)
	if newErr != nil {
		utils.Logger.Error("Build detected an error", "error", newErr)
		if castErr, ok := newErr.(*utils.SourceError); ok {
			return castErr
		}
		err = &utils.SourceError{
			Title:       "App failed to start up",
			Description: err.Error(),
		}
		return
	}

	if h.useProxy {
		h.app.Port = h.port
		runMode := h.runMode
		if !h.config.HistoricMode {
			// Recalulate run mode based on the config
			var paths []byte
			if len(h.app.PackagePathMap)>0 {
				paths, _ = json.Marshal(h.app.PackagePathMap)
			}
			runMode = fmt.Sprintf(`{"mode":"%s", "specialUseFlag":%v,"packagePathMap":%s}`, h.app.Paths.RunMode, h.config.Verbose, string(paths))

		}
		if err2 := h.app.Cmd(runMode).Start(h.config); err2 != nil {
			utils.Logger.Error("Could not start application", "error", err2)
			if err,k :=err2.(*utils.SourceError);k {
				return err
			}
			return &utils.SourceError{
				Title:       "App failed to start up",
				Description: err2.Error(),
			}
		}
	} else {
		h.app = nil
	}

	return
}

// WatchDir method returns false to file matches with doNotWatch
// otheriwse true
func (h *Harness) WatchDir(info os.FileInfo) bool {
	return !utils.ContainsString(doNotWatch, info.Name())
}

// WatchFile method returns true given filename HasSuffix of ".go"
// otheriwse false - implements revel.DiscerningListener
func (h *Harness) WatchFile(filename string) bool {
	return strings.HasSuffix(filename, ".go")
}

// Run the harness, which listens for requests and proxies them to the app
// server, which it runs and rebuilds as necessary.
func (h *Harness) Run() {
	var paths []string
	if h.paths.Config.BoolDefault("watch.gopath", false) {
		gopaths := filepath.SplitList(build.Default.GOPATH)
		paths = append(paths, gopaths...)
	}
	paths = append(paths, h.paths.CodePaths...)
	h.watcher = watcher.NewWatcher(h.paths, false)
	h.watcher.Listen(h, paths...)
	go h.Refresh()
	// h.watcher.Notify()

	if h.useProxy {
		go func() {
			// Check the port to start on a random port
			if h.paths.HTTPPort == 0 {
				h.paths.HTTPPort = getFreePort()
			}
			addr := fmt.Sprintf("%s:%d", h.paths.HTTPAddr, h.paths.HTTPPort)
			utils.Logger.Infof("Proxy server is listening on %s", addr)


			var err error
			if h.paths.HTTPSsl {
				err = http.ListenAndServeTLS(
					addr,
					h.paths.HTTPSslCert,
					h.paths.HTTPSslKey,
					h)
			} else {
				err = http.ListenAndServe(addr, h)
			}
			if err != nil {
				utils.Logger.Error("Failed to start reverse proxy:", "error", err)
			}
		}()

	}

	// Make a new channel to listen for the interrupt event
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	// Kill the app and exit
	if h.app != nil {
		h.app.Kill()
	}
	os.Exit(1)
}

// Find an unused port
func getFreePort() (port int) {
	conn, err := net.Listen("tcp", ":0")
	if err != nil {
		utils.Logger.Fatal("Unable to fetch a freee port address", "error", err)
	}

	port = conn.Addr().(*net.TCPAddr).Port
	err = conn.Close()
	if err != nil {
		utils.Logger.Fatal("Unable to close port", "error", err)
	}
	return port
}

// proxyWebsocket copies data between websocket client and server until one side
// closes the connection.  (ReverseProxy doesn't work with websocket requests.)
func (h *Harness) proxyWebsocket(w http.ResponseWriter, r *http.Request, host string) {
	var (
		d   net.Conn
		err error
	)
	if h.paths.HTTPSsl {
		// since this proxy isn't used in production,
		// it's OK to set InsecureSkipVerify to true
		// no need to add another configuration option.
		d, err = tls.Dial("tcp", host, &tls.Config{InsecureSkipVerify: true})
	} else {
		d, err = net.Dial("tcp", host)
	}
	if err != nil {
		http.Error(w, "Error contacting backend server.", 500)
		utils.Logger.Error("Error dialing websocket backend ", "host", host, "error", err)
		return
	}
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Not a hijacker?", 500)
		return
	}
	nc, _, err := hj.Hijack()
	if err != nil {
		utils.Logger.Error("Hijack error", "error", err)
		return
	}
	defer func() {
		if err = nc.Close(); err != nil {
			utils.Logger.Error("Connection close error", "error", err)
		}
		if err = d.Close(); err != nil {
			utils.Logger.Error("Dial close error", "error", err)
		}
	}()

	err = r.Write(d)
	if err != nil {
		utils.Logger.Error("Error copying request to target", "error", err)
		return
	}

	errc := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errc <- err
	}
	go cp(d, nc)
	go cp(nc, d)
	<-errc
}
