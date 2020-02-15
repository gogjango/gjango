package mockdb

import (
	"context"

	"github.com/calvinchengx/gin-go-pg/model"
)

// Account database mock
type Account struct {
	CreateFn                  func(context.Context, *model.User) error
	CreateAndVerifyFn         func(context.Context, *model.User) (*model.Verification, error)
	ChangePasswordFn          func(context.Context, *model.User) error
	FindVerificationTokenFn   func(context.Context, string) (*model.Verification, error)
	DeleteVerificationTokenFn func(context.Context, *model.Verification) error
}

// Create mock
func (a *Account) Create(c context.Context, usr *model.User) error {
	return a.CreateFn(c, usr)
}

// CreateAndVerify mock
func (a *Account) CreateAndVerify(c context.Context, usr *model.User) (*model.Verification, error) {
	return a.CreateAndVerifyFn(c, usr)
}

// ChangePassword mock
func (a *Account) ChangePassword(c context.Context, usr *model.User) error {
	return a.ChangePasswordFn(c, usr)
}

// FindVerificationToken mock
func (a *Account) FindVerificationToken(c context.Context, token string) (*model.Verification, error) {
	return a.FindVerificationTokenFn(c, token)
}

// DeleteVerificationToken mock
func (a *Account) DeleteVerificationToken(c context.Context, v *model.Verification) error {
	return a.DeleteVerificationTokenFn(c, v)
}
