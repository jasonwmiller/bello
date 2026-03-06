package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	goParser "go/parser"
	"go/scanner"
	"go/token"
	"strings"

	"github.com/minions/bello/pkg/lexer"
)

// Parser converts Bello source into a Go-compatible AST representation.
type Parser struct {
	filename string
	input    string
	lex      *lexer.Lexer
}

func New(filename, src string) *Parser {
	return &Parser{
		filename: filename,
		input:    src,
		lex:      lexer.New(filename, src),
	}
}

// Parse returns a parser File containing a translated Go AST and a Bello-mirrored AST.
func (p *Parser) Parse() (*File, error) {
	translated, err := p.translateToGo()
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	node, gerr := goParser.ParseFile(fset, p.filename+".go", translated, goParser.AllErrors|goParser.ParseComments)
	if gerr != nil {
		return nil, formatParseErr(gerr, p.filename)
	}

	astFile := &File{
		GoFile:     node,
		Filename:   p.filename,
		Translated: translated,
		Posn:       lexer.Position{Filename: p.filename, Line: 1, Column: 1},
	}
	if node.Name != nil {
		astFile.Package = &Ident{
			BaseNode: BaseNode{Posn: toPos(fset, node.Name.Pos())},
			Name:     node.Name.Name,
		}
	}
	for _, imp := range node.Imports {
		importSpec := &ImportSpec{
			BaseNode: BaseNode{Posn: toPos(fset, imp.Path.Pos())},
			Path:     strings.Trim(imp.Path.Value, "\""),
		}
		if imp.Name != nil && imp.Name.Name != "" {
			importSpec.Name = &Ident{
				BaseNode: BaseNode{Posn: toPos(fset, imp.Name.Pos())},
				Name:     imp.Name.Name,
			}
		}
		astFile.Imports = append(astFile.Imports, importSpec)
	}

	astFile.Decls = append(astFile.Decls, convertDecls(node.Decls, fset)...)
	return astFile, nil
}

func (p *Parser) translateToGo() (string, error) {
	var out bytes.Buffer
	var prev lexer.TokenType = -1

	tok := p.lex.Next()
	for tok.Type != lexer.EOF {
		if tok.Type == lexer.ILLEGAL {
			return "", tokenError(tok, "illegal token")
		}

		lit := translateLit(tok)
		if needSpace(prev, tok.Type) {
			out.WriteByte(' ')
		}
		out.WriteString(lit)
		prev = tok.Type
		tok = p.lex.Next()
	}
	return out.String(), nil
}

var keywordMap = map[string]string{
	"kampung":   "package",
	"muak":      "import",
	"banana":    "func",
	"bapple":    "return",
	"pooka":     "var",
	"gelato":    "const",
	"luk":       "type",
	"kampai":    "struct",
	"buddies":   "interface",
	"papoy":     "map",
	"po":        "if",
	"ka":        "else",
	"tulaliloo": "for",
	"tikali":    "range",
	"buttom":    "break",
	"bajo":      "continue",
	"bee":       "switch",
	"doh":       "case",
	"meh":       "default",
	"underpa":   "go",
	"tatata":    "chan",
	"culo":      "select",
	"tank_yu":   "defer",
	"patalaki":  "fallthrough",
	"waaah":     "goto",
	"dala":      "make",
	"pwede":     "new",
}

var predeclaredMap = map[string]string{
	"me":      "int",
	"me8":     "int8",
	"me16":    "int16",
	"me32":    "int32",
	"me64":    "int64",
	"ti":      "uint",
	"ti8":     "uint8",
	"ti16":    "uint16",
	"ti32":    "uint32",
	"ti64":    "uint64",
	"la32":    "float32",
	"la64":    "float64",
	"butt":    "bool",
	"bababa":  "string",
	"todo":    "any",
	"whaaat":  "error",
	"si":      "true",
	"naga":    "false",
	"hana":    "nil",
	"mamamia": "iota",
	"baboi":   "append",
	"para_tu": "len",
	"stupa":   "cap",
	"cierro":  "close",
	"yeet":    "delete",
	"mimik":   "copy",
	"BEE_DOH": "panic",
	"gelatin": "recover",
	"poopaye": "println",
}

