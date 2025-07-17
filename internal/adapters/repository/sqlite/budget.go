package sqlite

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
	
	"github.com/bootstrappedsoftware/rousseau_toolbox/internal/domain"
)

// Repo implements the usecase.BudgetRepository interface using SQLite.
type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

func (r *Repo) Create(name string) error {
	_, err := r.DB.Exec(`INSERT INTO budgets (name, created_at) VALUES (?, ?)`, name, time.Now())
	return err
}

func (r *Repo) List() ([]domain.Budget, error) {
	rows, err := r.DB.Query(`SELECT id, name, created_at FROM budgets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []domain.Budget
	for rows.Next() {
		var b domain.Budget
		var id int
		if err := rows.Scan(&id, &b.Name, &b.CreatedAt); err != nil {
			return nil, err
		}
		b.ID = strconv.Itoa(id)
		budgets = append(budgets, b)
	}
	return budgets, nil
}

func (r *Repo) Delete(id string) error {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid budget ID: %v", err)
	}
	
	result, err := r.DB.Exec(`DELETE FROM budgets WHERE id = ?`, intID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("budget with ID %s not found", id)
	}
	
	return nil
}
