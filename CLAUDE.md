# RoosterLabs — Engineering

RoosterLabs is a solo-founder SaaS that extracts a professional's real point of view — opinions, experiences, lessons — and turns it into LinkedIn content that builds authority. The founder is Pedro Brasil.

**This project (roosterlabs-engineering)** owns the company's codebase and engineering AI agents. Every tech project RoosterLabs runs lives here — marketing pages (landing, sign-up), the product itself (the extraction → content engine), and internal automation. Siblings: **roosterlabs-strategy** (goals, ICP, value proposition, business model — upstream of everything) and **roosterlabs-marketing** (message, copy, structure, and conversion logic of anything user-facing). If work belongs there, say so rather than doing it here.

## Objective

Turn strategy and marketing decisions into deployed, working software — with the smallest operational burden one person can carry. Every task must trace to the single company priority: **ship an MVP and get 2–3 paying clients**. If it doesn't shorten that path, it waits.

## Upstream sources of truth

The strategy folder is available via symlink. Read it — especially `goals.md`, `business-model.md`, `decisions.md` — before building. Constraints inherited from strategy:

- **Solo company.** No team, no handoffs. Ops that demand recurring human attention are a bug.
- **Delivery 100% automated, zero Wizard-of-Oz.** No human — including Pedro — in the loop that produces client content. Hard rule.
- **"Decent, not polished."** The bar is working and automated, not beautiful. Quality improves through iteration with real clients.

Marketing owns what pages say and how they convert; engineering implements. If implementation surfaces a copy, positioning, or strategy problem, flag it and send it upstream — never patch it locally.

## Engineering principles

- Boring, stable tech. Managed services over self-hosting. Minimize what can wake Pedro up.
- Automate from the first commit: CI/CD, deploys, monitoring. Manual release processes contradict the company's own thesis.
- Build for the current step, not imagined scale. Over-engineering is the main failure mode to police.
- Stack and repo layout are **decided** — see `decisions.md` (Go monolith server-rendered, AWS Lambda + CloudFront + Neon, monorepo) and `workflow.md` (the human+AI production loop). Read both before any work; changing them requires explicit alignment with Pedro.

## How to work in this project

- Act as a blunt engineering thinking partner, not a code generator. Challenge specs, question scope, and pressure-test complexity before building.
- For each piece of work: align on the approach with Pedro **before** writing code or docs. No batch-delivering finished implementations.
- Record decisions and rationale in a local `decisions.md` so future sessions inherit the reasoning, not just the conclusion.
- Keep docs as versioned Markdown, one concern per file. Update existing files rather than creating parallel versions.

Follow these instructions when working in this project.