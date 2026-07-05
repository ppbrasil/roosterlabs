# Decision Log — Engineering

Newest first. Every entry: decision, rationale, what would reverse it.
Upstream: `_strategy/decisions.md` (solo, delivery 100% automatizada, MVP + 2–3 clientes em ~90 dias).

## 2026-07-05 — Pipeline homem+AI: loop mínimo, zero agentes novos

- **Decisão:** o processo de produção é o loop definido em `workflow.md` (spec → build → test → review → deploy → learn). Nenhum agente/skill novo é criado antecipadamente.
- **Rationale:** criar um framework de agentes antes do primeiro deploy é over-engineering de processo. Agente/skill nasce quando um passo do loop repetir e doer.
- **Reversed if:** um passo do loop mostrar atrito recorrente (≥3 ocorrências) → automatizar exatamente aquele passo.

## 2026-07-05 — Devcontainer como ambiente padrão

- **Decisão:** `.devcontainer/` na raiz define o ambiente de dev (humano e AI), igual ao CI.
- **Rationale:** reproduzibilidade máquina-nova-em-minutos; paridade dev/CI; agentes AI rodam em ambiente idêntico. Custo baixo (um JSON + Docker Desktop). Coerente com container-first.
- **Reversed if:** manutenção do container passar a custar mais que o setup manual de Go (improvável).

## 2026-07-05 — Monorepo único (`roosterlabs-engineering`)

- **Decisão:** todo código de engenharia (landing, produto, automação, infra, docs de engenharia) vive neste repo. Landing não é projeto separado: é o primeiro pacote em `internal/` e as primeiras rotas do binário.
- **Rationale:** repo-por-serviço resolve problemas de organização (times, cadências, acessos) que não temos. Monorepo dá commit atômico, um só CI/lint/convenção e contexto total para o par AI. Estratégia e marketing ficam fora (docs upstream, ciclo próprio).
- **Reversed if:** um componente precisar de ciclo de vida ou acesso próprio (open-source de uma parte, contractor com acesso parcial).

## 2026-07-05 — Cloud: AWS. Deploy interino: Lambda (Web Adapter) + CloudFront + Neon Postgres

- **Decisão:** cloud de referência é AWS. Interino (até VPC própria): o container roda em Lambda via Lambda Web Adapter, exposto por CloudFront (domínio próprio + cache de borda). Banco: Neon (Postgres serverless, free tier), substituído por RDS quando a VPC própria existir.
- **Rationale:**
  - AWS: a rede de socorro de Pedro é letrada em AWS — para solo founder, "quem me ajuda" pesa mais que elegância marginal de serviço. Conhecimento operacional transfere para a VPC futura.
  - Lambda + Web Adapter: o código continua um servidor HTTP comum (zero acoplamento ao Lambda), custo ~US$0 no free tier, escala a zero.
  - CloudFront: obrigatório para domínio próprio em Function URL; de graça vira CDN (landing servida da borda, Lambda só em POST/cache miss).
  - Neon: Postgres puro, saída = `pg_dump`. Não usar features proprietárias (branching).
  - Alternativas avaliadas: Cloud Run + Neon (venceria por simplicidade/custo, perdeu no critério rede), Railway (conveniente, mas camada intermediária com conhecimento operacional descartável), Fly.io (free tier morto, Postgres não gerenciado), Encore (ver entrada própria).
- **Restrição arquitetural aceita:** app stateless; trabalho de background = jobs explícitos disparados por scheduler/fila (nunca goroutines que sobrevivem ao request). 15 min por invocação; Durable Functions se precisar de mais.
- **Reversed if:** VPC própria de pé → ECS/Fargate + RDS, mesmo container. Cold start ou limites do Lambda doerem antes disso → App Runner/Fargate.

## 2026-07-05 — Encore: avaliado e adiado

- **Decisão:** não adotar Encore (framework nem Encore Cloud) agora.
- **Rationale:** conflita com o critério "menos caixa preta dentro do código" (primitivas e codegen do framework); Pro (US$39/mês) cria dependência invertida — infra na nossa conta, conhecimento operacional deles. Nosso app é um monolito com poucas rotas; o forte do Encore (multi-serviço, filas, eventos) não é nosso problema atual.
- **Reversed if:** o produto evoluir para múltiplos serviços com filas/eventos → reavaliar; a automação dele passa a pagar a mágica.

## 2026-07-05 — Stack: Go, monolito server-rendered, Postgres, Docker

- **Decisão:** TypeScript/Next descartados; stack é **Go** (um binário), server-rendered com `html/template` + HTMX para interatividade, acesso a dados via `sqlc` (SQL explícito tipado, sem ORM), **Postgres** único banco, **Docker** desde o commit 1.
- **Rationale:** critério de Pedro: manutenção > velocidade de escrita; menos bugs, menos caixa preta. No loop homem+AI o gargalo é a *revisão* humana de código escrito por AI — explicitude barateia a revisão; a verbosidade quem paga é a AI. Go: compilação pega classes inteiras de bugs, toolchain de manutenção forte (`vet`, `golangci-lint`), dependências mínimas, performance e custo bons. Portabilidade: container + connection string = contrato; roda igual em Lambda hoje e VPC própria amanhã. Uma linguagem serve landing, produto (motor de extração; SDKs LLM oficiais em Go) e automação — mínima diversidade, paradigma único.
- **Custos aceitos:** scaffolding inicial mais lento que frameworks web; HTMX é convenção menor que React (mitigado: UI do produto é forms/aprovações, server-rendered).
- **Reversed if:** o produto exigir UI rica de verdade (editor complexo, tempo real) → discutir front dedicado nessa hora, não antes.

## 2026-07-05 — Rascunho da landing (Cloudflare) descontinuado

- **Decisão:** o draft em `landing/` (HTML estático + Pages Functions + D1, nunca deployado) está **deprecated**. A landing será refeita como rotas do binário Go. Assets visuais (`rooster-full.webp`, `og.png`, `favicon.svg`) e aprendizados de copy serão salvos na reconstrução; a copy oficial vem de `roosterlabs-marketing`.
- **Rationale:** o rascunho embutia decisões de stack não registradas e desalinhadas com a fundação (o rabo abanando o cachorro). Estratégia → stack → implementação, nessa ordem.
- **Reversed if:** N/A. Remover `landing/` após salvar assets/copy na tarefa de reconstrução.
