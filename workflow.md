# Workflow — o loop de produção homem+AI

Princípio de arquitetura do processo: **o barramento entre projetos são arquivos versionados com estados explícitos**, não conversas. Strategy, marketing e engineering são sessões separadas sem memória compartilhada; o que não está em arquivo evapora. Pedro é o roteador entre projetos e o gate humano em três pontos: seleção do épico, aprovação do escopo/DoD, merge do PR.

## O ciclo de épicos (macro)

Todo trabalho nasce de um épico. Três tipos válidos (decisão 2026-07-07; ver `guardrails.md`):

1. **Outcome** — resultado de negócio mensurável.
2. **Output** — entrega concreta e verificável (válido em fase inicial, onde ainda não há tráfego para medir outcome).
3. **Backpack-relief** — alívio de dívida registrada (técnica, design, UX, copy). Cada item puxado de um backlog nomeia o **padrão comprometido que viola** e o **canal de erosão que alimenta** (`guardrails.md`); item que não aponta guardrail não entra. Escopo congela no aceite; DoD = desvio eliminado, verificado. Propriedade não muda: o dono do domínio define o certo, o épico executa.

"Solução técnica como tema" continua sendo defeito — o que legitima trabalho técnico é dívida registrada contra guardrail, nunca o tema. O épico é um arquivo em `epics/` (template em `epics/TEMPLATE.md`) com máquina de estados:

```
proposto → escopado → aceito → em-execução → concluído
              ↑___________|        (loop de revisão)
```

| # | Etapa | Quem executa | Peça | Estado resultante |
|---|---|---|---|---|
| 1 | Definir | Pedro + strategy/marketing | skill `definir-epico` | `proposto` |
| 2 | Escopar | engineering | skill `escopar-epico` | `escopado` |
| 3 | Validar | Pedro + strategy/marketing | skill `validar-escopo` | `aceito` (DoD congela) ou volta a `proposto` com notas |
| 4 | Quebrar | engineering | skill `quebrar-epico` | `em-execução` + micro-tarefas |
| 5 | Build/Test/Review/Deploy | engineering (por tarefa) | Claude Code + agente `revisor` | tarefas concluídas |
| 6 | Fechar | engineering | skill `fechar-epico` | `concluído` + handoff do próximo ciclo |

### 1. Definir (strategy/marketing)
- Épico descreve **outcome, output ou backpack-relief** (tipos acima). Em fase inicial, um output concreto e verificável — "landing no ar em roosterlabs.com.br capturando leads" — é épico válido; não forçar métrica onde ainda não há tráfego para medi-la. Backpack-relief nasce dos backlogs de dívida (engineering: `docs/backlog.md`; marketing: equivalente lá), nunca de tema genérico.
- Copy, estrutura e conversão vêm de `roosterlabs-marketing`; requisitos de negócio de `_strategy`. Engenharia não inventa nenhum dos dois.

### 2–3. Escopar ↔ Validar (o loop de transparência)
- Coding não é o gargalo — **o foco do escopo não é enxugar entrega, é dar transparência sobre qual será a mudança no sistema**. Entregável central: a seção "Mudança no sistema" do épico + atualização proposta de `docs/architecture.md` (diagramas/docs refletindo o antes → depois: componentes, rotas, dados).
- Engineering desafia apenas o que parecer **long shot** (aposta desproporcional ao retorno, ou sem traço à prioridade única) — não corta escopo por reflexo.
- Rascunho de DoD acompanha. Strategy/marketing valida: a mudança está compreensível? O DoD entrega o output/outcome? Aprova ou pede revisão (notas no arquivo).
- **No aceite o DoD congela.** Mudou o DoD = novo loop de validação, explícito.

### 4. Quebrar (engineering)
- **Unidades tão mínimas quanto possível** para conter alucinação — mas a unidade mínima é **um comportamento observável**, não um arquivo (um comportamento em Go legitimamente toca handler + registro + template + teste; forçar um arquivo por tarefa quebra "main deployável").
- Anti-alucinação, cada tarefa declara: **blast radius** (arquivos/funções que espera tocar e que deve ler antes — sair da lista durante a implementação = parar e atualizar a tarefa), plano de teste com edge cases enumerados antes de código, traço ao item do DoD, e orçamento de diff (~150 linhas, excluindo gerado). **Proibido drive-by**: melhoria fora do comportamento vira tarefa própria.

### 5. Build / Test / Review / Deploy (por tarefa)
- Branch por tarefa. **Testes primeiro** (do plano de teste), depois implementação. Commits pequenos.
- Local/devcontainer: `make test` antes de qualquer push. CI roda build + vet + test + lint em todo PR; PR vermelho não é revisado.
- Review em camadas: (1) agente `revisor` — contexto limpo, ataca o diff contra tarefa, DoD e edge cases; (2) CI verde; (3) Pedro revisa o diff e roda local. `/security-review` obrigatório quando tocar input de usuário, auth ou dados de cliente.
- Merge em `main` → deploy automático em produção. Rollback = revert + merge.
- Problema de copy/posicionamento descoberto no meio? Não corrige local — nota no épico, sobe upstream.

### 6. Fechar (engineering) — e repetir
Invocada quando a última tarefa foi deployada e a hipótese é "DoD cumprido". Fecha **por evidência, não por sensação**:
- Verifica o DoD item a item com prova real (testa a URL de produção, mede, consulta). Item sem evidência não fecha; item falho gera tarefas-lacuna e o épico permanece `em-execução`.
- **Atualiza toda a documentação do produto** afetada (README, infra/, docs de domínio) para refletir o estado entregue.
- Roteia aprendizados: engenharia → issues; negócio → notas endereçadas a strategy; **impacto em GTM → nota explícita endereçada a marketing** para ajuste quando necessário.
- Arquiva em `epics/done/` e escreve o handoff — entregas, evidências, aprendizados, candidatos a próximo ponto de melhoria — que Pedro leva à sessão de strategy/marketing para o `definir-epico` recomeçar o ciclo.

## Papéis

| Quem | Faz | Não faz |
|---|---|---|
| Pedro | seleciona épico, aprova escopo/DoD, revisa PR, roteia entre projetos, opera contas | escrever boilerplate; produzir conteúdo de cliente (regra dura da estratégia) |
| AI (Claude Code / Cowork) | desafia specs, escopa, quebra, implementa, testa, revisa (agente), mantém docs | mergear sem review; congelar DoD sozinha; inventar copy/posicionamento |
| Pipeline (Actions) | test, lint, build, deploy, (futuro) monitoramento | — |

## Onde vivem as peças

- Skills e agente: `.claude/skills/` e `.claude/agents/` neste repo (versionados). No Claude Code funcionam nativamente; para usar nas sessões Cowork de strategy/marketing, Pedro instala as skills correspondentes em Settings > Capabilities a partir destes arquivos-fonte.
- Épicos: `epics/` neste repo. Strategy/marketing precisam da pasta acessível nas sessões deles (mesmo esquema do symlink `_strategy`, direção inversa).

## Quando criar uma peça nova

As peças acima existem porque o protocolo entre projetos precisa existir *antes* do uso — sem ele o processo evapora entre sessões. Qualquer peça além destas volta à regra: só quando um passo repetir com atrito ≥3 vezes, automatizando aquele passo específico, com registro em `decisions.md`. Nunca criar agente por completude.
