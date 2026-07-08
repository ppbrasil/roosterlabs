locals {
  tags = {
    Project = var.project_name
    Managed = "terraform"
  }
}

resource "aws_ecr_repository" "server" {
  name                 = "${var.project_name}-server"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = local.tags
}

resource "aws_iam_role" "lambda_exec" {
  name = "${var.project_name}-lambda-exec"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
      Action = "sts:AssumeRole"
    }]
  })

  tags = local.tags
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "server" {
  function_name = "${var.project_name}-server"
  package_type  = "Image"
  role          = aws_iam_role.lambda_exec.arn
  image_uri     = "${aws_ecr_repository.server.repository_url}:latest"
  timeout       = 30
  memory_size   = 512

  environment {
    variables = {
      DATABASE_URL  = var.database_url
      BASE_URL      = "https://${var.domain_name}"
      CONTACT_EMAIL = "contact@roosterlabs.com.br"
    }
  }

  tags = local.tags
}

resource "aws_lambda_function_url" "server" {
  function_name      = aws_lambda_function.server.function_name
  authorization_type = "NONE"
}

resource "aws_lambda_permission" "function_url_public" {
  statement_id           = "AllowPublicFunctionURLInvoke"
  action                 = "lambda:InvokeFunctionUrl"
  function_name          = aws_lambda_function.server.function_name
  principal              = "*"
  function_url_auth_type = "NONE"
}

# Desde out/2025 a AWS exige as DUAS permissões para Function URL pública:
# InvokeFunctionUrl (acima) E InvokeFunction (abaixo), condicionada a
# lambda:InvokedViaFunctionUrl. Foi adicionada via CLI na primeira provisão
# (drift — ver decisions.md e backlog); este recurso a codifica. Em conta
# nova, `terraform apply` já cria as duas; em conta existente, importar
# antes de aplicar (ver infra/README.md) para não recriar/duplicar.
resource "aws_lambda_permission" "function_url_invoke" {
  statement_id              = "AllowPublicFunctionURLInvokeFunction"
  action                    = "lambda:InvokeFunction"
  function_name             = aws_lambda_function.server.function_name
  principal                 = "*"
  invoked_via_function_url  = true
}

# Redirect 301 www -> apex na borda. Não dá para fazer isso no Go: o
# origin_request_policy abaixo (Managed-AllViewerExceptHostHeader) existe
# para NUNCA repassar o Host real do visitante ao Lambda (é o que evita o
# 403 da Function URL — ver infra/README.md), então o servidor jamais veria
# se a requisição veio de "www." ou do apex. CloudFront Function roda antes
# dessa policy, na borda, sem invocar o Lambda — mais barato e mais simples
# que Lambda@Edge para um redirect deste tamanho.
resource "aws_cloudfront_function" "www_redirect" {
  name    = "${var.project_name}-www-redirect"
  runtime = "cloudfront-js-2.0"
  comment = "301 www.${var.domain_name} -> ${var.domain_name}, path e query preservados"
  publish = true
  code    = <<-EOT
    function handler(event) {
      var request = event.request;
      var host = request.headers.host.value;
      if (host.substring(0, 4) !== 'www.') {
        return request;
      }
      var apex = host.substring(4);
      var keys = Object.keys(request.querystring);
      var qs = '';
      if (keys.length > 0) {
        var parts = [];
        for (var i = 0; i < keys.length; i++) {
          var key = keys[i];
          var entry = request.querystring[key];
          if (entry.multiValue) {
            for (var j = 0; j < entry.multiValue.length; j++) {
              parts.push(encodeURIComponent(key) + '=' + encodeURIComponent(entry.multiValue[j].value));
            }
          } else {
            parts.push(encodeURIComponent(key) + (entry.value !== undefined ? '=' + encodeURIComponent(entry.value) : ''));
          }
        }
        qs = '?' + parts.join('&');
      }
      return {
        statusCode: 301,
        statusDescription: 'Moved Permanently',
        headers: {
          location: { value: 'https://' + apex + request.uri + qs }
        }
      };
    }
  EOT
}

resource "aws_cloudfront_distribution" "landing" {
  enabled             = true
  default_root_object = ""
  # www. entra na mesma distribuição (cert ACM já é wildcard); o servidor
  # Go responde com 301 para o apex baseado no Host (ver internal/server) —
  # mais simples que uma segunda distribuição só para redirect.
  aliases = [var.domain_name, "www.${var.domain_name}"]

  origin {
    domain_name = trimsuffix(replace(aws_lambda_function_url.server.function_url, "https://", ""), "/")
    origin_id   = "lambda-url-origin"

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS", "PUT", "PATCH", "POST", "DELETE"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "lambda-url-origin"

    cache_policy_id = "658327ea-f89d-4fab-a63d-7e88639e58f6"
    # Managed-AllViewerExceptHostHeader: function URLs reject requests whose
    # Host header differs from the URL's own domain, so CloudFront must NOT
    # forward the viewer's Host (Managed-AllViewer does, and breaks with 403).
    origin_request_policy_id = "b689b0a8-53d0-40ab-baf2-68738e2966ac"
    viewer_protocol_policy   = "redirect-to-https"

    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.www_redirect.arn
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = var.acm_certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }

  tags = local.tags
}

resource "aws_route53_record" "apex" {
  zone_id = var.hosted_zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.landing.domain_name
    zone_id                = aws_cloudfront_distribution.landing.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "www" {
  zone_id = var.hosted_zone_id
  name    = "www.${var.domain_name}"
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.landing.domain_name
    zone_id                = aws_cloudfront_distribution.landing.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = ["sts.amazonaws.com"]

  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1"
  ]

  tags = local.tags
}

resource "aws_iam_role" "github_actions_deploy" {
  name = "${var.project_name}-github-actions-deploy"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Federated = aws_iam_openid_connect_provider.github.arn
      }
      Action = "sts:AssumeRoleWithWebIdentity"
      Condition = {
        StringEquals = {
          "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
        }
        StringLike = {
          "token.actions.githubusercontent.com:sub" = "repo:${var.github_owner}/${var.github_repo}:*"
        }
      }
    }]
  })

  tags = local.tags
}

resource "aws_iam_role_policy" "github_actions_deploy" {
  name = "${var.project_name}-github-actions-deploy"
  role = aws_iam_role.github_actions_deploy.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:PutImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "lambda:UpdateFunctionCode",
          "lambda:GetFunction",
          "lambda:GetFunctionConfiguration"
        ]
        Resource = aws_lambda_function.server.arn
      },
      {
        # Sem isso, todo deploy atualiza o Lambda mas o visitante continua
        # vendo a versao anterior por ate 24h (TTL default da cache_policy
        # CachingOptimized, ja que o Go nao manda Cache-Control) — gap achado
        # na verificacao de producao do epico 002 (T19). create-invalidation
        # com "/*" no deploy.yml fecha isso sem passo manual.
        Effect = "Allow"
        Action = [
          "cloudfront:CreateInvalidation"
        ]
        Resource = aws_cloudfront_distribution.landing.arn
      }
    ]
  })
}