func translateLit(tok lexer.Token) string {
	if repl, ok := keywordMap[tok.Lit]; ok {
		return repl
	}
	if repl, ok := predeclaredMap[tok.Lit]; ok {
		return repl
	}
	return tok.Lit
}

func needSpace(prev lexer.TokenType, curr lexer.TokenType) bool {
	if prev == -1 || prev == lexer.SEMICOLON {
		return false
	}
	if curr == lexer.EOF {
		return false
	}
	if curr == lexer.LPAREN || curr == lexer.LBRACK || curr == lexer.LBRACE {
		return isWord(prev)
	}
	if prev == lexer.LPAREN || prev == lexer.LBRACK || prev == lexer.LBRACE || prev == lexer.PERIOD {
		return false
	}
	if curr == lexer.PERIOD || curr == lexer.COMMA || curr == lexer.COLON || curr == lexer.RPAREN ||
		curr == lexer.RBRACK || curr == lexer.RBRACE || curr == lexer.SEMICOLON ||
		curr == lexer.INC || curr == lexer.DEC || curr == lexer.ELLIPSIS {
		return false
	}
	if prev == lexer.PERIOD {
		return false
	}
	if isWord(prev) && isWord(curr) {
		return true
	}
	return false
}

func isWord(tt lexer.TokenType) bool {
	switch tt {
	case lexer.IDENT, lexer.INT, lexer.FLOAT, lexer.IMAGINARY, lexer.RUNE, lexer.STRING,
		lexer.KAMPUNG, lexer.MUAK, lexer.BANANA, lexer.BAPPLE, lexer.POOKA, lexer.GELATO,
		lexer.LUK, lexer.KAMPAI, lexer.BUDDIES, lexer.PAPOY, lexer.PO, lexer.KA, lexer.TULALILOO,
		lexer.TIKALI, lexer.BUTTOM, lexer.BAJO, lexer.BEE, lexer.DOH, lexer.MEH, lexer.UNDERPA,
		lexer.TATATA, lexer.CULO, lexer.TANK_YU, lexer.PATALAKI, lexer.WAAAH, lexer.DALA, lexer.PWEDE:
		return true
	default:
		return false
	}
}

func tokenError(tok lexer.Token, msg string) error {
	return fmt.Errorf("BEE DOH! %s:%d:%d — %s", tok.Pos.Filename, tok.Pos.Line, tok.Pos.Column, msg)
}

func formatParseErr(err error, filename string) error {
	if list, ok := err.(scanner.ErrorList); ok && len(list) > 0 {
		first := list[0]
		return fmt.Errorf("BEE DOH! %s:%d:%d — %s", filename, first.Pos.Line, first.Pos.Column, strings.TrimSuffix(first.Msg, "\n"))
	}

	msg := strings.TrimSpace(err.Error())
	if msg == "" {
		return fmt.Errorf("BEE DOH! %s — go parse failed", filename)
	}
	return fmt.Errorf("BEE DOH! %s — %s", filename, msg)
}

func toPos(fset *token.FileSet, p token.Pos) lexer.Position {
	pos := fset.Position(p)
	return lexer.Position{Filename: pos.Filename, Offset: pos.Offset, Line: pos.Line, Column: pos.Column}
}

func convertDecls(decls []ast.Decl, fset *token.FileSet) []Decl {
	var out []Decl
	for _, decl := range decls {
		out = append(out, convertDecl(decl, fset)...)
	}
	return out
}

func convertDecl(decl ast.Decl, fset *token.FileSet) []Decl {
	switch x := decl.(type) {
	case *ast.FuncDecl:
		return []Decl{toFuncDecl(x, fset)}
	case *ast.GenDecl:
		return toGenDecl(x, fset)
	default:
		return nil
	}
}

func toGenDecl(decl *ast.GenDecl, fset *token.FileSet) []Decl {
	out := make([]Decl, 0, len(decl.Specs))
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.ValueSpec:
			if decl.Tok == token.VAR {
				out = append(out, toVarDecl(s, fset))
				continue
			}
			if decl.Tok == token.CONST {
				out = append(out, toConstDecl(s, fset))
				continue
			}
		case *ast.TypeSpec:
			if decl.Tok == token.TYPE {
				out = append(out, toTypeDecl(s, fset))
			}
		}
	}
	return out
}

