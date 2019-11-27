package league

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
	GetStats(ctx context.Context, owner string, repos []string) (*TeamStats, error)
}

func NewServer(tmpl *template.Template, league League) *Server {
	return &Server{
		tmpl:   tmpl,
		league: league,
	}
}

func (this *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/integrations" {
		owner := req.URL.Query().Get("owner")
		repos := req.URL.Query()["repo"]

		if owner == "" || len(repos) == 0 {
			http.Error(res, "Please provide both 'owner' and 'repo' query string values", http.StatusBadRequest)
			return
		}

		stats, err := this.league.GetStats(req.Context(), owner, repos)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		err = this.tmpl.Execute(res, stats)

		if err != nil {
			http.Error(res, fmt.Sprintf("Problem rendering stats %this", err), http.StatusInternalServerError)
			return
		}
	}
}
