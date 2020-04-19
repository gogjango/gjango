package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogjango/gjango/apperr"
	"github.com/gogjango/gjango/model"
	"github.com/gogjango/gjango/secret"
)

// Service represents the account application service
type Service struct {
	accountRepo model.AccountRepo
	userRepo    model.UserRepo
	rbac        model.RBACService
	secret      secret.Service
}

// NewAccountService creates a new account application service
func NewAccountService(userRepo model.UserRepo, accountRepo model.AccountRepo, rbac model.RBACService, secret secret.Service) *Service {
	return &Service{
		accountRepo: accountRepo,
		userRepo:    userRepo,
		rbac:        rbac,
		secret:      secret,
	}
}

// Create creates a new user account
func (s *Service) Create(c *gin.Context, u *model.User) error {
	if !s.rbac.AccountCreate(c, u.RoleID, u.CompanyID, u.LocationID) {
		return apperr.Forbidden
	}
	u.Password = s.secret.HashPassword(u.Password)
	u, err := s.accountRepo.Create(u)
	return err
}

// ChangePassword changes user's password
func (s *Service) ChangePassword(c *gin.Context, oldPass, newPass string, id int) error {
	if !s.rbac.EnforceUser(c, id) {
		return apperr.Forbidden
	}
	u, err := s.userRepo.View(id)
	if err != nil {
		return err
	}
	if !s.secret.HashMatchesPassword(u.Password, oldPass) {
		return apperr.New(http.StatusBadGateway, "old password is not correct")
	}
	u.Password = s.secret.HashPassword(newPass)
	return s.accountRepo.ChangePassword(u)
}
