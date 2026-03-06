# Bello Language Reference

Bello is a Go-compatible language using Minion-speak tokens.
The compiler performs a source-to-source pass to Go, then invokes Go tooling.

## Reserved keywords

| Bello | Go meaning |
| --- | --- |
| `kampung` | `package` |
| `muak` | `import` |
| `banana` | `func` |
| `bapple` | `return` |
| `pooka` | `var` |
| `gelato` | `const` |
| `luk` | `type` |
| `kampai` | `struct` |
| `buddies` | `interface` |
| `po` | `if` |
| `ka` | `else` |
| `tulaliloo` | `for` |
| `tikali` | `range` |
| `buttom` | `break` |
| `bajo` | `continue` |
| `bee` | `switch` |
| `doh` | `case` |
| `meh` | `default` |
| `underpa` | `go` |
| `tatata` | `chan` |
| `culo` | `select` |
| `tank_yu` | `defer` |
| `patalaki` | `fallthrough` |
| `waaah` | `goto` |
| `dala` | `make` |
| `pwede` | `new` |

Notes:
- The parser also maps predeclared identifiers to Go built-ins listed below.
- `jefe` is rewritten to `main` when used as package or entry function name.

## Predeclared names

| Bello | Go |
| --- | --- |
| `me`, `me8`, `me16`, `me32`, `me64` | `int`, `int8`, `int16`, `int32`, `int64` |
| `ti`, `ti8`, `ti16`, `ti32`, `ti64` | `uint`, `uint8`, `uint16`, `uint32`, `uint64` |
| `la32`, `la64` | `float32`, `float64` |
| `butt` | `bool` |
| `bababa` | `string` |
| `todo` | `any` |
| `whaaat` | `error` |
| `si`, `naga`, `hana`, `mamamia` | `true`, `false`, `nil`, `iota` |
| `baboi`, `para_tu`, `stupa`, `cierro`, `yeet`, `mimik` | `append`, `len`, `cap`, `close`, `delete`, `copy` |
| `BEE_DOH`, `gelatin` | `panic`, `recover` |
| `poopaye` | `println` |

## Minion stdlib imports

| Bello import | Go import |
| --- | --- |
| `boca` | `fmt` |
| `casa` | `os` |
| `tubo` | `io` |
| `tubo_gordo` | `bufio` |
| `la_red` | `net/http` |
| `amigos` | `sync` |
| `amigos/tikitik` | `sync/atomic` |
| `tic_toc` | `time` |
| `bababas` | `strings` |
| `bababas/morph` | `strconv` |
| `kepala` | `math` |
| `kepala/loco` | `math/rand` |
| `kanpai` | `testing` |
| `kotak` | `encoding/json` |
| `pwesto` | `context` |
| `libretto` | `log` |
| `whaaats` | `errors` |
| `pila` | `sort` |
| `doodle` | `regexp` |
| `jalan` | `path/filepath` |
| `bandera` | `flag` |

Selector-call rewrites are applied for minion packages (for example `boca.blabla` -> `fmt.Printf`).

## Module descriptor (`bello.🍑`)

Grammar is intentionally small:

```text
modulo <module-path>
bello <go-version>
necesita (optional require block)
cambio (optional replace block)
```

Examples map to the corresponding `go.mod` statements:
- `modulo` -> `module`
- `bello 1.24` -> `go 1.24`
- `necesita` -> `require`
- `cambio` -> `replace`

## CLI quick map

- `bello papala <file.🍌>`: compile and run one file
- `bello construccion [dir]`: transpile and `go build`
- `bello kanpai [dir]`: transpile and `go test`
- `bello sniff [dir]`: transpile and `go vet`
- `bello bonito <file.🍌>`: print normalized Bello source
- `bello modulo init <name>`: write `bello.🍑`

Error messages are prefixed with `BEE DOH!` and formatted as:

```text
BEE DOH! <file>:<line>:<col> — <message>
```
