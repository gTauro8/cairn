# Active Context

> Altissima frequenza di aggiornamento. Importato automaticamente a ogni sessione.

## Cosa si sta facendo ora

Corretto in Plan→Act→Validate→Sync il bug di dogfooding sull'ordine dei flag: `cairn add
"testo" --tags a,b` ora fallisce rumorosamente (nessuna riga scritta) invece di ingoiare
`--tags` nel testo. Vedi `progress.md` per i dettagli e i criteri validati.

## Ultima decisione presa

Fix minimo e mirato (rilevare flag noti finiti tra i positional via `fs.VisitAll` ed errore
esplicito) invece di un parser permutante fatto in casa o di rendere `--text` obbligatorio —
entrambe le alternative scartate per complessità sproporzionata o regressione sul caso comune
(nota senza tag). Relazioni tipizzate tra note restano comunque rimandate (vedi `progress.md`
§ Deferred), non toccate in questo incremento.

## Prossimo passo

Nessun task di codice aperto. Riprendere il dogfooding reale di `cairn add`/`cairn log`, o
aprire un nuovo giro di Plan su uno degli altri due candidati già in backlog (prossimo
incremento funzionale su `cairn`, oppure RFC-0001) quando l'utente decide quale.

## Blocchi/domande aperte

- Nessuno al momento. (`.cairn/` non è più in `.gitignore` — decisione confermata
  dall'utente il 2026-07-10: le note vanno versionate da subito, coerente con la visione di
  Cairn (§1). `.cairn/` risulta ora untracked in `git status`, non ancora aggiunto/committato
  — lascio quel passo esplicito all'utente.)