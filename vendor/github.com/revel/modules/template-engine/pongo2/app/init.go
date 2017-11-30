package pongo2

import (
	"io"
	"strings"

	p2 "github.com/flosch/pongo2"
	"github.com/revel/revel"
	"github.com/tylerb/gls"
)

// Adapter for HAML Templates.
type PongoTemplate struct {
	template *p2.Template
	engine   *PongoEngine
	*revel.TemplateView
}
type Pongo2BaseTag struct {
	field string
}

func (node *Pongo2BaseTag) GetField(ctx *p2.ExecutionContext) (value interface{}, found bool) {
	value, found = ctx.Public[node.field]
	if !found {
		value, found = ctx.Private[node.field]
	}
	if found {
		if wrapped, ok := value.(*p2.Value); ok {
			value = wrapped.Interface()
		}
	}
	return
}

type INodeImplied struct {
	Exec func(*p2.ExecutionContext, p2.TemplateWriter) *p2.Error
}

func (i *INodeImplied) Execute(ctx *p2.ExecutionContext, w p2.TemplateWriter) *p2.Error {
	return i.Exec(ctx, w)

}
func (tmpl PongoTemplate) Name() string {
	return tmpl.TemplateName
}
func getContext() map[string]interface{} {
	return gls.Get("data").(map[string]interface{})
}

// return a 'revel.Template' from HAML's template.
func (tmpl PongoTemplate) Render(wr io.Writer, arg interface{}) (err error) {
	gls.With(gls.Values(map[interface{}]interface{}{"data": arg}), func() {
		err = tmpl.template.ExecuteWriter(p2.Context(arg.(map[string]interface{})), wr)
		if nil != err {
			if e, ok := err.(*p2.Error); ok {
				rerr := &revel.Error{
					Title:       "Template Execution Error",
					Path:        tmpl.TemplateName,
					Description: e.Error(),
					Line:        e.Line,
				}
				if revel.DevMode {
					rerr.SourceLines = tmpl.Content()
				}
				err = rerr
			}
		}
	})
	return err
}

// There is only a single instance of the PongoEngine initialized
type PongoEngine struct {
	loader                *revel.TemplateLoader
	templateSetBybasePath map[string]*p2.TemplateSet
	templates             map[string]*PongoTemplate
	CaseInsensitive       bool
}

func (i *PongoEngine) ConvertPath(path string) string {
	if i.CaseInsensitive {
		return strings.ToLower(path)
	}
	return path
}

func (i *PongoEngine) Handles(templateView *revel.TemplateView) bool {
	return revel.EngineHandles(i, templateView)
}

func (engine *PongoEngine) ParseAndAdd(baseTemplate *revel.TemplateView) error {
	templateSet := engine.templateSetBybasePath[baseTemplate.BasePath]
	if nil == templateSet {
		templateSet = p2.NewSet(baseTemplate.BasePath, p2.MustNewLocalFileSystemLoader(baseTemplate.BasePath))
		engine.templateSetBybasePath[baseTemplate.BasePath] = templateSet
	}

	tpl, err := templateSet.FromBytes(baseTemplate.FileBytes)
	if nil != err {
		_, line, description := parsePongo2Error(err)
		return &revel.Error{
			Title:       "Template Compilation Error",
			Path:        baseTemplate.FilePath,
			Description: description,
			Line:        line,
			SourceLines: strings.Split(string(baseTemplate.FileBytes), "\n"),
		}
	}
	baseTemplate.TemplateName = engine.ConvertPath(baseTemplate.TemplateName)
	engine.templates[baseTemplate.TemplateName] = &PongoTemplate{
		template:     tpl,
		engine:       engine,
		TemplateView: baseTemplate}
	return nil
}
func (engine *PongoEngine) Name() string {
	return "pongo2"
}

func parsePongo2Error(err error) (templateName string, line int, description string) {
	pongoError := err.(*p2.Error)
	if nil != pongoError {
		return pongoError.Filename, pongoError.Line, pongoError.Error()
	}
	return "Unknown error", 0, err.Error()
}

func (engine *PongoEngine) Lookup(templateName string) revel.Template {
	tpl, found := engine.templates[engine.ConvertPath(templateName)]
	if !found {
		return nil
	}
	return tpl
}
func (engine *PongoEngine) Event(action int, i interface{}) {
	if action == revel.TEMPLATE_REFRESH_REQUESTED {
		// At this point all the templates have been passed into the
		engine.templateSetBybasePath = map[string]*p2.TemplateSet{}
		engine.templates = map[string]*PongoTemplate{}
		engine.CaseInsensitive = revel.Config.BoolDefault("pongo2.template.caseinsensitive", true)
	}
}

func init() {
	revel.RegisterTemplateLoader("pongo2", func(loader *revel.TemplateLoader) (revel.TemplateEngine, error) {
		return &PongoEngine{
			loader:                loader,
			templateSetBybasePath: map[string]*p2.TemplateSet{},
			templates:             map[string]*PongoTemplate{},
		}, nil
	})
	/*
	   // TODO Dynamically call all the built in functions, PR welcome
	   for key,templateFunction := range revel.TemplateFuncs {
	       p2.RegisterTag(key, func(doc *p2.Parser, start *p2.Token, arguments *p2.Parser) (p2.INodeTag, *p2.Error) {
	           evals := []p2.IEvaluator{}
	           for arguments.Remaining() > 0 {
	               expr, err := arguments.ParseExpression()
	               evals = append(evals, expr)
	               if err != nil {
	                   return  nil, err
	               }
	           }

	       return &INodeImplied{Exec: func(ctx *p2.ExecutionContext,w p2.TemplateWriter) *p2.Error {
	           args := make([]interface{}, len(evals))
	           for i, ev := range evals {
	               obj, err := ev.Evaluate(ctx)
	               if err != nil {
	                   return err
	               }
	               args[i] = obj
	           }

	           v:= &tagURLForNode{evals}
	           reflect.MakeFunc ....
	           return v.Execute(ctx,w)
	       }}, nil

	       })

	   }
	*/
}
