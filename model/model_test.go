package model_test

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/calvinchengx/gin-go-pg/mock"
	"github.com/calvinchengx/gin-go-pg/model"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/stretchr/testify/assert"
)

func TestBeforeInsert(t *testing.T) {
	base := &model.Base{
		ID: 1,
	}
	base.BeforeInsert(nil)
	if base.CreatedAt.IsZero() {
		t.Errorf("CreatedAt was not changed")
	}
	if base.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt was not changed")
	}
}

func TestBeforeUpdate(t *testing.T) {
	base := &model.Base{
		ID:        1,
		CreatedAt: mock.TestTime(2000),
	}
	base.BeforeUpdate(nil)
	if base.UpdatedAt == mock.TestTime(2001) {
		t.Errorf("UpdatedAt was not changed")
	}

}

func TestDelete(t *testing.T) {
	baseModel := &model.Base{
		ID:        1,
		CreatedAt: mock.TestTime(2000),
		UpdatedAt: mock.TestTime(2001),
	}
	baseModel.Delete()
	if baseModel.DeletedAt.IsZero() {
		t.Errorf("DeletedAt not changed")
	}

}

func TestDatabase(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	projectRoot := filepath.Dir(d)
	tmpDir := path.Join(projectRoot, "tmp")
	testConfig := embeddedpostgres.DefaultConfig().
		Username("db_test_user").
		Password("db_test_password").
		Database("db_test_database").
		Version("12.1.0").
		RuntimePath(tmpDir).
		Port(9876)

	postgres := embeddedpostgres.NewDatabase(testConfig)
	err := postgres.Start()
	assert.Equal(t, err, nil)
	_ = postgres.Stop()
}
