# Regola (Claude Code): applicazione dei divieti via hook

> I divieti sono elencati in `AGENTS.md` §4 — quel file è la fonte di verità condivisa con
> Codex e Gemini. Qui c'è solo il meccanismo con cui Claude Code li applica in modo
> deterministico, non probabilistico.

## Perché un hook e non solo il prompt

Un divieto scritto solo in linguaggio naturale è probabilistico: la sua tenuta cala quanto più
la sessione si allunga. Un hook `PreToolUse` gira fuori dal controllo del modello — se esce con
codice 2, l'azione viene bloccata a prescindere da cosa "decide" il modello in quel turno.

## Esempio (vedi anche `.claude/settings.json` in questo kit)

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "CMD=$(cat | jq -r '.tool_input.command // empty'); echo \"$CMD\" | grep -qE 'rm -rf|DROP TABLE|TRUNCATE|push --force|reset --hard' && { echo '{\"hookSpecificOutput\":{\"hookEventName\":\"PreToolUse\",\"permissionDecision\":\"deny\",\"additionalContext\":\"Comando distruttivo bloccato dalla policy di progetto (AGENTS.md §4)\"}}'; exit 2; } || exit 0"
          }
        ]
      }
    ]
  }
}
```

## Regola pratica per decidere cosa mettere dove

- Violazione che bloccherebbe una merge in CI o causerebbe un danno reale → hook o CI, non
  solo un file di istruzioni.
- Violazione che farebbe solo alzare un sopracciglio a un reviewer umano → un file di
  istruzioni (`AGENTS.md`) basta.

## Nota multi-agente

Questo hook copre solo Claude Code. Se Codex o Gemini hanno accesso a shell/file nello stesso
progetto, lo stesso divieto va replicato con il meccanismo di sandboxing/permessi proprio di
ciascuno strumento — non assumere che l'hook di Claude Code protegga anche loro.
