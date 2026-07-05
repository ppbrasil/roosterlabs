---
name: escopar-epico
description: Escopa um épico proposto — dá transparência sobre a mudança no sistema (docs/diagramas), trade-offs e rascunho de DoD. Use na sessão de engineering quando Pedro disser "escopar épico NNN" ou houver épico em estado proposto para trabalhar.
---

# Escopar Épico

Você transforma um épico `proposto` em `escopado`. **Coding não é o gargalo — o objetivo do escopo não é enxugar a entrega, é dar transparência sobre qual será a mudança no sistema.** O entregável central é Pedro (e strategy/marketing) entenderem exatamente o que vai mudar antes de qualquer código.

## Passos

1. Leia o épico, `decisions.md`, `workflow.md`, `docs/architecture.md` e o código existente relevante. Se houver notas de revisão de um loop anterior (Log de validação), elas são o ponto de partida.
2. Desafie **apenas o que for long shot**: aposta desproporcional ao retorno, ou sem traço à prioridade única (MVP + 2–3 clientes pagantes). Não corte escopo por reflexo; registre desafios e respostas.
3. Proponha o escopo como opções com trade-offs quando houver escolha real, com recomendação e rationale. Respeite a stack decidida; desvio de `decisions.md` exige alinhamento explícito com Pedro antes.
4. **Escreva a seção "Mudança no sistema"** — o coração do escopo: antes → depois de componentes, rotas e dados, e a atualização proposta de `docs/architecture.md` (diagramas mermaid + prosa). Prepare a atualização do doc junto; ela é mergeada com o épico aceito.
5. Rascunhe o DoD: itens verificáveis com evidência objetiva. Inclua sempre os transversais: testes no CI, deploy em produção, `docs/architecture.md` refletindo o entregue.
6. Estado → `escopado`; avise Pedro que aguarda validação na sessão de strategy/marketing.

## Guardrails

- Escopo cuja "Mudança no sistema" um leitor não-técnico não acompanha está mal escrito — transparência é o critério de qualidade número um.
- Não congele o DoD — congelamento é ato do `validar-escopo`.
- Não comece a implementar nem quebrar em tarefas.
- Se faltar insumo de copy/posicionamento, marque o épico como bloqueado por upstream — não invente.
