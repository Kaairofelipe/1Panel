package repo_test

import (
	"testing"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/app/repo"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	// Create a temporary database for testing using in-memory SQLite
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize test db: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&model.Command{})
	if err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}

	global.DB = db

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})
}

func TestCommandRepo_CreateAndGet(t *testing.T) {
	setupTestDB(t)
	r := repo.NewICommandRepo()

	cmd := &model.Command{
		Name:    "Test Command",
		Command: "echo 'hello'",
		GroupID: 1,
		Type:    "shell",
	}

	// Test Create
	err := r.Create(cmd)
	if err != nil {
		t.Fatalf("Failed to create command: %v", err)
	}
	if cmd.ID == 0 {
		t.Errorf("Expected command ID to be set, got 0")
	}

	// Test Get with WithByInfo
	fetchedCmd, err := r.Get(r.WithByInfo("Test Command"))
	if err != nil {
		t.Fatalf("Failed to get command: %v", err)
	}

	if fetchedCmd.Name != cmd.Name {
		t.Errorf("Expected command name %s, got %s", cmd.Name, fetchedCmd.Name)
	}

	// Test Get with no matched WithByInfo
	_, err = r.Get(r.WithByInfo("Non existent"))
	if err == nil {
		t.Errorf("Expected error for non-existent command, got nil")
	}
}

func TestCommandRepo_Update(t *testing.T) {
	setupTestDB(t)
	r := repo.NewICommandRepo()

	cmd := &model.Command{
		Name:    "Test Command",
		Command: "echo 'hello'",
		GroupID: 1,
		Type:    "shell",
	}

	err := r.Create(cmd)
	if err != nil {
		t.Fatalf("Failed to create command: %v", err)
	}

	vars := map[string]interface{}{
		"name": "Updated Command",
	}

	err = r.Update(cmd.ID, vars)
	if err != nil {
		t.Fatalf("Failed to update command: %v", err)
	}

	fetchedCmd, err := r.Get(r.WithByInfo("Updated Command"))
	if err != nil {
		t.Fatalf("Failed to get command: %v", err)
	}

	if fetchedCmd.Name != "Updated Command" {
		t.Errorf("Expected updated name %s, got %s", "Updated Command", fetchedCmd.Name)
	}
}

func TestCommandRepo_UpdateGroup(t *testing.T) {
	setupTestDB(t)
	r := repo.NewICommandRepo()

	cmd1 := &model.Command{
		Name:    "Test Command 1",
		Command: "echo 'hello'",
		GroupID: 1,
		Type:    "shell",
	}
	cmd2 := &model.Command{
		Name:    "Test Command 2",
		Command: "echo 'world'",
		GroupID: 1,
		Type:    "shell",
	}

	r.Create(cmd1)
	r.Create(cmd2)

	err := r.UpdateGroup(1, 2)
	if err != nil {
		t.Fatalf("Failed to update group: %v", err)
	}

	fetchedCmd1, _ := r.Get(r.WithByInfo("Test Command 1"))
	if fetchedCmd1.GroupID != 2 {
		t.Errorf("Expected GroupID 2, got %d", fetchedCmd1.GroupID)
	}

	fetchedCmd2, _ := r.Get(r.WithByInfo("Test Command 2"))
	if fetchedCmd2.GroupID != 2 {
		t.Errorf("Expected GroupID 2, got %d", fetchedCmd2.GroupID)
	}
}

func TestCommandRepo_List(t *testing.T) {
	setupTestDB(t)
	r := repo.NewICommandRepo()

	cmd1 := &model.Command{
		Name:    "Test Command 1",
		Command: "echo 'hello'",
		GroupID: 1,
		Type:    "shell",
	}
	cmd2 := &model.Command{
		Name:    "Test Command 2",
		Command: "echo 'world'",
		GroupID: 2,
		Type:    "shell",
	}

	r.Create(cmd1)
	r.Create(cmd2)

	cmds, err := r.List()
	if err != nil {
		t.Fatalf("Failed to list commands: %v", err)
	}

	if len(cmds) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(cmds))
	}

	cmds, err = r.List(r.WithByInfo("Test Command 1"))
	if err != nil {
		t.Fatalf("Failed to list commands with filter: %v", err)
	}

	if len(cmds) != 1 {
		t.Errorf("Expected 1 command with filter, got %d", len(cmds))
	}
}

func TestCommandRepo_Page(t *testing.T) {
	setupTestDB(t)
	r := repo.NewICommandRepo()

	for i := 0; i < 15; i++ {
		cmd := &model.Command{
			Name:    "Test Command",
			Command: "echo 'hello'",
			GroupID: 1,
			Type:    "shell",
		}
		r.Create(cmd)
	}

	// The interface is Page(limit, offset int) but the implementation actually expects page and size.
	// So we pass 1 for page (limit), and 10 for size (offset).
	count, cmds, err := r.Page(1, 10)
	if err != nil {
		t.Fatalf("Failed to page commands: %v", err)
	}

	if count != 15 {
		t.Errorf("Expected total count 15, got %d", count)
	}

	if len(cmds) != 10 {
		t.Errorf("Expected 10 commands in first page, got %d", len(cmds))
	}

	count2, cmds2, err2 := r.Page(2, 10)
	if err2 != nil {
		t.Fatalf("Failed to page commands: %v", err2)
	}

	if count2 != 15 {
		t.Errorf("Expected total count 15, got %d", count2)
	}

	if len(cmds2) != 5 {
		t.Errorf("Expected 5 commands in second page, got %d", len(cmds2))
	}
}

func TestCommandRepo_Delete(t *testing.T) {
	setupTestDB(t)
	r := repo.NewICommandRepo()

	cmd := &model.Command{
		Name:    "Test Command",
		Command: "echo 'hello'",
		GroupID: 1,
		Type:    "shell",
	}

	r.Create(cmd)

	err := r.Delete(r.WithByInfo("Test Command"))
	if err != nil {
		t.Fatalf("Failed to delete command: %v", err)
	}

	cmds, _ := r.List()
	if len(cmds) != 0 {
		t.Errorf("Expected 0 commands after deletion, got %d", len(cmds))
	}
}
