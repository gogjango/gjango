package service

import (
	"net/http"
	"strconv"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	ur := r.Group("/users")
	ur.GET("", list)
	ur.GET("/:id", view)
	ur.PATCH("/:id", update)
	ur.DELETE("/:id", delete)
}

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

func list(c *gin.Context) {
	// retrieve from database
	c.JSON(http.StatusOK, "list of users")
}

func view(c *gin.Context) {
	id, err := ID(c)
	if err != nil {
		return
	}
	// retrieve from database
	c.JSON(http.StatusOK, "view user "+strconv.Itoa(id))
}

func update(c *gin.Context) {
	// update database
	c.JSON(http.StatusOK, "update user")
}

func delete(c *gin.Context) {
	// delete user from database
	c.JSON(http.StatusOK, "delete user")
}