func toVarDecl(v *ast.ValueSpec, fset *token.FileSet) Decl {
	return &VarDecl{
		BaseNode: BaseNode{Posn: toPos(fset, v.Pos())},
		Names:    toIdentList(v.Names, fset),
		Type:     toExpr(v.Type, fset),
		Values:   toExprs(v.Values, fset),
	}
}

func toConstDecl(v *ast.ValueSpec, fset *token.FileSet) Decl {
	return &ConstDecl{
		BaseNode: BaseNode{Posn: toPos(fset, v.Pos())},
		Names:    toIdentList(v.Names, fset),
		Type:     toExpr(v.Type, fset),
		Values:   toExprs(v.Values, fset),
	}
}

func toTypeDecl(t *ast.TypeSpec, fset *token.FileSet) Decl {
	return &TypeDecl{
		BaseNode: BaseNode{Posn: toPos(fset, t.Pos())},
		Name:     toIdent(t.Name),
		Type:     toExpr(t.Type, fset),
	}
}

func toFuncDecl(f *ast.FuncDecl, fset *token.FileSet) Decl {
	out := &FuncDecl{
		BaseNode:  BaseNode{Posn: toPos(fset, f.Pos())},
		Recv:      toFieldList(f.Recv, fset),
		Name:      toIdent(f.Name),
		Type:      toFuncType(f.Type, fset),
		Body:      toBlock(f.Body, fset),
		TypeParams: nil,
	}
	if f.Type != nil && f.Type.TypeParams != nil {
		out.TypeParams = toFieldList(f.Type.TypeParams, fset)
	}
	return out
}

func toBlock(b *ast.BlockStmt, fset *token.FileSet) *BlockStmt {
	if b == nil {
		return nil
	}
	return &BlockStmt{
		BaseNode: BaseNode{Posn: toPos(fset, b.Pos())},
		List:     toStmts(b.List, fset),
	}
}

func toStmts(stmts []ast.Stmt, fset *token.FileSet) []Stmt {
	out := make([]Stmt, 0, len(stmts))
	for _, s := range stmts {
		if d, ok := s.(*ast.DeclStmt); ok {
			for _, dec := range convertDecl(d.Decl, fset) {
				out = append(out, &DeclStmt{
					BaseNode: BaseNode{Posn: toPos(fset, d.Pos())},
					Decl:     dec,
				})
			}
			continue
		}
		stmt := toStmt(s, fset)
		if stmt != nil {
			out = append(out, stmt)
		}
	}
	return out
}

