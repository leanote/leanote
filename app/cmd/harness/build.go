// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package harness

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/leanote/leanote/app/cmd/parser2"
	"github.com/revel/cmd/model"
	"github.com/revel/cmd/parser"
	"github.com/revel/cmd/utils"
)

var importErrorPattern = regexp.MustCompile("cannot find package \"([^\"]+)\"")

type ByString []*model.TypeInfo

func (c ByString) Len() int {
	return len(c)
}
func (c ByString) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c ByString) Less(i, j int) bool {
	return c[i].String() < c[j].String()
}

// Build the app:
// 1. Generate the the main.go file.
// 2. Run the appropriate "go build" command.
// Requires that revel.Init has been called previously.
// Returns the path to the built binary, and an error if there was a problem building it.
func Build(c *model.CommandConfig, paths *model.RevelContainer) (_ *App, err error) {
	// First, clear the generated files (to avoid them messing with ProcessSource).
	cleanSource(paths, "tmp", "routes")

	var sourceInfo *model.SourceInfo

	if c.HistoricBuildMode {
		sourceInfo, err = parser.ProcessSource(paths)
	} else {
		sourceInfo, err = parser2.ProcessSource(paths)
	}
	if err != nil {
		return
	}

	// Add the db.import to the import paths.
	if dbImportPath, found := paths.Config.String("db.import"); found {
		sourceInfo.InitImportPaths = append(sourceInfo.InitImportPaths, strings.Split(dbImportPath, ",")...)
	}

	// Sort controllers so that file generation is reproducible
	controllers := sourceInfo.ControllerSpecs()
	sort.Stable(ByString(controllers))

	// Generate two source files.
	templateArgs := map[string]interface{}{
		"ImportPath":     paths.ImportPath,
		"Controllers":    controllers,
		"ValidationKeys": sourceInfo.ValidationKeys,
		"ImportPaths":    calcImportAliases(sourceInfo),
		"TestSuites":     sourceInfo.TestSuites(),
	}

	// Generate code for the main, run and routes file.
	// The run file allows external programs to launch and run the application
	// without being the main thread
	cleanSource(paths, "tmp", "routes")

	if err = genSource(paths, "tmp", "main.go", RevelMainTemplate, templateArgs); err != nil {
		return
	}
	if err = genSource(paths, filepath.Join("tmp", "run"), "run.go", RevelRunTemplate, templateArgs); err != nil {
		return
	}
	if err = genSource(paths, "routes", "routes.go", RevelRoutesTemplate, templateArgs); err != nil {
		return
	}

	utils.Logger.Warn("gen tmp/main.go, tmp/run/run.go, routes/routes.go success!!")

	return // 改了这里

	// Read build config.
	buildTags := paths.Config.StringDefault("build.tags", "")

	// Build the user program (all code under app).
	// It relies on the user having "go" installed.
	goPath, err := exec.LookPath("go")
	if err != nil {
		utils.Logger.Fatal("Go executable not found in PATH.")
	}

	// Binary path is a combination of target/app directory, app's import path and its name.
	binName := filepath.Join("target", "app", paths.ImportPath, filepath.Base(paths.BasePath))

	// Change binary path for Windows build
	goos := runtime.GOOS
	if goosEnv := os.Getenv("GOOS"); goosEnv != "" {
		goos = goosEnv
	}
	if goos == "windows" {
		binName += ".exe"
	}

	gotten := make(map[string]struct{})
	contains := func(s []string, e string) bool {
		for _, a := range s {
			if a == e {
				return true
			}
		}
		return false
	}

	if len(c.GoModFlags) > 0 {
		for _, gomod := range c.GoModFlags {
			goModCmd := exec.Command(goPath, append([]string{"mod"}, strings.Split(gomod, " ")...)...)
			utils.CmdInit(goModCmd, !c.Vendored, c.AppPath)
			output, err := goModCmd.CombinedOutput()
			utils.Logger.Info("Gomod applied ", "output", string(output))

			// If the build succeeded, we're done.
			if err != nil {
				utils.Logger.Error("Gomod Failed continuing ", "error", err, "output", string(output))
			}
		}
	}

	for {
		appVersion := getAppVersion(paths)
		if appVersion == "" {
			appVersion = "noVersionProvided"
		}

		buildTime := time.Now().UTC().Format(time.RFC3339)
		versionLinkerFlags := fmt.Sprintf("-X '%s/app.AppVersion=%s' -X '%s/app.BuildTime=%s'",
			paths.ImportPath, appVersion, paths.ImportPath, buildTime)

		// Append any build flags specified, they will override existing flags
		flags := []string{}
		if len(c.BuildFlags) == 0 {
			flags = []string{
				"build",
				"-ldflags", versionLinkerFlags,
				"-tags", buildTags,
				"-o", binName}
		} else {
			if !contains(c.BuildFlags, "build") {
				flags = []string{"build"}
			}
			if !contains(flags, "-ldflags") {
				ldflags := "-ldflags= " + versionLinkerFlags
				// Add user defined build flags
				for i := range c.BuildFlags {
					ldflags += " -X '" + c.BuildFlags[i] + "'"
				}
				flags = append(flags, ldflags)
			}
			if !contains(flags, "-tags") && buildTags != "" {
				flags = append(flags, "-tags", buildTags)
			}
			if !contains(flags, "-o") {
				flags = append(flags, "-o", binName)
			}
		}

		// Note: It's not applicable for filepath.* usage
		flags = append(flags, path.Join(paths.ImportPath, "app", "tmp"))

		buildCmd := exec.Command(goPath, flags...)
		if !c.Vendored {
			// This is Go main path
			gopath := c.GoPath
			for _, o := range paths.ModulePathMap {
				gopath += string(filepath.ListSeparator) + o.Path
			}

			buildCmd.Env = append(os.Environ(),
				"GOPATH=" + gopath,
			)
		}
		utils.CmdInit(buildCmd, !c.Vendored, c.AppPath)

		utils.Logger.Info("Exec:", "args", buildCmd.Args, "working dir", buildCmd.Dir)
		output, err := buildCmd.CombinedOutput()

		// If the build succeeded, we're done.
		if err == nil {
			utils.Logger.Info("Build successful continuing")
			return NewApp(binName, paths, sourceInfo.PackageMap), nil
		}

		// Since there was an error, capture the output in case we need to report it
		stOutput := string(output)
		utils.Logger.Infof("Got error on build of app %s", stOutput)

		// See if it was an import error that we can go get.
		matches := importErrorPattern.FindAllStringSubmatch(stOutput, -1)
		utils.Logger.Info("Build failed checking for missing imports", "message", stOutput, "missing_imports", len(matches))
		if matches == nil {
			utils.Logger.Info("Build failed no missing imports", "message", stOutput)
			return nil, newCompileError(paths, output)
		}
		utils.Logger.Warn("Detected missing packages, importing them", "packages", len(matches))
		for _, match := range matches {
			// Ensure we haven't already tried to go get it.
			pkgName := match[1]
			utils.Logger.Info("Trying to import ", "package", pkgName)
			if _, alreadyTried := gotten[pkgName]; alreadyTried {
				utils.Logger.Error("Failed to import ", "package", pkgName)
				return nil, newCompileError(paths, output)
			}
			gotten[pkgName] = struct{}{}
			if err := c.PackageResolver(pkgName); err != nil {
				utils.Logger.Error("Unable to resolve package", "package", pkgName, "error", err)
				return nil, newCompileError(paths, []byte(err.Error()))
			}
		}

		// Success getting the import, attempt to build again.
	}

	// TODO remove this unreachable code and document it
	utils.Logger.Fatal("Not reachable")
	return nil, nil
}

