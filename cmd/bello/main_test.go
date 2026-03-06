package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var (
	builtOnce sync.Once
	buildErr  error
)

func TestBelloCLI_Bonito(t *testing.T) {
	goBin := mustResolveGoBinaryForTest(t)
	repoRoot := mustFindRepoRoot(t)

	fixture := filepath.Join(repoRoot, "testdata", "hello.🍌")
	out := runBelloCommand(t, goBin, repoRoot, nil, "bonito", fixture)
	if !strings.Contains(out, "kampung") {
		t.Fatalf("bonito output missing keyword mapping: %q", out)
	}
	if !strings.Contains(out, "blabla") {
		t.Fatalf("bonito output missing stdlib call rewrite: %q", out)
	}
}

func TestBelloCLI_ProjectCommands(t *testing.T) {
	goBin := mustResolveGoBinaryForTest(t)
	repoRoot := mustFindRepoRoot(t)

	workDir := t.TempDir()
	fixtureSrc := filepath.Join(repoRoot, "testdata", "hello.🍌")
	targetSrc := filepath.Join(workDir, "hello.🍌")
	if err := copyFile(fixtureSrc, targetSrc); err != nil {
		t.Fatalf("prepare fixture: %v", err)
	}

	out := runBelloCommand(t, goBin, workDir, nil, "papala", targetSrc)
	if strings.TrimSpace(out) != "bello" {
		t.Fatalf("papala output = %q", out)
	}

	cmds := []string{"construccion", "kanpai", "sniff"}
	for _, cmd := range cmds {
		out = runBelloCommand(t, goBin, workDir, nil, cmd)
		if strings.Contains(out, "BEE DOH!") {
			t.Fatalf("%s failed with error: %q", cmd, out)
		}
	}
}

func TestBelloCLI_MockModuleInit(t *testing.T) {
	goBin := mustResolveGoBinaryForTest(t)
	workDir := t.TempDir()
	out := runBelloCommand(t, goBin, workDir, map[string]string{"GO111MODULE": "off"}, "modulo", "init", "example.com/bello/demo")
	if !strings.Contains(out, "created bello.🍑") {
		t.Fatalf("modulo init output unexpected: %q", out)
	}
	if _, err := os.Stat(filepath.Join(workDir, "bello.🍑")); err != nil {
		t.Fatalf("expected bello.🍑 file in %s", workDir)
	}

}

func runBelloCommand(t *testing.T, goBinary, workDir string, env map[string]string, args ...string) string {
	t.Helper()

	binary := buildBelloBinary(t, goBinary)
	cmdArgs := append([]string{}, args...)
	cmd := exec.Command(binary, cmdArgs...)
	cmd.Dir = workDir
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GO_BIN="+goBinary)
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run bello %s: %v\noutput=%s", strings.Join(args, " "), err, string(out))
	}
	return strings.TrimSpace(string(out))
}

func buildBelloBinary(t *testing.T, goBinary string) string {
	t.Helper()

	root := repoRootForTests()
	output := filepath.Join(root, ".tools", "bello-integration")
	builtOnce.Do(func() {
		if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
			buildErr = err
			return
		}
		cmd := exec.Command(goBinary, "build", "-o", output, "./cmd/bello")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil {
			buildErr = fmt.Errorf("build bello binary: %w: %s", err, strings.TrimSpace(string(out)))
			return
		}
	})
	if buildErr != nil {
		t.Fatal(buildErr)
	}
	return output
}

func mustResolveGoBinaryForTest(t *testing.T) string {
	t.Helper()
	goBinary := resolveGoBinary()
	if err := exec.Command(goBinary, "version").Run(); err != nil {
		t.Skip("go tool unavailable for CLI integration tests:", err)
	}
	return goBinary
}

func mustFindRepoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cwd: %v", err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}
	t.Fatalf("could not find repository root from %s", cwd)
	return ""
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

// repoRootForTests captures the package directory at runtime without additional filesystem scans.
func repoRootForTests() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
}
