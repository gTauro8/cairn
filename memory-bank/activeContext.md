# Active Context

> Altissima frequenza di aggiornamento. Importato automaticamente a ogni sessione.

## Cosa si sta facendo ora

Pausa consapevole sugli incrementi guidati da spec: dogfooding reale di `cairn add`/`log` sul
repo stesso, iniziato il 2026-07-10 (una nota già inserita a mano dall'utente, quattro
aggiunte da me durante la sessione). Vedi `.cairn/log.jsonl` per le note reali e
`progress.md` per il bug trovato nell'uso vero.

## Ultima decisione presa

Relazioni tipizzate tra note restano rimandate (vedi `progress.md` § Deferred). Nessuna
nuova decisione architetturale oggi: la sessione è passata da "costruire il terzo
incremento" a "usare i primi due sul serio e vedere cosa si rompe".

## Prossimo passo

Continuare a usare `cairn add`/`cairn log` per note reali nei prossimi giorni, senza pilotare
l'uso con una spec. Il bug di ordine dei flag (`--tags` dopo il testo, vedi `progress.md` §
Bug noti) resta aperto e non corretto di proposito — è un segnale da dogfooding, non un
incidente da patchare subito. Decidere se/come correggerlo solo quando emergono altri segnali
simili o quando l'utente lo richiede esplicitamente in Plan mode.

## Blocchi/domande aperte

- Nessuno al momento. (`.cairn/` non è più in `.gitignore` — decisione confermata
  dall'utente il 2026-07-10: le note vanno versionate da subito, coerente con la visione di
  Cairn (§1). `.cairn/` risulta ora untracked in `git status`, non ancora aggiunto/committato
  — lascio quel passo esplicito all'utente.)