package parser2

import (
	"github.com/revel/cmd/logger"
	"github.com/revel/cmd/model"
	"github.com/revel/cmd/utils"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"golang.org/x/tools/go/packages"
	"os"
	"path/filepath"
	"strings"
)

type (
	SourceProcessor struct {
		revelContainer      *model.RevelContainer
		log                 logger.MultiLogger
		packageList         []*packages.Package
		importMap           map[string]string
		packageMap          map[string]string
		sourceInfoProcessor *SourceInfoProcessor
		sourceInfo          *model.SourceInfo
	}
)

func ProcessSource(revelContainer *model.RevelContainer) (sourceInfo *model.SourceInfo, compileError error) {
	utils.Logger.Info("ProcessSource")
	processor := NewSourceProcessor(revelContainer)
	compileError = processor.parse()
	sourceInfo = processor.sourceInfo
	if compileError == nil {
		processor.log.Infof("From parsers : Structures:%d InitImports:%d ValidationKeys:%d %v", len(sourceInfo.StructSpecs), len(sourceInfo.InitImportPaths), len(sourceInfo.ValidationKeys), sourceInfo.PackageMap)
	}

	return
}

func NewSourceProcessor(revelContainer *model.RevelContainer) *SourceProcessor {
	s := &SourceProcessor{revelContainer:revelContainer, log:utils.Logger.New("parser", "SourceProcessor")}
	s.sourceInfoProcessor = NewSourceInfoProcessor(s)
	return s
}

func (s *SourceProcessor) parse() (compileError error) {
	print("Parsing packages, (may require download if not cached)...")
	if compileError = s.addPackages(); compileError != nil {
		return
	}
	println(" Completed")
	if compileError = s.addImportMap(); compileError != nil {
		return
	}
	if compileError = s.addSourceInfo(); compileError != nil {
		return
	}

	s.sourceInfo.PackageMap = map[string]string{}
	getImportFromMap := func(packagePath string) string {
		for path := range s.packageMap {
			if strings.Index(path, packagePath) == 0 {
				fullPath := s.packageMap[path]
				return fullPath[:(len(fullPath) - len(path) + len(packagePath))]
			}
		}
		return ""
	}
	s.sourceInfo.PackageMap[model.RevelImportPath] = getImportFromMap(model.RevelImportPath)
	s.sourceInfo.PackageMap[s.revelContainer.ImportPath] = getImportFromMap(s.revelContainer.ImportPath)
	for _, module := range s.revelContainer.ModulePathMap {
		s.sourceInfo.PackageMap[module.ImportPath] = getImportFromMap(module.ImportPath)
	}

	return
}

// 这两个方法来自util

// Shortcut to fsWalk
func (s *SourceProcessor) Walk(root string, walkFn filepath.WalkFunc) error {
	return s.fsWalk(root, root, walkFn)
}

// Walk the path tree using the function
// Every file found will call the function
func (s *SourceProcessor) fsWalk(fname string, linkName string, walkFn filepath.WalkFunc) error {
	fsWalkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var name string
		name, err = filepath.Rel(fname, path)
		if err != nil {
			return err
		}

		path = filepath.Join(linkName, name)

		// 改了这里
		if strings.Contains(path, "/leanote/public") ||
			strings.Contains(path, "/leanote/files") ||
			strings.Contains(path, "/leanote/doc") ||
			strings.Contains(path, "/leanote/logs") ||
			strings.Contains(path, "/leanote/build") ||
			strings.Contains(path, "/leanote/target") {
			s.log.Warn("public 或 files 不要处理", "path", path)
			return filepath.SkipDir
		}

		if err == nil && info.Mode() & os.ModeSymlink == os.ModeSymlink {
			var symlinkPath string
			symlinkPath, err = filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}

			// https://github.com/golang/go/blob/master/src/path/filepath/path.go#L392
			info, err = os.Lstat(symlinkPath)

			if err != nil {
				return walkFn(path, info, err)
			}

			if info.IsDir() {
				return s.fsWalk(symlinkPath, path, walkFn)
			}
		}

		return walkFn(path, info, err)
	}
	err := filepath.Walk(fname, fsWalkFunc)
	return err
}

