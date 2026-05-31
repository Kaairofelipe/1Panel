Title: "🧪 Add tests for BackupRepo.Get and BackupRepo.Page"

🎯 **What:** The testing gap addressed
The `core/app/repo/backup.go` file lacked tests for the DB methods `Get` and `Page`. These methods perform DB queries via GORM that are crucial for backup operations and retrieving pages of backups.

📊 **Coverage:** What scenarios are now tested
- `BackupRepo.Get`: success, record not found error, and with global DB options applied.
- `BackupRepo.Page`: success for first and second pages to verify pagination limits and offsets, database error during the count phase, and database error during the fetch phase.
- Replaced actual database interactions by mocking `gorm.DB` using `github.com/DATA-DOG/go-sqlmock`.

✨ **Result:** The improvement in test coverage
The coverage for `BackupRepo` logic has significantly increased, ensuring no DB regressions and improving code reliability.
