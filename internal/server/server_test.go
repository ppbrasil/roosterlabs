package server

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/roosterlabs/roosterlabs-engineering/internal/leads"
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
	// v0.4: tagline "Você. AImplificado." substitui "Auto-Autenticidade"
	// (épico 002, T9) — checa a copy nova, não a antiga.
	if !strings.Contains(rec.Body.String(), "Você. AImplificado.") {
		t.Fatal("index page does not mention the v0.4 tagline")
	}
	if strings.Contains(rec.Body.String(), "AUTO-AUTENTICIDADE") {
		t.Fatal("index page still renders retired v0.3 copy (Auto-Autenticidade)")
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
	if !strings.Contains(rec.Body.String(), "You. AImplified.") {
		t.Fatal("english page does not mention the v0.4 tagline")
	}
	if strings.Contains(rec.Body.String(), "AUTO-AUTHENTICITY") {
		t.Fatal("english page still renders retired v0.3 copy (Auto-Authenticity)")
	}
	if !strings.Contains(rec.Body.String(), `id="lead-form"`) {
		t.Fatal("english page rendered without embedded lead form (partial render?)")
	}
}

// TestOGImagePerLanguage cobre o T15 do épico 002: og:image era um SVG que
// não renderizava preview no LinkedIn/WhatsApp (G1 do épico 001). Checa que
// cada idioma aponta pro PNG certo, que as metatags width/height/type saem
// certas, e que o arquivo é servido com o Content-Type correto.
func TestOGImagePerLanguage(t *testing.T) {
	h := newTestHandler(t)

	cases := []struct {
		name     string
		path     string
		wantFile string
	}{
		{"pt", "/", "og.png"},
		{"en", "/en/", "og-en.png"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			if rec.Code != http.StatusOK {
				t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
			}
			body := rec.Body.String()

			if !strings.Contains(body, `og:image" content="https://roosterlabs.com.br/static/`+tc.wantFile+`"`) {
				t.Fatalf("og:image does not point to /static/%s", tc.wantFile)
			}
			if !strings.Contains(body, `og:image:width" content="1200"`) {
				t.Fatal("missing og:image:width=1200")
			}
			if !strings.Contains(body, `og:image:height" content="630"`) {
				t.Fatal("missing og:image:height=630")
			}
			if !strings.Contains(body, `og:image:type" content="image/png"`) {
				t.Fatal("missing og:image:type=image/png")
			}
		})
	}

	// O arquivo precisa existir de fato e ser servido como PNG — sem isso as
	// metatags acima apontam pro nada.
	for _, file := range []string{"og.png", "og-en.png"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/static/"+file, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /static/%s: got status %d, want %d", file, rec.Code, http.StatusOK)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "image/png" {
			t.Fatalf("GET /static/%s: got Content-Type %q, want image/png", file, ct)
		}
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
		{"choice": {"Consultor ou advisor independente"}, "lang": {"pt-BR"}, "page_path": {"/"}},
		{"choice": {"Gerar negócios"}, "lang": {"pt-BR"}, "page_path": {"/"}},
		{"choice": {"Parado"}, "lang": {"pt-BR"}, "page_path": {"/"}},
		{"choice": {"Tempo"}, "lang": {"pt-BR"}, "page_path": {"/"}},
		{"email": {"contact@roosterlabs.com.br"}, "linkedin_handle": {"roosterlabs"}, "lang": {"pt-BR"}, "page_path": {"/"}},
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

// newBrokenTemplateApp constrói um *app cujo conjunto de templates não tem
// nenhum dos nomes usados em produção — qualquer ExecuteTemplate falha com
// "no such template", simulando o cenário de erro no meio do render sem
// precisar corromper os templates reais.
func newBrokenTemplateApp() *app {
	return &app{
		tmpl:         template.Must(template.New("empty").Parse("")),
		store:        leads.NewMemoryStore(),
		baseURL:      "https://roosterlabs.com.br",
		contactEmail: "contact@roosterlabs.com.br",
	}
}

func TestRenderLandingTemplateErrorDoesNotWritePartialBody(t *testing.T) {
	a := newBrokenTemplateApp()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	a.renderLanding(rec, req, "pt-BR", "index.html.tmpl", "/")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if strings.Contains(rec.Body.String(), "<html") {
		t.Fatal("response should never contain partial HTML from a failed template render")
	}
}

