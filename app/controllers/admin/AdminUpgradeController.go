package admin

import (
	"github.com/revel/revel"
//	"encoding/json"
//	. "github.com/leanote/leanote/app/lea"
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