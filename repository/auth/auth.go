package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
	"golang.org/x/crypto/bcrypt"
)

// NewAuthService creates new auth service
func NewAuthService(userRepo model.UserRepo, accountRepo model.AccountRepo, jwt JWT) *Service {
	return &Service{
		userRepo:    userRepo,
		accountRepo: accountRepo,
		jwt:         jwt,
	}
}

// Service represents the auth application service
type Service struct {
	userRepo    model.UserRepo
	accountRepo model.AccountRepo
	jwt         JWT
}

// JWT represents jwt interface
type JWT interface {
	GenerateToken(*model.User) (string, string, error)
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

// Signup returns any error from creating a new user in our database
func (s *Service) Signup(c *gin.Context) error {
	username := c.GetString("username")
	user, err := s.userRepo.FindByUsername(c, username)
	if err == nil {
		// no user will be created since it already exists
		return errors.New("user exists")
	}
	return s.accountRepo.Create(c, user)
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
