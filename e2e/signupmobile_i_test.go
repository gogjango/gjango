package e2e_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/calvinchengx/gin-go-pg/request"
	"github.com/stretchr/testify/assert"
)

func (suite *E2ETestSuite) TestSignupMobile() {
	t := suite.T()
	ts := httptest.NewServer(suite.r)
	defer ts.Close()

	urlSignupMobile := ts.URL + "/signup/m"

	req := &request.MobileSignup{
		CountryCode: "+65",
		Mobile:      "91919191",
	}
	b, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(urlSignupMobile, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Nil(t, err)
}
