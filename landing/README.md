# RoosterLabs — Landing Page

Static landing (PT-BR na raiz, EN em `/en/`) + captura de leads via Cloudflare Pages Functions + D1.
Fonte da copy: `roosterlabs-marketing/landing-page.md` (v0.3). Identidade: `roosterlabs-marketing/visual-identity.md` (v0.3).

## Estrutura

```
public/            → site estático (deploy output)
  index.html       → PT-BR
  en/index.html    → EN
  assets/          → css, js, imagens, favicon
functions/
  api/subscribe.js → POST /api/subscribe → grava lead no D1
schema.sql         → tabela leads
wrangler.toml      → config Pages + binding D1
```

## Deploy (primeira vez)

```bash
npm i -g wrangler
wrangler login
wrangler d1 create roosterlabs-leads          # copiar database_id para wrangler.toml
wrangler d1 execute roosterlabs-leads --remote --file=schema.sql
wrangler pages deploy public --project-name=roosterlabs-landing
```

Depois: domínio custom `roosterlabs.com.br` em Pages → Custom domains (DNS no Registro.br aponta para o Cloudflare).

## Deploys seguintes

```bash
wrangler pages deploy public --project-name=roosterlabs-landing
```

## Consultar leads

```bash
wrangler d1 execute roosterlabs-leads --remote --command "SELECT * FROM leads ORDER BY created_at DESC"
```

## Pendências

- E-mail de contato no footer (placeholder ausente de propósito).
- Cloudflare Web Analytics: habilitar no dashboard e colar o beacon nos dois HTML.
- Eventos de abandono por etapa do form (hoje só capturamos submissões completas).
