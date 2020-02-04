package service

import (
	"net/http"
	"strconv"

	"github.com/calvinchengx/gin-go-pg/request"
	"github.com/gin-gonic/gin"
)

// UserRouter declares the orutes for users router group
func UserRouter(r *gin.RouterGroup) {
	ur := r.Group("/users")
	ur.GET("", list)
	ur.GET("/:id", view)
	ur.PATCH("/:id", update)
	ur.DELETE("/:id", delete)
}

func list(c *gin.Context) {
	// retrieve from database
	c.JSON(http.StatusOK, "list of users")
}

func view(c *gin.Context) {
	id, err := request.ID(c)
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
