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
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/revel/revel"
)

var importErrorPattern = regexp.MustCompile("cannot find package \"([^\"]+)\"")

// Build the app:
// 1. Generate the the main.go file.
// 2. Run the appropriate "go build" command.
// Requires that revel.Init has been called previously.
// Returns the path to the built binary, and an error if there was a problem building it.
func Build(buildFlags ...string) (app *App, compileError *revel.Error) {
	// First, clear the generated files (to avoid them messing with ProcessSource).
	cleanSource("tmp", "routes")

	sourceInfo, compileError := ProcessSource(revel.CodePaths)
	if compileError != nil {
		return nil, compileError
	}

	// Add the db.import to the import paths.
	if dbImportPath, found := revel.Config.String("db.import"); found {
		sourceInfo.InitImportPaths = append(sourceInfo.InitImportPaths, strings.Split(dbImportPath,",")...)
	}

	// Generate two source files.
	templateArgs := map[string]interface{}{
		"Controllers":    sourceInfo.ControllerSpecs(),
		"ValidationKeys": sourceInfo.ValidationKeys,
		"ImportPaths":    calcImportAliases(sourceInfo),
		"TestSuites":     sourceInfo.TestSuites(),
	}
	genSource("tmp", "main.go", RevelMainTemplate, templateArgs)
	genSource("routes", "routes.go", RevelRoutesTemplate, templateArgs)

	// Read build config.
	buildTags := revel.Config.StringDefault("build.tags", "")

	// Build the user program (all code under app).
	// It relies on the user having "go" installed.
	goPath, err := exec.LookPath("go")
	if err != nil {
		revel.RevelLog.Fatalf("Go executable not found in PATH.")
	}

	// Detect if deps tool should be used (is there a vendor folder ?)
	useVendor := revel.DirExists(filepath.Join(revel.BasePath, "vendor"))
	basePath := revel.BasePath
	for !useVendor {
		basePath = filepath.Dir(basePath)
		found := false
		// Check to see if we are still in the GOPATH
		for _, path := range filepath.SplitList(build.Default.GOPATH) {
			if strings.HasPrefix(basePath, path) {
				found = true
				break
			}
		}
		if !found {
			break
		} else {
			useVendor = revel.DirExists(filepath.Join(basePath, "vendor"))
		}
	}

	var depPath string
	if useVendor {
		revel.RevelLog.Info("Vendor folder detected, scanning for deps in path")
		depPath, err = exec.LookPath("dep")
		if err != nil {
			// Do not halt build unless a new package needs to be imported
			revel.RevelLog.Warn("Build: `dep` executable not found in PATH, but vendor folder detected." +
				"Packages can only be added automatically to the vendor folder using the `dep` tool. " +
				"You can install the `dep` tool by doing a `go get -u github.com/golang/dep/cmd/dep`")
		}
	} else {
		revel.RevelLog.Info("No vendor folder detected, not using dependency manager to import files")
	}

	pkg, err := build.Default.Import(revel.ImportPath, "", build.FindOnly)
	if err != nil {
		revel.RevelLog.Fatal("Failure importing", "path", revel.ImportPath)
	}

	// Binary path is a combination of $GOBIN/revel.d directory, app's import path and its name.
	binName := filepath.Join(pkg.BinDir, "revel.d", revel.ImportPath, filepath.Base(revel.BasePath))

	// Change binary path for Windows build
	goos := runtime.GOOS
	if goosEnv := os.Getenv("GOOS"); goosEnv != "" {
		goos = goosEnv
	}
	if goos == "windows" {
		binName += ".exe"
	}

	gotten := make(map[string]struct{})
	for {
		appVersion := getAppVersion()

		buildTime := time.Now().UTC().Format(time.RFC3339)
		versionLinkerFlags := fmt.Sprintf("-X %s/app.AppVersion=%s -X %s/app.BuildTime=%s",
			revel.ImportPath, appVersion, revel.ImportPath, buildTime)

		flags := []string{
			"build",
			"-i",
			"-ldflags", versionLinkerFlags,
			"-tags", buildTags,
			"-o", binName}

		// Add in build flags
		flags = append(flags, buildFlags...)

		// This is Go main path
		// Note: It's not applicable for filepath.* usage
		flags = append(flags, path.Join(revel.ImportPath, "app", "tmp"))

		buildCmd := exec.Command(goPath, flags...)
		revel.RevelLog.Debug("Exec:", "args", buildCmd.Args)
		output, err := buildCmd.CombinedOutput()

		// If the build succeeded, we're done.
		if err == nil {
			return NewApp(binName), nil
		}
		revel.RevelLog.Error(string(output))

		// See if it was an import error that we can go get.
		matches := importErrorPattern.FindAllStringSubmatch(string(output), -1)
		if matches == nil {
			return nil, newCompileError(output)
		}
		for _, match := range matches {
			// Ensure we haven't already tried to go get it.
			pkgName := match[1]
			if _, alreadyTried := gotten[pkgName]; alreadyTried {
				return nil, newCompileError(output)
			}
			gotten[pkgName] = struct{}{}

			// Execute "go get <pkg>"
			// Or dep `dep ensure -add <pkg>` if it is there
			var getCmd *exec.Cmd
			if useVendor {
				if depPath == "" {
					revel.RevelLog.Error("Build: Vendor folder found, but the `dep` tool was not found, " +
						"if you use a different vendoring (package management) tool please add the following packages by hand, " +
						"or install the `dep` tool into your gopath by doing a `go get -u github.com/golang/dep/cmd/dep`. " +
						"For more information and usage of the tool please see http://github.com/golang/dep")
					for _, pkg := range matches {
						revel.RevelLog.Error("Missing package", "package", pkg[1])
					}
				}
				getCmd = exec.Command(depPath, "ensure", "-add", pkgName)
			} else {
				getCmd = exec.Command(goPath, "get", pkgName)
			}
			revel.RevelLog.Debug("Exec:", "args", getCmd.Args)
			getOutput, err := getCmd.CombinedOutput()
			if err != nil {
				revel.RevelLog.Error(string(getOutput))
				return nil, newCompileError(output)
			}
		}

		// Success getting the import, attempt to build again.
	}

	// TODO remove this unreachable code and document it
	revel.RevelLog.Fatalf("Not reachable")
	return nil, nil
}

