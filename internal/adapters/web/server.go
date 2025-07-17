package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bootstrappedsoftware/rousseau_toolbox/internal/usecase"
)

type Server struct {
	svc *usecase.BudgetService
}

func NewServer(svc *usecase.BudgetService) *Server {
	return &Server{svc: svc}
}

func (s *Server) routes() {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/add", s.handleAdd)
	http.HandleFunc("/budget/", s.handleBudgetDelete)
}

func (s *Server) Listen(addr string) error {
	s.routes()
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	budgets, err := s.svc.ListBudgets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	component := Dashboard(budgets)
	component.Render(r.Context(), w)
}

func (s *Server) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	name := r.Form.Get("name")
	if name == "" {
		http.Error(w, "Budget name is required", http.StatusBadRequest)
		return
	}
	
	err := s.svc.CreateBudget(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return the updated budget list for HTMX
	budgets, err := s.svc.ListBudgets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	component := BudgetList(budgets)
	component.Render(r.Context(), w)
}

func (s *Server) handleBudgetDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Extract budget ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/budget/")
	budgetID := path
	
	if budgetID == "" {
		http.Error(w, "Budget ID is required", http.StatusBadRequest)
		return
	}
	
	err := s.svc.DeleteBudget(budgetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete budget: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Return the updated budget list for HTMX
	budgets, err := s.svc.ListBudgets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	component := BudgetList(budgets)
	component.Render(r.Context(), w)
}
