# Guardrails — valores e regras que validam trabalho

Registro fechado. Um item de dívida (ou justificativa de épico backpack-relief) **só é válido se apontar para um guardrail desta lista**. Não aponta? Duas saídas: propor refinamento desta lista (decisão explícita, registrada em `decisions.md`) ou o item não é dívida — é preferência, e não entra.

## Valores (herdados da estratégia, operacionalizados aqui)

| Valor | O que protege | Origem |
|---|---|---|
| **Mochila leve** | velocidade contínua: dívida acumulada é peso que taxa todo passo futuro | esta decisão (2026-07-07) |
| **Decent, not polished** | foco: qualidade sobe por iteração com clientes reais, não por antecipação | `_strategy/decisions.md` |
| **Operação solo sem atenção recorrente** | o tempo do Pedro: ops manuais repetidas são bug, não rotina | `_strategy/decisions.md` |
| **Delivery 100% automatizada** | a tese da empresa: nenhum humano no loop de conteúdo de cliente | `_strategy/decisions.md` (regra dura) |
| **Menos caixa-preta** | revisabilidade: o gargalo do loop homem+AI é a revisão humana | `decisions.md` (stack) |

*Tensão assumida:* "mochila leve" e "decent, not polished" se limitam mutuamente — dívida válida é desvio de um **padrão comprometido** (abaixo), e os padrões já embutem o nível "decent". Polimento além do padrão não é dívida.

## Canais de erosão (dívida sem métrica ainda é dívida — se nomear o canal)

| Família de dívida | Erode | Chega no negócio como |
|---|---|---|
| Técnica | velocidade e confiabilidade de tudo que vem depois | menos iterações, incidentes, deploy travado |
| Design | credibilidade | menos vendas |
| UX | satisfação do cliente | menos retenção |
| Copy | clareza da promessa | menos vendas, mais suporte, menos retenção |

## Padrões comprometidos (dívida = desvio de um destes)

| Domínio | Padrão | Dono |
|---|---|---|
| Código/arquitetura | `docs/conventions.md`, decisões de stack em `decisions.md` | engineering |
| Infra/ops | `infra/README.md` (runbook), decisões de infra em `decisions.md` | engineering |
| Identidade visual | `roosterlabs-marketing/visual-identity.md` | marketing |
| Copy/tom | copy aprovada (`landing-page.md`) + `comms-foundations.md` | marketing |
| Comportamento/UX | specs de comportamento nos docs de marketing (ex.: spec do form) | marketing define, engineering implementa |
| Negócio | `_strategy/goals.md`, `business-model.md`, `decisions.md` | strategy |

**Não existe padrão para o desvio apontado?** Então não é dívida — é decisão faltando. O trabalho vai para o dono do domínio criar o padrão primeiro.

## Como usar (regra operacional)

1. Item de dívida entra num backlog (`docs/backlog.md` aqui; equivalente em marketing para design/UX/copy) nomeando: **padrão violado** + **canal de erosão**.
2. Épico `backpack-relief` (ver `workflow.md`) puxa itens registrados; escopo congela no aceite; DoD = desvio eliminado, verificado.
3. Propriedade não muda: engineering não inventa copy/design; o dono define, o épico executa.
4. Esta lista muda por decisão explícita, nunca por precedente silencioso.
