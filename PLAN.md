# Bello Compiler Implementation Plan

## Scope and completion target
Build a fully working Bello transpiler pipeline and CLI per the grammar/spec in `bello.bnf` and `bello.spec.md`.

## Repository bootstrap
- Initialize Go module `github.com/minions/bello`
- Create package directories:
  - `cmd/bello`
  - `pkg/lexer`
  - `pkg/parser`
  - `pkg/transformer`
  - `pkg/emitter`
  - `pkg/module`
- Add `testdata/` fixture set and sample tests
- Add repository hygiene and tooling ignores via `.gitignore` (`.jj/`, Go artifacts, IDE/temp files)

## Phase 1 — Shared token model (`pkg/lexer/token.go`)
1. Define token type enum and `String()`/`Literal()` helpers
2. Define operator/delimiter tokens
3. Define keyword token list and keyword lookup map
4. Add unit test for keyword token mapping

## Phase 2 — Lexer (`pkg/lexer/lexer.go`)
1. Implement one-pass greedy scanner with position tracking
2. Implement longest-match operators and delimiters
3. Implement semicolon insertion
4. Implement comments with newline preservation
5. Implement numbers:
   - decimal / binary / octal / hex integers
   - decimal and hex floats
   - imaginary suffix
   - underscore separators
6. Implement interpreted and raw strings
7. Expose iterator API matching parser expectation (`NextToken()` etc.)
8. Add tests for:
   - full tokenization of all fixtures
   - semicolon insertion behavior
   - `tatata<-` and `<-tatata`
   - edge strings and escapes

## Phase 3 — AST (`pkg/parser/ast.go`)
1. Add position-aware node base
2. Add nodes for declarations, statements, expressions, and types listed in AGENTS/spec
3. Add helpers for generic lists and pretty-print utility hooks for parser round-trip tests

## Phase 4 — Parser (`pkg/parser/parser.go`)
1. Recursive descent parser over token stream
2. Implement 5-level precedence climb for expressions:
   `LogicalOr -> LogicalAnd -> Comparison -> Addition -> Multiply -> Unary -> Primary`
3. Implement statements/declarations per BNF
4. Disambiguation implementations:
   - composite literal ambiguity in `po/tulaliloo/bee/culo`
   - `bee` -> switch vs type switch via `.(luk)` lookahead
   - `tulaliloo` condition/clause split by semicolon
   - `dala/pwede` first-arg type parse
5. Parse errors formatted as:
   `BEE DOH! <file>:<line>:<col> — <message>`
6. Add recovery: skip to `;`, `}` or next top-level keyword
7. Parser tests from fixtures and key error cases

## Phase 5 — Transformer (`pkg/transformer`)
1. Implement keyword-to-go keyword map
2. Implement predeclared identifier map
3. Implement stdlib package rewrite map
4. Implement stdlib selector rewrite (for Minion imports)
5. Implement `jefe` rewrite rules:
   - package name `jefe` -> `main`
   - function name in package `jefe` -> `main`
6. Walk AST depth-first and emit Go AST nodes
7. Add transformer golden tests using fixture translation expectations

## Phase 6 — Emitter (`pkg/emitter/emitter.go`)
1. Use `go/ast`, `go/token`, `go/format` to render clean Go source
2. Write temp output and return mapping for file positions
3. Add tests to ensure emission compiles at minimum AST-level

## Phase 7 — Module parser (`pkg/module/module.go`)
1. Implement line-oriented parser for `bello.🍑`
2. Parse directives to module model and render `go.mod` text
3. Add tests for require/replace/single-line and grouped forms

## Phase 8 — CLI (`cmd/bello/main.go`)
1. Implement commands:
   - `papala`, `construccion`, `kanpai`, `bonito`, `dame`, `modulo`, `sniff`, `splain`
2. Wire read->lex->parse->transform->emit pipeline
3. Implement temporary output staging and external command runner (`go build/test/run/vet/get`)
4. On tool errors, remap positions via emitted mapping and emit `BEE DOH!` prefix

