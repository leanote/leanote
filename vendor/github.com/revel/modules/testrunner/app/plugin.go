package app

import (
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(func() {
		revel.AppLog.Info("Go to /@tests to run the tests.")
	})
}
