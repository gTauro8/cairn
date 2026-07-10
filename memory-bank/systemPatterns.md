# System Patterns — Cairn

> Aggiornamento: medio, quando si introduce un nuovo modulo o pattern. Lettura on-demand prima
> di aggiungere codice nuovo, non a ogni sessione.

## Architettura ad alto livello

[Descrizione testuale o diagramma dei componenti principali e di come comunicano tra loro]

## Pattern obbligatori

- [Es. "tutte le scritture su disco passano dal livello X", "niente logica di business nei
  controller/handler", ecc.]
- [...]

## Decisioni architetturali (ADR sintetiche)

| Data | Decisione | Alternative scartate | Motivazione |
|---|---|---|---|
| 2026-07-10 | Saltare la validazione manuale pre-codice e passare direttamente all'implementazione in Go | Validazione manuale di 1 settimana (pianificata e poi non eseguita) | Semplice impazienza con il processo: dopo settimane su naming, prompt architetturale e governance, la spinta a vedere codice reale ha prevalso sul piano originale. Il vincolo del §2 di `AGENTS.md` resta comunque valido — lo si verifica ora osservando il prodotto reale invece che con un test a parte. |



## Anti-pattern noti in questo progetto (cosa NON fare)

- [Es. "non usare la libreria X per Y: causa il problema Z, vedi ADR del ..."]
