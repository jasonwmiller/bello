# Bello — Implementation Agent Guide

## Project Overview

Bello is a compiled language with Go semantics and Minion-speak syntax. The compiler is a **source-to-source transpiler**: Bello source (`.🍌`) is lexed, parsed into an AST, transformed to Go AST, emitted as Go source, then compiled with `go build`.

### Reference Files

- `bello.spec.md` — full language specification (keywords, types, stdlib mappings, conventions)
- `bello.bnf` — implementation-ready EBNF grammar with precedence levels and disambiguation rules

Read both files before starting any implementation work. The BNF is the parser's source of truth. The spec is the source of truth for keyword mappings, stdlib mappings, and language semantics.

---

## Architecture

```
 .🍌 source
    |
    v
 [Lexer] ──> token stream
    |
    v
 [Parser] ──> Bello AST
    |
    v
 [Transformer] ──> Go AST
    |
    v
 [Emitter] ──> .go source
    |
    v
 [go build] ──> binary
```

### Directory Structure

```
bello/
  cmd/
    bello/          — CLI entry point
      main.go
  pkg/
    lexer/          — tokenizer
      lexer.go
      lexer_test.go
      token.go      — token types and keyword table
    parser/         — recursive descent parser
      parser.go
      parser_test.go
      ast.go        — AST node types
    transformer/    — Bello AST -> Go AST
      transformer.go
      transformer_test.go
      keywords.go   — Bello->Go keyword mapping table
      stdlib.go     — Bello->Go stdlib package/method mapping
    emitter/        — Go AST -> Go source text
      emitter.go
      emitter_test.go
    module/         — bello.🍑 file parser
      module.go
      module_test.go
  testdata/         — .🍌 test fixtures
  bootstrap/        — committed Bello bootstrap seed source mirror
    go.mod
    cmd/
      bello/        — bootstrap compiler entry
        main.🍌
    pkg/           — translator packages mirrored in Bello
      lexer/
      parser/
      emitter/
      transformer/
      module/
  tools/           — repository utility tools
    bootstrap_seed.go
  examples/        — runnable Bello example programs
  .github/         — CI/release automation
  bello.spec.md
  bello.bnf
  go.mod
```

---

## Implementation Order

Build and test each phase completely before moving to the next. Each phase has a clear input/output contract.

### Phase 1: Token Types (`pkg/lexer/token.go`)

Define all token types. This is the shared vocabulary between lexer and parser.

**Token categories:**
- **Keywords** (25): `kampung`, `muak`, `banana`, `bapple`, `pooka`, `gelato`, `luk`, `kampai`, `buddies`, `papoy`, `po`, `ka`, `tulaliloo`, `tikali`, `buttom`, `bajo`, `bee`, `doh`, `meh`, `underpa`, `tatata`, `culo`, `tank_yu`, `patalaki`, `waaah`, `dala`, `pwede`
- **Literals**: INT, FLOAT, IMAGINARY, RUNE, STRING
- **Identifiers**: IDENT (covers predeclared types/functions/constants too)
- **Operators** (see BNF section 9): `+`, `-`, `*`, `/`, `%`, `&`, `|`, `^`, `<<`, `>>`, `&^`, `&&`, `||`, `<-`, `++`, `--`, `==`, `!=`, `<`, `<=`, `>`, `>=`, `=`, `:=`, `+=`, `-=`, `*=`, `/=`, `%=`, `&=`, `|=`, `^=`, `<<=`, `>>=`, `&^=`, `!`
- **Delimiters**: `(`, `)`, `[`, `]`, `{`, `}`, `,`, `.`, `;`, `:`, `...`
- **Special**: EOF, ILLEGAL, COMMENT

Build a keyword lookup map: `string -> TokenType`. The lexer checks every IDENT against this map.

**Test**: unit test that every keyword string maps to the correct token type.

### Phase 2: Lexer (`pkg/lexer/lexer.go`)

Implement a greedy, single-pass lexer. Key behaviors from the BNF:

