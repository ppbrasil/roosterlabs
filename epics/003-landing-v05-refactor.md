# Épico 003 — Landing v0.5: hero novo + identidade Albert Sans

**Estado:** proposto <!-- proposto | escopado | aceito | em-execução | concluído -->
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

_A preencher por engineering._

### Mudança no sistema

_A preencher por engineering._

## DoD — Definition of Done (congela no aceite)

_A preencher no escopo/aceite._

## Log de validação (skill `validar-escopo`)

_A preencher._

## Quebra (skill `quebrar-epico`)

_A preencher._

## Fechamento (skill `fechar-epico`)

_A preencher no fechamento._
