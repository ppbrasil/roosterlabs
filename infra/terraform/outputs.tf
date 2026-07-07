output "ecr_repository_url" {
  value       = aws_ecr_repository.server.repository_url
  description = "URL do repositorio ECR"
}

output "lambda_function_name" {
  value       = aws_lambda_function.server.function_name
  description = "Nome da funcao Lambda"
}

output "cloudfront_domain_name" {
  value       = aws_cloudfront_distribution.landing.domain_name
  description = "Dominio CloudFront"
}

output "github_actions_role_arn" {
  value       = aws_iam_role.github_actions_deploy.arn
  description = "Role ARN para GitHub Actions assumir via OIDC"
}
