# Active Context

> Altissima frequenza di aggiornamento. Importato automaticamente a ogni sessione.

## Cosa si sta facendo ora

Handoff multi-agente chiuso (vedi sotto): sottotask di Codex e Gemini rivisti, committati e
mergiati in `main`. Nessun task di codice aperto al momento.

## Ultima decisione presa

Cattura automatica da commit scelta come "solo commit marcati esplicitamente con un trailer",
non mirror di ogni commit (§1: sarebbe rumore, non conoscenza). PR-based capture e un comando
`cairn hook install` restano deferred (vedi `techContext.md`).

## Handoff chiuso — 2026-07-10, Codex e Gemini → main

Primo dispatch multi-agente del progetto, portato a termine. Note operative per il prossimo:

- **Bug Orca**: il repo Cairn è registrato con `kind: "folder"` (probabilmente perché
  aggiunto a Orca prima del primo commit) — `orca worktree create --repo id:dee56f51-...` NON
  isola davvero, riusa la stessa directory. Verificare `orca repo show --repo id:dee56f51-...`
  prima di un prossimo dispatch: se `kind` è ancora `"folder"`, creare il worktree a mano
  (`git worktree add`) e registrarlo con `orca repo add --path <nuovo-path>` (rileva `git`
  correttamente). Non risolvibile via CLI (nessun comando per correggere `kind` in place).
- **Codex e Gemini non fanno commit da soli** anche quando esplicitamente istruiti a "non
  fare push né merge" — hanno scritto i file nel worktree ma lasciato tutto non tracciato. Il
  coordinatore ha dovuto committare lui stesso su ciascun branch prima del merge. Da tenere a
  mente per il prossimo dispatch: non assumere che "non toccare main" implichi "committa sul
  tuo branch" — vanno probabilmente istruiti esplicitamente a committare.
- **Codex non aveva una shell disponibile** nel suo ambiente (MCP `shell-tool` fallito
  all'avvio) — non ha potuto eseguire `go build`/`go vet`/`go test` da solo. Il coordinatore
  ha validato al posto suo prima del merge. Verificare se è una limitazione sistemica di
  questo setup Orca+Codex o solo di questa sessione, prima di assumere di poterlo evitare.
- Merge fatti (`--no-ff`, uno per sottotask): `cmd/cairn/main_test.go` (7 test, tutti verdi
  dopo build/vet/test del coordinatore) e `memory-bank/rfc/0001-livello-di-conoscenza.md`.
  Nessun conflitto. Worktree e branch locali rimossi dopo il merge
  (`git worktree remove` + `git branch -d`). Le due voci repo temporanee in Orca
  (`7adac43e-...`, `83b94d08-...`) non sono rimovibili da CLI (nessun `orca repo rm`) — restano
  come voci stale che puntano a directory ormai cancellate, da pulire a mano dalla UI se dà
  fastidio, non urgente.

## Prossimo passo

Nessun task di codice aperto. In sospeso: push su `main` (non ancora fatto, da confermare),
ed eventuali limature di stile sull'RFC-0001 (attribuzione, tono) se l'utente le vuole.

## Blocchi/domande aperte

- Nessuno.