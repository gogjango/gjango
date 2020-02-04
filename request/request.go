package request

import (
	"net/http"
	"strconv"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/gin-gonic/gin"
)

// ID returns id url parameter.
// In case of conversion error to int, request will be aborted with StatusBadRequest.
func ID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return 0, apperr.BadRequest
	}
	return id, nil
}
