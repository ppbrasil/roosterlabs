# Arquitetura — estado atual

Mantido pela skill `escopar-epico` (propõe a mudança) e verificado pelo `fechar-epico` (o doc reflete o entregue). Diagrama descreve o sistema **como está em produção** + a mudança aceita em curso, quando houver.

## Visão geral

```mermaid
flowchart LR
    V[Visitante] --> CF[CloudFront\nroosterlabs.com.br + www.\ncache de borda para GET\ninvalidado a cada deploy]
    CF -->|POST form/beacon| L[AWS Lambda\nWeb Adapter]
    CF -.viewer-request.-> WWW[CloudFront Function\n301 www -> apex]
    subgraph container [Container Go — binário único]
        L --> S[internal/server\nrotas + templates PT/EN]
        S --> LD[internal/leads\nleads + funnel events]
    end
    LD --> DB[(Neon Postgres)]
    GH[GitHub Actions\nCI verde em main] -.build+deploy+invalidate.-> L
    FONTS[Google Fonts CDN\nFraunces/Inter/Plex Mono] -.decorativo, fora do\ncaminho crítico.-> V
```

**Estado entregue:** épico 001 (landing bilíngue em produção, captura de leads, funil instrumentado) e épico 002 (copy v0.4, débito técnico de infra/segurança) — ver `epics/done/001-landing-page.md` e `epics/002-landing-v04-sync.md`.

## Rotas

| Rota | O que faz |
|---|---|
| `GET /` | landing PT-BR (copy v0.4) |
| `GET /en/` | landing EN (copy v0.4) |
| `POST /form/{step}` | processa etapa do carrossel HTMX e retorna próximo fragmento |
| `POST /event/view` | grava pageview first-party (idioma + path + UTMs lidos de `location.search` no cliente, não do HTML cacheado) |
| `GET /healthz` | health check |
| `GET /static/*` | assets embutidos no binário (inclui `og.png`/`og-en.png`, `htmx.min.js` vendorado) |

## Dados

Schema inicial no Postgres (Neon):

- `leads`: lead consolidado por token (perfil, objetivo, maturidade, desafio, email, linkedin, idioma, UTMs, timestamps).
- `funnel_events`: eventos de funil (`view`, `answer`, `submit`) com `step`, `payload` JSON e metadados de aquisição.

Migração fonte da verdade: `internal/leads/migrations/001_init.sql`.

## Detalhes entregues no épico 002 (débito técnico + sync v0.4)

- **Fontes:** Fraunces/Inter/IBM Plex Mono via Google Fonts CDN (`<link>` no `<head>`, `preconnect`) — exceção deliberada à regra de "nenhum CDN de terceiro no caminho crítico", já que é decorativo (fallback de CSS se falhar); diferente do htmx, que é vendorado localmente por ser caminho crítico do form.
- **OG image por idioma:** `ogImageFile(lang)` em `internal/server/server.go` escolhe `og.png` (PT) ou `og-en.png` (EN); metatags incluem `og:image:width=1200`/`height=630`/`type=image/png`. Fecha o G1 do épico 001 (SVG não renderizava preview em LinkedIn/WhatsApp).
- **UTM à prova de cache:** `web/static/app.js` lê `location.search` no cliente e sobrescreve incondicionalmente os hidden inputs do form (passo 1) e o payload do beacon — nunca confia no HTML servido pelo CloudFront, que cacheia ignorando query string. Fecha o G2 do épico 001.
- **Deploy CI-gated:** `deploy.yml` dispara via `workflow_run` do workflow `CI` (só com `conclusion == 'success'`), não mais por `push` direto — um lint quebrado não vai mais ao ar.
- **Cache do CloudFront invalidada a cada deploy:** último passo do `deploy.yml` roda `create-invalidation --paths "/*"` — sem isso, merges em `main` ficavam invisíveis em produção por até 24h (TTL default da cache policy). Ver `infra/README.md`.
- **Permissão dupla do Lambda Function URL codificada no Terraform:** `aws_lambda_permission.function_url_invoke` (era drift manual via CLI desde o épico 001).
- **`www.` redireciona para o apex:** via CloudFront Function (`viewer-request`, borda), não pelo Go — ver `infra/README.md`.
- **Operação sem root:** usuário IAM `ppbrasil-admin` (`AdministratorAccess` + `SignInLocalDevelopmentAccess`) substitui o uso de credenciais root da conta AWS.