// Try to define a version string for the compiled app
// The following is tried (first match returns):
// - Read a version explicitly specified in the APP_VERSION environment
//   variable
// - Read the output of "git describe" if the source is in a git repository
// If no version can be determined, an empty string is returned.
func getAppVersion() string {
	if version := os.Getenv("APP_VERSION"); version != "" {
		return version
	}

	// Check for the git binary
	if gitPath, err := exec.LookPath("git"); err == nil {
		// Check for the .git directory
		gitDir := filepath.Join(revel.BasePath, ".git")
		info, err := os.Stat(gitDir)
		if (err != nil && os.IsNotExist(err)) || !info.IsDir() {
			return ""
		}
		gitCmd := exec.Command(gitPath, "--git-dir="+gitDir, "describe", "--always", "--dirty")
		revel.RevelLog.Debug("Exec:", "args", gitCmd.Args)
		output, err := gitCmd.Output()

		if err != nil {
			revel.RevelLog.Warn("Cannot determine git repository version:", "error", err)
			return ""
		}

		return "git-" + strings.TrimSpace(string(output))
	}

	return ""
}

func cleanSource(dirs ...string) {
	for _, dir := range dirs {
		cleanDir(dir)
	}
}

func cleanDir(dir string) {
	revel.RevelLog.Info("Cleaning dir " + dir)
	tmpPath := filepath.Join(revel.AppPath, dir)
	f, err := os.Open(tmpPath)
	if err != nil {
		if !os.IsNotExist(err) {
			revel.RevelLog.Error("Failed to clean dir:", "error", err)
		}
	} else {
		defer func() {
			_ = f.Close()
		}()

		infos, err := f.Readdir(0)
		if err != nil {
			if !os.IsNotExist(err) {
				revel.RevelLog.Error("Failed to clean dir:", "error", err)
			}
		} else {
			for _, info := range infos {
				pathName := filepath.Join(tmpPath, info.Name())
				if info.IsDir() {
					err := os.RemoveAll(pathName)
					if err != nil {
						revel.RevelLog.Error("Failed to remove dir:", "error", err)
					}
				} else {
					err := os.Remove(pathName)
					if err != nil {
						revel.RevelLog.Error("Failed to remove file:", "error", err)
					}
				}
			}
		}
	}
}