func toStmt(stmt ast.Stmt, fset *token.FileSet) Stmt {
	switch x := stmt.(type) {
	case *ast.BlockStmt:
		return toBlock(x, fset)
	case *ast.ReturnStmt:
		return &ReturnStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Results:  toExprs(x.Results, fset),
		}
	case *ast.IfStmt:
		return &IfStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Init:     toStmt(x.Init, fset),
			Cond:     toExpr(x.Cond, fset),
			Body:     toBlock(x.Body, fset),
			Else:     toStmt(x.Else, fset),
		}
	case *ast.ForStmt:
		return &ForStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Init:     toStmt(x.Init, fset),
			Cond:     toExpr(x.Cond, fset),
			Post:     toStmt(x.Post, fset),
			Body:     toBlock(x.Body, fset),
		}
	case *ast.RangeStmt:
		return &ForStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Range: &RangeClause{
				BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
				Key:      toExpr(x.Key, fset),
				Value:    toExpr(x.Value, fset),
				X:        toExpr(x.X, fset),
			},
			Body: toBlock(x.Body, fset),
		}
	case *ast.SwitchStmt:
		return &SwitchStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Init:     toStmt(x.Init, fset),
			Tag:      toExpr(x.Tag, fset),
			Body:     toSwitchClauses(x.Body.List, fset),
		}
	case *ast.TypeSwitchStmt:
		ts := &TypeSwitchStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Body:     toTypeSwitchClauses(x.Body.List, fset),
		}
		if x.Init != nil {
			ts.Init = toStmt(x.Init, fset)
		}
		if x.Assign != nil {
			ts.Assign = toStmt(x.Assign, fset)
		}
		return ts
	case *ast.SelectStmt:
		return &SelectStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Body:     toCommClauses(x.Body.List, fset),
		}
	case *ast.GoStmt:
		if x.Call == nil {
			return nil
		}
		return &GoStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Call:     toCallExpr(x.Call, fset),
		}
	case *ast.DeferStmt:
		if x.Call == nil {
			return nil
		}
		return &DeferStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Call:     toCallExpr(x.Call, fset),
		}
	case *ast.BranchStmt:
		return &BranchStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Tok:      x.Tok.String(),
			Label:    toIdent(x.Label),
		}
	case *ast.LabeledStmt:
		return &LabeledStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Label:    toIdent(x.Label),
			Stmt:     toStmt(x.Stmt, fset),
		}
	case *ast.AssignStmt:
		return &AssignStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Lhs:      toExprs(x.Lhs, fset),
			Op:       x.Tok.String(),
			Rhs:      toExprs(x.Rhs, fset),
		}
	case *ast.SendStmt:
		return &SendStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Chan:     toExpr(x.Chan, fset),
			Value:    toExpr(x.Value, fset),
		}
	case *ast.IncDecStmt:
		return &IncDecStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			X:        toExpr(x.X, fset),
			Tok:      x.Tok.String(),
		}
	case *ast.ExprStmt:
		return &ExprStmt{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			X:        toExpr(x.X, fset),
		}
	case *ast.EmptyStmt:
		return nil
	case *ast.DeclStmt:
		return nil
	case nil:
		return nil
	default:
		return nil
	}
}

func toSwitchClauses(clauses []ast.Stmt, fset *token.FileSet) []Stmt {
	out := make([]Stmt, 0, len(clauses))
	for _, stmt := range clauses {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}
		out = append(out, &CaseClause{
			BaseNode: BaseNode{Posn: toPos(fset, cc.Pos())},
			Values:  toExprs(cc.List, fset),
			Body:    toStmts(cc.Body, fset),
			Default: len(cc.List) == 0,
		})
	}
	return out
}

func toTypeSwitchClauses(clauses []ast.Stmt, fset *token.FileSet) []Stmt {
	out := make([]Stmt, 0, len(clauses))
	for _, stmt := range clauses {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}
		out = append(out, &CaseClause{
			BaseNode: BaseNode{Posn: toPos(fset, cc.Pos())},
			Values:   toExprs(cc.List, fset),
			Body:     toStmts(cc.Body, fset),
			Default:  len(cc.List) == 0,
		})
	}
	return out
}

func toCommClauses(clauses []ast.Stmt, fset *token.FileSet) []Stmt {
	out := make([]Stmt, 0, len(clauses))
	for _, stmt := range clauses {
		cc, ok := stmt.(*ast.CommClause)
		if !ok {
			continue
		}
		out = append(out, &CommClause{
			BaseNode: BaseNode{Posn: toPos(fset, cc.Pos())},
			Comm:     toStmt(cc.Comm, fset),
			Body:     toStmts(cc.Body, fset),
		})
	}
	return out
}

func toExprs(exprs []ast.Expr, fset *token.FileSet) []Expr {
	out := make([]Expr, 0, len(exprs))
	for _, e := range exprs {
		if ce := toExpr(e, fset); ce != nil {
			out = append(out, ce)
		}
	}
	return out
}

func toIdentList(list []*ast.Ident, fset *token.FileSet) []*Ident {
	out := make([]*Ident, 0, len(list))
	for _, id := range list {
		out = append(out, toIdent(id))
	}
	return out
}

func toIdent(id *ast.Ident) *Ident {
	if id == nil {
		return nil
	}
	return &Ident{Name: id.Name}
}

func toFieldList(fl *ast.FieldList, fset *token.FileSet) *FieldList {
	if fl == nil {
		return nil
	}
	out := &FieldList{BaseNode: BaseNode{Posn: toPos(fset, fl.Pos())}}
	for _, f := range fl.List {
		out.List = append(out.List, toField(f, fset))
	}
	return out
}

