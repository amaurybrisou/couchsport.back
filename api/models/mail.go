package models

import (
	"bytes"
	"html/template"
)

//Mail models definition
type Mail struct {
	From         string
	To           []string
	subject      string
	mime         string
	body         string
	template     string
	templateVars map[string]string
}

const (
	mime = "MIME-version: 1.0;\r\nContent-Type: text/html; charset=UTF-8\r\n"
)

//NewMail create a new mail
func NewMail(from string, to []string, subject string, template string, templateVars map[string]string) *Mail {
	return &Mail{
		From:         from,
		To:           to,
		subject:      subject,
		template:     template,
		mime:         mime,
		templateVars: templateVars,
	}
}

//GetBody returns the formatted body
func (me *Mail) GetBody() (string, error) {
	err := me.parseTemplate(me.template, me.templateVars)
	if err != nil {
		return "", err
	}
	return "From: " + me.From + "\r\nTo: " + me.To[0] + "\r\nSubject: " + me.subject + "\r\n" + mime + "\r\n" + me.body, nil
}

func (me *Mail) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}

	me.body = buffer.String()

	return nil
}