1. **Greedy matching** — longest token wins. `<-` is one token, not `<` then `-`. `&^` is one token, not `&` then `^`.
2. **Semicolon insertion** — after each newline, if the previous token was one of these, insert a `;` token:
   - Any IDENT
   - Any literal (INT, FLOAT, IMAGINARY, RUNE, STRING)
   - Keywords: `buttom`, `bajo`, `bapple`, `patalaki`
   - Punctuation: `++`, `--`, `)`, `}`, `]`
3. **Comments** — `//` to newline, `/* */` (non-nesting). Treat as whitespace but preserve newline significance for semicolons.
4. **String literals** — interpreted (`"..."` with escapes) and raw (`` `...` `` with no escapes, can span lines).
5. **Number literals** — decimal, binary (`0b`/`0B`), octal (`0o`/`0O`), hex (`0x`/`0X`), floats (decimal and hex with `p`/`P` exponent), imaginary (`i` suffix). Underscores allowed as separators.

**Test**: lex complete Bello programs from `testdata/`. Verify token streams including auto-inserted semicolons. Test edge cases: `tatata<-` (keyword + operator), `<-tatata` (operator + keyword), emoji in string literals.

### Phase 3: AST (`pkg/parser/ast.go`)

Define AST node types mirroring the BNF productions. Every named production in the BNF is a node type or is inlined into a parent node.

Key node types:
- `File` (SourceFile) — package name, imports, declarations
- `ImportSpec` — path, alias
- `VarDecl`, `ConstDecl`, `TypeDecl`
- `FuncDecl` — receiver, name, type params, signature, body
- `FieldList`, `Field` — for params, results, struct fields
- **Statements**: `BlockStmt`, `ReturnStmt`, `IfStmt`, `ForStmt`, `SwitchStmt`, `TypeSwitchStmt`, `SelectStmt`, `GoStmt`, `DeferStmt`, `BranchStmt` (break/continue/goto/fallthrough), `LabeledStmt`, `AssignStmt`, `SendStmt`, `IncDecStmt`, `ExprStmt`
- **Expressions**: `BinaryExpr` (with operator token), `UnaryExpr`, `CallExpr`, `IndexExpr`, `SliceExpr`, `SelectorExpr`, `TypeAssertExpr`, `Ident`, `BasicLit`, `CompositeLit`, `FuncLit`
- **Types**: `ArrayType`, `SliceType`, `MapType`, `ChanType`, `PointerType`, `FuncType`, `StructType`, `InterfaceType`

Every node stores source position (file, line, column) for error reporting.

### Phase 4: Parser (`pkg/parser/parser.go`)

Recursive descent parser consuming the token stream. Follow the BNF exactly.

**Critical disambiguation rules** (from BNF):
1. **Operator precedence** — use the 5-level precedence climbing grammar: `Expression -> LogicalOrExpr -> LogicalAndExpr -> ComparisonExpr -> AdditionExpr -> MultiplyExpr -> UnaryExpr -> PrimaryExpr`. Do NOT use Pratt parsing unless you map it to these exact levels.
2. **Composite literal ambiguity** — in `po`, `tulaliloo`, `bee`, `culo` bodies: `{` after expression is ALWAYS a block, never a composite literal. Composite literals in these contexts must be parenthesized.
3. **SwitchStmt vs TypeSwitchStmt** — both start with `bee`. Look ahead for `.(luk)` pattern before `{`. If found, parse as TypeSwitchStmt. Otherwise, SwitchStmt.
4. **ForCondition vs ForClause** — if a semicolon follows the first expression/statement, it's a ForClause. No semicolons = ForCondition or infinite loop.
5. **`dala`/`pwede` calls** — first argument is a Type, not an Expression. When parsing `Arguments` after these keywords, attempt Type parse first before falling back to Expression.

**Error recovery**: on parse error, emit a `BEE DOH!` message in the format from spec section 11:
```
BEE DOH! <file>:<line>:<col> — <message>
```
Skip to the next synchronization point (`;`, `}`, or next top-level keyword).

**Test**: parse every example from the spec. Round-trip: parse then pretty-print AST and verify structure.

### Phase 5: Transformer (`pkg/transformer/`)

