package leads

import (
	"context"
	"fmt"
	"sync"
)

type memoryStore struct {
	mu      sync.Mutex
	answers map[string]map[int]map[string]string
	leads   map[string]FinalizeInput
}

// NewMemoryStore retorna uma implementação em memória para testes e dev local.
func NewMemoryStore() Store {
	return &memoryStore{
		answers: map[string]map[int]map[string]string{},
		leads:   map[string]FinalizeInput{},
	}
}

func (m *memoryStore) RecordViewEvent(_ context.Context, _ ViewEvent) error {
	return nil
}

func (m *memoryStore) RecordFormAnswer(_ context.Context, ev FormAnswerEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ev.Step < 1 || ev.Step > 4 {
		return fmt.Errorf("invalid step %d", ev.Step)
	}
	if _, ok := m.answers[ev.Token]; !ok {
		m.answers[ev.Token] = map[int]map[string]string{}
	}
	cp := map[string]string{}
	for k, v := range ev.Payload {
		cp[k] = v
	}
	m.answers[ev.Token][ev.Step] = cp
	return nil
}

func (m *memoryStore) FinalizeLead(_ context.Context, in FinalizeInput) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	steps, ok := m.answers[in.Token]
	if !ok {
		return fmt.Errorf("no answers found for token")
	}
	for i := 1; i <= 4; i++ {
		if _, ok := steps[i]; !ok {
			return fmt.Errorf("missing answer for step %d", i)
		}
	}
	m.leads[in.Token] = in
	return nil
}

func (m *memoryStore) Close() error {
	return nil
}
