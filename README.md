# roosterlabs-engineering

Monorepo de engenharia da RoosterLabs. Um binário Go serve tudo: landing hoje, produto amanhã.

- **Decisões e rationale:** `decisions.md`
- **Como trabalhamos (loop homem+AI):** `workflow.md`
- **Infra/deploy:** `infra/README.md`

## Rodar local

```bash
DATABASE_URL="postgres://..." make run    # opcional; sem DATABASE_URL usa store em memória
make test   # build + vet + testes (igual ao CI)
```

Ou abra no devcontainer (VS Code: "Reopen in Container") — ambiente idêntico ao CI.

## Layout

```
cmd/server/     main do binário único
internal/       domínio (server, leads, futuro: extraction, ...)
web/            templates + estáticos (embutidos no binário)
infra/          desenho e setup do deploy
docs/           arquitetura do sistema (mantida pelo protocolo de épicos)
epics/          épicos (protocolo em workflow.md; template incluso)
.claude/        skills do pipeline + agente revisor (versionados)
landing/        DEPRECATED — rascunho antigo (Cloudflare); ver decisions.md
```

## Variáveis de ambiente

- `PORT` (default: `8080`)
- `DATABASE_URL` (opcional em dev, obrigatório em produção)
- `BASE_URL` (default: `https://roosterlabs.com.br`)
- `CONTACT_EMAIL` (default: `contact@roosterlabs.com.br`)
