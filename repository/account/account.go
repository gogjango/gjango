package account

import (
	"net/http"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository/auth"
	"github.com/gin-gonic/gin"
)

// Service represents the account application service
type Service struct {
	accountRepo model.AccountRepo
	userRepo    model.UserRepo
	rbac        model.RBACService
}

// NewAccountService creates a new account application service
func NewAccountService(userRepo model.UserRepo, accountRepo model.AccountRepo, rbac model.RBACService) *Service {
	return &Service{
		accountRepo: accountRepo,
		userRepo:    userRepo,
		rbac:        rbac,
	}
}

// Create creates a new user account
func (s *Service) Create(c *gin.Context, u *model.User) error {
	if !s.rbac.AccountCreate(c, u.RoleID, u.CompanyID, u.LocationID) {
		return apperr.Forbidden
	}
	u.Password = auth.HashPassword(u.Password)
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
	if !auth.HashMatchesPassword(u.Password, oldPass) {
		return apperr.New(http.StatusBadGateway, "old password is not correct")
	}
	u.Password = auth.HashPassword(newPass)
	return s.accountRepo.ChangePassword(c, u)
}
