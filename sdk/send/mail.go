package send

import (
	"crypto/tls"
	"os"

	"github.com/wellopen/gout/model/sdk"
	"gopkg.in/gomail.v2"
)

func QQMailSendFile(config sdk.Config) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.From...)
	m.SetHeader("To", config.To...)
	m.SetHeader("Subject", config.Subject...)
	m.SetBody("text/html", config.Body)
	f, err := os.OpenFile(config.Path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	h := make(map[string][]string, 0)
	h["Content-Type"] = []string{`application/octet-stream; charset=utf-8; name="` + f.Name() + `"`} //要设置这个，否则中文会乱码
	fileSetting := gomail.SetHeader(h)
	m.Attach(f.Name(), fileSetting)
	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
