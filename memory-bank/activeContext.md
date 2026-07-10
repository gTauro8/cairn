# Active Context

> Altissima frequenza di aggiornamento. Importato automaticamente a ogni sessione.

## Cosa si sta facendo ora

Secondo incremento su `cairn` implementato: tag opzionali sulle note. `cairn add --tags
a,b,c "testo"` aggiunge un campo `tags` (omitempty, additivo, nessuna migrazione richiesta);
`cairn log --tag x` filtra per tag. `cairn add "testo"` senza flag resta identico a v0 —
retrocompatibilità verificata manualmente riga per riga. Vedi `progress.md` per lo stato
dettagliato.

## Ultima decisione presa

Relazioni tipizzate tra note (supersede/deprecated-by) rimandate esplicitamente: richiedono
riferirsi all'ID di una nota esistente, quindi un passaggio di lookup che rallenta la cattura
— segnale contro secondo il vincolo del §2 di `AGENTS.md`. Si riprende in mano solo se il
dogfooding mostra un bisogno reale, non per completezza dello schema.

## Prossimo passo

Nessun incremento di codice deciso — in attesa di uso reale di `cairn add`/`cairn log --tag`
sul progetto stesso, o di una nuova richiesta in modalità Plan (`AGENTS.md` §5).

## Blocchi/domande aperte

- Nessuno al momento. (Nota: l'ADR del 2026-07-10 in `systemPatterns.md` ha ora la
  motivazione compilata — la voce precedente che la segnava "in compilazione" era stale ed è
  stata rimossa qui.)