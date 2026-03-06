package main

import (
	"errors"
	"fmt"
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
	case "construccion":
		runProjectCommand(expectArgOrDot(2), "build")
	case "kanpai":
		runProjectCommand(expectArgOrDot(2), "test")
	case "sniff":
		runProjectCommand(expectArgOrDot(2), "vet")
	case "bonito":
		runBonito(expectArg("bonito"))
	case "dame":
		runGoCommand("get", expectArg("dame"))
	case "modulo":
		if len(os.Args) < 4 || os.Args[2] != "init" {
			fail("usage: bello modulo init <name>")
		}
		runModuloInit(os.Args[3])
	case "splain":
		fmt.Println("Bello transpiler: .🍌 -> Go source")
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
	fmt.Println("bello papala file.🍌")
	fmt.Println("bello construccion [dir]")
	fmt.Println("bello kanpai [dir]")
	fmt.Println("bello bonito file.🍌")
	fmt.Println("bello dame pkg")
	fmt.Println("bello modulo init name")
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

	cmd := exec.Command(resolveGoBinary(), "run", ".")
	cmd.Dir = res.Workdir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		msg, n := parseGoToolError([]byte(stderr.String()), res.Maps)
		failWithSummary(msg, n)
	}
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
	res, err := buildProjectFromBello(path)
	if err != nil {
		fail(err.Error())
	}
	defer os.RemoveAll(res.Workdir)

	cmd := exec.Command(resolveGoBinary(), action, "./...")
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
	root, err := filepath.Abs(path)
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
