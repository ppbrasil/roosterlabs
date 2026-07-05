---
name: escopar-epico
description: Escopa um épico proposto — desafia, propõe implementação com trade-offs e rascunho de DoD. Use na sessão de engineering quando Pedro disser "escopar épico NNN" ou houver épico em estado proposto para trabalhar.
---

# Escopar Épico

Você transforma um épico `proposto` em `escopado`: desafio + proposta de implementação + rascunho de DoD, escritos no próprio arquivo do épico.

## Passos

1. Leia o épico, `decisions.md`, `workflow.md` e o código existente relevante. Se houver notas de revisão de um loop anterior (Log de validação), elas são o ponto de partida.
2. **Desafie antes de propor**: o outcome serve a prioridade única? O fora-de-escopo está completo? Há caminho mais curto para a mesma métrica? Registre desafios e respostas na seção "Escopo proposto".
3. Proponha o escopo de implementação como **opções com trade-offs** (mínimo 2 quando houver escolha real), com recomendação e rationale. Respeite a stack decidida; desvio de `decisions.md` exige alinhamento explícito com Pedro antes.
4. Rascunhe o DoD: itens verificáveis com evidência objetiva, cada um traçável ao outcome/métrica. Inclua sempre os itens transversais: testes passando no CI, deploy em produção, documentação atualizada.
5. Atualize o estado para `escopado` e avise Pedro que o épico aguarda validação na sessão de strategy/marketing.

## Guardrails

- Escopo que não cita trade-offs é escopo mal feito — sempre explicitar o que se perde na opção recomendada.
- Não congele o DoD — congelamento é ato do `validar-escopo`.
- Não comece a implementar nem quebrar em tarefas.
- Se faltar insumo de copy/posicionamento, marque o épico como bloqueado por upstream — não invente.
