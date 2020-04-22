package main

import (
	"fmt"

	"github.com/gogjango/gjango"
)

func main() {
	gjango.New().
		WithRoutes(&MyServices{}).
		Run()
}

// MyServices implements github.com/gogjango/gjango/route.ServicesI
type MyServices struct{}

// SetupRoutes is our implementation of custom routes
func (s *MyServices) SetupRoutes() {
	fmt.Println("set up our custom routes!")
}