// Try to define a version string for the compiled app
// The following is tried (first match returns):
// - Read a version explicitly specified in the APP_VERSION environment
//   variable
// - Read the output of "git describe" if the source is in a git repository
// If no version can be determined, an empty string is returned.
func getAppVersion(paths *model.RevelContainer) string {
	if version := os.Getenv("APP_VERSION"); version != "" {
		return version
	}

	// Check for the git binary
	if gitPath, err := exec.LookPath("git"); err == nil {
		// Check for the .git directory
		gitDir := filepath.Join(paths.BasePath, ".git")
		info, err := os.Stat(gitDir)
		if (err != nil && os.IsNotExist(err)) || !info.IsDir() {
			return ""
		}
		gitCmd := exec.Command(gitPath, "--git-dir=" + gitDir, "--work-tree=" + paths.BasePath, "describe", "--always", "--dirty")
		utils.Logger.Info("Exec:", "args", gitCmd.Args)
		output, err := gitCmd.Output()

		if err != nil {
			utils.Logger.Error("Cannot determine git repository version:", "error", err)
			return ""
		}

		return "git-" + strings.TrimSpace(string(output))
	}

	return ""
}

func cleanSource(paths *model.RevelContainer, dirs ...string) {
	for _, dir := range dirs {
		cleanDir(paths, dir)
	}
}

