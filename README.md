# Kit CLAUDE.md / AGENTS.md modulare — Cairn

Governance per agenti di coding, pensata per funzionare sia con un solo agente (Claude Code)
sia con più agenti in parallelo (Claude Code + Codex + Gemini, es. via un ADE come Orca).

## Struttura

```
.
├── AGENTS.md                      # FONTE DI VERITÀ condivisa — Codex la legge nativamente,
│                                   # Gemini CLI via .gemini/settings.json, Claude Code via @import
├── CLAUDE.md                      # adapter sottile per Claude Code: importa @AGENTS.md
│                                   # + le poche cose davvero specifiche di Claude Code
├── .claude/
│   ├── settings.json               # hook PreToolUse: blocca comandi distruttivi (solo Claude Code)
│   └── rules/
│       ├── security.md            # meccanismo dell'hook + nota "non protegge Codex/Gemini"
│       └── multi-agent.md         # ruolo di Claude Code (coordinatore/worker) nell'orchestrazione
├── .gemini/
│   └── settings.json               # dice a Gemini CLI di leggere AGENTS.md invece di GEMINI.md
└── memory-bank/
    ├── projectBrief.md            # statico — requisiti e confini del progetto
    ├── productContext.md          # bassa frequenza — perché esiste il prodotto
    ├── systemPatterns.md          # media frequenza — architettura e pattern obbligatori
    ├── techContext.md             # bassa frequenza — stack, comandi, debito tecnico
    ├── activeContext.md           # ALTA frequenza — stato della sessione corrente (importato)
    └── progress.md                # ALTA frequenza — avanzamento e bug noti (importato)
```

## Come funziona la condivisione tra i tre agenti

| Agente | Come legge `AGENTS.md` |
|---|---|
| **Codex** | Nativamente, senza configurazione. |
| **Gemini CLI** | Tramite `.gemini/settings.json` (già incluso in questo kit), che gli dice di usare `AGENTS.md` come file di contesto al posto del suo `GEMINI.md` di default. |
| **Claude Code** | Tramite `@AGENTS.md` in cima a `CLAUDE.md`. |

**Regola d'oro anti-drift:** le regole vere si scrivono **solo** in `AGENTS.md`. `CLAUDE.md` e
`.claude/rules/*.md` aggiungono esclusivamente ciò che è genuinamente specifico di Claude Code
(l'hook, il comportamento dopo `/compact`, il ruolo nell'orchestrazione). Se ti accorgi di
scrivere la stessa regola in due posti, quella regola andava in `AGENTS.md`.

## Come iniziare

1. Copia l'intera struttura nella radice del repo di Cairn.
2. Compila `AGENTS.md` e i sei file in `memory-bank/` con i dati reali del progetto —
   sostituisci i placeholder `[...]`.
3. **Non creare un file `GEMINI.md`** nella stessa cartella: se esiste insieme ad `AGENTS.md`,
   Gemini CLI dà precedenza a `GEMINI.md` e ignora la configurazione in `.gemini/settings.json`.
4. Committa tutto in Git.

## Come lavorare con Claude Code + Codex + Gemini in parallelo (es. con Orca)

1. **Login una tantum** in ciascun CLI sulla tua macchina (`claude`, `codex`, `gemini`) — gli
   ADE come Orca usano le tue sottoscrizioni già configurate, non serve un account separato.
2. **Un worktree Git per sottotask**, non per agente fisso: se oggi tocca a Codex il modulo X
   e domani a Gemini il modulo Y, ognuno lavora nel proprio worktree isolato.
3. **Scegli un coordinatore** per lo spec corrente (di norma Claude Code, che segue
   `.claude/rules/multi-agent.md`). Il coordinatore scompone lo spec in sottotask con
   dipendenze e li assegna, invece di far lavorare tutti e tre sullo stesso file
   contemporaneamente.
4. **`memory-bank/activeContext.md` e `progress.md` li scrive solo il coordinatore**, dopo il
   merge di ciascun worktree — gli altri due agenti li trattano come sola lettura durante il
   loro sottotask (regola in `AGENTS.md` §5).
5. **Handoff esplicito**: quando un sottotask passa da un agente all'altro, chi lo consegna
   scrive un blocco "Handoff" in `activeContext.md` (da chi, a chi, cosa aspettarsi, cosa non
   toccare) prima che l'altro agente parta.
6. **Gate umano obbligatorio** prima di: merge in `main`, modifiche a `systemPatterns.md`,
   qualunque azione tra i divieti di `AGENTS.md` §4 — indipendentemente da quale agente la
   propone.
7. **La sicurezza non è automaticamente condivisa tra i tre.** L'hook che blocca i comandi
   distruttivi in questo kit vale solo per Claude Code. Se dai a Codex o Gemini accesso a
   shell/file, replica lo stesso divieto con il meccanismo di sandboxing/permessi proprio di
   ciascuno strumento — controlla la loro documentazione, non riusare la configurazione di
   Claude Code assumendo che valga anche per loro.

## Perché questa struttura

- `AGENTS.md` come standard cross-tool evita la deriva tipica del mantenere più file di
  istruzioni paralleli che iniziano a divergere: una fonte di verità, adapter sottili per
  strumento.
- I sei file di `memory-bank/` restano divisi per **frequenza di aggiornamento**: solo i due
  "caldi" (`activeContext.md`, `progress.md`) vengono importati automaticamente; gli altri
  quattro si leggono on-demand.
- I divieti di sicurezza sono raddoppiati con un controllo deterministico (hook) per lo
  strumento che lo supporta in questo kit (Claude Code): un'istruzione in un prompt è
  probabilistica, un hook con `exit 2` no.

## Cosa non è incluso (di proposito)

L'idea di un hub di telemetria/osservabilità che intercetta gli span degli agenti e li instrada
verso APM esistenti (Datadog, Honeycomb, ecc.) è un'estensione infrastrutturale separata dalla
governance di base — non inclusa qui. Se ti serve davvero, la strutturo come progetto a sé.
