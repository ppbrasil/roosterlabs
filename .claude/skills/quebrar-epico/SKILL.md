---
name: quebrar-epico
description: Quebra um épico aceito em unidades mínimas com blast radius declarado, plano de teste e edge cases. Use na sessão de engineering quando Pedro disser "quebrar épico NNN" ou houver épico em estado aceito sem quebra.
---

# Quebrar Épico

Você transforma um épico `aceito` em micro-tarefas implementáveis. O DoD está congelado — a quebra serve a ele, não o reinterpreta. Objetivo: **unidades tão mínimas quanto possível para conter alucinação**, com desvios detectáveis.

## A unidade mínima

Uma tarefa = **um comportamento observável** (ou um passo preparatório), nunca "e também". Não é "um arquivo": um comportamento em Go legitimamente toca handler + registro de rota + template + teste — esse conjunto É a unidade mínima. Forçar um arquivo por tarefa cria estados intermediários que não compilam e quebra "main sempre deployável".

## Passos

1. Leia o épico (Mudança no sistema + DoD congelado) e o código existente.
2. Quebre em tarefas onde cada uma: entrega um comportamento, cabe num diff que Pedro lê em <15 min, deixa `main` deployável ao ser mergeada, e traça a ≥1 item do DoD.
3. Para **cada** tarefa escreva, antes de qualquer código:
   - **Blast radius declarado**: arquivos/funções que a tarefa espera tocar E arquivos que deve ler antes de começar.
   - **Plano de teste com edge cases enumerados**: entradas inválidas, vazias, duplicadas, limites, concorrência, falha de dependência externa — o máximo que o contexto justificar.
   - **Traço ao DoD**: qual item atende.
   - **Orçamento de diff**: ~150 linhas (excluindo código gerado e testdata); estimativa acima disso = quebrar mais.
4. Ordene por dependência. Registre a quebra na seção "Quebra" do épico e crie as issues no GitHub (`gh issue create`), uma por tarefa. Sem GitHub disponível, a lista no épico é a fonte.
5. Estado → `em-execução`. Confirme a ordem de ataque com Pedro.

## Guardrails

- Tarefa sem plano de teste ou sem blast radius não existe — ambos são escritos na quebra, não na implementação.
- Na implementação, precisar tocar fora do blast radius = **parar e atualizar a tarefa** (e reavaliar a quebra) — nunca seguir em silêncio. O `revisor` verifica isso.
- **Proibido drive-by**: refactor/melhoria fora do comportamento vira tarefa própria, mesmo com a mão no arquivo.
- Nenhuma tarefa pode exigir mudança no DoD; se a quebra revelar problema no DoD, pare e devolva ao loop de validação.
- Não implemente nada nesta skill.
