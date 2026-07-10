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

## In corso

- [ ] Nessun task di codice aperto al momento

## Da fare (prossimo)

- [ ] Estensione IDE (VS Code) per rendere le note visibili accanto al codice — ora che
  esiste `files`, si può scoping seriamente; Plan dedicato quando richiesto
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