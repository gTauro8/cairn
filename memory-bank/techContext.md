# Tech Context — Cairn

> Aggiornamento: basso, solo con migrazioni o aggiornamenti di dipendenze.

## Stack (adottato il 2026-07-10, senza il passaggio di validazione manuale originariamente
## pianificato — vedi ADR in `systemPatterns.md`)

- Linguaggio: Go
- Componenti previsti: SQLite, Tree-sitter
- Non sono state riconsiderate rispetto all'ipotesi iniziale: sono lo stack di partenza per
  l'implementazione, da correggere in corsa se il vincolo del §2 di `AGENTS.md` segnala attrito.

## Comandi

- Build: `go build ./...`
- Test: `go test ./...` (`cmd/cairn/main_test.go`, copre add/log/tag/files/ordine flag/cwd
  esatta — vedi `progress.md`)
- Lint: `go vet ./...` (nessun linter esterno introdotto ancora)
- Run locale: `go run ./cmd/cairn add [--tags a,b,c] [--files a,b] <testo>` oppure
  `go run ./cmd/cairn log [--tag x] [--file p]`

## Cattura automatica da commit (hook `post-commit`)

Convenzione: un commit i cui trailer contengono `Cairn-Note: true` genera automaticamente
una nota `cairn` (testo = subject del commit + hash corto, tag = valore di `Cairn-Tags:` più
il tag automatico `git`, `files` = i file effettivamente cambiati nel commit, popolato in
automatico via `git diff-tree --name-only`). Commit senza quel trailer: nessun effetto,
nessun output — non è un mirror di ogni commit, solo quelli marcati esplicitamente come
conoscenza da ricordare.

Esempio di messaggio di commit che genera una nota:
```
fix: cairn add rileva --tags dopo il testo

Spiegazione del bug e della scelta di design...

Cairn-Note: true
Cairn-Tags: bug,ux
```

Setup (una volta per clone/macchina — `core.hooksPath` è locale, non versionato):
```
git config core.hooksPath .githooks
```
Lo script è in `.githooks/post-commit`, versionato. Richiede il binario `cairn` sul `PATH`
oppure compilato in `<repo>/cairn` (`go build -o cairn ./cmd/cairn`) — altrimenti stampa un
warning su stderr e non fa fallire il commit.

Fuori scope (deferred, 2026-07-10): cattura da PR (richiederebbe una GitHub Action, non un
hook locale — infrastruttura diversa, si valuta solo se il trailer sui commit si dimostra
utile nell'uso reale); un comando `cairn hook install` che automatizzi il `git config` sopra.

**Nota tecnica (bug reale trovato in dogfooding, 2026-07-10):** il rilevamento NON usa
`git interpret-trailers`, cerca le righe `Cairn-Note:`/`Cairn-Tags:` ovunque nel messaggio con
un grep diretto. `git interpret-trailers --parse` riconosce solo l'ultimo paragrafo contiguo
come blocco trailer: una riga vuota tra `Cairn-Tags:` e un trailer successivo (es.
`Co-Authored-By:`) fa ignorare silenziosamente i trailer Cairn-*. Il primo commit reale con
questa feature (f7b5f89) ha perso la nota per questo motivo, prima del fix (39771a9).

## Restrizioni ambientali

- Sviluppato/validato su macOS (darwin/arm64), Go 1.25.3. Nessuna dipendenza da rete per il
  funzionamento di `cairn` stesso (solo build/toolchain Go richiedono l'ambiente di sviluppo
  standard). Non ancora testato su Linux/Windows — nessun uso di API specifiche di macOS nel
  codice, ma non verificato attivamente.
- L'hook `post-commit` richiede una shell POSIX (`/bin/sh`) e i comandi `git`/`grep` standard.

## Debito tecnico noto

| Elemento | Impatto | Quando affrontarlo |
|---|---|---|
| Nessun linter esterno oltre `go vet` | Possibili incoerenze di stile non catturate | Se il codice cresce oltre un singolo file, valutare `golangci-lint` |
| Nessuna verifica automatica multi-OS (solo macOS validato) | Un comportamento specifico di macOS potrebbe passare inosservato | Se/quando il progetto acquisisce contributor su altri OS |
| `cairn hook install` non esiste: setup dell'hook è un `git config` manuale per clone | Rischio che l'hook resti disattivato su una macchina/clone senza che nessuno se ne accorga | Se il progetto acquisisce più contributor/macchine (vedi `techContext.md` § Cattura automatica) |