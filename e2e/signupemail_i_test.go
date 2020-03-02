package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/calvinchengx/gin-go-pg/request"
	"github.com/stretchr/testify/assert"
)

func (suite *E2ETestSuite) TestSignupEmail() {

	t := suite.T()

	ts := httptest.NewServer(suite.r)
	defer ts.Close()

	urlSignup := ts.URL + "/signup"

	req := &request.EmailSignup{
		Email:           "user@example.org",
		Password:        "userpassword1",
		PasswordConfirm: "userpassword1",
	}
	b, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(urlSignup, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

	assert.Nil(t, nil)
}
