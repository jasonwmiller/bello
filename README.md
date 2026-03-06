# Bello

Bello is a source-to-source transpiler for the Bello language (Minion-speak syntax) targeting Go.

## Quick demo

```banana
kampung jefe

banana jefe() {
	poopaye("i love this language")
}
```

```bash
go run ./cmd/bello papala testdata/hello.🍌
# output
bello
```

It takes `.🍌` source files, parses them through a Go-backed flow, transforms syntax/features to Go equivalents, emits Go source, and then runs Go tooling (`go run`, `go build`, `go test`, `go vet`, or `go get`).

## Repository Layout

- `cmd/bello` - CLI entrypoint
- `pkg/lexer` - tokenizer and token model
- `pkg/parser` - source translation + Go parse + Bello-mirrored AST bridge
- `pkg/transformer` - keyword/builtin/stdlib and import rewriting
- `pkg/emitter` - Go emission (via `go/format`)
- `pkg/module` - parser for `bello.🍑` module descriptors
- `testdata` - reference `.🍌` fixtures

### Example programs

- `examples/hello.🍌`
- `examples/calculator.🍌`
- `examples/loops.🍌`
- `examples/structs.🍌`
- `examples/stdlib.🍌`
- `examples/agent.🍌`
- `examples/slackbot.🍌`
- `examples/webserver.🍌`
- `examples/grpc.🍌`
- `examples/tui.🍌`
- `examples/crypto.🍌`
- `examples/banana_detector.🍌`
- `examples/file_watcher.🍌`
- `examples/cancellation.🍌`
- `examples/cache.🍌`
- `examples/pipeline.🍌`
- `examples/http_json.🍌`
- `examples/generic_stack.🍌`
- `examples/error_wrapping.🍌`
- `examples/minion_postal.🍌`
- `examples/minion_heartbeat.🍌`
- `examples/minion_gateway.🍌`
- `examples/minion_notebook.🍌`
- `examples/minion_guardian.🍌`
- `examples/minion_mischief.🍌`
- `examples/minion_vibes.🍌`
- `examples/snek.🍌`
- `examples/http3/http3_server.🍌`
- `examples/http3/http3_client.🍌`

Language reference: [LANGUAGE.md](/gfs/git/bello/LANGUAGE.md)
LLM context file: [llms.txt](/gfs/git/bello/llms.txt)

## Requirements

- Go toolchain (1.23+), discovered via `PATH` or fallback paths
- Supported platforms can execute commands through Go toolchain wrappers in CLI

## Install and shell integration

- Installed mode (recommended):

```bash
go install ./cmd/bello

# verify installed command
bello splain
```

Once installed, `bello` directly runs the native translator and Go tooling (`go run`, `go build`, `go test`, `go vet`, `go get`) for each command.

- Source mode (repo checkout):

```bash
# quick local entrypoint (temporary)
alias bello='go run /absolute/path/to/repo/cmd/bello'
```

Enable shell autocomplete:

```bash
# bash
source <(bello completion)

# zsh
source <(bello completion zsh)

# fish
bello completion fish | source
```

## Building and Running

### Run a single file

```bash
# translate + run
 go run ./cmd/bello papala path/to/file.🍌
```

### HTTP/3 one-shot demo

```bash
# Terminal 1: start banana server (runs forever)
go run ./cmd/bello papala examples/http3/http3_server.🍌 > /tmp/bello_http3_server.log 2>&1 & \
SERVER_PID=$!

# Terminal 2: run banana client
sleep 2
go run ./cmd/bello construccion examples/http3/http3_client.🍌

kill "$SERVER_PID"
```

### Build project

```bash
# build all Bello files in current module/directories
 go run ./cmd/bello construccion
# or explicit path
 go run ./cmd/bello construccion ./some/dir
```

### Bootstrap/self-host check

```bash
# uses a checked-in minion seed in bootstrap/src (or generate temporarily),
# builds a native bootstrap compiler, validates it by running
# `bello construccion` on the same tree, and can launch additional commands.
go run ./cmd/bello bootstrap .
```

### Bootstrap then execute through generated compiler

```bash
# bootstrap with native compiler and immediately run an example
go run ./cmd/bello bootstrap-run . papala examples/hello.🍌

# bootstrap and build the same project
go run ./cmd/bello bootstrap-run . construccion .
```

### Promote bootstrapped compiler to active compiler

```bash
go run ./cmd/bello selfhost .

# make self-hosted compiler active for current shell
export BELLO_SELF_HOST_BIN=$PWD/.bello/bello

# then run like normal
./.bello/bello papala examples/hello.🍌
```

`bello bootstrap` (or `bello boosta`) is the bootstrap lane for the next phase:
- use `bootstrap/src` (preferred) as a prebuilt minion seed, or generate one on the fly,
- build `cmd/bello` with the current native translator,
- run the newly built compiler through `construccion` on the same seed tree.

