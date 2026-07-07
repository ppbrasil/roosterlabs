variable "project_name" {
  type        = string
  description = "Prefixo dos recursos AWS"
  default     = "roosterlabs"
}

variable "aws_region" {
  type        = string
  description = "Regiao AWS"
  default     = "us-east-1"
}

variable "domain_name" {
  type        = string
  description = "Dominio principal da landing"
  default     = "roosterlabs.com.br"
}

variable "hosted_zone_id" {
  type        = string
  description = "Hosted zone no Route53 para o dominio"
}

variable "database_url" {
  type        = string
  description = "DATABASE_URL do Neon"
  sensitive   = true
}

variable "acm_certificate_arn" {
  type        = string
  description = "ARN do certificado ACM em us-east-1 para o dominio"
}

variable "github_owner" {
  type        = string
  description = "Owner no GitHub para OIDC"
}

variable "github_repo" {
  type        = string
  description = "Repositorio no GitHub para OIDC"
}
