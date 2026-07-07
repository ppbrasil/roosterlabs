package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	h, err := NewHandler(Config{})
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
	return h
}

func TestHealthz(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("got body %q, want %q", rec.Body.String(), "ok")
	}
}

func TestIndexRenders(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("got Content-Type %q, want text/html", ct)
	}
	if !strings.Contains(rec.Body.String(), "AUTO-AUTENTICIDADE") {
		t.Fatal("index page does not mention RoosterLabs")
	}
	// O form do passo 1 é embutido no fim da página; se o template falhar no
	// meio da renderização, o corpo fica truncado e este check pega o bug
	// (checar só o headline deixa passar renderização parcial).
	if !strings.Contains(rec.Body.String(), `id="lead-form"`) {
		t.Fatal("index page rendered without embedded lead form (partial render?)")
	}
}

func TestEnglishIndexRenders(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/en/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "AUTO-AUTHENTICITY") {
		t.Fatal("english page did not render expected heading")
	}
	if !strings.Contains(rec.Body.String(), `id="lead-form"`) {
		t.Fatal("english page rendered without embedded lead form (partial render?)")
	}
}

func TestStaticFileServed(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/static/robots.txt", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "User-agent") {
		t.Fatal("robots.txt content not served")
	}
}

func TestUnknownPathIs404(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nope", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestViewEventEndpoint(t *testing.T) {
	h := newTestHandler(t)
	body := []byte(`{"path":"/en/","language":"en"}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/event/view", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestFormFlowToSuccess(t *testing.T) {
	h := newTestHandler(t)
	cookie := ""

	steps := []url.Values{
		{"choice": {"Consultor, advisor ou fractional"}, "lang": {"pt-BR"}, "page_path": {"/"}},
		{"choice": {"Gerar negocios"}, "lang": {"pt-BR"}, "page_path": {"/"}},
		{"choice": {"Parado"}, "lang": {"pt-BR"}, "page_path": {"/"}},
		{"choice": {"Tempo"}, "lang": {"pt-BR"}, "page_path": {"/"}},
		{"email": {"contact@roosterlabs.com.br"}, "linkedin": {"https://www.linkedin.com/in/roosterlabs"}, "lang": {"pt-BR"}, "page_path": {"/"}},
	}

	for i, vals := range steps {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/form/"+strconv.Itoa(i+1), strings.NewReader(vals.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if cookie != "" {
			req.Header.Set("Cookie", cookie)
		}
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("step %d got status %d, want %d", i+1, rec.Code, http.StatusOK)
		}
		if c := rec.Header().Get("Set-Cookie"); c != "" {
			cookie = c
		}
	}
}

func TestFormStepValidation(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/form/1", strings.NewReader("lang=pt-BR&page_path=/"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "Selecione uma opcao") {
		t.Fatal("expected validation message")
	}
}
