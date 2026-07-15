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
- [x] Cattura automatica da commit: hook `post-commit` (`.githooks/post-commit`) genera una
  nota `cairn` quando il commit ha il trailer `Cairn-Note: true` (testo = subject + hash
  corto, tag = `Cairn-Tags:` + tag automatico `git`); commit senza trailer restano silenziosi
  (niente mirror di ogni commit). Setup: `git config core.hooksPath .githooks`. Dettagli e
  convenzione in `techContext.md`. Validati in scratch clone: commit con trailer → nota
  corretta; commit senza trailer → nessuna nota; binario `cairn` assente → warning, commit
  comunque riuscito — 2026-07-10
- [x] Fix hook post-commit: il primo commit reale con questa feature (f7b5f89) NON ha
  generato la nota attesa — `git interpret-trailers --parse` riconosce solo l'ultimo
  paragrafo contiguo come blocco trailer, e la riga vuota prima di `Co-Authored-By:` faceva
  ignorare silenziosamente `Cairn-Note`/`Cairn-Tags`. Scoperto committando la feature stessa,
  corretto cercando le righe `Cairn-*` ovunque nel messaggio invece di passare da
  `git interpret-trailers` — 2026-07-10, ridogfoodato con successo sul commit di fix (39771a9)

- [x] Primo dispatch multi-agente (§7): test automatici da Codex (`cmd/cairn/main_test.go`,
  7 test case, tutti verdi) e RFC-0001 da Gemini (`memory-bank/rfc/0001-livello-di-conoscenza.md`)
  rivisti, committati (nessuno dei due aveva committato da solo — fatto dal coordinatore) e
  mergiati in `main` con `--no-ff`, nessun conflitto — 2026-07-10. Note operative sul processo
  in `activeContext.md` § Handoff chiuso
- [x] Memory bank compilata per intero: `projectBrief.md`, `productContext.md`,
  `systemPatterns.md` (erano placeholder), `techContext.md` aggiornato — 2026-07-10
- [x] Push su `main` di tutto il lavoro della sessione (`040a87a..f0f1ac7`, 12 commit) — 2026-07-10
- [x] Campo `files` sulle note (prerequisito per una futura estensione IDE): `cairn add
  --files a,b "testo"` (additivo, `omitempty`, nessuna lowercase perché i path sono
  case-sensitive), `cairn log --file p` per filtrare. L'hook `post-commit` lo popola in
  automatico dai file effettivamente cambiati nel commit (`git diff-tree --name-only`) — a
  differenza di `--tags`/`--files` su `add` manuale, qui l'inferenza è sicura per costruzione.
  Deciso di NON referenziare righe/range specifici: si sposterebbero a ogni refactor
  diventando silenziosamente sbagliati (rischio di freschezza, §2) — solo file interi.
  `misplacedFlag()` protegge anche `--files` senza modifiche, essendo già generica su
  `fs.VisitAll`. Validati: 2 nuovi test automatici, più hook testato in scratch clone
  (auto-popolamento da commit reale, filtro `--file`, retrocompat su `add` senza `--files`,
  rifiuto di `--files` dopo il testo) — 2026-07-10
- [x] Primo incremento v0.1: `cairn log --json` emette le righe JSONL originali complete e
  resta combinabile con `--tag`/`--file`; `cairn check` accumula e segnala JSON non valido,
  campi obbligatori (`id`, `ts`, `text`) mancanti/vuoti e file referenziati inesistenti, con
  exit non-zero senza modificare il log. Tre nuovi test automatici; build, vet, test e check
  sul log reale verdi — 2026-07-15, modifiche ancora da committare
- [x] Provenienza strutturata additiva: campi opzionali `source`/`commit`; le nuove note
  manuali ricevono `source: manual`, mentre l'hook passa `source: git` e l'hash completo del
  commit tramite i nuovi flag di `cairn add`. `cairn check` segnala `source: git` senza
  `commit`, lasciando valide le note legacy. Due nuovi test automatici; build/vet/test verdi;
  hook validato con commit reale in scratch e hash corrispondente — 2026-07-15, modifiche
  ancora da committare
- [x] Onboarding v0.1: `cairn hook install` idempotente dalla root Git, con rifiuto sicuro di
  `core.hooksPath` o `post-commit` preesistenti e diversi. Parsing/cattura spostati in
  `cairn hook run`; lo script versionato è un wrapper minimale che non blocca mai il commit.
  Test automatici per installazione, idempotenza, conflitti, parsing e cattura — 2026-07-15
- [x] README riscritto come documentazione del prodotto; aggiunti Apache-2.0, `CHANGELOG.md`,
  `RELEASING.md`, `cairn version`, Makefile con build riproducibili/cross-build e workflow CI
  macOS/Linux (`actions/checkout@v6`, `actions/setup-go@v6`) — 2026-07-15
- [x] Protocollo di dogfooding v0.1 definito: almeno 10 sessioni/14 giorni, criteri misurabili
  per cattura, consultazione, freschezza e attrito; estensione IDE dietro gate empirico —
  `memory-bank/dogfooding-v0.1.md`, 2026-07-15
- [x] Validazione locale completa del blocco v0.1: `make verify`, quattro cross-build
  macOS/Linux amd64/arm64, versione `v0.1.0` incorporata e scratch end-to-end (`hook install`,
  commit non marcato silenzioso, commit marcato con provenienza/file/tag corretti, `check`
  verde) — 2026-07-15
- [x] Gate umano ricevuto e `systemPatterns.md` sincronizzato col refactoring dell'hook —
  2026-07-15

## In corso

- [ ] Nessun task di codice aperto al momento

## Da fare (prossimo)

- [ ] Commit su `main` autorizzato ma bloccato dal sandbox (`.git/index.lock` non scrivibile):
  eseguirlo da un terminale con permesso Git, escludendo `.DS_Store`
- [ ] Eseguire il dogfooding v0.1 (10 sessioni su 14 giorni) e registrarne l'esito
- [ ] Estensione IDE (VS Code) solo se il dogfooding mostra utilità delle note ma scarsa
  discoverability accanto al codice
- [ ] Eventuali limature di stile su RFC-0001 (attribuzione, tono) se l'utente le vuole

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
