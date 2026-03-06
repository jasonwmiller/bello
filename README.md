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
- `examples/http3/http3_server.🍌`
- `examples/http3/http3_client.🍌`

Language reference: [LANGUAGE.md](/gfs/git/bello/LANGUAGE.md)
LLM context file: [llms.txt](/gfs/git/bello/llms.txt)

## Requirements

- Go toolchain (1.23+), discovered via `PATH` or fallback paths
- Supported platforms can execute commands through Go toolchain wrappers in CLI

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

### REPL

```bash
# start a simple prompt
 go run ./cmd/bello repl
```

Supported commands:
- `/chiku` — show prompt help (minion speak)
- `/bapple` — leave REPL (alias: `return`)

Legacy aliases still work for compatibility:
- `/help` — show prompt help
- `/quit`, `/exit` — leave REPL

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

`/.github/workflows/ci.yml` runs `go test ./...` on pushes and pull requests.

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
