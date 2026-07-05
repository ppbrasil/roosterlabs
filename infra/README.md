# Infra — deploy interino (AWS Lambda + CloudFront + Neon)

Racional em `../decisions.md`. Estado: **não provisionado ainda**.

## O desenho

```
visitante → CloudFront (cache borda, domínio roosterlabs.com.br)
              → Lambda Function URL → container (Web Adapter → servidor Go)
                                          → Neon Postgres (connection string via env)
```

## Setup primeira vez (Pedro — exige credenciais)

1. Conta AWS + usuário/role administrativo com MFA.
2. GitHub: criar repo e fazer push deste diretório.
3. AWS: OIDC provider para GitHub Actions (deploy sem secrets estáticos).
4. ECR: repositório `roosterlabs-server`.
5. Lambda: função a partir da imagem ECR + Function URL.
6. CloudFront: distribuição apontando para a Function URL; cache para GET, bypass para POST; domínio + certificado (ACM).
7. Neon: projeto free tier → `DATABASE_URL` como env do Lambda (via Secrets Manager quando houver dado real).

Passos 3–6 serão codificados (IaC) na primeira sessão de deploy — nada manual permanece.
A automação do deploy entra em `.github/workflows/` na mesma sessão.

## VPC própria (futuro)

Mesmo container → ECS/Fargate; Neon → RDS (`pg_dump`/restore). CloudFront permanece.
