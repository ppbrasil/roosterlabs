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
7. Registrar output `cloudfront_distribution_id` como secret `CLOUDFRONT_DISTRIBUTION_ID` no GitHub (deploy.yml invalida a cache a cada deploy — épico 002, T19).
8. Criar secret `AWS_REGION` se quiser sobrescrever `us-east-1`.

## Deploy automático

Workflow: `.github/workflows/deploy.yml`.

Fluxo em `main`:

1. executa build + vet + test;
2. assume role via OIDC;
3. builda/pusha imagem no ECR;
4. atualiza imagem da Lambda;
5. aguarda função atualizada;
6. invalida a cache do CloudFront (`/*`) — sem isso o visitante veria a versão anterior por até 24h (épico 002, T19; ver "Cache do CloudFront" abaixo);
7. imprime resumo.

## `www.` redireciona para o apex (épico 002)

`www.roosterlabs.com.br` é servido pela mesma distribuição CloudFront (alias adicional, cert ACM já é wildcard) e recebe um `A` record próprio no Route53. O redirect 301 é feito por uma **CloudFront Function** (`aws_cloudfront_function.www_redirect`, evento `viewer-request`) — **não** pelo servidor Go: o `origin_request_policy_id` (`Managed-AllViewerExceptHostHeader`, ver Gotchas abaixo) existe justamente para nunca repassar o `Host` real ao Lambda, então o Go jamais veria se a requisição veio de `www.` ou do apex. A função roda na borda, sem invocar o Lambda — mais barato e mais simples que Lambda@Edge para um redirect deste tamanho.

## Operação de dados (Neon)

Aplicar schema inicial no SQL editor do Neon:

- `internal/leads/migrations/001_init.sql`

Queries operacionais ficam em `internal/leads/queries.sql`.

## Runbook — aprendido na primeira provisão (2026-07-07)

### Credenciais AWS no devcontainer

`aws login` (browser) renova o profile, mas Terraform lê **env vars**, que expiram por conta própria. Qualquer `ExpiredToken`:

```bash
unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN AWS_CREDENTIAL_EXPIRATION
aws login
eval "$(aws configure export-credentials --format env)"
aws sts get-caller-identity   # sanity check
```

Sem o `unset`, o `export-credentials` re-exporta as env vars mortas (precedência). Rebuild do devcontainer apaga `~/.aws` → refazer `aws login`.

### Primeira provisão (chicken-and-egg da imagem)

O `aws_lambda_function` exige a imagem no ECR **antes** de existir. Bootstrap único:

```bash
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account>.dkr.ecr.us-east-1.amazonaws.com
docker build --platform linux/amd64 --provenance=false --sbom=false -t <ecr>/roosterlabs-server:latest .
docker push <ecr>/roosterlabs-server:latest
terraform apply ...   # retomará de onde parou
```

- `--platform linux/amd64`: obrigatório em máquina ARM (Lambda é x86_64; imagem errada = crash em runtime).
- `--provenance=false --sbom=false`: Lambda rejeita manifest com attestations ("image manifest, config or layer media type not supported").
- Push de `:latest` **não** atualiza a função — Lambda fixa o digest; `aws lambda update-function-code --image-uri ...` (é o que o CI faz).

### Gotchas AWS que custaram horas (não re-debugar)

| Sintoma | Causa | Fix |
|---|---|---|
| 403 `AccessDeniedException` via CloudFront | política `Managed-AllViewer` encaminha o header `Host` do visitante; Function URL rejeita Host ≠ o dela | origin request policy `Managed-AllViewerExceptHostHeader` (`b689b0a8-...`) — já no Terraform |
| 403 direto na Function URL, resource policy correta | desde out/2025 URL pública exige **duas** permissões: `lambda:InvokeFunctionUrl` E `lambda:InvokeFunction` (condição `lambda:InvokedViaFunctionUrl`) | codificado no Terraform desde o épico 002 (`aws_lambda_permission.function_url_invoke`) — se a conta já tem a permissão manual do épico 001, importar antes de aplicar (ver abaixo) |
| `Runtime.InvalidEntrypoint` / `fork/exec: permission denied` com binário e permissões corretos | distroless `:nonroot` no Lambda (workdir `/home/nonroot` inacessível ao usuário sandbox) | usar variante root — ver `decisions.md` 2026-07-07 |

### Importar a permissão manual do Lambda (épico 002, uma vez por conta existente)

Pré-requisito: provider AWS `>= 6.28.0` (ver `decisions.md`, 2026-07-07) — sem isso o `aws_lambda_permission.function_url_invoke` não aplica (`invoked_via_function_url` não existe em versões anteriores). Rodar primeiro:

```bash
terraform init -upgrade
```

Contas provisionadas antes do épico 002 têm a segunda permissão (`AllowPublicFunctionURLInvokeFunction`) criada via CLI, fora do estado. Importar antes do próximo `terraform apply` para não recriar/duplicar:

```bash
terraform import aws_lambda_permission.function_url_invoke roosterlabs-server/AllowPublicFunctionURLInvokeFunction
terraform plan   # deve sair limpo, sem diff
```

Conta nova: `terraform apply` já cria o recurso — nenhum import necessário.

**Nota (2026-07-07):** a primeira tentativa de aplicar este recurso usou `function_url_auth_type = "NONE"` em vez de `invoked_via_function_url = true` — a API da AWS rejeita `function_url_auth_type` para a ação `lambda:InvokeFunction` (`InvalidParameterValueException`). Corrigido; `main.tf` já reflete a versão certa.

### Credenciais de operação: sem root (épico 002, T2)

A conta AWS não usa mais o usuário root para operação do dia a dia. Existe um usuário IAM `ppbrasil-admin` com política `AdministratorAccess` + `SignInLocalDevelopmentAccess` (necessária para `aws login` via browser) e MFA configurado; o profile `default` do `aws login` aponta pra esse usuário. Sanity check: `aws sts get-caller-identity` deve retornar `arn:aws:iam::385129732099:user/ppbrasil-admin`, nunca a conta root. Se algum dia retornar root, pare e recrie o profile — root não deve ser usado para nada além do bootstrap inicial da conta.

### Cache do CloudFront

**Automático desde o épico 002 (T19):** todo deploy invalida `/*` como último passo (`deploy.yml`) — descoberto porque a `cache_policy_id` (`CachingOptimized`) usa TTL default de 24h na ausência de `Cache-Control` do Go, então merges em `main` ficavam invisíveis em produção por até um dia sem esse passo.

Precisa forçar manualmente (fora do fluxo de deploy normal, ex.: mudou um asset estático sem tocar código)?

```bash
aws cloudfront create-invalidation --distribution-id <ID> --paths "/*"
```

`terraform output cloudfront_distribution_id` dá o ID (também é o secret `CLOUDFRONT_DISTRIBUTION_ID` do GitHub).

Atenção: o cache key ignora query string — HTML cacheado pode carregar os UTMs do primeiro visitante entre invalidações (gap G2 do épico 001, mitigado no client-side por `app.js`, não pela invalidação).

## VPC própria (futuro)

Mesmo container → ECS/Fargate; Neon → RDS (`pg_dump`/restore). CloudFront permanece. Voltar imagem para distroless `:nonroot`.
