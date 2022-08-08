package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	m "example.com/m/v2/models"
)
//create template struct and function to create it
type  templateData struct {
	BillingAdminName string
	BillingAdminMail string
	CompanyName      string
	BillingStartDate time.Time
	Projects         []m.Project
}

func NewTemplateRequest(temp) *templateData {
	return &mailRequest{
		to:      to,
		subject: subject,
		body:    body,
	}
}


var auth smtp.Auth
// add secret to email server 
func SendEmailToBillingAdmins(billingAccount m.BillingAccount) {
	auth = smtp.PlainAuth("", "123456testest123456@gmail.com", "123test123", "smtp.gmail.com")

	for i := 0; i < len(billingAccount.BillingAdmins); i++ {
		//change struct 
		templateData := {
			BillingAdminName: billingAccount.BillingAdmins[i].Name,
			BillingAdminMail: billingAccount.BillingAdmins[i].Email,
			CompanyName:      billingAccount.Company.Name,
			BillingStartDate: billingAccount.BillingStartDate,
			Projects:         billingAccount.Projects,
		}
		r := NewMailRequest([]string{templateData.BillingAdminMail}, "Hello Junk!", "Hello, World!")
		fmt.Print(r.ParseTemplate("templates/mail.html", templateData))
		if err := r.ParseTemplate("templates/mail.html", templateData); err == nil {

			ok, err := r.SendEmail()
			fmt.Println(ok)
			fmt.Println(err)
		}

	}

}

//Request struct
type mailRequest struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewMailRequest(to []string, subject, body string) *mailRequest {
	return &mailRequest{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *mailRequest) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, "dhanush@geektrust.in", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mailRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
