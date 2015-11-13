package lea

import (
	"github.com/revel/revel"
	"net/smtp"
	"strings"
)

// 发送邮件
var host = "smtp.ym.163.com"
var port = "25"
var username = "noreply@leanote.com"
var password = "---"

func InitEmail() {
	config := revel.Config
	host, _ = config.String("email.host")
	port, _ = config.String("email.port")
	username, _ = config.String("email.username")
	password, _ = config.String("email.password")
}

var bodyTpl = `
	<html>
	<body>
		<div style="width: 600px; margin:auto; border-radius:5px; border: 1px solid #ccc; padding: 20px;">
			<div>
				<div>
					<div style="float:left; height: 40px;">
						<a href="http://leanote.com" style="font-size: 24px">leanote</a>
					</div>
					<div style="float:left; height:40px; line-height:40px;">
						&nbsp;&nbsp;| &nbsp;<span style="font-size:14px">$title</span>
					</div>
					<div style="clear:both"></div>
				</div>
			</div>
			<hr style="border:none;border-top: 1px solid #ccc"/>
			<div style="margin-top: 20px; font-size: 14px;">
				$body
			</div>

			<div id="leanoteFooter" style="margin-top: 30px; border-top: 1px solid #ccc">
				<style>
					#leanoteFooter {
						color: #666;
						font-size: 12px;
					}
					#leanoteFooter a {
						color: #666;
						font-size: 12px;
					}
				</style>
				<a href="http://leanote.com">leanote</a>, your own cloud note!
			</div>
		</div>
	</body>
	</html>
`

func SendEmailOld(to, subject, body string) bool {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", username, password, hp[0])

	var content_type string

	mailtype := "html"
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	//body = strings.Replace(bodyTpl, "$body", body, 1)
	//body = strings.Replace(body, "$title", title, 1)

	msg := []byte("To: " + to + "\r\nFrom: " + username + "<" + username + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host+":"+port, auth, username, send_to, msg)

	if err != nil {
		Log(err)
		return false
	}
	return true
}

func SendToLeanoteOld(subject, title, body string) {
	to := "leanote@leanote.com"
	SendEmailOld(to, subject, body)
}
