package stores

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/smtp"
)

type request struct {
	from    string
	to      []string
	subject string
	body    string
}

type mailStore struct {
	Email, Password, Server string
	Port                    int
	r                       *request
}

const (
	mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func (me *mailStore) newRequest(to []string, subject string) {
	r := &request{
		to:      to,
		subject: subject,
	}
	me.r = r
}

func (me mailStore) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	me.r.body = buffer.String()
	return nil
}

func (me *mailStore) sendMail() error {
	body := "To: " + me.r.to[0] + "\r\nSubject: " + me.r.subject + "\r\n" + mime + "\r\n" + me.r.body
	SMTP := fmt.Sprintf("%s:%d", me.Server, me.Port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", me.Email, me.Password, me.Server), me.Email, me.r.to, []byte(body)); err != nil {
		return err
	}

	me.r = nil

	return nil
}

func (me *mailStore) send(templateName string, items interface{}) error {
	err := me.parseTemplate(templateName, items)
	if err != nil {
		return err
	}
	if err := me.sendMail(); err != nil {
		return err
	}

	return err
}

func (me *mailStore) AccountAutoCreated(email, password string) {
	me.newRequest([]string{email}, "Your account has been created")
	if err := me.send("api/templates/mail/account_created.html", map[string]string{"email": email, "password": password}); err != nil {
		log.Error(err)
	}
}
