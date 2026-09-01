# Memovee CLI Bootstrap and Tama Integration Orchestrator

Status: planned initial implementation for the standalone `memovee` command.
The directory is intentionally a separate project from the Memovee Phoenix
application.

This document is authoritative for CLI ownership, project bootstrap, command
behavior, state, secret handling, orchestration, packaging, and delivery. The
application-side contracts live in:

```text
kritama/memovee/wip/memovee-cli-bootstrap-contract.md
upmaru/tama/wip/memovee-cli-bootstrap-contract.md
```

The existing Memovee OAuth and Tama `/mcp/app` WIPs remain authoritative for
protocol and application behavior. The CLI coordinates those applications; it
does not replace either application or `tama_oauth`.

## Accepted decisions

- Repository: separate `kritama/memovee-cli` repository.
- Executable name: `memovee`.
- Implementation language: Go.
- Initial platform: Linux amd64 and arm64, including Omarchy installation.
- Additional release targets: macOS amd64/arm64 and other platforms supported
  by the release pipeline when tests are available.
- Architecture: centralized orchestration with decentralized application
  authority and service-owned secret custody.
- Tama source development remains owned by Tama Kit.
- The initial cross-service capability is exactly Tama `/mcp/app` with scope
  `mcp.message`.

The repository name may include `-cli`; human commands, documentation, process
names, packages, and completion scripts use `memovee`.

## Goal

Provide one idempotent control-plane command that can bootstrap Memovee,
prepare Tama, construct their exact trust topology, verify both sides, activate
the integration deliberately, diagnose drift, and roll exposure back safely.

The CLI must support two initial workflows:

1. source development using local Memovee and Tama checkouts; and
2. packaged installation using service-owned environment or secret files.

Omarchy is one distribution channel, not an architectural dependency. The
same binary must work when installed manually or by another package manager.

## Non-goals

The first release does not:

- implement OAuth, JWT verification, refresh rotation, introspection, or MCP;
- move Tama Kit into this repository;
- rewrite Tama Kit's managed `.envrc`;
- create one shared secret file for both services;
- store application secrets in CLI state;
- provision arbitrary cloud infrastructure or mutate a secret manager without
  an explicit output adapter and user request;
- create OAuth clients on behalf of ChatGPT, Codex, or another MCP host;
- forward Tama MCP bearer tokens to Memovee APIs;
- infer allowed browser/MCP origins from service origins;
- silently enable a public resource after preparation; or
- require an interactive terminal for automation.

## System ownership

| Component | Responsibility |
| --- | --- |
| `memovee` CLI | Topology, planning, owner-specific rendering, orchestration, verification, activation, rollback, doctor output |
| Memovee | OAuth policy, Actors, grants, signing-key validation, metadata, JWKS, introspection, token lifecycle |
| Tama | Protected-resource policy, introspection-key validation, JWKS, principal, recipients, MCP server and `message` |
| Tama Kit | Tama source checkout, generated base `.envrc`, pgvector Compose database, Mix setup, test foundation |
| `tama_oauth` | Shared policy-free OAuth/JWT/JWKS mechanics and conformance fixtures |

The CLI consumes versioned non-secret contracts from Memovee and Tama. A
successful plan does not override application startup validation; the running
applications remain the final authority.

## Project bootstrap

Initialize the repository as a Go module using the organization repository
path and keep the runtime dependency set intentionally small.

Proposed layout:

```text
cmd/memovee/main.go
internal/command/
internal/contract/
internal/topology/
internal/plan/
internal/apply/
internal/state/
internal/secret/
internal/safefs/
internal/process/
internal/adapter/memovee/
internal/adapter/tama/
internal/adapter/tamakit/
internal/verify/
internal/output/
internal/version/
contracts/testdata/
integration/
wip/
```

Responsibilities are:

- `command`: parsing, usage, stable exit codes, confirmation policy;
- `contract`: version negotiation and application contract validation;
- `topology`: canonical origins, resources, endpoints, ports, and identifiers;
- `plan`: pure desired-state and diff construction with no writes;
- `apply`: ordered transactional execution and rollback journal;
- `state`: public installation metadata and safe digests only;
- `secret`: key generation and secret-owner routing;
- `safefs`: atomic writes, permissions, ownership, path and symlink safety;
- `process`: allowlisted subprocess execution with redacted environments;
- `adapter/*`: application and Tama Kit boundaries;
- `verify`: bounded metadata, JWKS, introspection, health, and acceptance checks;
- `output`: human and deterministic JSON renderers; and
- `version`: CLI build version and supported contract matrix.

