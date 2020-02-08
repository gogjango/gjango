package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/calvinchengx/gin-go-pg/mock"
	"github.com/calvinchengx/gin-go-pg/mock/mockdb"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository/account"
	"github.com/calvinchengx/gin-go-pg/service"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		name        string
		req         string
		wantStatus  int
		wantResp    *model.User
		accountRepo *mockdb.Account
		rbac        *mock.RBAC
	}{
		{
			name:       "Invalid request",
			req:        `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter1234","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":3}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on userSvc",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":2}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c *gin.Context, roleID, companyID, locationID int) bool {
					return false
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":2}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c *gin.Context, roleID, companyID, locationID int) bool {
					return true
				},
			},
			accountRepo: &mockdb.Account{
				CreateFn: func(c context.Context, usr *model.User) error {
					usr.ID = 1
					usr.CreatedAt = mock.TestTime(2018)
					usr.UpdatedAt = mock.TestTime(2018)
					return nil
				},
			},
			wantResp: &model.User{
				Base: model.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2018),
					UpdatedAt: mock.TestTime(2018),
				},
				FirstName:  "John",
				LastName:   "Doe",
				Username:   "juzernejm",
				Email:      "johndoe@gmail.com",
				CompanyID:  1,
				LocationID: 2,
			},
			wantStatus: http.StatusOK,
		},
	}
	gin.SetMode(gin.TestMode)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rg := r.Group("/v1")
			accountService := account.NewAccountService(nil, tt.accountRepo, tt.rbac)
			service.AccountRouter(accountService, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/users"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(model.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
