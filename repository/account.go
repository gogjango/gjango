package repository

import (
	"net/http"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/go-pg/pg/v9/orm"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// NewAccountRepo returns an AccountRepo instance
func NewAccountRepo(db orm.DB, log *zap.Logger) *AccountRepo {
	return &AccountRepo{db, log}
}

// AccountRepo represents the client for the user table
type AccountRepo struct {
	db  orm.DB
	log *zap.Logger
}

// Create creates a new user in our database
func (a *AccountRepo) Create(u *model.User) (*model.User, error) {
	user := new(model.User)
	sql := `SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`
	res, err := a.db.Query(user, sql, u.Username, u.Email, u.CountryCode, u.Mobile)
	if err != nil {
		a.log.Error("AccountRepo Error: ", zap.Error(err))
		return nil, apperr.DB
	}
	if res.RowsReturned() != 0 {
		return nil, apperr.New(http.StatusBadRequest, "User already exists.")
	}
	if err := a.db.Insert(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return nil, apperr.DB
	}
	return u, nil
}

// CreateAndVerify creates a new user in our database, and generates a verification token.
// User active being false until after verification.
func (a *AccountRepo) CreateAndVerify(u *model.User) (*model.Verification, error) {
	user := new(model.User)
	sql := `SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`
	res, err := a.db.Query(user, sql, u.Username, u.Email, u.CountryCode, u.Mobile)
	if err != nil {
		a.log.Error("AccountRepo Error: ", zap.Error(err))
		return nil, apperr.DB
	}
	if res.RowsReturned() != 0 {
		return nil, apperr.New(http.StatusBadRequest, "User already exists.")
	}
	if err := a.db.Insert(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return nil, apperr.DB
	}
	v := new(model.Verification)
	v.UserID = u.ID
	v.Token = uuid.NewV4().String()
	if err := a.db.Insert(v); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return nil, apperr.DB
	}
	return v, nil
}

// CreateWithMobile creates a new user in our database with country code and mobile number
func (a *AccountRepo) CreateWithMobile(u *model.User) error {
	user := new(model.User)
	sql := `SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`
	res, err := a.db.Query(user, sql, u.Username, u.Email, u.CountryCode, u.Mobile)
	if err != nil {
		a.log.Error("AccountRepo Error: ", zap.Error(err))
		return apperr.DB
	}
	if res.RowsReturned() != 0 && user.Active == true {
		return apperr.NewStatus(http.StatusConflict) // user already exists and is already verified
	}
	if res.RowsReturned() != 0 {
		return apperr.BadRequest // user already exists but is not yet verified
	}
	if err := a.db.Insert(u); err != nil {
		a.log.Warn("AccountRepo error: ", zap.Error(err))
		return apperr.DB
	}
	return nil
}

// ChangePassword changes user's password
func (a *AccountRepo) ChangePassword(u *model.User) error {
	u.Update()
	_, err := a.db.Model(u).Column("password", "updated_at").WherePK().Update()
	if err != nil {
		a.log.Warn("AccountRepo Error: ", zap.Error(err))
	}
	return err
}

// FindVerificationToken retrieves an existing verification token
func (a *AccountRepo) FindVerificationToken(token string) (*model.Verification, error) {
	var v = new(model.Verification)
	sql := `SELECT * FROM verifications WHERE (token = ? and deleted_at IS NULL)`
	_, err := a.db.QueryOne(v, sql, token)
	if err != nil {
		a.log.Warn("AccountRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.NotFound
	}
	return v, nil
}

// DeleteVerificationToken sets deleted_at for an existing verification token
func (a *AccountRepo) DeleteVerificationToken(v *model.Verification) error {
	_, err := a.db.Model(v).Column("deleted_at").WherePK().Update()
	if err != nil {
		a.log.Warn("AccountRepo Error", zap.Error(err))
		return apperr.DB
	}
	return err
}
