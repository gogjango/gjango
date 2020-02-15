package mail

import (
	"os"

	"github.com/calvinchengx/gin-go-pg/config"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/sendgrid/sendgrid-go"
	m "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Send email with sendgrid
func Send(subject string, fromName string, fromEmail string, toName string, toEmail string, content string) error {
	from := m.NewEmail(fromName, fromEmail)
	to := m.NewEmail(toName, toEmail)
	plainTextContent := content
	message := m.NewSingleEmail(from, subject, to, plainTextContent, content)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return err
	}
	return nil
}

// SendWithDefaults assumes some defaults for sending out email with sendgrid
func SendWithDefaults(subject, toEmail, content string) error {
	c := config.GetMailConfig()
	err := Send(subject, c.Name, c.Email, toEmail, toEmail, content)
	if err != nil {
		return err
	}
	return nil
}

// SendVerificationEmail assumes defaults and generates a verification email
func SendVerificationEmail(toEmail string, v *model.Verification) error {
	c := config.GetSiteConfig()
	url := c.ExternalURL + "/verification/" + v.Token
	content := "Click on this to verify your account " + url
	err := SendWithDefaults("Verification Email", toEmail, content)
	if err != nil {
		return err
	}
	return nil
}
