package pongo2

import (
	"fmt"
	"html"

	p2 "github.com/flosch/pongo2"
	"github.com/revel/revel"
)

type tagOptionNode struct {
	value p2.IEvaluator
	label string
	Pongo2BaseTag
}

func (node *tagOptionNode) Execute(ctx *p2.ExecutionContext, writer p2.TemplateWriter) *p2.Error {
	fieldObj, _ := node.GetField(ctx)

	if nil == fieldObj {
		return ctx.Error("field '"+node.field+"' tagOptionNode is missing.", nil)
	}
	field, _ := fieldObj.(*revel.Field)
	if nil == field {
		return ctx.Error(fmt.Sprintf("field '"+node.field+"' isn't Field - %T.", fieldObj), nil)
	}

	val, err := node.value.Evaluate(ctx)
	if err != nil {
		return err
	}
	val_str := val.String()

	selected := ""
	if field.Flash() == val_str || (field.Flash() == "" && field.Value() == val_str) {
		selected = " selected"
	}

	fmt.Fprintf(writer, `<option value="%s"%s>%s</option>`,
		html.EscapeString(val_str), selected, html.EscapeString(node.label))
	return nil
}

// tagURLForParser implements a {% urlfor %} tag.
//
// urlfor takes one argument for the controller, as well as any number of key/value pairs for additional URL data.
// Example: {% urlfor "UserController.View" ":slug" "oal" %}
func tagOptionParser(doc *p2.Parser, start *p2.Token, arguments *p2.Parser) (p2.INodeTag, *p2.Error) {
	var field string
	typeToken := arguments.MatchType(p2.TokenIdentifier)
	if typeToken != nil {
		field = typeToken.Val
	} else if sToken := arguments.MatchType(p2.TokenString); nil != sToken {
		field = sToken.Val
	} else {
		return nil, arguments.Error("Expected an identifier or string.", nil)
	}

	expr, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}

	var v *tagOptionNode
	if sToken := arguments.MatchType(p2.TokenString); nil != sToken {
		v = &tagOptionNode{Pongo2BaseTag: Pongo2BaseTag{field: field},
			value: expr,
			label: sToken.Val}
	} else {
		return nil, arguments.Error("Expected an string.", nil)
	}
	return &INodeImplied{Exec: func(ctx *p2.ExecutionContext, w p2.TemplateWriter) *p2.Error {
		return v.Execute(ctx, w)
	}}, nil
}

func init() {
	p2.RegisterTag("option", tagOptionParser)
}
