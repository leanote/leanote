package pongo2

import (
	p2 "github.com/flosch/pongo2"
	"reflect"
)

type tagAppendNode struct {
	isPublic         bool
	name             string
	objectEvaluators []p2.IEvaluator
}

func (node *tagAppendNode) Execute(ctx *p2.ExecutionContext, writer p2.TemplateWriter) *p2.Error {
	var values []interface{}
	var reflectValues reflect.Value
	found := false
	if o, found := ctx.Public[node.name]; found && nil != o {
		values, _ = o.([]interface{})
		if nil == values {
			reflectValues = reflect.ValueOf(o)
			if reflectValues.Kind() == reflect.Ptr {
				reflectValues = reflectValues.Elem()
			}
			if reflectValues.Kind() != reflect.Slice {
				return &p2.Error{Sender: "'" + node.name + "' isn't a slice."}
			}
		}
	}

	for _, ev := range node.objectEvaluators {
		obj, err := ev.Evaluate(ctx)
		if err != nil {
			return err
		}
		if !found || reflectValues.IsNil() {
			values = append(values, obj)
		} else {
			reflectValues = reflect.AppendSlice(reflectValues, reflect.ValueOf(obj))
		}
	}

	if !found || reflectValues.IsNil() {
		ctx.Public[node.name] = values
	} else {
		ctx.Public[node.name] = reflectValues.Interface()
	}
	return nil
}

// tagURLForParser implements a {% urlfor %} tag.
//
// urlfor takes one argument for the controller, as well as any number of key/value pairs for additional URL data.
// Example: {% urlfor "UserController.View" ":slug" "oal" %}
func tagAppendParser(doc *p2.Parser, start *p2.Token, arguments *p2.Parser) (p2.INodeTag, *p2.Error) {
	var name string
	var isPublic bool
	// Parse variable name
	typeToken := arguments.MatchType(p2.TokenIdentifier)
	if typeToken != nil {
		name = typeToken.Val
	} else if sToken := arguments.MatchType(p2.TokenString); nil != sToken {
		name = sToken.Val
	} else {
		return nil, arguments.Error("Expected an identifier or string.", nil)
	}

	evals := []p2.IEvaluator{}
	for arguments.Remaining() > 0 {
		expr, err := arguments.ParseExpression()
		if err != nil {
			return nil, err
		}
		evals = append(evals, expr)
	}

	return &INodeImplied{Exec: func(ctx *p2.ExecutionContext, w p2.TemplateWriter) *p2.Error {

		node := &tagAppendNode{isPublic, name, evals}
		return node.Execute(ctx, w)
	}}, nil

}

func init() {
	p2.RegisterTag("append", tagAppendParser)
}
