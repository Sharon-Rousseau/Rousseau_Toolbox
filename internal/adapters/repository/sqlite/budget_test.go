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
	
	// Test Create and List
	if err := repo.Create("test"); err != nil {
		t.Fatal(err)
	}
	budgets, err := repo.List()
	if err != nil || len(budgets) != 1 {
		t.Fatal("expected one budget")
	}
	
	// Test that ID is converted to string
	budget := budgets[0]
	if budget.ID == "" {
		t.Fatal("expected budget ID to be a non-empty string")
	}
	if budget.Name != "test" {
		t.Fatal("expected budget name to be 'test'")
	}
	
	// Test Delete
	if err := repo.Delete(budget.ID); err != nil {
		t.Fatal(err)
	}
	budgets, err = repo.List()
	if err != nil || len(budgets) != 0 {
		t.Fatal("expected zero budgets after delete")
	}
	
	// Test Delete with invalid ID
	if err := repo.Delete("invalid"); err == nil {
		t.Fatal("expected error when deleting with invalid ID")
	}
	
	// Test Delete with non-existent ID
	if err := repo.Delete("999"); err == nil {
		t.Fatal("expected error when deleting non-existent budget")
	}
}
