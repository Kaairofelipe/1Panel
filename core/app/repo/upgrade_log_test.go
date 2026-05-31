package repo

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// AnyTime can match any time.Time value in sqlmock
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func setupTestDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	// create sqlmock
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}

	// Mock the sqlite_version() query that gorm driver executes
	mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{"sqlite_version()"}).AddRow("3.35.0"))

	// create gorm db
	gormDB, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{
		AllowGlobalUpdate: true, // Allow global updates for delete tests
	})
	if err != nil {
		t.Fatalf("Failed to open gorm db: %v", err)
	}

	// save original and replace with mock
	originalDB := global.DB
	global.DB = gormDB

	cleanup := func() {
		global.DB = originalDB
		sqlDB.Close()
	}

	return mock, cleanup
}

func TestUpgradeLogRepo_Get(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewIUpgradeLogRepo()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "node_id", "old_version", "new_version", "backup_file"}).
			AddRow(1, 10, "v1.0.0", "v1.1.0", "/backup/file.tar.gz")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `upgrade_logs` ORDER BY `upgrade_logs`.`id` LIMIT 1")).
			WillReturnRows(rows)

		log, err := repo.Get()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if log.ID != 1 {
			t.Errorf("Expected log ID 1, got %d", log.ID)
		}
		if log.NodeID != 10 {
			t.Errorf("Expected node ID 10, got %d", log.NodeID)
		}
		if log.OldVersion != "v1.0.0" {
			t.Errorf("Expected old version v1.0.0, got %s", log.OldVersion)
		}
		if log.NewVersion != "v1.1.0" {
			t.Errorf("Expected new version v1.1.0, got %s", log.NewVersion)
		}
	})

	t.Run("With Options", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "node_id"}).AddRow(2, 20)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `upgrade_logs` WHERE node_id = ? ORDER BY `upgrade_logs`.`id` LIMIT 1")).
			WithArgs(uint(20)).
			WillReturnRows(rows)

		log, err := repo.Get(repo.WithByNodeID(20))

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if log.ID != 2 {
			t.Errorf("Expected log ID 2, got %d", log.ID)
		}
		if log.NodeID != 20 {
			t.Errorf("Expected node ID 20, got %d", log.NodeID)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `upgrade_logs` ORDER BY `upgrade_logs`.`id` LIMIT 1")).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.Get()

		if err != gorm.ErrRecordNotFound {
			t.Errorf("Expected ErrRecordNotFound, got %v", err)
		}
	})
}

func TestUpgradeLogRepo_List(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewIUpgradeLogRepo()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "node_id", "old_version", "new_version"}).
			AddRow(1, 10, "v1.0.0", "v1.1.0").
			AddRow(2, 10, "v1.1.0", "v1.2.0")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `upgrade_logs`")).
			WillReturnRows(rows)

		logs, err := repo.List()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(logs) != 2 {
			t.Errorf("Expected 2 logs, got %d", len(logs))
		}
		if logs[0].ID != 1 || logs[1].ID != 2 {
			t.Errorf("Expected log IDs 1 and 2, got %d and %d", logs[0].ID, logs[1].ID)
		}
	})

	t.Run("With Options", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "node_id"}).
			AddRow(1, 10).
			AddRow(2, 10)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `upgrade_logs` WHERE node_id = ?")).
			WithArgs(uint(10)).
			WillReturnRows(rows)

		logs, err := repo.List(repo.WithByNodeID(10))

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(logs) != 2 {
			t.Errorf("Expected 2 logs, got %d", len(logs))
		}
		for _, log := range logs {
			if log.NodeID != 10 {
				t.Errorf("Expected node ID 10, got %d", log.NodeID)
			}
		}
	})

	t.Run("Empty List", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "node_id"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `upgrade_logs`")).
			WillReturnRows(rows)

		logs, err := repo.List()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(logs) != 0 {
			t.Errorf("Expected 0 logs, got %d", len(logs))
		}
	})
}

