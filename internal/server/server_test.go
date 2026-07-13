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
	body := rec.Body.String()
	// v0.5 (épico 003, T2): tagline promovida a H1, só o "AI" em âmbar —
	// o span quebra a string contígua, então o assert mira o markup exato.
	if !strings.Contains(body, `<h1>Você. <span class="ai-accent">AI</span>mplificado.</h1>`) {
		t.Fatal("index page missing v0.5 hero H1 (tagline with amber AI prefix)")
	}
	if strings.Contains(body, `class="eyebrow"`) {
		t.Fatal("index page still renders the eyebrow (retired in v0.5)")
	}
	if strings.Contains(body, "Amplificar você não é escrever por você") {
		t.Fatal("index page still renders the v0.4 hero H1")
	}
	if !strings.Contains(body, "A RoosterLabs garante a sua cadência de conteúdo no LinkedIn") {
		t.Fatal("index page missing v0.5 hero sub")
	}
	if !strings.Contains(body, "O conteúdo é seu. O trabalho pesado, nosso.") {
		t.Fatal("index page missing v0.5 hero closing line")
	}
	if !strings.Contains(body, `href="#lead-form"`) || !strings.Contains(body, "Quero uma vaga no alfa") {
		t.Fatal("index page missing hero CTA anchored to #lead-form")
	}
	if strings.Contains(body, "AUTO-AUTENTICIDADE") {
		t.Fatal("index page still renders retired v0.3 copy (Auto-Autenticidade)")
	}
	// Decisão do aceite do épico 003: title e og:description NÃO mudam na
	// v0.5 (achado do revisor da T2 — sem estes asserts, uma edição futura
	// do <head> passaria verde).
	if !strings.Contains(body, "<title>RoosterLabs · Você. AImplificado.</title>") {
		t.Fatal("index page title must keep the tagline (unchanged in v0.5)")
	}
	if !strings.Contains(body, `og:description" content="Saber nunca foi o seu problema. Transformar isso em conteúdo, sim."`) {
		t.Fatal("index page og:description must keep the problem H2 (unchanged in v0.5)")
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
	enBody := rec.Body.String()
	if !strings.Contains(enBody, `<h1>You. <span class="ai-accent">AI</span>mplified.</h1>`) {
		t.Fatal("english page missing v0.5 hero H1 (tagline with amber AI prefix)")
	}
	if strings.Contains(enBody, `class="eyebrow"`) {
		t.Fatal("english page still renders the eyebrow (retired in v0.5)")
	}
	if strings.Contains(enBody, "Amplifying you isn't writing for you") {
		t.Fatal("english page still renders the v0.4 hero H1")
	}
	if !strings.Contains(enBody, "RoosterLabs guarantees your LinkedIn content cadence") {
		t.Fatal("english page missing v0.5 hero sub")
	}
	if !strings.Contains(enBody, "The content is yours. The heavy lifting, ours.") {
		t.Fatal("english page missing v0.5 hero closing line")
	}
	if !strings.Contains(enBody, `href="#lead-form"`) || !strings.Contains(enBody, "Get an alpha seat") {
		t.Fatal("english page missing hero CTA anchored to #lead-form")
	}
	if strings.Contains(enBody, "AUTO-AUTHENTICITY") {
		t.Fatal("english page still renders retired v0.3 copy (Auto-Authenticity)")
	}
	if !strings.Contains(enBody, "<title>RoosterLabs · You. AImplified.</title>") {
		t.Fatal("english page title must keep the tagline (unchanged in v0.5)")
	}
	if !strings.Contains(enBody, `og:description" content="Knowing was never your problem. Turning it into content is."`) {
		t.Fatal("english page og:description must keep the problem H2 (unchanged in v0.5)")
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
			// v0.5 (épico 003, T4): ✓/✗ semânticos — cada marca embrulhada
			// no span da classe certa, sem marca fora de span.
			if got := strings.Count(body, `<span class="mark-ok">✓</span>`); got != 12 {
				t.Fatalf("got %d ok marks in mark-ok spans, want 12", got)
			}
			if got := strings.Count(body, `<span class="mark-no">✗</span>`); got != 6 {
				t.Fatalf("got %d cross marks in mark-no spans, want 6", got)
			}
			if got := strings.Count(body, "✓"); got != 12 {
				t.Fatalf("got %d total check marks, want 12 (mark outside span?)", got)
			}
			if got := strings.Count(body, "✗"); got != 6 {
				t.Fatalf("got %d total cross marks, want 6 (mark outside span?)", got)
			}
		})
	}

	// Guarda do CSS (padrão T1/T3): tokens semânticos definidos e aplicados —
	// e só em contexto de comparação (regra da identidade v0.5).
	h := newTestHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/static/site.css", nil))
	css := rec.Body.String()
	for _, want := range []string{"--ok-500: #82b88a", "--no-500: #c96f6f", ".mark-ok", ".mark-no"} {
		if !strings.Contains(css, want) {
			t.Fatalf("/static/site.css: missing %q (semantic ok/no tokens, v0.5)", want)
		}
	}
	// Regra 6 da identidade: ok/no só em comparação. Único consumidor
	// permitido de cada token é a própria classe .mark-* (achado do revisor
	// da T4 — presença não impede var(--ok-500) vazar para CTA/link).
	for _, token := range []string{"var(--ok-500)", "var(--no-500)"} {
		if got := strings.Count(css, token); got != 1 {
			t.Fatalf("/static/site.css: %s used %d times, want exactly 1 (comparison marks only)", token, got)
		}
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

// TestSiteWatermark cobre o T3 + a emenda T8 do épico 003: galo decorativo
// atrás de todo o conteúdo (backdrop da página, não mais ornamento do hero)
// nas duas rotas (aria-hidden, alt vazio) e asset .webp servido com o
// Content-Type certo.
func TestSiteWatermark(t *testing.T) {
	h := newTestHandler(t)

	for _, path := range []string{"/", "/en/"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		body := rec.Body.String()

		if !strings.Contains(body, `class="site-watermark"`) {
			t.Fatalf("%s: missing site watermark element", path)
		}
		if !strings.Contains(body, `src="/static/rooster-watermark.webp" alt="" aria-hidden="true"`) {
			t.Fatalf("%s: site watermark must be decorative (empty alt + aria-hidden)", path)
		}
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/static/rooster-watermark.webp", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("watermark asset: got status %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/webp" {
		t.Fatalf("watermark asset: got Content-Type %q, want image/webp", ct)
	}

	// Guarda do CSS (achado do revisor da T3, mesmo padrão da T1; atualizado
	// na T8): sem isto, reverter o bloco do watermark deixa a <img> órfã em
	// fluxo, opaca, empurrando o layout — com CI verde. `position: fixed`
	// trava o comportamento de backdrop (a emenda T8): se alguém voltar para
	// absolute/quina, o teste pega.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/static/site.css", nil))
	css := rec.Body.String()
	if !strings.Contains(css, ".site-watermark") || !strings.Contains(css, "opacity: 0.16") {
		t.Fatal("/static/site.css: missing site watermark styles (opacity 16%)")
	}
	if !strings.Contains(css, "position: fixed") {
		t.Fatal("/static/site.css: site watermark must be fixed (backdrop da página, emenda T8)")
	}
}

// TestLandingAmendmentsT8T9 cobre as emendas de 2026-07-13: favicon = ícone da
// marca (item 2) e os dois primeiros H2 em duas linhas com a 2ª em âmbar (T9).
// A copy dos H2 em `landing-page.md`/`visual-identity.md` ainda vai ser
// sincronizada por marketing (ver nota no épico 003); este guard trava o que
// já está em produção para que a sync não remova as quebras sem querer.
func TestLandingAmendmentsT8T9(t *testing.T) {
	h := newTestHandler(t)

	// Favicon aponta para o ícone da marca (mesmo do rodapé), não o SVG
	// provisório.
	for _, path := range []string{"/", "/en/"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		body := rec.Body.String()
		if !strings.Contains(body, `rel="icon" href="/static/rooster-icon.png"`) {
			t.Fatalf("%s: favicon deve ser /static/rooster-icon.png (item 2)", path)
		}
	}

	cases := []struct {
		path, want string
	}{
		{"/", `Saber nunca foi o seu problema.<br><span class="h2-accent">Transformar isso em conteúdo, sim.</span>`},
		{"/", `Do gatilho à publicação,<br><span class="h2-accent">com você no centro.</span>`},
		{"/en/", `Knowing was never your problem.<br><span class="h2-accent">Turning it into content is.</span>`},
		{"/en/", `From trigger to published,<br><span class="h2-accent">with you at the center.</span>`},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, c.path, nil))
		if !strings.Contains(rec.Body.String(), c.want) {
			t.Fatalf("%s: H2 âmbar (T9) fora do esperado, faltou: %s", c.path, c.want)
		}
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/static/site.css", nil))
	if css := rec.Body.String(); !strings.Contains(css, ".h2-accent") {
		t.Fatal("/static/site.css: missing .h2-accent (2ª linha âmbar do H2, T9)")
	}
}

