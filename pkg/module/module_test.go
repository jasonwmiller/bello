package module

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseModuleFileWithBlocks(t *testing.T) {
	src := `modulo example.com/minions/demo
bello 1.22

necesita (
	github.com/example/lib v1.2.3
)

cambio example.com/a v1.0.0 => example.com/b v2.0.0
`
	root := t.TempDir()
	modPath := filepath.Join(root, "bello.🍑")
	if err := os.WriteFile(modPath, []byte(src), 0o644); err != nil {
		t.Fatalf("write module fixture: %v", err)
	}

	mf, err := Parse(modPath)
	if err != nil {
		t.Fatalf("parse module: %v", err)
	}

	if mf.ModulePath != "example.com/minions/demo" {
		t.Fatalf("module path = %q", mf.ModulePath)
	}
	if mf.GoVersion != "1.22" {
		t.Fatalf("go version = %q", mf.GoVersion)
	}
	if got := len(mf.Requires); got != 1 {
		t.Fatalf("requires length = %d", got)
	}
	if mf.Requires[0] != (Require{Module: "github.com/example/lib", Version: "v1.2.3"}) {
		t.Fatalf("unexpected require: %+v", mf.Requires[0])
	}
	if got := len(mf.Replaces); got != 1 {
		t.Fatalf("replaces length = %d", got)
	}
	if mf.Replaces[0] != (Replace{OldMod: "example.com/a", OldVer: "v1.0.0", NewMod: "example.com/b", NewVer: "v2.0.0"}) {
		t.Fatalf("unexpected replace: %+v", mf.Replaces[0])
	}
}

func TestRenderGoModRoundTrip(t *testing.T) {
	mf := &ModuleFile{
		ModulePath: "example.com/demo",
		GoVersion:  "1.22",
		Requires:   []Require{{Module: "github.com/a", Version: "v1.0.0"}},
		Replaces:   []Replace{{OldMod: "github.com/x", OldVer: "v0.1.0", NewMod: "github.com/y", NewVer: "v0.2.0"}},
	}

	out := mf.RenderGoMod()
	if got := expectedTextContains(out, []string{
		"module example.com/demo",
		"go 1.22",
		"require github.com/a v1.0.0",
		"replace github.com/x v0.1.0 => github.com/y v0.2.0",
	}); !got {
		t.Fatalf("rendered go.mod unexpected:\n%s", out)
	}
}

func expectedTextContains(text string, lines []string) bool {
	for _, line := range lines {
		if !containsLine(text, line) {
			return false
		}
	}
	return true
}

func containsLine(text, line string) bool {
	return strings.Contains(text, line)
}
