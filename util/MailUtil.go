package util

import (
	"gopkg.in/gomail.v2"
)

// Mail 邮件模型
type Mail struct {
	To      string `json:"email"` // 收件人
	Subject string // 主题
	Body    string // 内容
}

var password = "fstrpvwdziuibibd" // 邮箱smtp授权码

// SendMail 发送邮件
func (mail *Mail) SendMail() error {
	m := gomail.NewMessage()

	//发送人
	m.SetHeader("From", "325936216@qq.com")
	//接收人
	m.SetHeader("To", mail.To)
	//抄送人
	//m.SetAddressHeader("Cc", "xxx@qq.com", "xiaozhujiao")
	//主题
	m.SetHeader("Subject", mail.Subject)
	//内容
	m.SetBody("text/html", mail.Body)
	//附件
	//m.Attach("./myIpPic.png")

	//拿到token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer("smtp.qq.com", 587, "325936216@qq.com", password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
