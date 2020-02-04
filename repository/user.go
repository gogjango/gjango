package repository

import (
	"context"

	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/model"
)

type UserRepo struct {
	db  *pg.DB
	log *zap.Logger
}

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
