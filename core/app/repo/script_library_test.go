package repo

import (
	"regexp"
	"testing"

	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	// SQLite dialect needs these queries to be mocked during gorm.Open
	mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{"sqlite_version()"}).AddRow("3.38.5"))

	dialector := sqlite.Dialector{
		Conn: sqlDB,
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	global.DB = gormDB
	return mock, nil
}

func TestScriptRepo_GetList(t *testing.T) {
	mock, err := setupTestDB()
	assert.NoError(t, err)

	repo := &ScriptRepo{}

	t.Run("GetList successfully", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries`")).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "script"}).
				AddRow(1, "Test Script 1", "echo 'hello'").
				AddRow(2, "Test Script 2", "echo 'world'"))

		scripts, err := repo.GetList()

		assert.NoError(t, err)
		assert.Len(t, scripts, 2)
		if len(scripts) > 0 {
			assert.Equal(t, uint(1), scripts[0].ID)
			assert.Equal(t, "Test Script 1", scripts[0].Name)
		}
		if len(scripts) > 1 {
			assert.Equal(t, uint(2), scripts[1].ID)
			assert.Equal(t, "Test Script 2", scripts[1].Name)
		}

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetList with DBOption", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries` WHERE id = ?")).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "script"}).
				AddRow(1, "test script", "echo 'hello'"))

		dummyOpt := func(db *gorm.DB) *gorm.DB {
			return db.Where("id = ?", 1)
		}

		scripts, err := repo.GetList(dummyOpt)

		assert.NoError(t, err)
		assert.Len(t, scripts, 1)
		if len(scripts) > 0 {
			assert.Equal(t, uint(1), scripts[0].ID)
			assert.Equal(t, "test script", scripts[0].Name)
		}

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetList with error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `script_libraries`")).
			WillReturnError(gorm.ErrInvalidDB)

		scripts, err := repo.GetList()

		assert.ErrorIs(t, err, gorm.ErrInvalidDB)
		assert.Empty(t, scripts)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
