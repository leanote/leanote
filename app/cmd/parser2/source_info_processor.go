package parser2

import (
	"github.com/revel/cmd/utils"
	"golang.org/x/tools/go/packages"
	"github.com/revel/cmd/model"
	"go/ast"
	"go/token"
	"strings"
	"path/filepath"
	"github.com/revel/cmd/logger"
)

type (
	SourceInfoProcessor struct {
		sourceProcessor *SourceProcessor
	}
)

func NewSourceInfoProcessor(sourceProcessor *SourceProcessor) *SourceInfoProcessor {
	return &SourceInfoProcessor{sourceProcessor:sourceProcessor}
}

func (s *SourceInfoProcessor) processPackage(p *packages.Package) (sourceInfo *model.SourceInfo) {
	sourceInfo = &model.SourceInfo{
		ValidationKeys: map[string]map[int]string{},
	}
	var (
		isController = strings.HasSuffix(p.PkgPath, "/controllers") ||
			strings.Contains(p.PkgPath, "/controllers/")
		isTest = strings.HasSuffix(p.PkgPath, "/tests") ||
			strings.Contains(p.PkgPath, "/tests/")
		methodMap = map[string][]*model.MethodSpec{}
	)
	localImportMap := map[string]string{}
	log := s.sourceProcessor.log.New("package", p.PkgPath)
	log.Info("Processing package")
	for _, tree := range p.Syntax {
		for _, decl := range tree.Decls {

			s.sourceProcessor.packageMap[p.PkgPath] = filepath.Dir(p.Fset.Position(decl.Pos()).Filename)
			if !s.addImport(decl, p, localImportMap, log) {
				continue
			}
			spec, found := s.getStructTypeDecl(decl, p.Fset)
			//log.Info("Checking file","filename", p.Fset.Position(decl.Pos()).Filename,"found",found)
			if found {
				if isController || isTest {
					controllerSpec := s.getControllerSpec(spec, p, localImportMap)
					sourceInfo.StructSpecs = append(sourceInfo.StructSpecs, controllerSpec)
				}
			} else {
				// Not a type definition, this could be a method for a controller try to extract that
				// Func declaration?
				funcDecl, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				// This could be a controller action endpoint, check and add if needed
				if isController &&
					funcDecl.Recv != nil && // Must have a receiver
					funcDecl.Name.IsExported() && // be public
					funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) == 1 {
					// return one result
					if m, receiver := s.getControllerFunc(funcDecl, p, localImportMap); m != nil {
						methodMap[receiver] = append(methodMap[receiver], m)
						log.Info("Added method map to ", "receiver", receiver, "method", m.Name)
					}
				}
				// Check for validation
				if lineKeyMap := s.getValidation(funcDecl, p); len(lineKeyMap) > 1 {
					sourceInfo.ValidationKeys[p.PkgPath + "." + s.getFuncName(funcDecl)] = lineKeyMap
				}
				if funcDecl.Name.Name == "init" {
					sourceInfo.InitImportPaths = append(sourceInfo.InitImportPaths, p.PkgPath)
				}
			}
		}
	}

	// Add the method specs to the struct specs.
	for _, spec := range sourceInfo.StructSpecs {
		spec.MethodSpecs = methodMap[spec.StructName]
	}

	return
}
// Scan app source code for calls to X.Y(), where X is of type *Validation.
//
// Recognize these scenarios:
// - "Y" = "Validation" and is a member of the receiver.
//   (The common case for inline validation)
// - "X" is passed in to the func as a parameter.
//   (For structs implementing Validated)
//
// The line number to which a validation call is attributed is that of the
// surrounding ExprStmt.  This is so that it matches what runtime.Callers()
// reports.
//
// The end result is that we can set the default validation key for each call to
// be the same as the local variable.
func (s *SourceInfoProcessor) getValidation(funcDecl *ast.FuncDecl, p *packages.Package) (map[int]string) {
	var (
		lineKeys = make(map[int]string)

		// Check the func parameters and the receiver's members for the *revel.Validation type.
		validationParam = s.getValidationParameter(funcDecl)
	)

	ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
		// e.g. c.Validation.Required(arg) or v.Required(arg)
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		// e.g. c.Validation.Required or v.Required
		funcSelector, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		switch x := funcSelector.X.(type) {
		case *ast.SelectorExpr: // e.g. c.Validation
			if x.Sel.Name != "Validation" {
				return true
			}

		case *ast.Ident: // e.g. v
			if validationParam == nil || x.Obj != validationParam {
				return true
			}

		default:
			return true
		}

		if len(callExpr.Args) == 0 {
			return true
		}

		// Given the validation expression, extract the key.
		key := callExpr.Args[0]
		switch expr := key.(type) {
		case *ast.BinaryExpr:
			// If the argument is a binary expression, take the first expression.
			// (e.g. c.Validation.Required(myName != ""))
			key = expr.X
		case *ast.UnaryExpr:
			// If the argument is a unary expression, drill in.
			// (e.g. c.Validation.Required(!myBool)
			key = expr.X
		case *ast.BasicLit:
			// If it's a literal, skip it.
			return true
		}

		if typeExpr := model.NewTypeExprFromAst("", key); typeExpr.Valid {
			lineKeys[p.Fset.Position(callExpr.End()).Line] = typeExpr.TypeName("")
		} else {
			s.sourceProcessor.log.Error("Error: Failed to generate key for field validation. Make sure the field name is valid.", "file", p.PkgPath,
				"line", p.Fset.Position(callExpr.End()).Line, "function", funcDecl.Name.String())
		}
		return true
	})

	return lineKeys

}
// Check to see if there is a *revel.Validation as an argument.
func (s *SourceInfoProcessor)  getValidationParameter(funcDecl *ast.FuncDecl) *ast.Object {
	for _, field := range funcDecl.Type.Params.List {
		starExpr, ok := field.Type.(*ast.StarExpr) // e.g. *revel.Validation
		if !ok {
			continue
		}

		selExpr, ok := starExpr.X.(*ast.SelectorExpr) // e.g. revel.Validation
		if !ok {
			continue
		}

		xIdent, ok := selExpr.X.(*ast.Ident) // e.g. rev
		if !ok {
			continue
		}

		if selExpr.Sel.Name == "Validation" && s.sourceProcessor.importMap[xIdent.Name] == model.RevelImportPath {
			return field.Names[0].Obj
		}
	}
	return nil
}
func (s *SourceInfoProcessor) getControllerFunc(funcDecl *ast.FuncDecl, p *packages.Package, localImportMap map[string]string) (method *model.MethodSpec, recvTypeName string) {
	selExpr, ok := funcDecl.Type.Results.List[0].Type.(*ast.SelectorExpr)
	if !ok {
		return
	}
	if selExpr.Sel.Name != "Result" {
		return
	}
	if pkgIdent, ok := selExpr.X.(*ast.Ident); !ok || s.sourceProcessor.importMap[pkgIdent.Name] != model.RevelImportPath {
		return
	}
	method = &model.MethodSpec{
		Name: funcDecl.Name.Name,
	}

	// Add a description of the arguments to the method.
	for _, field := range funcDecl.Type.Params.List {
		for _, name := range field.Names {
			var importPath string
			typeExpr := model.NewTypeExprFromAst(p.Name, field.Type)
			if !typeExpr.Valid {
				utils.Logger.Warn("Warn: Didn't understand argument '%s' of action %s. Ignoring.", name, s.getFuncName(funcDecl))
				return // We didn't understand one of the args.  Ignore this action.
			}
			// Local object
			if typeExpr.PkgName == p.Name {
				importPath = p.PkgPath
			} else if typeExpr.PkgName != "" {
				var ok bool
				if importPath, ok = localImportMap[typeExpr.PkgName]; !ok {
					if importPath, ok = s.sourceProcessor.importMap[typeExpr.PkgName]; !ok {
						utils.Logger.Error("Unable to find import", "importMap", s.sourceProcessor.importMap, "localimport", localImportMap)
						utils.Logger.Fatalf("Failed to find import for arg of type: %s , %s", typeExpr.PkgName, typeExpr.TypeName(""))
					}
				}
			}
			method.Args = append(method.Args, &model.MethodArg{
				Name:       name.Name,
				TypeExpr:   typeExpr,
				ImportPath: importPath,
			})
		}
	}

	// Add a description of the calls to Render from the method.
	// Inspect every node (e.g. always return true).
	method.RenderCalls = []*model.MethodCall{}
	ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
		// Is it a function call?
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Is it calling (*Controller).Render?
		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// The type of the receiver is not easily available, so just store every
		// call to any method called Render.
		if selExpr.Sel.Name != "Render" {
			return true
		}

		// Add this call's args to the renderArgs.
		pos := p.Fset.Position(callExpr.Lparen)
		methodCall := &model.MethodCall{
			Line:  pos.Line,
			Names: []string{},
		}
		for _, arg := range callExpr.Args {
			argIdent, ok := arg.(*ast.Ident)
			if !ok {
				continue
			}
			methodCall.Names = append(methodCall.Names, argIdent.Name)
		}
		method.RenderCalls = append(method.RenderCalls, methodCall)
		return true
	})

	var recvType = funcDecl.Recv.List[0].Type
	if recvStarType, ok := recvType.(*ast.StarExpr); ok {
		recvTypeName = recvStarType.X.(*ast.Ident).Name
	} else {
		recvTypeName = recvType.(*ast.Ident).Name
	}
	return
}
func (s *SourceInfoProcessor) getControllerSpec(spec *ast.TypeSpec, p *packages.Package, localImportMap map[string]string) (controllerSpec *model.TypeInfo) {
	structType := spec.Type.(*ast.StructType)

	// At this point we know it's a type declaration for a struct.
	// Fill in the rest of the info by diving into the fields.
	// Add it provisionally to the Controller list -- it's later filtered using field info.
	controllerSpec = &model.TypeInfo{
		StructName:  spec.Name.Name,
		ImportPath:  p.PkgPath,
		PackageName: p.Name,
	}
	log := s.sourceProcessor.log.New("file", p.Fset.Position(spec.Pos()).Filename, "position", p.Fset.Position(spec.Pos()).Line)
	for _, field := range structType.Fields.List {
		// If field.Names is set, it's not an embedded type.
		if field.Names != nil {
			continue
		}

		// A direct "sub-type" has an ast.Field as either:
		//   Ident { "AppController" }
		//   SelectorExpr { "rev", "Controller" }
		// Additionally, that can be wrapped by StarExprs.
		fieldType := field.Type
		pkgName, typeName := func() (string, string) {
			// Drill through any StarExprs.
			for {
				if starExpr, ok := fieldType.(*ast.StarExpr); ok {
					fieldType = starExpr.X
					continue
				}
				break
			}

			// If the embedded type is in the same package, it's an Ident.
			if ident, ok := fieldType.(*ast.Ident); ok {
				return "", ident.Name
			}

			if selectorExpr, ok := fieldType.(*ast.SelectorExpr); ok {
				if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok {
					return pkgIdent.Name, selectorExpr.Sel.Name
				}
			}
			return "", ""
		}()

		// If a typename wasn't found, skip it.
		if typeName == "" {
			continue
		}

		// Find the import path for this type.
		// If it was referenced without a package name, use the current package import path.
		// Else, look up the package's import path by name.
		var importPath string
		if pkgName == "" {
			importPath = p.PkgPath
		} else {
			var ok bool
			if importPath, ok = localImportMap[pkgName]; !ok {
				log.Debug("Debug: Unusual, failed to find package locally ", "package", pkgName, "type", typeName, "map", s.sourceProcessor.importMap, "usedin", )
				if importPath, ok = s.sourceProcessor.importMap[pkgName]; !ok {
					log.Error("Error: Failed to find import path for ", "package", pkgName, "type", typeName, "map", s.sourceProcessor.importMap, "usedin", )
					continue
				}
			}
		}

		controllerSpec.EmbeddedTypes = append(controllerSpec.EmbeddedTypes, &model.EmbeddedTypeName{
			ImportPath: importPath,
			StructName: typeName,
		})
	}
	s.sourceProcessor.log.Info("Added controller spec", "name", controllerSpec.StructName, "package", controllerSpec.ImportPath)
	return
}
func (s *SourceInfoProcessor) getStructTypeDecl(decl ast.Decl, fset *token.FileSet) (spec *ast.TypeSpec, found bool) {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return
	}

	if genDecl.Tok != token.TYPE {
		return
	}

	if len(genDecl.Specs) == 0 {
		utils.Logger.Warn("Warn: Surprising: %s:%d Decl contains no specifications", fset.Position(decl.Pos()).Filename, fset.Position(decl.Pos()).Line)
		return
	}

	spec = genDecl.Specs[0].(*ast.TypeSpec)
	_, found = spec.Type.(*ast.StructType)

	return

}
func (s *SourceInfoProcessor) getFuncName(funcDecl *ast.FuncDecl) string {
	prefix := ""
	if funcDecl.Recv != nil {
		recvType := funcDecl.Recv.List[0].Type
		if recvStarType, ok := recvType.(*ast.StarExpr); ok {
			prefix = "(*" + recvStarType.X.(*ast.Ident).Name + ")"
		} else {
			prefix = recvType.(*ast.Ident).Name
		}
		prefix += "."
	}
	return prefix + funcDecl.Name.Name
}
func (s *SourceInfoProcessor) addImport(decl ast.Decl, p *packages.Package, localImportMap map[string]string, log logger.MultiLogger) (shouldContinue bool) {
	shouldContinue = true
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return
	}

	if genDecl.Tok == token.IMPORT {
		shouldContinue = false
		for _, spec := range genDecl.Specs {
			importSpec := spec.(*ast.ImportSpec)
			//fmt.Printf("*** import specification %#v\n", importSpec)
			var pkgAlias string
			if importSpec.Name != nil {
				pkgAlias = importSpec.Name.Name
				if pkgAlias == "_" {
					continue
				}
			}
			quotedPath := importSpec.Path.Value           // e.g. "\"sample/app/models\""
			fullPath := quotedPath[1 : len(quotedPath) - 1] // Remove the quotes
			if pkgAlias == "" {
				pkgAlias = fullPath
				if index := strings.LastIndex(pkgAlias, "/"); index > 0 {
					pkgAlias = pkgAlias[index + 1:]
				}
			}
			localImportMap[pkgAlias] = fullPath
		}

	}
	return
}