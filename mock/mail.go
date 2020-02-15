package mock

import "github.com/calvinchengx/gin-go-pg/model"

// Mail mock
type Mail struct {
	SendVerificationEmailFn func(string, *model.Verification) error
}

// SendVerificationEmail mock
func (m *Mail) SendVerificationEmail(toEmail string, v *model.Verification) error {
	return m.SendVerificationEmailFn(toEmail, v)
}
