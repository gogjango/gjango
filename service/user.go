package service

import (
	"net/http"
	"strconv"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository/user"
	"github.com/calvinchengx/gin-go-pg/request"
	"github.com/gin-gonic/gin"
)

// User represents the user http service
type User struct {
	svc *user.Service
}

// UserRouter declares the orutes for users router group
func UserRouter(svc *user.Service, r *gin.RouterGroup) {
	u := User{
		svc: svc,
	}
	ur := r.Group("/users")
	ur.GET("", u.list)
	ur.GET("/:id", u.view)
	ur.PATCH("/:id", u.update)
	ur.DELETE("/:id", u.delete)
}

type listResponse struct {
	Users []model.User `json:"users"`
	Page  int          `json:"page"`
}

func (u *User) list(c *gin.Context) {
	p, err := request.Paginate(c)
	if err != nil {
		return
	}
	result, err := u.svc.List(c, &model.Pagination{
		Limit: p.Limit, Offset: p.Offset,
	})
	if err != nil {
		apperr.Response(c, err)
		return
	}
	c.JSON(http.StatusOK, listResponse{
		Users: result,
		Page:  p.Page,
	})
}

func (u *User) view(c *gin.Context) {
	id, err := request.ID(c)
	if err != nil {
		return
	}
	// retrieve from database
	c.JSON(http.StatusOK, "view user "+strconv.Itoa(id))
}

func (u *User) update(c *gin.Context) {
	// update database
	c.JSON(http.StatusOK, "update user")
}

func (u *User) delete(c *gin.Context) {
	// delete user from database
	c.JSON(http.StatusOK, "delete user")
}
