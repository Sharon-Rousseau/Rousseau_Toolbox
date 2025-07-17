package database

import (
	"database/sql"
	"io/fs"
	"os"
	"sort"
)

// RunMigrations executes all SQL files in the migrations directory.
func RunMigrations(db *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := fs.ReadFile(os.DirFS(dir), e.Name())
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(data)); err != nil {
			return err
		}
	}
	return nil
}
