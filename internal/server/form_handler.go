package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/roosterlabs/roosterlabs-engineering/internal/leads"
)

func (a *app) handleFormStep(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form payload", http.StatusBadRequest)
		return
	}

	step, err := strconv.Atoi(r.PathValue("step"))
	if err != nil || step < 1 || step > 5 {
		http.Error(w, "invalid step", http.StatusBadRequest)
		return
	}

	token := a.ensureToken(w, r)
	lang := leads.NormalizeLanguage(r.FormValue("lang"))
	if lang == "pt-BR" && strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.FormValue("page_path"))), "/en") {
		lang = "en"
	}
	pagePath := strings.TrimSpace(r.FormValue("page_path"))
	if pagePath == "" {
		if lang == "en" {
			pagePath = "/en/"
		} else {
			pagePath = "/"
		}
	}

	utm := utmFromForm(r)
	values := map[string]string{}
	for k, arr := range r.Form {
		if len(arr) > 0 {
			values[k] = arr[0]
		}
	}

	if step < 5 {
		errMsg := validateStep(step, values)
		if errMsg != "" {
			a.renderFormStep(w, step, pageData{
				Lang:         lang,
				PagePath:     pagePath,
				UTM:          utm,
				ContactEmail: a.contactEmail,
				Error:        errMsg,
				Values:       values,
			})
			return
		}

		ctx, cancel := withTimeoutCtx(r)
		defer cancel()
		if err := a.store.RecordFormAnswer(ctx, leads.FormAnswerEvent{
			Token:    token,
			Step:     step,
			Language: lang,
			PagePath: pagePath,
			UTM:      utm,
			Payload: map[string]string{
				"choice": strings.TrimSpace(values["choice"]),
				"other":  strings.TrimSpace(values["other"]),
			},
		}); err != nil {
			http.Error(w, "failed to record answer", http.StatusInternalServerError)
			return
		}

		a.renderFormStep(w, step+1, pageData{
			Lang:         lang,
			PagePath:     pagePath,
			UTM:          utm,
			ContactEmail: a.contactEmail,
			Values:       map[string]string{},
		})
		return
	}

	errMsg := validateFinal(values)
	if errMsg != "" {
		a.renderFormStep(w, 5, pageData{
			Lang:         lang,
			PagePath:     pagePath,
			UTM:          utm,
			ContactEmail: a.contactEmail,
			Error:        errMsg,
			Values:       values,
		})
		return
	}

	// validateFinal já confirmou que normalizeLinkedInURL não erra para este
	// valor — o erro aqui só seria alcançável por uma corrida de dados
	// improvável entre validação e uso; nesse caso extremo, falha explícito
	// em vez de gravar lead com LinkedIn vazio.
	linkedinURL, err := normalizeLinkedInURL(values["linkedin_handle"])
	if err != nil {
		http.Error(w, "invalid linkedin handle", http.StatusInternalServerError)
		return
	}

	ctx, cancel := withTimeoutCtx(r)
	defer cancel()
	if err := a.store.FinalizeLead(ctx, leads.FinalizeInput{
		Token:    token,
		Language: lang,
		PagePath: pagePath,
		UTM:      utm,
		Email:    strings.TrimSpace(values["email"]),
		LinkedIn: linkedinURL,
	}); err != nil {
		http.Error(w, "failed to persist lead", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := a.tmpl.ExecuteTemplate(&buf, formTemplateName(lang, 6), pageData{
		Lang:         lang,
		PagePath:     pagePath,
		UTM:          utm,
		ContactEmail: a.contactEmail,
	}); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

func (a *app) renderFormStep(w http.ResponseWriter, step int, data pageData) {
	if data.Values == nil {
		data.Values = map[string]string{}
	}

	var buf bytes.Buffer
	if err := a.tmpl.ExecuteTemplate(&buf, formTemplateName(data.Lang, step), data); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

func formTemplateName(lang string, step int) string {
	suffix := "pt"
	if leads.NormalizeLanguage(lang) == "en" {
		suffix = "en"
	}
	return fmt.Sprintf("form_step_%d_%s", step, suffix)
}

func validateStep(step int, values map[string]string) string {
	choice := strings.TrimSpace(values["choice"])
	other := strings.TrimSpace(values["other"])
	if choice == "" {
		return "Selecione uma opcao para continuar."
	}
	if step == 1 || step == 2 {
		if strings.EqualFold(choice, "Outro") || strings.EqualFold(choice, "Other") {
			if other == "" {
				return "Preencha o campo de detalhe em \"Outro\" para continuar."
			}
		}
	}
	return ""
}

func validateFinal(values map[string]string) string {
	email := strings.TrimSpace(values["email"])
	if _, err := mail.ParseAddress(email); err != nil {
		return "Informe um e-mail valido."
	}
	if _, err := normalizeLinkedInURL(values["linkedin_handle"]); err != nil {
		return "Informe seu usuario do LinkedIn."
	}
	return ""
}

// normalizeLinkedInURL constrói a URL canônica do LinkedIn a partir do
// handle (formato do form v0.4: campo com prefixo fixo linkedin.com/in/,
// épico 002, T13). Aceita defensivamente que o usuário tenha colado a URL
// inteira no campo — nesse caso extrai só o handle antes de remontar,
// para não duplicar o prefixo (ex.: "https://www.linkedin.com/in/foo" e
// "linkedin.com/in/foo/" viram ambos ".../in/foo").
func normalizeLinkedInURL(raw string) (string, error) {
	handle := strings.TrimSpace(raw)
	handle = strings.TrimPrefix(handle, "https://")
	handle = strings.TrimPrefix(handle, "http://")
	handle = strings.TrimPrefix(handle, "www.")
	handle = strings.TrimPrefix(handle, "linkedin.com/in/")
	handle = strings.TrimPrefix(handle, "linkedin.com/")
	if i := strings.IndexAny(handle, "/?"); i >= 0 {
		handle = handle[:i]
	}
	handle = strings.TrimSpace(handle)

	if handle == "" {
		return "", fmt.Errorf("linkedin handle vazio")
	}

	return "https://www.linkedin.com/in/" + handle, nil
}
