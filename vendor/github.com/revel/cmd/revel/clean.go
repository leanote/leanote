// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
)

var cmdClean = &Command{
	UsageLine: "clean [import path]",
	Short:     "clean a Revel application's temp files",
	Long: `
Clean the Revel web application named by the given import path.

For example:

    revel clean github.com/revel/examples/chat

It removes the app/tmp and app/routes directory.
`,
}

func init() {
	cmdClean.Run = cleanApp
}

func cleanApp(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, cmdClean.Long)
		return
	}

	appPkg, err := build.Import(args[0], "", build.FindOnly)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Abort: Failed to find import path:", err)
		return
	}

	purgeDirs := []string{
		filepath.Join(appPkg.Dir, "app", "tmp"),
		filepath.Join(appPkg.Dir, "app", "routes"),
	}

	for _, dir := range purgeDirs {
		fmt.Println("Removing:", dir)
		err = os.RemoveAll(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Abort:", err)
			return
		}
	}
}
