// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"go/build"
	"strconv"

	"github.com/leanote/leanote/cmd/harness"
	"github.com/revel/revel"
	// "fmt"

	"path/filepath"
)

var cmdRun = &Command{
	UsageLine: "run [import path] [run mode] [port]",
	Short:     "run a Revel application",
	Long: `
Run the Revel web application named by the given import path.

For example, to run the chat room sample application:

    revel run github.com/revel/examples/chat dev

The run mode is used to select which set of app.conf configuration should
apply and may be used to determine logic in the application itself.

Run mode defaults to "dev".

You can set a port as an optional third parameter.  For example:

    revel run github.com/revel/examples/chat prod 8080`,
}

// RunArgs holds revel run parameters
type RunArgs struct {
	ImportPath string
	Mode       string
	Port       int
}

func init() {
	cmdRun.Run = runApp
}

func parseRunArgs(args []string) *RunArgs {
	inputArgs := RunArgs{
		ImportPath: importPathFromCurrentDir(),
		Mode:       DefaultRunMode,
		Port:       revel.HTTPPort,
	}
	switch len(args) {
	case 3:
		// Possibile combinations
		// revel run [import-path] [run-mode] [port]
		port, err := strconv.Atoi(args[2])
		if err != nil {
			errorf("Failed to parse port as integer: %s", args[2])
		}
		inputArgs.ImportPath = args[0]
		inputArgs.Mode = args[1]
		inputArgs.Port = port
	case 2:
		// Possibile combinations
		// 1. revel run [import-path] [run-mode]
		// 2. revel run [import-path] [port]
		// 3. revel run [run-mode] [port]
		if _, err := build.Import(args[0], "", build.FindOnly); err == nil {
			// 1st arg is the import path
			inputArgs.ImportPath = args[0]
			if port, err := strconv.Atoi(args[1]); err == nil {
				// 2nd arg is the port number
				inputArgs.Port = port
			} else {
				// 2nd arg is the run mode
				inputArgs.Mode = args[1]
			}
		} else {
			// 1st arg is the run mode
			port, err := strconv.Atoi(args[1])
			if err != nil {
				errorf("Failed to parse port as integer: %s", args[1])
			}
			inputArgs.Mode = args[0]
			inputArgs.Port = port
		}
	case 1:
		// Possibile combinations
		// 1. revel run [import-path]
		// 2. revel run [port]
		// 3. revel run [run-mode]
		if _, err := build.Import(args[0], "", build.FindOnly); err == nil {
			// 1st arg is the import path
			inputArgs.ImportPath = args[0]
		} else if port, err := strconv.Atoi(args[0]); err == nil {
			// 1st arg is the port number
			inputArgs.Port = port
		} else {
			// 1st arg is the run mode
			inputArgs.Mode = args[0]
		}
	}

	return &inputArgs
}


// findSrcPaths uses the "go/build" package to find the source root for Revel
// and the app.
func findSrcPaths(importPath string) (appSourcePath string) {
	var (
		gopaths = filepath.SplitList(build.Default.GOPATH)
		goroot  = build.Default.GOROOT
	)

	if len(gopaths) == 0 {
		revel.ERROR.Fatalln("GOPATH environment variable is not set. ",
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.")
	}

	if revel.ContainsString(gopaths, goroot) {
		revel.ERROR.Fatalf("GOPATH (%s) must not include your GOROOT (%s). "+
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.",
			gopaths, goroot)
	}

	appPkg, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		revel.ERROR.Fatalln("Failed to import", importPath, "with error:", err)
	}

	return appPkg.SrcRoot
}

func runApp(args []string) {
	runArgs := parseRunArgs(args)

	// Find and parse app.conf
	// fmt.Println(runArgs.ImportPath + "/vendor")
	// revel.Init(runArgs.Mode, runArgs.ImportPath, runArgs.ImportPath + "/vendor")
	srcPath := findSrcPaths(runArgs.ImportPath)
	srcPath = ""
	revel.Init(runArgs.Mode, runArgs.ImportPath, srcPath)
	revel.LoadMimeConfig()

	// fallback to default port
	if runArgs.Port == 0 {
		runArgs.Port = revel.HTTPPort
	}

	revel.INFO.Printf("Running %s (%s) in %s mode\n", revel.AppName, revel.ImportPath, runArgs.Mode)
	revel.TRACE.Println("Base path:", revel.BasePath)

	// If the app is run in "watched" mode, use the harness to run it.
	if revel.Config.BoolDefault("watch", true) && revel.Config.BoolDefault("watch.code", true) {
		revel.TRACE.Println("Running in watched mode.")
		revel.HTTPPort = runArgs.Port
		harness.NewHarness().Run() // Never returns.
	}

	// Else, just build and run the app.
	revel.TRACE.Println("Running in live build mode.")
	app, err := harness.Build()
	if err != nil {
		errorf("Failed to build app: %s", err)
	}
	app.Port = runArgs.Port
	app.Cmd().Run()
}