// Using the packages.Load function load all the packages and type specifications (forces compile).
// this sets the SourceProcessor.packageList         []*packages.Package
func (s *SourceProcessor) addPackages() (err error) {
	allPackages := []string{model.RevelImportPath + "/..."}
	for _, module := range s.revelContainer.ModulePathMap {
		allPackages = append(allPackages, module.ImportPath + "/...") // +"/app/controllers/...")
	}
	s.log.Info("Reading packages", "packageList", allPackages)
	//allPackages = []string{s.revelContainer.ImportPath + "/..."} //+"/app/controllers/..."}

	config := &packages.Config{
		// ode: packages.NeedSyntax | packages.NeedCompiledGoFiles,
		Mode:
		packages.NeedTypes | // For compile error
			packages.NeedDeps | // To load dependent files
			packages.NeedName | // Loads the full package name
			packages.NeedSyntax, // To load ast tree (for end points)
		//Mode:	packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
		//	packages.NeedImports | packages.NeedDeps | packages.NeedExportsFile |
		//	packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo |
		//	packages.NeedTypesSizes,

		//Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedExportsFile | packages.NeedFiles |
		//	packages.NeedCompiledGoFiles | packages.NeedTypesSizes |
		//	packages.NeedSyntax | packages.NeedCompiledGoFiles ,
		//Mode:  packages.NeedSyntax | packages.NeedCompiledGoFiles |  packages.NeedName | packages.NeedFiles |
		//	packages.LoadTypes | packages.NeedTypes | packages.NeedDeps,  //, // |
		// packages.NeedTypes, // packages.LoadTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		//packages.LoadSyntax | packages.NeedDeps,
		Dir:s.revelContainer.AppPath,
	}
	s.packageList, err = packages.Load(config, allPackages...)
	s.log.Info("Loaded modules ", "len results", len(s.packageList), "error", err)

	// Now process the files in the aap source folder	s.revelContainer.ImportPath + "/...",
	err = s.Walk(s.revelContainer.BasePath, s.processPath)
	s.log.Info("Loaded apps and modules ", "len results", len(s.packageList), "error", err)
	return
}

// This callback is used to build the packages for the "app" package. This allows us to
// parse the source files without doing a full compile on them
// This callback only processes folders, so any files passed to this will return a nil
func (s *SourceProcessor) processPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		s.log.Error("Error scanning app source:", "error", err)
		return nil
	}

	// Ignore files and folders not marked tmp (since those are generated)
	if !info.IsDir() || info.Name() == "tmp" {
		return nil
	}

	// Real work for processing the folder
	pkgImportPath := s.revelContainer.ImportPath
	appPath := s.revelContainer.BasePath
	if appPath != path {
		pkgImportPath = s.revelContainer.ImportPath + "/" + filepath.ToSlash(path[len(appPath) + 1:])
	}
	s.log.Info("Processing source package folder", "package", pkgImportPath, "path", path)

	// Parse files within the path.
	var pkgMap map[string]*ast.Package
	fset := token.NewFileSet()
	pkgMap, err = parser.ParseDir(
		fset,
		path,
		func(f os.FileInfo) bool {
			return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && strings.HasSuffix(f.Name(), ".go")
		},
		0)

	if err != nil {
		if errList, ok := err.(scanner.ErrorList); ok {
			var pos = errList[0].Pos
			newError := &utils.SourceError{
				SourceType:  ".go source",
				Title:       "Go Compilation Error",
				Path:        pos.Filename,
				Description: errList[0].Msg,
				Line:        pos.Line,
				Column:      pos.Column,
				SourceLines: utils.MustReadLines(pos.Filename),
			}

			errorLink := s.revelContainer.Config.StringDefault("error.link", "")
			if errorLink != "" {
				newError.SetLink(errorLink)
			}
			return newError
		}

		// This is exception, err already checked above. Here just a print
		ast.Print(nil, err)
		s.log.Fatal("Failed to parse dir", "error", err)
	}

	// Skip "main" packages.
	delete(pkgMap, "main")

	// Ignore packages that end with _test
	// These cannot be included in source code that is not generated specifically as a test
	for i := range pkgMap {
		if len(i) > 6 {
			if string(i[len(i) - 5:]) == "_test" {
				delete(pkgMap, i)
			}
		}
	}

	// If there is no code in this directory, skip it.
	if len(pkgMap) == 0 {
		return nil
	}

	// There should be only one package in this directory.
	if len(pkgMap) > 1 {
		for i := range pkgMap {
			println("Found duplicate packages in single directory ", i)
		}
		utils.Logger.Fatal("Most unexpected! Multiple packages in a single directory:", "packages", pkgMap)
	}

	// At this point there is only one package in the pkgs map,
	p := &packages.Package{}
	p.PkgPath = pkgImportPath
	p.Fset = fset
	for _, pkg := range pkgMap {
		p.Name = pkg.Name
		s.log.Info("Found package", "pkg.Name", pkg.Name, "p.Name", p.PkgPath)
		for filename, astFile := range pkg.Files {
			p.Syntax = append(p.Syntax, astFile)
			p.GoFiles = append(p.GoFiles, filename)
		}
	}
	s.packageList = append(s.packageList, p)

	return nil
}

// This function is used to populate a map so that we can lookup controller embedded types in order to determine
// if a Struct inherits from from revel.Controller
func (s *SourceProcessor) addImportMap() (err error) {
	s.importMap = map[string]string{}
	s.packageMap = map[string]string{}
	for _, p := range s.packageList {

		if len(p.Errors) > 0 {
			// Generate a compile error
			for _, e := range p.Errors {
				s.log.Info("While reading packages encountered import error ignoring ", "PkgPath", p.PkgPath, "error", e)
			}
		}
		for _, tree := range p.Syntax {
			s.importMap[tree.Name.Name] = p.PkgPath
		}
	}
	return
}

func (s *SourceProcessor) addSourceInfo() (err error) {
	for _, p := range s.packageList {
		if sourceInfo := s.sourceInfoProcessor.processPackage(p); sourceInfo != nil {
			if s.sourceInfo != nil {
				s.sourceInfo.Merge(sourceInfo)
			} else {
				s.sourceInfo = sourceInfo
			}
		}
	}
	return
}
