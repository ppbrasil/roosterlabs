# Backlog de engenharia

Registro de dívida técnica/infra. Itens entram nomeando **padrão violado + canal de erosão** (ver `guardrails.md`) e saem: como tarefas dentro de épicos, via épico **backpack-relief** (`workflow.md`), ou como manutenção pontual alinhada com Pedro. Dívidas de design/UX/copy têm registro equivalente em `roosterlabs-marketing`. Remover ao concluir.

## Segurança / ops — fazer logo

- [ ] **Rotacionar senha do Neon.** Exposta em texto plano numa sessão de trabalho (output de `aws lambda get-function-configuration`, 2026-07-07). Rotacionar no console do Neon → atualizar `TF_VAR_database_url` → `terraform apply`.
- [ ] **Parar de operar com root da AWS.** `aws login` hoje entrega credenciais root. Criar IAM user admin (ou Identity Center) e usar ele no dia a dia. Origem: fechamento 001.

## Infra / pipeline

- [ ] **Deploy deve esperar CI.** Workflows `CI` e `Deploy` correm em paralelo; lint quebrado foi ao ar em 2026-07-07. Unificar num workflow só (test → deploy) ou usar `workflow_run`. Origem: T11 do épico 001 previa esse bloqueio.
- [ ] **Codificar no Terraform a permissão manual do Lambda.** `lambda:InvokeFunction` com condição `lambda:InvokedViaFunctionUrl` foi adicionada via CLI (exigência AWS pós-out/2025 para Function URLs públicas). Verificar suporte do provider (`invoked_via_function_url`) e eliminar o drift. Origem: debug 403 de 2026-07-07.
- [ ] **`www.roosterlabs.com.br`.** Não resolve. Adicionar alias no CloudFront + record no Route 53 (cert wildcard já cobre). Decidir: servir direto ou redirect 301 para o apex.
- [ ] **Bump das versões de Actions.** Warnings de deprecação do Node 20 (`actions/checkout@v4`, `setup-go@v5`, `configure-aws-credentials@v4`, `golangci-lint-action@v6`).
- [ ] **Estado do Terraform é local.** Um laptop perdido = estado perdido (recuperável via import, mas doloroso). Backend S3 quando houver segundo motivo para mexer na infra.

## Código

- [ ] **Renderizar templates em buffer antes de escrever.** Erro de template hoje produz HTML parcial + `superfluous WriteHeader` (visto em produção pré-fix do `pageData`). Render em `bytes.Buffer`, escrever só se sucesso.
- [ ] **Unificar `pageData`/`formTemplateData`.** `pageData` virou superconjunto para o template embutido do form; funciona, mas é acoplamento implícito entre structs e templates. Avaliar um único view-model.

## Guardrails dos itens acima

Todos violam padrões comprometidos de engineering (`docs/conventions.md`, decisões de stack/infra) e alimentam o canal de erosão **técnico** (velocidade/confiabilidade — ver `guardrails.md`). Exceções nomeadas: senha exposta e credenciais root violam o padrão de ops seguro (canal: confiabilidade + risco existencial da conta única).
