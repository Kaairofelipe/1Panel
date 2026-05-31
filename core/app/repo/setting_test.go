package repo

import (
	"testing"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Migrate the schema
	err = db.AutoMigrate(&model.Setting{})
	if err != nil {
		panic(err)
	}
	return db
}

func TestSettingRepo_List(t *testing.T) {
	global.DB = setupTestDB()
	repo := NewISettingRepo()

	global.DB.Create(&model.Setting{Key: "key1", Value: "value1"})
	global.DB.Create(&model.Setting{Key: "key2", Value: "value2"})

	settings, err := repo.List()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(settings) != 2 {
		t.Fatalf("expected 2 settings, got %d", len(settings))
	}

	opt := func(db *gorm.DB) *gorm.DB {
		return db.Where("key = ?", "key1")
	}

	settings, err = repo.List(opt)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(settings) != 1 {
		t.Fatalf("expected 1 setting, got %d", len(settings))
	}
	if settings[0].Key != "key1" {
		t.Errorf("expected key1, got %s", settings[0].Key)
	}
}

func TestSettingRepo_Create(t *testing.T) {
	global.DB = setupTestDB()
	repo := NewISettingRepo()

	err := repo.Create("new_key", "new_value")
	if err != nil {
		t.Fatalf("expected no error on Create, got %v", err)
	}

	var setting model.Setting
	if err := global.DB.Where("key = ?", "new_key").First(&setting).Error; err != nil {
		t.Fatalf("failed to find created setting: %v", err)
	}
	if setting.Value != "new_value" {
		t.Errorf("expected new_value, got %s", setting.Value)
	}

	// Test cache was set
	val, found := settingCache.Get("new_key")
	if !found {
		t.Errorf("expected key to be in cache")
	} else if val.(string) != "new_value" {
		t.Errorf("expected cached value new_value, got %v", val)
	}
}

func TestSettingRepo_Get(t *testing.T) {
	global.DB = setupTestDB()
	repo := NewISettingRepo()

	global.DB.Create(&model.Setting{Key: "get_key", Value: "get_value"})

	opt := func(db *gorm.DB) *gorm.DB {
		return db.Where("key = ?", "get_key")
	}

	setting, err := repo.Get(opt)
	if err != nil {
		t.Fatalf("expected no error on Get, got %v", err)
	}
	if setting.Value != "get_value" {
		t.Errorf("expected get_value, got %s", setting.Value)
	}

	// Ensure not found returns error
	optNotFound := func(db *gorm.DB) *gorm.DB {
		return db.Where("key = ?", "missing_key")
	}
	_, err = repo.Get(optNotFound)
	if err == nil {
		t.Errorf("expected error when setting not found")
	}
}

func TestSettingRepo_GetValueByKey(t *testing.T) {
	global.DB = setupTestDB()
	repo := NewISettingRepo()

	global.DB.Create(&model.Setting{Key: "val_key", Value: "val_value"})

	// Clear cache to test DB query
	settingCache.Flush()

	val, err := repo.GetValueByKey("val_key")
	if err != nil {
		t.Fatalf("expected no error on GetValueByKey, got %v", err)
	}
	if val != "val_value" {
		t.Errorf("expected val_value, got %s", val)
	}

	// Change DB directly, and get again to test cache
	global.DB.Model(&model.Setting{}).Where("key = ?", "val_key").Update("value", "new_val")

	valFromCache, err := repo.GetValueByKey("val_key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should still be val_value due to cache
	if valFromCache != "val_value" {
		t.Errorf("expected cached val_value, got %s", valFromCache)
	}

	_, err = repo.GetValueByKey("missing_val_key")
	if err == nil {
		t.Errorf("expected error when key not found")
	}
}

func TestSettingRepo_Update(t *testing.T) {
	global.DB = setupTestDB()
	repo := NewISettingRepo()

	global.DB.Create(&model.Setting{Key: "update_key", Value: "initial"})

	err := repo.Update("update_key", "updated")
	if err != nil {
		t.Fatalf("expected no error on Update, got %v", err)
	}

	var setting model.Setting
	global.DB.Where("key = ?", "update_key").First(&setting)
	if setting.Value != "updated" {
		t.Errorf("expected DB to have updated value, got %s", setting.Value)
	}

	// Cache updated
	val, found := settingCache.Get("update_key")
	if !found {
		t.Errorf("expected key in cache")
	} else if val.(string) != "updated" {
		t.Errorf("expected updated cache value, got %v", val)
	}
}

func TestSettingRepo_UpdateOrCreate(t *testing.T) {
	global.DB = setupTestDB()
	repo := NewISettingRepo()

	// Test Create
	err := repo.UpdateOrCreate("uoc_key", "value1")
	if err != nil {
		t.Fatalf("expected no error on UpdateOrCreate (create), got %v", err)
	}

	var setting model.Setting
	global.DB.Where("key = ?", "uoc_key").First(&setting)
	if setting.Value != "value1" {
		t.Errorf("expected value1, got %s", setting.Value)
	}

	// Test Update
	err = repo.UpdateOrCreate("uoc_key", "value2")
	if err != nil {
		t.Fatalf("expected no error on UpdateOrCreate (update), got %v", err)
	}

	global.DB.Where("key = ?", "uoc_key").First(&setting)
	if setting.Value != "value2" {
		t.Errorf("expected value2, got %s", setting.Value)
	}
}

func TestSettingRepo_DefaultMenu(t *testing.T) {
	global.DB = setupTestDB()
	repo := NewISettingRepo()

	global.DB.Create(&model.Setting{Key: "HideMenu", Value: "old_menus"})

	err := repo.DefaultMenu()
	if err != nil {
		t.Fatalf("expected no error on DefaultMenu, got %v", err)
	}

	var setting model.Setting
	global.DB.Where("key = ?", "HideMenu").First(&setting)

	// The helper.LoadMenus() might return an empty string or standard menu based on environment
	// We just ensure it's not "old_menus" which means it executed the update.
	if setting.Value == "old_menus" {
		t.Errorf("expected DefaultMenu to update HideMenu, it did not")
	}
}
