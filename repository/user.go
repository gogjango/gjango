package repository

import (
	"context"

	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
)

// NewUserRepo returns a new UserRepo instance
func NewUserRepo(db *pg.DB, log *zap.Logger) *UserRepo {
	return &UserRepo{db, log}
}

// UserRepo is the client for our user model
type UserRepo struct {
	db  *pg.DB
	log *zap.Logger
}

// FindByUsername queries for a single user by username
func (u *UserRepo) FindByUsername(c context.Context, username string) (*model.User, error) {
	user := new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, username)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.NotFound
	}
	return user, nil
}

// FindByToken queries for single user by token
func (u *UserRepo) FindByToken(c context.Context, token string) (*model.User, error) {
	var user = new(model.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)`
	_, err := u.db.QueryOne(user, sql, token)
	if err != nil {
		u.log.Warn("UserRepo Error", zap.String("Error:", err.Error()))
		return nil, apperr.NotFound
	}
	return user, nil
}

// UpdateLogin updates last login and refresh token for user
func (u *UserRepo) UpdateLogin(c context.Context, user *model.User) error {
	_, err := u.db.Model(user).Column("last_login", "token").WherePK().Update()
	if err != nil {
		u.log.Warn("UserRepo Error", zap.Error(err))
	}
	return err
}
