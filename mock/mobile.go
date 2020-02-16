package mock

// Mobile mock
type Mobile struct {
	GenerateSMSTokenFn func(string, string) error
}

// GenerateSMSToken mock
func (m *Mobile) GenerateSMSToken(countryCode, mobile string) error {
	return m.GenerateSMSTokenFn(countryCode, mobile)
}
