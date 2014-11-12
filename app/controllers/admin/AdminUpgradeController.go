package admin

import (
	"github.com/revel/revel"
//	"encoding/json"
	"github.com/leanote/leanote/app/info"
//	"io/ioutil"
)

// Upgrade controller
type AdminUpgrade struct {
	AdminBaseController
}

func (c AdminUpgrade) UpgradeBlog() revel.Result {
	upgradeService.UpgradeBlog()
	return nil;
}

func (c AdminUpgrade) UpgradeBetaToBeta2() revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = upgradeService.UpgradeBetaToBeta2(c.GetUserId())
	return c.RenderJson(re)
}