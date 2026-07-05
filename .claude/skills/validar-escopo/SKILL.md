---
name: validar-escopo
description: Valida o escopo/DoD de um épico escopado com olhos de negócio. Use nas sessões de strategy/marketing quando Pedro disser "validar épico NNN" ou houver épico em estado escopado.
---

# Validar Escopo

Você revisa um épico `escopado` do ponto de vista de negócio e ou o aceita (DoD congela) ou devolve com pedido de revisão. Engenharia já fez o trabalho técnico — seu papel é outro: **o DoD, se cumprido, entrega o outcome?**

## Passos

1. Leia o épico inteiro e o contexto de strategy/marketing que o originou.
2. Cheque, item a item: (a) o DoD mede o outcome ou mede atividade? (b) a métrica de sucesso está coberta? (c) o escopo recomendado sacrifica algo que o negócio não pode sacrificar (dos trade-offs listados)? (d) algo essencial do outcome ficou fora do DoD?
3. Discuta os achados com Pedro. A decisão é dele.
4. Registre no "Log de validação": data + `aprovado` ou `revisão solicitada` + notas específicas (o que mudar e por quê).
5. Se aprovado: estado → `aceito`, e anote no DoD que está congelado. Se revisão: estado → `proposto`, e as notas orientam o próximo `escopar-epico`.

## Guardrails

- Não redesenhe a solução técnica — isso é loop de engineering. Aponte o problema de negócio, não a correção técnica.
- Pedido de revisão sem nota específica e acionável não vale — "não gostei" bloqueia o loop.
- DoD aceito congela: qualquer mudança posterior exige novo ciclo de validação registrado.
