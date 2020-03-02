package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/request"
	"golang.org/x/crypto/bcrypt"
)

// NewAuthService creates new auth service
func NewAuthService(userRepo model.UserRepo, accountRepo model.AccountRepo, jwt JWT, mail Mail, mobile Mobile) *Service {
	return &Service{
		userRepo:    userRepo,
		accountRepo: accountRepo,
		jwt:         jwt,
		mail:        mail,
		mobile:      mobile,
	}
}

// Service represents the auth application service
type Service struct {
	userRepo    model.UserRepo
	accountRepo model.AccountRepo
	jwt         JWT
	mail        Mail
	mobile      Mobile
}

// JWT represents jwt interface
type JWT interface {
	GenerateToken(*model.User) (string, string, error)
}

// Mail represents mail interface
type Mail interface {
	SendVerificationEmail(string, *model.Verification) error
}

// Mobile represents mobile interface
type Mobile interface {
	GenerateSMSToken(string, string) error
	CheckCode(string, string, string) error
}

// Authenticate tries to authenticate the user provided by username and password
func (s *Service) Authenticate(c context.Context, email, password string) (*model.AuthToken, error) {
	u, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if !HashMatchesPassword(u.Password, password) {
		return nil, apperr.New(http.StatusNotFound, "Username or password does not exist")
	}
	if !u.Active {
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
func (s *Service) Refresh(c context.Context, token string) (*model.RefreshToken, error) {
	user, err := s.userRepo.FindByToken(token)
	if err != nil {
		return nil, err
	}
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

// VerifyMobile verifies the mobile verification code, i.e. (6-digit) code
func (s *Service) VerifyMobile(c context.Context, countryCode, mobile, code string) error {
	// send code to twilio
	err := s.mobile.CheckCode(countryCode, mobile, code)
	if err != nil {
		return err
	}
	user, err := s.userRepo.FindByMobile(countryCode, mobile)
	if err != nil {
		return err
	}
	// if it code is approved, make user active
	user.Active = true
	_, err = s.userRepo.Update(user)
	if err != nil {
		return err
	}
	return nil
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
	err = s.mail.SendVerificationEmail(e.Email, v)
	if err != nil {
		apperr.Response(c, err)
		return err
	}
	return nil
}

// SignupMobile returns any error from creating a new user in our database with a mobile number
func (s *Service) SignupMobile(c *gin.Context, m *request.MobileSignup) error {
	// find by countryCode and mobile
	u, err := s.userRepo.FindByMobile(m.CountryCode, m.Mobile)
	if err == nil { // user already exists
		return apperr.NewStatus(http.StatusConflict)
	}
	// create and verify
	err = s.accountRepo.CreateWithMobile(u)
	if err != nil {
		return err
	}
	// generate sms token
	err = s.mobile.GenerateSMSToken(m.CountryCode, m.Mobile)
	if err != nil {
		apperr.Response(c, err)
		return err
	}
	return nil
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

// HashMatchesPassword matches hash with password. Returns true if hash and password match.
func HashMatchesPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
