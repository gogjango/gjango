package main

import (
	"fmt"

	"github.com/calvinchengx/gin-go-pg/config"
	mw "github.com/calvinchengx/gin-go-pg/middleware"
	"github.com/calvinchengx/gin-go-pg/repository"
	"github.com/calvinchengx/gin-go-pg/repository/account"
	"github.com/calvinchengx/gin-go-pg/repository/auth"
	"github.com/calvinchengx/gin-go-pg/repository/user"
	"github.com/calvinchengx/gin-go-pg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	r := gin.Default()
	mw.Add(r, cors.Default())

	c, _ := config.Load("dev")
	jwt := mw.NewJWT(c.JWT)

	db := config.GetConnection()
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	userRepo := repository.NewUserRepo(db, log)
	accountRepo := repository.NewAccountRepo(db, log)
	rbac := repository.NewRBACService(userRepo)

	authService := auth.NewAuthService(userRepo, jwt)
	accountService := account.NewAccountService(userRepo, accountRepo, rbac)
	userService := user.NewUserService(userRepo, authService, rbac)
	fmt.Println(userService)

	service.AuthRouter(authService, r)

	v1Router := r.Group("/v1")
	v1Router.Use(jwt.MWFunc())

	service.AccountRouter(accountService, v1Router)
	service.UserRouter(userService, v1Router)

	r.Run()
}
