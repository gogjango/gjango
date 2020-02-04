package auth

import (
	"context"
	"net/http"

	"github.com/rs/xid"

	"github.com/calvinchengx/gin-go-pg/apperr"
	mw "github.com/calvinchengx/gin-go-pg/middleware"
	"github.com/calvinchengx/gin-go-pg/model"
	"golang.org/x/crypto/bcrypt"
)

// NewAuthService creates new auth service
func NewAuthService(userRepo model.UserRepo, jwt *mw.JWT) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

// AuthService represents the auth application service
type AuthService struct {
	userRepo model.UserRepo
	jwt      *mw.JWT
}

// Authenticate tries to authenticate the user provided by username and password
func (s *AuthService) Authenticate(c context.Context, username, password string) (*model.AuthToken, error) {
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
func (s *AuthService) Refresh(c context.Context, token string) (*model.RefreshToken, error) {
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

// HashPassword hashes the password using bcrypt
func HashPassword(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

// HashMatchesPassword matches hash with password. Returns true if hash and password match.
func HashMatchesPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
