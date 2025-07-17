package usecase

import "github.com/bootstrappedsoftware/rousseau_toolbox/internal/domain"

// BudgetRepository provides access to the budget storage.
type BudgetRepository interface {
	Create(name string) error
	List() ([]domain.Budget, error)
	Delete(id string) error
}

type BudgetService struct {
	repo BudgetRepository
}

func NewBudgetService(r BudgetRepository) *BudgetService {
	return &BudgetService{repo: r}
}

func (s *BudgetService) CreateBudget(name string) error {
	return s.repo.Create(name)
}

func (s *BudgetService) ListBudgets() ([]domain.Budget, error) {
	return s.repo.List()
}

func (s *BudgetService) DeleteBudget(id string) error {
	return s.repo.Delete(id)
}
