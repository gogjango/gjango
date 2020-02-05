package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"github.com/calvinchengx/gin-go-pg/apperr"
	mw "github.com/calvinchengx/gin-go-pg/middleware"
	"github.com/calvinchengx/gin-go-pg/model"
	"golang.org/x/crypto/bcrypt"
)

// NewAuthService creates new auth service
func NewAuthService(userRepo model.UserRepo, jwt *mw.JWT) *Service {
	return &Service{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

// Service represents the auth application service
type Service struct {
	userRepo model.UserRepo
	jwt      *mw.JWT
}

// Authenticate tries to authenticate the user provided by username and password
func (s *Service) Authenticate(c context.Context, username, password string) (*model.AuthToken, error) {
	u, err := s.userRepo.FindByUsername(c, username)
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
	if err := s.userRepo.UpdateLogin(c, u); err != nil {
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
	user, err := s.userRepo.FindByToken(c, token)
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

// HashPassword hashes the password using bcrypt
func HashPassword(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

// HashMatchesPassword matches hash with password. Returns true if hash and password match.
func HashMatchesPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
