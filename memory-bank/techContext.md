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

## Restrizioni ambientali

- [Es. versione minima runtime, OS supportati, rete disponibile in sandbox, ecc.]

## Debito tecnico noto

| Elemento | Impatto | Quando affrontarlo |
|---|---|---|
| Nessun test automatico (solo validazione manuale) | Rischio di regressioni silenziose sugli incrementi futuri | Prima del terzo incremento funzionale, o appena la validazione manuale richiede più di ~2 minuti |