package repo

import (
	"testing"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.LoginLog{}, &model.OperationLog{})
	if err != nil {
		t.Fatalf("failed to auto migrate: %v", err)
	}

	return db
}

func TestLogRepo_CreateLoginLog(t *testing.T) {
	db := setupTestDB(t)
	originalDB := global.DB
	global.DB = db
	defer func() {
		global.DB = originalDB
	}()

	repo := NewILogRepo()

	logEntry := &model.LoginLog{
		IP:      "192.168.1.1",
		Address: "Local Network",
		Agent:   "Mozilla/5.0",
		Status:  "Success",
		Message: "Login successful",
	}

	err := repo.CreateLoginLog(logEntry)
	if err != nil {
		t.Fatalf("CreateLoginLog failed: %v", err)
	}

	var count int64
	db.Model(&model.LoginLog{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 login log, got %d", count)
	}

	var fetchedLog model.LoginLog
	db.First(&fetchedLog)
	if fetchedLog.IP != "192.168.1.1" {
		t.Errorf("Expected IP 192.168.1.1, got %s", fetchedLog.IP)
	}
}

func TestLogRepo_CleanLogin(t *testing.T) {
	db := setupTestDB(t)
	originalDB := global.DB
	global.DB = db
	defer func() {
		global.DB = originalDB
	}()

	repo := NewILogRepo()

	// Insert some test data
	db.Create(&model.LoginLog{IP: "1.1.1.1"})
	db.Create(&model.LoginLog{IP: "2.2.2.2"})

	err := repo.CleanLogin()
	if err != nil {
		t.Fatalf("CleanLogin failed: %v", err)
	}

	var count int64
	db.Model(&model.LoginLog{}).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 login logs after clean, got %d", count)
	}
}

func TestLogRepo_PageLoginLog(t *testing.T) {
	db := setupTestDB(t)
	originalDB := global.DB
	global.DB = db
	defer func() {
		global.DB = originalDB
	}()

	repo := NewILogRepo()

	// Insert test data
	db.Create(&model.LoginLog{IP: "1.1.1.1", Status: "Success"})
	db.Create(&model.LoginLog{IP: "2.2.2.2", Status: "Failed"})
	db.Create(&model.LoginLog{IP: "3.3.3.3", Status: "Success"})

	// Test pagination without options
	count, logs, err := repo.PageLoginLog(1, 2)
	if err != nil {
		t.Fatalf("PageLoginLog failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs on first page, got %d", len(logs))
	}

	// Test with WithByStatus option
	var logRepo *LogRepo
	logRepo = repo.(*LogRepo)
	count, logs, err = repo.PageLoginLog(1, 10, logRepo.WithByStatus("Success"))
	if err != nil {
		t.Fatalf("PageLoginLog with WithByStatus failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2 for 'Success' status, got %d", count)
	}
}

func TestLogRepo_CreateOperationLog(t *testing.T) {
	db := setupTestDB(t)
	originalDB := global.DB
	global.DB = db
	defer func() {
		global.DB = originalDB
	}()

	repo := NewILogRepo()

	logEntry := &model.OperationLog{
		IP:       "10.0.0.1",
		Source:   "Admin",
		Method:   "GET",
		Path:     "/api/test",
		DetailEN: "Test operation",
	}

	err := repo.CreateOperationLog(logEntry)
	if err != nil {
		t.Fatalf("CreateOperationLog failed: %v", err)
	}

	var count int64
	db.Model(&model.OperationLog{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 operation log, got %d", count)
	}
}

func TestLogRepo_CleanOperation(t *testing.T) {
	db := setupTestDB(t)
	originalDB := global.DB
	global.DB = db
	defer func() {
		global.DB = originalDB
	}()

	repo := NewILogRepo()

	db.Create(&model.OperationLog{IP: "1.1.1.1"})

	err := repo.CleanOperation()
	if err != nil {
		t.Fatalf("CleanOperation failed: %v", err)
	}

	var count int64
	db.Model(&model.OperationLog{}).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 operation logs after clean, got %d", count)
	}
}

func TestLogRepo_PageOperationLog(t *testing.T) {
	db := setupTestDB(t)
	originalDB := global.DB
	global.DB = db
	defer func() {
		global.DB = originalDB
	}()

	repo := NewILogRepo()

	db.Create(&model.OperationLog{Source: "Admin", DetailEN: "Login"})
	db.Create(&model.OperationLog{Source: "User", DetailEN: "Logout"})

	count, logs, err := repo.PageOperationLog(1, 10)
	if err != nil {
		t.Fatalf("PageOperationLog failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}
}

func TestLogRepo_Options(t *testing.T) {
	db := setupTestDB(t)
	repo := &LogRepo{}

	db.Create(&model.OperationLog{Source: "Admin", IP: "192.168.1.1", DetailZH: "登录", DetailEN: "Login"})
	db.Create(&model.OperationLog{Source: "System", IP: "10.0.0.1", DetailZH: "退出", DetailEN: "Logout"})

	tests := []struct{
		name string
		opt global.DBOption
		expectedCount int64
	}{
		{
			name: "WithBySource",
			opt: repo.WithBySource("Admin"),
			expectedCount: 1,
		},
		{
			name: "WithBySource empty",
			opt: repo.WithBySource(""),
			expectedCount: 2,
		},
		{
			name: "WithByIP",
			opt: repo.WithByIP("192.168.1"),
			expectedCount: 1,
		},
		{
			name: "WithByLikeOperation ZH",
			opt: repo.WithByLikeOperation("登录"),
			expectedCount: 1,
		},
		{
			name: "WithByLikeOperation EN",
			opt: repo.WithByLikeOperation("Logout"),
			expectedCount: 1,
		},
		{
			name: "WithByLikeOperation empty",
			opt: repo.WithByLikeOperation(""),
			expectedCount: 2,
		},
		{
			name: "WithByStatus empty",
			opt: repo.WithByStatus(""),
			expectedCount: 2, // No change to query
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := db.Model(&model.OperationLog{})
			query = tt.opt(query)

			var count int64
			query.Count(&count)

			if count != tt.expectedCount {
				t.Errorf("Expected count %d, got %d", tt.expectedCount, count)
			}
		})
	}
}
