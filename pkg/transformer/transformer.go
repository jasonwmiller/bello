package transformer

import (
	"fmt"
	"go/ast"
	"go/format"
	goParser "go/parser"
	"go/scanner"
	"go/token"
	"strconv"
	"strings"

	"github.com/jasonwmiller/bello/pkg/lexer"
	"github.com/jasonwmiller/bello/pkg/parser"
)

// PositionMap stores output line mapping metadata.
type PositionMap struct {
	GoFile      string
	BelloFile   string
	LineOffset  int
}

func (m *PositionMap) Remap(_ string, line, col int) lexer.Position {
	return lexer.Position{Filename: m.BelloFile, Line: line + m.LineOffset, Column: col}
}

// Transform rewrites a parsed Bello file to Go AST.
func Transform(f *parser.File) (*ast.File, *PositionMap, error) {
	if f == nil {
		return nil, nil, fmt.Errorf("nil parser file")
	}
	if f.GoFile == nil {
		return nil, nil, fmt.Errorf("empty parser output")
	}

	g := f.GoFile
	aliases := map[string]string{}

	for _, imp := range g.Imports {
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}
		if goPath, ok := stdlibPackageMap[path]; ok {
			alias := path
			if imp.Name != nil && imp.Name.Name != "" && imp.Name.Name != "_" && imp.Name.Name != "." {
				alias = imp.Name.Name
			}
			imp.Path.Value = strconv.Quote(goPath)
			if imp.Name == nil {
				imp.Name = &ast.Ident{Name: alias}
			}
			aliases[alias] = goPath
		}
	}

	mapMain := false
	if g.Name != nil && g.Name.Name == "jefe" {
		g.Name.Name = "main"
		mapMain = true
	}

	ast.Walk(&rewriteVisitor{aliases: aliases, mapMain: mapMain}, g)

	return g, &PositionMap{GoFile: "", BelloFile: f.Filename, LineOffset: 0}, nil
}

// For emergency direct translation cases without parser AST.
func FallbackParse(filename string, src string) (*ast.File, error) {
	fset := token.NewFileSet()
	return goParser.ParseFile(fset, filename+".go", src, 0)
}

type rewriteVisitor struct {
	aliases map[string]string
	mapMain bool
}

func (v *rewriteVisitor) Visit(n ast.Node) ast.Visitor {
	switch x := n.(type) {
	case *ast.Ident:
		if repl, ok := builtinMap[x.Name]; ok {
			x.Name = repl
		}
	case *ast.SelectorExpr:
		pkg, ok := x.X.(*ast.Ident)
		if ok {
			if goPath, ok := v.aliases[pkg.Name]; ok {
				if repl, ok := rewriteMethodAlias(goPath, x.Sel.Name); ok {
					x.Sel.Name = repl
				}
			}
		}
	case *ast.FuncDecl:
		if v.mapMain && x.Name != nil && x.Name.Name == "jefe" {
			x.Name.Name = "main"
		}
	}
	return v
}

// RewriteGoSource is a helper for non-CLI callers when parser output is already
// available as source text and only textual rewrite is needed.
func RewriteGoSource(src string) ([]byte, error) {
	fset := token.NewFileSet()
	tree, err := goParser.ParseFile(fset, "bello_tmp.go", src, 0)
	if err != nil {
		return nil, err
	}
	ast.Walk(&rewriteVisitor{aliases: map[string]string{}}, tree)
	var out strings.Builder
	if err := format.Node(&out, fset, tree); err != nil {
		return nil, err
	}
	return []byte(out.String()), nil
}

// RewriteGoToBelloSource returns a Bello-formatted source from Go source.
func RewriteGoToBelloSource(src string) ([]byte, error) {
	goSrc, err := RewriteGoSource(src)
	if err != nil {
		return nil, err
	}
	return reverseToBello(string(goSrc))
}

// reverseToBello rewrites Go keywords and predeclared identifiers back to Bello words.
func reverseToBello(src string) ([]byte, error) {
	data := []byte(src)
	fset := token.NewFileSet()
	f := fset.AddFile("bello_tmp.go", -1, len(data))
	var scan scanner.Scanner
	scan.Init(f, data, nil, scanner.ScanComments)

	type tokenSpan struct {
		Offset int
		Token  token.Token
		Text   string
	}

	var spans []tokenSpan
	for {
		pos, tok, lit := scan.Scan()
		if tok == token.EOF {
			break
		}
		offset := f.Offset(pos)
		text := lit
		if text == "" {
			text = tok.String()
		}
		spans = append(spans, tokenSpan{Offset: offset, Token: tok, Text: text})
	}

	var out strings.Builder
	last := 0
	for _, span := range spans {
		if span.Offset < 0 || span.Offset > len(data) {
			return nil, fmt.Errorf("invalid token offset during bonito rewrite")
		}
		out.Write(data[last:span.Offset])
		out.WriteString(rewriteBelloToken(span.Token, span.Text))
		last = span.Offset + len(span.Text)
	}
	if last < len(data) {
		out.Write(data[last:])
	}
	return []byte(out.String()), nil
}

func rewriteBelloToken(tok token.Token, text string) string {
	if tok.IsKeyword() {
		if repl, ok := goKeywordInverse[text]; ok {
			return repl
		}
	}
	if tok == token.IDENT {
		if repl, ok := goBuiltinInverse[text]; ok {
			return repl
		}
	}
	return text
}