Use Go standard-library primitives wherever they are sufficient. Any CLI,
JOSE/JWK, terminal, or release dependency must be maintained, pinned, reviewed,
and justified by reducing security-sensitive handwritten behavior.

## Command surface

The initial command family is:

```text
memovee setup
memovee dev setup
memovee integration tama plan
memovee integration tama prepare
memovee integration tama verify
memovee integration tama activate
memovee integration tama status
memovee integration tama disable
memovee doctor
memovee version
```

Global automation options:

```text
--json
--no-color
--non-interactive
--yes
--config <path>
```

Mutation commands also support `--dry-run` where meaningful. `--json`, dry-run,
and non-TTY output must be deterministic, ANSI-free, progress-free, and free
of secret values.

### `memovee setup`

Prepare a packaged or current-directory Memovee installation. It discovers or
accepts an explicit installation profile, validates prerequisites, constructs
a plan, and writes only Memovee-owned files.

It must not enable the Tama integration merely because Tama is discovered.

### `memovee dev setup`

Prepare coordinated source checkouts. The initial interface accepts explicit
paths and supports safe sibling discovery:

```text
memovee dev setup \
  --memovee-path /path/to/kritama/memovee \
  --tama-path /path/to/upmaru/tama \
  --memovee-port 4000 \
  --tama-port 4001 \
  --tama-postgres-port 55432
```

The command:

1. validates both repository identities and contract versions;
2. invokes `tama-kit dev setup <path> --port 4001 --postgres-port 55432
   --json` through the Tama Kit adapter;
3. prepares Memovee's existing database and source environment without
   assuming host PostgreSQL when the repository declares Compose ownership;
4. writes separate integration fragments for Memovee and Tama;
5. preserves existing secrets on rerun;
6. leaves both integrations in `prepared` mode unless activation is requested
   by the separate explicit command; and
7. prints exact next commands without printing secret values.

The CLI parses Tama Kit JSON and must not scrape human terminal output. A Tama
Kit version that lacks the required JSON contract or `--port` support fails
with an actionable compatibility error.

### `plan`

Planning is read-only. It validates inputs and contracts, computes canonical
values, inspects current public state, and reports proposed operations.

It must not:

- generate or rotate a private key;
- write a file;
- start a container or service;
- run a migration;
- restart an application; or
- make a remote mutation.

When a new secret will be needed, the plan reports its owner, destination, and
reason, never a placeholder secret value.

### `prepare`

Preparation writes owner-specific configuration and moves both applications to
the non-user-facing `prepared` state. It may start or restart local services
when the selected adapter explicitly owns them.

Preparation is idempotent and preserves valid existing keypairs. It never
rotates a key only because the command was rerun.

### `verify`

Verification performs bounded, read-only protocol and readiness checks:

- application contract compatibility;
- exact origin, issuer, resource, and endpoint agreement;
- Memovee public access-token signing key presence;
- Tama public introspection assertion key presence;
- authenticated inactive-token introspection probe;
- Tama App metadata and route absence while prepared;
- configured allowed-origin syntax;
- root Space recipient availability when recipients were supplied; and
- no resource advertisement before activation.

### `activate`

Activation is a distinct command because it changes externally observable
authorization and resource availability.

It requires a successful current preparation/verification record, checks that
configuration digests have not changed, shows the activation order, and asks
for confirmation in an interactive terminal. Automation must pass `--yes`.

Activation performs:

1. Tama `prepared -> enabled`;
2. Tama restart and protected-resource metadata/challenge verification;
3. Memovee `prepared -> enabled`;
4. Memovee restart and authorization-server resource advertisement
   verification; and
5. optional guided live MCP client acceptance.

If step 3 or 4 fails, the CLI returns Memovee to prepared and Tama to prepared
where the adapter can do so safely. It never weakens authentication to keep the
route available.

### `status` and `doctor`

`status` reports the desired installation record and current public
observations. `doctor` explains drift and failing readiness checks.

Safe output includes:

- component and contract versions;
- paths and service origins;
- lifecycle states recorded by the local plan;
- public key IDs and algorithms;
- endpoint reachability and bounded reason enums;
- resource and issuer agreement;
- recipient availability; and
- whether a restart or activation is required.

Never include raw environment values, private keys, setup tokens, database
passwords, compact assertions, bearer tokens, response bodies containing
credentials, or subprocess environments.

### `disable`

