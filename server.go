package ci_league

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
)

type Server struct {
	tmpl                *template.Template
	integrationsService IntegrationsService
}

type IntegrationsService interface {
	GetIntegrations(ctx context.Context, owner string, repo string) (TeamIntegrations, error)
}

func NewServer(tmpl *template.Template, service IntegrationsService) *Server {
	return &Server{
		tmpl:                tmpl,
		integrationsService: service,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")

	if owner == "" || repo == "" {
		http.Error(w, "Please provide both 'owner' and 'repo' query string values", http.StatusBadRequest)
		return
	}

	teamIntegrations, err := s.integrationsService.GetIntegrations(r.Context(), owner, repo)

	if err != nil {
		http.Error(w, fmt.Sprintf("Problem getting integrations %s",err), http.StatusInternalServerError)
		return
	}

	err = s.tmpl.Execute(w, teamIntegrations)

	if err != nil {
		http.Error(w, fmt.Sprintf("Problem rendering integrations %s", err), http.StatusInternalServerError)
	}
}
