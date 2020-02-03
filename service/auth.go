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
	c.JSON(http.StatusOK, "logging in")
}

func refresh(c *gin.Context) {
	token := c.Param("token")
	c.JSON(http.StatusOK, token)
}