Disable exposure in this order:

1. Memovee `enabled -> prepared` to stop advertisement and new issuance;
2. Tama `enabled -> prepared` to remove metadata, route, and server child;
3. verify both exposure surfaces are unavailable; and
4. retain public overlap keys until all bounded token, assertion, cache, skew,
   and deployment windows have elapsed.

A separate explicit cleanup may later move prepared installations to disabled.
Disable must not delete application data, grants, messages, or secret files.

## Canonical topology model

Represent topology once and project it into both application contracts.

Required logical values:

```text
memovee_origin
tama_origin
tama_mcp_app_resource
memovee_authorization_server_metadata
memovee_jwks_uri
memovee_introspection_endpoint
tama_jwks_uri
tama_introspection_client_id
access_token_signing_algorithm
introspection_signing_algorithm
allowed_origins
recipient_identifiers
```

Derived identifiers:

```text
tama_mcp_app_resource = <tama_origin>/mcp/app
memovee_authorization_server_metadata =
  <memovee_origin>/.well-known/oauth-authorization-server
memovee_jwks_uri = <memovee_origin>/.well-known/jwks.json
memovee_introspection_endpoint = <memovee_origin>/auth/introspections
tama_jwks_uri = <tama_origin>/.well-known/jwks.json
tama_introspection_client_id = <tama_origin>/mcp/app/introspection
```

The local default is Memovee `http://127.0.0.1:4000` and Tama
`http://127.0.0.1:4001`. Production requires exact HTTPS origins.

Do not create trailing-slash aliases, alternate audiences, request-derived
fallbacks, or a generic Memovee API token name.

## Configuration projection

The same topology produces separate owner-specific outputs.

### Memovee-owned output

```text
MEMOVEE_TAMA_MCP_APP_MODE
MEMOVEE_OAUTH_ISSUER
MEMOVEE_TAMA_MCP_APP_RESOURCE
MEMOVEE_OAUTH_SIGNING_ALGORITHM
MEMOVEE_OAUTH_SIGNING_KEY_ID
MEMOVEE_OAUTH_PRIVATE_SIGNING_KEY
MEMOVEE_OAUTH_PUBLIC_SIGNING_KEYS
MEMOVEE_TAMA_INTROSPECTION_CLIENT_ID
MEMOVEE_TAMA_INTROSPECTION_JWKS_URI
```

### Tama-owned output

```text
TAMA_MCP_APP_MODE
TAMA_MCP_APP_RESOURCE
TAMA_MCP_APP_ALLOWED_ORIGINS
TAMA_MCP_APP_AUTHORIZATION_SERVER
TAMA_MCP_APP_JWKS_URI
TAMA_MCP_APP_INTROSPECTION_ENDPOINT
TAMA_MCP_APP_SIGNING_ALGORITHMS
TAMA_MCP_APP_INTROSPECTION_CLIENT_ID
TAMA_MCP_APP_INTROSPECTION_SIGNING_ALGORITHM
TAMA_MCP_APP_INTROSPECTION_SIGNING_KEY_ID
TAMA_MCP_APP_INTROSPECTION_PRIVATE_KEY
TAMA_MCP_APP_INTROSPECTION_PUBLIC_KEYS
```

One exact resource, issuer, client ID, and algorithm selection must agree
across both projections. Each application independently validates its output.

## Secret ownership and key generation

Never generate a shared symmetric key for this integration.

The CLI generates or imports two distinct asymmetric keypairs:

| Keypair | Private-key destination | Public-key publication |
| --- | --- | --- |
| Memovee access-token signing | Memovee-owned secret output | Memovee JWKS |
| Tama introspection assertion signing | Tama-owned secret output | Tama JWKS |

Key generation uses the operating system cryptographic random source and
initially produces RSA keys suitable for `RS256`. Key IDs must be unique,
stable across idempotent reruns, free of secret material, and include enough
versioning information to support rotation.

The private keys may exist in CLI process memory only while being delivered to
their final owner-specific destinations. They must not be written to a common
temporary file, CLI state, logs, crash reports, command arguments, or JSON
output.

Existing valid keys are preserved. Rotation is a later explicit operation with
publish-new, verify-new, switch-signer, overlap, and remove-old phases.

## State model

Store only non-secret reconciliation state under the platform's XDG state
directory or an explicitly selected state path.

An installation record contains:

