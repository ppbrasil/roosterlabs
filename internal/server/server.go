// Package server monta o http.Handler da aplicação: rotas, templates e
// arquivos estáticos. Toda rota nova entra aqui.
package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/roosterlabs/roosterlabs-engineering/web"
)

// NewHandler constrói o handler raiz com todas as rotas registradas.
func NewHandler() (http.Handler, error) {
	tmpl, err := template.ParseFS(web.Templates, "templates/*.html.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	staticFS, err := fs.Sub(web.Static, "static")
	if err != nil {
		return nil, fmt.Errorf("static fs: %w", err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, _ *http.Request) {
		// Placeholder até a landing definitiva (copy vem de roosterlabs-marketing).
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "index.html.tmpl", nil); err != nil {
			slog.Error("rendering index", "err", err)
		}
	})

	return mux, nil
}
