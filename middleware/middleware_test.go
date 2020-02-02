package middleware_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ribice/gorsk-gin/cmd/api/mw"
)

func TestAdd(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mw.Add(r, gin.Logger())
}
