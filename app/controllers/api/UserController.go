package api

import (
	"github.com/revel/revel"
//	"encoding/json"
//	"gopkg.in/mgo.v2/bson"
	. "github.com/leanote/leanote/app/lea"
//	"github.com/leanote/leanote/app/info"
//	"github.com/leanote/leanote/app/types"
//	"io/ioutil"
//	"fmt"
//	"math"
//	"os"
//	"path"
//	"strconv"
)

type ApiUser struct {
	*revel.Controller
}

// 修改用户名, 需要重置session
func (c ApiUser) Info() revel.Result {
	Log("APIUser");
	return c.RenderTemplate("home/index.html");
//	return nil;
}