// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"runtime"

	"github.com/revel/revel"
)

var cmdVersion = &Command{
	UsageLine: "version",
	Short:     "displays the Revel Framework and Go version",
	Long: `
Displays the Revel Framework and Go version.

For example:

    revel version
`,
}

func init() {
	cmdVersion.Run = versionApp
}

func versionApp(args []string) {
	fmt.Printf("Version(s):")
	fmt.Printf("\n   Revel v%v (%v)", revel.Version, revel.BuildDate)
	fmt.Printf("\n   %s %s/%s\n\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
