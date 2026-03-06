# Bello Language Specification

**Version 0.1.0** — Draft, 2026-03-06

**Bello** — a compiled language with Go semantics and Minion-speak syntax.

> "Bello!" — Minion greeting

---

## 1. Design Philosophy

Bello is a statically typed, compiled language that mirrors Go's semantics — goroutines, channels, garbage collection, interfaces, multiple return values, and packages — but replaces all keywords and standard library names with Minion-speak derived from Italian, Spanish, French, English, Korean, Japanese, and gibberish as spoken by the Minions.

### Goals

- **Familiar semantics**: anyone who knows Go can read Bello after learning the keyword table
- **Compiled**: transpiles to Go as an initial backend, native backend later
- **Fun**: error messages, tooling, and idioms all stay in-character
- **Practical**: full interop with Go's ecosystem via the transpiler

---

## 2. Keywords

| Concept         | Go              | Bello         | Minion Origin                        |
|-----------------|-----------------|---------------|--------------------------------------|
| package         | `package`       | `kampung`     | Malay/Indonesian — village, community|
| import          | `import`        | `muak`        | Minion kiss sound — bringing in love |
| func            | `func`          | `banana`      | the Minion sacred fruit              |
| return          | `return`        | `bapple`      | Minion word — giving back            |
| var             | `var`           | `pooka`       | Minion exclamation                   |
| const           | `const`         | `gelato`      | Italian — frozen, unchanging         |
| type            | `type`          | `luk`         | Minion — look, define what you see   |
| struct          | `struct`        | `kampai`      | Japanese — cheers, a toast to types  |
| interface       | `interface`     | `buddies`     | Minion — friends who share behavior  |
| map             | `map`           | `papoy`       | Minion exclamation                   |
| if              | `if`            | `po`          | Minion conditional sound             |
| else            | `else`          | `ka`          | Minion continuation                  |
| for             | `for`           | `tulaliloo`   | from "tulaliloo ti amo"              |
| range           | `range`         | `tikali`      | Minion — tick through each one       |
| break           | `break`         | `buttom`      | Minion — bottom, the end             |
| continue        | `continue`      | `bajo`        | Spanish — under, keep going below    |
| switch          | `switch`        | `bee`         | Minion decision sound                |
| case            | `case`          | `doh`         | Minion reaction — "doh!"             |
| default         | `default`       | `meh`         | universal indifference               |
| go (goroutine)  | `go`            | `underpa`     | Minion — underwear, just go!         |
| chan            | `chan`           | `tatata`      | Minion — rapid communication         |
| select          | `select`        | `culo`        | Minion/Spanish — pick one            |
| defer           | `defer`         | `tank_yu`     | Minion — thank you, do it later      |
| make            | `make`          | `dala`        | Minion — make/do                     |
| new             | `new`           | `pwede`       | Filipino — can/possible              |
| nil             | `nil`           | `hana`        | Korean — one/nothing (context)       |
| true            | `true`          | `si`          | Spanish/Italian — yes                |
| false           | `false`         | `naga`        | Minion — no                          |
| error           | `error`         | `whaaat`      | Minion confusion                     |
| panic           | `panic`         | `BEE_DOH`     | Minion alarm siren                   |
| recover         | `recover`       | `gelatin`     | Minion — jelly, bounce back          |
| append          | `append`        | `baboi`       | Minion — add more                    |
| len             | `len`           | `para_tu`     | Minion — "for you", how much         |
| cap             | `cap`           | `stupa`       | Minion — capacity                    |
| fallthrough     | `fallthrough`   | `patalaki`    | Minion — fall down                   |
| println (builtin)| `println`      | `poopaye`     | Minion goodbye — quick debug output  |
| main            | `main`          | `jefe`        | Spanish — boss, the main one         |
| close           | `close`         | `cierro`      | Spanish — close/shut                 |
| delete          | `delete`        | `yeet`        | internet Minion energy               |
| copy            | `copy`          | `mimik`       | Minion mimicry                       |
| goto            | `goto`          | `waaah`       | Minion cry — jump somewhere crying   |
| iota            | `iota`          | `mamamia`     | Italian — auto-incrementing awe      |

---

## 3. Types

### 3.1 Primitive Types

