# Active Context

> Altissima frequenza di aggiornamento. Importato automaticamente a ogni sessione.

## Cosa si sta facendo ora

Blocco tecnico v0.1 completato e pronto per il commit su `main`: `cairn hook
install/run`, wrapper `post-commit` minimale, `cairn version`, README di prodotto, licenza
Apache-2.0, Makefile/release cross-platform, CI macOS/Linux, changelog, procedura di release e
protocollo di dogfooding. `make verify` e quattro cross-build verdi; onboarding validato
end-to-end in scratch con commit marcato e non marcato. Il gate al commit è stato ricevuto,
ma il sandbox corrente non può creare `.git/index.lock`; nessun file è stato staged.

## Ultima decisione presa

La logica di cattura Git vive ora nel binario Go (`cairn hook run`); lo script installato è un
wrapper POSIX stabile che converte ogni errore in warning per non bloccare mai un commit.
`hook install` rifiuta di sovrascrivere configurazioni o hook diversi: l'alternativa di
generare l'intero hook shell avrebbe duplicato la logica e aumentato il rischio di drift.

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

Eseguire il commit autorizzato da un ambiente con scrittura su `.git`, escludendo `.DS_Store`,
poi avviare il dogfooding: almeno 10 sessioni su 14 giorni, secondo
`memory-bank/dogfooding-v0.1.md`. Nessuna estensione IDE e nessuna release v0.1.0 prima
dell'esito.

## Blocchi/domande aperte

- Il sandbox della sessione consente di leggere `.git` ma non di creare `.git/index.lock`:
  `git add` fallisce prima di modificare l'indice. Serve un terminale con permesso Git.
- Dogfooding non iniziato; v0.1.0 non va dichiarata rilasciata prima dell'esito.
