# Infra — deploy interino (AWS Lambda + CloudFront + Neon)

Racional em `../decisions.md`.

## O desenho

```
visitante → CloudFront (cache borda GET + domínio roosterlabs.com.br)
              → Lambda Function URL → container Go (Lambda Web Adapter)
                                          → Neon Postgres (DATABASE_URL)
```

## IaC (Terraform)

Código em `infra/terraform/` provisiona:

- ECR (`roosterlabs-server`)
- Lambda (package image) + Function URL
- CloudFront + alias de domínio
- Route53 record de apex
- OIDC provider + role de deploy para GitHub Actions

Arquivos principais:

- `providers.tf`
- `variables.tf`
- `main.tf`
- `outputs.tf`

### Variáveis obrigatórias

- `hosted_zone_id`
- `database_url`
- `acm_certificate_arn`
- `github_owner`
- `github_repo`

## Passo-a-passo de primeira provisão

1. Configurar credenciais AWS locais (`aws configure`) com permissão de admin inicial.
2. Entrar em `infra/terraform/`.
3. Rodar `terraform init`.
4. Rodar `terraform plan -var='hosted_zone_id=...' -var='database_url=...' -var='acm_certificate_arn=...' -var='github_owner=...' -var='github_repo=...'`.
5. Rodar `terraform apply` com os mesmos `-var`.
6. Registrar output `github_actions_role_arn` como secret `AWS_DEPLOY_ROLE_ARN` no GitHub.
7. Criar secret `AWS_REGION` se quiser sobrescrever `us-east-1`.

## Deploy automático

Workflow: `.github/workflows/deploy.yml`.

Fluxo em `main`:

1. executa build + vet + test;
2. assume role via OIDC;
3. builda/pusha imagem no ECR;
4. atualiza imagem da Lambda;
5. aguarda função atualizada e imprime resumo.

## Operação de dados (Neon)

Aplicar schema inicial no SQL editor do Neon:

- `internal/leads/migrations/001_init.sql`

Queries operacionais ficam em `internal/leads/queries.sql`.

## VPC própria (futuro)

Mesmo container → ECS/Fargate; Neon → RDS (`pg_dump`/restore). CloudFront permanece.
