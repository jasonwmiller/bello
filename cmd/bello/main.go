package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/minions/bello/pkg/emitter"
	"github.com/minions/bello/pkg/module"
	"github.com/minions/bello/pkg/parser"
	"github.com/minions/bello/pkg/transformer"
)

type buildResult struct {
	Workdir string
	Maps    map[string]*transformer.PositionMap
}

var goBinary string

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "papala":
		runPapala(expectArg("papala"))
	case "repl":
		runRepl()
	case "chiku":
		runRepl()
	case "construccion":
		runProjectCommand(expectArgOrDot(2), "build")
	case "kanpai":
		runProjectCommand(expectArgOrDot(2), "test")
	case "sniff":
		runProjectCommand(expectArgOrDot(2), "vet")
	case "bonito":
		runBonito(expectArg("bonito"))
	case "bootstrap", "boosta":
		runBootstrap(expectArgOrDot(2))
	case "dame":
		runGoCommand("get", expectArg("dame"))
	case "modulo":
		if len(os.Args) < 4 || os.Args[2] != "init" {
			fail("usage: bello modulo init <name>")
		}
		runModuloInit(os.Args[3])
	case "splain":
		fmt.Println("Bello transpiler: .🍌 -> Go source")
	case "completion":
		shell := "bash"
		if len(os.Args) > 2 {
			shell = os.Args[2]
		}
		runCompletion(shell)
	default:
		printUsage()
		os.Exit(1)
	}
}

func requireGoTool() {
	bin := resolveGoBinary()
	if err := exec.Command(bin, "version").Run(); err != nil {
		var ee *exec.Error
		if errors.As(err, &ee) && ee.Err == exec.ErrNotFound {
			fail("BEE DOH! -:1:1 — go tool not found in PATH")
		}
		fail("BEE DOH! -:1:1 — go tool unavailable: " + err.Error())
	}
}

func resolveGoBinary() string {
	if goBinary != "" {
		return goBinary
	}
	if v := strings.TrimSpace(os.Getenv("GO_BIN")); v != "" {
		goBinary = v
		return goBinary
	}
	if p, err := exec.LookPath("go"); err == nil {
		goBinary = p
		return goBinary
	}
	for _, p := range []string{"/usr/local/go/bin/go", "/usr/bin/go"} {
		if _, err := os.Stat(p); err == nil {
			goBinary = p
			return goBinary
		}
	}
	goBinary = "go"
	return goBinary
}

func printUsage() {
	fmt.Println("bello papala file.🍌 [arg ...]")
	fmt.Println("bello repl")
	fmt.Println("bello chiku")
	fmt.Println("bello construccion [dir]")
	fmt.Println("bello kanpai [dir]")
	fmt.Println("bello bootstrap [dir]")
	fmt.Println("bello boosta [dir]")
	fmt.Println("bello bonito file.🍌")
	fmt.Println("bello dame pkg")
	fmt.Println("bello modulo init name")
	fmt.Println("bello completion [bash|zsh|fish]")
	fmt.Println("bello sniff [dir]")
	fmt.Println("bello splain")
}

func expectArg(cmd string) string {
	if len(os.Args) < 3 {
		fail(cmd + " requires an argument")
	}
	return os.Args[2]
}

func expectArgOrDot(i int) string {
	if len(os.Args) > i {
		return os.Args[i]
	}
	return "."
}

func runPapala(file string) {
	requireGoTool()
	res, err := transpileSingle(file)
	if err != nil {
		fail(err.Error())
	}
	defer os.RemoveAll(res.Workdir)

	args := []string{"run", "."}
	if len(os.Args) > 3 {
		args = append(args, os.Args[3:]...)
	}

	cmd := exec.Command(resolveGoBinary(), args...)
	cmd.Dir = res.Workdir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	var stderr strings.Builder
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	err = cmd.Run()
	if err != nil {
		msg, n := parseGoToolError([]byte(stderr.String()), res.Maps)
		failWithSummary(msg, n)
	}
}

func runRepl() {
	requireGoTool()
	fmt.Println("bello repl — use /chiku for help, /bapple to bounce")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("🍌> ")
		if !scanner.Scan() {
			return
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		switch strings.ToLower(line) {
		case "/bapple", "/return", "/quit", "/exit":
			return
		case "/chiku", "/help":
			fmt.Println("bello commands: /chiku (help), /bapple (return)")
			fmt.Println("type a full Bello statement, or a full statement block")
			continue
		}

		if err := runReplLine(line); err != nil {
			fmt.Println(err)
		}
	}
}