- installation ID and schema version;
- CLI version;
- Memovee, Tama, Tama Kit, and application-contract versions;
- selected profile and adapter types;
- canonical public topology;
- owner-specific destination paths;
- public key IDs, algorithms, and public-key digests;
- desired lifecycle states;
- safe file digests and ownership markers;
- last successful prepare, verify, activate, and disable checkpoints; and
- restart requirements and safe error enums.

It must not contain private JWK parameters, environment-file content, database
credentials, setup tokens, access/refresh tokens, client assertions, browser
sessions, or provisioning credentials.

State writes are atomic and mode `0600` even though the state is non-secret.
Corrupt or unknown future schemas fail closed with a recovery explanation.

## Managed-file safety

Every write follows these invariants:

- plan root-anchored ignore rules before writing a repository-local secret;
- reject a secret file already tracked by Git;
- resolve the expected project root and reject path escape;
- inspect every existing ancestor with `lstat` and reject symbolic-link or
  non-directory ancestors;
- reject a symbolic-link destination;
- preserve owner and mode for safe existing managed files;
- use same-filesystem temporary files, `fsync` as appropriate, and atomic
  rename;
- write private material with mode `0600` from creation;
- preserve valid existing secrets and detect user drift;
- attach an ownership marker and digest without embedding a secret digest in
  user-visible output; and
- roll back all completed writes if post-write validation fails.

The CLI must not take ownership of an unmarked user file merely because its
contents resemble a generated template.

## Single-writer file contract

For source development:

| File | Writer |
| --- | --- |
| Memovee committed `.envrc` | Memovee repository |
| Memovee `.memovee.integration.env` | Memovee CLI |
| Tama `.envrc` | Tama Kit |
| Tama `.tama.dev.postgres.env` | Tama Kit |
| Tama `.tama.memovee.integration.env` | Memovee CLI |
| CLI installation state | Memovee CLI |

The committed/generated base env files load the optional integration fragments
by fixed root-relative path. This avoids concurrent ownership and allows each
tool to upgrade its files independently.

Packaged adapters render equivalent owner-specific files or secret payloads.
They never combine both services' private keys.

## Application-contract negotiation

Each application commits a non-secret versioned bootstrap contract. The CLI
has a tested compatibility matrix rather than accepting unknown structures.

Negotiation rules:

- exact major contract versions must match;
- the CLI may support multiple known minor versions explicitly;
- an unknown required field, state, algorithm, or constraint fails planning;
- application versions outside the declared compatible range fail with the
  required upgrade direction;
- the Memovee and Tama compatibility identifiers must agree;
- contract fixtures are copied only into CLI testdata with source version and
  digest, not edited as a third authority; and
- live startup and public endpoint verification remain required after a
  contract-valid render.

## Adapter contracts

### Memovee source adapter

- Verify the Mix application and repository identity.
- Read the committed bootstrap contract.
- Plan the owner-specific integration fragment.
- Use the repository's declared Compose and Mix setup commands.
- Load environment without echoing values.
- Run only explicit allowlisted commands.
- Report structured safe status.

### Tama Kit adapter

- Discover `tama-kit` and inspect its version.
- Require the supported `dev setup --json` contract.
- Pass paths and ports as separate arguments, never secrets.
- Parse exactly one JSON result document.
- Treat stderr as diagnostic text and redact unexpected environment values.
- Do not import Tama Kit's JavaScript modules or scrape human output.
- Do not duplicate its PostgreSQL, Terraform, or managed `.envrc` logic.

### Packaged service adapter

- Render to an explicit destination selected by the user or package profile.
- Preserve service ownership and permissions.
- Separate rendering from service restart.
- Show the exact restart action in the plan.
- Never assume root authority; request it only for a known system-owned target.

### HTTP verifier

- Use explicit connect, receive, and overall deadlines.
- Reject redirects unless the exact contract permits one.
- Bound response bytes before decoding.
- Require expected content types and JSON shapes.
- Bind JWKS fetches to configured origins.
- Never include bearer tokens or assertions in errors, traces, or cache keys.
- Report safe stage and reason enums.

## Orchestration state machine

The coordinated states are:

```text
unconfigured
    -> prepared
    -> verified
    -> enabled
    -> accepted

enabled
    -> prepared (disable exposure)
    -> disabled (later cleanup after overlap windows)
```

Transitions are resumable. The CLI records a checkpoint only after the
transition and its verification pass. A failed step does not advance desired
state.

### Prepare sequence