func TestRenderFormStepTemplateErrorDoesNotWritePartialBody(t *testing.T) {
	a := newBrokenTemplateApp()
	rec := httptest.NewRecorder()

	a.renderFormStep(rec, 1, pageData{Lang: "pt-BR", PagePath: "/"})

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if strings.Contains(rec.Body.String(), `id="lead-form"`) {
		t.Fatal("response should never contain partial form HTML from a failed template render")
	}
}

// Regressão do achado do épico 002 (T7): antes do render em buffer,
// ensureToken rodava DEPOIS do ExecuteTemplate(w, ...) escrever direto na
// resposta — como isso já dispara o WriteHeader implícito, o Set-Cookie
// que ensureToken tentava adicionar nunca saía no GET / (headers já
// fechados). Trava esse comportamento.
func TestIndexSetsSessionCookie(t *testing.T) {
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	cookie := rec.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, "rlsid=") {
		t.Fatalf("expected Set-Cookie with rlsid on GET /, got %q", cookie)
	}
}

func TestPositioningTableRendersAllRows(t *testing.T) {
	cases := []struct {
		name string
		path string
		rows []string
	}{
		{
			name: "pt",
			path: "/",
			rows: []string{
				"Escreve o post",
				"A partir do seu take real, não de template",
				"Gatilhos + linha editorial próprios",
				"Consistência como parte do serviço",
				"Na sua voz — nem genérica, nem emprestada",
				"Preço de software",
			},
		},
		{
			name: "en",
			path: "/en/",
			rows: []string{
				"Writes the post",
				"From your real take, not a template",
				"Triggers + an editorial line of your own",
				"Consistency as part of the service",
				"In your voice — not generic, not borrowed",
				"Software price",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			body := rec.Body.String()
			if !strings.Contains(body, `class="compare-table"`) {
				t.Fatal("positioning table not rendered")
			}
			for _, row := range tc.rows {
				if !strings.Contains(body, row) {
					t.Fatalf("missing comparison row %q", row)
				}
			}
			if got := strings.Count(body, "✓"); got != 12 {
				t.Fatalf("got %d check marks, want 12", got)
			}
			if got := strings.Count(body, "✗"); got != 6 {
				t.Fatalf("got %d cross marks, want 6", got)
			}
		})
	}
}

func TestAlphaAccessSectionRenders(t *testing.T) {
	h := newTestHandler(t)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	body := rec.Body.String()
	if !strings.Contains(body, "Vagas alfa abertas — poucas, de propósito.") {
		t.Fatal("missing PT v0.4 alpha-access heading")
	}
	if !strings.Contains(body, "Depois vêm as vagas beta") {
		t.Fatal("missing PT v0.4 alpha-access paragraph")
	}

	recEN := httptest.NewRecorder()
	h.ServeHTTP(recEN, httptest.NewRequest(http.MethodGet, "/en/", nil))
	bodyEN := recEN.Body.String()
	if !strings.Contains(bodyEN, "Alpha seats open — few, on purpose.") {
		t.Fatal("missing EN v0.4 alpha-access heading")
	}
	if !strings.Contains(bodyEN, "Beta seats come next") {
		t.Fatal("missing EN v0.4 alpha-access paragraph")
	}
}

