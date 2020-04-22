package gjango

import (
	"github.com/gogjango/gjango/cmd"
	"github.com/gogjango/gjango/route"
)

// New creates a new Gjango instance
func New() *Gjango {
	return &Gjango{}
}

// Gjango allows us to specify customizations, such as custom route services
type Gjango struct {
	RouteServices []route.ServicesI
}

// WithRoutes is the builder method for us to add in custom route services
func (g *Gjango) WithRoutes(RouteServices []route.ServicesI) *Gjango {
	return &Gjango{RouteServices}
}

// Run executes our gjango functions or servers
func (g *Gjango) Run() {
	cmd.Execute(g.RouteServices)
}
