package repo

import (
	"regexp"
	"testing"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, got error: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		AllowGlobalUpdate: true, // Allow deletes without where clause for tests
	})

	if err != nil {
		t.Fatalf("Failed to initialize gorm db, got error: %v", err)
	}

	return gormDB, mock
}

func TestScriptRepo_Get(t *testing.T) {
	db, mock := setupTestDB(t)

	global.DB = db

	repo := NewIScriptRepo()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries` ORDER BY `script_libraries`.`id` LIMIT ?")).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "script"}).
				AddRow(1, "Test Script", "echo 'hello'"))

		script, err := repo.Get()
		assert.NoError(t, err)
		assert.Equal(t, uint(1), script.ID)
		assert.Equal(t, "Test Script", script.Name)
		assert.Equal(t, "echo 'hello'", script.Script)
	})

	t.Run("with options", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries` WHERE name LIKE ? OR description LIKE ? ORDER BY `script_libraries`.`id` LIMIT ?")).
			WithArgs("%test%", "%test%", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "script"}).
				AddRow(2, "Test Script 2", "echo 'hello 2'"))

		script, err := repo.Get(repo.WithByInfo("test"))
		assert.NoError(t, err)
		assert.Equal(t, uint(2), script.ID)
		assert.Equal(t, "Test Script 2", script.Name)
		assert.Equal(t, "echo 'hello 2'", script.Script)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries` ORDER BY `script_libraries`.`id` LIMIT ?")).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.Get()
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestScriptRepo_GetList(t *testing.T) {
	db, mock := setupTestDB(t)
	global.DB = db
	repo := NewIScriptRepo()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries`")).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "script"}).
				AddRow(1, "Test Script 1", "echo 'hello 1'").
				AddRow(2, "Test Script 2", "echo 'hello 2'"))

		scripts, err := repo.GetList()
		assert.NoError(t, err)
		assert.Len(t, scripts, 2)
		assert.Equal(t, uint(1), scripts[0].ID)
		assert.Equal(t, "Test Script 2", scripts[1].Name)
	})

	t.Run("with options", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries` WHERE name LIKE ? OR description LIKE ?")).
			WithArgs("%test%", "%test%").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "script"}).
				AddRow(1, "Test Script 1", "echo 'hello 1'"))

		scripts, err := repo.GetList(repo.WithByInfo("test"))
		assert.NoError(t, err)
		assert.Len(t, scripts, 1)
		if len(scripts) > 0 {
			assert.Equal(t, uint(1), scripts[0].ID)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries`")).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "script"}))

		scripts, err := repo.GetList()
		assert.NoError(t, err)
		assert.Len(t, scripts, 0)
	})
}

func TestScriptRepo_Page(t *testing.T) {
	db, mock := setupTestDB(t)
	global.DB = db
	repo := NewIScriptRepo()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `script_libraries`")).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries` LIMIT ?")).
			WithArgs(10).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "script"}).
				AddRow(1, "Test Script 1", "echo 'hello 1'").
				AddRow(2, "Test Script 2", "echo 'hello 2'"))

		count, scripts, err := repo.Page(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
		assert.Len(t, scripts, 2)
	})

	t.Run("with options", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `script_libraries` WHERE name LIKE ? OR description LIKE ?")).
			WithArgs("%test%", "%test%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries` WHERE name LIKE ? OR description LIKE ? LIMIT ?")).
			WithArgs("%test%", "%test%", 10).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "script"}).
				AddRow(1, "Test Script 1", "echo 'hello 1'"))

		count, scripts, err := repo.Page(1, 10, repo.WithByInfo("test"))
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
		assert.Len(t, scripts, 1)
	})

	t.Run("with offset", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `script_libraries`")).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(15))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries` LIMIT ? OFFSET ?")).
			WithArgs(10, 10).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "script"}).
				AddRow(11, "Test Script 11", "echo 'hello 11'"))

		count, scripts, err := repo.Page(2, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(15), count)
		assert.Len(t, scripts, 1)
	})
}

func TestScriptRepo_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	global.DB = db
	repo := NewIScriptRepo()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `script_libraries`")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Test Script", false, "echo 'hello'", "", false, "").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		script := &model.ScriptLibrary{
			Name:   "Test Script",
			Script: "echo 'hello'",
		}
		err := repo.Create(script)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), script.ID)
	})
}

func TestScriptRepo_Update(t *testing.T) {
	db, mock := setupTestDB(t)
	global.DB = db
	repo := NewIScriptRepo()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `script_libraries`")).
			WithArgs("Updated Name", sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		vars := map[string]interface{}{
			"name": "Updated Name",
		}
		err := repo.Update(1, vars)
		assert.NoError(t, err)
	})
}

func TestScriptRepo_UpdateGroup(t *testing.T) {
	db, mock := setupTestDB(t)
	global.DB = db
	repo := NewIScriptRepo()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `script_libraries`")).
			WithArgs(2, sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateGroup(1, 2)
		assert.NoError(t, err)
	})
}

func TestScriptRepo_Delete(t *testing.T) {
	db, mock := setupTestDB(t)
	global.DB = db
	repo := NewIScriptRepo()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `script_libraries`")).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete()
		assert.NoError(t, err)
	})

	t.Run("with options", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `script_libraries` WHERE name LIKE ? OR description LIKE ?")).
			WithArgs("%test%", "%test%").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(repo.WithByInfo("test"))
		assert.NoError(t, err)
	})
}

func TestScriptRepo_SyncAll(t *testing.T) {
	db, mock := setupTestDB(t)
	global.DB = db
	repo := NewIScriptRepo()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `script_libraries` WHERE is_system = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `script_libraries`")).
			WillReturnResult(sqlmock.NewResult(2, 2))

		mock.ExpectCommit()

		scripts := []model.ScriptLibrary{
			{Name: "System Script 1", Script: "echo 'sys1'", IsSystem: true},
			{Name: "System Script 2", Script: "echo 'sys2'", IsSystem: true},
		}

		err := repo.SyncAll(scripts)
		assert.NoError(t, err)
	})

	t.Run("delete error", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `script_libraries` WHERE is_system = ?")).
			WithArgs(1).
			WillReturnError(gorm.ErrInvalidDB)

		mock.ExpectRollback()

		scripts := []model.ScriptLibrary{
			{Name: "System Script 1", Script: "echo 'sys1'", IsSystem: true},
		}

		err := repo.SyncAll(scripts)
		assert.ErrorIs(t, err, gorm.ErrInvalidDB)
	})
}