func cleanDir(paths *model.RevelContainer, dir string) {
	utils.Logger.Info("Cleaning dir ", "dir", dir)
	tmpPath := filepath.Join(paths.AppPath, dir)
	f, err := os.Open(tmpPath)
	if err != nil {
		if !os.IsNotExist(err) {
			utils.Logger.Error("Failed to clean dir:", "error", err)
		}
	} else {
		defer func() {
			_ = f.Close()
		}()

		infos, err := f.Readdir(0)
		if err != nil {
			if !os.IsNotExist(err) {
				utils.Logger.Fatal("Failed to clean dir:", "error", err)
			}
		} else {
			for _, info := range infos {
				pathName := filepath.Join(tmpPath, info.Name())
				if info.IsDir() {
					err := os.RemoveAll(pathName)
					if err != nil {
						utils.Logger.Fatal("Failed to remove dir:", "error", err)
					}
				} else {
					err := os.Remove(pathName)
					if err != nil {
						utils.Logger.Fatal("Failed to remove file:", "error", err)
					}
				}
			}
		}
	}
}

// genSource renders the given template to produce source code, which it writes
// to the given directory and file.
func genSource(paths *model.RevelContainer, dir, filename, templateSource string, args map[string]interface{}) error {

	return utils.GenerateTemplate(filepath.Join(paths.AppPath, dir, filename), templateSource, args)
}

// Looks through all the method args and returns a set of unique import paths
// that cover all the method arg types.
// Additionally, assign package aliases when necessary to resolve ambiguity.
func calcImportAliases(src *model.SourceInfo) map[string]string {
	aliases := make(map[string]string)
	typeArrays := [][]*model.TypeInfo{src.ControllerSpecs(), src.TestSuites()}
	for _, specs := range typeArrays {
		for _, spec := range specs {
			addAlias(aliases, spec.ImportPath, spec.PackageName)

			for _, methSpec := range spec.MethodSpecs {
				for _, methArg := range methSpec.Args {
					if methArg.ImportPath == "" {
						continue
					}

					addAlias(aliases, methArg.ImportPath, methArg.TypeExpr.PkgName)
				}
			}
		}
	}

	// Add the "InitImportPaths", with alias "_"
	for _, importPath := range src.InitImportPaths {
		if _, ok := aliases[importPath]; !ok {
			aliases[importPath] = "_"
		}
	}

	return aliases
}

// Adds an alias to the map of alias names
func addAlias(aliases map[string]string, importPath, pkgName string) {
	alias, ok := aliases[importPath]
	if ok {
		return
	}
	alias = makePackageAlias(aliases, pkgName)
	aliases[importPath] = alias
}

// Generates a package alias
func makePackageAlias(aliases map[string]string, pkgName string) string {
	i := 0
	alias := pkgName
	for containsValue(aliases, alias) || alias == "revel" {
		alias = fmt.Sprintf("%s%d", pkgName, i)
		i++
	}
	return alias
}

// Returns true if this value is in the map
func containsValue(m map[string]string, val string) bool {
	for _, v := range m {
		if v == val {
			return true
		}
	}
	return false
}

// Parse the output of the "go build" command.
// Return a detailed Error.
func newCompileError(paths *model.RevelContainer, output []byte) *utils.SourceError {
	errorMatch := regexp.MustCompile(`(?m)^([^:#]+):(\d+):(\d+:)? (.*)$`).
		FindSubmatch(output)
	if errorMatch == nil {
		errorMatch = regexp.MustCompile(`(?m)^(.*?):(\d+):\s(.*?)$`).FindSubmatch(output)

		if errorMatch == nil {
			utils.Logger.Error("Failed to parse build errors", "error", string(output))
			return &utils.SourceError{
				SourceType:  "Go code",
				Title:       "Go Compilation Error",
				Description: "See console for build error.",
			}
		}

		errorMatch = append(errorMatch, errorMatch[3])

		utils.Logger.Error("Build errors", "errors", string(output))
	}

	findInPaths := func(relFilename string) string {
		// Extract the paths from the gopaths, and search for file there first
		gopaths := filepath.SplitList(build.Default.GOPATH)
		for _, gp := range gopaths {
			newPath := filepath.Join(gp, "src", paths.ImportPath, relFilename)
			println(newPath)
			if utils.Exists(newPath) {
				return newPath
			}
		}
		newPath, _ := filepath.Abs(relFilename)
		utils.Logger.Warn("Could not find in GO path", "file", relFilename)
		return newPath
	}


	// Read the source for the offending file.
	var (
		relFilename = string(errorMatch[1]) // e.g. "src/revel/sample/app/controllers/app.go"
		absFilename = findInPaths(relFilename)
		line, _ = strconv.Atoi(string(errorMatch[2]))
		description = string(errorMatch[4])
		compileError = &utils.SourceError{
			SourceType:  "Go code",
			Title:       "Go Compilation Error",
			Path:        relFilename,
			Description: description,
			Line:        line,
		}
	)

	errorLink := paths.Config.StringDefault("error.link", "")

	if errorLink != "" {
		compileError.SetLink(errorLink)
	}

	fileStr, err := utils.ReadLines(absFilename)
	if err != nil {
		compileError.MetaError = absFilename + ": " + err.Error()
		utils.Logger.Info("Unable to readlines " + compileError.MetaError, "error", err)
		return compileError
	}

	compileError.SourceLines = fileStr
	return compileError
}

