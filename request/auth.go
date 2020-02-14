package request

import (
	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/gin-gonic/gin"
)

// Credentials stores the username and password provided in the request
type Credentials struct {
	Username string `json:"username" binding:"required"`
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
