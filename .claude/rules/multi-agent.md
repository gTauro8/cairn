# Regola (Claude Code): ruolo nell'orchestrazione multi-agente

> Il contratto di coordinamento condiviso con Codex e Gemini è in `AGENTS.md` §5 — quella è
> la fonte di verità che vale per tutti e tre gli agenti. Qui c'è solo come Claude Code lo
> esegue in pratica dentro un ADE multi-agente come Orca.

## Se Claude Code è il coordinatore del task

- Usa `/orchestrate` (o la relativa skill di orchestrazione) per scomporre lo spec in
  sottotask con dipendenze (DAG), invece di dispatchare lavoro non strutturato.
- Assegna i sottotask a Codex/Gemini nei rispettivi worktree — non scrivere direttamente nei
  worktree di un altro agente.
- Attendi il segnale di completamento di un sottotask prima di considerarlo chiuso e passare
  al successivo nel DAG.
- Inserisci un gate decisionale (conferma umana esplicita) prima di: merge in `main`,
  modifiche a `memory-bank/systemPatterns.md`, qualunque azione elencata in `AGENTS.md` §4.
  Questo vale a prescindere da quale agente ha prodotto la modifica.
- Alla chiusura di ciascun sottotask ricevuto in consegna, sei tu (coordinatore) ad aggiornare
  `activeContext.md`/`progress.md` dopo il merge — non lasciare che ogni agente scriva la
  memoria condivisa per conto proprio e in parallelo.

## Se un altro agente coordina e Claude Code è worker su un sottotask

- Tratta `memory-bank/activeContext.md` e `progress.md` come sola lettura per la durata del
  sottotask: non sovrascriverli. Riporta l'esito al coordinatore e lascia che sia lui a
  consolidare la memoria condivisa dopo il merge.
- Resta dentro i confini del tuo worktree: nessuna operazione Git che tocchi `main` o altri
  worktree senza che il coordinatore lo richieda esplicitamente.
- Se ricevi un blocco "Handoff" in `activeContext.md` scritto da un altro agente, leggilo
  prima di iniziare: contiene cosa aspettarti e cosa non toccare.

## Nota

L'hook di sicurezza in `.claude/rules/security.md` resta attivo indipendentemente dal ruolo
(coordinatore o worker): l'orchestrazione multi-agente non è un'eccezione ai divieti di
`AGENTS.md` §4.
