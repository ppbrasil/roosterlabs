# roosterlabs-engineering

Monorepo de engenharia da RoosterLabs. Um binário Go serve tudo: landing hoje, produto amanhã.

- **Decisões e rationale:** `decisions.md`
- **Como trabalhamos (loop homem+AI):** `workflow.md`
- **Infra/deploy:** `infra/README.md`

## Rodar local

```bash
make run    # servidor em http://localhost:8080
make test   # build + vet + testes (igual ao CI)
```

Ou abra no devcontainer (VS Code: "Reopen in Container") — ambiente idêntico ao CI.

## Layout

```
cmd/server/     main do binário único
internal/       domínio (server, futuro: landing, extraction, ...)
web/            templates + estáticos (embutidos no binário)
infra/          desenho e setup do deploy
landing/        DEPRECATED — rascunho antigo (Cloudflare); ver decisions.md
```
