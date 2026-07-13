# Épico 003 — Landing v0.5: hero novo + identidade Albert Sans

**Estado:** em-execução <!-- proposto | escopado | aceito | em-execução | concluído -->
**Origem:** marketing · 2026-07-09
**Prioridade servida:** MVP + 2–3 clientes pagantes (as reações de prospects que alimentam o Objetivo 1 do GTM OKR — `roosterlabs-marketing/gtm-okrs.md` — precisam ser colhidas sobre a mensagem v0.5 e a identidade que Pedro de fato aprovou, não sobre a versão anterior)

## Output / Outcome (por que este épico existe)

**Output:** `roosterlabs.com.br` (raiz PT-BR e `/en/`) servindo fielmente a copy **v0.5** (`roosterlabs-marketing/landing-page.md`) e a identidade visual **v0.5** (`roosterlabs-marketing/visual-identity.md`) — hero novo com a tagline como H1, tipografia display Albert Sans, galo como marca d'água do hero, tabela de comparação com ✓/✗ semânticos coloridos, rodapé com ícone + wordmark e OG images v0.5 por rota.

**Outcome:** nenhum novo — como o 002, este épico é pré-requisito para os KRs do Objetivo 1 medirem a mensagem real.

- **Critério de sucesso:**
  1. `GET /` e `GET /en/` servem a copy v0.5 fiel (conferência linha a linha contra `landing-page.md` v0.5): hero novo — H1 "Você. AImplificado."/"You. AImplified." **sem eyebrow**, com **apenas o prefixo "AI" em âmbar** (regra de marca: em Albert Sans, palavra inteira numa cor lê "Almplificado"), sub novo, linha de fechamento própria destacada ("O conteúdo é seu. O trabalho pesado, nosso." / "The content is yours. The heavy lifting, ours."); demais seções idênticas à v0.4 já em produção.
  2. Display em produção usa **Albert Sans carregada de verdade** (não fallback): peso 300 só em texto ≥36px, 400 abaixo disso (guardrail do peso fino, `visual-identity.md` v0.5); Fraunces sai por completo; corpo segue Inter; meta segue IBM Plex Mono.
  3. Galo como **marca d'água do hero**: opacidade 16%, integrado ao glow âmbar, atrás do texto, lado direito; nunca compete com o H1; em mobile recua para trás do CTA ou some (spec em `visual-identity.md` v0.5).
  4. Tabela de comparação com ✓ em `ok-500` (`#82B88A`) e ✗ em `no-500` (`#C96F6F`) — tokens novos da paleta v0.2.
  5. Rodapé com ícone (`roosterlabs-marketing/brand/rooster-icon-transparent.png`) + wordmark "Rooster"+"Labs" em Albert Sans 500 ("Labs" em âmbar).
  6. `og:image` v0.5 por rota (`/` → `og.png`, `/en/` → `og-en.png` — **já regenerados por marketing em 2026-07-09**, `roosterlabs-marketing/brand/web/`); preview correto no Post Inspector do LinkedIn e no WhatsApp.
  7. Title/meta description revisados onde a copy v0.5 os afeta (title mantém a tagline; og:description acompanha o novo sub — conferir contra `landing-page.md`).
- **Restrições:**
  - Empresa solo — nenhuma operação recorrente nova; custo de infra segue ≤ ~US$1/mês.
  - Copy, estrutura, lógica de conversão e identidade visual são propriedade de marketing; problema encontrado na execução sobe como nota, não vira patch local.
  - **Pré-condição de processo:** épico 002 fechado (`fechar-epico`) e arquivado em `done/` antes do aceite deste — as 19 tarefas constam concluídas, falta o fechamento formal. Verificar no fechamento do 002 as pendências que eram de Pedro (secret `CLOUDFRONT_DISTRIBUTION_ID`, push do T19).
  - Dependências de marketing **já resolvidas:** OGs v0.5 gerados; ícone com fundo transparente entregue (`brand/rooster-icon-transparent.png`, 1024px).
  - ™ de "You. AImplified" segue pendente de busca formal INPI/USPTO — não bloqueia.
