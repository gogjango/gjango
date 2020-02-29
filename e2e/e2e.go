package e2e

import "github.com/calvinchengx/gin-go-pg/model"

// GetModels retrieve models
func GetModels() []interface{} {
	return model.Models
}