func TestNormalizeLinkedInURL(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "handle simples", raw: "joao-silva", want: "https://www.linkedin.com/in/joao-silva"},
		{name: "handle vazio", raw: "   ", wantErr: true},
		{name: "URL completa colada por engano", raw: "https://www.linkedin.com/in/joao-silva", want: "https://www.linkedin.com/in/joao-silva"},
		{name: "sem esquema, com www", raw: "www.linkedin.com/in/joao-silva", want: "https://www.linkedin.com/in/joao-silva"},
		{name: "com barra final", raw: "linkedin.com/in/joao-silva/", want: "https://www.linkedin.com/in/joao-silva"},
		{name: "com query string colada", raw: "joao-silva?trk=public_profile", want: "https://www.linkedin.com/in/joao-silva"},
		{name: "espacos ao redor", raw: "  joao-silva  ", want: "https://www.linkedin.com/in/joao-silva"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := normalizeLinkedInURL(tc.raw)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("normalizeLinkedInURL(%q) = %q, want error", tc.raw, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeLinkedInURL(%q) unexpected error: %v", tc.raw, err)
			}
			if got != tc.want {
				t.Fatalf("normalizeLinkedInURL(%q) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}

func TestFormStep5UsesFixedLinkedInPrefix(t *testing.T) {
	h := newTestHandler(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/form/5", strings.NewReader("lang=pt-BR&page_path=/"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, `class="input-prefix"`) {
		t.Fatal("step 5 does not render the fixed linkedin.com/in/ prefix")
	}
	if !strings.Contains(body, `name="linkedin_handle"`) {
		t.Fatal("step 5 does not have a linkedin_handle input")
	}
	if strings.Contains(body, `name="linkedin"`) {
		t.Fatal("step 5 still has the old full-URL linkedin input")
	}
}

// postFormStep submete um passo do carrossel e devolve o corpo do
// fragmento seguinte, propagando o cookie de sessão entre chamadas.
func postFormStep(t *testing.T, h http.Handler, cookie *string, step int, vals url.Values) string {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/form/"+strconv.Itoa(step), strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if *cookie != "" {
		req.Header.Set("Cookie", *cookie)
	}
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("step %d got status %d, want %d (body: %s)", step, rec.Code, http.StatusOK, rec.Body.String())
	}
	if c := rec.Header().Get("Set-Cookie"); c != "" {
		*cookie = c
	}
	return rec.Body.String()
}

func TestFormOptionsMatchV04(t *testing.T) {
	h := newTestHandler(t)

	// GET / só embute o passo 1 do carrossel — passos 2+ só aparecem na
	// resposta do POST do passo anterior (fragmento htmx).
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	step1PT := rec.Body.String()
	if strings.Contains(step1PT, "fractional") {
		t.Fatal("PT step 1 should not mention \"fractional\" per v0.4 (EN-only term)")
	}

	cookie := ""
	step2PT := postFormStep(t, h, &cookie, 1, url.Values{
		"choice": {"Consultor ou advisor independente"}, "lang": {"pt-BR"}, "page_path": {"/"},
	})
	if !strings.Contains(step2PT, "Construir audiência") {
		t.Fatal("PT step 2 missing v0.4 option \"Construir audiência\"")
	}

	step3PT := postFormStep(t, h, &cookie, 2, url.Values{
		"choice": {"Gerar negócios"}, "lang": {"pt-BR"}, "page_path": {"/"},
	})
	if !strings.Contains(step3PT, "Só comento e reajo a posts dos outros") {
		t.Fatal("PT step 3 missing v0.4 option about reacting/commenting only")
	}

	recEN := httptest.NewRecorder()
	h.ServeHTTP(recEN, httptest.NewRequest(http.MethodGet, "/en/", nil))
	step1EN := recEN.Body.String()
	if !strings.Contains(step1EN, "Consultant, advisor or fractional") {
		t.Fatal("EN step 1 should keep \"fractional\" per v0.4")
	}

	cookieEN := ""
	step2EN := postFormStep(t, h, &cookieEN, 1, url.Values{
		"choice": {"Consultant, advisor or fractional"}, "lang": {"en"}, "page_path": {"/en/"},
	})
	if !strings.Contains(step2EN, "Build an audience") {
		t.Fatal("EN step 2 missing v0.4 option \"Build an audience\"")
	}

	step3EN := postFormStep(t, h, &cookieEN, 2, url.Values{
		"choice": {"Generate business"}, "lang": {"en"}, "page_path": {"/en/"},
	})
	if !strings.Contains(step3EN, "I only react and comment on others") {
		t.Fatal("EN step 3 missing v0.4 option about reacting/commenting only")
	}
}

func TestFontsLoadedFromGoogleFonts(t *testing.T) {
	h := newTestHandler(t)

	for _, path := range []string{"/", "/en/"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		body := rec.Body.String()

		if !strings.Contains(body, "fonts.googleapis.com/css2") {
			t.Fatalf("%s: missing Google Fonts stylesheet link", path)
		}
		for _, family := range []string{"family=Fraunces", "family=Inter", "family=IBM+Plex+Mono"} {
			if !strings.Contains(body, family) {
				t.Fatalf("%s: missing font family %q in Google Fonts link", path, family)
			}
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
