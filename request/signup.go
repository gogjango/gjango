package request

import (
	"net/http"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/gin-gonic/gin"
)

// Signup contains the user signup request
type Signup struct {
	Username        string `json:"username" binding:"required,min=3,alphanum"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

// AccountSignup validates user signup request
func AccountSignup(c *gin.Context) (*Signup, error) {
	var r Signup
	if err := c.ShouldBindJSON(&r); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	if r.Password != r.PasswordConfirm {
		err := apperr.New(http.StatusBadRequest, "passwords do not match")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return nil, err
	}
	return &r, nil
}
