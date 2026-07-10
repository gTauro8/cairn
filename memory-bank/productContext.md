# Product Context — Cairn

> Aggiornamento: basso. Rileggi quando la motivazione di una feature non è chiara.

## Problema reale che risolviamo

Il codice spiega il *come*; il *perché* — i compromessi accettati, le strade scartate, le
assunzioni di business, le scoperte fatte durante il debug — resta disperso in canali che
marciscono: wiki e ADR manuali richiedono uno sforzo cosciente estraneo al flusso di lavoro,
quindi vengono rimandati o mai scritti (vedi RFC-0001 §1). Quando le persone che hanno quel
contesto in testa se ne vanno o dimenticano, il software si irrigidisce e il debito tecnico
cresce per paura di toccare codice "di cui non si capisce il senso".

## Utenti e casi d'uso principali

1. Sviluppatore che scrive codice → annota una decisione mentre la prende (`cairn add`) →
   la nota resta accanto al codice, versionata, interrogabile in futuro.
2. Sviluppatore che committa → marca il commit con `Cairn-Note: true` quando vale la pena
   ricordarlo → la nota viene catturata automaticamente, zero passi extra (vedi `techContext.md`).
3. Chiunque (umano o agente AI) riprenda in mano il progetto dopo una pausa → consulta
   `cairn log`/`cairn log --tag x` → ricostruisce il perché di una scelta senza dover
   ricostruire la cronologia della chat o scandire l'intera codebase (AGENTS.md §4).
4. Agente AI che coordina o esegue un sottotask → legge/scrive note come client di pari
   livello rispetto agli umani, non privilegiato (AGENTS.md §1).

## Flussi chiave (UX/DX)

- **Cattura manuale:** `cairn add [--tags a,b,c] "testo"` → riga append-only in
  `.cairn/log.jsonl` nella cwd esatta del progetto.
- **Cattura automatica da commit:** commit con trailer `Cairn-Note: true` (+ `Cairn-Tags:`
  opzionale) → hook `post-commit` genera la nota da solo, nessun comando aggiuntivo.
- **Consultazione:** `cairn log` / `cairn log --tag x` → lettura testuale, filtrabile per tag.

## Metriche di successo

- Qualitativa (nessuna metrica quantitativa ancora raccolta): una decisione presa durante una
  sessione di lavoro viene catturata in quella stessa sessione, non rimandata a "dopo" — il
  criterio è osservare se l'attrito reale fa saltare l'abitudine (vincolo AGENTS.md §2).
- Segnale negativo da cercare attivamente: righe `.cairn/log.jsonl` che si accumulano senza
  mai essere consultate (`cairn log --tag` mai usato) indicherebbe cattura senza freschezza —
  non ancora osservato, da monitorare mano a mano che le note aumentano.
