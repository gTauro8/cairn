# System Patterns — Cairn

> Aggiornamento: medio, quando si introduce un nuovo modulo o pattern. Lettura on-demand prima
> di aggiungere codice nuovo, non a ogni sessione.

## Architettura ad alto livello

Un solo binario Go (`cmd/cairn/main.go`, package `main`, nessun layering interno finché non
serve): due comandi, `add` e `log`, che scrivono/leggono un log append-only in formato JSON
Lines (`.cairn/log.jsonl`), cercato solo nella cwd esatta del progetto (nessuna risalita alle
directory padre). Nessun server, nessun database: SQLite e Tree-sitter restano nello stack
ipotizzato in `techContext.md` ma non sono ancora stati introdotti — la CLI a file piatto
basta finché non emerge un bisogno reale (query complesse, concorrenza) che li giustifichi.

Cattura automatica opzionale via hook Git (`.githooks/post-commit`, wrapper shell POSIX): lo
script individua il binario e delega a `cairn hook run`; parsing dei trailer, raccolta dei file
e provenienza vivono nel binario Go, che riusa direttamente la stessa logica di `add`. Il
wrapper converte ogni errore in warning e non può far fallire il commit. `cairn hook install`
installa/configura il wrapper senza sovrascrivere hook o `core.hooksPath` differenti.

## Pattern obbligatori

- **Schema JSONL solo additivo:** nuovi campi sempre `omitempty`/opzionali (vedi `tags`);
  nessuna riga esistente deve richiedere una migrazione quando si aggiunge un campo. Se un
  cambiamento richiedesse di riscrivere le righe esistenti, è quasi certamente il segnale che
  lo schema è sbagliato, non solo il codice che lo scrive.
- **Cattura mai implicita:** sia manuale (`cairn add`) sia automatica (hook), una nota nasce
  solo da un'azione intenzionale — mai un mirror automatico di ogni evento (es. ogni commit).
  Mirrorare tutto sarebbe rumore, non conoscenza (vedi Anti-pattern sotto).
- **stdlib prima di dipendenze esterne:** parsing flag/CLI, trailer git, generazione ID — tutto
  con la libreria standard di Go o strumenti git nativi, finché non dimostrano di non bastare.
- **`.cairn/` cercato solo nella cwd esatta:** nessun comportamento implicito di risalita alle
  directory padre (a differenza di `.git/`) — deciso per tenere il comportamento prevedibile
  senza introdurre configurazione.
- **Ogni incremento validato prima del merge**, manualmente o con test automatici, in una
  directory/worktree isolata — mai contro lo stato reale del repository principale.

## Decisioni architetturali (ADR sintetiche)

| Data | Decisione | Alternative scartate | Motivazione |
|---|---|---|---|
| 2026-07-10 | Saltare la validazione manuale pre-codice e passare direttamente all'implementazione in Go | Validazione manuale di 1 settimana (pianificata e poi non eseguita) | Semplice impazienza con il processo: dopo settimane su naming, prompt architetturale e governance, la spinta a vedere codice reale ha prevalso sul piano originale. Il vincolo del §2 di `AGENTS.md` resta comunque valido — lo si verifica ora osservando il prodotto reale invece che con un test a parte. |
| 2026-07-10 | `.cairn/` cercato solo nella cwd esatta, niente risalita alle directory padre | Walk-up come `.git/` | Meno codice in v0; coerente col principio "niente config implicita" già scelto per il resto della CLI. |
| 2026-07-10 | Tag opzionali sulle note come campo additivo (`omitempty`), non relazioni tipizzate tra note | Relazioni tipizzate (`supersede`/`deprecated by`) | Le relazioni richiedono conoscere l'ID di una nota esistente (lookup prima della scrittura): attrito reale contro il vincolo del §2. Rimandate finché il dogfooding non mostra un bisogno concreto (vedi `progress.md` § Deferred). |
| 2026-07-10 | Fix del bug sull'ordine dei flag: `cairn add` rileva un flag noto finito tra gli argomenti posizionali ed errore esplicito, invece di un parser permutante fatto in casa o di un `--text` obbligatorio | Parser permutante custom; `--text` obbligatorio | Il parser custom è complessità sproporzionata per un solo flag; `--text` obbligatorio regredirebbe il caso comune (nota senza tag). La soluzione scelta trasforma una perdita di dati silenziosa in un errore visibile, senza toccare l'ordine richiesto dei flag. |
| 2026-07-10 | Cattura automatica solo da commit esplicitamente marcati (trailer `Cairn-Note: true`), non mirror di ogni commit | Mirror di ogni commit; cattura da PR (GitHub Action) | Git già versiona ogni commit — mirrorarli tutti sarebbe rumore, non conoscenza (visione in AGENTS.md §1). La cattura da PR è infrastruttura diversa (CI, non un hook locale), rimandata finché il trailer sui commit non dimostra di funzionare nell'uso reale. |
| 2026-07-10 | Rilevamento dei trailer `Cairn-*` via grep diretto sul messaggio, non via `git interpret-trailers` | `git interpret-trailers --parse` | Quel comando riconosce solo l'ultimo paragrafo contiguo come blocco trailer: una riga vuota prima di `Co-Authored-By:` faceva ignorare silenziosamente `Cairn-Note`/`Cairn-Tags`, con perdita di dati scoperta solo committando la feature stessa (vedi Anti-pattern sotto). |
| 2026-07-10 | Campo `files` sulle note riferisce solo file interi, mai righe/range specifici | Riferimenti a riga/range per un'eventuale estensione IDE | Un numero di riga si sposta a ogni refactor e diventerebbe silenziosamente sbagliato — rischio di freschezza (§2) peggiore del non avere il riferimento. Popolato manualmente su `cairn add --files`, automaticamente dall'hook `post-commit` (i file del commit sono per definizione il contesto, inferenza sicura). |
| 2026-07-15 | Logica dell'hook nel binario Go; `post-commit` ridotto a wrapper installabile e non bloccante | Generare e mantenere l'intero script shell; configurazione manuale per clone | Una sola implementazione del parsing evita drift e dipendenze da `grep`/`sed`; `hook install` riduce l'attrito ma rifiuta conflitti invece di sovrascrivere configurazioni dell'utente. Il compromesso `post-commit` lascia il log modificato fino al commit successivo e va misurato nel dogfooding. |

## Anti-pattern noti in questo progetto (cosa NON fare)

- Non mirrorare ogni commit (o ogni evento in genere) in una nota Cairn: il valore di Cairn è
  la conoscenza *attorno* al codice, non un doppione di ciò che git già versiona. Cattura solo
  ciò che è marcato esplicitamente come degno di nota.
- Non usare `git interpret-trailers --parse` per rilevare i trailer `Cairn-*` se nel messaggio
  possono comparire altri trailer dopo una riga vuota (es. `Co-Authored-By:`): riconosce solo
  l'ultimo paragrafo contiguo e ignora silenziosamente tutto il resto. Il parser Go cerca le
  righe `Cairn-*` ovunque nel messaggio. Bug reale, corretto originariamente nel commit
  `39771a9` — vedi anche `techContext.md`.
- Non introdurre relazioni tipizzate tra note "per completezza dello schema": vedi ADR sopra,
  vanno introdotte solo a fronte di un bisogno concreto osservato nell'uso reale.
- Non assumere che un agente worker (Codex/Gemini) committi da solo sul proprio branch anche
  se istruito a "non toccare main" — nella pratica ha lasciato i file non tracciati nel
  worktree. Il coordinatore deve istruire esplicitamente a committare e verificare prima del
  merge (vedi Handoff in `activeContext.md`, 2026-07-10).
