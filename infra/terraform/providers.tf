terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source = "hashicorp/aws"
      # >= 6.28.0: primeira versão com `invoked_via_function_url` em
      # aws_lambda_permission (épico 002, T3 — ver decisions.md). Checado
      # contra o guia oficial de upgrade v5->v6: nenhum recurso usado neste
      # repo está na lista de breaking changes.
      version = ">= 6.28.0, < 7.0.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}
