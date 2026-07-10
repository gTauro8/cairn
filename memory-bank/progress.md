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
- [x] Tag opzionali sulle note: `cairn add --tags a,b,c "testo"` (campo `tags` additivo,
  omitempty), `cairn log --tag x` per filtrare — 2026-07-10, retrocompatibilità con le righe
  JSONL senza `tags` validata manualmente in scratch dir
- [x] Fix bug ordine flag: `cairn add "testo" --tags a,b` ora fallisce rumorosamente (nessuna
  riga scritta) invece di ingoiare `--tags` nel testo — 2026-07-10, via `misplacedFlag()` che
  deriva i nomi dei flag noti da `fs.VisitAll` (si auto-aggiorna se in futuro si aggiungono
  altri flag ad `add`). Validati: ordine corretto invariato, ordine sbagliato ora errore +
  exit 1, nessun falso positivo su testo libero con un trattino dentro

## In corso

- [ ] Nessun task di codice aperto al momento — in attesa di uso reale o nuova richiesta

## Da fare (prossimo)

- [ ] Decidere il prossimo incremento su `cairn` (uso reale di `add`/`log --tag` su note vere
  del progetto, o altra estensione)
- [ ] RFC-0001: "perché il software ha bisogno di un livello di conoscenza"

## Deferred (non riproporre senza un motivo nuovo)

- Relazioni tipizzate tra note (es. "supersede", "deprecated by"): rimandate il 2026-07-10.
  Motivo: richiedono riferirsi all'ID di una nota esistente, quindi un lookup che rallenta la
  cattura rispetto a v0/v0.1 — segnale contro secondo il vincolo del §2 di `AGENTS.md`. Da
  riprendere solo quando il dogfooding mostra un bisogno reale e concreto (es. una nota che
  *effettivamente* ne sostituisce un'altra e la mancanza del collegamento crea confusione),
  non per completezza dello schema.

## Bug noti

- Nessuno aperto. (Il bug sull'ordine di `--tags` — vedi `.cairn/log.jsonl` per la nota
  originale e quella di chiusura — è stato corretto il 2026-07-10: ora `cairn add "testo"
  --tags a,b` fallisce con errore invece di ingoiare il flag nel testo.)