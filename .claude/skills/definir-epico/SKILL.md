---
name: definir-epico
description: Define um novo épico como outcome de negócio. Use nas sessões de strategy/marketing quando Pedro disser "novo épico", "próximo épico", "definir épico" ou selecionar um ponto de melhoria do handoff de um épico fechado.
---

# Definir Épico

Você conduz Pedro na definição de um épico — um **outcome de negócio**, nunca uma solução técnica. O resultado é um arquivo novo em `epics/` (repo roosterlabs-engineering) no estado `proposto`, seguindo `epics/TEMPLATE.md`.

## Passos

1. Leia o contexto upstream: `goals.md`, `icp.md`, `value-proposition.md` (strategy) e, se o épico for de página/conversão, os docs de marketing. Leia também o handoff do último épico fechado em `epics/done/`, se existir.
2. Desafie a seleção: este épico serve a prioridade única (MVP + 2–3 clientes pagantes)? Se não, diga e proponha alternativas do handoff/backlog.
3. Entreviste Pedro até ter: outcome observável, métrica de sucesso mensurável, restrições, fora-de-escopo explícito.
4. Escreva o arquivo `epics/NNN-nome.md` (próximo número livre) preenchendo apenas: cabeçalho, Outcome, métrica, restrições, fora-de-escopo. Estado: `proposto`.

## Guardrails

- **Proibido** prescrever solução técnica ("usar HTMX", "criar tabela X"). Se Pedro ditar solução, registre como restrição apenas se for inegociável — e pergunte por quê.
- Métrica de sucesso vaga ("melhorar conversão") não fecha o passo 3 — exija número ou comparador.
- Não preencha as seções de escopo, DoD, quebra ou fechamento — pertencem às etapas seguintes.
