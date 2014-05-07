package binder

import (
	"github.com/revel/revel"
	"github.com/leanote/leanote/app/info"
	"github.com/leanote/leanote/app/controllers"
	"fmt"
	"reflect"
	"strings"
)

// leanote binder struct
// rewrite revel struct binder
// not need the struct name as prefix, 
// eg:
// type Note struct {Name}
// func (c Controller) List(note Note) revel.Result {}
// in revel you must pass the note.Name as key, now you just pass Name

// for test
// map[string][string]
var MSSBinder = revel.Binder{
	Bind: func(params *revel.Params, name string, typ reflect.Type) reflect.Value {
		var (
			result    = reflect.MakeMap(typ)
			keyType   = typ.Key()
			valueType = typ.Elem()
		)
		for paramName, values := range params.Values {
			key := paramName // [len(name)+1 : len(paramName)-1]
			// 每一个值 values[0]
			result.SetMapIndex(revel.BindValue(key, keyType), revel.BindValue(values[0], valueType))
		}
		return result
	},
	
	Unbind: func(output map[string]string, name string, val interface{}) {
		mapValue := reflect.ValueOf(val)
		for _, key := range mapValue.MapKeys() {
			revel.Unbind(output, fmt.Sprintf("%v", key.Interface()),
				mapValue.MapIndex(key).Interface())
		}	
	},
}

// struct需要., a.b = "life"
// a: contoller是形参名
// 修改, 默认就是a.b, 传b
func nextKey(key string) string {
	fieldLen := strings.IndexAny(key, ".[")
	if fieldLen == -1 {
		return key
	}
	return key[:fieldLen]
}
var leanoteStructBinder = revel.Binder{
	Bind: func(params *revel.Params, name string, typ reflect.Type) reflect.Value {
		result := reflect.New(typ).Elem()
		fieldValues := make(map[string]reflect.Value)
		for key, _ := range params.Values {
			// 这里, 如果没有点, 默认就是a.
			// life
			var suffix string
			var noPrefix = false
			if !strings.HasPrefix(key, name + ".") {
				noPrefix = true
				suffix = key
	//			continue
			} else {
				// Get the name of the struct property.
				// Strip off the prefix. e.g. foo.bar.baz => bar.baz
				suffix = key[len(name)+1:]
			}
	
			fieldName := nextKey(suffix) // e.g. bar => "bar", bar.baz => "bar", bar[0] => "bar"
			fieldLen := len(fieldName)
			
			if _, ok := fieldValues[fieldName]; !ok {
				// Time to bind this field.  Get it and make sure we can set it.
				fieldValue := result.FieldByName(fieldName)
				if !fieldValue.IsValid() {
					continue
				}
				if !fieldValue.CanSet() {
					continue
				}
				var boundVal reflect.Value
				// 没有name前缀
				if(noPrefix) {
					boundVal = revel.Bind(params, key[:fieldLen], fieldValue.Type())
				} else {
					boundVal = revel.Bind(params, key[:len(name)+1+fieldLen], fieldValue.Type())
				}
				fieldValue.Set(boundVal)
				fieldValues[fieldName] = boundVal
			}
		}
		return result
	},
	Unbind: func(output map[string]string, name string, iface interface{}) {
		val := reflect.ValueOf(iface)
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			structField := typ.Field(i)
			fieldValue := val.Field(i)
			// PkgPath is specified to be empty exactly for exported fields.
			if structField.PkgPath == "" {
				revel.Unbind(output, fmt.Sprintf("%s.%s", name, structField.Name), fieldValue.Interface())
			}
		}
	},
}

func init() {
	revel.TypeBinders[reflect.TypeOf(info.UserBlogBase{})] = leanoteStructBinder
	revel.TypeBinders[reflect.TypeOf(info.UserBlogComment{})] = leanoteStructBinder
	revel.TypeBinders[reflect.TypeOf(info.UserBlogStyle{})] = leanoteStructBinder
	revel.TypeBinders[reflect.TypeOf(info.Notebook{})] = leanoteStructBinder
	revel.TypeBinders[reflect.TypeOf(controllers.NoteOrContent{})] = leanoteStructBinder
}