| Go Type     | Bello Type | Notes                    |
|-------------|------------|--------------------------|
| `int`       | `me`       | platform-sized integer   |
| `int8`      | `me8`      |                          |
| `int16`     | `me16`     |                          |
| `int32`     | `me32`     |                          |
| `int64`     | `me64`     |                          |
| `uint`      | `ti`       | platform-sized unsigned  |
| `uint8`     | `ti8`      |                          |
| `uint16`    | `ti16`     |                          |
| `uint32`    | `ti32`     |                          |
| `uint64`    | `ti64`     |                          |
| `float32`   | `la32`     |                          |
| `float64`   | `la64`     |                          |
| `bool`      | `butt`     | si / naga                |
| `string`    | `bababa`   | Minion chatter           |
| `byte`      | `ti8`      | alias                    |
| `rune`      | `me32`     | alias                    |
| `any`       | `todo`     | empty interface          |

### 3.2 Composite Types

```
[]T              — slice of T
[N]T             — array of N elements of type T
papoy[K]V        — map from K to V
tatata T         — channel of T
tatata<- T       — send-only channel
<-tatata T       — receive-only channel
*T               — pointer to T
```

### 3.3 Zero Values

| Type      | Zero Value |
|-----------|------------|
| `me`      | `0`        |
| `la64`    | `0.0`      |
| `butt`    | `naga`     |
| `bababa`  | `""`       |
| pointer   | `hana`     |
| slice     | `hana`     |
| `papoy`   | `hana`     |
| `tatata`  | `hana`     |

---

## 4. Syntax

### 4.1 Package Declaration

Every Bello source file begins with a package declaration.

```bello
kampung jefe
```

### 4.2 Imports

```bello
muak "boca"

muak (
    "boca"
    "casa"
    "tic_toc"
)
```

### 4.3 Variables & Constants

```bello
// Explicit type
pooka nome bababa = "Bob"
pooka ojos me = 2

// Short declaration (type inferred)
nome := "Bob"
ojos := 2

// Constants
gelato BANANAS = 42
gelato PI la64 = 3.14159

// Constant block
gelato (
    A = 1
    B = 2
    C = 3
)

// Auto-incrementing with mamamia
gelato (
    Sunday    = mamamia  // 0
    Monday               // 1
    Tuesday              // 2
    Wednesday            // 3
    Thursday             // 4
    Friday               // 5
    Saturday             // 6
)
```

### 4.4 Functions

```bello
// Basic function
banana add(a me, b me) me {
    bapple a + b
}

// Multiple return values — assuming muak "boca"
banana divide(a la64, b la64) (la64, whaaat) {
    po b == 0.0 {
        bapple 0.0, boca.bee_doh_f("division by zero: %f / %f", a, b)
    }
    bapple a / b, hana
}

// Variadic
banana sum(nums ...me) me {
    total := 0
    tulaliloo _, n := tikali nums {
        total += n
    }
    bapple total
}

// Anonymous function
f := banana(x me) me { bapple x * x }
```

### 4.5 Control Flow

#### If / Else

```bello
po x > 0 {
    poopaye("positive")
} ka po x == 0 {
    poopaye("zero")
} ka {
    poopaye("negative")
}

// With init statement
po err := doThing(); err != hana {
    poopaye(err)
}
```

#### For Loop

Bello uses `tulaliloo` for all loop forms, just as Go uses `for`.

```bello
// C-style
tulaliloo i := 0; i < 10; i++ {
    poopaye(i)
}

// While-style
tulaliloo x > 0 {
    x--
}

// Infinite
tulaliloo {
    // ...
    buttom
}

// Range
tulaliloo i, v := tikali mySlice {
    poopaye(i, v)
}

// Range over map
tulaliloo k, v := tikali myPapoy {
    poopaye(k, v)
}

// Range over channel
tulaliloo msg := tikali ch {
    poopaye(msg)
}
```

#### Switch

```bello
bee fruit {
doh "banana":
    poopaye("BANANA!")
doh "apple", "papaya":
    poopaye("also good")
    patalaki
meh:
    poopaye("meh")
}

// Type switch
bee v := x.(luk) {
doh me:
    poopaye("is me")
doh bababa:
    poopaye("is bababa")
}
```

### 4.6 Structs

