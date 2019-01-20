package stores

import (
	"crypto/tls"
	"fmt"
	"github.com/goland-amaurybrisou/couchsport/api/models"
	log "github.com/sirupsen/logrus"
	"net/smtp"
)

type mailStore struct {
	Email, Password, Server string
	Port                    int
	mail                    *models.Mail
}

func (me *mailStore) sendMail() error {
	body, err := me.mail.GetBody()
	if err != nil {
		return err
	}

	smtpServer := fmt.Sprintf("%s:%d", me.Server, me.Port)

	if err := smtp.SendMail(smtpServer, smtp.PlainAuth("", me.Email, me.Password, me.Server), me.Email, me.mail.To, []byte(body)); err != nil {
		return err
	}

	me.mail = nil

	return nil
}

func (me *mailStore) sendMailTLS() error {
	body, err := me.mail.GetBody()
	if err != nil {
		return err
	}
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
	if err = c.Mail(me.mail.From); err != nil {
		log.Error(err)
	}

	if err = c.Rcpt(me.mail.To[0]); err != nil {
		log.Error(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Error(err)
	}

	_, err = w.Write([]byte(body))
	if err != nil {
		log.Error(err)
	}

	err = w.Close()
	if err != nil {
		log.Error(err)
	}

	c.Quit()

	log.Printf("email sent to %s", me.mail.To[0])

	me.mail = nil

	return nil
}

func (me *mailStore) send(tls bool) error {

	if tls {
		if err := me.sendMailTLS(); err != nil {
			return err
		}
	} else {
		if err := me.sendMail(); err != nil {
			return err
		}
	}

	return nil
}

func (me *mailStore) AccountAutoCreated(email, password string) {
	me.mail = models.NewMail(
		me.Email,
		[]string{email},
		"Your account has been created", "api/templates/mail/account_created.html",
		map[string]string{"email": email, "password": password},
	)
	log.Printf("sending 'AccountAutoCreated' email to %s", email)
	if err := me.send(true); err != nil {
		log.Error(err)
	}
}
