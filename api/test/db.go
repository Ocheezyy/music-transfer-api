package test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ocheezyy/music-transfer-api/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMockDB returns a mock DB and the sqlmock instance
func NewMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		helpers.CoreLogError("NewMockDB", fmt.Sprintf("failed to create sqlmock: %s", err), true)
		t.Fatalf("failed to create sqlmock: %s", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		helpers.CoreLogError("NewMockDB", fmt.Sprintf("failed to open gorm DB: %s", err), true)
	}

	return gormDB, mock
}
