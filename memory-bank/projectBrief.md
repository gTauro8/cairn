# Project Brief — Cairn

> Statico. Modifica solo in caso di pivot strategico. Lettura consigliata una sola volta a
> inizio lavoro, o quando lo scope sembra ambiguo.

## Perché esiste questo progetto

Git versiona il codice sorgente; Cairn versiona la *conoscenza* attorno al codice — decisioni
architetturali, alternative scartate, regole di business, incidenti, vincoli, debito tecnico —
come artefatto esplicito, versionato, interrogabile, con provenienza e relazioni. Oggi questa
conoscenza vive sparsa tra README, ADR sparsi, Confluence, Slack, Notion, issue, PR, memoria
delle persone e prompt LLM usa-e-getta: se le persone che l'hanno in testa se ne vanno o
dimenticano, la conoscenza si perde anche se il codice resta. Vedi anche RFC-0001
(`memory-bank/rfc/0001-livello-di-conoscenza.md`) per l'argomentazione estesa, con prove
concrete tratte dallo sviluppo di Cairn stesso.

## Obiettivi (in ordine di priorità)

1. Validare che la conoscenza si possa catturare con attrito quasi zero (vincolo centrale,
   AGENTS.md §2) — precede ogni altro obiettivo, si verifica osservando l'uso reale.
2. Restare utile anche se ogni LLM esistente sparisse domani: nessuna funzionalità core può
   dipendere strutturalmente da un modello specifico.
3. Trattare umani, agenti AI, CI/CD e strumenti di sviluppo come client di pari livello — non
   ottimizzare per un solo tipo di consumatore della conoscenza catturata.

## Confini invalicabili (esplicitamente fuori scope)

- Non è un framework AI, non un chatbot, non un sistema di memoria per agenti, non un vector DB.
- Gli agenti AI sono uno dei client possibili, non il client privilegiato.
- Nessuna funzionalità core deve dipendere strutturalmente da un modello LLM specifico.

## Vincoli non negoziabili

- Local-first e git-native: la conoscenza vive nello stesso repository del codice — si clona,
  si dirama e si fonde insieme ad esso (nessun servizio esterno obbligatorio per l'uso base).
- Cattura append-only, zero-friction: lo strumento primario di input è una CLI leggera, un
  singolo comando per registrare un'osservazione — mai un passaggio che richieda di lasciare
  il terminale o il flusso di lavoro in corso.
- Tecnologia noiosa e collaudata di default (stdlib prima di librerie esterne, git trailer
  prima di formati custom) — la novità deve ripagare chiaramente il suo costo (AGENTS.md §3).
- Nessuna migrazione forzata dello schema esistente senza un motivo forte: le estensioni sono
  additive (campi `omitempty`), non breaking.

## Chi decide in caso di ambiguità

- Il maintainer umano che guida le sessioni con gli agenti, in ultima istanza. Le decisioni
  tecniche quotidiane seguono la persona architetturale di AGENTS.md §3 senza bisogno di
  conferma esplicita per ogni scelta; i gate che richiedono conferma umana esplicita sono
  elencati in AGENTS.md §7 (merge in `main`, modifiche a `systemPatterns.md`, azioni del §6).
