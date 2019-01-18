package stores

import (
	"bytes"
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/smtp"
)

type request struct {
	from    string
	to      []string
	subject string
	mime    string
	body    string
}

type mailStore struct {
	Email, Password, Server string
	Port                    int
	r                       *request
}

const (
	mime = "MIME-version: 1.0;\r\nContent-Type: text/html; charset=UTF-8\r\n"
)

func (me *mailStore) newRequest(to []string, subject string) {
	r := &request{
		from:    me.Email,
		to:      to,
		subject: subject,
		mime:    mime,
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

func (me *mailStore) sendMailTLS() error {
	message := "From: " + me.r.from + "\r\nTo: " + me.r.to[0] + "\r\nSubject: " + me.r.subject + "\r\n" + mime + "\r\n" + me.r.body
	smtpServer := fmt.Sprintf("%s:%d", me.Server, me.Port)

	auth := smtp.PlainAuth("", me.Email, me.Password, me.Server)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         me.Server,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", smtpServer, tlsconfig)
	if err != nil {
		log.Error(err)
	}

	c, err := smtp.NewClient(conn, me.Server)
	if err != nil {
		log.Error(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Error(err)
	}

	// To && From
	if err = c.Mail(me.r.from); err != nil {
		log.Error(err)
	}

	if err = c.Rcpt(me.r.to[0]); err != nil {
		log.Error(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Error(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Error(err)
	}

	err = w.Close()
	if err != nil {
		log.Error(err)
	}

	c.Quit()

	log.Printf("email sent to %s", me.r.to[0])

	me.r = nil

	return nil
}

func (me *mailStore) send(templateName string, items interface{}) error {
	err := me.parseTemplate(templateName, items)
	if err != nil {
		return err
	}
	if err := me.sendMailTLS(); err != nil {
		return err
	}

	return err
}

func (me *mailStore) AccountAutoCreated(email, password string) {
	me.newRequest([]string{email}, "Your account has been created")
	log.Printf("sending 'AccountAutoCreated' email to %s", email)
	if err := me.send("api/templates/mail/account_created.html", map[string]string{"email": email, "password": password}); err != nil {
		log.Error(err)
	}
}
