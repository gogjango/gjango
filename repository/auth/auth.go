package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"github.com/gogjango/gjango/apperr"
	"github.com/gogjango/gjango/mail"
	"github.com/gogjango/gjango/mobile"
	"github.com/gogjango/gjango/model"
	"github.com/gogjango/gjango/request"
	"github.com/gogjango/gjango/secret"
)

// NewAuthService creates new auth service
func NewAuthService(userRepo model.UserRepo, accountRepo model.AccountRepo, jwt JWT, m mail.Service, mob mobile.Service) *Service {
	return &Service{userRepo, accountRepo, jwt, m, mob}
}

// Service represents the auth application service
type Service struct {
	userRepo    model.UserRepo
	accountRepo model.AccountRepo
	jwt         JWT
	m           mail.Service
	mob         mobile.Service
}

// JWT represents jwt interface
type JWT interface {
	GenerateToken(*model.User) (string, string, error)
}

// Authenticate tries to authenticate the user provided by username and password
func (s *Service) Authenticate(c context.Context, email, password string) (*model.AuthToken, error) {
	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, apperr.Unauthorized
	}
	if !secret.New().HashMatchesPassword(u.Password, password) {
		return nil, apperr.Unauthorized
	}
	// user must be active and verified. Active is enabled/disabled by superadmin user. Verified depends on user verifying via /verification/:token or /mobile/verify
	if !u.Active || !u.Verified {
		return nil, apperr.Unauthorized
	}
	token, expire, err := s.jwt.GenerateToken(u)
	if err != nil {
		return nil, apperr.Unauthorized
	}
	u.UpdateLastLogin()
	u.Token = xid.New().String()
	if err := s.userRepo.UpdateLogin(u); err != nil {
		return nil, err
	}
	return &model.AuthToken{
		Token:        token,
		Expires:      expire,
		RefreshToken: u.Token,
	}, nil
}

// Refresh refreshes jwt token and puts new claims inside
func (s *Service) Refresh(c context.Context, refreshToken string) (*model.RefreshToken, error) {
	user, err := s.userRepo.FindByToken(refreshToken)
	if err != nil {
		return nil, err
	}
	// this is our re-generated JWT
	token, expire, err := s.jwt.GenerateToken(user)
	if err != nil {
		return nil, apperr.Generic
	}
	return &model.RefreshToken{
		Token:   token,
		Expires: expire,
	}, nil
}

// Verify verifies the (verification) token and deletes it
func (s *Service) Verify(c context.Context, token string) error {
	v, err := s.accountRepo.FindVerificationToken(token)
	if err != nil {
		return err
	}
	err = s.accountRepo.DeleteVerificationToken(v)
	if err != nil {
		return err
	}
	return nil
}

// MobileVerify verifies the mobile verification code, i.e. (6-digit) code
func (s *Service) MobileVerify(c context.Context, countryCode, mobile, code string, signup bool) (*model.AuthToken, error) {
	// send code to twilio
	err := s.mob.CheckCode(countryCode, mobile, code)
	if err != nil {
		return nil, err
	}
	u, err := s.userRepo.FindByMobile(countryCode, mobile)
	if err != nil {
		return nil, err
	}
	if signup { // signup case, make user verified and active
		u.Verified = true
		u.Active = true
	} else { // login case, update user's last_login attribute
		u.UpdateLastLogin()
	}
	u, err = s.userRepo.Update(u)
	if err != nil {
		return nil, err
	}

	// generate jwt and return
	token, expire, err := s.jwt.GenerateToken(u)
	if err != nil {
		return nil, apperr.Unauthorized
	}
	u.UpdateLastLogin()
	u.Token = xid.New().String()
	if err := s.userRepo.UpdateLogin(u); err != nil {
		return nil, err
	}
	return &model.AuthToken{
		Token:        token,
		Expires:      expire,
		RefreshToken: u.Token,
	}, nil
}

// User returns user data stored in jwt token
func (s *Service) User(c *gin.Context) *model.AuthUser {
	id := c.GetInt("id")
	companyID := c.GetInt("company_id")
	locationID := c.GetInt("location_id")
	user := c.GetString("username")
	email := c.GetString("email")
	role := c.MustGet("role").(int8)
	return &model.AuthUser{
		ID:         id,
		Username:   user,
		CompanyID:  companyID,
		LocationID: locationID,
		Email:      email,
		Role:       model.AccessRole(role),
	}
}

// Signup returns any error from creating a new user in our database
func (s *Service) Signup(c *gin.Context, e *request.EmailSignup) error {
	_, err := s.userRepo.FindByEmail(e.Email)
	if err == nil { // user already exists
		return apperr.NewStatus(http.StatusConflict)
	}
	v, err := s.accountRepo.CreateAndVerify(&model.User{Email: e.Email, Password: e.Password})
	if err != nil {
		return err
	}
	err = s.m.SendVerificationEmail(e.Email, v)
	if err != nil {
		apperr.Response(c, err)
		return err
	}
	return nil
}

// Mobile returns any error from creating a new user in our database with a mobile number
func (s *Service) Mobile(c *gin.Context, m *request.MobileSignup) error {
	// find by countryCode and mobile
	_, err := s.userRepo.FindByMobile(m.CountryCode, m.Mobile)
	if err == nil { // user already exists
		return apperr.New(http.StatusConflict, "User already exists.")
	}
	// create and verify
	user := &model.User{
		CountryCode: m.CountryCode,
		Mobile:      m.Mobile,
	}
	err = s.accountRepo.CreateWithMobile(user)
	if err != nil {
		return err
	}
	// generate sms token
	err = s.mob.GenerateSMSToken(m.CountryCode, m.Mobile)
	if err != nil {
		apperr.Response(c, err)
		return err
	}
	return nil
}