## Phase 9 — Fixtures and integration
1. Add all required `.🍌` fixtures from AGENTS
2. Add CLI-end-to-end test paths and build checks
3. Round-trip tests: parse -> format -> reparse AST

## Phase 10 — Self-host bootstrap
1. Add deterministic boosta command for seeded Bello compiler boosta.
2. Prefer committed `bootstrap/src` as boosta input when present.
3. Require committed `bootstrap/src` as the boosta source of truth for build operations.
4. Keep repeatable seed refresh tooling documented for intentional updates.
5. Add a self-host validation pass: native bootstrap build compiles Bello from `./cmd/bello` and runs `construccion` to verify it.
6. Add `boosta-run` command that builds a bootstrap compiler and immediately runs a requested Bello subcommand on the source tree.
7. Add `micasa` activation flow: build/install self-hosted compiler into `.bello/bello` for opt-in native replacement.

## Current implementation status

- Lexer and token model: implemented and ready.
- Parser: implemented as Bello token translation + Go parser + AST conversion layer (now builds `File.Decls` and statement/expression/type nodes, including switch/type-switch/select forms).
- Transformer: implemented; stdlib package/method rewrite and `jefe` handling wired.
- Emitter: implemented with `go/format` output.
- Module parser: implemented and supports `modulo`, `bello`, `necesita`, `cambio` (single-line and grouped forms).
- CLI: command pipeline (`papala`, `construccion`, `kanpai`, `sniff`, `bonito`, `dame`, `modulo init`, `splain`) wired.
- CLI module bootstrap now copies local `go.sum` alongside `go.mod` for project and single-file builds.
- Added self-host activation support via `.bello/bello` and `BELLO_SELF_HOST_BIN`.
- Fixtures: full fixture set has been added in `testdata`.
- Example portfolio expanded with `agent`, `slackbot`, `webserver`, `grpc`, `tui`, `crypto`, and `banana_detector`; `stdlib`, `tui`, and `banana_detector` are currently corrected for compile-safe behavior.
- HTTP/3 examples now include local module metadata in `examples/http3` so they build with `bello construccion`.
- Added `snek.🍌` and documented its usage/output in README.
- Docs: `README.md` added for setup, commands, and workflow notes.
- Validation tests: lexer/parser/module/transformer tests pass in local environment.
- Language docs and llms context files are now available (`LANGUAGE.md`, `llms.txt`).

### Open work before delivery
- Finalize `.gitignore` enforcement for all working dirs and wire it into CI/repro checks. *(partially complete; .gitignore exists and is being tracked)*
- Add CLI end-to-end command tests (`bello papala`, `bello construccion`, `bello kanpai`, `bello sniff`, `bello bonito`) via fixture projects. *(added and validated in local run)*
- harden parser recovery and error messages to fully match BNF-level recovery requirements;
- continue refining `bonito` formatting parity (spacing/comments and edge-case conversions) while keeping AST round-trip stable.

### Delivery status
- CLI commands currently verified: `bello papala`, `bello construccion`, `bello kanpai`, `bello sniff`, `bello bonito`, `bello dame`, `bello modulo init`, `bello boosta`, `bello boosta-run`, `bello micasa`, and `bello splain`; examples execute successfully through generated compiler paths.
- All `examples/*.🍌` and nested `examples/http3/*.🍌` build successfully with `bello construccion`.
- `bootstrap/src` is committed as a seed mirror and is used directly by `bello boosta`.
- Seed refresh tooling is documented and available via `go run ./tools/bootstrap_seed.go --source .`.
- GitHub Actions workflow added for:
  - native CI + native bootstrap pass
  - full self-hosted test + example build validation
  - tag-triggered release packaging from `.bello/bello` for Linux/macOS `amd64`/`arm64`
  - release artifact health check (`bello splain`) during packaging

## Completion criteria
- All tests passing locally with fixtures and CLI commands exercised
- `bello construccion` works on all fixture inputs
- `bello bonito` round-trip is stable
- `jefe` translation and Minion stdlib rewriting are verified
