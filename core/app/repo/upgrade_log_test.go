package repo

import (
	"testing"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.UpgradeLog{})
	assert.NoError(t, err)

	return db
}

func TestUpgradeLogRepo_CreateAndGet(t *testing.T) {
	oldDB := global.DB
	defer func() { global.DB = oldDB }()

	db := setupTestDB(t)
	global.DB = db

	repo := NewIUpgradeLogRepo()

	log := &model.UpgradeLog{
		NodeID:     1,
		OldVersion: "v1.0",
		NewVersion: "v1.1",
		BackupFile: "backup1.zip",
	}

	err := repo.Create(log)
	assert.NoError(t, err)
	assert.NotZero(t, log.ID)

	retrievedLog, err := repo.Get(repo.WithByNodeID(1))
	assert.NoError(t, err)
	assert.Equal(t, log.OldVersion, retrievedLog.OldVersion)
	assert.Equal(t, log.NewVersion, retrievedLog.NewVersion)

	retrievedLog2, err := repo.Get(repo.WithByUpgradeVersion("v1.0", "v1.1"))
	assert.NoError(t, err)
	assert.Equal(t, log.NodeID, retrievedLog2.NodeID)
}

func TestUpgradeLogRepo_ListAndPage(t *testing.T) {
	oldDB := global.DB
	defer func() { global.DB = oldDB }()

	db := setupTestDB(t)
	global.DB = db

	repo := NewIUpgradeLogRepo()

	log1 := &model.UpgradeLog{
		NodeID:     1,
		OldVersion: "v1.0",
		NewVersion: "v1.1",
	}
	log2 := &model.UpgradeLog{
		NodeID:     2,
		OldVersion: "v1.1",
		NewVersion: "v1.2",
	}

	repo.Create(log1)
	repo.Create(log2)

	logs, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, logs, 2)

	count, pagedLogs, err := repo.Page(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
	assert.Len(t, pagedLogs, 1)
	assert.Equal(t, log1.OldVersion, pagedLogs[0].OldVersion)

	count2, pagedLogs2, err := repo.Page(2, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count2)
	assert.Len(t, pagedLogs2, 1)
	assert.Equal(t, log2.OldVersion, pagedLogs2[0].OldVersion)
}

func TestUpgradeLogRepo_Delete(t *testing.T) {
	oldDB := global.DB
	defer func() { global.DB = oldDB }()

	db := setupTestDB(t)
	global.DB = db

	repo := NewIUpgradeLogRepo()

	log := &model.UpgradeLog{
		NodeID:     1,
		OldVersion: "v1.0",
		NewVersion: "v1.1",
	}

	repo.Create(log)

	err := repo.Delete(repo.WithByNodeID(1))
	assert.NoError(t, err)

	_, err = repo.Get(repo.WithByNodeID(1))
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
