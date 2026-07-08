// Package server monta o http.Handler da aplicação: rotas, templates e
// arquivos estáticos. Toda rota nova entra aqui.
package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/roosterlabs/roosterlabs-engineering/internal/leads"
	"github.com/roosterlabs/roosterlabs-engineering/web"
)

const sessionCookieName = "rlsid"

type app struct {
	tmpl         *template.Template
	store        leads.Store
	baseURL      string
	contactEmail string
}

// pageData é o único tipo de dados passado a qualquer template desta
// aplicação — landing e todos os passos do form (épico 002, T8: antes
// existiam dois tipos quase idênticos, pageData e formTemplateData, só por
// causa da landing embutir o passo 1 do form via {{template "form_step_1_*" .}}).
// BaseURL e OGImageURL só são usados pelos templates de landing; os passos
// do form ignoram esses dois campos.
type pageData struct {
	Lang         string
	PagePath     string
	BaseURL      string
	OGImageURL   string
	ContactEmail string
	UTM          leads.UTM
	Error        string
	Values       map[string]string
}

// Config parametriza a construção do handler para ambientes diferentes.
type Config struct {
	Store        leads.Store
	BaseURL      string
	ContactEmail string
}

// NewHandler constrói o handler raiz com todas as rotas registradas.
func NewHandler(cfg Config) (http.Handler, error) {
	tmpl, err := template.ParseFS(web.Templates, "templates/*.html.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	store := cfg.Store
	if store == nil {
		store = leads.NewMemoryStore()
	}

	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = "https://roosterlabs.com.br"
	}

	contactEmail := strings.TrimSpace(cfg.ContactEmail)
	if contactEmail == "" {
		contactEmail = "contact@roosterlabs.com.br"
	}

	a := &app{
		tmpl:         tmpl,
		store:        store,
		baseURL:      strings.TrimRight(baseURL, "/"),
		contactEmail: contactEmail,
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

	mux.HandleFunc("GET /{$}", a.handleLandingPT)
	mux.HandleFunc("GET /en/", a.handleLandingEN)
	mux.HandleFunc("POST /event/view", a.handleViewEvent)
	mux.HandleFunc("POST /form/{step}", a.handleFormStep)

	return mux, nil
}

func (a *app) handleLandingPT(w http.ResponseWriter, r *http.Request) {
	a.renderLanding(w, r, "pt-BR", "index.html.tmpl", "/")
}

func (a *app) handleLandingEN(w http.ResponseWriter, r *http.Request) {
	a.renderLanding(w, r, "en", "index.en.html.tmpl", "/en/")
}

func (a *app) renderLanding(w http.ResponseWriter, r *http.Request, lang, templateName, pagePath string) {
	// ensureToken tem que rodar antes de qualquer Write em w: o cookie de
	// sessão (Set-Cookie) só é aceito pelo cliente se sair ANTES do corpo
	// (ver "Achado" abaixo). Também é por isso que renderizamos em buffer:
	// se o template falhar no meio, w nunca recebe bytes parciais — ou sai
	// a página inteira, ou sai um 500 limpo.
	_ = a.ensureToken(w, r)

	utm := utmFromURL(r.URL)
	data := pageData{
		Lang:         leads.NormalizeLanguage(lang),
		PagePath:     pagePath,
		BaseURL:      a.baseURL,
		OGImageURL:   a.baseURL + "/static/" + ogImageFile(lang),
		ContactEmail: a.contactEmail,
		UTM:          utm,
		Values:       map[string]string{},
	}

	var buf bytes.Buffer
	if err := a.tmpl.ExecuteTemplate(&buf, templateName, data); err != nil {
		slog.Error("rendering landing", "template", templateName, "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

func (a *app) ensureToken(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie(sessionCookieName); err == nil && strings.TrimSpace(c.Value) != "" {
		return c.Value
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		fallback := make([]byte, 16)
		for i := range fallback {
			fallback[i] = byte(i + 1)
		}
		b = fallback
	}
	token := hex.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 365,
	})
	return token
}

func withTimeoutCtx(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 3*time.Second)
}

// ogImageFile devolve o PNG de OG certo por idioma (épico 002, T15 — fecha
// G1 do épico 001: og:image era SVG e não renderizava preview no
// LinkedIn/WhatsApp). Assets gerados por marketing em
// roosterlabs-marketing/brand/web/, copiados para web/static/.
func ogImageFile(lang string) string {
	if leads.NormalizeLanguage(lang) == "en" {
		return "og-en.png"
	}
	return "og.png"
}

func utmFromURL(u *url.URL) leads.UTM {
	q := u.Query()
	return leads.UTM{
		Source:   strings.TrimSpace(q.Get("utm_source")),
		Medium:   strings.TrimSpace(q.Get("utm_medium")),
		Campaign: strings.TrimSpace(q.Get("utm_campaign")),
		Term:     strings.TrimSpace(q.Get("utm_term")),
		Content:  strings.TrimSpace(q.Get("utm_content")),
	}
}

func utmFromForm(r *http.Request) leads.UTM {
	return leads.UTM{
		Source:   strings.TrimSpace(r.FormValue("utm_source")),
		Medium:   strings.TrimSpace(r.FormValue("utm_medium")),
		Campaign: strings.TrimSpace(r.FormValue("utm_campaign")),
		Term:     strings.TrimSpace(r.FormValue("utm_term")),
		Content:  strings.TrimSpace(r.FormValue("utm_content")),
	}
}
