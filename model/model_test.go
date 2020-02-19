package model_test

import (
	"testing"

	"github.com/calvinchengx/gin-go-pg/mock"
	"github.com/calvinchengx/gin-go-pg/model"
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
