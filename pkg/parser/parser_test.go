package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFixtureFunctions(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "functions.🍌")
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	p := New(path, string(src))
	file, err := p.Parse()
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}

	if file.Package == nil || file.Package.Name != "jefe" {
		t.Fatalf("unexpected package name: %#v", file.Package)
	}
	if len(file.Decls) == 0 {
		t.Fatalf("expected declarations")
	}
}

func TestParseSemanticallyInvalidProgram(t *testing.T) {
	src := "kampung\nbanana main() {"
	p := New("bad.🍌", src)
	if _, err := p.Parse(); err == nil {
		t.Fatal("expected parse error")
	} else if !strings.Contains(err.Error(), "BEE DOH!") {
		t.Fatalf("wrong error format: %v", err)
	}
}

