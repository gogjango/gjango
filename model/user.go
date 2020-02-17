package model

import (
	"context"
	"time"
)

// User represents user domain model
type User struct {
	Base
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

// UserRepo represents user database interface (the repository)
type UserRepo interface {
	View(context.Context, int) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	FindByEmail(context.Context, string) (*User, error)
	FindByMobile(context.Context, string, string) (*User, error)
	FindByToken(context.Context, string) (*User, error)
	UpdateLogin(context.Context, *User) error
	List(context.Context, *ListQuery, *Pagination) ([]User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, *User) error
}

// AccountRepo represents account database interface (the repository)
type AccountRepo interface {
	Create(context.Context, *User) error
	CreateAndVerify(context.Context, *User) (*Verification, error)
	CreateWithMobile(context.Context, *User) error
	ChangePassword(context.Context, *User) error
	FindVerificationToken(context.Context, string) (*Verification, error)
	DeleteVerificationToken(context.Context, *Verification) error
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
