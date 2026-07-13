# Prompt de design — Loop de extração (MVP)

> Brief completo para uma AI de design (ex.: Claude Design) produzir um **wireframe anotado low-fi**. O valor está neste brief; a AI de design só renderiza. Origem da tese: `_strategy/product.md`. Escopo e restrições vêm da estratégia — não os expanda.

## Contexto

RoosterLabs extrai o ponto de vista real de um profissional sênior e o transforma em conteúdo de LinkedIn que constrói autoridade. O produto **é a extração** (não a escrita). Este é o loop central: o autor traz (ou recebe) um tema, é entrevistado, e sai um post que só ele poderia ter escrito — que ele aprova e publica.

## Objetivo do artefato

Um **wireframe anotado, low-fi, mobile-first**, cobrindo todas as telas e estados do loop abaixo. Cada tela vem com uma **anotação do comportamento** por trás dela — o que o agente faz, não só o que aparece. Caixa e rótulo, sem polimento visual: a inteligência da interação é o que importa, a UI é a metade barata.

## Escopo (MVP — não ultrapassar)

Um ciclo: **um tema → um post**, com dois gates de validação. Fora de escopo: geração de múltiplos ângulos personalizados a partir de histórico, perfis de estilo elaborados, integrações de publicação, billing/onboarding de conta.

## Restrições (da estratégia — inegociáveis)

- **100% automatizado.** O único humano no loop é o autor. Nenhum humano (nem o Pedro) produz conteúdo.
- **Baixo custo de tempo** de um sênior ocupado. Fricção é inimiga; cada toque a mais precisa se pagar.
- **Áudio é a entrada primária do autor; o agente responde por texto.** Não é voz-para-voz.
- **Decent, not polished.** Funcional e claro, não bonito.
- **Barra de qualidade:** "só esse autor poderia ter postado isso" — medida por ele aprovar e publicar.

## O fluxo, etapa a etapa (com comportamento)

1. **Entrada de tema.** Pede um tema ao autor. Se ele traz um, segue. Se chega em branco, o sistema **sugere/extrai motes** — cientes do nicho, de preferência puxados da experiência recente ("o que você fez/viu essa semana que ensinou algo?"). O autor escolhe/confirma. *Comportamento:* mote é andaime, não produto; nunca genérico a ponto de gerar POV genérico.

2. **Entrevista (extração).** O agente pergunta por texto; o autor responde por áudio. O agente opera com **dois chapéus**: cabeça de jornalista (pensa no formato final e na estrutura) + profundidade de campo (conhecimento latente do modelo) para **enxergar a estrutura que emerge, identificar as lacunas e mirar as perguntas nelas**. Curta. *Comportamento:* para quando a estrutura não tem mais lacuna crítica — não estica a entrevista à toa (custo de tempo).

3. **Geração da arquitetura de informação.** Ao final, o agente gera a **estrutura do texto** (esqueleto/IA) em texto. *Comportamento:* estrutura antes de prosa — o esqueleto é validado antes de se gastar na redação.

4. **Gate 1 — validação da estrutura (áudio).** O autor aprova a estrutura ou pede mudanças por áudio. Mudança → volta e refina a estrutura. *Comportamento:* gate barato que evita um texto longo errado.

5. **Redação no tom.** Com a estrutura aprovada, o agente escreve o post, **bebendo do estilo extraído da fala do autor** — traduzindo o conteúdo falado em voz escrita de post, sem imitar o padrão oral (muleta, digressão, cadência).

6. **Gate 2 — validação do texto (áudio).** O autor aprova o texto ou pede mudanças por áudio. Mudança → volta e refina. Aprovado → pronto para publicar. *Comportamento:* a aprovação final antes de publicar **no nome do autor** é dele, sempre — é controle/reputação, e é o sinal que mede a qualidade.

**Transversal — acúmulo de estilo:** palavras/jargões/padrões capturados nesta conversa alimentam um perfil por autor que cresce com o tempo. No MVP, só capture a semente barata; não desenhe telas de gestão de perfil.

## Estados e telas a cobrir

- Entrada de tema: (a) autor traz o tema; (b) autor em branco → sistema sugere motes.
- Turno de entrevista: pergunta do agente (texto) + captura de áudio do autor (gravar/reenviar).
- Revisão da estrutura: estrutura exibida + aprovar / pedir mudança por áudio.
- Revisão do texto: post exibido + aprovar / pedir mudança por áudio.
- Concluído: post pronto para publicar (o publish em si está fora de escopo — só o estado).
- Edge: falha de áudio/transcrição; autor abandona e volta depois (o loop é assíncrono).

## Saída esperada da AI de design

Wireframe low-fi mobile-first, uma tela por estado acima, cada uma com a anotação de comportamento correspondente. Marque claramente os dois gates e os pontos de loop (pedir mudança → refinar).
