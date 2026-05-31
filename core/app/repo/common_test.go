package repo

import (
	"strings"
	"testing"

	"github.com/1Panel-dev/1Panel/core/constant"
	"github.com/1Panel-dev/1Panel/core/utils/re"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DummyModel struct {
	ID        uint
	GroupID   uint
	Name      string
	Type      string
	Addr      string
	Key       string
	Status    string
	Node      string
	CreatedAt string
}

func init() {
	re.Init()
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
	db.AutoMigrate(&DummyModel{})
	return db
}

func getSQL(db *gorm.DB) string {
	stmt := db.Session(&gorm.Session{DryRun: true}).Find(&DummyModel{}).Statement
	return stmt.SQL.String()
}

func checkSQLContains(t *testing.T, db *gorm.DB, expected string) {
	t.Helper()
	sql := getSQL(db)
	if !strings.Contains(sql, expected) {
		t.Errorf("Expected SQL to contain %q, but got %q", expected, sql)
	}
}

func TestWithByID(t *testing.T) {
	db := setupTestDB(t)
	db = WithByID(1)(db)
	checkSQLContains(t, db, "id = ?")
}

func TestWithByGroupID(t *testing.T) {
	db := setupTestDB(t)
	db = WithByGroupID(2)(db)
	checkSQLContains(t, db, "group_id = ?")
}

func TestWithByIDs(t *testing.T) {
	db := setupTestDB(t)
	db = WithByIDs([]uint{1, 2, 3})(db)
	checkSQLContains(t, db, "id in (?,?,?)")
}

func TestWithByName(t *testing.T) {
	db := setupTestDB(t)
	db = WithByName("test")(db)
	checkSQLContains(t, db, "`name` = ?")
}

func TestWithoutByName(t *testing.T) {
	db := setupTestDB(t)
	db = WithoutByName("test")(db)
	checkSQLContains(t, db, "`name` != ?")
}

func TestWithByType(t *testing.T) {
	db := setupTestDB(t)
	db = WithByType("typeA")(db)
	checkSQLContains(t, db, "`type` = ?")
}

func TestWithByAddr(t *testing.T) {
	db := setupTestDB(t)
	db = WithByAddr("127.0.0.1")(db)
	checkSQLContains(t, db, "addr = ?")
}

func TestWithByKey(t *testing.T) {
	db := setupTestDB(t)
	db = WithByKey("keyA")(db)
	checkSQLContains(t, db, "key = ?")
}

func TestWithByStatus(t *testing.T) {
	db := setupTestDB(t)
	db = WithByStatus("active")(db)
	checkSQLContains(t, db, "status = ?")
}

func TestWithByNode(t *testing.T) {
	db := setupTestDB(t)
	db = WithByNode("node1")(db)
	checkSQLContains(t, db, "node = ?")
}

func TestWithOrderDesc(t *testing.T) {
	db := setupTestDB(t)
	db = WithOrderDesc("id")(db)
	checkSQLContains(t, db, "ORDER BY id desc")
}

func TestWithOrderAsc(t *testing.T) {
	db := setupTestDB(t)
	db = WithOrderAsc("id")(db)
	checkSQLContains(t, db, "ORDER BY id asc")
}

func TestWithOrderRuleBy(t *testing.T) {
	tests := []struct {
		name     string
		orderBy  string
		order    string
		expected string
	}{
		{
			name:     "valid field, asc",
			orderBy:  "id",
			order:    constant.OrderAsc,
			expected: "ORDER BY id asc",
		},
		{
			name:     "valid field, desc",
			orderBy:  "name",
			order:    constant.OrderDesc,
			expected: "ORDER BY name desc",
		},
		{
			name:     "camelCase createdAt mapped to snake_case",
			orderBy:  "createdAt",
			order:    constant.OrderDesc,
			expected: "ORDER BY created_at desc",
		},
		{
			name:     "invalid field defaults to created_at",
			orderBy:  "invalid field;",
			order:    constant.OrderAsc,
			expected: "ORDER BY created_at asc",
		},
		{
			name:     "invalid order defaults to desc",
			orderBy:  "id",
			order:    "invalid_order",
			expected: "ORDER BY created_at desc",
		},
		{
			name:     "invalid field and invalid order default to created_at desc",
			orderBy:  "invalid field;",
			order:    "invalid_order",
			expected: "ORDER BY created_at desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			db = WithOrderRuleBy(tt.orderBy, tt.order)(db)
			checkSQLContains(t, db, tt.expected)
		})
	}
}
