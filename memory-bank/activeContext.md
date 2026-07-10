# Active Context

> Altissima frequenza di aggiornamento. Importato automaticamente a ogni sessione.

## Cosa si sta facendo ora

Primo dispatch multi-agente reale del progetto (§7): Claude Code (coordinatore) ha assegnato
due sottotask indipendenti a Codex e Gemini via Orca, 2026-07-10. Vedi blocco Handoff sotto.

## Ultima decisione presa

Fix minimo e mirato sul bug ordine flag (vedi `progress.md`) chiuso in precedenza in questa
sessione. Relazioni tipizzate tra note restano rimandate (vedi `progress.md` § Deferred).

## Handoff — 2026-07-10, da Claude Code (coordinatore) a Codex e Gemini (worker)

**Scoperta operativa importante prima del dispatch:** il repo Cairn era registrato in Orca
con `kind: "folder"` invece di `"git"` (probabilmente perché aggiunto a Orca prima del primo
`git init`/commit) — `orca worktree create --repo id:dee56f51-...` su questo repo NON crea un
vero worktree isolato, riusa la stessa directory della sessione principale (verificato:
`git worktree list` non mostrava nulla di nuovo dopo la creazione). **Se in futuro si rifà un
dispatch su Cairn tramite Orca, verificare prima `orca repo show --repo id:dee56f51-8...` →
se `kind` è ancora `"folder"`, non usare la scorciatoia `orca worktree create --repo id:...`:
creare invece il worktree git a mano (`git worktree add ...`) e registrarlo come repo
separato con `orca repo add --path <nuovo-path>`** (questo rileva correttamente `kind: git`).
Non ho tentato di correggere la registrazione originale del repo Cairn in Orca — nessun
comando CLI visto lo permette senza rimuoverlo e riaggiungerlo, azione che avrebbe effetti
sugli altri worktree Orca già collegati a quell'id repo, non presa senza chiedere.

**Sottotask 1 — Codex**, worktree git isolato `/Users/giuseppetauro/Development/Cairn-codex-tests`
(branch `codex/cairn-tests`, repo Orca id `7adac43e-fadd-466d-8c46-8de3e57e7a4d`, terminale
`term_1c08a39b-b642-41d0-a301-2d68dc6a530b`): scrivere `cmd/cairn/main_test.go` (solo stdlib
`testing`) che copre gli scenari già validati a mano su `add`/`log`, tag, filtro, e il fix
sull'ordine dei flag. Istruito a non toccare `memory-bank/`, non fare push/merge, e a
proporre un piano breve prima di scrivere codice (sta seguendo la pipeline correttamente).

**Sottotask 2 — Gemini**, worktree git isolato
`/Users/giuseppetauro/Development/Cairn-gemini-rfc0001` (branch `gemini/rfc-0001`, repo Orca
id `83b94d08-0be8-4fb6-8dac-0f7fad54564b`, terminale `term_dfc50609-69be-4970-af8b-250459edf0fc`):
bozza di `memory-bank/rfc/0001-livello-di-conoscenza.md`, basata su AGENTS.md + memory-bank +
`.cairn/log.jsonl` come prova concreta vissuta. Istruito a non toccare codice né
`activeContext.md`/`progress.md`, non fare push/merge.

**Cosa NON toccare finché questo Handoff è aperto:** i due branch `codex/cairn-tests` e
`gemini/rfc-0001` e le rispettive directory sono di competenza esclusiva dell'agente
assegnato. Il merge in `main` resta un gate umano (§7) — nessuno dei due deve farlo da solo.

## Prossimo passo

Rivedere l'output di Codex e Gemini quando segnalano di aver finito (o su richiesta
dell'utente), poi valutare merge in `main` — con conferma esplicita dell'utente, come da gate
del §7. Dopo il merge, il coordinatore (io) aggiorna `activeContext.md`/`progress.md` e può
rimuovere i due worktree (`git worktree remove ...` + `orca repo` cleanup).

## Blocchi/domande aperte

- Nessuno sullo stato del repo principale. (`.cairn/` è versionato dal 2026-07-10, vedi commit
  precedenti.) Aperto solo l'esito dei due sottotask in corso sopra.