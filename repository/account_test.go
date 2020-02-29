package repository_test

import (
	"net/http"
	"testing"

	"github.com/calvinchengx/gin-go-pg/apperr"
	"github.com/calvinchengx/gin-go-pg/mockgopg"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Emulate database error
func TestCreateAndVerifyDBError(t *testing.T) {
	db, mock, err := mockgopg.NewGoPGDBTest()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	log, _ := zap.NewDevelopment()
	accountRepo := repository.NewAccountRepo(db, log)

	u := &model.User{
		Username:    "hello",
		Email:       "hello@world.org",
		CountryCode: "+65",
		Mobile:      "91919191",
	}
	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u.Username, u.Email, u.CountryCode, u.Mobile).
		Returns(mockgopg.NewResult(0, 0, nil), apperr.DB)
	v, err := accountRepo.CreateAndVerify(u)
	assert.Nil(t, v)
	assert.Equal(t, apperr.DB, err)
}

func TestCreateAndVerify2(t *testing.T) {

	db, mock, err := mockgopg.NewGoPGDBTest()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	log, _ := zap.NewDevelopment()
	accountRepo := repository.NewAccountRepo(db, log)

	u2 := &model.User{
		Username:    "hello",
		Email:       "hello@world.org",
		CountryCode: "+65",
		Mobile:      "92929292",
	}
	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u2.Username, u2.Email, u2.CountryCode, u2.Mobile).
		Returns(mockgopg.NewResult(1, 1, u2), nil)
	v, err := accountRepo.CreateAndVerify(u2)
	assert.Nil(t, v)
	assert.Equal(t, apperr.New(http.StatusBadRequest, "User already exists."), err)
}

func TestCreateAndVerify3(t *testing.T) {

	db, mock, err := mockgopg.NewGoPGDBTest()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	log, _ := zap.NewDevelopment()
	accountRepo := repository.NewAccountRepo(db, log)

	u2 := &model.User{
		Username:    "hello",
		Email:       "hello@world.org",
		CountryCode: "+65",
		Mobile:      "92929292",
	}

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u2.Username, u2.Email, u2.CountryCode, u2.Mobile).
		Returns(mockgopg.NewResult(1, 1, u2), nil)
	v, err := accountRepo.CreateAndVerify(u2)
	assert.Nil(t, v)
	assert.Equal(t, apperr.New(http.StatusBadRequest, "User already exists."), err)
}

func TestCreateAndVerifyDBErrOnInsertUser(t *testing.T) {
	db, mock, err := mockgopg.NewGoPGDBTest()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	log, _ := zap.NewDevelopment()
	accountRepo := repository.NewAccountRepo(db, log)

	u2 := &model.User{
		Username:    "hello",
		Email:       "hello@world.org",
		CountryCode: "+65",
		Mobile:      "92929292",
	}

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u2.Username, u2.Email, u2.CountryCode, u2.Mobile).
		Returns(mockgopg.NewResult(0, 0, nil), apperr.NotFound)

	mock.ExpectInsert(u2).
		Returns(nil, apperr.DB)

	v, err := accountRepo.CreateAndVerify(u2)
	assert.Nil(t, v)
	assert.Equal(t, apperr.DB, err)
}

func TestCreateAndVerifyDBErrOnInsertVerififcation(t *testing.T) {
	db, mock, err := mockgopg.NewGoPGDBTest()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	log, _ := zap.NewDevelopment()
	accountRepo := repository.NewAccountRepo(db, log)

	u2 := &model.User{
		Username:    "hello",
		Email:       "hello@world.org",
		CountryCode: "+65",
		Mobile:      "92929292",
	}

	mock.ExpectQuery(`SELECT id FROM users WHERE username = ? OR email = ? OR (country_code = ? AND mobile = ?) AND deleted_at IS NULL`).
		WithArgs(u2.Username, u2.Email, u2.CountryCode, u2.Mobile).
		Returns(mockgopg.NewResult(0, 0, nil), apperr.NotFound)

	mock.ExpectInsert(u2).
		Returns(nil, nil)

	v := new(model.Verification)
	mock.ExpectInsert(v).
		Returns(nil, apperr.DB)

	v, err = accountRepo.CreateAndVerify(u2)
	assert.Nil(t, v)
	assert.Equal(t, apperr.DB, err)
}
