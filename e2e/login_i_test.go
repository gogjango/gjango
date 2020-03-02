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

func (suite *E2ETestSuite) TestLogin() {

	t := suite.T()
	fmt.Println("TestLogin")

	ts := httptest.NewServer(suite.r)
	defer ts.Close()

	url := ts.URL + "/login"

	req := &request.Credentials{
		Email:    "superuser@example.org",
		Password: "testpassword",
	}
	b, err := json.Marshal(req)
	assert.Nil(t, err)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
