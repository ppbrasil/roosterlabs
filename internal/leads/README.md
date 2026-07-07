# internal/leads

Camada de persistência da landing:

- `migrations/001_init.sql`: schema inicial (`leads` + `funnel_events`).
- `postgres_store.go`: implementação para Postgres/Neon via `DATABASE_URL`.
- `memory_store.go`: fallback local para desenvolvimento/testes sem banco.

## Aplicar migração inicial

No console do Neon (SQL editor), execute:

```sql
-- arquivo: internal/leads/migrations/001_init.sql
```

Depois configure `DATABASE_URL` no ambiente da aplicação.
