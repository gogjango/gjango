package repository

import (
	"context"
	"fmt"
	"net/http"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/go-pg/pg/v9"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// NewAccountRepo returns an AccountRepo instance
func NewAccountRepo(db *pg.DB, log *zap.Logger) *AccountRepo {
	return &AccountRepo{db, log}
}

// AccountRepo represents the client for the user table
type AccountRepo struct {
	db  *pg.DB
	log *zap.Logger
}

// Create creates a new user in our database
func (a *AccountRepo) Create(c context.Context, u *model.User) error {
	user := new(model.User)
	res, err := a.db.Query(user, "SELECT id FROM users WHERE username = ? or email = ? AND deleted_at IS NULL", u.Username, u.Email)
	if err != nil {
		a.log.Error("AccountRepo Error: ", zap.Error(err))
		return apperr.DB
	}
	if res.RowsReturned() != 0 {
		return apperr.New(http.StatusBadRequest, "Username or email already exists.")
	}
	if err := a.db.Insert(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
	}
	return nil
}

// CreateAndVerify creates a new user in our database, and generates a verification token.
// User active being false until after verification.
func (a *AccountRepo) CreateAndVerify(c context.Context, u *model.User) error {
	user := new(model.User)
	fmt.Println(user.Active)
	res, err := a.db.Query(user, "SELECT id FROM users WHERE username = ? or email = ? AND deleted_at IS NULL", u.Username, u.Email)
	if err != nil {
		a.log.Error("AccountRepo Error: ", zap.Error(err))
		return apperr.DB
	}
	if res.RowsReturned() != 0 {
		return apperr.New(http.StatusBadRequest, "Username or email already exists.")
	}

	if err := a.db.Insert(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
	}

	v := new(model.Verification)
	v.UserID = u.ID
	v.Token = uuid.NewV4().String()
	if err := a.db.Insert(v); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
	}
	return nil
}

// ChangePassword changes user's password
func (a *AccountRepo) ChangePassword(c context.Context, u *model.User) error {
	_, err := a.db.Model(u).Column("password", "updated_at").WherePK().Update()
	if err != nil {
		a.log.Warn("AccountRepo Error: ", zap.Error(err))
	}
	return err
}