// genSource renders the given template to produce source code, which it writes
// to the given directory and file.
func genSource(dir, filename, templateSource string, args map[string]interface{}) {
	sourceCode := revel.ExecuteTemplate(
		template.Must(template.New("").Parse(templateSource)),
		args)

	// Create a fresh dir.
	cleanSource(dir)
	tmpPath := filepath.Join(revel.AppPath, dir)
	err := os.Mkdir(tmpPath, 0777)
	if err != nil && !os.IsExist(err) {
		revel.RevelLog.Fatalf("Failed to make '%v' directory: %v", dir, err)
	}

	// Create the file
	file, err := os.Create(filepath.Join(tmpPath, filename))
	if err != nil {
		revel.RevelLog.Fatalf("Failed to create file: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err = file.WriteString(sourceCode); err != nil {
		revel.RevelLog.Fatalf("Failed to write to file: %v", err)
	}
}

// Looks through all the method args and returns a set of unique import paths
// that cover all the method arg types.
// Additionally, assign package aliases when necessary to resolve ambiguity.
func calcImportAliases(src *SourceInfo) map[string]string {
	aliases := make(map[string]string)
	typeArrays := [][]*TypeInfo{src.ControllerSpecs(), src.TestSuites()}
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

func addAlias(aliases map[string]string, importPath, pkgName string) {
	alias, ok := aliases[importPath]
	if ok {
		return
	}
	alias = makePackageAlias(aliases, pkgName)
	aliases[importPath] = alias
}

func makePackageAlias(aliases map[string]string, pkgName string) string {
	i := 0
	alias := pkgName
	for containsValue(aliases, alias) || alias == "revel" {
		alias = fmt.Sprintf("%s%d", pkgName, i)
		i++
	}
	return alias
}

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
func newCompileError(output []byte) *revel.Error {
	errorMatch := regexp.MustCompile(`(?m)^([^:#]+):(\d+):(\d+:)? (.*)$`).
		FindSubmatch(output)
	if errorMatch == nil {
		errorMatch = regexp.MustCompile(`(?m)^(.*?)\:(\d+)\:\s(.*?)$`).FindSubmatch(output)

		if errorMatch == nil {
			revel.RevelLog.Error("Failed to parse build errors", "error", string(output))
			return &revel.Error{
				SourceType:  "Go code",
				Title:       "Go Compilation Error",
				Description: "See console for build error.",
			}
		}

		errorMatch = append(errorMatch, errorMatch[3])

		revel.RevelLog.Error("Build errors", "errors", string(output))
	}

	// Read the source for the offending file.
	var (
		relFilename    = string(errorMatch[1]) // e.g. "src/revel/sample/app/controllers/app.go"
		absFilename, _ = filepath.Abs(relFilename)
		line, _        = strconv.Atoi(string(errorMatch[2]))
		description    = string(errorMatch[4])
		compileError   = &revel.Error{
			SourceType:  "Go code",
			Title:       "Go Compilation Error",
			Path:        relFilename,
			Description: description,
			Line:        line,
		}
	)

	errorLink := revel.Config.StringDefault("error.link", "")

	if errorLink != "" {
		compileError.SetLink(errorLink)
	}

	fileStr, err := revel.ReadLines(absFilename)
	if err != nil {
		compileError.MetaError = absFilename + ": " + err.Error()
		revel.RevelLog.Error(compileError.MetaError)
		return compileError
	}

	compileError.SourceLines = fileStr
	return compileError
}

// RevelMainTemplate template for app/tmp/main.go
const RevelMainTemplate = `// GENERATED CODE - DO NOT EDIT
package main

import (
	"flag"
	"reflect"
	"github.com/revel/revel"{{range $k, $v := $.ImportPaths}}
	{{$v}} "{{$k}}"{{end}}
	"github.com/revel/revel/testing"
)

var (
	runMode    *string = flag.String("runMode", "", "Run mode.")
	port       *int    = flag.Int("port", 0, "By default, read from app.conf")
	importPath *string = flag.String("importPath", "", "Go Import Path for the app.")
	srcPath    *string = flag.String("srcPath", "", "Path to the source root.")

	// So compiler won't complain if the generated code doesn't reference reflect package...
	_ = reflect.Invalid
)

func main() {
	flag.Parse()
	revel.Init(*runMode, *importPath, *srcPath)
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

	revel.Run(*port)
}
`

// RevelRoutesTemplate template for app/conf/routes
const RevelRoutesTemplate = `// GENERATED CODE - DO NOT EDIT
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
