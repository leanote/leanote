// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"strings"

	"fmt"
	"github.com/leanote/leanote/app/cmd/harness" // 只改了这个
	"github.com/revel/cmd/model"
	"github.com/revel/cmd/utils"
)

var cmdBuild = &Command{
	UsageLine: "revel build [-r [run mode]] [import path] [target path] ",
	Short:     "build a Revel application (e.g. for deployment)",
	Long: `
Build the Revel web application named by the given import path.
This allows it to be deployed and run on a machine that lacks a Go installation.

For example:

    revel build github.com/revel/examples/chat /tmp/chat

`,
}

func init() {
	cmdBuild.RunWith = buildApp
	cmdBuild.UpdateConfig = updateBuildConfig
}

// The update config updates the configuration command so that it can run
func updateBuildConfig(c *model.CommandConfig, args []string) bool {
	c.Index = model.BUILD
	if c.Build.TargetPath == "" {
		c.Build.TargetPath = "target"
	}
	if len(args) == 0 && c.Build.ImportPath != "" {
		return true
	}
	// If arguments were passed in then there must be two
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "%s\n%s", cmdBuild.UsageLine, cmdBuild.Long)
		return false
	}

	c.Build.ImportPath = args[0]
	c.Build.TargetPath = args[1]
	if len(args) > 2 {
		c.Build.Mode = args[2]
	}
	return true
}

// The main entry point to build application from command line
func buildApp(c *model.CommandConfig) (err error) {

	appImportPath, destPath, mode := c.ImportPath, c.Build.TargetPath, DefaultRunMode
	if len(c.Build.Mode) > 0 {
		mode = c.Build.Mode
	}

	// Convert target to absolute path
	c.Build.TargetPath, _ = filepath.Abs(destPath)
	c.Build.Mode = mode
	c.Build.ImportPath = appImportPath

	revel_paths, err := model.NewRevelPaths(mode, appImportPath, c.AppPath, model.NewWrappedRevelCallback(nil, c.PackageResolver))
	if err != nil {
		return
	}

	if err = buildSafetyCheck(destPath); err != nil {
		return
	}

	// Ensure the application can be built, this generates the main file
	app, err := harness.Build(c, revel_paths)
	if err != nil {
		return err
	}
	// Copy files
	// Included are:
	// - run scripts
	// - binary
	// - revel
	// - app

	return // 改了这里

	packageFolders, err := buildCopyFiles(c, app, revel_paths)
	if err != nil {
		return
	}
	err = buildCopyModules(c, revel_paths, packageFolders, app)
	if err != nil {
		return
	}
	err = buildWriteScripts(c, app)
	if err != nil {
		return
	}
	return
}

// Copy the files to the target
func buildCopyFiles(c *model.CommandConfig, app *harness.App, revel_paths *model.RevelContainer) (packageFolders []string, err error) {
	appImportPath, destPath := c.ImportPath, c.Build.TargetPath

	// Revel and the app are in a directory structure mirroring import path
	srcPath := filepath.Join(destPath, "src")
	destBinaryPath := filepath.Join(destPath, filepath.Base(app.BinaryPath))
	tmpRevelPath := filepath.Join(srcPath, filepath.FromSlash(model.RevelImportPath))
	if err = utils.CopyFile(destBinaryPath, filepath.Join(revel_paths.BasePath, app.BinaryPath)); err != nil {
		return
	}
	utils.MustChmod(destBinaryPath, 0755)

	// Copy the templates  from the revel
	if err = utils.CopyDir(filepath.Join(tmpRevelPath, "conf"), filepath.Join(revel_paths.RevelPath, "conf"), nil); err != nil {
		return
	}
	if err = utils.CopyDir(filepath.Join(tmpRevelPath, "templates"), filepath.Join(revel_paths.RevelPath, "templates"), nil); err != nil {
		return
	}

	// Get the folders to be packaged
	packageFolders = strings.Split(revel_paths.Config.StringDefault("package.folders", "conf,public,app/views"), ",")
	for i, p := range packageFolders {
		// Clean spaces, reformat slash to filesystem
		packageFolders[i] = filepath.FromSlash(strings.TrimSpace(p))
	}

	if c.Build.CopySource {
		err = utils.CopyDir(filepath.Join(srcPath, filepath.FromSlash(appImportPath)), revel_paths.BasePath, nil)
		if err != nil {
			return
		}
	} else {
		for _, folder := range packageFolders {
			err = utils.CopyDir(
				filepath.Join(srcPath, filepath.FromSlash(appImportPath), folder),
				filepath.Join(revel_paths.BasePath, folder),
				nil)
			if err != nil {
				return
			}
		}
	}

	return
}

