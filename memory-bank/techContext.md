# Tech Context — Cairn

> Aggiornamento: basso, solo con migrazioni o aggiornamenti di dipendenze.

## Stack (adottato il 2026-07-10, senza il passaggio di validazione manuale originariamente
## pianificato — vedi ADR in `systemPatterns.md`)

- Linguaggio: Go
- Componenti previsti: SQLite, Tree-sitter
- Non sono state riconsiderate rispetto all'ipotesi iniziale: sono lo stack di partenza per
  l'implementazione, da correggere in corsa se il vincolo del §2 di `AGENTS.md` segnala attrito.

## Comandi

- Build locale: `make build` (binario `./cairn`, versione `dev` per default)
- Verifica completa: `make verify`
- Release locale: `make release VERSION=v0.1.0` (macOS/Linux, amd64/arm64, CGO disabilitato)
- Test: `go test ./...` (`cmd/cairn/main_test.go`, copre add/log/tag/files/ordine flag/cwd
  esatta/output JSON/check/hook install e cattura Git — vedi `progress.md`)
- Lint: `go vet ./...` (nessun linter esterno introdotto ancora)
- Run locale: `go run ./cmd/cairn add [--tags a,b,c] [--files a,b] [--source s]
  [--commit hash] <testo>`,
  `go run ./cmd/cairn log [--tag x] [--file p] [--json]`,
  `go run ./cmd/cairn check`, `go run ./cmd/cairn hook install`, oppure
  `go run ./cmd/cairn version`

`cairn log --json` restituisce le righe JSONL originali, senza rimarshal: in questo modo non
perde campi additivi sconosciuti a una versione precedente del binario. `cairn check` è
read-only e restituisce exit non-zero se trova JSON non valido, campi obbligatori mancanti o
vuoti, oppure riferimenti in `files` che non esistono nella cwd.

## Cattura automatica da commit (hook `post-commit`)

Convenzione: un commit i cui trailer contengono `Cairn-Note: true` genera automaticamente
una nota `cairn` (testo = subject del commit + hash corto, tag = valore di `Cairn-Tags:` più
il tag automatico `git`, `files` = i file effettivamente cambiati nel commit, `source` =
`git`, `commit` = hash completo). Commit senza quel trailer: nessun effetto, nessun output —
non è un mirror di ogni commit, solo quelli marcati esplicitamente come conoscenza da
ricordare. Le note create direttamente da `cairn add` usano `source: manual` per default.

Esempio di messaggio di commit che genera una nota:
```
fix: cairn add rileva --tags dopo il testo

Spiegazione del bug e della scelta di design...

Cairn-Note: true
Cairn-Tags: bug,ux
```

Setup (una volta per clone/macchina — `core.hooksPath` è locale, non versionato):
```
cairn hook install
```
Il comando va eseguito dalla root Git esatta, installa il wrapper versionato e configura
`core.hooksPath=.githooks`. È idempotente e non sovrascrive hook/configurazioni differenti.
Il wrapper richiede `cairn` sul `PATH` oppure in `<repo>/cairn`; delega parsing e cattura a
`cairn hook run`, convertendo ogni errore in warning per non far fallire il commit.

Fuori scope: cattura da PR (infrastruttura CI distinta, da valutare solo dopo il dogfooding
dell'hook locale).

**Nota tecnica (bug reale trovato in dogfooding, 2026-07-10):** il parser Go cerca le righe
`Cairn-Note:`/`Cairn-Tags:` ovunque nel messaggio e non usa `git interpret-trailers --parse`.
Quest'ultimo riconosceva solo l'ultimo paragrafo contiguo e aveva perso una nota quando un
`Co-Authored-By:` era separato da una riga vuota (commit di fix 39771a9).

## Restrizioni ambientali

- Sviluppato e validato localmente su macOS (darwin/arm64), Go 1.25.3. La CI è configurata per
  macOS e Linux; l'esecuzione remota inizierà dopo il commit/push del workflow.
- L'hook `post-commit` richiede una shell POSIX (`/bin/sh`) e Git. Windows non è supportato
  nella v0.1.

## Debito tecnico noto

| Elemento | Impatto | Quando affrontarlo |
|---|---|---|
| Nessun linter esterno oltre `go vet` | Possibili incoerenze di stile non catturate | Se il codice cresce oltre un singolo file, valutare `golangci-lint` |
| Wrapper hook POSIX | Windows non può usare la cattura automatica v0.1 | Valutare solo a fronte di utenti Windows reali |
| Nota creata dopo il commit | `.cairn/log.jsonl` resta modificato fino al commit successivo | Misurare nel dogfooding; cambiare hook solo se crea attrito reale |
