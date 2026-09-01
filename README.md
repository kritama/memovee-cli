# Memovee CLI

`memovee` is the standalone control-plane command for bootstrapping Memovee,
preparing Tama, verifying their trust topology, and deliberately activating or
disabling the Tama `/mcp/app` integration.

The project is in its initial implementation phase. The authoritative design
is in [`wip/bootstrap-orchestration.md`](wip/bootstrap-orchestration.md).

## Development

The project targets Go 1.25. Install the pinned toolchain with
[`mise`](https://mise.jdx.dev/), then run:

```sh
mise install
go fmt ./...
go test ./...
```

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the Git Flow workflow.

## Commands

The current repository foundation provides human and deterministic JSON
version output:

```sh
go run ./cmd/memovee version
go run ./cmd/memovee --json version
```

Global options must appear before the command. Run `memovee help` for the
current command surface.

## Exit codes

Exit codes are stable API values for automation:

| Code | Category |
| ---: | --- |
| 0 | success |
| 2 | usage |
| 10 | prerequisite |
| 11 | contract |
| 12 | configuration |
| 13 | ownership |
| 14 | secret |
| 15 | process |
| 16 | network |
| 17 | verification |
| 18 | activation |
| 19 | rollback |
| 20 | internal |

## License

Licensed under the Apache License, Version 2.0. See [`LICENSE`](LICENSE).
