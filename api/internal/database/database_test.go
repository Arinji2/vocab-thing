package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupDatabase(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")

	db, err := SetupDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to setup database: %v", err)
	}
	defer db.Close()
	t.Run("WAL Mode", func(t *testing.T) {
		var journalMode string
		err = db.QueryRow("PRAGMA journal_mode;").Scan(&journalMode)
		if err != nil {
			t.Fatalf("Failed to query journal mode: %v", err)
		}
		if journalMode != "wal" {
			t.Errorf("Expected journal mode to be 'wal', got '%s'", journalMode)
		}
	})

	t.Run("Foregin Keys", func(t *testing.T) {
		var foreignKeys int
		err = db.QueryRow("PRAGMA foreign_keys;").Scan(&foreignKeys)
		if err != nil {
			t.Fatalf("Failed to query foreign keys setting: %v", err)
		}
		if foreignKeys != 1 {
			t.Errorf("Expected foreign keys to be enabled (1), got %d", foreignKeys)
		}
	})
}
