package pongo2

import (
	p2 "github.com/flosch/pongo2"
	"github.com/revel/revel"
)

func init() {
	p2.RegisterFilter("field", func(in *p2.Value, param *p2.Value) (out *p2.Value, err *p2.Error) {
		if nil == in.Interface() || in.String() == "" {
			return nil, &p2.Error{Sender: "field argument must is string"}
		}
		newField := revel.NewField(in.String(), getContext())
		return p2.AsValue(newField), nil
	})
}
