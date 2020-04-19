package request

import (
	"github.com/gin-gonic/gin"
	"github.com/gogjango/gjango/apperr"
)

// Credentials stores the username and password provided in the request
type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login parses out the username and password in gin's request context, into Credentials
func Login(c *gin.Context) (*Credentials, error) {
	cred := new(Credentials)
	if err := c.ShouldBindJSON(cred); err != nil {
		apperr.Response(c, err)
		return nil, err
	}
	return cred, nil
}
