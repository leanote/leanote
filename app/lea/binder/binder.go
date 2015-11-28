package binder

import (
	"github.com/leanote/leanote/app/info"
	"github.com/revel/revel"
	//	"github.com/leanote/leanote/app/controllers/api"
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
	// name == "noteOrContent"
	Bind: func(params *revel.Params, name string, typ reflect.Type) reflect.Value {
		result := reflect.New(typ).Elem() // 创建一个该类型的, 然后其field从所有的param去取
		fieldValues := make(map[string]reflect.Value)
		//		fmt.Println(name)
		// fmt.Println(typ) // api.NoteFiles
		// name = files[0], files[1], noteContent
		//		fmt.Println(params.Values)
		/*
			map[Title:[test1] METHOD:[POST] NotebookId:[54c4f51705fcd14031000002]
			files[1][FileId]:[]
			controller:[note]
			files[1][LocalFileId]:[54c7ae27d98d0329dd000000] files[1][HasBody]:[true] files[0][FileId]:[] files[0][LocalFileId]:[54c7ae855e94ea2dba000000] action:[addNote] Content:[<p>lifedddddd</p><p><img src="app://leanote/data/54bdc65599c37b0da9000002/images/1422368307147_2.png" alt="" data-mce-src="app://leanote/data/54bdc65599c37b0da9000002/images/1422368307147_2.png" style="display: block; margin-left: auto; margin-right: auto;"></p><p><img src="http://127.0.0.1:8008/api/file/getImage?fileId=54c7ae27d98d0329dd000000" alt="" data-mce-src="http://127.0.0.1:8008/api/file/getImg?fileId=54c7ae27d98d0329dd000000"></p><p><br></p><p><img src="http://127.0.0.1:8008/api/file/getImage?fileId=54c7ae855e94ea2dba000000" alt="" data-mce-src="http://127.0.0.1:8008/api/file/getImage?fileId=54c7ae855e94ea2dba000000" style="display: block; margin-left: auto; margin-right: auto;"></p><p><br></p><p><br></p>] IsBlog:[false] token:[user1]
			files[0][HasBody]:[true]]
		*/
		nameIsSlice := strings.Contains(name, "[")
		//		fmt.Println(params.Values["files[1]"])
		//		fmt.Println(params.Values["Title"])
		for key, _ := range params.Values { // Title, Content, Files
			// 这里, 如果没有点, 默认就是a.
			// life
			//			fmt.Println("key:" + key); // files[0][LocalFileId]
			//			fmt.Println("name:" + name); // files[0][LocalFileId]
			var suffix string
			var noPrefix = false
			if nameIsSlice && strings.HasPrefix(key, name) {
				suffix = key[len(name)+1 : len(key)-1] // files[0][LocalFileId] 去掉 => LocalFileId
			} else if !strings.HasPrefix(key, name+".") {
				noPrefix = true
				suffix = key
				//			continue
			} else {
				// Get the name of the struct property.
				// Strip off the prefix. e.g. foo.bar.baz => bar.baz
				suffix = key[len(name)+1:]
			}
			//			fmt.Println(suffix);

			fieldName := nextKey(suffix) // e.g. bar => "bar", bar.baz => "bar", bar[0] => "bar"
			//			fmt.Println(fieldName);
			fieldLen := len(fieldName)

			if _, ok := fieldValues[fieldName]; !ok {
				// Time to bind this field.  Get it and make sure we can set it.
				fieldName = strings.Title(fieldName) // 传过来title, 但struct是Title
				//				fmt.Println("xx: " + fieldName)
				fieldValue := result.FieldByName(fieldName)
				//				fmt.Println(fieldValue)
				if !fieldValue.IsValid() {
					continue
				}
				if !fieldValue.CanSet() {
					continue
				}
				var boundVal reflect.Value
				// 没有name前缀
				if noPrefix {
					// life
					//					fmt.Println("<<")
					//					fmt.Println(strings.Title(key[:fieldLen]));
					boundVal = revel.Bind(params, key[:fieldLen], fieldValue.Type())
				} else {
					//					fmt.Println("final")
					//					fmt.Println(key[:len(name)+1+fieldLen]) // files[0][HasBody
					if nameIsSlice {
						fieldLen += 1
					}
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
	revel.TypeBinders[reflect.TypeOf(info.UserAccount{})] = leanoteStructBinder
	revel.TypeBinders[reflect.TypeOf(info.NoteOrContent{})] = leanoteStructBinder
	revel.TypeBinders[reflect.TypeOf(info.ApiNote{})] = leanoteStructBinder
	revel.TypeBinders[reflect.TypeOf(info.NoteFile{})] = leanoteStructBinder
}
