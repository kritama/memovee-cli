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
go test ./...
```

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the Git Flow workflow.

## License

Licensed under the Apache License, Version 2.0. See [`LICENSE`](LICENSE).