func TestUpgradeLogRepo_Clean(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewIUpgradeLogRepo()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("delete from upgrade_logs;")).
			WillReturnResult(sqlmock.NewResult(0, 5))

		err := repo.(*UpgradeLogRepo).Clean()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestUpgradeLogRepo_Create(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewIUpgradeLogRepo()

	t.Run("Success", func(t *testing.T) {
		log := &model.UpgradeLog{
			NodeID:     10,
			OldVersion: "v1.0.0",
			NewVersion: "v1.1.0",
		}

		mock.ExpectBegin()
		// Gorm 1.22+ on sqlite uses QueryRow with RETURNING id
		mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `upgrade_logs` (`created_at`,`updated_at`,`node_id`,`old_version`,`new_version`,`backup_file`) VALUES (?,?,?,?,?,?) RETURNING `id`")).
			WithArgs(AnyTime{}, AnyTime{}, 10, "v1.0.0", "v1.1.0", "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(log)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if log.ID != 1 {
			t.Errorf("Expected ID to be populated, got %d", log.ID)
		}
	})
}

func TestUpgradeLogRepo_Save(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewIUpgradeLogRepo()

	t.Run("Success", func(t *testing.T) {
		log := &model.UpgradeLog{
			NodeID:     10,
			OldVersion: "v1.0.0",
			NewVersion: "v1.1.0",
		}
		log.ID = 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `upgrade_logs` SET `created_at`=?,`updated_at`=?,`node_id`=?,`old_version`=?,`new_version`=?,`backup_file`=? WHERE `id` = ?")).
			WithArgs(AnyTime{}, AnyTime{}, 10, "v1.0.0", "v1.1.0", "", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.(*UpgradeLogRepo).Save(log)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestUpgradeLogRepo_Delete(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewIUpgradeLogRepo()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `upgrade_logs`")).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("With Options", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `upgrade_logs` WHERE node_id = ?")).
			WithArgs(uint(10)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(repo.WithByNodeID(10))

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestUpgradeLogRepo_Page(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewIUpgradeLogRepo()

	t.Run("Success", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `upgrade_logs`")).
			WillReturnRows(countRows)

		dataRows := sqlmock.NewRows([]string{"id", "node_id"}).
			AddRow(1, 10).
			AddRow(2, 10)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `upgrade_logs` LIMIT 2")).
			WillReturnRows(dataRows)

		count, logs, err := repo.Page(1, 2)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if count != 10 {
			t.Errorf("Expected count 10, got %d", count)
		}
		if len(logs) != 2 {
			t.Errorf("Expected 2 logs, got %d", len(logs))
		}
	})

	t.Run("With Options", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(5)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `upgrade_logs` WHERE node_id = ?")).
			WithArgs(uint(20)).
			WillReturnRows(countRows)

		dataRows := sqlmock.NewRows([]string{"id", "node_id"}).
			AddRow(1, 20).
			AddRow(2, 20)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `upgrade_logs` WHERE node_id = ? LIMIT 2 OFFSET 2")).
			WithArgs(uint(20)).
			WillReturnRows(dataRows)

		count, logs, err := repo.Page(2, 2, repo.WithByNodeID(20))

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if count != 5 {
			t.Errorf("Expected count 5, got %d", count)
		}
		if len(logs) != 2 {
			t.Errorf("Expected 2 logs, got %d", len(logs))
		}
	})
}

func TestUpgradeLogRepo_WithByUpgradeVersion(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewIUpgradeLogRepo()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "node_id", "old_version", "new_version"}).
			AddRow(1, 10, "v1.0.0", "v1.1.0")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `upgrade_logs` WHERE old_version = ? AND new_version = ? ORDER BY `upgrade_logs`.`id` LIMIT 1")).
			WithArgs("v1.0.0", "v1.1.0").
			WillReturnRows(rows)

		log, err := repo.Get(repo.WithByUpgradeVersion("v1.0.0", "v1.1.0"))

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if log.ID != 1 {
			t.Errorf("Expected log ID 1, got %d", log.ID)
		}
	})
}
