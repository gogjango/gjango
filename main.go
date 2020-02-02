package main

import (
	"os"

	mw "github.com/calvinchengx/gin-go-pg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	handleMigration(os.Args)

	r := gin.Default()

	mw.Add(r, cors.Default())

	r.Run()
}
