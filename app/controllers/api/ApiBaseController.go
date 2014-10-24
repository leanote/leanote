package api

import (
//	"github.com/revel/revel"
//	"gopkg.in/mgo.v2/bson"
//	"encoding/json"
//	. "github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/controllers"
//	"io/ioutil"
//	"fmt"
//	"math"
//	"strconv"
//	"strings"
)

// 公用Controller, 其它Controller继承它
type ApiBaseContrller struct {
	controllers.BaseController // 不能用*BaseController
}