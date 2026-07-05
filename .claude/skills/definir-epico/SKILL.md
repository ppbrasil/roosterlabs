---
name: definir-epico
description: Define um novo épico como output e/ou outcome de negócio. Use nas sessões de strategy/marketing quando Pedro disser "novo épico", "próximo épico", "definir épico" ou selecionar um ponto de melhoria do handoff de um épico fechado.
---

# Definir Épico

Você conduz Pedro na definição de um épico — um **output e/ou outcome de negócio**, nunca uma solução técnica. O resultado é um arquivo novo em `epics/` (repo roosterlabs-engineering) no estado `proposto`, seguindo `epics/TEMPLATE.md`.

## Output vs Outcome — pragmatismo de fase inicial

- **Outcome**: resultado de negócio mensurável ("N leads/semana", "primeiro cliente pagante").
- **Output**: entrega concreta e verificável ("landing no ar em roosterlabs.com.br capturando leads").
- Em fase inicial, onde o outcome maior é "o primeiro dólar", **output é épico válido** — não force métrica de outcome onde ainda não há tráfego/clientes para medi-la. Exija pelo menos um dos dois, bem definido.

## Passos

1. Leia o contexto upstream: `goals.md`, `icp.md`, `value-proposition.md` (strategy) e, se o épico for de página/conversão, os docs de marketing. Leia também o handoff do último épico fechado em `epics/done/`, se existir.
2. Desafie a seleção: este épico serve a prioridade única (MVP + 2–3 clientes pagantes)? Se não, diga e proponha alternativas do handoff/backlog.
3. Entreviste Pedro até ter: output verificável e/ou outcome com métrica, restrições, fora-de-escopo explícito.
4. Escreva o arquivo `epics/NNN-nome.md` (próximo número livre) preenchendo apenas: cabeçalho, Output/Outcome, critério de sucesso, restrições, fora-de-escopo. Estado: `proposto`.

## Guardrails

- **Proibido** prescrever solução técnica ("usar HTMX", "criar tabela X"). Se Pedro ditar solução, registre como restrição apenas se for inegociável — e pergunte por quê.
- Critério de sucesso vago ("melhorar conversão", "site melhor") não fecha o passo 3 — exija output verificável ou número/comparador.
- Não preencha as seções de escopo, DoD, quebra ou fechamento — pertencem às etapas seguintes.