func runReplLine(line string) error {
	replSrc := "kampung jefe\n\nbanana main() {\n\t" + line + "\n}\n"
	tempFile, err := os.CreateTemp("", "bello-repl-*.🍌")
	if err != nil {
		return fmt.Errorf("BEE DOH! -:1:1 — cannot create repl file: %v", err)
	}
	path := tempFile.Name()
	if _, err := tempFile.WriteString(replSrc); err != nil {
		tempFile.Close()
		os.Remove(path)
		return fmt.Errorf("BEE DOH! -:1:1 — cannot write repl file: %v", err)
	}
	tempFile.Close()
	defer os.Remove(path)

	res, err := transpileSingle(path)
	if err != nil {
		return err
	}
	defer os.RemoveAll(res.Workdir)

	cmd := exec.Command(resolveGoBinary(), "run", ".")
	cmd.Dir = res.Workdir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	var stderr strings.Builder
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	err = cmd.Run()
	if err != nil {
		msg, n := parseGoToolError([]byte(stderr.String()), res.Maps)
		if n > 0 {
			return fmt.Errorf("POOPAYE! compilation naga success. %d whaaat found.\n%s", n, msg)
		}
		return fmt.Errorf(msg)
	}
	return nil
}

func runBonito(file string) {
	src, err := os.ReadFile(file)
	if err != nil {
		fail(err.Error())
	}
	p := parser.New(file, string(src))
	f, err := p.Parse()
	if err != nil {
		fail(err.Error())
	}

	goSrc, err := transformer.RewriteGoToBelloSource(f.Translated)
	if err != nil {
		goSrc = []byte(f.Translated)
	}
	os.Stdout.Write(goSrc)
}

func runProjectCommand(path string, action string) {
	requireGoTool()
	workOnDir := true
	var res *buildResult
	var err error
	if filepath.Ext(path) == ".🍌" {
		workOnDir = false
		res, err = transpileSingle(path)
	} else {
		res, err = buildProjectFromBello(path)
	}
	if err != nil {
		fail(err.Error())
	}
	defer os.RemoveAll(res.Workdir)

	cmdArgs := []string{action, "./..."}
	if !workOnDir {
		cmdArgs = []string{action, "."}
	}
	cmd := exec.Command(resolveGoBinary(), cmdArgs...)
	cmd.Dir = res.Workdir
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		msg, n := parseGoToolError([]byte(out.String()), res.Maps)
		failWithSummary(msg, n)
	}
	if strings.TrimSpace(out.String()) != "" {
		fmt.Println(out.String())
	}
}

func runBootstrap(root string) {
	requireGoTool()
	absRoot := root
	if absRoot == "" {
		absRoot = "."
	}
	var err error
	absRoot, err = filepath.Abs(absRoot)
	if err != nil {
		fail("BEE DOH! -:1:1 — cannot resolve bootstrap root: " + err.Error())
	}

	workspace, err := os.MkdirTemp("", "bello-bootstrap-")
	if err != nil {
		fail("BEE DOH! -:1:1 — cannot create bootstrap workspace: " + err.Error())
	}
	defer os.RemoveAll(workspace)

	seed, hasSeed := locateBootstrapSeed(absRoot)
	if hasSeed {
		fmt.Println("bello bootstrap: using prebuilt minion seed in", seed)
		if err := copyBootstrapSource(seed, workspace); err != nil {
			fail("BEE DOH! -:1:1 — cannot copy bootstrap seed source: " + err.Error())
		}
	} else {
		fmt.Println("bello bootstrap: generating Bello bootstrap source from", absRoot)
		if err := copyBootstrapModuleFiles(absRoot, workspace); err != nil {
			fail("BEE DOH! -:1:1 — cannot write bootstrap module files: " + err.Error())
		}
		if err := convertGoSourcesToBello(absRoot, workspace); err != nil {
			fail(err.Error())
		}
	}

	binPath := filepath.Join(workspace, "bello.bootstrap")
	fmt.Println("bello bootstrap: building bootstrap compiler with native translator")
	if err := buildBootstrapBinary(workspace, binPath); err != nil {
		fail(err.Error())
	}

	fmt.Println("bello bootstrap: validating bootstrap compiler self-host pass")
	if _, err := runBelloBinaryCommand(binPath, workspace, "construccion", "."); err != nil {
		fail(err.Error())
	}
	fmt.Println("bello bootstrap: self-host validation complete")
}

