package model

import (
	"time"
)

// User represents user domain model
type User struct {
	Base
	ID          int        `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Username    string     `json:"username"`
	Password    string     `json:"-"`
	Email       string     `json:"email"`
	Mobile      string     `json:"mobile,omitempty"`
	CountryCode string     `json:"country_code,omitempty"`
	Address     string     `json:"address,omitempty"`
	LastLogin   *time.Time `json:"last_login,omitempty"`
	Active      bool       `json:"active"`
	Token       string     `json:"-"`
	Role        *Role      `json:"role,omitempty"`
	RoleID      int        `json:"-"`
	CompanyID   int        `json:"company_id"`
	LocationID  int        `json:"location_id"`
}

// UpdateLastLogin updates last login field
func (u *User) UpdateLastLogin() {
	t := time.Now()
	u.LastLogin = &t
}

// Delete updates the deleted_at field
func (u *User) Delete() {
	t := time.Now()
	u.DeletedAt = &t
}

// Update updates the updated_at field
func (u *User) Update() {
	t := time.Now()
	u.UpdatedAt = t
}

// UserRepo represents user database interface (the repository)
type UserRepo interface {
	View(int) (*User, error)
	FindByUsername(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByMobile(string, string) (*User, error)
	FindByToken(string) (*User, error)
	UpdateLogin(*User) error
	List(*ListQuery, *Pagination) ([]User, error)
	Update(*User) (*User, error)
	Delete(*User) error
}

// AccountRepo represents account database interface (the repository)
type AccountRepo interface {
	Create(*User) (*User, error)
	CreateAndVerify(*User) (*Verification, error)
	CreateWithMobile(*User) error
	ChangePassword(*User) error
	FindVerificationToken(string) (*Verification, error)
	DeleteVerificationToken(*Verification) error
}

// AuthUser represents data stored in JWT token for user
type AuthUser struct {
	ID         int
	CompanyID  int
	LocationID int
	Username   string
	Email      string
	Role       AccessRole
}
