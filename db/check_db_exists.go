package db

import "os"

func AdjustDBPath(dbPath string) (string, bool) {
	var exists = true

	if _, err := os.Stat(dbPath); err != nil {
		exists = false
	}

	if !exists {
		dbPath = dbPath + "?mode=rwc"
	}

	return dbPath, exists
}
