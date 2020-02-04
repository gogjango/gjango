package controller

import (
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository"
	"github.com/gin-gonic/gin"
)

// AccountService represents the account application service
type AccountService struct {
	accountRepo *repository.AccountRepo
}

// NewAccountService creates a new account application service
func NewAccountService(accountRepo *repository.AccountRepo) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

// Create creates a new user account
func (s *AccountService) Create(c *gin.Context, u *model.User) error {
	return s.accountRepo.Create(c, u)
}