Seed layout:

- `bootstrap/src/go.mod`
- `bootstrap/src/cmd/bello/main.🍌`
- `bootstrap/src/pkg/...` translator packages

To regenerate seed files manually:

```bash
go run ./tools/bootstrap_seed.go --source .
```

### REPL

```bash
# start a simple prompt
 go run ./cmd/bello repl
# or minion way
 go run ./cmd/bello chiku
```

Supported commands:
- `/chiku` — show prompt help (minion speak)
- `/bapple` — leave REPL (alias: `return`)

Legacy aliases still work for compatibility:
- `/help` — show prompt help
- `/quit`, `/exit` — leave REPL

## Example: snek script

```bash
go run ./cmd/bello papala examples/snek.🍌 foo BAR baz
```

Expected output:
```text
snek script woke up with 4 arguments
slice 1 : foo
slice 2 : bar
slice 3 : baz
```

 ### Test package

```bash
 go run ./cmd/bello kanpai
```

### Vet project

```bash
 go run ./cmd/bello sniff
```

### Format output (Bello source pretty print pass)

```bash
 go run ./cmd/bello bonito path/to/file.🍌
```

### Module and support commands

```bash
# write go dependency package file
 go run ./cmd/bello dame github.com/some/pkg

# create module file
 go run ./cmd/bello modulo init module/name

# show short help text
 go run ./cmd/bello splain

# shell completion
 go run ./cmd/bello completion [bash|zsh|fish]

# run bootstrap validation lane
 go run ./cmd/bello bootstrap [dir]
 go run ./cmd/bello boosta [dir]
 go run ./cmd/bello bootstrap-run [dir] <command> [args...]
 go run ./cmd/bello boosta-run [dir] <command> [args...]
 go run ./cmd/bello selfhost [dir]
```

## Behavior and mappings

Bello syntax is translated with a keyword/predeclared mapping layer in the translator and then parsed as Go.

- Program/command keywords are normalized to Go (`kampung -> package`, `banana -> func`, etc.).
- Types and predeclared identifiers map to Go (`me -> int`, `bababa -> string`, etc.).
- Minion stdlib package/method rewrites are applied in transformer stage.

## Tests

```bash
/usr/local/go/bin/go test ./...
```

## Supported CI

`/.github/workflows/ci.yml` runs:
- `go test ./...`
- `go run ./cmd/bello bootstrap .`
- `go run ./cmd/bello selfhost .`
- Self-hosted verification with `.bello/bello`:
  - `kanpai examples/hello.🍌`
  - `bootstrap .`
  - compile checks for every `examples/**/*.🍌` file

> Note: when `BELLO_SELF_HOST_BIN` points to `.bello/bello` (or it is found by walking parent directories), `bello` automatically routes regular commands through the self-hosted compiler. This keeps day-to-day usage non-interactive and fully automated.

## Release

Pushing tags matching `v*` runs release packaging from the self-hosted compiler:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Artifacts:
- `dist/bello-<tag>-linux-amd64.tar.gz`
- `dist/bello-<tag>-linux-amd64.sha256`
- `dist/bello-<tag>-linux-arm64.tar.gz`
- `dist/bello-<tag>-linux-arm64.sha256`
- `dist/bello-<tag>-darwin-amd64.tar.gz`
- `dist/bello-<tag>-darwin-amd64.sha256`
- `dist/bello-<tag>-darwin-arm64.tar.gz`
- `dist/bello-<tag>-darwin-arm64.sha256`

Linux `amd64` is the x86_64 release artifact and `arm64` targets ARM-based runners.

## Development plan status

- `PLAN.md` tracks implementation status and open work.
- Current implementation includes lexer, parser bridge, transformer, emitter, module parser, fixtures, CLI command wiring, and regression checks.

### Current capabilities and remaining gaps

- ✅ Lexer, parser bridge, transformer, emitter, module parser, and CLI command dispatch are implemented.
- ✅ Non-interactive compile/run flows are working via `papala`, `construccion`, `kanpai`, `sniff`, and `dame`.
- ✅ Example suite includes runtime-oriented programs and build-only programs (HTTP server, grpc, async patterns, stdlib cases, etc.).
- ⚠️ `construccion` on a mixed directory of standalone examples may still fail if multiple `main` programs coexist in that directory.
- ⚠️ Parser behavior is intentionally conservative around ambiguous constructs (for example, composite-vs-block edge cases), so some exotic syntax may still need follow-up tests before claiming broad spec parity.

### Active conventions

- Tooling errors are surfaced as compiler-style diagnostics with `BEE DOH! ...` format.
- `jefe` is rewritten to `main` in package/function positions as required by translator rules.

## Notes

- Use `GO_BIN` environment variable to force a specific Go executable if needed.
