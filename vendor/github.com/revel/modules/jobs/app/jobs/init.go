package jobs

import (
	"github.com/revel/revel"
)

var jobLog = revel.AppLog

func init() {
	revel.RegisterModuleInit(func(m *revel.Module){
		jobLog = m.Log
	})
}
