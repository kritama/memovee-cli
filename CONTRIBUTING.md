# Contributing

## Toolchain

Memovee CLI targets Go 1.25. Use the version pinned in `mise.toml`.

```sh
mise install
go test ./...
```

## Branching

This repository uses Git Flow:

- `main` is the production branch;
- `develop` is the integration branch;
- new work uses `feature/<name>` branches from `develop`;
- releases use `release/<version>` branches; and
- urgent production fixes use `hotfix/<name>` branches.

Start feature work with:

```sh
git flow feature start <name>
```

Before committing, format and test the complete module:

```sh
gofmt -w .
go test ./...
```
