package repo

import (
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database", err)
	}

	return mock, gormDB
}

func TestBackupRepo_Get(t *testing.T) {
	mock, gormDB := setupMockDB(t)
	oldDB := global.DB
	global.DB = gormDB
	t.Cleanup(func() { global.DB = oldDB })
	repo := NewIBackupRepo()

	t.Run("Success", func(t *testing.T) {
		expectedBackup := model.BackupAccount{
			BaseModel: model.BaseModel{ID: 1},
			Name: "test-backup",
			Type: "s3",
		}

		rows := sqlmock.NewRows([]string{"id", "name", "type", "created_at", "updated_at"}).
			AddRow(expectedBackup.ID, expectedBackup.Name, expectedBackup.Type, time.Now(), time.Now())

		mock.ExpectQuery("^SELECT \\* FROM `backup_accounts` ORDER BY `backup_accounts`.`id` LIMIT \\?$").WillReturnRows(rows)

		backup, err := repo.Get()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if backup.ID != expectedBackup.ID {
			t.Errorf("expected %v, got %v", expectedBackup, backup)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("^SELECT \\* FROM `backup_accounts` ORDER BY `backup_accounts`.`id` LIMIT \\?$").WillReturnError(gorm.ErrRecordNotFound)
		_, err := repo.Get()
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("expected ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("WithOptions", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(2)
		mock.ExpectQuery("^SELECT \\* FROM `backup_accounts` WHERE id = \\? ORDER BY `backup_accounts`.`id` LIMIT \\?$").
			WithArgs(2, 1).
			WillReturnRows(rows)

		opt := func(db *gorm.DB) *gorm.DB { return db.Where("id = ?", 2) }
		backup, err := repo.Get(opt)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if backup.ID != 2 {
			t.Errorf("expected ID 2, got %v", backup.ID)
		}
	})
}

func TestBackupRepo_List(t *testing.T) {
	mock, gormDB := setupMockDB(t)
	oldDB := global.DB
	global.DB = gormDB
	t.Cleanup(func() { global.DB = oldDB })
	repo := NewIBackupRepo()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("^SELECT \\* FROM `backup_accounts`$").WillReturnRows(rows)
		backups, err := repo.List()
		if err != nil || len(backups) != 1 || backups[0].ID != 1 {
			t.Errorf("unexpected error or result")
		}
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("^SELECT \\* FROM `backup_accounts`$").WillReturnError(errors.New("db error"))
		_, err := repo.List()
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("WithOptions", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(2)
		mock.ExpectQuery("^SELECT \\* FROM `backup_accounts` WHERE type = \\?$").WithArgs("s3").WillReturnRows(rows)
		opt := func(db *gorm.DB) *gorm.DB { return db.Where("type = ?", "s3") }
		backups, err := repo.List(opt)
		if err != nil || len(backups) != 1 || backups[0].ID != 2 {
			t.Errorf("unexpected error or result")
		}
	})
}

func TestBackupRepo_Page(t *testing.T) {
	mock, gormDB := setupMockDB(t)
	oldDB := global.DB
	global.DB = gormDB
	t.Cleanup(func() { global.DB = oldDB })
	repo := NewIBackupRepo()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `backup_accounts`$").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectQuery("^SELECT \\* FROM `backup_accounts` LIMIT \\?( OFFSET \\?)?$").WithArgs(10).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		count, backups, err := repo.Page(1, 10)
		if err != nil || count != 1 || len(backups) != 1 || backups[0].ID != 1 {
			t.Errorf("unexpected error or result")
		}
	})

	t.Run("WithOptions", func(t *testing.T) {
		mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `backup_accounts` WHERE type = \\?$").WithArgs("s3").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectQuery("^SELECT \\* FROM `backup_accounts` WHERE type = \\? LIMIT \\?( OFFSET \\?)?$").WithArgs("s3", 10, 10).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
		opt := func(db *gorm.DB) *gorm.DB { return db.Where("type = ?", "s3") }
		count, backups, err := repo.Page(2, 10, opt)
		if err != nil || count != 1 || len(backups) != 1 || backups[0].ID != 2 {
			t.Errorf("unexpected error or result")
		}
	})
}

func TestBackupRepo_Create(t *testing.T) {
	mock, gormDB := setupMockDB(t)
	oldDB := global.DB
	global.DB = gormDB
	t.Cleanup(func() { global.DB = oldDB })
	repo := NewIBackupRepo()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO `backup_accounts`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		err := repo.Create(&model.BackupAccount{})
		if err != nil {
			t.Errorf("unexpected error")
		}
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO `backup_accounts`").WillReturnError(errors.New("db error"))
		mock.ExpectRollback()
		err := repo.Create(&model.BackupAccount{})
		if err == nil {
			t.Errorf("expected error")
		}
	})
}

func TestBackupRepo_Save(t *testing.T) {
	mock, gormDB := setupMockDB(t)
	oldDB := global.DB
	global.DB = gormDB
	t.Cleanup(func() { global.DB = oldDB })
	repo := NewIBackupRepo()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE `backup_accounts`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		err := repo.Save(&model.BackupAccount{BaseModel: model.BaseModel{ID: 1}})
		if err != nil {
			t.Errorf("unexpected error")
		}
	})
}

func TestBackupRepo_Delete(t *testing.T) {
	mock, gormDB := setupMockDB(t)
	oldDB := global.DB
	global.DB = gormDB
	t.Cleanup(func() { global.DB = oldDB })
	repo := NewIBackupRepo()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM `backup_accounts` WHERE id = \\?$").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		err := repo.Delete(func(db *gorm.DB) *gorm.DB { return db.Where("id = ?", 1) })
		if err != nil {
			t.Errorf("unexpected error")
		}
	})
}
