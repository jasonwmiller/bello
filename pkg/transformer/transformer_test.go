package transformer

import (
	"go/ast"
	"path/filepath"
	"testing"

	"github.com/minions/bello/pkg/parser"
)

func TestTransformRewritesMinionImportAndIdentifiers(t *testing.T) {
	src := `kampung jefe

muak "boca"

banana jefe() {
	poopaye("hello")
}`
	p := parser.New(filepath.Join("..", "..", "testdata", "stdio_smoke.🍌"), src)
	f, err := p.Parse()
	if err != nil {
		t.Fatalf("parse source: %v", err)
	}

	gf, pm, err := Transform(f)
	if err != nil {
		t.Fatalf("transform: %v", err)
	}
	if pm == nil || pm.BelloFile != filepath.Join("..", "..", "testdata", "stdio_smoke.🍌") {
		t.Fatalf("position map lost")
	}
	if gf.Name == nil || gf.Name.Name != "main" {
		t.Fatalf("package name not rewritten: %#v", gf.Name)
	}
	if len(gf.Imports) != 1 {
		t.Fatalf("imports = %d", len(gf.Imports))
	}
	if !containsStringLiteral(gf.Imports[0].Path.Value, "fmt") {
		t.Fatalf("import not rewritten: %q", gf.Imports[0].Path.Value)
	}

	callFound := false
	ast.Inspect(gf, func(n ast.Node) bool {
		if call, ok := n.(*ast.SelectorExpr); ok {
			if call.Sel.Name == "Println" && call.X != nil {
				callFound = true
				return false
			}
		}
		return true
	})
	if !callFound {
		t.Fatalf("expected rewritten Println selector")
	}
}

func containsStringLiteral(s, want string) bool {
	if len(s) < 2 {
		return false
	}
	return s[1:len(s)-1] == want
}
