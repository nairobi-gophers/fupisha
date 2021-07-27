package provider

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/nairobi-gophers/fupisha/config"
	"github.com/pkg/errors"
	"github.com/vanng822/go-premailer/premailer"
	"gopkg.in/mail.v2"
	"jaytaylor.com/html2text"
)

// var (
// 	templates *template.Template
// )

type Mailer struct {
	client   *mail.Dialer
	template *template.Template
	from     Email
}

//NewEmailWithSMTP is a constructor function that initializes and returns a ready to use mailer object, an error interface otherwise.
func NewMailerWithSMTP(cfg *config.Config, tplDir string) (*Mailer, error) {
	//parse templates here, if err we fail early and return.
	tpl, err := parseTemplates(tplDir)
	if err != nil {
		return nil, err
	}

	dialer := mail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password)

	m := &Mailer{
		client:   dialer,
		from:     NewEmail(cfg.SMTP.FromName, cfg.SMTP.FromAddress),
		template: tpl,
	}

	m.client.StartTLSPolicy = mail.MandatoryStartTLS

	d, err := m.client.Dial()
	if err != nil {
		return nil, err
	}
	defer d.Close()

	return m, nil
}

func (m Mailer) send(email interface{}) error {

	msg := mail.NewMessage()
	if em, ok := email.(*message); ok {
		msg.SetHeader("From", em.from.Address)
		msg.SetHeader("To", em.to.Address)
		msg.SetHeader("Subject", em.subject)
		msg.SetBody("text/plain", em.text)
		msg.AddAlternative("text/html", em.html)
		return m.client.DialAndSend(msg)
	}

	return errors.New("unsupported type.expected *message to be type email")
}

//SendVerifyNotification sends the verify email notification to the user's email address.
func (m Mailer) SendVerifyNotification(address string, content VerifyEmailContent) error {
	msg := &message{
		from:     m.from,
		to:       NewEmail("", address),
		subject:  "Verify Account",
		template: "verify",
		content:  content,
	}

	if err := msg.parse(m.template); err != nil {
		return err
	}

	return m.send(msg)
}

func (m Mailer) SendWelcomeNotification(address string, content WelcomeEmailContent) error {
	msg := &message{
		from:     m.from,
		to:       NewEmail("", address),
		subject:  "Welcome to Fupisha!",
		template: "welcome",
		content:  content,
	}

	if err := msg.parse(m.template); err != nil {
		return err
	}

	return m.send(msg)
}

func parseTemplates(tplDir string) (*template.Template, error) {

	templates := template.New("").Funcs(fMap)

	if err := filepath.Walk(tplDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".html" {
			_, err = templates.ParseFiles(path)
			return err
		}
		return err
	}); err != nil {
		return nil, err
	}

	return templates, nil
}

var fMap = template.FuncMap{
	"formatAsDuration": formatAsDuration,
}

func formatAsDuration(t time.Time) string {
	dur := time.Until(t)
	hours := int(dur.Hours())
	mins := int(dur.Minutes())

	var v string
	if hours != 0 {
		v += strconv.Itoa(hours) + " hours and "
	}
	v += strconv.Itoa(mins) + " minutes"
	return v
}

//Email holds the recipients email address and name
type Email struct {
	Name    string
	Address string
}

//NewEmail returns an email address
func NewEmail(name, address string) Email {
	return Email{
		name,
		address,
	}
}

type message struct {
	from     Email
	to       Email
	subject  string
	template string
	content  interface{}
	html     string
	text     string
}

func (m *message) parse(templates *template.Template) error {
	buf := new(bytes.Buffer)

	if err := templates.ExecuteTemplate(buf, m.template, m.content); err != nil {
		return err
	}

	prem, err := premailer.NewPremailerFromString(buf.String(), premailer.NewOptions())
	if err != nil {
		return err
	}

	html, err := prem.Transform()
	if err != nil {
		return err
	}

	m.html = html
	text, err := html2text.FromString(buf.String(), html2text.Options{PrettyTables: true})
	if err != nil {
		return err
	}
	m.text = text
	return nil
}

//VerifyEmailContent provides the values to be displayed in the verify email template.
type VerifyEmailContent struct {
	// Email              string
	SiteURL            string
	SiteName           string
	VerificationExpiry time.Time
	VerificationURL    string
}

//WelcomeEmailContent provides the values to be displayed in the welcome email template.
type WelcomeEmailContent struct {
	SiteURL  string
	SiteName string
	LoginURL string
}