Convert Bello AST to Go AST. This is primarily a mapping operation.

**`keywords.go`** — a map from Bello keyword/builtin to Go equivalent:
```
kampung -> package     banana -> func       bapple -> return
pooka   -> var         gelato -> const      luk    -> type
kampai  -> struct      buddies -> interface  papoy  -> map
po      -> if          ka     -> else       tulaliloo -> for
tikali  -> range       buttom -> break      bajo   -> continue
bee     -> switch      doh    -> case       meh    -> default
underpa -> go          tatata -> chan        culo   -> select
tank_yu -> defer       patalaki -> fallthrough      waaah -> goto
dala    -> make        pwede  -> new
```

Predeclared identifiers:
```
me -> int    me8 -> int8    me16 -> int16   me32 -> int32   me64 -> int64
ti -> uint   ti8 -> uint8   ti16 -> uint16  ti32 -> uint32  ti64 -> uint64
la32 -> float32   la64 -> float64   butt -> bool   bababa -> string
todo -> any   whaaat -> error
si -> true   naga -> false   hana -> nil   mamamia -> iota
baboi -> append   para_tu -> len   stupa -> cap   cierro -> close
yeet -> delete   mimik -> copy   BEE_DOH -> panic   gelatin -> recover
poopaye -> println
```

**`stdlib.go`** — maps Minion package imports and method calls to Go equivalents. Two cases:

1. **Package import rewriting**: `"boca"` -> `"fmt"`, `"casa"` -> `"os"`, etc. (full table in spec section 7.2)
2. **Method name rewriting**: when the import is a Minion package, rewrite method calls. `boca.poopaye(...)` -> `fmt.Println(...)`, `boca.blabla(...)` -> `fmt.Printf(...)`, etc. (full tables in spec sections 7.3-7.9)

Go-mode imports (`"fmt"`, `"os"`, etc.) pass through unchanged — no method rewriting.

**Transformer walk**:
1. Walk the AST depth-first
2. Rewrite every Bello identifier node using the keyword/type/builtin maps
3. Rewrite import paths using the stdlib package map
4. For Minion-mode imports, rewrite selector expressions (`boca.poopaye` -> `fmt.Println`)
5. `jefe` as a package name -> `main`, `jefe` as a function name -> `main`

**Test**: transform Bello ASTs from the spec examples and verify the resulting Go AST matches expected Go code.

### Phase 6: Emitter (`pkg/emitter/emitter.go`)

Print the Go AST as valid, `gofmt`-formatted Go source. Use `go/printer` or `go/format` from the Go stdlib if using `go/ast` nodes, or write a simple recursive printer if using custom AST nodes.

The emitter should:
1. Write the Go source to a temp directory
2. Preserve a mapping from Go source positions back to Bello source positions (for error remapping)

**Test**: emit Go source from transformed ASTs, run `go build` on the output, verify it compiles.

### Phase 7: CLI (`cmd/bello/main.go`)

Wire everything together. Implement the toolchain commands from spec section 10:

| Command | Action |
|---|---|
| `bello papala file.🍌` | transpile + `go run` |
| `bello construccion [dir]` | transpile + `go build` |
| `bello kanpai [dir]` | transpile test files + `go test` |
| `bello sniff [dir]` | transpile + `go vet` |
| `bello bonito file.🍌` | format (parse + pretty-print Bello) |
| `bello dame pkg` | `go get pkg` |
| `bello modulo init name` | create `bello.🍑` |
| `bello splain` | show docs |
| `bello boosta [dir]` | bootstrap translator build + self-host validation |
| `bello boosta-run [dir] <command> [args...]` | build bootstrap translator then immediately run command |
| `bello micasa [dir]` | promote bootstrapped native compiler to `.bello/bello` |
| `bello completion [bash|zsh|fish]` | print shell completion scripts |

No legacy command aliases are kept; only canonical names above are supported.

**Error remapping**: when `go build` or `go test` reports errors, map Go source positions back to Bello source positions using the position map from the emitter. Reformat as `BEE DOH! <file>:<line>:<col> — <message>`.

