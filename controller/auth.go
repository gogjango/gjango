package controller

import (
	"context"
	"net/http"

	"github.com/rs/xid"

	"github.com/calvinchengx/gin-go-pg/apperr"
	mw "github.com/calvinchengx/gin-go-pg/middleware"
	"github.com/calvinchengx/gin-go-pg/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo model.UserRepo
	jwt      mw.JWT
}

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
	// TODO: finish this
	// if err := s.userRepo.UpdateLogin
	return &model.AuthToken{
		Token:        token,
		Expires:      expire,
		RefreshToken: u.Token,
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
