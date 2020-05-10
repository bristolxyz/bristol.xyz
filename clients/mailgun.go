package clients

import (
	"context"

	"github.com/bristolxyz/bristol.xyz/env"
	"github.com/mailgun/mailgun-go/v4"
)

// MailgunClient defines the client which is used for mailgun.
var MailgunClient = mailgun.NewMailgun(env.RequiredEnvs("MAILGUN_DOMAIN")["MAILGUN_DOMAIN"], env.RequiredEnvs("MAILGUN_KEY")["MAILGUN_KEY"])

// Sets Mailgun to the EU region.
func init() {
	MailgunClient.SetAPIBase(mailgun.APIBaseEU)
}

// SendEmail is used to send a email.
func SendEmail(ToAddress, Subject, HTML string) error {
	m := MailgunClient.NewMessage(env.RequiredEnvs("FROM_ADDRESS")["FROM_ADDRESS"], Subject, "", ToAddress)
	m.SetHtml(HTML)
	_, _, err := MailgunClient.Send(context.TODO(), m)
	return err
}
