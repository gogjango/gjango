package mockdb

import (
	"context"

	"github.com/calvinchengx/gin-go-pg/model"
)

// User database mock
type User struct {
	ViewFn           func(int) (*model.User, error)
	FindByUsernameFn func(context.Context, string) (*model.User, error)
	FindByEmailFn    func(string) (*model.User, error)
	FindByMobileFn   func(string, string) (*model.User, error)
	FindByTokenFn    func(context.Context, string) (*model.User, error)
	UpdateLoginFn    func(context.Context, *model.User) error
	ListFn           func(context.Context, *model.ListQuery, *model.Pagination) ([]model.User, error)
	DeleteFn         func(context.Context, *model.User) error
	UpdateFn         func(context.Context, *model.User) (*model.User, error)
}

// View mock
func (u *User) View(id int) (*model.User, error) {
	return u.ViewFn(id)
}

// FindByUsername mock
func (u *User) FindByUsername(c context.Context, username string) (*model.User, error) {
	return u.FindByUsernameFn(c, username)
}

// FindByEmail mock
func (u *User) FindByEmail(email string) (*model.User, error) {
	return u.FindByEmailFn(email)
}

// FindByMobile mock
func (u *User) FindByMobile(countryCode, mobile string) (*model.User, error) {
	return u.FindByMobileFn(countryCode, mobile)
}

// FindByToken mock
func (u *User) FindByToken(c context.Context, token string) (*model.User, error) {
	return u.FindByTokenFn(c, token)
}

// UpdateLogin mock
func (u *User) UpdateLogin(c context.Context, usr *model.User) error {
	return u.UpdateLoginFn(c, usr)
}

// List mock
func (u *User) List(c context.Context, lq *model.ListQuery, p *model.Pagination) ([]model.User, error) {
	return u.ListFn(c, lq, p)
}

// Delete mock
func (u *User) Delete(c context.Context, usr *model.User) error {
	return u.DeleteFn(c, usr)
}

// Update mock
func (u *User) Update(c context.Context, usr *model.User) (*model.User, error) {
	return u.UpdateFn(c, usr)
}
