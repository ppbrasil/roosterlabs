---
name: quebrar-epico
description: Quebra um épico aceito em micro-tarefas com plano de teste e edge cases. Use na sessão de engineering quando Pedro disser "quebrar épico NNN" ou houver épico em estado aceito sem quebra.
---

# Quebrar Épico

Você transforma um épico `aceito` em micro-tarefas implementáveis. O DoD está congelado — a quebra serve a ele, não o reinterpreta.

## Passos

1. Leia o épico (escopo recomendado + DoD congelado) e o código existente.
2. Quebre em micro-tarefas onde cada uma: cabe num PR revisável (regra prática: diff que Pedro lê em <15 min), deixa `main` deployável ao ser mergeada, e traça a ≥1 item do DoD.
3. Para **cada** tarefa escreva, antes de qualquer código: descrição, **plano de teste com edge cases enumerados** (entradas inválidas, vazias, duplicadas, limites, concorrência, falha de dependência externa — o máximo que o contexto justificar), e o item do DoD que atende.
4. Ordene por dependência. Registre a quebra na seção "Quebra" do épico e crie as issues no GitHub (`gh issue create`), uma por tarefa. Sem GitHub disponível, a lista no épico é a fonte.
5. Estado → `em-execução`. Confirme a ordem de ataque com Pedro.

## Guardrails

- Tarefa sem plano de teste não existe — o plano é escrito na quebra, não na implementação.
- Nenhuma tarefa pode exigir mudança no DoD; se a quebra revelar problema no DoD, pare e devolva ao loop de validação.
- Não implemente nada nesta skill.
