package repo

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	mock.ExpectQuery("select sqlite_version\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"sqlite_version()"}).AddRow("3.36.0"))

	dialector := sqlite.Dialector{Conn: sqlDB}
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	oldDB := global.AgentDB
	global.AgentDB = db

	t.Cleanup(func() {
		global.AgentDB = oldDB
	})

	return mock, nil
}

func TestAgentRepo_GetWebsiteSSL(t *testing.T) {
	mock, err := setupTestDB(t)
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "primary_domain", "provider"}).
			AddRow(1, "example.com", "letsencrypt")

		mock.ExpectQuery("^SELECT .* FROM `website_ssls`.*").
			WillReturnRows(rows)

		repo := NewIAgentRepo()
		ssl, err := repo.GetWebsiteSSL()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if ssl.PrimaryDomain != "example.com" {
			t.Errorf("expected PrimaryDomain 'example.com', got '%s'", ssl.PrimaryDomain)
		}
		if ssl.Provider != "letsencrypt" {
			t.Errorf("expected Provider 'letsencrypt', got '%s'", ssl.Provider)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery("^SELECT .* FROM `website_ssls`.*").
			WillReturnError(errors.New("db error"))

		repo := NewIAgentRepo()
		_, err := repo.GetWebsiteSSL()

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "db error" {
			t.Errorf("expected error 'db error', got '%v'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})
}

func TestAgentRepo_GetCA(t *testing.T) {
	mock, err := setupTestDB(t)
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "csr"}).
			AddRow(1, "TestCA", "dummy-csr")

		mock.ExpectQuery("^SELECT .* FROM `website_cas`.*").
			WillReturnRows(rows)

		repo := NewIAgentRepo()
		ca, err := repo.GetCA()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if ca.Name != "TestCA" {
			t.Errorf("expected Name 'TestCA', got '%s'", ca.Name)
		}
		if ca.CSR != "dummy-csr" {
			t.Errorf("expected CSR 'dummy-csr', got '%s'", ca.CSR)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery("^SELECT .* FROM `website_cas`.*").
			WillReturnError(errors.New("db error"))

		repo := NewIAgentRepo()
		_, err := repo.GetCA()

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "db error" {
			t.Errorf("expected error 'db error', got '%v'", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})
}
