package repo

import (
	"testing"
	"time"

	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCommandRepo_Get(t *testing.T) {
	// Setup mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Setup GORM with the mock DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("failed to setup GORM with mock DB: %v", err)
	}

	// Backup and restore global DB instance
	originalDB := global.DB
	defer func() {
		global.DB = originalDB
	}()
	global.DB = gormDB

	repo := NewICommandRepo()

	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `commands` ORDER BY `commands`.`id` LIMIT (\\?)?$").
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "type", "name", "group_id", "command"}).
				AddRow(1, now, now, "test-type", "test-name", 2, "echo 'test'"))

		cmd, err := repo.Get()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if cmd.ID != 1 || cmd.Name != "test-name" {
			t.Errorf("unexpected command returned: %+v", cmd)
		}
	})

	t.Run("with options", func(t *testing.T) {
		// LIMIT parameter needs to be accounted for in the arguments list for sqlmock
		mock.ExpectQuery("SELECT \\* FROM `commands` WHERE name like \\? or command like \\? ORDER BY `commands`.`id` LIMIT (\\?)?$").
			WithArgs("%test%", "%test%", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "type", "name", "group_id", "command"}).
				AddRow(1, now, now, "test-type", "test-name", 2, "echo 'test'"))

		cmd, err := repo.Get(repo.WithByInfo("test"))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if cmd.ID != 1 || cmd.Name != "test-name" {
			t.Errorf("unexpected command returned: %+v", cmd)
		}
	})

	t.Run("not found", func(t *testing.T) {
		// No args provided to Get, expect limit 1
		mock.ExpectQuery("SELECT \\* FROM `commands` ORDER BY `commands`.`id` LIMIT (\\?)?$").
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.Get()
		if err != gorm.ErrRecordNotFound {
			t.Errorf("expected ErrRecordNotFound, got %v", err)
		}
	})

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
