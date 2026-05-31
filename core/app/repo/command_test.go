package repo

import (
	"testing"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestCommandRepo_Page(t *testing.T) {
	// Backup and restore global.DB
	originalDB := global.DB
	defer func() {
		global.DB = originalDB
	}()

	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	global.DB = db

	err = db.AutoMigrate(&model.Command{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	repo := NewICommandRepo()

	// Insert dummy data
	commands := []*model.Command{
		{Name: "cmd1", GroupID: 1, Command: "echo 1", Type: "bash"},
		{Name: "cmd2", GroupID: 1, Command: "echo 2", Type: "bash"},
		{Name: "cmd3", GroupID: 2, Command: "echo 3", Type: "bash"},
	}
	for _, cmd := range commands {
		err := repo.Create(cmd)
		if err != nil {
			t.Fatalf("failed to create command: %v", err)
		}
	}

	// Test Page 1, Size 2
	count, result, err := repo.Page(1, 2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
	if result[0].Name != "cmd1" {
		t.Errorf("expected first item to be cmd1, got %s", result[0].Name)
	}

	// Test Page 2, Size 2
	count, result, err = repo.Page(2, 2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}
	if result[0].Name != "cmd3" {
		t.Errorf("expected first item to be cmd3, got %s", result[0].Name)
	}

	// Test with Option (WithByInfo)
	count, result, err = repo.Page(1, 10, repo.WithByInfo("cmd2"))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 item, got %d", len(result))
	}
	if result[0].Name != "cmd2" {
		t.Errorf("expected item to be cmd2, got %s", result[0].Name)
	}
}

func TestCommandRepo_Get(t *testing.T) {
	// Backup and restore global.DB
	originalDB := global.DB
	defer func() {
		global.DB = originalDB
	}()

	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	global.DB = db

	err = db.AutoMigrate(&model.Command{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	repo := NewICommandRepo()

	// Insert dummy data
	commands := []*model.Command{
		{Name: "cmd1", GroupID: 1, Command: "echo 1", Type: "bash"},
		{Name: "cmd2", GroupID: 1, Command: "echo 2", Type: "bash"},
	}
	for _, cmd := range commands {
		err := repo.Create(cmd)
		if err != nil {
			t.Fatalf("failed to create command: %v", err)
		}
	}

	// Test Get first item
	cmd, err := repo.Get()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if cmd.Name != "cmd1" {
		t.Errorf("expected cmd1, got %s", cmd.Name)
	}

	// Test Get with Option
	cmd, err = repo.Get(repo.WithByInfo("cmd2"))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if cmd.Name != "cmd2" {
		t.Errorf("expected cmd2, got %s", cmd.Name)
	}

	// Test Get not found
	_, err = repo.Get(repo.WithByInfo("non-existent"))
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCommandRepo_List(t *testing.T) {
	// Backup and restore global.DB
	originalDB := global.DB
	defer func() {
		global.DB = originalDB
	}()

	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	global.DB = db

	err = db.AutoMigrate(&model.Command{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	repo := NewICommandRepo()

	// Insert dummy data
	commands := []*model.Command{
		{Name: "cmd1", GroupID: 1, Command: "echo 1", Type: "bash"},
		{Name: "cmd2", GroupID: 1, Command: "echo 2", Type: "bash"},
	}
	for _, cmd := range commands {
		err := repo.Create(cmd)
		if err != nil {
			t.Fatalf("failed to create command: %v", err)
		}
	}

	// Test List without options
	list, err := repo.List()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}

	// Test List with Option
	list, err = repo.List(repo.WithByInfo("cmd2"))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 item, got %d", len(list))
	}
	if list[0].Name != "cmd2" {
		t.Errorf("expected cmd2, got %s", list[0].Name)
	}
}

func TestCommandRepo_Update(t *testing.T) {
	// Backup and restore global.DB
	originalDB := global.DB
	defer func() {
		global.DB = originalDB
	}()

	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	global.DB = db

	err = db.AutoMigrate(&model.Command{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	repo := NewICommandRepo()

	// Insert dummy data
	cmd := &model.Command{Name: "cmd1", GroupID: 1, Command: "echo 1", Type: "bash"}
	err = repo.Create(cmd)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	// Test Update
	err = repo.Update(cmd.ID, map[string]interface{}{"name": "cmd1_updated"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	updatedCmd, err := repo.Get(repo.WithByInfo("cmd1_updated"))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if updatedCmd.Name != "cmd1_updated" {
		t.Errorf("expected cmd1_updated, got %s", updatedCmd.Name)
	}
}

func TestCommandRepo_UpdateGroup(t *testing.T) {
	// Backup and restore global.DB
	originalDB := global.DB
	defer func() {
		global.DB = originalDB
	}()

	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	global.DB = db

	err = db.AutoMigrate(&model.Command{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	repo := NewICommandRepo()

	// Insert dummy data
	commands := []*model.Command{
		{Name: "cmd1", GroupID: 1, Command: "echo 1", Type: "bash"},
		{Name: "cmd2", GroupID: 1, Command: "echo 2", Type: "bash"},
		{Name: "cmd3", GroupID: 2, Command: "echo 3", Type: "bash"},
	}
	for _, cmd := range commands {
		err := repo.Create(cmd)
		if err != nil {
			t.Fatalf("failed to create command: %v", err)
		}
	}

	// Test UpdateGroup
	err = repo.UpdateGroup(1, 3)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify update
	list, err := repo.List(func(g *gorm.DB) *gorm.DB {
		return g.Where("group_id = ?", 3)
	})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

func TestCommandRepo_Delete(t *testing.T) {
	// Backup and restore global.DB
	originalDB := global.DB
	defer func() {
		global.DB = originalDB
	}()

	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	global.DB = db

	err = db.AutoMigrate(&model.Command{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	repo := NewICommandRepo()

	// Insert dummy data
	commands := []*model.Command{
		{Name: "cmd1", GroupID: 1, Command: "echo 1", Type: "bash"},
		{Name: "cmd2", GroupID: 1, Command: "echo 2", Type: "bash"},
	}
	for _, cmd := range commands {
		err := repo.Create(cmd)
		if err != nil {
			t.Fatalf("failed to create command: %v", err)
		}
	}

	// Test Delete
	err = repo.Delete(func(g *gorm.DB) *gorm.DB {
		return g.Where("name = ?", "cmd1")
	})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Verify delete
	list, err := repo.List()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 item, got %d", len(list))
	}
	if list[0].Name != "cmd2" {
		t.Errorf("expected cmd2, got %s", list[0].Name)
	}
}

func TestCommandRepo_WithByInfoEmpty(t *testing.T) {
	// Backup and restore global.DB
	originalDB := global.DB
	defer func() {
		global.DB = originalDB
	}()

	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	global.DB = db

	err = db.AutoMigrate(&model.Command{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	repo := NewICommandRepo()

	// Insert dummy data
	commands := []*model.Command{
		{Name: "cmd1", GroupID: 1, Command: "echo 1", Type: "bash"},
		{Name: "cmd2", GroupID: 1, Command: "echo 2", Type: "bash"},
	}
	for _, cmd := range commands {
		err := repo.Create(cmd)
		if err != nil {
			t.Fatalf("failed to create command: %v", err)
		}
	}

	// Test List with empty Option
	list, err := repo.List(repo.WithByInfo(""))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}
