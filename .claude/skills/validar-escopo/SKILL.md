---
name: validar-escopo
description: Valida o escopo/DoD de um épico escopado com olhos de negócio. Use nas sessões de strategy/marketing quando Pedro disser "validar épico NNN" ou houver épico em estado escopado.
---

# Validar Escopo

Você revisa um épico `escopado` do ponto de vista de negócio e ou o aceita (DoD congela) ou devolve com pedido de revisão. O papel **não é enxugar a entrega** — coding não é o gargalo. As perguntas são: **a mudança no sistema está compreensível? O DoD, se cumprido, entrega o output/outcome?**

## Passos

1. Leia o épico inteiro — em especial a seção "Mudança no sistema" — e o contexto de strategy/marketing que o originou.
2. Cheque: (a) a mudança descrita é compreensível para quem não é dev? Se não, isso é defeito do escopo; (b) o DoD cobre o critério de sucesso (output verificável e/ou métrica)? (c) dos trade-offs listados, a recomendação sacrifica algo que o negócio não pode sacrificar? (d) algo parece long shot que engineering deixou passar?
3. Discuta os achados com Pedro. A decisão é dele.
4. Registre no "Log de validação": data + `aprovado` ou `revisão solicitada` + notas específicas (o que mudar e por quê).
5. Se aprovado: estado → `aceito`, e anote no DoD que está congelado. Se revisão: estado → `proposto`, e as notas orientam o próximo `escopar-epico`.

## Guardrails

- Não redesenhe a solução técnica — isso é loop de engineering. Aponte o problema de negócio, não a correção técnica.
- Não peça corte de escopo por princípio — só quando algo for long shot ou não servir ao output/outcome.
- Pedido de revisão sem nota específica e acionável não vale — "não gostei" bloqueia o loop.
- DoD aceito congela: qualquer mudança posterior exige novo ciclo de validação registrado.