```bello
luk Minion kampai {
    Nome    bababa
    Ojos    me
    Tall    butt
}

// Instantiation
bob := Minion{
    Nome: "Bob",
    Ojos: 2,
    Tall: naga,
}

// Pointer
kevin := &Minion{Nome: "Kevin", Ojos: 2, Tall: si}
```

### 4.7 Methods

```bello
banana (m Minion) Greet() bababa {
    bapple "Bello! Me " + m.Nome + "!"
}

// Pointer receiver
banana (m *Minion) GrowTall() {
    m.Tall = si
}
```

### 4.8 Interfaces

```bello
luk Speaker buddies {
    Speak() bababa
}

luk Greeter buddies {
    Greet()
}

// todo is a builtin type (equivalent to Go's any/interface{})
// No need to define it — just use it:
// banana doStuff(val todo) { ... }
```

A type satisfies an interface implicitly — no `implements` keyword, just like Go.

> **Note:** `luk` serves double duty — it's used for type declarations (`luk Minion kampai`) and in type switches (`bee v := x.(luk)`). This mirrors Go where `type` is used in both contexts.

### 4.9 Goroutines

```bello
underpa doWork()

underpa banana() {
    poopaye("running in background!")
}()
```

### 4.10 Channels

```bello
// Unbuffered
ch := dala(tatata me)

// Buffered
ch := dala(tatata bababa, 10)

// Send
ch <- "banana"

// Receive
msg := <-ch

// Close
cierro(ch)
```

### 4.11 Select

```bello
culo {
doh msg := <-ch1:
    poopaye("from ch1:", msg)
doh msg := <-ch2:
    poopaye("from ch2:", msg)
doh <-tic_toc.apres(tic_toc.Tic):
    poopaye("timeout!")
meh:
    poopaye("nothing ready")
}
```

### 4.12 Defer / Panic / Recover

```bello
banana readFile(path bababa) {
    f, err := casa.buuka(path)
    po err != hana {
        BEE_DOH(err)
    }
    tank_yu f.cierro()

    // ... read file ...
}

banana safeCall() {
    tank_yu banana() {
        po r := gelatin(); r != hana {
            poopaye("Recovered:", r)
        }
    }()

    BEE_DOH("AAAAAAA!")
}
```

### 4.13 Error Handling

