package pongo2

import (
	"github.com/revel/revel"

	p2 "github.com/flosch/pongo2"
)

type tagURLForNode struct {
	objectEvaluators []p2.IEvaluator
}

func (node *tagURLForNode) Execute(ctx *p2.ExecutionContext, writer p2.TemplateWriter) *p2.Error {
	args := make([]string, len(node.objectEvaluators))
	for i, ev := range node.objectEvaluators {
		obj, err := ev.Evaluate(ctx)
		if err != nil {
			return err
		}
		args[i] = obj.String()
	}

	params := make([]interface{}, len(args))
	params[0] = args[0]
	for i := range params[1:] {
		params[i+1] = args[i+1]
	}

	url, err := revel.ReverseURL(params...)
	if nil != err {
		return ctx.Error(err.Error(), nil)
	}

	writer.WriteString(string(url))
	return nil
}

// tagURLForParser implements a {% urlfor %} tag.
//
// urlfor takes one argument for the controller, as well as any number of key/value pairs for additional URL data.
// Example: {% urlfor "UserController.View" ":slug" "oal" %}
func tagURLForParser(doc *p2.Parser, start *p2.Token, arguments *p2.Parser) (p2.INodeTag, *p2.Error) {
	evals := []p2.IEvaluator{}
	for arguments.Remaining() > 0 {
		expr, err := arguments.ParseExpression()
		evals = append(evals, expr)
		if err != nil {
			return nil, err
		}
	}

	if len(evals) <= 0 {
		return nil, arguments.Error("URL takes one argument for the controller and any number of optional value.", nil)
	}
	return &INodeImplied{Exec: func(ctx *p2.ExecutionContext, w p2.TemplateWriter) *p2.Error {

		v := &tagURLForNode{evals}
		return v.Execute(ctx, w)
	}}, nil
}

func init() {
	p2.RegisterTag("url", tagURLForParser)
}