1. inspect contracts, versions, paths, ports, current files, and services;
2. construct a pure plan;
3. generate only missing owner-specific keys in memory;
4. apply ignore and environment operations transactionally;
5. set Memovee to prepared and start/restart it;
6. set Tama to prepared and start/restart it;
7. verify file permissions and safe public state; and
8. record the prepared checkpoint without secrets.

### Verify sequence

1. verify Memovee issuer and metadata endpoints;
2. verify Memovee JWKS contains its expected public signing key;
3. verify Tama JWKS contains its expected introspection public key;
4. verify Tama App metadata and route remain unavailable;
5. perform a fresh Tama client assertion against Memovee introspection using a
   deliberately invalid bounded token value;
6. require an authenticated inactive response;
7. verify exact public topology agreement; and
8. record public observations and digests.

### Activate sequence

1. reject stale verification or changed configuration digests;
2. verify recipients and explicit allowed origins;
3. enable and restart Tama;
4. verify App metadata, challenge, and route behavior;
5. enable and restart Memovee;
6. verify resource advertisement;
7. guide or run the authorized live acceptance flow; and
8. record acceptance without retaining credentials.

## Error and exit-code contract

Define stable categories and nonzero exit codes for:

```text
usage
prerequisite
contract
configuration
ownership
secret
process
network
verification
activation
rollback
internal
```

JSON output uses one versioned envelope for success and failure. Human output
may add color and stage progress, but both renderers expose the same safe
semantic result.

Errors must identify what the user can do next without revealing sensitive
values. A subprocess failure reports its allowlisted command name and exit
status, not its entire environment.

## Human and automation output

Interactive output should provide:

- clear stage-based progress;
- a concise plan before mutation;
- colored success/warning/error states when supported;
- copyable commands and URLs on single logical lines;
- an explicit activation confirmation; and
- actionable recovery or resume instructions.

Automation output must be deterministic and non-interactive. JSON and non-TTY
modes never show progress animation, ANSI escapes, token-bearing setup URLs,
private keys, secret values, or prompts.

## Packaging and installation

Release the CLI independently from the Phoenix application.

The release pipeline must:

- build reproducible versioned binaries for supported OS/architecture pairs;
- run unit, race, integration, and contract tests before release;
- publish SHA-256 checksums and provenance/SBOM artifacts;
- sign release artifacts using the organization release policy;
- embed version, revision, build time, and supported contract versions;
- generate shell completions and concise install documentation; and
- verify a clean installation in an isolated environment.

Omarchy packaging installs the released `memovee` binary and completion files.
It must not clone application repositories, create integration secrets, or
activate services during package installation. Users run `memovee setup` after
installation.

Other installation paths may include verified release archives and package
recipes. The CLI must not require Go on the target machine.

## Test strategy

### Unit tests

- topology derivation and exact URI validation;
- contract negotiation and unknown-version rejection;
- plan purity and deterministic diffs;
- secret-owner routing and public JWK derivation;
- state serialization without forbidden fields;
- stable exit codes and JSON envelopes;
- redaction and hostile error values; and
- lifecycle transition guards.

### Filesystem security tests

- mode `0600` from initial creation;
- tracked-secret rejection;
- root-anchored ignore ordering after negations;
- symbolic-link destination and ancestor rejection;
- non-directory ancestor rejection;
- path traversal and project-root escape;
- user drift and unknown ownership;
- atomic replacement and rollback;
- existing-secret preservation; and
- external canaries remain untouched.

### Process tests

- exact allowlisted executable and arguments;
- paths containing spaces and Unicode;
- no shell interpolation;
- secret-free command line and diagnostic output;
- Tama Kit JSON success, failure, malformed output, and version mismatch; and
- cancellation and timeout cleanup.

### Protocol verification tests

- metadata and JWKS success;
- origin, resource, client-ID, algorithm, and `kid` mismatch;
- redirect, timeout, DNS/connect failure, content type, malformed JSON, and
  oversized response;
- authenticated inactive-token introspection success;
- assertion unknown `kid`, invalid signature, expiry, and replay;
- key rotation overlap and unknown-`kid` refresh; and
- no raw assertion or token in any result.

### Cross-repository tests

- source setup with Memovee 4000 and Tama 4001;
- independent database setup through repository-owned mechanisms;
- fresh `unconfigured -> prepared -> verified`;
- idempotent rerun with unchanged secrets;
- `verified -> enabled -> accepted`;
- refresh across at least two access-token expiries;
- revocation and Actor deactivation fail closed;
- `/mcp/system` and `/mcp/app` token isolation;
- `enabled -> prepared` rollback; and
- upgrade from each supported application-contract minor version.

