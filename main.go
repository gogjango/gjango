package main

import (
	"github.com/calvinchengx/gin-go-pg/config"
	"github.com/calvinchengx/gin-go-pg/controller"
	mw "github.com/calvinchengx/gin-go-pg/middleware"
	"github.com/calvinchengx/gin-go-pg/repository"
	"github.com/calvinchengx/gin-go-pg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	r := gin.Default()
	mw.Add(r, cors.Default())

	c, _ := config.Load("dev")
	db := config.GetConnection()
	log, _ := zap.NewDevelopment()
	defer log.Sync()
	userRepo := repository.NewUserRepo(db, log)
	jwt := mw.NewJWT(c.JWT)
	authService := controller.NewAuthService(userRepo, jwt)

	service.AuthRouter(authService, r)

	v1Router := r.Group("/v1")
	service.AccountRouter(controller.NewAccountService(), v1Router)
	service.UserRouter(v1Router)

	r.Run()
}
