package email_utils

import (
	"douyin-backend/app/global/variable"
	"fmt"

	"gopkg.in/gomail.v2"
)

// SendVerificationCode 发送验证码邮件
func SendVerificationCode(email, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>",
		variable.ConfigYml.GetString("Email.FromName"),
		variable.ConfigYml.GetString("Email.From")))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "验证码")
	m.SetBody("text/html", fmt.Sprintf("您的验证码是：<b>%s</b>，%d分钟内有效。",
		code,
		variable.ConfigYml.GetInt("Email.VerificationCodeExpire")/60))

	d := gomail.NewDialer(
		variable.ConfigYml.GetString("Email.Host"),
		variable.ConfigYml.GetInt("Email.Port"),
		variable.ConfigYml.GetString("Email.Username"),
		variable.ConfigYml.GetString("Email.Password"),
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}
