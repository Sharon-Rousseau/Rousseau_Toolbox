package sqlite

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRepo(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("CREATE TABLE budgets (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, created_at DATETIME)"); err != nil {
		t.Fatal(err)
	}
	repo := NewRepo(db)
	if err := repo.Create("test"); err != nil {
		t.Fatal(err)
	}
	budgets, err := repo.List()
	if err != nil || len(budgets) != 1 {
		t.Fatal("expected one budget")
	}
}
