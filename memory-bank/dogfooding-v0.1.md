# Dogfooding v0.1

## Stato

Non iniziato. Questo documento definisce il gate empirico prima della release v0.1.0 e prima
di investire in un'estensione IDE.

## Durata minima

Almeno 10 sessioni di sviluppo distribuite su 14 giorni di calendario. Il periodo parte dal
primo uso intenzionale della build candidata, non dalla scrittura di questo documento.

## Cosa registrare

Per ogni sessione aggiungere una riga alla tabella senza telemetria automatica.

| Data | Decisione catturabile? | Catturata nella sessione? | Consultazione reale? | `cairn check` | Attrito/note |
|---|---|---|---|---|---|

Una "decisione catturabile" è una scelta, alternativa scartata, scoperta di debug o vincolo
che sarebbe costoso ricostruire dal solo codice. Non contare commit ordinari o riassunti già
evidenti dal diff.

## Criteri di successo

Tutti i criteri devono essere soddisfatti:

1. almeno l'80% delle decisioni catturabili viene registrato nella stessa sessione;
2. almeno tre consultazioni nascono da un bisogno reale, e almeno due evitano di ricostruire
   manualmente il contesto;
3. `cairn check` resta verde oppure ogni segnalazione risulta utile e viene gestita entro la
   sessione successiva;
4. non più di una sessione su dieci abbandona la cattura per attrito;
5. il comportamento `post-commit` non lascia sistematicamente note non versionate o un
   working tree fastidiosamente sporco.

## Esiti possibili

- **Promuovi v0.1.0:** tutti i criteri sono soddisfatti.
- **Itera sulla CLI/hook:** la cattura o la freschezza falliscono; correggere il problema e
  ripetere il periodo interessato.
- **Riconsidera il prodotto:** le note si accumulano ma non vengono consultate.

## Gate estensione IDE

L'estensione VS Code si pianifica soltanto se il dogfooding mostra che le note sono utili ma
difficili da scoprire accanto al codice. In quel caso deve restare un client sottile di
`cairn log --json`, senza storage o regole di dominio proprie.