Run `go test ./...`, the race detector on supported platforms, static analysis,
format checks, vulnerability/dependency review, integration tests, and
`git diff --check` before release.

## Implementation phases

### Phase 0: repository foundation

- Initialize Git, Go module, license, README, contribution instructions, and
  supported toolchain.
- Add command entrypoint, version output, structured errors, exit codes,
  human/JSON output, and CI.
- Add release builds and checksum verification without publishing a stable
  release yet.

### Phase 1: contracts and pure planning

- Implement topology and application-contract models.
- Add version negotiation and the supported compatibility matrix.
- Implement `integration tama plan`, `status`, and initial `doctor` as
  read-only commands.
- Add golden JSON and redaction tests.

### Phase 2: safe owner-specific rendering

- Implement safe filesystem operations, state, key generation, and separate
  Memovee/Tama projections.
- Implement prepared-mode rendering and transactional rollback.
- Add all filesystem and secret-safety tests before invoking services.

### Phase 3: source-development adapters

- Implement the Memovee source adapter.
- Implement the Tama Kit subprocess/JSON adapter.
- Support Memovee 4000, Tama 4001, and explicit port overrides.
- Add `memovee dev setup`, prepare, and verify against live local checkouts.

### Phase 4: activation and rollback

- Implement guarded activation, restart adapters, public verification, and
  disable-to-prepared rollback.
- Require current verification digests and explicit confirmation.
- Complete cross-service identity, refresh, revocation, and deactivation tests.

### Phase 5: packaged installation and Omarchy delivery

- Add packaged environment-file/service adapters.
- Publish signed multi-architecture release artifacts and checksums.
- Add the Omarchy package recipe as a thin binary installer.
- Validate install, setup, upgrade, doctor, disable, and uninstall boundaries
  in a clean environment.

### Phase 6: explicit rotation and additional adapters

- Add access-token and introspection key-rotation workflows with overlap.
- Add external secret-store outputs only when their authorization and rollback
  contracts are explicit.
- Add other operating-system package recipes without changing core behavior.

## Cross-repository delivery order

1. Memovee and Tama agree on contract version 1, lifecycle names, topology,
   variable ownership, and readiness semantics.
2. Memovee implements disabled/prepared/enabled and publishes its contract.
3. Tama implements disabled/prepared/enabled, closes the public-key staging
   gap, and publishes its contract.
4. Tama Kit adds explicit Tama HTTP-port selection and the optional
   integration-fragment include while preserving its single-writer ownership.
5. Memovee CLI implements contract negotiation and read-only planning.
6. Memovee CLI implements owner-specific preparation and live verification.
7. All three repositories pass the prepared-state integration suite.
8. Memovee CLI implements guarded activation and rollback.
9. Live MCP acceptance passes before the first stable CLI release or Omarchy
   package is advertised.

Do not merge an application into a state that advertises the resource before
its counterpart and the compatible CLI path are ready for coordinated rollout.

## Acceptance criteria

The initial CLI is complete when:

- the binary installs independently and reports its supported contracts;
- source setup delegates Tama-specific work to Tama Kit through JSON;
- packaged setup does not require Node, Go, Elixir, or Erlang at CLI runtime;
- one canonical topology produces exact, agreeing application projections;
- Memovee and Tama private keys remain in separate owner-specific files;
- plans are read-only and never generate secret material;
- preparation is transactional, idempotent, resumable, and secret-preserving;
- every managed path is protected from tracking, drift, traversal, and
  symlinked ancestors;
- prepared mode verifies both public keys and authenticated introspection while
  user authorization and `/mcp/app` remain unavailable;
- activation requires current verification and explicit authority;
- rollback removes advertisement and route exposure before removing trust
  material;
- doctor output identifies drift without exposing credentials;
- JSON/non-TTY output is deterministic and secret-free;
- real MCP acceptance verifies Actor identity, refresh, revocation, and Actor
  deactivation;
- `/mcp/system` remains isolated and unchanged; and
- CLI, Memovee, Tama, Tama Kit, and cross-service verification suites are
  green.

## Later decisions

The following remain explicit later design work rather than implicit first
release behavior:

- production service manager adapters beyond the first packaged profile;
- remote host orchestration;
- cloud secret-manager integrations;
- unattended key rotation policy;
- automatic root Space graph provisioning beyond current Tama Kit contracts;
- Windows support and installation packaging; and
- a plugin API for third-party deployment adapters.