### Phase 8: Module File Parser (`pkg/module/module.go`)

Separate parser for `bello.🍑` files. Uses the module file grammar from BNF section 14. Line-oriented (newlines are terminators, not semicolons).

Maps to `go.mod`:
- `modulo` -> `module`
- `bello X.Y` -> `go X.Y`
- `necesita` -> `require`
- `cambio` -> `replace`

---

## Testing Strategy

### Test Fixtures (`testdata/`)

Create `.🍌` files for every example in the spec:
- `hello.🍌` — minimal kampung jefe + poopaye
- `functions.🍌` — basic, multi-return, variadic, anonymous
- `control_flow.🍌` — po/ka, tulaliloo (all 4 forms), bee/doh/meh
- `structs.🍌` — luk/kampai, methods, embedding
- `interfaces.🍌` — buddies, type assertion, type switch
- `concurrency.🍌` — underpa, tatata, culo
- `generics.🍌` — type params, constraints
- `error_handling.🍌` — whaaat, BEE_DOH, gelatin
- `http_server.🍌` — spec section 12 complete example
- `worker_pool.🍌` — spec section 13 complete example
- `stdlib_minion.🍌` — boca/casa/amigos imports
- `stdlib_go.🍌` — fmt/os/sync imports (Go mode)
- `stdlib_mixed.🍌` — both modes in one file

### Test Levels

1. **Lexer tests** — token stream verification, semicolon insertion, edge cases
2. **Parser tests** — AST structure for each fixture, error recovery
3. **Transformer tests** — Go AST output matches expected Go code
4. **End-to-end tests** — `bello construccion` on each fixture produces a working binary
5. **Round-trip tests** — `bello bonito` on a file then re-parse produces identical AST

---

## Error Message Convention

All compiler/toolchain errors follow the format from spec section 11:

```
BEE DOH! <file>:<line>:<col> — <message>
```

Summary line at the end of a failed compilation:

```
POOPAYE! compilation naga success. <N> whaaat found.
```

Use clear technical descriptions in the message body. The `BEE DOH!` prefix is the character flavor. The content must be useful.

---

## Key Design Decisions

1. **Go stdlib for Go AST**: use `go/ast`, `go/token`, `go/printer`, `go/format` from the Go standard library for the Go-side AST and emission. This gives you correct Go output for free.

2. **Bello AST is separate**: do NOT try to reuse `go/ast` for Bello's AST. Bello has different keywords and needs source-position tracking back to `.🍌` files. Build a custom Bello AST that structurally mirrors `go/ast` for easy transformation.

3. **The transformer is a 1:1 tree walk**: every Bello AST node maps to exactly one Go AST node. There are no desugaring steps or complex transformations. This keeps the transpiler simple and debuggable.

4. **Predeclared identifiers are NOT special in the parser**: `si`, `naga`, `hana`, `BEE_DOH`, etc. are parsed as regular identifiers. The transformer rewrites them. Only `dala` and `pwede` need special parser handling (Type as first arg).

5. **Emoji file extensions**: use the actual emoji in filenames. The filesystem, Go's `os.Open`, and most modern terminals handle this fine. If portability is a concern, accept `.banana`, `.boom`, `bello.peach`, `bello.poo` as ASCII fallbacks.

6. **`jefe` is not a keyword**: it's a conventional identifier (like Go's `main`). The transformer rewrites `jefe` -> `main` only when it appears as a package name or a function name in package `jefe`.

---

## Implementation Language

The transpiler is written in **Go**. This gives:
- Direct access to `go/ast`, `go/printer`, `go/format`, `go/build`
- Easy invocation of `go build`, `go test`, `go vet` as subprocesses
- Path toward Phase 2 (self-hosting: rewrite in Bello, which transpiles to Go)

Use Go modules. The module path is `github.com/minions/bello` (or whatever the actual repo path is).

---

## Non-Goals for Phase 1

- No type checking (let `go build` handle it)
- No optimization passes
- No LSP/editor integration
- No custom runtime or GC
- No native code generation
- No cross-compilation beyond what `go build` provides
