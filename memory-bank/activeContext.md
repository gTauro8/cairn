# Active Context

> Altissima frequenza di aggiornamento. Importato automaticamente a ogni sessione.

## Cosa si sta facendo ora

Primo pezzo di codice Go implementato: comando `cairn` con `add`/`log` che scrive/legge un log
append-only di note in `.cairn/log.jsonl` (cwd esatta, niente walk-up — deciso in Plan). Vedi
`progress.md` per lo stato dettagliato.

## Ultima decisione presa

`.cairn/` viene cercato solo nella cwd esatta, senza risalita delle directory padre (analogo
a git ma senza walk-up). Motivazione: meno codice in v0, coerente col "niente config" già
messo fuori scope nella spec del primo pezzo.

## Prossimo passo

Decidere il secondo incremento (es. primo uso reale di `cairn add` per annotare una decisione
vera del progetto, oppure estendere `log`/`add` prima di quello). Restare in modalità Plan
prima di scrivere altro codice, come da `AGENTS.md` §5.

## Blocchi/domande aperte

- La motivazione dell'ADR del 2026-07-10 in `systemPatterns.md` è ancora in compilazione da
  parte dell'utente insieme a Giuseppe — non toccare quel file di iniziativa propria finché
  non è stata inserita.