package model

func init() {
	Register(&Location{})
}

// Location represents company location model
type Location struct {
	Base
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	Address   string `json:"address"`
	CompanyID int    `json:"company_id"`
}
