package repo

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("3.35.0"))
	dialector := sqlite.Dialector{
		Conn: sqlDB,
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Fatalf("failed to create gorm DB: %v", err)
	}

	return gormDB, mock, sqlDB
}

func TestLogRepo_CleanLogin(t *testing.T) {
	gormDB, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()
	global.DB = gormDB

	repo := NewILogRepo()

	mock.ExpectExec("delete from login_logs;").WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.CleanLogin()
	if err != nil {
		t.Errorf("CleanLogin() error = %v, wantErr nil", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLogRepo_CreateLoginLog(t *testing.T) {
	gormDB, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()
	global.DB = gormDB

	repo := NewILogRepo()


	log := &model.LoginLog{
		IP:        "127.0.0.1",
		Address:   "Local",
		Status:    "Success",
		Message:   "Logged in",
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `login_logs`")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := repo.CreateLoginLog(log)
	if err != nil {
		t.Errorf("CreateLoginLog() error = %v, wantErr nil", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLogRepo_PageLoginLog(t *testing.T) {
	gormDB, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()
	global.DB = gormDB

	repo := NewILogRepo()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `login_logs`")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `login_logs` LIMIT 10")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ip", "address", "status", "message"}).
			AddRow(1, "127.0.0.1", "Local", "Success", "Logged in"))

	count, logs, err := repo.PageLoginLog(1, 10)
	if err != nil {
		t.Errorf("PageLoginLog() error = %v, wantErr nil", err)
	}
	if count != 1 {
		t.Errorf("PageLoginLog() count = %v, want 1", count)
	}
	if len(logs) != 1 {
		t.Errorf("PageLoginLog() logs count = %v, want 1", len(logs))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLogRepo_CleanOperation(t *testing.T) {
	gormDB, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()
	global.DB = gormDB

	repo := NewILogRepo()

	mock.ExpectExec("delete from operation_logs").WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.CleanOperation()
	if err != nil {
		t.Errorf("CleanOperation() error = %v, wantErr nil", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLogRepo_CreateOperationLog(t *testing.T) {
	gormDB, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()
	global.DB = gormDB

	repo := NewILogRepo()

	log := &model.OperationLog{
		Source:   "System",
		IP:       "127.0.0.1",
		DetailZH: "测试",
		DetailEN: "Test",
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `operation_logs`")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := repo.CreateOperationLog(log)
	if err != nil {
		t.Errorf("CreateOperationLog() error = %v, wantErr nil", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLogRepo_PageOperationLog(t *testing.T) {
	gormDB, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()
	global.DB = gormDB

	repo := NewILogRepo()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `operation_logs`")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `operation_logs` LIMIT 10")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "source", "ip", "detail_zh", "detail_en"}).
			AddRow(1, "System", "127.0.0.1", "测试", "Test"))

	count, logs, err := repo.PageOperationLog(1, 10)
	if err != nil {
		t.Errorf("PageOperationLog() error = %v, wantErr nil", err)
	}
	if count != 1 {
		t.Errorf("PageOperationLog() count = %v, want 1", count)
	}
	if len(logs) != 1 {
		t.Errorf("PageOperationLog() logs count = %v, want 1", len(logs))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLogRepo_WithByStatus(t *testing.T) {
	gormDB, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	repo := &LogRepo{}

	opt := repo.WithByStatus("Success")
	db := opt(gormDB.Model(&model.LoginLog{}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `login_logs` WHERE status = ?")).
		WithArgs("Success").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var logs []model.LoginLog
	err := db.Find(&logs).Error
	if err != nil {
		t.Errorf("WithByStatus() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLogRepo_WithBySource(t *testing.T) {
	gormDB, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	repo := &LogRepo{}

	opt := repo.WithBySource("System")
	db := opt(gormDB.Model(&model.OperationLog{}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `operation_logs` WHERE source = ?")).
		WithArgs("System").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var logs []model.OperationLog
	err := db.Find(&logs).Error
	if err != nil {
		t.Errorf("WithBySource() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLogRepo_WithByIP(t *testing.T) {
	gormDB, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	repo := &LogRepo{}

	opt := repo.WithByIP("127.0")
	db := opt(gormDB.Model(&model.LoginLog{}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `login_logs` WHERE ip LIKE ?")).
		WithArgs("%127.0%").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var logs []model.LoginLog
	err := db.Find(&logs).Error
	if err != nil {
		t.Errorf("WithByIP() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestLogRepo_WithByLikeOperation(t *testing.T) {
	gormDB, mock, sqlDB := setupTestDB(t)
	defer sqlDB.Close()

	repo := &LogRepo{}

	opt := repo.WithByLikeOperation("test")
	db := opt(gormDB.Model(&model.OperationLog{}))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `operation_logs` WHERE detail_zh LIKE ? OR detail_en LIKE ?")).
		WithArgs("%test%", "%test%").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var logs []model.OperationLog
	err := db.Find(&logs).Error
	if err != nil {
		t.Errorf("WithByLikeOperation() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