func runCompletion(shell string) {
	switch strings.ToLower(shell) {
	case "bash", "sh", "":
		fmt.Println(`# Install shell completion for bash:
#   source <(bello completion)
_bello_complete() {
	local cur prev
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"
	COMPREPLY=()
	if [[ ${COMP_CWORD} -eq 1 ]]; then
		COMPREPLY=( $(compgen -W "papala repl chiku construccion kanpai sniff bonito bootstrap boosta completion dame modulo splain" -- "$cur") )
		return 0
	fi
	case "$prev" in
		modulo)
			COMPREPLY=( $(compgen -W "init" -- "$cur") )
			return 0
			;;
		construccion|kanpai|sniff|bootstrap|boosta|completion|papala|bonito|dame)
			COMPREPLY=( $(compgen -f -- "$cur") )
			return 0
			;;
	esac
	COMPREPLY=( $(compgen -f -- "$cur") )
}
complete -F _bello_complete bello` + "\n")
		return
	case "zsh":
		fmt.Println(`# Install shell completion for zsh:
#   autoload -Uz compinit && compinit
#   source <(bello completion zsh)
_bello() {
  local -a cmds
  cmds=("papala" "repl" "chiku" "construccion" "kanpai" "sniff" "bonito" "bootstrap" "boosta" "completion" "dame" "modulo" "splain")
  _describe 'bello command' cmds
}
compdef _bello bello`)
		return
	case "fish":
		fmt.Println(`# Install shell completion for fish:
#   bello completion fish | source
complete -c bello -n "__fish_use_subcommand" -a "papala repl chiku construccion kanpai sniff bonito bootstrap boosta completion dame modulo splain"`)
		return
	default:
		fail("unsupported shell for completion: " + shell)
	}
}

func locateBootstrapSeed(root string) (string, bool) {
	candidate := filepath.Join(root, "bootstrap", "src")
	if isBootstrapSeedDir(candidate) {
		return candidate, true
	}
	if isBootstrapSeedDir(root) {
		return root, true
	}
	return "", false
}

func isBootstrapSeedDir(dir string) bool {
	if _, err := os.Stat(dir); err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		return false
	}
	_, err := os.Stat(filepath.Join(dir, "cmd", "bello", "main.🍌"))
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(dir, "pkg", "parser", "parser.🍌"))
	return err == nil
}

func copyBootstrapSource(sourceRoot, workspace string) error {
	return filepath.WalkDir(sourceRoot, func(path string, d os.DirEntry, errIn error) error {
		if errIn != nil {
			return errIn
		}
		rel, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		dst := filepath.Join(workspace, rel)
		if d.IsDir() {
			if shouldSkipBootstrapDir(path) {
				return filepath.SkipDir
			}
			return os.MkdirAll(dst, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0o644)
	})
}

func runGoCommand(args ...string) {
	requireGoTool()
	cmd := exec.Command(resolveGoBinary(), args...)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fail("BEE DOH! " + strings.TrimSpace(out.String()))
	}
	if strings.TrimSpace(out.String()) != "" {
		fmt.Print(out.String())
	}
}

func runModuloInit(name string) {
	mod := &module.ModuleFile{
		ModulePath: name,
		GoVersion:  "1.23",
	}
	if err := os.WriteFile("bello.🍑", []byte(mod.RenderGoMod()), 0o644); err != nil {
		fail(err.Error())
	}
	fmt.Println("created bello.🍑")
}

func transpileSingle(path string) (*buildResult, error) {
	tmpDir, err := os.MkdirTemp("", "bello-single-")
	if err != nil {
		return nil, err
	}
	if err := writeSingleModuleRoot(tmpDir, filepath.Dir(path)); err != nil {
		return nil, err
	}
	goPath, pm, err := transpile(path, tmpDir)
	if err != nil {
		return nil, err
	}
	return &buildResult{
		Workdir: filepath.Dir(goPath),
		Maps:    map[string]*transformer.PositionMap{goPath: pm},
	}, nil
}

