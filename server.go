package ci_league

import (
	"html/template"
	"net/http"
)

type Server struct {
	tmpl        *template.Template
	repo        string
	githubToken string
	idMappings  map[string]string
}

func NewServer(tmpl *template.Template, repo string, githubToken string, idMappings map[string]string) *Server {
	return &Server{tmpl: tmpl, repo: repo, githubToken: githubToken, idMappings: idMappings}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := s.tmpl.Execute(w, GetIntegrations(r.Context(), s.repo, s.githubToken, s.idMappings))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
