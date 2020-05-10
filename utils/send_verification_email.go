package utils

import (
	"errors"

	"github.com/bristolxyz/bristol.xyz/clients"
	"github.com/bristolxyz/bristol.xyz/models"
	"github.com/getsentry/sentry-go"
)

// SendVerificationEmail allows you to send a verification e-mail to a user.
func SendVerificationEmail(u *models.User) error {
	// Check the user is not nil.
	if u == nil {
		return errors.New("user should be set")
	}

	// Get the verification ID from the models function.
	VerificationCode, err := models.NewVerificationCode(u)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	// Send the email.
	URL := "https://bristol.xyz/verify/" + VerificationCode.ID
	Email := "Thanks for signing up to bristol.xyz! The following URL will let you verify your user to recieve notifications: <a href=\"" + URL + "\">" + URL + "</a>"
	return clients.SendEmail(u.Email, "[bristol.xyz] Verification", Email)
}
