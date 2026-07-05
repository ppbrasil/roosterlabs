---
name: fechar-epico
description: Fecha um épico por evidência — verifica DoD, atualiza documentação, roteia aprendizados e notifica marketing sobre impacto em GTM. Use na sessão de engineering quando Pedro disser "fechar épico NNN" após a última tarefa deployada.
---

# Fechar Épico

Você fecha um épico `em-execução` cuja hipótese é "DoD cumprido". Fechamento é **por evidência, não por sensação**.

## Passos

1. **Verifique o DoD item a item com prova real**: teste a URL de produção, rode a query, meça o número. Cole a evidência ao lado de cada item na seção "Fechamento". Item sem evidência possível = item mal escrito; discuta com Pedro antes de dar por cumprido.
2. Se algum item falhar: crie as tarefas-lacuna (mesmo formato do `quebrar-epico`), mantenha o estado `em-execução` e pare aqui.
3. **Atualize toda a documentação do produto afetada**: README, `infra/`, docs de domínio, comentários de arquitetura. O repo deve descrever o sistema como ele ficou, não como era.
4. Roteie aprendizados, escrevendo na seção "Fechamento":
   - Engenharia → issues no repo.
   - Negócio (ICP, willingness to pay, comportamento) → nota endereçada a **strategy**.
   - **Impacto em GTM (páginas, copy, funil, canais) → nota explícita endereçada a marketing**, com o que mudou e o que pode exigir ajuste. Se não houver impacto, registre "nenhum".
5. Estado → `concluído`; mova o arquivo para `epics/done/`.
6. Escreva o **handoff**: o que foi entregue, evidências-chave, aprendizados, candidatos a próximo ponto de melhoria — pronto para Pedro levar à sessão de strategy/marketing e recomeçar o ciclo com `definir-epico`.

## Guardrails

- Nunca marque item de DoD sem evidência colada.
- Não selecione o próximo épico — isso é decisão de strategy/marketing com Pedro; você só deixa a decisão pronta.
- As notas para strategy/marketing ficam no épico (o barramento é o arquivo); avise Pedro de que existem.

## Nota sobre documentação

O passo 3 (atualizar documentação) inclui obrigatoriamente `docs/architecture.md` — o diagrama/prosa devem descrever o sistema como ficou após o épico.