func toField(f *ast.Field, fset *token.FileSet) *Field {
	field := &Field{
		BaseNode: BaseNode{Posn: toPos(fset, f.Pos())},
		Names:    toIdentList(f.Names, fset),
		Tag:      nil,
		Ellipsis: false,
	}
	if f.Tag != nil {
		field.Tag = toBasicLit(f.Tag, fset)
	}

	switch t := f.Type.(type) {
	case *ast.Ellipsis:
		field.Ellipsis = true
		field.Type = toExpr(t.Elt, fset)
	default:
		field.Type = toExpr(f.Type, fset)
	}
	return field
}

func toCallExpr(expr ast.Expr, fset *token.FileSet) *CallExpr {
	c, ok := expr.(*ast.CallExpr)
	if !ok {
		if expr == nil {
			return nil
		}
		return &CallExpr{BaseNode: BaseNode{Posn: toPos(fset, expr.Pos())}, Fun: toExpr(expr, fset)}
	}
	return &CallExpr{
		BaseNode: BaseNode{Posn: toPos(fset, c.Pos())},
		Fun:      toExpr(c.Fun, fset),
		Args:     toExprs(c.Args, fset),
		Ellipsis: c.Ellipsis.IsValid(),
	}
}

func toExpr(e ast.Expr, fset *token.FileSet) Expr {
	if e == nil {
		return nil
	}

	switch x := e.(type) {
	case *ast.Ident:
		return &Ident{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Name:     x.Name,
		}
	case *ast.BasicLit:
		return toBasicLit(x, fset)
	case *ast.BinaryExpr:
		return &BinaryExpr{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			X:        toExpr(x.X, fset),
			Op:       toLexerToken(x.Op),
			Y:        toExpr(x.Y, fset),
		}
	case *ast.UnaryExpr:
		return &UnaryExpr{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Op:       toLexerToken(x.Op),
			X:        toExpr(x.X, fset),
		}
	case *ast.CallExpr:
		return toCallExpr(x, fset)
	case *ast.IndexExpr:
		return &IndexExpr{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			X:        toExpr(x.X, fset),
			Index:    toExpr(x.Index, fset),
		}
	case *ast.SliceExpr:
		return &SliceExpr{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			X:        toExpr(x.X, fset),
			Low:      toExpr(x.Low, fset),
			High:     toExpr(x.High, fset),
			Max:      toExpr(x.Max, fset),
			Full:     x.Slice3,
		}
	case *ast.SelectorExpr:
		return &SelectorExpr{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			X:        toExpr(x.X, fset),
			Sel:      toExpr(x.Sel, fset),
		}
	case *ast.TypeAssertExpr:
		return &TypeAssertExpr{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			X:        toExpr(x.X, fset),
			Type:     toExpr(x.Type, fset),
		}
	case *ast.CompositeLit:
		lit := &CompositeLit{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Type:     toExpr(x.Type, fset),
		}
		for _, el := range x.Elts {
			lit.Elts = append(lit.Elts, toExpr(el, fset))
		}
		return lit
	case *ast.FuncLit:
		return &FuncLit{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Type:     toFuncType(x.Type, fset),
			Body:     toBlock(x.Body, fset),
		}
	case *ast.ParenExpr:
		return toExpr(x.X, fset)
	case *ast.ArrayType:
		if x.Len == nil {
			return &SliceType{
				BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
				Elt:      toExpr(x.Elt, fset),
			}
		}
		return &ArrayType{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Len:      toExpr(x.Len, fset),
			Elt:      toExpr(x.Elt, fset),
		}
	case *ast.MapType:
		return &MapType{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Key:      toExpr(x.Key, fset),
			Value:    toExpr(x.Value, fset),
		}
	case *ast.ChanType:
		dir := ""
		switch x.Dir {
		case ast.SEND:
			dir = "send"
		case ast.RECV:
			dir = "recv"
		default:
			dir = ""
		}
		return &ChanType{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Dir:      dir,
			Value:    toExpr(x.Value, fset),
		}
	case *ast.StarExpr:
		return &PointerType{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Elt:      toExpr(x.X, fset),
		}
	case *ast.FuncType:
		return toFuncType(x, fset)
	case *ast.StructType:
		return &StructType{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Fields:   toFieldList(x.Fields, fset),
		}
	case *ast.InterfaceType:
		return &InterfaceType{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Methods:  toFieldList(x.Methods, fset),
		}
	case *ast.KeyValueExpr:
		return &KeyValueExpr{
			BaseNode: BaseNode{Posn: toPos(fset, x.Pos())},
			Key:      toExpr(x.Key, fset),
			Value:    toExpr(x.Value, fset),
		}
	default:
		return nil
	}
}

