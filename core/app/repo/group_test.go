package repo

import (
	"regexp"
	"testing"
	"time"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, got error: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm db, got error: %v", err)
	}

	global.DB = db
	return db, mock
}

func TestGroupRepo_Get(t *testing.T) {
	_, mock := setupTestDB(t)

	repo := NewIGroupRepo()

	now := time.Now()
	expectedGroup := model.Group{
		BaseModel: model.BaseModel{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		IsDefault: true,
		Name:      "Test Group",
		Type:      "website",
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "is_default", "name", "type"}).
		AddRow(expectedGroup.ID, expectedGroup.CreatedAt, expectedGroup.UpdatedAt, expectedGroup.IsDefault, expectedGroup.Name, expectedGroup.Type)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups" ORDER BY "groups"."id" LIMIT $1`)).
        WithArgs(1).
		WillReturnRows(rows)

	group, err := repo.Get()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if group.ID != expectedGroup.ID {
		t.Errorf("Expected ID %d, got %d", expectedGroup.ID, group.ID)
	}
	if group.Name != expectedGroup.Name {
		t.Errorf("Expected Name %s, got %s", expectedGroup.Name, group.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGroupRepo_Get_WithOpts(t *testing.T) {
	_, mock := setupTestDB(t)

	repo := NewIGroupRepo()

	now := time.Now()
	expectedGroup := model.Group{
		BaseModel: model.BaseModel{
			ID:        2,
			CreatedAt: now,
			UpdatedAt: now,
		},
		IsDefault: true,
		Name:      "Default Website Group",
		Type:      "website",
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "is_default", "name", "type"}).
		AddRow(expectedGroup.ID, expectedGroup.CreatedAt, expectedGroup.UpdatedAt, expectedGroup.IsDefault, expectedGroup.Name, expectedGroup.Type)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups" WHERE is_default = $1 ORDER BY "groups"."id" LIMIT $2`)).
		WithArgs(true, 1).
		WillReturnRows(rows)

	group, err := repo.Get(repo.WithByDefault(true))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if group.ID != expectedGroup.ID {
		t.Errorf("Expected ID %d, got %d", expectedGroup.ID, group.ID)
	}
	if !group.IsDefault {
		t.Errorf("Expected IsDefault to be true, got %v", group.IsDefault)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGroupRepo_Get_NotFound(t *testing.T) {
	_, mock := setupTestDB(t)

	repo := NewIGroupRepo()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups" ORDER BY "groups"."id" LIMIT $1`)).
        WithArgs(1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.Get()
	if err != gorm.ErrRecordNotFound {
		t.Fatalf("Expected gorm.ErrRecordNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGroupRepo_GetList(t *testing.T) {
	_, mock := setupTestDB(t)

	repo := NewIGroupRepo()

	now := time.Now()
	expectedGroups := []model.Group{
		{
			BaseModel: model.BaseModel{ID: 1, CreatedAt: now, UpdatedAt: now},
			IsDefault: true, Name: "Group 1", Type: "website",
		},
		{
			BaseModel: model.BaseModel{ID: 2, CreatedAt: now, UpdatedAt: now},
			IsDefault: false, Name: "Group 2", Type: "database",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "is_default", "name", "type"}).
		AddRow(expectedGroups[0].ID, expectedGroups[0].CreatedAt, expectedGroups[0].UpdatedAt, expectedGroups[0].IsDefault, expectedGroups[0].Name, expectedGroups[0].Type).
		AddRow(expectedGroups[1].ID, expectedGroups[1].CreatedAt, expectedGroups[1].UpdatedAt, expectedGroups[1].IsDefault, expectedGroups[1].Name, expectedGroups[1].Type)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups"`)).
		WillReturnRows(rows)

	groups, err := repo.GetList()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(groups) != 2 {
		t.Fatalf("Expected 2 groups, got %d", len(groups))
	}
	if groups[0].ID != expectedGroups[0].ID || groups[1].ID != expectedGroups[1].ID {
		t.Errorf("Returned groups do not match expected")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGroupRepo_GetList_WithOpts(t *testing.T) {
	_, mock := setupTestDB(t)

	repo := NewIGroupRepo()

	now := time.Now()
	expectedGroups := []model.Group{
		{
			BaseModel: model.BaseModel{ID: 1, CreatedAt: now, UpdatedAt: now},
			IsDefault: true, Name: "Group 1", Type: "website",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "is_default", "name", "type"}).
		AddRow(expectedGroups[0].ID, expectedGroups[0].CreatedAt, expectedGroups[0].UpdatedAt, expectedGroups[0].IsDefault, expectedGroups[0].Name, expectedGroups[0].Type)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "groups" WHERE is_default = $1`)).
		WithArgs(true).
		WillReturnRows(rows)

	groups, err := repo.GetList(repo.WithByDefault(true))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(groups))
	}
	if groups[0].ID != expectedGroups[0].ID {
		t.Errorf("Returned groups do not match expected")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGroupRepo_Create(t *testing.T) {
	_, mock := setupTestDB(t)

	repo := NewIGroupRepo()

	group := &model.Group{
		IsDefault: true,
		Name:      "New Group",
		Type:      "website",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "groups" ("created_at","updated_at","is_default","name","type") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), group.IsDefault, group.Name, group.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(group)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if group.ID != 1 {
		t.Errorf("Expected group ID to be 1, got %d", group.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGroupRepo_Update(t *testing.T) {
	_, mock := setupTestDB(t)

	repo := NewIGroupRepo()

	vars := map[string]interface{}{
		"name": "Updated Group",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "groups" SET "name"=$1,"updated_at"=$2 WHERE id = $3`)).
		WithArgs(vars["name"], sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(1, vars)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGroupRepo_Delete(t *testing.T) {
	_, mock := setupTestDB(t)

	repo := NewIGroupRepo()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "groups" WHERE is_default = $1`)).
		WithArgs(true).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(repo.WithByDefault(true))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGroupRepo_CancelDefault(t *testing.T) {
	_, mock := setupTestDB(t)

	repo := NewIGroupRepo()

	groupType := "website"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "groups" SET "is_default"=$1,"updated_at"=$2 WHERE is_default = $3 AND type = $4`)).
		WithArgs(0, sqlmock.AnyArg(), 1, groupType).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.CancelDefault(groupType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
