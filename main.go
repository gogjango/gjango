package main

import (
	"strconv"

	"github.com/calvinchengx/gin-go-pg/config"
	"github.com/calvinchengx/gin-go-pg/mail"
	mw "github.com/calvinchengx/gin-go-pg/middleware"
	"github.com/calvinchengx/gin-go-pg/mobile"
	"github.com/calvinchengx/gin-go-pg/repository"
	"github.com/calvinchengx/gin-go-pg/repository/account"
	"github.com/calvinchengx/gin-go-pg/repository/auth"
	"github.com/calvinchengx/gin-go-pg/repository/user"
	"github.com/calvinchengx/gin-go-pg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

func main() {
	r := gin.Default()

	// middleware
	mw.Add(r, cors.Default())

	// load configuration
	c, _ := config.Load("dev")
	jwt := mw.NewJWT(c.JWT)
	m := mail.NewMail(config.GetMailConfig(), config.GetSiteConfig())
	mobile := mobile.NewMobile(config.GetTwilioConfig())
	db := config.GetConnection()
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	// setup routes
	rs := &RouteServices{db, log, jwt, m, mobile, r}
	rs.setupV1Routes()

	// run with port from config
	port := ":" + strconv.Itoa(c.Server.Port)
	r.Run(port)
}

// RouteServices lets us bind specific services when setting up routes
type RouteServices struct {
	db     *pg.DB
	log    *zap.Logger
	jwt    *mw.JWT
	m      *mail.Mail
	mobile *mobile.Mobile
	r      *gin.Engine
}

func (s *RouteServices) setupV1Routes() {
	// database logic
	userRepo := repository.NewUserRepo(s.db, s.log)
	accountRepo := repository.NewAccountRepo(s.db, s.log)
	rbac := repository.NewRBACService(userRepo)

	// service logic
	authService := auth.NewAuthService(userRepo, accountRepo, s.jwt, s.m, s.mobile)
	accountService := account.NewAccountService(userRepo, accountRepo, rbac)
	userService := user.NewUserService(userRepo, authService, rbac)

	// no prefix, no jwt
	service.AuthRouter(authService, s.r)

	// prefixed with /v1 and protected by jwt
	v1Router := s.r.Group("/v1")
	v1Router.Use(s.jwt.MWFunc())
	service.AccountRouter(accountService, v1Router)
	service.UserRouter(userService, v1Router)
}
