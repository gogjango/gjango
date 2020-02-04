package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthRouter creates new auth http service
func AuthRouter(r *gin.Engine) {
	r.POST("/login", login)
	r.GET("/refresh/:token", refresh)
}

func login(c *gin.Context) {
	// cred, err := request.Login(c)
	// if err != nil {
	// 	return
	// }
	// r, err := a.svc.Authenticate(c, cred.Username, cred.Password)
	// if err != nil {
	// apperr.Response(c, err)
	// return
	// }
	c.JSON(http.StatusOK, "logging in")
}

func refresh(c *gin.Context) {
	token := c.Param("token")
	c.JSON(http.StatusOK, token)
}
