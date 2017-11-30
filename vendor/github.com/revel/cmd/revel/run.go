// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"go/build"
	"strconv"

	"github.com/revel/cmd/harness"
	"github.com/revel/revel"
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
		// Possible combinations
		// revel run [import-path] [run-mode] [port]
		port, err := strconv.Atoi(args[2])
		if err != nil {
			errorf("Failed to parse port as integer: %s", args[2])
		}
		inputArgs.ImportPath = args[0]
		inputArgs.Mode = args[1]
		inputArgs.Port = port
	case 2:
		// Possible combinations
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
		// Possible combinations
		// 1. revel run [import-path]
		// 2. revel run [port]
		// 3. revel run [run-mode]
		_, err := build.Import(args[0], "", build.FindOnly)
		if err != nil {
			revel.RevelLog.Warn("Unable to run using an import path, assuming import path is working directory %s %s", "Argument", args[0], "error", err.Error())
		}
		println("Trying to build with", args[0], err)
		if err == nil {
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

func runApp(args []string) {
	runArgs := parseRunArgs(args)

	// Find and parse app.conf
	revel.Init(runArgs.Mode, runArgs.ImportPath, "")
	revel.LoadMimeConfig()

	// fallback to default port
	if runArgs.Port == 0 {
		runArgs.Port = revel.HTTPPort
	}

	revel.RevelLog.Infof("Running %s (%s) in %s mode\n", revel.AppName, revel.ImportPath, runArgs.Mode)
	revel.RevelLog.Debug("Base path:", "path", revel.BasePath)

	// If the app is run in "watched" mode, use the harness to run it.
	if revel.Config.BoolDefault("watch", true) && revel.Config.BoolDefault("watch.code", true) {
		revel.RevelLog.Debug("Running in watched mode.")
		revel.HTTPPort = runArgs.Port
		harness.NewHarness().Run() // Never returns.
	}

	// Else, just build and run the app.
	revel.RevelLog.Debug("Running in live build mode.")
	app, err := harness.Build()
	if err != nil {
		errorf("Failed to build app: %s", err)
	}
	app.Port = runArgs.Port
	app.Cmd().Run()
}
