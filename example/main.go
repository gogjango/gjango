package main

import (
	"github.com/gogjango/gjango"
	"github.com/gogjango/gjango/route"
)

func main() {
	gjango.New().
		WithRoutes(route.ServicesI{}).
		Run()
}
