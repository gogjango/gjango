package mail

import "github.com/calvinchengx/gin-go-pg/model"

// Service is the interface to access our Mail
type Service interface {
	Send(subject string, toName string, toEmail string, content string) error
	SendWithDefaults(subject, toEmail, content string) error
	SendVerificationEmail(toEmail string, v *model.Verification) error
}
