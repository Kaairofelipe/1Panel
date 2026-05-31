package repo

import (
	"regexp"
	"testing"

	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize gorm DB: %v", err)
	}

	return gormDB, mock
}

func TestBackupRepo_Get(t *testing.T) {
	gormDB, mock := setupTestDB(t)
	global.DB = gormDB

	repo := NewIBackupRepo()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `backup_accounts` ORDER BY `backup_accounts`.`id` LIMIT ?")).
            WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type"}).AddRow(1, "test_backup", "S3"))

		backup, err := repo.Get()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if backup.ID != 1 || backup.Name != "test_backup" {
			t.Errorf("unexpected result: %+v", backup)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `backup_accounts` ORDER BY `backup_accounts`.`id` LIMIT ?")).
            WithArgs(1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.Get()
		if err != gorm.ErrRecordNotFound {
			t.Errorf("expected ErrRecordNotFound, got %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("WithOption", func(t *testing.T) {
		opt := func(db *gorm.DB) *gorm.DB {
			return db.Where("id = ?", 2)
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `backup_accounts` WHERE id = ? ORDER BY `backup_accounts`.`id` LIMIT ?")).
			WithArgs(2, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(2, "backup_2"))

		backup, err := repo.Get(opt)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if backup.ID != 2 || backup.Name != "backup_2" {
			t.Errorf("unexpected result: %+v", backup)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestBackupRepo_Page(t *testing.T) {
	repo := NewIBackupRepo()

	t.Run("Success First Page", func(t *testing.T) {
        gormDB, mock := setupTestDB(t)
	    global.DB = gormDB

		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `backup_accounts`")).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `backup_accounts` LIMIT ?")).
            WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "b1"))

		count, list, err := repo.Page(1, 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if count != 2 {
			t.Errorf("expected count 2, got %d", count)
		}
		if len(list) != 1 || list[0].ID != 1 {
			t.Errorf("unexpected list: %+v", list)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Success Second Page", func(t *testing.T) {
        gormDB, mock := setupTestDB(t)
	    global.DB = gormDB

		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `backup_accounts`")).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `backup_accounts` LIMIT ? OFFSET ?")).
            WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(2, "b2"))

		count, list, err := repo.Page(2, 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if count != 2 {
			t.Errorf("expected count 2, got %d", count)
		}
		if len(list) != 1 || list[0].ID != 2 {
			t.Errorf("unexpected list: %+v", list)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

    t.Run("Count Error", func(t *testing.T) {
        gormDB, mock := setupTestDB(t)
	    global.DB = gormDB

		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `backup_accounts`")).
			WillReturnError(gorm.ErrInvalidDB)

		count, list, err := repo.Page(1, 1)
		if err != gorm.ErrInvalidDB {
			t.Errorf("expected ErrInvalidDB, got %v", err)
		}
        if count != 0 {
            t.Errorf("expected count 0, got %d", count)
        }
        if len(list) != 0 {
            t.Errorf("expected empty list, got %v", list)
        }

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Find Error", func(t *testing.T) {
        gormDB, mock := setupTestDB(t)
	    global.DB = gormDB

		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `backup_accounts`")).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `backup_accounts` LIMIT ?")).
            WithArgs(1).
			WillReturnError(gorm.ErrInvalidDB)

		count, list, err := repo.Page(1, 1)
		if err != gorm.ErrInvalidDB {
			t.Errorf("expected ErrInvalidDB, got %v", err)
		}
        if count != 2 {
            t.Errorf("expected count 2, got %d", count)
        }
        if len(list) != 0 {
            t.Errorf("expected empty list, got %v", list)
        }

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
