# Changelog

Il progetto segue [Semantic Versioning](https://semver.org/). Le modifiche non ancora
rilasciate restano sotto `Unreleased`.

## Unreleased

### Added

- CLI append-only con `add`, `log`, filtri per tag/file e output JSON Lines.
- `check` per integrità del log e riferimenti a file mancanti.
- Provenienza strutturata `manual`/`git` con hash completo del commit.
- `hook install` per configurare in sicurezza la cattura da commit marcati.
- `version`, build riproducibili e CI Linux/macOS.

### Changed

- La logica dell'hook Git vive nel binario; lo script `post-commit` è un wrapper minimale.

## 0.1.0 - non rilasciata

Prima release prevista dopo il completamento del dogfooding v0.1.
