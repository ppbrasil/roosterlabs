# Workflow — o loop de produção homem+AI

Um loop, seis passos. Pedro decide e revisa; a AI implementa; o pipeline executa.
Regra de ouro: nenhum passo novo, agente ou ferramenta entra no processo antes de um passo existente doer ≥3 vezes.

## O loop

### 1. Spec
- Toda unidade de trabalho nasce como issue curta (GitHub Issues): problema, resultado esperado, fora de escopo.
- Copy, estrutura e conversão de páginas vêm de `roosterlabs-marketing`; requisitos de negócio vêm de `_strategy`. Engenharia **não inventa** copy nem posicionamento — se faltar, a issue bloqueia e o pedido sobe upstream.
- A AI desafia a spec antes de aceitar (escopo, complexidade, traço até a prioridade única). Se não encurta o caminho para MVP + clientes pagantes, espera.

### 2. Build
- Branch por issue. AI implementa em sessão de Claude Code, commits pequenos e frequentes.
- Alinhamento de abordagem com Pedro **antes** de codar (este arquivo e `decisions.md` são o contrato; desvios exigem conversa).
- Código segue a stack decidida: Go stdlib-first, SQL explícito via sqlc, templates server-rendered. Dependência nova = justificativa no PR.

### 3. Test
- Local/devcontainer: `go build ./... && go vet ./... && go test ./...` antes de qualquer push.
- CI (GitHub Actions) roda build + vet + test + `golangci-lint` em todo PR. PR vermelho não é revisado.
- Cobertura pragmática: toda lógica com decisão (handlers com validação, motor de extração) tem teste; template estático não precisa.

### 4. Review
- Pedro revisa o diff do PR — o gargalo do loop é este passo; código explícito existe para baratear ele.
- Roda local (`make run` / devcontainer) para ver a mudança de verdade.
- `/security-review` obrigatório quando o PR toca input de usuário, auth ou dados de cliente.
- Achou problema de copy/posicionamento durante a review? Não corrige local — sobe para marketing/strategy.

### 5. Deploy
- Merge em `main` → GitHub Actions → build da imagem → deploy no Lambda → produção. Sem passo manual.
- Rollback = redeploy do commit anterior (revert + merge).

### 6. Learn
- Dados: leads no Postgres, métricas no CloudWatch, analytics da página.
- Aprendizado vira: issue (se é engenharia) ou nota upstream em strategy/marketing (se é negócio). Nada fica só na cabeça.

## Papéis

| Quem | Faz | Não faz |
|---|---|---|
| Pedro | decide spec, revisa PR, aprova deploy implícito no merge, opera contas (GitHub/AWS/Neon) | escrever boilerplate; produzir conteúdo de cliente (regra dura da estratégia) |
| AI (Claude Code / Cowork) | desafia spec, implementa, escreve testes, prepara PR, mantém docs | mergear sem review; decidir stack/escopo sozinha; inventar copy |
| Pipeline (Actions) | test, lint, build, deploy, (futuro) monitoramento | — |

## Quando criar um agente/skill

Só quando um passo do loop repetir com atrito ≥3 vezes. Então: automatizar **aquele passo específico**, registrar em `decisions.md`, versionar em `.claude/` neste repo. Nunca criar agente "por completude".
