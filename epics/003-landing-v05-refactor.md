# Épico 003 — Landing v0.5: hero novo + identidade Albert Sans

**Estado:** escopado <!-- proposto | escopado | aceito | em-execução | concluído -->
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

## DoD — Definition of Done (congela no aceite)

Rascunho — congela no `validar-escopo`.

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

_A preencher._

## Quebra (skill `quebrar-epico`)

_A preencher._

## Fechamento (skill `fechar-epico`)

_A preencher no fechamento._
