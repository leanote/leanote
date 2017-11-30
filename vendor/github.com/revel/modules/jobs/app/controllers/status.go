package controllers

import (
	"strings"

	"github.com/revel/cron"
	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
)

type Jobs struct {
	*revel.Controller
}

func (c Jobs) Status() revel.Result {
	remoteAddress := c.Request.RemoteAddr
	if revel.Config.BoolDefault("jobs.acceptproxyaddress", false) {
		if proxiedAddress := c.Request.GetHttpHeader("X-Forwarded-For"); proxiedAddress!="" {
			remoteAddress = proxiedAddress
		}
	}
	if !strings.HasPrefix(remoteAddress, "127.0.0.1") &&
		!strings.HasPrefix(remoteAddress, "::1") &&
		!strings.HasPrefix(remoteAddress, "[::1]") {
		return c.Forbidden("%s is not local", remoteAddress)
	}
	entries := jobs.MainCron.Entries()
	return c.Render(entries)
}

func init() {
	revel.TemplateFuncs["castjob"] = func(job cron.Job) *jobs.Job {
		return job.(*jobs.Job)
	}
}
