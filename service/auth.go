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
	r.POST("/signup/m", a.signupm) // mobile
	r.POST("/signup", a.signup)    // email
	r.POST("/login", a.login)
	r.GET("/refresh/:token", a.refresh)
	r.GET("/verification/:token", a.verify) // email
	r.POST("/verifycode", a.verifym)        // mobile
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

func (a *Auth) signup(c *gin.Context) {
	e, err := request.AccountSignup(c)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	err = a.svc.Signup(c, e)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

func (a *Auth) verify(c *gin.Context) {
	token := c.Param("token")
	err := a.svc.Verify(c, token)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (a *Auth) signupm(c *gin.Context) {
	m, err := request.AccountSignupMobile(c)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	err = a.svc.SignupMobile(c, m)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

func (a *Auth) verifym(c *gin.Context) {
	m, err := request.AccountVerifyMobile(c)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	err = a.svc.VerifyMobile(c, m.CountryCode, m.Mobile, m.Code)
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.Status(http.StatusOK)
}
