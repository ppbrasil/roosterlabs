package leads

import (
	"context"
	"strings"
)

// UTM agrupa metadados de aquisição de tráfego por lead/evento.
type UTM struct {
	Source   string
	Medium   string
	Campaign string
	Term     string
	Content  string
}

// ViewEvent representa o beacon de visita de página.
type ViewEvent struct {
	Token    string
	Language string
	PagePath string
	UTM      UTM
}

// FormAnswerEvent representa a resposta de uma etapa do formulário.
type FormAnswerEvent struct {
	Token    string
	Step     int
	Language string
	PagePath string
	UTM      UTM
	Payload  map[string]string
}

// FinalizeInput contém os dados da etapa final para consolidação do lead.
type FinalizeInput struct {
	Token    string
	Language string
	PagePath string
	UTM      UTM
	Email    string
	LinkedIn string
}

// Store define o contrato de persistência para leads e eventos de funil.
type Store interface {
	RecordViewEvent(ctx context.Context, ev ViewEvent) error
	RecordFormAnswer(ctx context.Context, ev FormAnswerEvent) error
	FinalizeLead(ctx context.Context, in FinalizeInput) error
	Close() error
}

// NormalizeLanguage reduz idiomas para os dois suportados pela landing.
func NormalizeLanguage(v string) string {
	lc := strings.ToLower(strings.TrimSpace(v))
	if lc == "en" || strings.HasPrefix(lc, "en-") {
		return "en"
	}
	return "pt-BR"
}
