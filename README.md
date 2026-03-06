# Bello

Bello is a source-to-source transpiler for the Bello language (Minion-speak syntax) targeting Go.

## Quick demo

```banana
kampung jefe

muak "boca"

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

## Requirements

- Go toolchain (1.23+), discovered via `PATH` or fallback paths
- Supported platforms can execute commands through Go toolchain wrappers in CLI

## Building and Running

### Run a single file

```bash
# translate + run
 go run ./cmd/bello papala path/to/file.🍌
```

### Build project

```bash
# build all Bello files in current module/directories
 go run ./cmd/bello construccion
# or explicit path
 go run ./cmd/bello construccion ./some/dir
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

### Active conventions

- Tooling errors are surfaced as compiler-style diagnostics with `BEE DOH! ...` format.
- `jefe` is rewritten to `main` in package/function positions as required by translator rules.

## Notes

- Use `GO_BIN` environment variable to force a specific Go executable if needed.
