package e2e

import (
	"fmt"

	"github.com/calvinchengx/gin-go-pg/manager"
)

// SetupDatabase creates the schema and populates it with data
func SetupDatabase(m *manager.Manager) {
	models := manager.GetModels()
	m.CreateSchema(models...)
	m.CreateRoles()
	superUser, err := m.CreateSuperAdmin("superuser@example.org", "testpassword")
	fmt.Println("superUser", superUser)
	fmt.Println("Is there an error?", err)
}
