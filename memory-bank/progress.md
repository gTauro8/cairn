# Progress

> Aggiornamento alto: a ogni chiusura di ticket/milestone. Importato automaticamente.

## Fatto

- [x] Naming del progetto: Cairn
- [x] Scaffolding di governance creato (AGENTS.md, CLAUDE.md, .claude/, .gemini/, memory-bank/)
- [x] Verificato che Claude Code, Codex e Gemini leggono AGENTS.md correttamente
- [x] Definiti durata e criterio di successo della validazione manuale (poi non eseguita)
- [x] Decisione di saltare la validazione manuale e passare all'implementazione diretta
- [x] Definizione del primo task di implementazione (Plan mode, vedi AGENTS.md §5) — 2026-07-10
- [x] Bootstrap del modulo Go: `go.mod`, `cmd/cairn/main.go`, comandi `cairn add`/`cairn log`
  su log append-only `.cairn/log.jsonl` (cwd esatta, no walk-up) — 2026-07-10, build e test
  funzionale validati manualmente in scratch dir

## In corso

- [ ] Nessun task di codice aperto al momento — in attesa di definire il secondo incremento

## Da fare (prossimo)

- [ ] Compilare la motivazione della decisione di saltare la validazione (systemPatterns.md —
  in corso lato utente/Giuseppe, non toccare)
- [ ] Decidere il secondo incremento su `cairn` (uso reale di `add` su una nota vera, o
  estensione di `add`/`log`)
- [ ] RFC-0001: "perché il software ha bisogno di un livello di conoscenza"

## Bug noti

- Nessuno: implementazione appena iniziata.