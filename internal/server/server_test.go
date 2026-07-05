package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	h, err := NewHandler()
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
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
	h, err := NewHandler()
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("got Content-Type %q, want text/html", ct)
	}
	if !strings.Contains(rec.Body.String(), "RoosterLabs") {
		t.Fatal("index page does not mention RoosterLabs")
	}
}

func TestStaticFileServed(t *testing.T) {
	h, err := NewHandler()
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
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
	h, err := NewHandler()
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nope", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}
