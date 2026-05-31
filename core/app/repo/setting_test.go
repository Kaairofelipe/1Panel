package repo

import (
	"errors"
	"testing"
	"time"

	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/patrickmn/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupSqlMock(t *testing.T) (sqlmock.Sqlmock, func()) {
	oldDB := global.DB

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to init gorm db: %v", err)
	}

	global.DB = gormDB

	// Reset cache
	settingCache = cache.New(5*time.Minute, 10*time.Minute)

	return mock, func() {
		db.Close()
		global.DB = oldDB
	}
}

func TestSettingRepo_Create(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		mockSetup func(sqlmock.Sqlmock)
		wantErr   bool
		err       error
	}{
		{
			name:  "success",
			key:   "test_key",
			value: "test_value",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `settings`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "db error",
			key:   "test_key",
			value: "test_value",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `settings`").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			err:     errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, cleanup := setupSqlMock(t)
			defer cleanup()

			repo := NewISettingRepo()

			tt.mockSetup(mock)
			err := repo.Create(tt.key, tt.value)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.err.Error() {
					t.Errorf("expected error %v, got %v", tt.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				// Verify cache
				val, found := settingCache.Get(tt.key)
				if !found {
					t.Errorf("expected setting to be in cache")
				} else if val.(string) != tt.value {
					t.Errorf("expected cache value %v, got %v", tt.value, val)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
