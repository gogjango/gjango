package service

import (
	"net/http"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/repository/auth"
	"github.com/calvinchengx/gin-go-pg/request"
	"github.com/gin-gonic/gin"
)

// AuthRouter creates new auth http service
func AuthRouter(svc *auth.Service, r *gin.Engine) {
	a := Auth{svc}
	r.POST("/login", a.login)
	r.GET("/refresh/:token", a.refresh)
}

// Auth represents auth http service
type Auth struct {
	svc *auth.Service
}

func (a *Auth) login(c *gin.Context) {
	cred, err := request.Login(c)
	if err != nil {
		return
	}
	r, err := a.svc.Authenticate(c, cred.Username, cred.Password)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, r)
}

func (a *Auth) refresh(c *gin.Context) {
	r, err := a.svc.Refresh(c, c.Param("token"))
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, r)
}
