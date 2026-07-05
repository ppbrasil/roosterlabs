# Épico NNN — <nome curto>

**Estado:** proposto <!-- proposto | escopado | aceito | em-execução | concluído -->
**Origem:** <strategy | marketing> · <data>
**Prioridade servida:** MVP + 2–3 clientes pagantes (se não servir, o épico não existe)

## Output / Outcome (por que este épico existe)

<Output: entrega concreta e verificável (válido em fase inicial). Outcome: resultado de negócio mensurável. Pelo menos um dos dois. Sem solução técnica aqui — descrever solução é defeito.>

- **Critério de sucesso:** <output verificável e/ou métrica de outcome>
- **Restrições:** <prazo, orçamento, dependências>
- **Fora de escopo:** <o que explicitamente NÃO entra>

## Escopo proposto (engineering — skill `escopar-epico`)

<Opções consideradas com trade-offs; recomendação com rationale. Desafios só ao que for long shot.>

### Mudança no sistema

<Antes → depois: componentes, rotas, dados afetados. Referência à atualização proposta de `docs/architecture.md` (diagramas). Este é o entregável central do escopo — transparência, não corte.>

## DoD — Definition of Done (congela no aceite)

<Cada item verificável com evidência objetiva. "Funciona" não é item; "POST /subscribe grava lead no Postgres de produção e retorna 200" é.>

- [ ] <item verificável>
- [ ] <item verificável>

## Log de validação (skill `validar-escopo`)

- <data> — <aprovado | revisão solicitada> — <notas>

## Quebra (skill `quebrar-epico`)

<Lista de micro-tarefas ou links para issues. Cada tarefa: descrição, plano de teste com edge cases, item do DoD que atende.>

## Fechamento (skill `fechar-epico`)

- **Verificação do DoD:** <item → evidência (URL testada, medição, query)>
- **Documentação atualizada:** <quais arquivos e por quê>
- **Aprendizados → engenharia:** <issues criadas>
- **Aprendizados → strategy:** <notas>
- **Impacto em GTM → marketing:** <nota para ajuste, ou "nenhum">
- **Handoff:** <resumo do entregue + candidatos a próximo ponto de melhoria>
