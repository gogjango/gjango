package mail

import (
	"os"

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
