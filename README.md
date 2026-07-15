# Cairn

Cairn conserva la conoscenza che vive attorno al codice: decisioni architetturali, alternative
scartate, vincoli, incidenti e debito tecnico. Le note sono locali al repository, versionate
con Git e consultabili senza servizi esterni o modelli AI.

> Stato: v0.1 in preparazione. Il formato e i comandi sono già utilizzabili, ma il criterio
> principale — cattura e freschezza con attrito quasi zero — è ancora in dogfooding.

## Cosa non è

Cairn non è un chatbot, un framework AI, un vector database o un sistema di memoria per
agenti. Umani, agenti, IDE e CI sono client di pari livello dello stesso log append-only.

## Build locale

Richiede Go 1.25 o successivo e Git.

```sh
make build
./cairn version
```

Il binario risultante è `./cairn`. Per usare l'hook Git, lascia il binario nella root del
repository oppure installalo in una directory presente nel `PATH`.

## Quick start

Registra una decisione manuale:

```sh
cairn add --tags architecture --files cmd/cairn/main.go "JSONL resta lo storage della v0.1"
```

Consulta tutte le note o applica filtri combinabili:

```sh
cairn log
cairn log --tag architecture
cairn log --file cmd/cairn/main.go
cairn log --tag architecture --file cmd/cairn/main.go
```

Per IDE, CI e altri client è disponibile l'output JSON Lines originale:

```sh
cairn log --json
```

Verifica integrità e freschezza deterministica:

```sh
cairn check
```

`check` restituisce un exit code non-zero se trova JSON non valido, campi obbligatori vuoti,
un riferimento a un file inesistente oppure una nota con `source: git` priva di commit. Non
riscrive mai il log e non tenta di indovinare rinomine o obsolescenza semantica.

## Cattura da commit Git

Installa l'hook dalla root esatta del repository:

```sh
cairn hook install
```

Il comando è idempotente. Non sovrascrive un `core.hooksPath` diverso o un hook `post-commit`
personalizzato: in quei casi si ferma e lascia la configurazione invariata.

Marca soltanto i commit che contengono conoscenza da conservare:

```text
refactor: separa il parsing dalla persistenza

La separazione mantiene sostituibile il formato senza cambiare la CLI.

Cairn-Note: true
Cairn-Tags: architecture,storage
```

L'hook crea una nota con:

- testo formato dal subject e dall'hash corto;
- tag dichiarati più `git`;
- file modificati dal commit;
- `source: git` e hash completo in `commit`.

I commit senza `Cairn-Note: true` o `Cairn-Note: yes` restano silenziosi. L'hook non fa mai
fallire un commit: se Cairn non è disponibile emette soltanto un warning.

### Conseguenza di `post-commit`

La nota nasce immediatamente dopo il commit che l'ha generata, quindi la modifica a
`.cairn/log.jsonl` non appartiene a quel commit: verrà versionata nel commit successivo. È un
compromesso intenzionalmente visibile durante il dogfooding; se causa note dimenticate o
working tree perennemente sporchi, il meccanismo dovrà essere rivisto prima della v1.

## Formato dei dati

Le note vivono in `.cairn/log.jsonl`, una riga JSON per nota. I campi correnti sono:

```json
{"id":"a1b2c3","ts":"2026-07-15T10:00:00Z","text":"decisione","tags":["architecture"],"files":["main.go"],"source":"manual","commit":""}
```

`id`, `ts` e `text` sono obbligatori. Gli altri campi sono additivi e opzionali, così le
versioni nuove non richiedono migrazioni delle righe esistenti. Le note manuali usano
`source: manual`; l'hook usa `source: git` e valorizza `commit`.

Cairn cerca `.cairn/` soltanto nella directory corrente: non risale implicitamente alle
directory padre. Esegui quindi i comandi dalla root del progetto.

## Sviluppo

```sh
make verify
```

Esegue test, `go vet` e build. La CI ripete gli stessi controlli su Linux e macOS. La procedura
di release è documentata in [RELEASING.md](RELEASING.md); il protocollo di validazione reale
è in [memory-bank/dogfooding-v0.1.md](memory-bank/dogfooding-v0.1.md).

## Limiti attuali

- nessuna relazione tipizzata tra note;
- nessuna ricerca full-text o UI IDE;
- hook Git basato su shell POSIX, quindi Windows non è ancora supportato;
- file rinominati o eliminati vengono segnalati da `cairn check` e richiedono valutazione
  umana;
- nessun server, database, sincronizzazione cloud o dipendenza da LLM.

## Licenza

Apache License 2.0. Vedi [LICENSE](LICENSE).
