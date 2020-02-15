package mail

import (
	"os"

	"github.com/calvinchengx/gin-go-pg/config"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/sendgrid/sendgrid-go"
	s "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// NewMail generates new Mail variable
func NewMail(mc *config.MailConfig, sc *config.SiteConfig) *Mail {
	return &Mail{
		ExternalURL: sc.ExternalURL,
		FromName:    mc.Name,
		FromEmail:   mc.Email,
	}
}

// Mail provides a mail service implementation
type Mail struct {
	ExternalURL string
	FromName    string
	FromEmail   string
}

// Send email with sendgrid
func (m *Mail) Send(subject string, toName string, toEmail string, content string) error {
	from := s.NewEmail(m.FromName, m.FromEmail)
	to := s.NewEmail(toName, toEmail)
	plainTextContent := content
	message := s.NewSingleEmail(from, subject, to, plainTextContent, content)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return err
	}
	return nil
}

// SendWithDefaults assumes some defaults for sending out email with sendgrid
func (m *Mail) SendWithDefaults(subject, toEmail, content string) error {
	err := m.Send(subject, toEmail, toEmail, content)
	if err != nil {
		return err
	}
	return nil
}

// SendVerificationEmail assumes defaults and generates a verification email
func (m *Mail) SendVerificationEmail(toEmail string, v *model.Verification) error {
	url := m.ExternalURL + "/verification/" + v.Token
	content := "Click on this to verify your account " + url
	err := m.SendWithDefaults("Verification Email", toEmail, content)
	if err != nil {
		return err
	}
	return nil
}