// TestFooterBrand cobre o T5 do épico 003: ícone + wordmark no rodapé das
// duas rotas, asset otimizado servido corretamente.
func TestFooterBrand(t *testing.T) {
	h := newTestHandler(t)

	for _, path := range []string{"/", "/en/"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		body := rec.Body.String()

		if !strings.Contains(body, `src="/static/rooster-icon.png" alt="" width="28" height="28"`) {
			t.Fatalf("%s: missing decorative footer icon", path)
		}
		if !strings.Contains(body, `<span class="wordmark">Rooster<span class="wordmark-labs">Labs</span></span>`) {
			t.Fatalf("%s: missing two-tone wordmark in footer", path)
		}
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/static/rooster-icon.png", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("footer icon asset: got status %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("footer icon asset: got Content-Type %q, want image/png", ct)
	}
	// Plano da T5: asset otimizado (derivado ~112px), nunca o master de
	// 1,6 MB por engano.
	if size := rec.Body.Len(); size > 30*1024 {
		t.Fatalf("footer icon asset: %d bytes, want <=30KB (optimized derivative)", size)
	}

	// Guarda do CSS (padrão T1/T3/T4): wordmark em Albert Sans 500 com
	// "Labs" em âmbar.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/static/site.css", nil))
	css := rec.Body.String()
	for _, want := range []string{".footer-brand", ".wordmark"} {
		if !strings.Contains(css, want) {
			t.Fatalf("/static/site.css: missing %q (footer brand, v0.5)", want)
		}
	}
	// Achado do revisor da T5: presença de seletor não impede regressão de
	// valor — 300 no wordmark ("não segura pequeno") ou "Labs" fora do âmbar
	// passariam verdes. Asserts miram os valores da spec.
	if !strings.Contains(css, ".wordmark {\n  font-family: \"Albert Sans\", Arial, sans-serif;\n  font-weight: 500;") {
		t.Fatal("/static/site.css: wordmark must be Albert Sans weight 500 (visual identity v0.5)")
	}
	if !strings.Contains(css, ".wordmark-labs {\n  color: var(--amber-500);") {
		t.Fatal("/static/site.css: wordmark 'Labs' must be amber (visual identity v0.5)")
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
		for _, family := range []string{"family=Albert+Sans:wght@300;400;500", "family=Inter", "family=IBM+Plex+Mono"} {
			if !strings.Contains(body, family) {
				t.Fatalf("%s: missing font family %q in Google Fonts link", path, family)
			}
		}
		if !strings.Contains(body, "display=swap") {
			t.Fatalf("%s: Google Fonts link missing display=swap", path)
		}
		// Fraunces aposentada na identidade v0.5 (guardrail da serifa, 2º acionamento).
		if strings.Contains(body, "Fraunces") {
			t.Fatalf("%s: Fraunces must not be referenced anymore (visual identity v0.5)", path)
		}
	}

	// O HTML carregar a fonte não basta: o CSS precisa usá-la (achado do
	// revisor da T1 — sem isso, reverter o site.css deixa o CI verde com os
	// headings caindo no fallback).
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/static/site.css", nil))
	css := rec.Body.String()
	if !strings.Contains(css, `"Albert Sans"`) {
		t.Fatal("/static/site.css: headings must use Albert Sans (visual identity v0.5)")
	}
	if strings.Contains(css, "Fraunces") {
		t.Fatal("/static/site.css: Fraunces must not be referenced anymore (visual identity v0.5)")
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
