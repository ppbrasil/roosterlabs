# Épico 002 — Sincronizar landing com v0.4 e fechar lacunas do 001

**Estado:** em-execução <!-- proposto | escopado | aceito | em-execução | concluído -->
**Origem:** marketing · 2026-07-07
**Prioridade servida:** MVP + 2–3 clientes pagantes (a mensagem em produção precisa ser a mensagem que Pedro de fato aprovou, sem bugs de credibilidade, para que as reações de prospects que alimentam o Objetivo 1 do GTM OKR — `roosterlabs-marketing/gtm-okrs.md` — meçam a coisa certa)
**Tipo:** output (sync v0.4) + **backpack-relief** (débitos técnicos do backlog unificado incorporados em 2026-07-07 por decisão explícita de Pedro — ver Escopo proposto)

## Output / Outcome (por que este épico existe)

**Output:** `roosterlabs.com.br` (raiz PT-BR e `/en/`) servindo fielmente a copy **v0.4** (`roosterlabs-marketing/landing-page.md`) e a identidade visual **v0.4** (`roosterlabs-marketing/visual-identity.md` — corpo em Inter, não mais Source Serif 4), sem o bug de acentuação PT hoje em produção, com preview de link (OG) correto para v0.4 e captura de UTM correta mesmo sob cache de borda. O Épico 001 fecha (11/11 itens do DoD por evidência) e arquiva em `epics/done/` quando este épico concluir.

**Outcome:** nenhum novo neste épico — ele é pré-requisito para os KRs do Objetivo 1 (10 reações qualificadas, baseline de conversão ≥8%, top 5 objeções/claims documentados) continuarem sendo coletados sobre a mensagem real, não uma versão desatualizada ou quebrada.

- **Critério de sucesso:**
  1. `GET /` e `GET /en/` servem a copy v0.4 fiel (conferência linha a linha contra `landing-page.md` v0.4, PT e EN — tagline "You. AImplified"/"Você. AImplificado.", H1 anti-categoria, H2 do problema reescrito, tabela de comparação ✓/✗, form revisado).
  2. Corpo de texto em produção usa Inter (guardrail da v0.4 disparado por leitura cansativa da serifa em fundo escuro — ver `roosterlabs-marketing/decisions.md` 2026-07-07).
  3. Acentuação PT correta em toda a página em produção (bug registrado no backlog unificado, `_strategy/backlog.md`, origem: revisão de Pedro 2026-07-07).
  4. `og:image` é PNG 1200×630 v0.4 nas duas rotas (`/` → `og.png`, `/en/` → `og-en.png`, ambos regerados com a headline v0.4); preview correto no Post Inspector do LinkedIn e no WhatsApp.
  5. UTM capturado corretamente mesmo em cache hit do CloudFront (teste com dois UTMs diferentes na mesma URL cacheada retornando valores distintos).
  6. Épico 001 tem os 11 itens do seu DoD fechados por evidência e é movido para `epics/done/`.
