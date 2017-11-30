// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"go/build"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/revel/revel"
)

var cmdNew = &Command{
	UsageLine: "new [path] [skeleton]",
	Short:     "create a skeleton Revel application",
	Long: `
New creates a few files to get a new Revel application running quickly.

It puts all of the files in the given import path, taking the final element in
the path to be the app name.

Skeleton is an optional argument, provided as an import path

For example:

    revel new import/path/helloworld

    revel new import/path/helloworld import/path/skeleton
`,
}

func init() {
	cmdNew.Run = newApp
}

var (

	// go related paths
	gopath  string
	gocmd   string
	srcRoot string

	// revel related paths
	revelPkg     *build.Package
	revelCmdPkg  *build.Package
	appPath      string
	appName      string
	basePath     string
	importPath   string
	skeletonPath string
)

func newApp(args []string) {
	// check for proper args by count
	if len(args) == 0 {
		errorf("No import path given.\nRun 'revel help new' for usage.\n")
	}
	if len(args) > 2 {
		errorf("Too many arguments provided.\nRun 'revel help new' for usage.\n")
	}

	revel.ERROR.SetFlags(log.LstdFlags)

	// checking and setting go paths
	initGoPaths()

	// checking and setting application
	setApplicationPath(args)

	// checking and setting skeleton
	setSkeletonPath(args)

	// copy files to new app directory
	copyNewAppFiles()

	// goodbye world
	fmt.Fprintln(os.Stdout, "Your application is ready:\n  ", appPath)
	fmt.Fprintln(os.Stdout, "\nYou can run it with:\n   revel run", importPath)
}

const alphaNumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func generateSecret() string {
	chars := make([]byte, 64)
	for i := 0; i < 64; i++ {
		chars[i] = alphaNumeric[rand.Intn(len(alphaNumeric))]
	}
	return string(chars)
}

// lookup and set Go related variables
func initGoPaths() {
	// lookup go path
	gopath = build.Default.GOPATH
	if gopath == "" {
		errorf("Abort: GOPATH environment variable is not set. " +
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.")
	}

	// check for go executable
	var err error
	gocmd, err = exec.LookPath("go")
	if err != nil {
		errorf("Go executable not found in PATH.")
	}

	// revel/revel#1004 choose go path relative to current working directory
	workingDir, _ := os.Getwd()
	goPathList := filepath.SplitList(gopath)
	for _, path := range goPathList {
		if strings.HasPrefix(strings.ToLower(workingDir), strings.ToLower(path)) {
			srcRoot = path
			break
		}

		path, _ = filepath.EvalSymlinks(path)
		if len(path) > 0 && strings.HasPrefix(strings.ToLower(workingDir), strings.ToLower(path)) {
			srcRoot = path
			break
		}
	}

	if len(srcRoot) == 0 {
		revel.ERROR.Fatalln("Abort: could not create a Revel application outside of GOPATH.")
	}

	// set go src path
	srcRoot = filepath.Join(srcRoot, "src")
}

func setApplicationPath(args []string) {
	var err error
	importPath = args[0]

	// revel/revel#1014 validate relative path, we cannot use built-in functions
	// since Go import path is valid relative path too.
	// so check basic part of the path, which is "."
	if filepath.IsAbs(importPath) || strings.HasPrefix(importPath, ".") {
		errorf("Abort: '%s' looks like a directory.  Please provide a Go import path instead.",
			importPath)
	}

	_, err = build.Import(importPath, "", build.FindOnly)
	if err == nil {
		errorf("Abort: Import path %s already exists.\n", importPath)
	}

	revelPkg, err = build.Import(revel.RevelImportPath, "", build.FindOnly)
	if err != nil {
		errorf("Abort: Could not find Revel source code: %s\n", err)
	}

	appPath = filepath.Join(srcRoot, filepath.FromSlash(importPath))
	appName = filepath.Base(appPath)
	basePath = filepath.ToSlash(filepath.Dir(importPath))

	if basePath == "." {
		// we need to remove the a single '.' when
		// the app is in the $GOROOT/src directory
		basePath = ""
	} else {
		// we need to append a '/' when the app is
		// is a subdirectory such as $GOROOT/src/path/to/revelapp
		basePath += "/"
	}
}

func setSkeletonPath(args []string) {
	var err error
	if len(args) == 2 { // user specified
		skeletonName := args[1]
		_, err = build.Import(skeletonName, "", build.FindOnly)
		if err != nil {
			// Execute "go get <pkg>"
			getCmd := exec.Command(gocmd, "get", "-d", skeletonName)
			fmt.Println("Exec:", getCmd.Args)
			getOutput, err := getCmd.CombinedOutput()

			// check getOutput for no buildible string
			bpos := bytes.Index(getOutput, []byte("no buildable Go source files in"))
			if err != nil && bpos == -1 {
				errorf("Abort: Could not find or 'go get' Skeleton  source code: %s\n%s\n", getOutput, skeletonName)
			}
		}
		// use the
		skeletonPath = filepath.Join(srcRoot, skeletonName)

	} else {
		// use the revel default
		revelCmdPkg, err = build.Import(RevelCmdImportPath, "", build.FindOnly)
		if err != nil {
			errorf("Abort: Could not find Revel Cmd source code: %s\n", err)
		}

		skeletonPath = filepath.Join(revelCmdPkg.Dir, "revel", "skeleton")
	}
}

func copyNewAppFiles() {
	var err error
	err = os.MkdirAll(appPath, 0777)
	panicOnError(err, "Failed to create directory "+appPath)

	_ = mustCopyDir(appPath, skeletonPath, map[string]interface{}{
		// app.conf
		"AppName":  appName,
		"BasePath": basePath,
		"Secret":   generateSecret(),
	})

	// Dotfiles are skipped by mustCopyDir, so we have to explicitly copy the .gitignore.
	gitignore := ".gitignore"
	mustCopyFile(filepath.Join(appPath, gitignore), filepath.Join(skeletonPath, gitignore))

}