func transpile(path, outDir string) (string, *transformer.PositionMap, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}
	p := parser.New(path, string(src))
	f, err := p.Parse()
	if err != nil {
		return "", nil, err
	}

	gf, pm, err := transformer.Transform(f)
	if err != nil {
		return "", nil, err
	}
	goPath, pm2, err := emitter.Emit(gf, path, outDir)
	if pm2 != nil {
		pm = pm2
	}
	return goPath, pm, err
}

func buildProjectFromBello(path string) (*buildResult, error) {
	rootPath := path
	if filepath.Ext(rootPath) == ".🍌" {
		rootPath = filepath.Dir(rootPath)
	}
	root, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	tmp, err := os.MkdirTemp("", "bello-build-")
	if err != nil {
		return nil, err
	}
	res := &buildResult{Workdir: tmp, Maps: map[string]*transformer.PositionMap{}}

	if err := prepareModuleRoot(root, tmp); err != nil {
		return nil, err
	}

	if err := filepath.WalkDir(root, func(p string, d os.DirEntry, errIn error) error {
		if errIn != nil {
			return errIn
		}
		rel, errRel := filepath.Rel(root, p)
		if errRel != nil {
			return errRel
		}
		if d.IsDir() {
			if rel == "." {
				return nil
			}
			if rel == ".git" || strings.HasPrefix(rel, ".git"+string(filepath.Separator)) {
				return filepath.SkipDir
			}
			return os.MkdirAll(filepath.Join(tmp, rel), 0o755)
		}
		if filepath.Ext(p) != ".🍌" {
			return nil
		}

		goPath, pm, err := transpile(p, filepath.Join(tmp, filepath.Dir(rel)))
		if err != nil {
			return err
		}
		if pm != nil {
			res.Maps[goPath] = pm
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}

func prepareModuleRoot(root, outDir string) error {
	belloModPath := filepath.Join(root, "bello.🍑")
	if moduleLike(belloModPath) {
		mf, err := module.Parse(belloModPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(outDir, "go.mod"), []byte(mf.RenderGoMod()), 0o644); err != nil {
			return err
		}
		return copyGoSum(root, outDir)
	}
	goModPath := filepath.Join(root, "go.mod")
	if moduleLike(goModPath) {
		b, err := os.ReadFile(goModPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(outDir, "go.mod"), b, 0o644); err != nil {
			return err
		}
		return copyGoSum(root, outDir)
	}
	return os.WriteFile(filepath.Join(outDir, "go.mod"), []byte((&module.ModuleFile{
		ModulePath: module.ModuleNameFromPath(root),
		GoVersion:  "1.24",
	}).RenderGoMod()), 0o644)
}

func writeSingleModuleRoot(workDir, srcRoot string) error {
	beloModPath := filepath.Join(srcRoot, "go.mod")
	if moduleLike(beloModPath) {
		b, err := os.ReadFile(beloModPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(workDir, "go.mod"), b, 0o644); err != nil {
			return err
		}
		return copyGoSum(srcRoot, workDir)
	}

	return os.WriteFile(filepath.Join(workDir, "go.mod"), []byte((&module.ModuleFile{
		ModulePath: module.ModuleNameFromPath(srcRoot),
		GoVersion:  getGoMajorMinor(),
	}).RenderGoMod()), 0o644)
}

func copyGoSum(srcRoot, outDir string) error {
	goSumPath := filepath.Join(srcRoot, "go.sum")
	if !moduleLike(goSumPath) {
		return nil
	}
	b, err := os.ReadFile(goSumPath)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "go.sum"), b, 0o644)
}

func getGoMajorMinor() string {
	v := runtime.Version()
	v = strings.TrimPrefix(v, "go")
	parts := strings.Split(v, ".")
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	return "1.24"
}

func moduleLike(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func parseGoToolError(output []byte, maps map[string]*transformer.PositionMap) (string, int) {
	txt := strings.TrimSpace(string(output))
	if txt == "" {
		return "BEE DOH! go command failed", 1
	}

	re := regexp.MustCompile(`(?m)^(.+?\.go):([0-9]+):([0-9]+)(?::[0-9]+)?:\s*(.*)$`)
	matches := re.FindAllStringSubmatch(txt, -1)
	if len(matches) == 0 {
		return "BEE DOH! " + txt, 1
	}

	var b strings.Builder
	for _, m := range matches {
		file, line, col := m[1], m[2], m[3]
		msg := strings.TrimSpace(m[4])
		ln, _ := strconv.Atoi(line)
		cl, _ := strconv.Atoi(col)
		if pm, ok := maps[file]; ok {
			pos := emitter.RemapError(pm, file, ln, cl)
			b.WriteString(fmt.Sprintf("BEE DOH! %s:%d:%d — %s\n", pos.Filename, pos.Line, pos.Column, msg))
			continue
		}
		b.WriteString(fmt.Sprintf("BEE DOH! %s:%d:%d — %s\n", file, ln, cl, msg))
	}

	return strings.TrimSpace(b.String()), len(matches)
}

func runBelloBinaryCommand(binary, workdir string, args ...string) (string, error) {
	cmd := exec.Command(binary, args...)
	cmd.Dir = workdir
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GO_BIN="+resolveGoBinary())

	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return strings.TrimSpace(out.String()), fmt.Errorf("BEE DOH! bootstrap runner failed: %w: %s", err, strings.TrimSpace(out.String()))
	}
	return strings.TrimSpace(out.String()), nil
}

func buildBootstrapBinary(sourceRoot, out string) error {
	res, err := buildProjectFromBello(sourceRoot)
	if err != nil {
		return fmt.Errorf("BEE DOH! bootstrap seed build failed: %v", err)
	}
	defer os.RemoveAll(res.Workdir)

	cmd := exec.Command(resolveGoBinary(), "build", "-o", out, "./cmd/bello")
	cmd.Dir = res.Workdir
	var outBuilder strings.Builder
	cmd.Stdout = &outBuilder
	cmd.Stderr = &outBuilder

	if err := cmd.Run(); err != nil {
		msg, n := parseGoToolError([]byte(outBuilder.String()), res.Maps)
		return wrapMappedError(msg, n, err)
	}
	return nil
}

func wrapMappedError(msg string, errCount int, err error) error {
	if errCount > 0 {
		return fmt.Errorf("BEE DOH! bootstrap translation failed: %s", msg)
	}
	return fmt.Errorf("BEE DOH! bootstrap translation failed: %v", err)
}

func copyBootstrapModuleFiles(sourceRoot, workspace string) error {
	goModPath := filepath.Join(sourceRoot, "go.mod")
	goSumPath := filepath.Join(sourceRoot, "go.sum")
	if moduleLike(goModPath) {
		b, err := os.ReadFile(goModPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(workspace, "go.mod"), b, 0o644); err != nil {
			return err
		}
		if moduleLike(goSumPath) {
			if b, err := os.ReadFile(goSumPath); err == nil {
				if err := os.WriteFile(filepath.Join(workspace, "go.sum"), b, 0o644); err != nil {
					return err
				}
			}
		}
		return nil
	}

	return os.WriteFile(filepath.Join(workspace, "go.mod"), []byte((&module.ModuleFile{
		ModulePath: module.ModuleNameFromPath(sourceRoot),
		GoVersion:  getGoMajorMinor(),
	}).RenderGoMod()), 0o644)
}

func convertGoSourcesToBello(sourceRoot, workspace string) error {
	return filepath.WalkDir(sourceRoot, func(path string, d os.DirEntry, errIn error) error {
		if errIn != nil {
			return errIn
		}
		if d.IsDir() {
			if shouldSkipBootstrapDir(path) {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		rel, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return err
		}
		clean := filepath.ToSlash(rel)
		parts := strings.Split(clean, "/")
		if len(parts) < 2 {
			return nil
		}
		if parts[0] != "cmd" && parts[0] != "pkg" {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		belloSrc, err := transformer.RewriteGoToBelloSource(string(src))
		if err != nil {
			return fmt.Errorf("BEE DOH! cannot translate %s: %v", path, err)
		}

		outRel := strings.TrimSuffix(rel, ".go") + ".🍌"
		dst := filepath.Join(workspace, outRel)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, belloSrc, 0o644); err != nil {
			return err
		}
		return nil
	})
}

func shouldSkipBootstrapDir(path string) bool {
	base := filepath.Base(path)
	if base == ".git" || base == ".jj" || base == ".tools" || strings.HasPrefix(base, ".") {
		return true
	}
	return false
}

func failWithSummary(msg string, errorCount int) {
	if errorCount > 0 {
		countText := "1"
		if errorCount > 1 {
			countText = fmt.Sprintf("%d", errorCount)
		}
		fmt.Printf("POOPAYE! compilation naga success. %s whaaat found.\n", countText)
	}
	fail(msg)
}

func fail(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
