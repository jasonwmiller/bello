package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBelloCLI_Bonito(t *testing.T) {
	goBin := mustResolveGoBinaryForTest(t)
	repoRoot := mustFindRepoRoot(t)

	fixture := filepath.Join(repoRoot, "testdata", "hello.🍌")
	out := runBelloCommand(t, goBin, repoRoot, nil, "bonito", fixture)
	if !strings.Contains(out, "kampung") {
		t.Fatalf("bonito output missing keyword mapping: %q", out)
	}
	if !strings.Contains(out, "poopaye") {
		t.Fatalf("bonito output missing built-in rewrite: %q", out)
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
	cmdPath := filepath.Join(repoRootForTests(), "cmd", "bello")
	cmd := exec.Command(goBinary, "run", cmdPath, args...)
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
	return filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filename)))
}
