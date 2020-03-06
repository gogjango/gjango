package server

import (
	"strconv"

	"github.com/calvinchengx/gin-go-pg/config"
	"github.com/calvinchengx/gin-go-pg/mail"
	mw "github.com/calvinchengx/gin-go-pg/middleware"
	"github.com/calvinchengx/gin-go-pg/mobile"
	"github.com/calvinchengx/gin-go-pg/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

// Run runs our API server
func Run() error {
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
	rs := &route.Services{
		DB:     db,
		Log:    log,
		JWT:    jwt,
		Mail:   m,
		Mobile: mobile,
		R:      r}
	rs.SetupV1Routes()

	// run with port from config
	port := ":" + strconv.Itoa(c.Server.Port)
	return r.Run(port)
}
