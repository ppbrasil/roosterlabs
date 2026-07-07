package server

import (
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"

	"github.com/roosterlabs/roosterlabs-engineering/internal/leads"
)

type formTemplateData struct {
	Lang         string
	PagePath     string
	UTM          leads.UTM
	ContactEmail string
	Error        string
	Values       map[string]string
}

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
			a.renderFormStep(w, step, formTemplateData{
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

		a.renderFormStep(w, step+1, formTemplateData{
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
		a.renderFormStep(w, 5, formTemplateData{
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
	if err := a.store.FinalizeLead(ctx, leads.FinalizeInput{
		Token:    token,
		Language: lang,
		PagePath: pagePath,
		UTM:      utm,
		Email:    strings.TrimSpace(values["email"]),
		LinkedIn: strings.TrimSpace(values["linkedin"]),
	}); err != nil {
		http.Error(w, "failed to persist lead", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.tmpl.ExecuteTemplate(w, formTemplateName(lang, 6), formTemplateData{
		Lang:         lang,
		PagePath:     pagePath,
		UTM:          utm,
		ContactEmail: a.contactEmail,
	}); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (a *app) renderFormStep(w http.ResponseWriter, step int, data formTemplateData) {
	if data.Values == nil {
		data.Values = map[string]string{}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.tmpl.ExecuteTemplate(w, formTemplateName(data.Lang, step), data); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
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
	linkedin := strings.TrimSpace(values["linkedin"])
	if _, err := mail.ParseAddress(email); err != nil {
		return "Informe um e-mail valido."
	}
	u, err := url.Parse(linkedin)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "Informe uma URL valida do LinkedIn."
	}
	if !strings.Contains(strings.ToLower(u.Host), "linkedin.com") {
		return "A URL precisa ser do dominio linkedin.com."
	}
	return ""
}

func hiddenUTM(utm leads.UTM) template.HTML {
	parts := []string{
		hiddenInput("utm_source", utm.Source),
		hiddenInput("utm_medium", utm.Medium),
		hiddenInput("utm_campaign", utm.Campaign),
		hiddenInput("utm_term", utm.Term),
		hiddenInput("utm_content", utm.Content),
	}
	return template.HTML(strings.Join(parts, "\n"))
}

func hiddenInput(name, value string) string {
	return fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`, template.HTMLEscapeString(name), template.HTMLEscapeString(value))
}
