package lea

import (
    "encoding/json"
    "fmt"
    "github.com/revel/revel"
)

func Log(msg string, i ...interface{}) {
    revel.AppLog.Info(msg, i...)
}

func LogW(msg string, i ...interface{}) {
    revel.AppLog.Warn(msg, i...)
}

func LogJ(i interface{}) {
    b, _ := json.MarshalIndent(i, "", " ")
    revel.AppLog.Info(string(b))
}

// 为test用
func L(i interface{}) {
    fmt.Println(i)
}

func LJ(i interface{}) {
    b, _ := json.MarshalIndent(i, "", " ")
    fmt.Println(string(b))
}
