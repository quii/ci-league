package ci_league

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
)

type Server struct {
	tmpl   *template.Template
	league League
}

type League interface {
	GetStats(ctx context.Context, owner string, repos []string) (TeamStats, error)
}

func NewServer(tmpl *template.Template, league League) *Server {
	return &Server{
		tmpl:   tmpl,
		league: league,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/integrations" {
		owner := r.URL.Query().Get("owner")
		repos := r.URL.Query()["repo"]

		if owner == "" || len(repos) == 0 {
			http.Error(w, "Please provide both 'owner' and 'repo' query string values", http.StatusBadRequest)
			return
		}

		stats, err := s.league.GetStats(r.Context(), owner, repos)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.tmpl.Execute(w, stats)

		if err != nil {
			http.Error(w, fmt.Sprintf("Problem rendering stats %s", err), http.StatusInternalServerError)
			return
		}
	}
}
