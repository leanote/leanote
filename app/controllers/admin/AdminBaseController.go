package admin

import (
//	"github.com/revel/revel"
//	"gopkg.in/mgo.v2/bson"
//	"encoding/json"
	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/controllers"
//	"io/ioutil"
//	"fmt"
//	"math"
//	"strconv"
	"strings"
)

// 公用Controller, 其它Controller继承它
type AdminBaseController struct {
	controllers.BaseController // 不能用*BaseController
}

// 得到sorterField 和 isAsc
// okSorter = ['email', 'username']
func (c AdminBaseController) getSorter(sorterField string, isAsc bool, okSorter []string) (string, bool){
	sorter := ""
	c.Params.Bind(&sorter, "sorter")
	if sorter == "" {
		return sorterField, isAsc;
	}
	
	// sorter形式 email-up, email-down
	s2 := strings.Split(sorter, "-")
	if len(s2) != 2 {
		return sorterField, isAsc;
	}
	
	// 必须是可用的sorter
	if okSorter != nil && len(okSorter) > 0 {
		if !InArray(okSorter, s2[0]) {
			return sorterField, isAsc;
		}
	}
	
	sorterField = strings.Title(s2[0])
	if s2[1] == "up" {
		isAsc = true
	} else {
		isAsc = false
	}
	c.RenderArgs["sorter"] = sorter
	return sorterField, isAsc;
}