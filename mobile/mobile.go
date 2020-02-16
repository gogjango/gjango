package mobile

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/calvinchengx/gin-go-pg/config"
)

// NewMobile creates a new mobile service implementation
func NewMobile(config *config.TwilioConfig) *Mobile {
	return &Mobile{config}
}

// Mobile provides a mobile service implementation
type Mobile struct {
	config *config.TwilioConfig
}

// GenerateSMSToken sends an sms token to the mobile numer
// func (m *Mobile) GenerateSMSToken(countryCode, mobile string) error
// m.GenerateSMSToken("+65", "90901299")
func (m *Mobile) GenerateSMSToken(countryCode, mobile string) error {
	apiURL := m.getTwilioVerifyURL()
	fmt.Println(apiURL)

	data := url.Values{}
	data.Set("To", countryCode+mobile)
	data.Set("Channel", "sms")

	u, _ := url.ParseRequestURI(apiURL)
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.SetBasicAuth(m.config.Account, m.config.Token)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	return err
}

func (m *Mobile) getTwilioVerifyURL() string {
	return "https://verify.twilio.com/v2/Services/" + m.config.Verify + "/Verifications"
}
