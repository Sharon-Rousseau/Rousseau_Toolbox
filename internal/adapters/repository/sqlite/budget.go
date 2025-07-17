package sqlite

import (
	"database/sql"
	"github.com/bootstrappedsoftware/rousseau_toolbox/internal/domain"
	"time"
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
		if err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}
