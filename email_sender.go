package goemail

import (
	"bytes"
	"html/template"
	"net/smtp"
)

// email model
type email struct {
	authEmail    string
	authPassword string
	//auth host, eg: smtp.gmail.com
	authHost string
	//address should include smtp provider port eg: "smtp.gmail.com:587" google smtp host
	address string
	from    string
	to      []string
	subject string
	body    string
}

// New function, for initialize email model
func New(to []string, address, from, subject, body, authEmail, authPassword, authHost string) emailSender {
	return &email{
		authEmail:    authEmail,
		authPassword: authPassword,
		authHost:     authHost,
		address:      address,
		to:           to,
		subject:      subject,
		body:         body,
	}
}

// Send function, for send email
func (e *email) send() error {
	//setup auth
	auth := smtp.PlainAuth("", e.authEmail, e.authPassword, e.authHost)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + e.subject + "!\n"
	msg := []byte(subject + mime + "\n" + e.body)

	if err := smtp.SendMail(e.address, auth, e.from, e.to, msg); err != nil {
		return err
	}
	return nil
}

// SetTemplate function for set and parse template to email body
func (e *email) setTemplate(templateFile string, data interface{}) error {
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	e.body = buf.String()
	return nil
}

// Execute function for execute EmailSender implementation
func Execute(u emailSender, fileName string, data interface{}) error {
	err := u.setTemplate(fileName, data)
	if err != nil {
		return err
	}

	err = u.send()
	if err != nil {
		return err
	}

	return nil
}
