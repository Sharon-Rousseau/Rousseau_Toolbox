package web

import (
	"html/template"
	"net/http"

	"github.com/bootstrappedsoftware/rousseau_toolbox/internal/usecase"
)

type Server struct {
	svc       *usecase.BudgetService
	templates *template.Template
}

func NewServer(svc *usecase.BudgetService) *Server {
	t := template.Must(template.ParseFS(templateFS, "dashboard.html"))
	return &Server{svc: svc, templates: t}
}

func (s *Server) routes() {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/add", s.handleAdd)
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
	s.templates.ExecuteTemplate(w, "dashboard.html", budgets)
}

func (s *Server) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		name := r.Form.Get("name")
		_ = s.svc.CreateBudget(name)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	s.templates.ExecuteTemplate(w, "dashboard.html", nil)
}
