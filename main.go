package main

import (
	"strconv"

	"github.com/calvinchengx/gin-go-pg/config"
	mw "github.com/calvinchengx/gin-go-pg/middleware"
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
	db := config.GetConnection()
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	// setup routes
	setupV1Routes(db, log, jwt, r)

	// run with port from config
	port := ":" + strconv.Itoa(c.Server.Port)
	r.Run(port)
}

func setupV1Routes(db *pg.DB, log *zap.Logger, jwt *mw.JWT, r *gin.Engine) {
	// database logic
	userRepo := repository.NewUserRepo(db, log)
	accountRepo := repository.NewAccountRepo(db, log)
	rbac := repository.NewRBACService(userRepo)

	// service logic
	authService := auth.NewAuthService(userRepo, accountRepo, jwt)
	accountService := account.NewAccountService(userRepo, accountRepo, rbac)
	userService := user.NewUserService(userRepo, authService, rbac)

	// no prefix, no jwt
	service.AuthRouter(authService, r)

	// prefixed with /v1 and protected by jwt
	v1Router := r.Group("/v1")
	v1Router.Use(jwt.MWFunc())
	service.AccountRouter(accountService, v1Router)
	service.UserRouter(userService, v1Router)
}