// RevelMainTemplate template for app/tmp/run/run.go
const RevelRunTemplate = `// GENERATED CODE - DO NOT EDIT
// This file is the run file for Revel.
// It registers all the controllers and provides details for the Revel server engine to
// properly inject parameters directly into the action endpoints.
package run

import (
	"reflect"
	"github.com/revel/revel"{{range $k, $v := $.ImportPaths}}
	{{$v}} "{{$k}}"{{end}}
	"github.com/revel/revel/testing"
)

var (
	// So compiler won't complain if the generated code doesn't reference reflect package...
	_ = reflect.Invalid
)

// Register and run the application
func Run(port int) {
	Register()
	revel.Run(port)
}

// Register all the controllers
func Register() {
	revel.AppLog.Info("Running revel server")
	{{range $i, $c := .Controllers}}
	revel.RegisterController((*{{index $.ImportPaths .ImportPath}}.{{.StructName}})(nil),
		[]*revel.MethodType{
			{{range .MethodSpecs}}&revel.MethodType{
				Name: "{{.Name}}",
				Args: []*revel.MethodArg{ {{range .Args}}
					&revel.MethodArg{Name: "{{.Name}}", Type: reflect.TypeOf((*{{index $.ImportPaths .ImportPath | .TypeExpr.TypeName}})(nil)) },{{end}}
				},
				RenderArgNames: map[int][]string{ {{range .RenderCalls}}
					{{.Line}}: []string{ {{range .Names}}
						"{{.}}",{{end}}
					},{{end}}
				},
			},
			{{end}}
		})
	{{end}}
	revel.DefaultValidationKeys = map[string]map[int]string{ {{range $path, $lines := .ValidationKeys}}
		"{{$path}}": { {{range $line, $key := $lines}}
			{{$line}}: "{{$key}}",{{end}}
		},{{end}}
	}
	testing.TestSuites = []interface{}{ {{range .TestSuites}}
		(*{{index $.ImportPaths .ImportPath}}.{{.StructName}})(nil),{{end}}
	}
}
`
const RevelMainTemplate = `// GENERATED CODE - DO NOT EDIT
// This file is the main file for Revel.
// It registers all the controllers and provides details for the Revel server engine to
// properly inject parameters directly into the action endpoints.
package main

import (
	"flag"
	"{{.ImportPath}}/app/tmp/run"
	"github.com/revel/revel"
)

var (
	runMode    *string = flag.String("runMode", "", "Run mode.")
	port       *int    = flag.Int("port", 0, "By default, read from app.conf")
	importPath *string = flag.String("importPath", "", "Go Import Path for the app.")
	srcPath    *string = flag.String("srcPath", "", "Path to the source root.")

)

func main() {
	flag.Parse()
	revel.Init(*runMode, *importPath, *srcPath)
	run.Run(*port)
}
`

// RevelRoutesTemplate template for app/conf/routes
const RevelRoutesTemplate = `// GENERATED CODE - DO NOT EDIT
// This file provides a way of creating URL's based on all the actions
// found in all the controllers.
package routes

import "github.com/revel/revel"

{{range $i, $c := .Controllers}}
type t{{.StructName}} struct {}
var {{.StructName}} t{{.StructName}}

{{range .MethodSpecs}}
func (_ t{{$c.StructName}}) {{.Name}}({{range .Args}}
		{{.Name}} {{if .ImportPath}}interface{}{{else}}{{.TypeExpr.TypeName ""}}{{end}},{{end}}
		) string {
	args := make(map[string]string)
	{{range .Args}}
	revel.Unbind(args, "{{.Name}}", {{.Name}}){{end}}
	return revel.MainRouter.Reverse("{{$c.StructName}}.{{.Name}}", args).URL
}
{{end}}
{{end}}
`
