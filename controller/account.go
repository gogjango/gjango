package controller

import (
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository"
	"github.com/gin-gonic/gin"
)

// AccountService represents the account application service
type AccountService struct {
	accountRepo *repository.AccountRepo
	userRepo    *repository.UserRepo
}

// NewAccountService creates a new account application service
func NewAccountService(accountRepo *repository.AccountRepo) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
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
