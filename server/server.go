package server

import (
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gogjango/gjango/config"
	"github.com/gogjango/gjango/mail"
	mw "github.com/gogjango/gjango/middleware"
	"github.com/gogjango/gjango/mobile"
	"github.com/gogjango/gjango/route"

	"go.uber.org/zap"
)

// Server holds all the routes and their services
type Server struct {
	RouteServices []route.ServicesI
}

// Run runs our API server
func (server *Server) Run(env string) error {

	// load configuration
	c := config.Load(env)
	j := config.LoadJWT(env)

	r := gin.Default()

	// middleware
	mw.Add(r, cors.Default())
	jwt := mw.NewJWT(j)
	m := mail.NewMail(config.GetMailConfig(), config.GetSiteConfig())
	mobile := mobile.NewMobile(config.GetTwilioConfig())
	db := config.GetConnection()
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	// setup default routes
	rsDefault := &route.Services{
		DB:     db,
		Log:    log,
		JWT:    jwt,
		Mail:   m,
		Mobile: mobile,
		R:      r}
	rsDefault.SetupV1Routes()

	// setup all custom/user-defined route services
	for _, rs := range server.RouteServices {
		rs.SetupRoutes()
	}

	// run with port from config
	port := ":" + strconv.Itoa(c.Server.Port)
	return r.Run(port)
}
