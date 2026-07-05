# Decision Log — Engineering

Newest first. Every entry: decision, rationale, what would reverse it.
Upstream: `_strategy/decisions.md` (solo, delivery 100% automatizada, MVP + 2–3 clientes em ~90 dias).

## 2026-07-05 — Paradigma arquitetural: regra de dependência, sem Clean Architecture nominal

- **Decisão:** o paradigma do código é a **regra de dependência** (negócio cego para infra; setas apontam do detalhe para o negócio), instanciada nas 6 regras de `docs/conventions.md`: pacote por domínio, direção de dependência, handler fino, SQL só via sqlc no domínio dono, interface só com segundo implementador, stateless. Clean Architecture **nominal** (anéis, use cases, repositories, interfaces por padrão) foi rejeitada.
- **Rationale:** a versão nominal exige 3–4 indireções para gravar um e-mail; o ganho dessas peças só existe com múltiplas implementações ou times paralelos — não temos nenhum dos dois. As 6 regras entregam o valor real (domínio testável sem HTTP, infra trocável — nosso contrato de portabilidade) com custo mínimo de leitura/revisão, o critério dominante de Pedro. Enforcement: `revisor` + review; lint mecânico (depguard) quando violação repetir ≥3.
- **Reversed if:** segundo implementador real de um port (ex.: segundo storage), ou o monolito crescer a ponto de times/AI paralelos colidirem — aí a cerimônia adicional se paga e vira decisão nova.

## 2026-07-05 — Skills de build por camada (front/back): rejeitado

- **Decisão:** não haverá skills de build especializadas por camada técnica ou linguagem. Build permanece genérico; contenção de contexto é responsabilidade do blast radius declarado por tarefa.
- **Rationale:** a stack tem um paradigma só (Go server-rendered) — não há camada para uma skill possuir; skill por camada reintroduziria no processo a diversidade eliminada na stack, criaria decisão de roteamento por tarefa e viés contra tarefas que cruzam camadas (a maioria). O blast radius reduz contexto cirurgicamente, por tarefa. Especialização futura legítima é por **disciplina** (ex.: prompt/eval do motor de extração), governada pela regra do atrito ≥3.
- **Reversed if:** surgir camada realmente distinta (front rico em JS, mobile) ou disciplina com convenções próprias e atrito comprovado.

## 2026-07-05 — Protocolo de épicos + 5 skills + agente revisor (supersede "zero agentes novos")

- **Decisão:** o trabalho é organizado em **épicos** (outcome de negócio, arquivo em `epics/` com estados `proposto → escopado → aceito → em-execução → concluído`). O pipeline ganha peças fixas: skills `definir-epico`, `escopar-epico`, `validar-escopo`, `quebrar-epico`, `fechar-epico` e o agente `revisor` (contexto limpo), versionados em `.claude/`. Detalhe em `workflow.md`.
- **Rationale:** a decisão anterior ("zero agentes novos") subestimava o problema real apontado por Pedro: o loop cruza três projetos sem memória compartilhada. O barramento são arquivos versionados com estados; sem ritual codificado, o protocolo evapora entre sessões. Strategy/marketing definem épicos como outcome; o detalhamento nasce do loop escopo↔validação com DoD congelado no aceite; quebra em micro-tarefas com plano de teste (edge cases) antes de código; fechamento por evidência inclui atualizar a documentação do produto e notificar marketing sobre impacto em GTM. Pedro permanece roteador e gate (seleção, aceite do DoD, merge).
- **O que segue proibido:** orquestrador autônomo, comunicação agente-a-agente em runtime, agentes por especialidade técnica. Peças além destas voltam à regra do atrito ≥3.
- **Reversed if:** o protocolo custar mais que o retrabalho que evita (épicos de um dia gastando dois em cerimônia) → simplificar estados/skills.

## 2026-07-05 — [superseded pela entrada acima] Pipeline homem+AI: loop mínimo, zero agentes novos

- **Decisão:** o processo de produção é o loop definido em `workflow.md` (spec → build → test → review → deploy → learn). Nenhum agente/skill novo é criado antecipadamente.
- **Rationale:** criar um framework de agentes antes do primeiro deploy é over-engineering de processo. Agente/skill nasce quando um passo do loop repetir e doer.
- **Reversed if:** um passo do loop mostrar atrito recorrente (≥3 ocorrências) → automatizar exatamente aquele passo. *(Revertida no mesmo dia: o gatilho não era atrito repetido, e sim a constatação de que o protocolo entre projetos precisa existir antes do uso.)*

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