- **Restrições:**
  - Empresa solo — nenhuma operação recorrente nova.
  - Copy, estrutura, lógica de conversão e identidade visual continuam propriedade de marketing; problema encontrado na execução sobe como nota, não vira patch local de engenharia.
  - **Dependência upstream (marketing) — RESOLVIDA em 2026-07-07:** `og.png`/`og-en.png` regerados com a headline v0.4 (script em `roosterlabs-marketing/brand/web/generate_og.py`). Segue como tarefa de engineering apenas servir os arquivos e ajustar as metatags (ver critério de sucesso #4).
  - ™ de "You. AImplified" segue pendente de busca formal INPI/USPTO (`landing-page.md`, pendências) — não bloqueia este épico.
  - Custo de infra segue ≤ ~US$1/mês (mesmo teto do 001).
- **Fora de escopo:**
  - Página de sign-up (Objetivo 2 do GTM OKR).
  - LinkedIn Ads (Objetivo 3).
  - Qualquer iteração de copy além da v0.4 já aprovada — este épico sincroniza, não redesenha a mensagem.
  - **Migração do estado do Terraform para backend S3** — item do backlog (`_strategy/backlog.md`, "Estado do Terraform é local"); decisão de Pedro 2026-07-07: adiar. O 002 já cria o "segundo motivo para mexer na infra" que o próprio backlog cita como gatilho (item de drift do Lambda, ver Escopo proposto), mas migração de backend é trabalho de infra maior sem urgência ligada aos KRs do Objetivo 1 — fica para épico de infra à parte.
  - Demais itens do backlog unificado não puxados explicitamente por Pedro nesta sessão de escopamento (2026-07-07) — ver decisão abaixo em Escopo proposto para o que entrou.

## Fontes da verdade (marketing — ler antes de escopar)

| Arquivo | O que fornece |
|---|---|
| `roosterlabs-marketing/landing-page.md` (v0.4) | Copy integral PT/EN aprovada — fonte da verdade para o diff contra a v0.3 em produção |
| `roosterlabs-marketing/visual-identity.md` (v0.4) | Corpo Inter (era Source Serif 4); demais tokens seguem v0.3 |
| `roosterlabs-marketing/decisions.md` | Decisões de 2026-07-07: copy v0.4, identidade v0.4, OG raster |
| `_strategy/backlog.md` | Itens que este épico precisa fechar: OG raster (bloqueia G1), acentuação PT quebrada, corpo divergente do padrão, UTM à prova de cache (G2) — mais os débitos incorporados em 2026-07-07 (ver Escopo proposto): deploy×CI, drift do Lambda no Terraform, `www.` não resolver, bump de Actions, senha do Neon, root da AWS, buffer de template, unificação de `pageData`. Fica de fora: migração do backend do Terraform para S3 (adiada). |
| `epics/001-landing-page.md` (seção Fechamento) | Evidências dos 9 itens já fechados, detalhe técnico de G1 e G2, e os 6 aprendizados de engenharia já roteados ao backlog |

Nota de contexto: os PNGs de OG hoje em `roosterlabs-marketing/brand/web/` (`og.png`, `og-en.png`) ainda trazem a headline v0.3 — regeneração para v0.4 é dependência de marketing (ver Restrições acima), não bloqueia o escopamento, mas bloqueia o fechamento.

## Escopo proposto (engineering — skill `escopar-epico`)

Escopado em 2026-07-07.

**Achados da leitura do código atual (não estavam registrados em nenhuma fonte):**

1. **O "bug de acentuação PT" não é bug de pipeline.** Não há problema de encoding nem de subset de fonte sem diacrítico — o texto PT foi digitado sem acento desde a primeira versão dos templates Go (`voce`, `nao`, `consistencia`, `autentico`, etc., em `index.html.tmpl` e `form_steps.html.tmpl`). Não existe causa raiz para investigar; é reescrever o texto certo, o que a sincronização v0.4 já faz.
2. **Nenhuma fonte da identidade carrega de verdade hoje.** `site.css` declara `font-family: "Source Serif 4"`, `"Fraunces"`, `"IBM Plex Mono"`, mas não existe `<link>` de Google Fonts nem `@font-face` em nenhum arquivo. A página inteira sempre renderizou em fallback do navegador (Georgia/Times/Menlo) — nunca foi de fato Source Serif 4. É candidato a explicar o próprio guardrail que disparou a troca para Inter: o "parece que a fonte não renderizou" pode ter sido literal. **Decisão alinhada com Pedro (2026-07-07): este épico carrega as três fontes de verdade** (Fraunces, Inter, IBM Plex Mono), não só troca a declaração do corpo — sem isso o DoD #2 não fecha por evidência honesta.
3. **A seção de comparação (tabela ✓/✗) não existe no código.** A landing em produção (v0.3) não tem markup de tabela nenhum — é seção nova da v0.4, não edição de texto existente.

**Desafios de long shot:** nenhum. Os itens do épico (sync de copy, troca de fonte, OG raster, UTM à prova de cache) têm blast radius pequeno e trace direto aos KRs do Objetivo 1 — nada aqui é aposta desproporcional.

**Decisão de engenharia proposta — fontes via Google Fonts `<link>` (CDN), não vendorizadas:** a proibição de CDN de terceiro (`decisions.md`, 2026-07-07) vale para dependências do caminho crítico de conversão (htmx); tipografia é decorativa — a página funciona (e o form submete) mesmo se o Google Fonts falhar, só degrada para fallback. Consistente com a exceção já registrada em `decisions.md` ("assets decorativos podem usar CDN se houver razão").

### Débito técnico incorporado (backpack-relief — decisão de Pedro, 2026-07-07)

Pedro decidiu puxar para este épico os débitos do backlog unificado (`_strategy/backlog.md`) nascidos do fechamento do 001, além dos já cobertos acima (OG raster, acentuação, corpo tipográfico, UTM×cache). Cada item aponta o guardrail que viola, por `guardrails.md`:

**Infra/pipeline** (mesmo domínio que o 002 não tocaria de outra forma, mas blast radius pequeno cada um):
- **Deploy deve esperar CI.** `[técnica · pipeline · confiabilidade do deploy]` Hoje os workflows de CI e deploy correm em paralelo — lint quebrado já foi ao ar (2026-07-07). Fix: `deploy.yml` passa a depender do sucesso de `ci.yml` (`workflow_run` ou job unificado).
- **Codificar no Terraform a permissão manual do Lambda.** `[técnica · IaC sem drift · confiabilidade]` `lambda:InvokeFunction` (condição `InvokedViaFunctionUrl`) foi adicionada via CLI fora do Terraform. Fix: `aws_lambda_permission` equivalente em `main.tf`, importar o recurso existente (`terraform import`) para não recriar.
- **`www.roosterlabs.com.br` não resolve.** `[técnica · infra da landing · credibilidade→vendas]` Decisão de engenharia proposta: 301 para o apex (mais simples que servir os dois; cert wildcard já cobre `www.`) — alias CloudFront + record Route53.
- **Bump das versões de GitHub Actions.** `[técnica · pipeline saudável · velocidade]` Warnings de deprecação do Node 20 nos workflows; subir para a versão suportada.

**Segurança/ops** (backlog marca como "urgente — não espera épico"; **Pedro decidiu manter dentro do DoD do 002, mas como prioridade de execução** — não ficam atrás do loop de validação de copy na ordem de ataque da quebra):
- **Rotacionar senha do Neon.** `[técnica · ops seguro · confiabilidade+risco da conta única]` Exposta em texto plano em sessão de trabalho (2026-07-07). Rotacionar no console Neon → atualizar `TF_VAR_database_url`/secret → `terraform apply`.
- **Parar de operar AWS com root.** `[técnica · ops seguro · risco da conta única]` `aws login` hoje entrega credenciais root. Criar IAM user/role admin para operação do dia a dia; reservar root só para as ações que a AWS exige.

**Código** (mesmos arquivos que o 002 já vai editar por causa da sync de copy/UTM):
- **Renderizar templates em buffer.** `[técnica · conventions (handler correto) · confiabilidade]` Erro de template hoje emite HTML parcial (`ExecuteTemplate` escreve direto em `w`). Fix: renderizar em `bytes.Buffer`, só escrever na resposta (`w.Write`) se não houver erro.
- **Unificar `pageData`/`formTemplateData`.** `[técnica · conventions · velocidade de manutenção]` Os dois structs em `internal/server` têm campos quase idênticos (`Lang`, `PagePath`, `UTM`, `ContactEmail`, `Error`, `Values`) e acoplamento implícito com os templates embutidos (`pageData` existe só porque a landing embute `form_step_1_*`). Fix: um struct único ou composição explícita, eliminando a duplicação.

**Adiado (decisão de Pedro):** migração do estado do Terraform para backend S3 — ver Fora de escopo.

### Mudança no sistema

**Antes:**
- Copy PT/EN hardcoded nos templates Go em v0.3, sem acentuação PT, sem seção de comparação, com "AUTO-AUTENTICIDADE™" no eyebrow/title/meta.
- `site.css` declara Source Serif 4 / Fraunces / IBM Plex Mono sem nenhum mecanismo de carregamento — renderização real é fallback do navegador para as três.
- `og:image` aponta para `/static/og-image.svg` (mesmo arquivo nas duas rotas) — SVG não renderiza preview em LinkedIn/WhatsApp.
- UTM é capturado só a partir de valores renderizados no servidor no momento do `GET /` (hidden inputs do form embutido + `data-utm-*` do body) — sob cache hit do CloudFront, esses valores são os do primeiro visitante que gerou aquele HTML cacheado, não os do visitante atual.
- `deploy.yml` e `ci.yml` rodam em paralelo (deploy não espera CI verde); permissão `lambda:InvokeFunction` existe na AWS mas não no Terraform (drift); `www.roosterlabs.com.br` sem record/alias; Actions em versões com warning de deprecação (Node 20); `ExecuteTemplate` escreve direto no `http.ResponseWriter` (erro no meio = HTML parcial); `pageData` e `formTemplateData` duplicam campos por acoplamento implícito com os templates embutidos; senha do Neon em texto plano numa sessão passada; operação AWS via credenciais root.

**Depois:**
- Templates atualizados com a copy v0.4 completa PT/EN: tagline "Você. AImplificado." / "You. AImplified.", H1 anti-categoria, H2 do problema reescrito, nova seção "Posicionamento" (tabela ✓/✗ de 6 linhas), nova seção "Acesso alfa" com o parágrafo próprio antes do form, opções de form revisadas por pergunta (incluindo prefixo fixo `linkedin.com/in/` no campo de LinkedIn), title/eyebrow/meta atualizados para a tagline v0.4. Acentuação PT correta em todo o texto.
- `index.html.tmpl`/`index.en.html.tmpl` ganham `<link>` do Google Fonts carregando Fraunces (peso 500–600), Inter e IBM Plex Mono (peso 500); `site.css` corpo passa de Source Serif 4 para Inter. As três famílias passam a renderizar de fato, não só na declaração CSS.
- `og.png` (PT) e `og-en.png` (EN) — já gerados por marketing em `roosterlabs-marketing/brand/web/` — copiados para `web/static/`; metatags trocam de `og-image.svg` único para arquivo por rota, com `og:image:width=1200`, `og:image:height=630`, `og:image:type=image/png`.
- `app.js` passa a ler `location.search` no carregamento da página e **sobrescrever incondicionalmente** (não só "se presente" — um visitante sem UTM precisa herdar campos vazios, não o UTM cacheado de outro visitante) os 5 campos `utm_*`: os hidden inputs do `#lead-form` (passo 1) e o payload do beacon `/event/view`. Como os passos seguintes do carrossel apenas ecoam de volta o que o passo anterior recebeu (`utmFromForm`), corrigir o passo 1 no cliente basta para o funil inteiro carregar o UTM real, mesmo servido de um HTML cacheado. Cache de borda (CloudFront) não muda — a correção é só no cliente.
- Épico 001 fecha os 11 itens do DoD por evidência (G1 e G2 fecham como subproduto direto dos pontos acima) e arquiva em `epics/done/`.
- `deploy.yml` só roda após `ci.yml` verde (job unificado ou `workflow_run`); `main.tf` ganha o recurso `aws_lambda_permission` da invocação via Function URL, importado do estado real (zero drift); `www.roosterlabs.com.br` responde com 301 para o apex; workflows de Actions atualizados para versões sem deprecação; `renderLanding`/`renderFormStep` escrevem em `bytes.Buffer` antes de gravar na resposta; `pageData`/`formTemplateData` unificados num único tipo; senha do Neon rotacionada e conexão via IAM/usuário admin não-root na AWS.

**Fora da mudança:** nenhuma alteração em estrutura/lógica de conversão além do que a v0.4 já especifica; nenhuma mudança de schema de dados (`leads`/`funnel_events` já capturam as variáveis certas, só os *labels* de opção mudam); CloudFront cache key continua ignorando query string (rejeitado incluir UTM no cache key — mata hit ratio, decisão já tomada no fechamento do 001).

**Atualização proposta de `docs/architecture.md`:** nota adicional na seção de rotas/dados (sem mudança de diagrama — rotas e dados não mudam) registrando: fontes carregadas via Google Fonts CDN (decorativo, fora do caminho crítico); OG por idioma; captura de UTM corrigida no cliente antes do envio, cache de borda inalterado; deploy passa a depender de CI verde; permissão do Lambda codificada no Terraform (drift zerado).

## DoD — Definition of Done (congelado no aceite de 2026-07-07)

| # | Item verificável | Evidência de fechamento |
|---|---|---|
| 1 | `GET /` e `GET /en/` servem a copy v0.4 fiel (PT e EN) — tagline, H1, H2, seção de comparação, seção de acesso alfa, form revisado | conferência linha a linha contra `landing-page.md` v0.4 em produção |
| 2 | Corpo usa Inter; Fraunces (display) e IBM Plex Mono (meta) carregam de verdade (não fallback) | screenshot + inspeção de fontes computadas/requisições de rede em produção |
| 3 | Acentuação PT correta em toda a página em produção | inspeção do texto renderizado na URL real |
| 4 | `og:image` é PNG 1200×630 v0.4 por rota (`/`→`og.png`, `/en/`→`og-en.png`) com `width`/`height`/`type` corretos | preview no Post Inspector do LinkedIn e no WhatsApp |
| 5 | UTM capturado corretamente mesmo em cache hit do CloudFront | teste com dois UTMs diferentes na mesma URL cacheada retornando valores distintos nos eventos/lead gravados |
| 6 | Épico 001 com os 11 itens do DoD fechados por evidência, arquivado em `epics/done/` | diff do fechamento do 001 + arquivo movido |
| 7 | **[priority — executar cedo na quebra]** Senha do Neon rotacionada; AWS operada por IAM admin, não root | console Neon + `terraform apply` pós-rotação; console AWS/IAM mostrando usuário não-root em uso |
| 8 | Deploy em `main` só ocorre com CI verde (lint/test quebrado não vai ao ar) | PR de teste com lint quebrado: deploy não dispara; CI verde: dispara |
| 9 | Permissão `lambda:InvokeFunction` (Function URL) codificada no Terraform, sem drift | `terraform plan` limpo (sem diff) contra o estado real da AWS |
| 10 | `www.roosterlabs.com.br` resolve (301 para o apex) | `curl -I https://www.roosterlabs.com.br` mostrando redirect |
| 11 | Actions sem warning de deprecação (Node 20) | log do Actions sem warning |
| 12 | Erro de template não emite HTML parcial (render em buffer) | teste simulando erro de template + verificação de resposta vazia/erro limpo |
| 13 | `pageData`/`formTemplateData` unificados (um tipo, sem duplicação de campos) | diff do código mostrando struct único |
| 14 (transversal) | Testes no CI verdes (build+vet+test+lint) em todo PR do épico | histórico do CI |
| 15 (transversal) | Deploy automático: merge em `main` → produção, sem passo manual | merge de teste + mudança visível em produção |
| 16 (transversal) | `docs/architecture.md` reflete o estado entregue | diff do doc no fechamento |
| 17 (transversal) | Custo de infra segue ≤ ~US$1/mês | fatura/console AWS + Neon no fechamento |

## Log de validação (skill `validar-escopo`)

- 2026-07-07 — **aprovado**
  - Mudança no sistema compreensível para não-dev; DoD cobre os 6 critérios de sucesso do output (copy v0.4, tipografia real, acentuação, OG por idioma, UTM à prova de cache, fechamento do 001) mais os 7 itens de débito técnico incorporados nesta sessão (infra/pipeline, segurança/ops, código).
  - Nenhum long shot identificado; nenhum trade-off da recomendação sacrifica algo que o negócio não possa sacrificar.
  - **DoD congelado nesta data.** Qualquer alteração posterior exige novo ciclo de validação registrado neste log.
  - Próximo passo: `quebrar-epico` (engineering) — estado → `em-execução`.

## Quebra (skill `quebrar-epico`)

Quebrado em 2026-07-07. `gh` não está disponível neste ambiente (mesma limitação do épico 001) — esta lista é a fonte oficial das tarefas.

**Ordem de ataque confirmada (dependência + prioridade combinadas):** T1 → T2 → T3 → T4 → T5 → T6 → T7 → T8 → T9 → T10 → T11 → T12 → T13 → T14 → T15 → T16 → T17 → T18.

Racional da ordem: T1/T2 primeiro por decisão explícita de Pedro (segurança/ops não fica atrás do loop de copy). T3–T6 (infra/pipeline) agrupados antes do código de app para não competir por revisão com o trabalho de conteúdo. T7/T8 (buffer render + unificação de structs) antes de qualquer tarefa de conteúdo porque ambas tocam as mesmas funções de renderização que T9–T13 vão usar — fazer depois duplicaria retrabalho. T9→T13 em sequência porque todas editam os mesmos arquivos de template (`index*.html.tmpl`, `form_steps.html.tmpl`) e cada uma parte do estado deixado pela anterior. T14/T15 (fontes, OG) depois do conteúdo para não colidir com os merges de copy. T16 (UTM) fecha os dois débitos técnicos (G1 via T15, G2 via T16) que T17 precisa para fechar o 001. T18 é sempre a última.

### Progresso de execução (2026-07-07)

- **T3 — feito (código).** `aws_lambda_permission.function_url_invoke` codificado em `main.tf` + passo de `terraform import` documentado em `infra/README.md`. **Pendente de Pedro:** rodar o import + `terraform apply` (credenciais AWS).
- **T4 — feito.** `deploy.yml` agora dispara via `workflow_run` do `CI`, só com `conclusion == 'success'`; checkout e tag de imagem usam `head_sha` do run de CI, não `github.sha`.
- **T5 — feito (código).** Descoberta durante a implementação: o plano original (redirect no Go via `Host`) não funciona — o `origin_request_policy_id` atual existe para nunca repassar o `Host` real ao Lambda (evita 403 da Function URL). Trocado para uma **CloudFront Function** (`viewer-request`, borda, sem invocar o Lambda) — ver nota na tarefa T5 acima. **Pendente de Pedro:** `terraform apply` (credenciais AWS).
- **T6 — feito, com ressalva.** `actions/checkout` e `actions/setup-go` bumpados (Node24, confirmados por busca). `golangci-lint-action` **mantido em v6/golangci-lint v1.64.8** de propósito: `v7+` exige migrar para `golangci-lint v2` (regras mais estritas, config nova) — fora do orçamento desta tarefa. Fica em Node20 até essa migração ser decidida à parte. Adicionado `.github/dependabot.yml` (github-actions, semanal) para não reabrir esse gap no futuro.
- **Restrição de ambiente:** este ambiente de execução não tem toolchain Go nem Terraform instalados (rede restrita a allowlist, sem privilégio para instalar) — T3/T5 (Terraform) e T7 em diante (Go) não puderam ser validados aqui com `terraform validate`/`make test`. Necessário rodar `make test` e `terraform validate`/`plan` no devcontainer antes de qualquer push.
- **T1/T3 — aplicados e confirmados em produção.** `terraform apply` rodou limpo depois de duas correções encontradas só na execução real (fora do que o escopo previa): (1) provider AWS precisou subir de `~> 5.0` para `>= 6.28.0, < 7.0.0` — o argumento `invoked_via_function_url` do `aws_lambda_permission` só existe a partir dessa versão (checado o guia de upgrade v5→v6 linha a linha contra os recursos deste repo, nenhum breaking change aplicável); (2) o argumento certo é `invoked_via_function_url = true`, não `function_url_auth_type` (a API da AWS rejeita esse último para a ação `InvokeFunction`). Ambas registradas em `decisions.md`. Rotação da senha do Neon confirmada por evidência real (submissão de teste do form + query no Neon), não só `/healthz`.
- **T7 — achado adicional durante a implementação.** `ensureToken` (grava o cookie de sessão `rlsid`) rodava DEPOIS de `ExecuteTemplate(w, ...)` escrever direto na resposta em `renderLanding` — como isso já dispara o `WriteHeader` implícito do Go, o `Set-Cookie` nunca saía no `GET /`/`GET /en/` (só nos passos do form, que já chamavam `ensureToken` antes de qualquer escrita). Corrigido como parte do mesmo fix de buffer (`ensureToken` agora roda antes de qualquer header/write); teste de regressão `TestIndexSetsSessionCookie` adicionado.
- **T8 — feito.** `formTemplateData` removido; `pageData` é agora o único tipo de dado de template, usado por landing e todos os passos do form.
- **T2 — feito.** `pedro-admin` criado no IAM (`AdministratorAccess` + `SignInLocalDevelopmentAccess`, MFA configurado), profile `default` do `aws login` migrado — `aws sts get-caller-identity` confirma `arn:aws:iam::385129732099:user/ppbrasil-admin`, não mais root.
- **T1/T2/T3/T4/T5/T6/T7/T8 — todos fechados.**
- **T9 — feito.** Hero, problema, como funciona, footer, title/meta em v0.4 (PT/EN) em `index.html.tmpl`/`index.en.html.tmpl`. Testes antigos que buscavam "AUTO-AUTENTICIDADE"/"AUTO-AUTHENTICITY" atualizados para checar a copy v0.4 e a ausência da copy retirada.
- **T10 — feito.** Seção "Posicionamento" (tabela ✓/✗, 6 linhas) adicionada nas duas páginas + CSS (`compare-table`/`compare-scroll`, com scroll horizontal interno em vez de estourar a página em mobile). Teste novo confere as 6 linhas e a contagem exata de ✓/✗ por idioma.
- **T11 — feito.** Seção "Acesso alfa" (H2 + parágrafo) adicionada antes do form embutido, PT/EN.
- **T12 — feito.** Opções de form atualizadas nas duas línguas: PT perde "fractional" no passo 1 (só EN mantém, por especificação v0.4); "Construir audiência" (PT)/"Build an audience" (EN) novo no passo 2; "Só comento e reajo a posts dos outros" (PT)/"I only react and comment on others' posts" (EN) novo no passo 3. Passo 6 (pós-envio) com a redação v0.4 exata. `TestFormFlowToSuccess` e `TestFormOptionsMatchV04` atualizados/adicionados.
- **T13 — feito.** Campo de LinkedIn no passo 5 trocado de URL completa para handle com prefixo fixo `linkedin.com/in/` (`.input-prefix-group`/`.input-prefix` no CSS). `validateFinal` usa a nova função `normalizeLinkedInURL` (reconstrói a URL canônica a partir do handle; aceita defensivamente que o usuário cole a URL inteira, extraindo só o handle). `net/url` removido do import (não usado mais). Testes de tabela cobrindo os edge cases do plano (handle simples, vazio, URL colada, barra final, query string colada, espaços).
- **T14 — feito.** Fontes carregadas via Google Fonts CDN (`Fraunces`/`Inter`/`IBM Plex Mono`, `preconnect` + `<link rel="stylesheet">`) nas duas páginas; `site.css` com `font-family: "Inter", Arial, sans-serif` no body (era serif genérico, não batia com o CSS que já referenciava as três famílias). Exceção de CDN de terceiro documentada inline no template (decorativo, fora do caminho crítico — diferente do htmx, vendorado). Teste `TestFontsLoadedFromGoogleFonts` novo.
- **T15 — feito.** `og.png`/`og-en.png` (gerados por marketing em `roosterlabs-marketing/brand/web/`) copiados para `web/static/`, substituindo o SVG que não gerava preview no LinkedIn/WhatsApp (fecha G1 do épico 001). Novo helper `ogImageFile(lang)` em `server.go` escolhe o PNG certo por idioma; `og:image:width`/`height`/`type` adicionados nas duas páginas. `og-image.svg` antigo não pôde ser removido (permissão do filesystem montado) — fica como arquivo morto inofensivo, remoção manual opcional. Teste novo `TestOGImagePerLanguage` cobre metatags corretas por idioma e `Content-Type: image/png` do arquivo servido.
- **T16 — feito.** `app.js` agora lê `location.search` diretamente (via `URLSearchParams`) em vez de confiar no dataset do `<body>` — que é HTML servido pelo CloudFront e pode ter sido cacheado pra outro visitante (a política de cache ignora query string). Sobrescreve incondicionalmente (mesmo para `""`) os 5 hidden inputs do form no passo 1 e usa os mesmos valores no payload do beacon. Passos 2+ herdam o valor corrigido porque `utmFromForm` (form_handler.go) só ecoa o que foi submetido no passo anterior — não precisou tocar o form_handler nem os templates. Sem teste automatizado: como já previsto no plano da tarefa, T16 é smoke manual + revisão de código (não há infra de teste JS no repo; criar uma só para isso seria over-engineering fora do orçamento).
- **T19 — aplicado.** Ao tentar fechar o épico 001 por evidência, `roosterlabs.com.br` ainda servia v0.3 mesmo com o Lambda já atualizado — `CachingOptimized` cacheia a página por até 24h sem `Cache-Control` do Go. Adicionado passo de `create-invalidation --paths "/*"` no fim do `deploy.yml` + permissão IAM + output `cloudfront_distribution_id`. `terraform apply` rodado com sucesso em 2026-07-08: permissão `cloudfront:CreateInvalidation` criada, `cloudfront_distribution_id = E217UBG541YE9Q`. **Pendente de Pedro:** registrar secret `CLOUDFRONT_DISTRIBUTION_ID` no GitHub + commit/push do código do T19 (main.tf/outputs.tf/deploy.yml/infra/README.md/decisions.md ainda não versionados). T17 só pode fechar depois do próximo deploy rodar com a invalidação.
- **T17 — feito.** Épico 001 fechado por evidência real em 2026-07-08 (LinkedIn Post Inspector + teste ao vivo de UTM no Chrome) — ver `epics/done/001-landing-page.md`.
- **T18 — feito.** `docs/architecture.md` reescrito (diagrama + seção "Detalhes entregues no épico 002"); `infra/README.md` corrigido (linha desatualizada sobre redirect www via Go — na real é CloudFront Function — e tabela de Gotchas que tinha ficado quebrada com linha órfã) + nova seção de runbook sobre operação sem root (T2). Grep confirmou nenhuma referência a root/senha antiga do Neon.
- Todas as 19 tarefas concluídas.

### Progresso de execução (2026-07-09 — tentativa de fechamento)

- **Pendências de Pedro do T19 — confirmadas resolvidas.** Secret `CLOUDFRONT_DISTRIBUTION_ID` registrado e funcionando: Deploy #10 (commit 63f9114) rodou com sucesso incluindo o passo "Invalidate CloudFront cache"; código do T19 commitado/pushado (2668bbb).
- **Achado no fechamento — DoD 11 falha por evidência.** O log de Actions ainda tinha **dois** warnings de deprecação Node 20: (1) `aws-actions/configure-aws-credentials@v4` no deploy — não coberto pelo T6; (2) `golangci/golangci-lint-action@v6` no CI — a exceção deliberada do T6 (documentada em `decisions.md` + ignore no Dependabot), que o texto congelado do DoD 11 não reflete. Geradas duas tarefas-lacuna:
- **T20 — feito.** PR #1 do Dependabot (`configure-aws-credentials` 4→6) mergeada em 2026-07-09 com autorização de Pedro, após verificação dos breaking changes de v5 (input handling booleano — não usamos) e v6 (Node 24 — runners hosted ok) contra o nosso uso (só `role-to-assume` + `aws-region`). Evidência de fechamento: Deploy #11 (2026-07-09, commit 9de4f29) — Success, zero annotations (o #10 tinha 1 warning). De quebra, reconfirma DoD 8 (disparo via CI verde) e DoD 15 (pipeline + invalidação).
- **T21 — tarefa-lacuna aberta (bloqueia DoD 11 e o fechamento do épico).** Migração golangci-lint v1.64.8 → v2 + `golangci-lint-action` v6 → atual, único caminho para zerar o warning restante. **Decisão de Pedro (2026-07-09): segurar o fechamento do 002 até esta migração** — oferecida a alternativa de fechar com exceção documentada via emenda ao DoD; recusada. Ao concluir, remover o `ignore` do `.github/dependabot.yml` (ver `decisions.md` 2026-07-08).
- **T21 — executada em 2026-07-10.** Migração menor que o orçado: **não existe `.golangci.yml` no repo** — nada de config para migrar; o default "standard" do v2 equivale ao default do v1 (gosimple absorvido pelo staticcheck). Mudanças: `ci.yml` → `golangci-lint-action@v9` + `version: v2.12.2`; `dependabot.yml` → ignore removido. O v2 pagou a entrada: apontou 1 `errcheck` real (`postgres_store.go` — `defer rows.Close()` sem tratamento), corrigido com descarte explícito comentado (`rows.Err()` já cobria iteração). `make test` + `make lint` verdes no devcontainer com v2.12.2. Evidência final do DoD 11 (log do CI sem warning) confirmada no push de fechamento.
- Evidências dos demais itens do DoD colhidas em 2026-07-09 e anotadas na seção Fechamento (pré-preenchida); épico permanece **em-execução** até T21.

### T1 — [Pedro] Rotacionar senha do Neon
- **Comportamento observável:** credencial do Postgres trocada; nenhuma sessão/cliente usa a senha antiga.
- **Blast radius (tocar):** console Neon (fora do repo); secret `DATABASE_URL` no GitHub Actions / `TF_VAR_database_url` local.
- **Blast radius (ler antes):** `infra/README.md` (operação de dados), `infra/terraform/variables.tf`.
- **Plano de teste:** conexão com a senha antiga falha após a rotação; conexão com a senha nova funciona (`curl /healthz` pós-deploy); `terraform apply` não gera diff além da var.
- **Traço ao DoD:** item 7.
- **Orçamento de diff:** 0 (ação operacional); nota no `infra/README.md` se o procedimento mudar.

### T2 — [Pedro] Eliminar uso de root da AWS
- **Comportamento observável:** operação do dia a dia (deploy, terraform, debug) usa usuário/role IAM admin, nunca root.
- **Blast radius (tocar):** console IAM (fora do repo); `infra/README.md` (runbook de credenciais).
- **Blast radius (ler antes):** `infra/README.md` (seção "Credenciais AWS no devcontainer").
- **Plano de teste:** `aws sts get-caller-identity` mostra ARN de usuário IAM, não `root`; o fluxo de `aws login` documentado no runbook usa o novo usuário.
- **Traço ao DoD:** item 7.
- **Orçamento de diff:** ~20 linhas (runbook em `infra/README.md`).

### T3 — Codificar no Terraform a permissão do Lambda (zerar drift)
- **Comportamento observável:** `terraform plan` não mostra diff para a permissão `lambda:InvokeFunction` (condição `InvokedViaFunctionUrl`) — ela passa a existir no estado, não só na AWS.
- **Blast radius (tocar):** `infra/terraform/main.tf` (novo `aws_lambda_permission`), `infra/README.md` (passo de `terraform import`).
- **Blast radius (ler antes):** `infra/README.md` (tabela de gotchas), `decisions.md` (entrada do container distroless/Lambda).
- **Plano de teste:** `terraform validate` limpo; `terraform import` do recurso existente não recria a permissão; `terraform plan` pós-import sem diff; em ambiente novo (sem o recurso), `terraform apply` cria a permissão corretamente.
- **Traço ao DoD:** item 9.
- **Orçamento de diff:** ~40 linhas.

### T4 — Deploy só ocorre com CI verde
- **Comportamento observável:** `deploy.yml` não roda (ou aborta) se `ci.yml` não passou para o commit em `main`.
- **Blast radius (tocar):** `.github/workflows/deploy.yml`, `.github/workflows/ci.yml`.
- **Blast radius (ler antes):** os dois workflows na íntegra.
- **Plano de teste:** commit com lint quebrado em `main` — deploy não dispara ou aborta; commit com CI verde — deploy dispara normalmente; sem brecha para deploy manual contornar o gate sem intenção explícita.
- **Traço ao DoD:** itens 8, 14.
- **Orçamento de diff:** ~60 linhas.

### T5 — `www.roosterlabs.com.br` redireciona 301 para o apex
- **Comportamento observável:** requisição a `www.roosterlabs.com.br` (qualquer path/query) responde 301 para o mesmo path/query em `https://roosterlabs.com.br`.
- **Correção durante a implementação (fora do blast radius original, registrada aqui):** o plano original tocava `internal/server.go` (redirect por `Host` no Go). Na implementação isso se mostrou impossível com a infra atual: o `origin_request_policy_id` do CloudFront (`Managed-AllViewerExceptHostHeader`) existe justamente para **nunca** repassar o `Host` real do visitante ao Lambda — é o que evita o 403 da Function URL (`infra/README.md`). O servidor Go, portanto, nunca veria se a requisição veio de `www.` ou do apex. Trocado para uma **CloudFront Function** (`viewer-request`, roda na borda antes da origin request policy, sem invocar o Lambda) — mais barato e mais simples que Lambda@Edge para um redirect deste tamanho. Nenhuma mudança em `internal/server/`.
- **Blast radius (tocar):** `infra/terraform/main.tf` (novo `aws_cloudfront_function.www_redirect` + `function_association` no `default_cache_behavior` + `www.` nos `aliases` do CloudFront + `aws_route53_record.www`).
- **Blast radius (ler antes):** `infra/terraform/main.tf` (cert já é wildcard; origin request policy existente).
- **Plano de teste:** `Host: www.roosterlabs.com.br` + `/en/?utm_source=x` → 301 para `https://roosterlabs.com.br/en/?utm_source=x` (path e query preservados, inclusive múltiplos parâmetros); `Host: roosterlabs.com.br` → sem redirect; TLS ok em `www.` (cert wildcard); `terraform validate`/`plan` limpo.
- **Traço ao DoD:** item 10.
- **Orçamento de diff:** ~60 linhas (só infra, sem código Go).

### T6 — Bump das versões de GitHub Actions
- **Comportamento observável:** workflows rodam sem warning de deprecação (Node 20).
- **Blast radius (tocar):** `.github/workflows/ci.yml`, `.github/workflows/deploy.yml`.
- **Blast radius (ler antes):** changelog das actions em uso (`checkout`, `setup-go` etc.).
- **Plano de teste:** execução sem warning no log; build/test/deploy continuam passando sem mudança de comportamento.
- **Traço ao DoD:** item 11.
- **Orçamento de diff:** ~20 linhas.

### T7 — Renderizar templates em buffer
- **Comportamento observável:** erro de execução de template nunca chega ao cliente como HTML parcial — resposta é 500 limpo ou sucesso completo.
- **Blast radius (tocar):** `internal/server/server.go` (`renderLanding`), `internal/server/form_handler.go` (`renderFormStep` e o render final do passo 6), `internal/server/server_test.go`.
- **Blast radius (ler antes):** as duas funções na íntegra.
- **Plano de teste:** template forçado a falhar no meio (nome inexistente injetado em teste) → resposta vazia + 500, nunca HTML cortado; render bem-sucedido → Content-Type e corpo idênticos ao comportamento atual.
- **Traço ao DoD:** item 12.
- **Orçamento de diff:** ~90 linhas.

### T8 — Unificar `pageData`/`formTemplateData`
- **Comportamento observável:** um único tipo serve landing e form steps; templates recebem os mesmos campos de antes, sem duplicação de struct.
- **Blast radius (tocar):** `internal/server/server.go`, `internal/server/form_handler.go`, `internal/server/events_handler.go` (se usar o tipo), `internal/server/server_test.go`.
- **Blast radius (ler antes):** todos os `.html.tmpl` que referenciam campos de `pageData`/`formTemplateData`.
- **Plano de teste:** todos os testes existentes de renderização (index, en, form flow, validação) continuam verdes sem mudar asserção; nenhum campo usado em template deixa de existir no tipo unificado.
- **Traço ao DoD:** item 13.
- **Orçamento de diff:** ~110 linhas.

### T9 — Copy v0.4 core (hero, problema, como funciona, footer, title/meta) — PT/EN
- **Comportamento observável:** `GET /` e `GET /en/` renderizam tagline, H1, H2 do problema e "como funciona" na redação v0.4; title/eyebrow/meta refletem a tagline v0.4 (sem "Auto-Autenticidade").
- **Blast radius (tocar):** `web/templates/index.html.tmpl`, `web/templates/index.en.html.tmpl`, `internal/server/server_test.go`.
- **Blast radius (ler antes):** `roosterlabs-marketing/landing-page.md` v0.4 (Hero, Problema, Como funciona, Footer), testes atuais que buscam "AUTO-AUTENTICIDADE".
- **Plano de teste:** tagline exata PT/EN; H1 com a resolução certa por idioma (PT "É refletir você." / EN "It's writing from you."); testes antigos atualizados para não falsear passar; nenhuma referência residual a "Auto-Autenticidade".
- **Traço ao DoD:** item 1.
- **Orçamento de diff:** ~150 linhas (dividir PT/EN se exceder).

### T10 — Seção "Posicionamento" (tabela ✓/✗) — PT/EN
- **Comportamento observável:** `GET /` e `GET /en/` renderizam a tabela de comparação de 6 linhas com os valores exatos da v0.4.
- **Blast radius (tocar):** `web/templates/index.html.tmpl`, `web/templates/index.en.html.tmpl`, `web/static/site.css`, `internal/server/server_test.go`.
- **Blast radius (ler antes):** `roosterlabs-marketing/landing-page.md` v0.4 (Posicionamento), `visual-identity.md` (spotlight/translucidez).
- **Plano de teste:** 6 linhas × 3 colunas presentes; ✓/✗ exatos por célula; sem overflow horizontal em viewport ~390px (verificação visual no fechamento).
- **Traço ao DoD:** item 1.
- **Orçamento de diff:** ~130 linhas.

### T11 — Seção "Acesso alfa" — PT/EN
- **Comportamento observável:** `GET /` e `GET /en/` renderizam a seção de acesso alfa (H2 + parágrafo) antes do form embutido.
- **Blast radius (tocar):** `web/templates/index.html.tmpl`, `web/templates/index.en.html.tmpl`.
- **Blast radius (ler antes):** `roosterlabs-marketing/landing-page.md` v0.4 (Acesso alfa).
- **Plano de teste:** texto exato PT/EN presente; ordem (acesso alfa → form) preservada; `id="lead-form"` continua presente logo em seguida.
- **Traço ao DoD:** item 1.
- **Orçamento de diff:** ~50 linhas.

### T12 — Opções de form v0.4 — PT/EN (passos 1–4 e mensagem final)
- **Comportamento observável:** cada pergunta do carrossel oferece exatamente as opções da v0.4 (novas: "Construir audiência"/"Build an audience" no P2; "Só comento e reajo..."/"I only react and comment..." no P3); PT remove "fractional" do P1; passo 6 usa a redação v0.4.
- **Blast radius (tocar):** `web/templates/form_steps.html.tmpl`, `internal/server/form_handler.go` (`validateStep`/`validateFinal`, se necessário), `internal/server/server_test.go`.
- **Blast radius (ler antes):** `roosterlabs-marketing/landing-page.md` v0.4 (Form, PT e EN), `TestFormFlowToSuccess` (usa valores literais de `choice`).
- **Plano de teste:** cada passo aceita as opções novas; `choice` antigo removido (`fractional` no PT-P1) não é mais esperado nos testes; `TestFormFlowToSuccess` atualizado com os novos valores continua passando; validação de "Outro" continua funcionando nos passos certos.
- **Traço ao DoD:** item 1.
- **Orçamento de diff:** ~100 linhas.

### T13 — Campo de LinkedIn com prefixo fixo `linkedin.com/in/`
- **Comportamento observável:** usuário digita só o handle; servidor grava/valida a URL completa (`https://www.linkedin.com/in/<handle>`).
- **Blast radius (tocar):** `web/templates/form_steps.html.tmpl` (passo 5, PT/EN), `internal/server/form_handler.go` (`validateFinal`/consolidação), `internal/server/server_test.go`.
- **Blast radius (ler antes):** `internal/leads/types.go`/`store.go` (formato esperado do campo LinkedIn hoje).
- **Plano de teste:** handle simples válido; handle vazio rejeitado; usuário cola URL completa por engano — normaliza em vez de duplicar prefixo; handle com caracteres especiais/espaços; handle com barra final.
- **Traço ao DoD:** item 1.
- **Orçamento de diff:** ~80 linhas.

### T14 — Carregamento real das fontes (Fraunces, Inter, IBM Plex Mono)
- **Comportamento observável:** `GET /` e `GET /en/` carregam as três fontes via Google Fonts; corpo do texto usa Inter de fato (não fallback).
- **Blast radius (tocar):** `web/templates/index.html.tmpl`, `web/templates/index.en.html.tmpl` (`<link>` no `<head>`), `web/static/site.css` (`font-family` do `body`).
- **Blast radius (ler antes):** `roosterlabs-marketing/visual-identity.md` v0.4 (pesos: Fraunces 500–600, IBM Plex Mono 500).
- **Plano de teste:** `<link>` presente com os três `family=` corretos e `display=swap`; página não quebra se o Google Fonts falhar (fallback do CSS continua); `site.css` não referencia mais "Source Serif 4" no corpo.
- **Traço ao DoD:** item 2.
- **Orçamento de diff:** ~40 linhas.

### T15 — OG assets por idioma + metatags
- **Comportamento observável:** `/` serve `og:image` = `og.png`; `/en/` serve `og:image` = `og-en.png`; ambos com `og:image:width`/`height`/`type`.
- **Blast radius (tocar):** copiar `roosterlabs-marketing/brand/web/og.png` e `og-en.png` para `web/static/`, `internal/server/server.go` (`OGImageURL` por idioma), `web/templates/index*.html.tmpl` (metatags), `internal/server/server_test.go`.
- **Blast radius (ler antes):** `roosterlabs-marketing/decisions.md` (spec do asset, 2026-07-07).
- **Plano de teste:** `og:image` de `/` aponta para `og.png`, de `/en/` para `og-en.png`; `width`/`height`/`type` corretos; arquivo servido com `Content-Type: image/png`; fallback documentado se o asset faltar.
- **Traço ao DoD:** item 4.
- **Orçamento de diff:** ~60 linhas (assets binários não contam).

### T16 — UTM à prova de cache (client-side)
- **Comportamento observável:** `location.search` da requisição atual sempre vence — hidden inputs do form (passo 1) e payload do beacon usam o UTM real do visitante, mesmo em cache hit do CloudFront.
- **Blast radius (tocar):** `web/static/app.js`.
- **Blast radius (ler antes):** `web/templates/form_steps.html.tmpl` (nomes dos 5 campos hidden), `internal/server/server.go` (`utmFromURL`), `internal/server/form_handler.go` (`utmFromForm` — confirma que os passos seguintes só ecoam o passo 1).
- **Plano de teste (smoke manual + revisão de código):** URL com todos os UTMs → hidden inputs e beacon usam esses valores, não o dataset renderizado; URL sem UTM → os 5 campos ficam vazios mesmo com dataset cacheado de outro visitante; URL com só 1 dos 5 UTMs → os outros 4 ficam vazios; comportamento idêntico em `/` e `/en/`.
- **Traço ao DoD:** item 5.
- **Orçamento de diff:** ~50 linhas.

### T19 — Invalidação de cache do CloudFront no deploy (achado durante T17)
- **Comportamento observável:** `deploy.yml` invalida `/*` na distribuição CloudFront como último passo; conteúdo novo visível em produção logo após o pipeline terminar, sem esperar o TTL default (24h) da cache policy.
- **Blast radius (tocar):** `infra/terraform/main.tf` (permissão IAM `cloudfront:CreateInvalidation`), `infra/terraform/outputs.tf` (output `cloudfront_distribution_id`), `.github/workflows/deploy.yml`, `infra/README.md`.
- **Blast radius (ler antes):** `decisions.md` (entrada de 2026-07-08).
- **Achado:** ao tentar verificar produção pro T17 (fechamento do épico 001), `roosterlabs.com.br` ainda servia a copy v0.3 mesmo com o Lambda já atualizado — `cache_policy_id` = `CachingOptimized`, TTL default 24h, Go nunca manda `Cache-Control`. Sem invalidação automática, todo deploy fica "no ar" no Lambda mas invisível ao visitante por até um dia.
- **Plano de teste:** Pedro registra o secret `CLOUDFRONT_DISTRIBUTION_ID`, `terraform apply` aplica a permissão IAM nova, próximo merge em `main` dispara o deploy e a invalidação; produção reflete a mudança em minutos, confirmado por `curl`/fetch.
- **Traço ao DoD:** item 6 (pré-requisito para fechar o épico 001 por evidência real).
- **Orçamento de diff:** ~40 linhas.

### T17 — Fechamento do épico 001 (skill `fechar-epico`)
- **Comportamento observável:** os 11 itens do DoD do épico 001 verificados por evidência em produção (G1 via T15, G2 via T16); arquivo movido para `epics/done/`.
- **Blast radius (tocar):** `epics/001-landing-page.md` (seção Fechamento, G1/G2 de pendente para fechado), mover para `epics/done/001-landing-page.md`.
- **Blast radius (ler antes):** `epics/001-landing-page.md` inteiro.
- **Plano de teste:** cada um dos 11 itens com evidência nova ou reconfirmada; nenhum item fecha "por sensação".
- **Traço ao DoD:** item 6.
- **Orçamento de diff:** renomeação + texto, sem código.

### T18 — Atualizar `docs/architecture.md` e `infra/README.md`
- **Comportamento observável:** os dois docs refletem o estado entregue (fontes via CDN, OG por idioma, UTM corrigido no cliente, deploy dependente de CI, permissão do Lambda no Terraform, runbook sem root).
- **Blast radius (tocar):** `docs/architecture.md`, `infra/README.md`.
- **Blast radius (ler antes):** mudanças consolidadas de T1–T17, seção "Mudança no sistema" deste épico.
- **Plano de teste:** instruções reproduzíveis; nenhuma referência a root/senha antiga do Neon; nenhuma rota/dado desatualizado.
- **Traço ao DoD:** item 16.
- **Orçamento de diff:** ~100 linhas.

### T21 — Migração golangci-lint v2 (tarefa-lacuna do fechamento)
- **Comportamento observável:** CI roda sem warning de deprecação Node 20; lint continua bloqueando PR vermelho.
- **Blast radius (tocar):** `.golangci.yml` (config nova do v2), `.github/workflows/ci.yml` (`golangci-lint-action` → versão atual), `.github/dependabot.yml` (remover o `ignore` de `golangci-lint-action >= 7`), código Go que as regras mais estritas do v2 apontarem.
- **Blast radius (ler antes):** guia oficial de migração v1→v2 do golangci-lint; `decisions.md` (entradas de 2026-07-08 sobre o ignore e o T6).
- **Plano de teste:** `make test` + lint local verdes no devcontainer; CI verde no PR; log do CI sem warning de Node 20; um erro de lint proposital ainda quebra o CI (gate preservado).
- **Traço ao DoD:** item 11.
- **Orçamento de diff:** ~60 linhas de config/workflow + fixes de lint que surgirem (se explodirem, parar e reavaliar com Pedro).

## Fechamento (skill `fechar-epico`)

**Estado em 2026-07-09: fechamento bloqueado por decisão de Pedro — DoD 11 pendente da T21 (migração golangci v2).** Todas as demais evidências já colhidas e registradas abaixo; ao fechar T21, reconferir só o item 11 e concluir.

| # | Item | Evidência (2026-07-09, salvo indicação) |
|---|---|---|
| 1 | Copy v0.4 fiel | Fetch de produção `/` e `/en/`: tagline, H1 anti-categoria, H2, tabela ✓/✗ 6 linhas, seção alfa, form v0.4 (PT sem "fractional"; EN com) — batem com a v0.4; testes de copy exata no CI verdes. `landing-page.md` upstream já avançou para v0.5 (épico 003) — conferência linha a linha feita contra os fragmentos v0.4 congelados neste épico + testes |
| 2 | Fontes reais | `document.fonts.check` em produção: Inter ✓, IBM Plex Mono ✓, Fraunces w600 ✓ (w500 não carregado por lazy-load — nenhum elemento o usa); `body` computado = `Inter, Arial, sans-serif` |
| 3 | Acentuação PT | Texto renderizado em produção correto ("Você", "consistência", "própria", etc.) |
| 4 | OG PNG por rota | Metatags em produção: `/`→`og.png`, `/en/`→`og-en.png`, `width/height/type` corretos; preview LinkedIn Post Inspector + WhatsApp evidenciado no fechamento do 001 (2026-07-08) |
| 5 | UTM × cache | Teste ao vivo: `/?utm_source=fech002&utm_medium=chrome&utm_campaign=dod5` → hidden inputs exatamente esses valores, `utm_term`/`utm_content` vazios (não herdou cache); idem em `/en/` |
| 6 | 001 fechado/arquivado | `epics/done/001-landing-page.md` (commits 40e2495, 63f9114) |
| 7 | Neon rotacionada; AWS sem root | Evidência registrada 2026-07-07/08: submissão de teste + query no Neon pós-rotação; `aws sts get-caller-identity` = `user/ppbrasil-admin` |
| 8 | Deploy só com CI verde | Deploy #10 "Triggered via workflow run" com `conclusion == success`; PR do golangci com CI vermelho fechada sem disparar deploy |
| 9 | Permissão Lambda no TF sem drift | `terraform plan` em 2026-07-10: **"No changes. Your infrastructure matches the configuration."** Achado no caminho: o único drift era estrutural — `image_uri` (`:latest` no TF vs. tag por SHA que o `deploy.yml` aponta a cada deploy); resolvido com `lifecycle { ignore_changes = [image_uri] }` no `aws_lambda_function` (TF = bootstrap; pipeline = dono da imagem), o que também elimina o risco de um `apply` re-apontar produção para tag mutável |
| 10 | www → 301 apex | Navegação real: `www.roosterlabs.com.br/en/?utm_source=www301&utm_medium=redirect` → apex `/en/` com os dois UTMs intactos nos hidden inputs (path + query preservados) |
| 11 | **PENDENTE — T21** | Warning do deploy morto pela T20 (PR #1 mergeada); warning do CI (`golangci-lint-action@v6`) só morre com a migração v2 |
| 12 | Render em buffer | `TestRenderLandingTemplateErrorDoesNotWritePartialBody` + `TestRenderFormStepTemplateErrorDoesNotWritePartialBody` no repo; CI verde |
| 13 | Struct único | `formTemplateData` ausente do código (grep); só `pageData` |
| 14 | CI verde em PRs | Histórico: 13 runs, falha só na PR do Dependabot que era para ser fechada |
| 15 | Deploy automático | Deploy #10 + invalidação CloudFront → produção servindo v0.4 minutos após merge |
| 16 | `docs/architecture.md` atualizado | T18 commitado em 2026-07-09 (estava só no working tree) |
| 17 | Custo ≤ ~US$1/mês | Fechado 2026-07-10 por leitura direta dos consoles: AWS Bills julho/2026 = **US$ 0,57** (Route 53 US$ 0,50; resto free tier) + Neon = **US$ 0** (plano Free, 1,85 CU-hrs, 32 MB). Achado colateral: billing era invisível para IAM users ("Activate IAM Access" desativado — leitura exigia root, contra o espírito do T2); IAM Access ao billing ativado em 2026-07-10 (console mostra "Activated") — checagens futuras de custo dispensam root |

_Notas de roteamento (strategy/marketing/issues) e handoff: a preencher no fechamento final, pós-T21._
