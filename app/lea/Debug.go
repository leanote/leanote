package lea

import (
	"encoding/json"
	"github.com/revel/revel"
	"fmt"
)

func Log(i interface{}) {
	revel.INFO.Println(i)
}

func LogJ(i interface{}) {
	b, _ := json.MarshalIndent(i, "", "	") 
	revel.INFO.Println(string(b))
}

// 为test用
func L(i interface{}) {
	fmt.Println(i)
}

func LJ(i interface{}) {
	b, _ := json.MarshalIndent(i, "", "	") 
	fmt.Println(string(b))
}
