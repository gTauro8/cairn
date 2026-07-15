# Release di Cairn

Questa procedura prepara gli artefatti ma non autorizza automaticamente tag, push o release
GitHub: sono gate umani separati.

## Prerequisiti

1. Il protocollo in `memory-bank/dogfooding-v0.1.md` ha raggiunto un esito esplicito.
2. `CHANGELOG.md` descrive la versione e non la indica più come "non rilasciata".
3. Il working tree contiene soltanto le modifiche intenzionali della release.

## Verifica

```sh
make verify
cairn check
```

## Artefatti

```sh
make release VERSION=v0.1.0
./dist/cairn_v0.1.0_darwin_arm64 version
```

Il target produce binari CGO-disabled per macOS/Linux su amd64/arm64, usando `-trimpath` e
`-buildvcs=false`. Generare e pubblicare i checksum con lo strumento SHA-256 disponibile
nell'ambiente di release.

## Gate manuali

Dopo aver verificato gli artefatti:

1. chiedere conferma per il commit finale su `main`;
2. chiedere conferma separata per creare il tag `v0.1.0`;
3. chiedere conferma separata per push e pubblicazione della release GitHub.

Non riutilizzare un tag esistente e non eseguire force-push.
