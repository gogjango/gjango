package mail

import (
	"os"

	"github.com/calvinchengx/gin-go-pg/config"
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

// DefaultSend assumes some defaults for sending out email with sendgrid
func DefaultSend(subject, toEmail, content string) {
	c := config.GetMailConfig()
	Send(subject, c.Name, c.Email, toEmail, toEmail, content)
}
