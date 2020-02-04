package repository

import (
	"context"
	"net/http"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/go-pg/pg/v9"
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

// ChangePassword changes user's password
func (a *AccountRepo) ChangePassword(c context.Context, usr *model.User) error {
	_, err := a.db.Model(usr).Column("password", "updated_at").WherePK().Update()
	if err != nil {
		a.log.Warn("AccountRepo Error: ", zap.Error(err))
	}
	return err
}
