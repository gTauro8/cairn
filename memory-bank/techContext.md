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
- Test: `go test ./...` (nessun test automatizzato ancora presente — solo validazione manuale
  finora, vedi `progress.md`)
- Lint: `go vet ./...` (nessun linter esterno introdotto ancora)
- Run locale: `go run ./cmd/cairn add [--tags a,b,c] <testo>` oppure
  `go run ./cmd/cairn log [--tag x]`

## Cattura automatica da commit (hook `post-commit`)

Convenzione: un commit i cui trailer contengono `Cairn-Note: true` genera automaticamente
una nota `cairn` (testo = subject del commit + hash corto, tag = valore di `Cairn-Tags:` più
il tag automatico `git`). Commit senza quel trailer: nessun effetto, nessun output — non è
un mirror di ogni commit, solo quelli marcati esplicitamente come conoscenza da ricordare.

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

## Restrizioni ambientali

- [Es. versione minima runtime, OS supportati, rete disponibile in sandbox, ecc.]

## Debito tecnico noto

| Elemento | Impatto | Quando affrontarlo |
|---|---|---|
| Nessun test automatico (solo validazione manuale) | Rischio di regressioni silenziose sugli incrementi futuri | Prima del terzo incremento funzionale, o appena la validazione manuale richiede più di ~2 minuti |