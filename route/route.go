package route

import (
	"github.com/calvinchengx/gin-go-pg/mail"
	mw "github.com/calvinchengx/gin-go-pg/middleware"
	"github.com/calvinchengx/gin-go-pg/mobile"
	"github.com/calvinchengx/gin-go-pg/repository"
	"github.com/calvinchengx/gin-go-pg/repository/account"
	"github.com/calvinchengx/gin-go-pg/repository/auth"
	"github.com/calvinchengx/gin-go-pg/repository/user"
	"github.com/calvinchengx/gin-go-pg/service"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

// NewServices creates a new router services
func NewServices(DB *pg.DB, Log *zap.Logger, JWT *mw.JWT, Mail mail.Service, Mobile mobile.Service, R *gin.Engine) *Services {
	return &Services{DB, Log, JWT, Mail, Mobile, R}
}

// Services lets us bind specific services when setting up routes
type Services struct {
	DB     *pg.DB
	Log    *zap.Logger
	JWT    *mw.JWT
	Mail   mail.Service
	Mobile mobile.Service
	R      *gin.Engine
}

// SetupV1Routes instances various repos and services and sets up the routers
func (s *Services) SetupV1Routes() {
	// database logic
	userRepo := repository.NewUserRepo(s.DB, s.Log)
	accountRepo := repository.NewAccountRepo(s.DB, s.Log)
	rbac := repository.NewRBACService(userRepo)

	// service logic
	authService := auth.NewAuthService(userRepo, accountRepo, s.JWT, s.Mail, s.Mobile)
	accountService := account.NewAccountService(userRepo, accountRepo, rbac)
	userService := user.NewUserService(userRepo, authService, rbac)

	// no prefix, no jwt
	service.AuthRouter(authService, s.R)

	// prefixed with /v1 and protected by jwt
	v1Router := s.R.Group("/v1")
	v1Router.Use(s.JWT.MWFunc())
	service.AccountRouter(accountService, v1Router)
	service.UserRouter(userService, v1Router)
}
