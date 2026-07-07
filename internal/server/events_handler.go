package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/roosterlabs/roosterlabs-engineering/internal/leads"
)

type viewPayload struct {
	Path       string `json:"path"`
	Language   string `json:"language"`
	UTMSource  string `json:"utm_source"`
	UTMMedium  string `json:"utm_medium"`
	UTMCampaign string `json:"utm_campaign"`
	UTMTerm    string `json:"utm_term"`
	UTMContent string `json:"utm_content"`
}

func (a *app) handleViewEvent(w http.ResponseWriter, r *http.Request) {
	token := a.ensureToken(w, r)
	payload := viewPayload{}
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid json payload", http.StatusBadRequest)
			return
		}
	}

	lang := leads.NormalizeLanguage(firstNonEmpty(payload.Language, r.URL.Query().Get("language")))
	if lang == "pt-BR" && strings.HasPrefix(strings.ToLower(strings.TrimSpace(payload.Path)), "/en") {
		lang = "en"
	}
	path := strings.TrimSpace(firstNonEmpty(payload.Path, r.URL.Query().Get("path")))
	if path == "" {
		if lang == "en" {
			path = "/en/"
		} else {
			path = "/"
		}
	}

	utm := leads.UTM{
		Source:   strings.TrimSpace(firstNonEmpty(payload.UTMSource, r.URL.Query().Get("utm_source"))),
		Medium:   strings.TrimSpace(firstNonEmpty(payload.UTMMedium, r.URL.Query().Get("utm_medium"))),
		Campaign: strings.TrimSpace(firstNonEmpty(payload.UTMCampaign, r.URL.Query().Get("utm_campaign"))),
		Term:     strings.TrimSpace(firstNonEmpty(payload.UTMTerm, r.URL.Query().Get("utm_term"))),
		Content:  strings.TrimSpace(firstNonEmpty(payload.UTMContent, r.URL.Query().Get("utm_content"))),
	}

	ctx, cancel := withTimeoutCtx(r)
	defer cancel()

	if err := a.store.RecordViewEvent(ctx, leads.ViewEvent{
		Token:    token,
		Language: lang,
		PagePath: path,
		UTM:      utm,
	}); err != nil {
		http.Error(w, "failed to record event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if s := strings.TrimSpace(value); s != "" {
			return s
		}
	}
	return ""
}