// Based on the section copy over the build modules
func buildCopyModules(c *model.CommandConfig, revel_paths *model.RevelContainer, packageFolders []string, app *harness.App) (err error) {
	destPath := filepath.Join(c.Build.TargetPath, "src")
	// Find all the modules used and copy them over.
	config := revel_paths.Config.Raw()

	// We should only copy over the section of options what the build is targeted for
	// We will default to prod
	moduleImportList := []string{}
	for _, section := range config.Sections() {
		// If the runmode is defined we will only import modules defined for that run mode
		if c.Build.Mode != "" && c.Build.Mode != section {
			continue
		}
		options, _ := config.SectionOptions(section)
		for _, key := range options {
			if !strings.HasPrefix(key, "module.") {
				continue
			}
			moduleImportPath, _ := config.String(section, key)
			if moduleImportPath == "" {
				continue
			}
			moduleImportList = append(moduleImportList, moduleImportPath)

		}
	}

	// Copy the the paths for each of the modules
	for _, importPath := range moduleImportList {
		fsPath := app.PackagePathMap[importPath]
		utils.Logger.Info("Copy files ", "to", filepath.Join(destPath, importPath), "from", fsPath)
		if c.Build.CopySource {
			err = utils.CopyDir(filepath.Join(destPath, importPath), fsPath, nil)
			if err != nil {
				return
			}
		} else {
			for _, folder := range packageFolders {
				err = utils.CopyDir(
					filepath.Join(destPath, importPath, folder),
					filepath.Join(fsPath, folder),
					nil)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

// Write the run scripts for the build
func buildWriteScripts(c *model.CommandConfig, app *harness.App) (err error) {
	tmplData := map[string]interface{}{
		"BinName":    filepath.Base(app.BinaryPath),
		"ImportPath": c.Build.ImportPath,
		"Mode":       c.Build.Mode,
	}

	err = utils.GenerateTemplate(
		filepath.Join(c.Build.TargetPath, "run.sh"),
		PACKAGE_RUN_SH,
		tmplData,
	)
	if err != nil {
		return
	}
	utils.MustChmod(filepath.Join(c.Build.TargetPath, "run.sh"), 0755)
	err = utils.GenerateTemplate(
		filepath.Join(c.Build.TargetPath, "run.bat"),
		PACKAGE_RUN_BAT,
		tmplData,
	)
	if err != nil {
		return
	}

	fmt.Println("Your application has been built in:", c.Build.TargetPath)

	return
}

// Checks to see if the target folder exists and can be created
func buildSafetyCheck(destPath string) error {

	// First, verify that it is either already empty or looks like a previous
	// build (to avoid clobbering anything)
	if utils.Exists(destPath) && !utils.Empty(destPath) && !utils.Exists(filepath.Join(destPath, "run.sh")) {
		return utils.NewBuildError("Abort: %s exists and does not look like a build directory.", "path", destPath)
	}

	if err := os.RemoveAll(destPath); err != nil && !os.IsNotExist(err) {
		return utils.NewBuildIfError(err, "Remove all error", "path", destPath)
	}

	if err := os.MkdirAll(destPath, 0777); err != nil {
		return utils.NewBuildIfError(err, "MkDir all error", "path", destPath)
	}
	return nil
}

const PACKAGE_RUN_SH = `#!/bin/sh

SCRIPTPATH=$(cd "$(dirname "$0")"; pwd)
"$SCRIPTPATH/{{.BinName}}" -importPath {{.ImportPath}} -srcPath "$SCRIPTPATH/src" -runMode {{.Mode}}
`
const PACKAGE_RUN_BAT = `@echo off

{{.BinName}} -importPath {{.ImportPath}} -srcPath "%CD%\src" -runMode {{.Mode}}
`
