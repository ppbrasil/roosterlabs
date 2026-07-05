# Convenções de arquitetura e código

Lei do repo. Toda sessão (humana ou AI) herda estas regras via `CLAUDE.md`; tarefas com blast radius em código as listam como leitura obrigatória; o agente `revisor` trata violação como achado bloqueante. Rationale e alternativas rejeitadas em `decisions.md`.

O princípio único por trás de tudo: **detalhes técnicos (HTTP, Postgres, Lambda) mudam; regras de negócio não mudam quando a infra muda. Logo, código de negócio não pode saber que a infra existe. Dependências apontam do detalhe para o negócio, nunca o contrário.**

## Regra 1 — Pacote por domínio, nunca por camada técnica

```
✅ internal/leads, internal/extraction, internal/billing
❌ internal/controllers, internal/repositories, internal/services, internal/utils
```

Quem procura "tudo sobre leads" acha em um lugar. `utils`/`helpers` é onde código órfão vai apodrecer — se não tem dono de domínio, questione se deve existir.

## Regra 2 — Direção de dependência

```
cmd/server → internal/server → internal/<domínio> → (stdlib e nada mais nosso)
                    ↑
                  web/ (templates/estáticos: só internal/server os conhece)
```

```go
// ✅ internal/server/server.go
import "github.com/roosterlabs/roosterlabs-engineering/internal/leads" // detalhe→negócio

// ❌ internal/leads/leads.go
import "net/http"                                                    // negócio conhecendo HTTP
import ".../internal/server"                                         // seta invertida
```

Teste de sanidade em review: leia os imports do diff. Domínio importando `net/http`, `html/template` ou `internal/server` = bloqueia.

## Regra 3 — Handler fino

Handler traduz: parseia request → chama domínio → renderiza resposta. Decisão de negócio em handler é defeito.

```go
// ✅
func handleSubscribe(w http.ResponseWriter, r *http.Request) {
    result, err := leads.Register(r.Context(), r.FormValue("email"))
    // ... escolher template conforme result/err
}

// ❌ regra de negócio (dedupe) morando no handler
func handleSubscribe(w http.ResponseWriter, r *http.Request) {
    if exists, _ := db.LeadExists(email); exists { ... }
}
```

## Regra 4 — SQL só via sqlc, dentro do domínio dono do dado

`internal/leads/queries.sql` é o único lugar com SQL da tabela de leads; sqlc gera as funções tipadas. Outro pacote querendo um lead **pede ao pacote `leads`** — nunca consulta a tabela.

```go
// ❌ internal/extraction montando SQL de tabela alheia
rows, _ := db.Query("SELECT email FROM leads WHERE ...")
```

## Regra 5 — Interface só com segundo implementador real

Interface nasce quando existe a segunda implementação ou a necessidade concreta de fake em teste — nunca especulativa ("um dia pode ter outro banco").

```go
// ❌ cerimônia especulativa
type LeadRepository interface { Save(Lead) error }   // com UM implementador

// ✅ função concreta até o segundo implementador aparecer
func (s *Store) Save(ctx context.Context, l Lead) error
```

## Regra 6 — Stateless entre requests

Nada de estado em memória que precise sobreviver a um request (cache global mutável, sessão em map, goroutine de background que "fica rodando"). Restrição do Lambda que já aceitamos — e que vale igual na VPC futura. Trabalho de background = job explícito disparado por scheduler/fila.

## Enforcement

1. Hoje: `revisor` (imports + estas regras no diff) e review do Pedro.
2. Quando a mesma violação ocorrer ≥3 vezes: vira lint mecânico no CI (depguard/golangci) — registrar em `decisions.md`.