func toBasicLit(lit *ast.BasicLit, fset *token.FileSet) *BasicLit {
	if lit == nil {
		return nil
	}
	out := &BasicLit{
		BaseNode: BaseNode{Posn: toPos(fset, lit.Pos())},
		Value:    lit.Value,
	}
	switch lit.Kind {
	case token.INT:
		out.Kind = lexer.INT
	case token.FLOAT:
		out.Kind = lexer.FLOAT
	case token.IMAG:
		out.Kind = lexer.IMAGINARY
	case token.CHAR:
		out.Kind = lexer.RUNE
	case token.STRING:
		out.Kind = lexer.STRING
	default:
		out.Kind = lexer.STRING
	}
	return out
}

func toFuncType(fn *ast.FuncType, fset *token.FileSet) *FuncType {
	if fn == nil {
		return nil
	}
	typ := &FuncType{
		BaseNode:   BaseNode{Posn: toPos(fset, fn.Pos())},
		TypeParams: toFieldList(fn.TypeParams, fset),
		Params:     toFieldList(fn.Params, fset),
		Results:    toFieldList(fn.Results, fset),
	}
	return typ
}

func toLexerToken(tok token.Token) lexer.TokenType {
	switch tok {
	case token.ADD:
		return lexer.ADD
	case token.SUB:
		return lexer.SUB
	case token.MUL:
		return lexer.MUL
	case token.QUO:
		return lexer.QUO
	case token.REM:
		return lexer.REM
	case token.AND:
		return lexer.AND
	case token.OR:
		return lexer.OR
	case token.XOR:
		return lexer.XOR
	case token.SHL:
		return lexer.SHL
	case token.SHR:
		return lexer.SHR
	case token.AND_NOT:
		return lexer.AND_NOT
	case token.LAND:
		return lexer.LAND
	case token.LOR:
		return lexer.LOR
	case token.ARROW:
		return lexer.ARROW
	case token.INC:
		return lexer.INC
	case token.DEC:
		return lexer.DEC
	case token.EQL:
		return lexer.EQL
	case token.NEQ:
		return lexer.NEQ
	case token.LSS:
		return lexer.LSS
	case token.LEQ:
		return lexer.LEQ
	case token.GTR:
		return lexer.GTR
	case token.GEQ:
		return lexer.GEQ
	case token.ASSIGN:
		return lexer.ASSIGN
	case token.DEFINE:
		return lexer.DEFINE
	case token.ADD_ASSIGN:
		return lexer.ADD_ASSIGN
	case token.SUB_ASSIGN:
		return lexer.SUB_ASSIGN
	case token.MUL_ASSIGN:
		return lexer.MUL_ASSIGN
	case token.QUO_ASSIGN:
		return lexer.QUO_ASSIGN
	case token.REM_ASSIGN:
		return lexer.REM_ASSIGN
	case token.AND_ASSIGN:
		return lexer.AND_ASSIGN
	case token.OR_ASSIGN:
		return lexer.OR_ASSIGN
	case token.XOR_ASSIGN:
		return lexer.XOR_ASSIGN
	case token.SHL_ASSIGN:
		return lexer.SHL_ASSIGN
	case token.SHR_ASSIGN:
		return lexer.SHR_ASSIGN
	case token.AND_NOT_ASSIGN:
		return lexer.AND_NOT_ASSIGN
	case token.NOT:
		return lexer.NOT
	case token.EOF:
		return lexer.EOF
	default:
		return lexer.IDENT
	}
}

// FallbackParse is a helper for tests or emergency direct fallback.
func FallbackParse(filename string, src string) (*ast.File, error) {
	fset := token.NewFileSet()
	gf, err := goParser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return nil, formatParseErr(err, filename)
	}
	return gf, nil
}
