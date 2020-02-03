package main

import (
	mw "github.com/calvinchengx/gin-go-pg/middleware"
	"github.com/calvinchengx/gin-go-pg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	mw.Add(r, cors.Default())

	v1Router := r.Group("/v1")
	service.UserRouter(v1Router)

	r.Run()
}
