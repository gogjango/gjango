package controller

import (
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/gin-gonic/gin"
)

// AccountService represents the account application service
type AccountService struct {
	accountRepo model.AccountRepo
	userRepo    model.UserRepo
	rbac        model.RBACService
}

// NewAccountService creates a new account application service
func NewAccountService(userRepo model.UserRepo, accountRepo model.AccountRepo, rbac model.RBACService) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		userRepo:    userRepo,
		rbac:        rbac,
	}
}

// Create creates a new user account
func (s *AccountService) Create(c *gin.Context, u *model.User) error {
	// TODO: implement RBAC
	return s.accountRepo.Create(c, u)
}

// ChangePassword changes user's password
func (s *AccountService) ChangePassword(c *gin.Context, oldPass, newPass string, id int) error {
	// TODO: implement RBAC
	u, err := s.userRepo.View(c, id)
	if err != nil {
		return err
	}
	return s.accountRepo.ChangePassword(c, u)
}
