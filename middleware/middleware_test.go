package middleware_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	mw "github.com/gogjango/gjango/middleware"
)

func TestAdd(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mw.Add(r, gin.Logger())
}