- **Fora de escopo:**
  - **Favicon novo** — a redução micro sólida do ícone ainda não existe (pendência de Pedro em `visual-identity.md`); o zigue-zague atual segue. Vira item quando o asset chegar.
  - Versão pen-700 do ícone para fundos claros (mesma pendência de asset).
  - Página de sign-up (Objetivo 2), LinkedIn Ads (Objetivo 3).
  - Qualquer iteração de copy além da v0.5 aprovada — sincroniza, não redesenha.
  - Teste A/B de hero — morto por decisão de Pedro (2026-07-09, `roosterlabs-marketing/decisions.md`).

## Fontes da verdade (marketing — ler antes de escopar)

| Arquivo | O que fornece |
|---|---|
| `roosterlabs-marketing/landing-page.md` (v0.5) | Copy do hero novo PT/EN + regra do "AI" âmbar; demais seções herdadas da v0.4 |
| `roosterlabs-marketing/visual-identity.md` (v0.5) | Albert Sans + guardrail do peso fino, paleta v0.2 (ok/no), spec do watermark, família do ícone, wordmark |
| `roosterlabs-marketing/comms-foundations.md` (v0.2) | Microcopy de form/erros/estados (do/don't novos) — usar se a execução tocar essas telas |
| `roosterlabs-marketing/decisions.md` | Decisões de 2026-07-09: copy v0.5, identidade v0.5, teste A/B morto, regra do "AI" âmbar |
| `roosterlabs-marketing/brand/` | `rooster-icon-transparent.png` (rodapé), `rooster-master-fullbody.png`/`rooster-full.webp` (watermark), `web/og.png`+`og-en.png` (OG v0.5), `web/generate_og.py` (reprodutível) |

## Escopo proposto (engineering — skill `escopar-epico`)

Escopado em 2026-07-09.

**Achados da leitura do código atual:**

1. **O hero em produção não tem CTA.** A copy v0.4 especificava "CTA: Quero uma vaga no alfa", mas o `<header>` implementado no 002 tem só eyebrow + H1 + sub — a conversão sempre dependeu do visitante rolar até o form. A v0.5 volta a especificar o CTA, e a spec do watermark o pressupõe ("em mobile recua para trás do CTA"). Proposta: botão-âncora no hero rolando até `#lead-form` — zero mudança de backend, fecha a lacuna herdada.
2. **Os ✓/✗ da tabela são texto puro nos `<td>`.** Para colorir com `ok-500`/`no-500` cada célula precisa de um `<span>` com classe semântica — mudança mecânica em 12 células × 2 templates + 2 tokens novos no CSS.
3. **O ícone entregue é 1024px / 1,6 MB.** Servir isso num rodapé de ~28px é desperdício; será gerado um derivado pequeno otimizado (ex. 112px, poucos KB) para `web/static/`. O master fica em `brand/` como fonte.
4. **Para o watermark, o asset viável é `brand/web/rooster-full.webp` (26 KB)** — os masters PNG (1,1–1,4 MB) não vão para a web. Asset próprio servido de `/static/`, vendorado como o htmx (caminho visual, não CDN).
5. **Title já é a tagline** ("RoosterLabs · Você. AImplificado.") — não muda. A `meta description` de hoje é o sub v0.4; com o sub novo, ela acompanha. `og:description` (H2 do problema) fica — o problema não mudou na v0.5. **Ponto para validação de marketing no aceite.**

**Desafios de long shot:** nenhum. Tudo aqui é sync de copy/identidade com blast radius pequeno e traço direto ao Objetivo 1. Os dois riscos reais são visuais, não técnicos: peso 300 sobre ink (mitigado pelo guardrail 300-só-≥36px, já codificado no CSS) e o watermark competir com o H1 (mitigado pela spec de 16% + verificação visual no DoD).

**Decisões de engenharia propostas:**

- **Albert Sans via Google Fonts** (pesos 300, 400, 500), mesma exceção decorativa registrada no 002/T14 — se o CDN falhar, cai em fallback sans e o form segue funcionando. Fraunces sai do `<link>` e do CSS por completo (aposentada por guardrail, `visual-identity.md` v0.5).
- **Watermark via CSS** (elemento decorativo posicionado no hero, `aria-hidden`, opacidade 16%, integrado ao glow existente): em viewport estreita, media query reposiciona atrás do CTA ou o esconde, conforme spec.
- **Regra do "AI" âmbar** como `<span>` no H1 (e onde mais a tagline aparecer em Albert Sans). O title/meta (texto puro) não têm cor — sem mudança lá.
- **OGs v0.5** copiados por cima dos v0.4 em `web/static/` (mesmos nomes `og.png`/`og-en.png`): o CloudFront invalida no deploy (T19/002), mas **o cache do LinkedIn é por URL** — o fechamento precisa de re-scrape no Post Inspector, não só de olhar a metatag.

**Nota upstream (marketing):** a pendência "Transições do form (UX entre etapas)" em `landing-page.md` diz "escopo do épico novo", mas a definição do 003 não a lista em critério nenhum. **Deixada fora deste escopo** — se Pedro quiser dentro, é decisão no aceite (e o DoD ganha item + spec de UX de marketing).

### Mudança no sistema

**Antes (produção hoje, v0.4):**
- Hero: eyebrow mono "Você. AImplificado." + H1 anti-categoria ("Amplificar você não é escrever por você...") + sub v0.4, sem CTA e sem nenhuma imagem.
- Display em Fraunces (serif); tabela de comparação com ✓/✗ na cor do texto comum; rodapé só texto; favicon zigue-zague; `og.png`/`og-en.png` com a composição v0.4.
- Nenhuma ilustração na página (regra v0.4, revogada pela v0.5).

**Depois:**
- Hero: H1 = "Você. AImplificado." / "You. AImplified." com **só o "AI" em âmbar**, sem eyebrow; sub v0.5 (cadência garantida + POV real + "para quem"); linha de fechamento própria destacada ("O conteúdo é seu. O trabalho pesado, nosso."); botão CTA rolando ao form; galo em marca d'água a 16% atrás do texto, lado direito, recuando/sumindo em mobile.
- Display em Albert Sans (300 só em texto grande, 400 abaixo de 36px — guardrail do peso fino); Fraunces removida do carregamento e do CSS.
- Tabela: ✓ verdes dessaturados (`ok-500`), ✗ vermelhos dessaturados (`no-500`) — tokens novos no CSS.
- Rodapé: ícone do galo (derivado otimizado) + wordmark "Rooster**Labs**" em Albert Sans 500 com "Labs" âmbar.
- `og.png`/`og-en.png` substituídos pela composição v0.5; `meta description` acompanha o sub novo.
- **Intocado:** problema, como funciona, posicionamento (texto), acesso alfa, form (perguntas, opções, prefixo LinkedIn), rotas, dados, eventos, infra, pipeline. Nenhuma mudança de backend além de servir 2–3 assets estáticos novos.

**Atualização proposta de `docs/architecture.md`:** sem mudança de diagrama (rotas e dados não mudam); na seção de frontend, trocar a menção a Fraunces por Albert Sans (mesma exceção de CDN decorativo), registrar os assets novos em `/static/` (watermark webp, ícone do rodapé) e a nota do cache de OG do LinkedIn (re-scrape via Post Inspector após troca de imagem com mesmo nome).

## DoD — Definition of Done (congelado no aceite de 2026-07-10)

| # | Item verificável | Evidência de fechamento |
|---|---|---|
| 1 | `GET /` e `GET /en/` servem o hero v0.5 fiel (H1 tagline sem eyebrow, "AI" âmbar isolado, sub, linha de fechamento, CTA) e as demais seções idênticas à v0.4 | conferência linha a linha contra `landing-page.md` v0.5 em produção; testes de copy no CI |
| 2 | Albert Sans carrega de verdade (não fallback); 300 apenas em texto ≥36px; nenhuma referência a Fraunces no repo | fontes computadas/requisições em produção + grep |
| 3 | Watermark do galo: 16%, atrás do texto, lado direito; não compete com o H1; comportamento mobile conforme spec | verificação visual em produção, desktop + ~390px |
| 4 | ✓ em `ok-500` e ✗ em `no-500` na tabela, nas duas rotas | inspeção de cor computada em produção |
| 5 | Rodapé com ícone + wordmark ("Labs" âmbar, Albert Sans 500), asset otimizado (≤ ~30 KB servido) | inspeção visual + tamanho da resposta |
| 6 | `og:image` v0.5 por rota, preview correto | LinkedIn Post Inspector (re-scrape) + WhatsApp |
| 7 | `meta description` = sub v0.5; title inalterado; `og:description` conferido contra decisão de marketing no aceite | diff das metatags + conferência |
| 8 (transversal) | Testes no CI verdes em todo PR do épico | histórico do CI |
| 9 (transversal) | Deploy automático: merge em `main` → produção visível (invalidação T19) | merge + produção em minutos |
| 10 (transversal) | `docs/architecture.md` reflete o estado entregue | diff do doc |
| 11 (transversal) | Custo de infra segue ≤ ~US$1/mês | console AWS + Neon |

## Log de validação (skill `validar-escopo`)

- 2026-07-10 — **aprovado** (validação de Pedro em 2026-07-10; aceite registrado após o cumprimento da pré-condição — épico 002 fechado 17/17 por evidência e arquivado em `done/` na mesma data, via T21).
  - Os três pontos deferidos ao aceite confirmados por Pedro como propostos no escopo: (1) CTA âncora no hero rolando ao `#lead-form` (fecha lacuna herdada da v0.4 que a spec do watermark pressupõe); (2) `meta description` acompanha o sub v0.5, `og:description` e title inalterados; (3) transições do form (UX entre etapas) **fora** deste épico — se marketing quiser, é épico/definição à parte com spec própria.
  - Nenhum long shot identificado no escopo; mudança 100% frontend/assets, sem rotas, dados ou infra.
  - **DoD congelado nesta data (11 itens).** Alteração posterior exige novo ciclo de validação registrado neste log.
  - Próximo passo: `quebrar-epico` (engineering) — estado → `em-execução`.

## Quebra (skill `quebrar-epico`)

Quebrado em 2026-07-10. `gh` indisponível neste ambiente (mesma limitação dos épicos 001/002) — esta lista é a fonte oficial das tarefas.

**Ordem de ataque proposta: T1 → T2 → T3 → T4 → T5 → T6 → T7.**

Racional: T1 (tipografia) primeiro porque todo o resto renderiza em Albert Sans — hero, wordmark e tabela dependem dela para verificação visual honesta. T2 (hero) é a maior mudança de conteúdo e define o layout sobre o qual T3 (watermark) se posiciona — nessa ordem se evita posicionar a marca d'água duas vezes. T4 (tabela) e T5 (rodapé) são independentes entre si, mas ambas tocam `site.css` — sequenciais para não competir por merge. T6 (OG) por último antes dos docs: não colide com nada e sua verificação (re-scrape) só faz sentido com tudo em produção. T7 (docs) é sempre a última. Restrição de ambiente herdada do 002: sem toolchain Go no sandbox — `make test` roda no devcontainer de Pedro antes de cada push.

### T1 — Albert Sans no lugar de Fraunces (display real, guardrail do peso fino)
- **Comportamento observável:** H1–H3 e wordmark renderizam em Albert Sans carregada de verdade; peso 300 só onde o texto computado é ≥36px (na prática: só o H1 no clamp cheio), 400 nos demais; nenhuma referência a Fraunces no repo (link, CSS ou testes).
- **Blast radius (tocar):** `web/templates/index.html.tmpl` + `index.en.html.tmpl` (`<link>` Google Fonts: sai Fraunces, entra `Albert+Sans:wght@300;400;500`), `web/static/site.css` (`h1,h2,h3` → Albert Sans; h1 peso 300, h2/h3 peso 400), `internal/server/server_test.go` (`TestFontsLoadedFromGoogleFonts`).
- **Blast radius (ler antes):** `roosterlabs-marketing/visual-identity.md` v0.5 (guardrail do peso fino, pesos por papel), `site.css` na íntegra.
- **Plano de teste:** `<link>` com `Albert+Sans` e sem `Fraunces` nas duas rotas; `display=swap` presente; grep de "Fraunces" no repo = zero; fallback sans-serif se o CDN falhar (página funcional); h2 em peso 400 (não 300) — inspeção visual no fechamento; Inter/IBM Plex Mono intactas.
- **Traço ao DoD:** item 2.
- **Orçamento de diff:** ~60 linhas.

### T2 — Hero v0.5 PT/EN (tagline como H1 com "AI" âmbar, sub, linha de fechamento, CTA)
- **Comportamento observável:** `GET /` e `GET /en/` renderizam o hero v0.5: H1 "Você. AImplificado."/"You. AImplified." com **apenas o prefixo "AI" em âmbar** (span), sem eyebrow; sub v0.5; linha de fechamento própria destacada em paper-100; botão CTA ("Quero uma vaga no alfa"/"Get an alpha seat") rolando até `#lead-form`; `meta description` = sub v0.5; title e og:description inalterados. Demais seções intocadas.
- **Blast radius (tocar):** `web/templates/index.html.tmpl`, `index.en.html.tmpl` (header + meta description), `web/static/site.css` (`.eyebrow` sai do hero — conferir se segue usada em outro lugar antes de remover a classe; estilos do span âmbar, da linha de fechamento e do CTA), `internal/server/server_test.go` (testes de copy do hero).
- **Blast radius (ler antes):** `roosterlabs-marketing/landing-page.md` v0.5 (Hero PT/EN, racional do "AI" âmbar), testes atuais que asseguram o H1 anti-categoria e o eyebrow.
- **Plano de teste:** H1 exato por idioma com span só no "AI" (não na palavra inteira — teste de markup); eyebrow ausente; sub e linha de fechamento exatos; CTA presente com `href="#lead-form"` e o form continua com esse id; meta description nova por idioma; testes antigos do H1 v0.4 atualizados para não falsear; H2 do problema e seções seguintes intactas (asserts existentes continuam passando).
- **Traço ao DoD:** itens 1, 7.
- **Orçamento de diff:** ~150 linhas (dividir PT/EN em dois commits se estourar).

### T3 — Watermark do galo no hero
- **Comportamento observável:** galo (asset web) atrás do texto do hero, lado direito, opacidade 16%, integrado ao glow âmbar existente; nunca sobrepõe o H1; em viewport ~390px recua para trás do CTA ou some (spec v0.5).
- **Blast radius (tocar):** copiar `roosterlabs-marketing/brand/web/rooster-full.webp` (26 KB) para `web/static/`, `web/templates/index*.html.tmpl` (elemento decorativo `aria-hidden="true"` no header), `web/static/site.css` (posicionamento absoluto, opacidade, media query), `internal/server/server_test.go` (asset servido + presença do elemento).
- **Blast radius (ler antes):** `visual-identity.md` v0.5 (spec do watermark: regras de competição com H1, mobile), `site.css` (gradientes do body — o glow âmbar fica a 15%/15%, o watermark vai à direita).
- **Plano de teste:** elemento presente nas duas rotas com `aria-hidden`; `Content-Type` correto do `.webp` servido; opacidade 16% no CSS; sem overflow horizontal em 390px; texto do H1 legível sobre o watermark (verificação visual desktop+mobile no fechamento); página íntegra se o asset 404 (é `<img>` decorativa/background, não quebra layout).
- **Traço ao DoD:** item 3.
- **Orçamento de diff:** ~60 linhas (asset binário não conta).

### T4 — ✓/✗ semânticos na tabela (tokens ok/no)
- **Comportamento observável:** todos os ✓ da tabela de comparação renderizam em `ok-500` (#82B88A) e todos os ✗ em `no-500` (#C96F6F), nas duas rotas; tokens novos disponíveis no CSS mas usados só em contexto de comparação (regra da identidade).
- **Blast radius (tocar):** `web/static/site.css` (`--ok-500`/`--no-500` no `:root` + classes `.mark-ok`/`.mark-no`), `web/templates/index.html.tmpl` + `index.en.html.tmpl` (12 células ✓/✗ por página ganham span), `internal/server/server_test.go` (teste da tabela atualizado para contar spans por classe).
- **Blast radius (ler antes):** `visual-identity.md` v0.5 (paleta, regra "só semântica"), teste atual da tabela (conta ✓/✗ literais).
- **Plano de teste:** contagem exata por classe e por idioma (12 ✓ / 6 ✗ na tabela v0.5, igual à v0.4); nenhum ✓/✗ fora de span; tokens não usados em CTA/link/destaque (grep); contraste dos tons dessaturados sobre ink legível (visual no fechamento).
- **Traço ao DoD:** item 4.
- **Orçamento de diff:** ~60 linhas.

### T5 — Rodapé com ícone + wordmark
- **Comportamento observável:** rodapé das duas rotas exibe o ícone do galo + wordmark "RoosterLabs" em Albert Sans 500 — "Rooster" em paper-100, "Labs" em âmbar — mantendo contato e link de idioma atuais.
- **Blast radius (tocar):** gerar derivado otimizado do ícone (~112px, a partir de `roosterlabs-marketing/brand/rooster-icon-transparent.png` 1024px/1,6 MB) em `web/static/`, `web/templates/index*.html.tmpl` (markup do footer), `web/static/site.css` (`.site-footer`), `internal/server/server_test.go`.
- **Blast radius (ler antes):** `visual-identity.md` v0.5 (wordmark: peso 500 obrigatório em tamanho pequeno; ícone ≥28px), `.site-footer` atual no CSS.
- **Plano de teste:** asset servido ≤ ~30 KB com `Content-Type` correto; `alt` adequado no ícone; wordmark com os dois tons (span); footer íntegro em 390px (ícone + texto sem quebra feia); e-mail de contato e link de idioma preservados (asserts existentes).
- **Traço ao DoD:** item 5.
- **Orçamento de diff:** ~70 linhas (asset binário não conta).

### T6 — OG images v0.5 por rota
- **Comportamento observável:** `og.png`/`og-en.png` em `web/static/` são os arquivos v0.5 (composição com tagline-H1 e "AI" âmbar); metatags já corretas desde o 002 — sem mudança de markup.
- **Blast radius (tocar):** copiar `roosterlabs-marketing/brand/web/og.png` + `og-en.png` (v0.5, gerados 2026-07-09) por cima dos v0.4 em `web/static/`.
- **Blast radius (ler antes):** `roosterlabs-marketing/decisions.md` (spec/evidência dos OGs v0.5).
- **Plano de teste:** arquivos servidos são os novos (tamanho/hash diferem dos v0.4); 1200×630 PNG; **pós-deploy: re-scrape no LinkedIn Post Inspector** (cache de OG do LinkedIn é por URL — trocar o arquivo não basta) + preview no WhatsApp, nas duas rotas.
- **Traço ao DoD:** item 6.
- **Orçamento de diff:** ~0 linhas de código (assets binários).

### T7 — Atualizar `docs/architecture.md`
- **Comportamento observável:** o doc descreve o frontend como ficou: Albert Sans via Google Fonts (mesma exceção decorativa), assets novos em `/static/` (watermark webp, ícone do rodapé), nota do re-scrape de OG no LinkedIn.
- **Blast radius (tocar):** `docs/architecture.md`.
- **Blast radius (ler antes):** seção "Mudança no sistema" deste épico; mudanças consolidadas de T1–T6.
- **Plano de teste:** nenhuma referência a Fraunces; sem mudança de diagrama (rotas/dados intactos); instruções reproduzíveis.
- **Traço ao DoD:** item 10.
- **Orçamento de diff:** ~40 linhas.

_Itens transversais do DoD (8, 9, 11) são verificados por tarefa (CI/deploy) e no fechamento (custo)._

### Progresso de execução

- **Desvio de processo registrado (2026-07-10, decisão de Pedro):** em vez de branch+PR por tarefa, o épico inteiro roda na branch `t1-albert-sans` com **um commit por tarefa** e uma única PR ao final. Mantidos: revisor de contexto limpo por tarefa, CI verde em cada push, revisão de Pedro commit a commit na PR. Ganho: deploy atômico da v0.5 (sem estados híbridos em produção). Racional: tarefas pequenas, tráfego ~zero, empresa solo.
- **T1 — feita (2026-07-10).** Albert Sans 300/400/500 via Google Fonts nas duas rotas; `site.css` com guardrail do peso fino codificado (300 só ≥36px; media query 600px derivada do clamp do h1 — **quem fizer T2 revisita se mexer no clamp**); Fraunces zerada de link/CSS/testes. Revisor aprovou; achado dele incorporado (teste também confere o `site.css` servido). Nota do revisor: "grep Fraunces = zero" vale para referências de USO — docs/architecture.md só atualiza na T7, e citações históricas em decisions.md/épicos ficam.
- **T2 — feita (2026-07-10).** Hero v0.5 nas duas rotas: H1 tagline com span `.ai-accent` só no "AI", sem eyebrow (classe removida do CSS — uso zero), sub v0.5, linha de fechamento `.hero-close` (paper-100, destaque), CTA `.hero-cta` âncora para `#lead-form` (mesmo desenho do botão do form); meta description = sub v0.5. T2 não tocou o clamp do h1 — media query da T1 segue válida. Revisor conferiu a copy **caractere a caractere** contra `landing-page.md` v0.5 (byte a byte, PT e EN) e aprovou; achado incorporado: asserts pinando title e og:description (decisão do aceite estava sem teste). Notas do revisor registradas: (a) dois focos âmbar no hero (AI + CTA) — legítimo pela identidade (papéis: CTA + uma palavra por headline), mas a T3 soma o glow do watermark: tensão para marketing avaliar com a página pronta, não patch local; (b) `.hero-cta` espelha o botão do form, ambos sem `:hover`/`:focus-visible` — melhoria conjunta futura, não drive-by; (c) `max-width: 18ch` do h1 — conferir tagline em uma linha na inspeção visual do fechamento.
- **Nota de histórico:** o commit da T1 não chegou a acontecer no terminal de Pedro (branch criada, commit não; pego pelo revisor da T2 via `git log`). T1+T2 entraram num **único commit combinado** — separação por hunk não valia o atrito; este log é o registro por tarefa.
- **Correção de processo (2026-07-10):** o CI só roda em `pull_request` e push na `main` — push em branch sem PR não valida nada. Ajuste no desvio registrado acima: **abrir a PR como draft já no primeiro push** (CI roda em cada push da branch; draft não convida merge); no final, tirar do draft e revisar commit a commit.
- **T4 — feita (2026-07-10).** Tokens `--ok-500`/`--no-500` (hex exatos da paleta v0.2) + classes `.mark-ok`/`.mark-no`; 12✓/6✗ embrulhados em span por página (matriz conferida célula a célula pelo revisor contra `landing-page.md` v0.5, PT e EN). Teste conta por classe E por total (marca fora de span falha), e — achado do revisor incorporado — trava o número de consumidores de cada token em exatamente 1 (regra 6 da identidade: ok/no nunca em CTA/link/destaque). ~55 linhas.
- **Nota de histórico (2ª):** T1–T4 entraram num **único commit combinado** — o `git add`/`commit` de Pedro falhou duas vezes por locks órfãos (`HEAD.lock`, `index.lock`) deixados por operações de git no sandbox do Cowork; sandbox parou de rodar git no repo (lição operacional). O registro por tarefa é este log; o `make test`+`lint` local passou com T1–T4 juntas antes do commit.
- **T5 — feita (2026-07-10).** `rooster-icon.png` (112px, 12,8 KB — RMSE 3e-13 contra resize limpo do master, alpha preservado, conferido pelo revisor) + `.footer-brand` com wordmark dois tons (Albert Sans 500; peso 500 já carregado desde a T1). `alt=""` deliberado: o wordmark textual ao lado já nomeia a marca (prática WAI; comentário nos templates). Teste cobre asset ≤30 KB (o master de 1,6 MB por engano falha), Content-Type, markup dois tons e — achados do revisor incorporados — os **valores** do CSS (peso 500, Labs âmbar), não só seletores.
- **Furo de processo registrado:** o commit rotulado "T1-T4" levou a T5 junta (o `git add -A` sugerido varreu trabalho já implementado e ainda não revisado — erro de sequenciamento da sessão, pego pelo revisor da T5). Conteúdo aprovado a posteriori; correção: amend da mensagem para "T1-T5" antes do merge. Lição: bloco de commit só depois do revisor da tarefa correspondente.
- **T6 — feita (2026-07-10).** `og.png`/`og-en.png` v0.5 copiados por cima dos v0.4 (hashes trocados: `61c4c7…`→`8254d1…` PT, `034cc8…`→`8acaba…` EN; 1200×630 PNG confirmado). Metatags inalteradas (já corretas desde o 002). **Pendente para o fechamento: re-scrape no LinkedIn Post Inspector + preview WhatsApp pós-deploy** (cache de OG do LinkedIn é por URL).
- **T7 — feita (2026-07-10).** `docs/architecture.md`: diagrama (fontes), rotas (v0.5), nova seção "Detalhes entregues no épico 003" (Albert Sans + guardrail, hero v0.5, assets novos, tokens semânticos, nota do re-scrape de OG); seção do 002 enxugada com remissão. Menções a Fraunces no doc são históricas (migração), consistente com a leitura do DoD 2 registrada na T1.
- **T3 — feita (2026-07-10).** `rooster-full.webp` (26 KB, 600×600, cópia byte a byte) servido como `/static/rooster-watermark.webp`; `<img>` decorativa (`alt=""` + `aria-hidden` + `loading="lazy"` — deliberado: em mobile o `display:none` evita até o download) no header das duas rotas; CSS: 16%, direita, `z-index` negativo (atrás do texto, na frente do fundo — stacking verificado pelo revisor), glow âmbar próprio via `::before`, sem offsets negativos (sem overflow horizontal, verificado 601–1000px), some ≤600px. Revisor aprovou; achado incorporado: guarda do CSS no teste (mesmo padrão da T1). Nota do revisor para o fechamento: a verificação visual do "nunca compete com o H1" deve incluir a faixa **601–960px** (H1 e galo se tocam no limite, galo atrás a 16%), não só 390px e desktop cheio. Diff ~75 linhas (orçamento "~60" — dentro do "~", registrado).

### Emendas pós-produção (2026-07-13 — sessão de correção da landing v0.5)

Pedro inspecionou a v0.5 no ar e abriu três correções. Item 2 (favicon) e T8 (watermark) são implementação dentro da latitude de engineering; T9 (H2 âmbar) é decisão de marca que Pedro tomou **como dono** (rota "decide-agora"), com sync obrigatório abaixo.

- **Favicon (item 2) — feita (2026-07-13).** `<link rel="icon">` deixou de apontar para o `favicon.svg` provisório (zigue-zague) e passa a usar `/static/rooster-icon.png` — o mesmo ícone do rodapé (T5), por decisão de Pedro. A identidade ainda prevê uma redução micro sólida; quando existir, troca aqui (comentário nos templates).
- **T8 — emenda (2026-07-13).** Watermark deixou de ser ornamento de quina do hero e virou **backdrop da página inteira**: `.hero-watermark` (absolute, dentro do `<header>`) → `.site-watermark` (fixed, nível de `<body>`, atrás de todo o conteúdo). `position:fixed; top:40px; left:50% + translateX(-50%); width:min(2000px,170vw)` (~2x); 16%; `z-index:-1`. **Mantido em mobile** (Pedro 2026-07-13, ajuste final): o esconde ≤600px do spec v0.5 foi removido para a experiência mobile ficar igual à do desktop; como 170vw excede a largura do celular, o galo aparece ampliado (peito/corpo). Crop deliberado no desktop: cabeça/crista preservadas, **rabo (esquerda) e pernas (base) saem do quadro** — não há 2x com cabeça inteira sem cortar a base; valores (16%, top:40px) aprovados por Pedro numa prévia HTML que espelha produção. Removido o glow âmbar `::before` do hero: com o galo central, um 2º blob âmbar competia com o CTA (identidade v0.5: "um foco âmbar por tela"). Isso **resolve a tensão que o revisor da T2/T3 tinha deixado em aberto** (AI + CTA + glow do watermark = 3 focos) — agora o hero volta a 2 focos (AI + CTA).
- **T9 — emenda (2026-07-13).** Os dois primeiros H2 quebrados em duas linhas, 2ª linha em âmbar via `.h2-accent`: PT "Saber nunca foi o seu problema." / "Transformar isso em conteúdo, sim." e "Do gatilho à publicação," / "com você no centro." (EN equivalente). Decisão de Pedro como dono da marca.

> **✅ NOTA DE SYNC PARA MARKETING (RESOLVIDA em 2026-07-13):** os dois arquivos abaixo foram atualizados para **v0.5.1** na mesma sessão (`visual-identity.md` e `landing-page.md` em `roosterlabs-marketing`), eliminando a dessincronia. Registro do que divergia e foi sincronizado:
> - **`visual-identity.md` v0.5** — a regra do âmbar hoje diz "uma palavra por headline no máximo" e "um foco âmbar por tela". A T9 põe **orações inteiras** em âmbar em dois H2. A regra precisa ser reescrita conscientemente (ex.: distinguir *headline do hero* de *H2 de seção*, ou permitir 2ª-linha-âmbar como padrão de ênfase). **Atualizar também o spec do watermark:** de "16% no hero, some em mobile" para "backdrop da página inteira, fixed, ~2x, 16%, ancorado no topo, **visível também em mobile**" (a regra "some em mobile" da v0.5 foi revertida por Pedro em 2026-07-13).
> - **`landing-page.md` v0.5** — documentar a apresentação dos H2 em duas linhas com 2ª linha em âmbar (PT/EN), senão o diff de copy da próxima sync remove as quebras.
>
> Enquanto não sincronizado, `.h2-accent` no `site.css` carrega o aviso inline apontando para cá.

## Fechamento (skill `fechar-epico`)

_A preencher no fechamento._
