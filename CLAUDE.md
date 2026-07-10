# CLAUDE.md — Cairn (adapter per Claude Code)

@AGENTS.md

> Le regole vere vivono in `AGENTS.md` — lette anche da Codex (nativamente) e da Gemini CLI
> (via `.gemini/settings.json`). Qui sotto solo ciò che è specifico di Claude Code. Non
> duplicare qui contenuto già presente in `AGENTS.md`.

## Import automatici aggiuntivi (caricati a ogni sessione)

@memory-bank/activeContext.md
@memory-bank/progress.md

> Nota: importiamo solo i due file "caldi". Gli altri quattro di `memory-bank/` (vedi
> `AGENTS.md` §2) vanno letti on-demand con lo strumento di lettura file, non importati,
> altrimenti gonfiano ogni sessione con contenuto che cambia raramente.

## Specifico di Claude Code

- Se disponibile, usa il **Plan Mode nativo** del client per la fase Plan (blocco reale di
  scrittura, non solo un'istruzione interpretabile).
- Il divieto sui comandi distruttivi (`AGENTS.md` §4) è applicato anche a livello di **hook**
  in `.claude/settings.json` — dettaglio in `.claude/rules/security.md`.
- Ruolo nell'orchestrazione multi-agente (es. con Orca) → `.claude/rules/multi-agent.md`.
- Dopo `/compact`, Claude Code ricarica questo file da disco ma **non** rilegge
  automaticamente `.claude/rules/*.md` finché non tocchi un file in quella sottocartella: se
  un'istruzione sembra "sparita" dopo una compattazione, è lì il motivo più probabile.

## Conferma di avvio

Dopo aver letto `AGENTS.md` + i due import, stampa `[MEMORIA PROGETTO: ATTIVA]` prima di
iniziare il task.