Errors follow Go convention: return `whaaat` as the last value. `whaaat` is a builtin interface type (equivalent to Go's `error`) — it does not need to be defined by the user.

```bello
// whaaat is a builtin:
//   luk whaaat buddies { Error() bababa }

// Custom error type
luk BananaError kampai {
    Code    me
    Message bababa
}

banana (e BananaError) Error() bababa {
    // uses fmt.Sprintf — assuming muak "boca"
    bapple boca.mumuak("BEE DOH %d: %s", e.Code, e.Message)
}

// Usage
banana peel(ripe me) (bababa, whaaat) {
    po ripe > 10 {
        bapple "", &BananaError{Code: 42, Message: "too ripe"}
    }
    bapple "yummy", hana
}
```

### 4.14 Pointers

```bello
pooka x me = 42
pooka p *me = &x

poopaye(*p)  // 42

*p = 100
poopaye(x)   // 100
```

### 4.15 Slices & Arrays

```bello
// Array (fixed size)
pooka arr [3]me = [3]me{1, 2, 3}

// Slice
nums := []me{1, 2, 3, 4, 5}
nums = baboi(nums, 6)

// Slice operations
sub := nums[1:3]

poopaye(para_tu(nums))  // length
poopaye(stupa(nums))    // capacity

// Make a slice
s := dala([]me, 0, 10)
```

### 4.16 Maps

```bello
// Literal
ages := papoy[bababa]me{
    "Bob":    12,
    "Kevin":  14,
    "Stuart": 11,
}

// Make
scores := dala(papoy[bababa]me)

// Access
age := ages["Bob"]

// Check existence
age, ok := ages["Dave"]
po !ok {
    poopaye("Dave naga found")
}

// Delete
yeet(ages, "Stuart")
```

### 4.17 Type Assertions & Type Switches

```bello
// Assertion
s := val.(bababa)

// Safe assertion
s, ok := val.(bababa)

// Type switch
bee v := val.(luk) {
doh me:
    poopaye("me:", v)
doh bababa:
    poopaye("bababa:", v)
meh:
    poopaye("whaaat is this")
}
```

### 4.18 Embedding

```bello
luk Animal kampai {
    Nome bababa
}

banana (a Animal) Speak() bababa {
    bapple a.Nome + " speaks"
}

luk Minion kampai {
    Animal              // embedded
    Ojos me
}

// Minion inherits Speak() from Animal
bob := Minion{
    Animal: Animal{Nome: "Bob"},
    Ojos:   2,
}
bob.Speak()  // "Bob speaks"
```

### 4.19 Generics

Bello supports type parameters, matching Go 1.18+ generics.

```bello
// Generic function
banana Max[T kepala.Ordered](a T, b T) T {
    po a > b {
        bapple a
    }
    bapple b
}

// Generic struct
luk Stack[T todo] kampai {
    items []T
}

banana (s *Stack[T]) Push(item T) {
    s.items = baboi(s.items, item)
}

banana (s *Stack[T]) Pop() (T, butt) {
    po para_tu(s.items) == 0 {
        pooka zero T
        bapple zero, naga
    }
    item := s.items[para_tu(s.items)-1]
    s.items = s.items[:para_tu(s.items)-1]
    bapple item, si
}

// Usage
s := Stack[me]{}
s.Push(42)
val, ok := s.Pop()
```

### 4.20 Goto

```bello
po somethingWrong {
    waaah cleanup
}

// ... normal path ...

cleanup:
    poopaye("cleaning up")
```

### 4.21 Type Conversions

Type conversions use the type name as a function, same as Go.

```bello
// Numeric conversions
x := me64(42)
y := la64(x)

// String <-> byte slice
bytes := []ti8("bello")
str := bababa(bytes)

// Rune conversion
r := me32('B')
```

---

## 5. Comments

```bello
// Single line comment — same as Go

/*
   Multi-line comment
   Banana banana banana
*/
```

---

## 6. Operators

All operators are identical to Go:

| Category    | Operators                          |
|-------------|------------------------------------|
| Arithmetic  | `+  -  *  /  %`                   |
| Comparison  | `==  !=  <  >  <=  >=`            |
| Logical     | `&&  \|\|  !`                      |
| Bitwise     | `&  \|  ^  &^  <<  >>`            |
| Assignment  | `=  :=  +=  -=  *=  /=  %=`       |
| Address     | `&  *`                             |
| Channel     | `<-`                               |
| Increment   | `++  --`                           |

---

## 7. Standard Library (`la_biblioteca`)

Bello supports two import modes. Importing via `la_biblioteca` (Minion names) gives you fully Minion-ified method and struct names. Importing Go packages directly gives you standard Go names. Both work, both can be mixed in the same file.

### 7.1 Import Modes

```bello
// Minion mode — full Minion method names
muak (
    "boca"
    "casa"
)
boca.poopaye("Bello!")
boca.blabla("Me %s, %d ojos\n", nome, ojos)    // Printf
boca.spitoo(casa.Bee_Doh, "BEE DOH\n")           // Fprintf

// Go mode — standard Go names, full interop
muak (
    "fmt"
    "os"
)
fmt.Println("Bello!")
fmt.Sprintf("Me %s, %d ojos", nome, ojos)
fmt.Fprintf(os.Stderr, "BEE DOH\n")

// Mix both in one file
muak (
    "boca"
    "fmt"
)
boca.poopaye("Minion style")
fmt.Println("Go style")
```

The transpiler maps Minion package + method names to their Go equivalents. Go imports pass through unchanged.

### 7.2 Package & Method Mapping

| Go Package       | Bello Package    | Description                        |
|------------------|------------------|------------------------------------|
| `fmt`            | `boca`           | formatted I/O (mouth)              |
| `os`             | `casa`           | OS operations (house)              |
| `io`             | `tubo`           | I/O primitives (tube)              |
| `bufio`          | `tubo_gordo`     | buffered I/O (fat tube)            |
| `net/http`       | `la_red`         | HTTP client/server (the net)       |
| `sync`           | `amigos`         | synchronization (friends)          |
| `sync/atomic`    | `amigos/tikitik` | atomic ops (Minion — tiny tiny)    |
| `time`           | `tic_toc`        | time and duration                  |
| `strings`        | `bababas`        | string manipulation (chatter)      |
| `strconv`        | `bababas/morph`  | string conversions (shapeshifting) |
| `math`           | `kepala`         | math (Malay — head/brain)          |
| `math/rand`      | `kepala/loco`    | random numbers (crazy head)        |
| `testing`        | `kanpai`         | test framework (Japanese — cheers) |
| `encoding/json`  | `kotak`          | JSON (Malay — box)                 |
| `context`        | `pwesto`         | context (Filipino — place/situation)|
| `log`            | `libretto`       | logging (Italian — little book)    |
| `errors`         | `whaaats`        | error utilities (Minion confusion) |
| `sort`           | `pila`           | sorting (Filipino — line up)       |
| `regexp`         | `doodle`         | regex (Minion scribbling)          |
| `path/filepath`  | `jalan`          | file paths (Malay — road/path)    |
| `flag`           | `bandera`        | CLI flags (Italian/Spanish — flag) |

### 7.3 boca (fmt) Method Names

| Go              | Bello (via `boca`)  | Minion Origin                      |
|-----------------|---------------------|------------------------------------|
| `Println`       | `poopaye`           | Minion goodbye — send output       |
| `Printf`        | `blabla`            | Minion chatter — formatted output  |
| `Sprintf`       | `mumuak`            | Minion kiss — whisper to string    |
| `Fprintf`       | `spitoo`            | Minion spit — write to a writer   |
| `Errorf`        | `bee_doh_f`         | alarm! — format an error          |
| `Scan`          | `huh`               | Minion — huh? (listen for input)  |
| `Scanf`         | `huh_huh`           | Minion — huh huh? (formatted)     |
| `Sscanf`        | `luk_luk`           | Minion — look look (from string)  |

### 7.4 casa (os) Method & Field Names

| Go              | Bello (via `casa`)  | Minion Origin                      |
|-----------------|---------------------|------------------------------------|
| `Open`          | `buuka`             | Minion — open up                   |
| `Create`        | `tada`              | Minion — ta-da! create something   |
| `Remove`        | `pchoo`             | Minion — laser blast, destroy      |
| `Exit`          | `bai_bai`           | Minion goodbye — exit process      |
| `Stdin`         | `Oreille`           | French — ear (input)               |
| `Stdout`        | `Boca`              | Minion/Italian — mouth (output)    |
| `Stderr`        | `Bee_Doh`           | Minion alarm — errors go here      |
| `Args`          | `Bagay`             | Filipino — things/stuff            |

**File object methods (returned by `casa.buuka`, `casa.tada`):**

| Go              | Bello               | Minion Origin                      |
|-----------------|---------------------|------------------------------------|
| `Close`         | `cierro`            | Minion/Spanish — shut it           |
| `Read`          | `nom_nom`           | Minion — consume/read bytes        |
| `Write`         | `spitoo`            | Minion — spit out bytes            |
| `Name`          | `nome`              | Minion — name                      |

> **Note:** `ReadAll` lives in `tubo` (io), not on file objects: `tubo.nom_nom_nom(reader)` maps to `io.ReadAll(reader)`.

### 7.5 amigos (sync) Method & Type Names

| Go              | Bello (via `amigos`) | Minion Origin                      |
|-----------------|----------------------|------------------------------------|
| `Mutex`         | `Jamu`               | Korean 자물 — lock (Minion-ified)  |
| `WaitGroup`     | `Chingus`            | Korean 친구 — friends (pluralized) |
| `Lock`          | `mwah`               | Minion kiss — lock it tight        |
| `Unlock`        | `bapapa`             | Minion — release!                  |
| `Add`           | `mas`                | Minion/Spanish — more!             |
| `Done`          | `listo`              | Minion/Spanish — ready!            |
| `Wait`          | `hmmmm`              | Minion — waiting patiently         |

### 7.6 la_red (net/http)

| Go                | Bello (via `la_red`)   | Minion Origin                      |
|-------------------|------------------------|------------------------------------|
| `HandleFunc`      | `ooh_ooh`              | Minion excitement — handle this!   |
| `ListenAndServe`  | `bello_bello`          | Minion greeting — welcome requests |
| `Get`             | `gimme`                | Minion — give me!                  |
| `Post`            | `takka`                | Minion — take this!                |
| `ResponseWriter`  | `Reponsu`              | Minion-ified French réponse        |
| `Request`         | `Juseyo`               | Korean 주세요 — please give me     |
| `Server`          | `Jefe_Red`             | Minion — boss of the net           |
| `StatusOK`        | `TodoBien`             | all good!                          |
| `StatusNotFound`  | `NagaAqui`             | naga here!                         |

### 7.7 tic_toc (time)

| Go              | Bello (via `tic_toc`) | Minion Origin                      |
|-----------------|-----------------------|------------------------------------|
| `Now`           | `nau`                 | Minion — now!                      |
| `Sleep`         | `zzzzz`               | Minion snoring                     |
| `After`         | `apres`               | French — after                     |
| `Second`        | `Tic`                 | one tick                           |
| `Minute`        | `Tic_Tic`             | many ticks                         |
| `Hour`          | `Tic_Tic_Tic`         | so many ticks                      |

### 7.8 kanpai (testing)

> **Casing convention for BEE_DOH variants:** The builtin panic is `BEE_DOH` (all caps — screaming). Formatted method variants use lowercase: `bee_doh_f`. Capital `BEE_DOH` as a method name (e.g., `t.BEE_DOH()`) maps to fatal/instant-kill methods that mirror the builtin's severity.

| Go              | Bello (via `kanpai`)  | Minion Origin                      |
|-----------------|-----------------------|------------------------------------|
| `T`             | `T`                   | kept short — test context          |
| `B`             | `B`                   | kept short — benchmark context     |
| `Errorf`        | `bee_doh_f`           | alarm! — test failure              |
| `Fatalf`        | `BEE_DOH_F`           | BIG alarm — fatal test failure     |
| `Fatal`         | `BEE_DOH`             | instant death                      |
| `Skip`          | `pfft`                | Minion dismissal — skip this       |
| `Run`           | `dale`                | Minion — go go, run a subtest      |
| `Log`           | `psst`                | Minion whisper — log info          |
| `Helper`        | `shh`                 | Minion — quiet, I'm a helper      |

### 7.9 whaaats (errors)

| Go              | Bello (via `whaaats`) | Minion Origin                      |
|-----------------|-----------------------|------------------------------------|
| `New`           | `uh_oh`               | Minion — new problem!              |
| `Is`            | `sama_sama`           | Minion — same same? (comparison)   |
| `As`            | `luk_como`            | Minion — look like (type assert)   |
| `Unwrap`        | `peela`               | Minion — peel the banana (unwrap)  |

### 7.10 Example: Same Program, Two Styles

```bello
// === Full Minion ===
kampung jefe

muak (
    "boca"
    "amigos"
)

banana jefe() {
    pooka wg amigos.Chingus
    wg.mas(1)
    underpa banana() {
        tank_yu wg.listo()
        boca.poopaye("Bello!")
    }()
    wg.hmmmm()
}
```

```bello
// === Go interop ===
kampung jefe

muak (
    "fmt"
    "sync"
)

banana jefe() {
    pooka wg sync.WaitGroup
    wg.Add(1)
    underpa banana() {
        tank_yu wg.Done()
        fmt.Println("Bello!")
    }()
    wg.Wait()
}
```

Both compile to identical Go output.

---

## 8. Testing

Test files end in `.💥`. Test functions start with `Kanpai`. Benchmarks start with `Rapido`.

```bello
kampung kepala

muak "kanpai"

banana KanpaiAdd(t *kanpai.T) {
    result := Add(2, 3)
    po result != 5 {
        t.bee_doh_f("expected 5, got %d, WHAAAT?!", result)
    }
}

banana RapidoAdd(b *kanpai.B) {
    tulaliloo i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

---

## 9. Modules & Packages

### 9.1 Module File: `bello.🍑`

Every module root contains exactly one `bello.🍑` file (the fixed filename, like Go's `go.mod`).

```
modulo github.com/minions/myproject

bello 1.0

necesita (
    github.com/minions/banana-lib v1.2.3
)
```

| Go term    | Bello term  |
|------------|-------------|
| `module`   | `modulo`    |
| `go`       | `bello`     |
| `require`  | `necesita`  |
| `replace`  | `cambio`    |

### 9.2 Package Visibility

Same rule as Go: identifiers starting with an uppercase letter are exported.

```bello
banana Greet() bababa { ... }   // exported
banana helper() bababa { ... }  // unexported
```

---

## 10. Toolchain

| Command                        | Description                           | Go Equivalent    |
|--------------------------------|---------------------------------------|------------------|
| `bello papala file.🍌`        | compile and run (papala = hurry!)     | `go run`         |
| `bello construccion`           | compile to binary (build it!)         | `go build`       |
| `bello kanpai ./...`           | run all tests (cheers!)               | `go test ./...`  |
| `bello bonito file.🍌`        | format source code (make it pretty)   | `gofmt`          |
| `bello dame pkg`               | fetch dependency (give me!)           | `go get`         |
| `bello modulo init name`       | initialize new module (bello.🍑)     | `go mod init`    |
| `bello sniff`                  | static analysis (smell the code)      | `go vet`         |
| `bello splain`                 | show documentation (explain!)         | `go doc`         |

---

## 11. Compiler Errors

Compiler errors are in Minion-speak for character but include clear technical information.

```
BEE DOH! file.🍌:12:5 — "bob" naga declared in this kampung
BEE DOH! file.🍌:7:20 — banana greet() bapple bababa, but got me
BEE DOH! file.🍌:3:1 — muak "boca" declared but naga used
POOPAYE! compilation naga success. 3 whaaat found.
```

### Error Format

```
BEE DOH! <file>:<line>:<col> — <message>
```

---

## 12. Complete Example: HTTP Server

```bello
kampung jefe

muak (
    "boca"
    "la_red"
)

banana handler(w la_red.Reponsu, r *la_red.Juseyo) {
    boca.spitoo(w, "Bello! You visited %s\n", r.URL.Path)
}

banana jefe() {
    la_red.ooh_ooh("/", handler)
    boca.poopaye("Server running on :8080")

    po err := la_red.bello_bello(":8080", hana); err != hana {
        BEE_DOH(err)
    }
}
```

---

## 13. Complete Example: Concurrent Worker Pool

```bello
kampung jefe

muak (
    "boca"
    "amigos"
)

banana worker(id me, jobs <-tatata me, results tatata<- me, wg *amigos.Chingus) {
    tank_yu wg.listo()
    tulaliloo j := tikali jobs {
        boca.blabla("Minion %d nom-nom job %d\n", id, j)
        results <- j * 2
    }
}

banana jefe() {
    numJobs := 10
    jobs := dala(tatata me, numJobs)
    results := dala(tatata me, numJobs)

    pooka wg amigos.Chingus
    tulaliloo w := 1; w <= 3; w++ {
        wg.mas(1)
        underpa worker(w, jobs, results, &wg)
    }

    tulaliloo j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    cierro(jobs)

    underpa banana() {
        wg.hmmmm()
        cierro(results)
    }()

    tulaliloo r := tikali results {
        boca.poopaye("Result:", r)
    }
}
```

---

## 14. Implementation Strategy

### Phase 1: Transpiler (Bello → Go)

The initial compiler is a source-to-source transpiler:

1. **Lexer**: tokenize Bello keywords, identifiers, literals, operators
2. **Parser**: build AST (structure mirrors Go's AST)
3. **Transformer**: map Bello AST nodes to Go AST nodes (keyword substitution + package name mapping)
4. **Emitter**: output valid Go source
5. **Build**: invoke `go build` on the generated Go source

This approach gives full Go ecosystem access from day one.

### Phase 2: Self-Hosting

Rewrite the transpiler in Bello itself.

### Phase 3: Native Backend (Optional)

Replace the Go backend with LLVM or a custom code generator for independent compilation.

---

## 15. File Extensions

| Extension       | Purpose              |
|-----------------|----------------------|
| `.🍌`           | source file          |
| `.💥`           | test file            |
| `bello.🍑`      | module definition — fixed filename (the bottom of it all) |
| `bello.💩`      | dependency checksums — fixed filename (trust but verify)  |

---

## 16. Reserved Words

All keywords listed in Section 2 are reserved and cannot be used as identifiers.

### Built-in Functions

| Go         | Bello      |
|------------|------------|
| `make`     | `dala`     |
| `new`      | `pwede`    |
| `append`   | `baboi`    |
| `len`      | `para_tu`  |
| `cap`      | `stupa`    |
| `close`    | `cierro`   |
| `delete`   | `yeet`     |
| `copy`     | `mimik`    |
| `panic`    | `BEE_DOH`  |
| `recover`  | `gelatin`  |
| `println`  | `poopaye`  |


---

*POOPAYE!*
