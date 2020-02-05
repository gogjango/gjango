package user

import (
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository/platform/query"
	"github.com/gin-gonic/gin"
)

// NewUserService create a new user application service
func NewUserService(userRepo model.UserRepo, auth model.AuthService, rbac model.RBACService) *Service {
	return &Service{
		userRepo: userRepo,
		auth:     auth,
		rbac:     rbac,
	}
}

// Service represents the user application service
type Service struct {
	userRepo model.UserRepo
	auth     model.AuthService
	rbac     model.RBACService
}

// List returns list of users
func (s *Service) List(c *gin.Context, p *model.Pagination) ([]model.User, error) {
	u := s.auth.User(c)
	q, err := query.List(u)
	if err != nil {
		return nil, err
	}
	return s.userRepo.List(c, q, p)
}
