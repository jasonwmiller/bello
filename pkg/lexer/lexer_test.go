package lexer

import (
	"reflect"
	"testing"
)

func collectTypes(tokens []Token) []TokenType {
	out := make([]TokenType, 0, len(tokens))
	for _, tok := range tokens {
		out = append(out, tok.Type)
	}
	return out
}

func TestKeywordLookup(t *testing.T) {
	tests := map[string]TokenType{
		"kampung":  KAMPUNG,
		"muak":     MUAK,
		"banana":   BANANA,
		"bapple":   BAPPLE,
		"pooka":    POOKA,
		"gelato":   GELATO,
		"luk":      LUK,
		"kampai":   KAMPAI,
		"buddies":  BUDDIES,
		"papoy":    PAPOY,
		"po":       PO,
		"ka":       KA,
		"tulaliloo": TULALILOO,
		"tikali":   TIKALI,
		"buttom":   BUTTOM,
		"bajo":     BAJO,
		"bee":      BEE,
		"doh":      DOH,
		"meh":      MEH,
		"underpa":  UNDERPA,
		"tatata":   TATATA,
		"culo":     CULO,
		"tank_yu":  TANK_YU,
		"patalaki": PATALAKI,
		"waaah":    WAAAH,
		"dala":     DALA,
		"pwede":    PWEDE,
	}

	for lit, expected := range tests {
		tok, ok := KeywordLookup(lit)
		if !ok {
			t.Fatalf("keyword %q was not found", lit)
		}
		if tok != expected {
			t.Fatalf("keyword %q mapped to %v, expected %v", lit, tok, expected)
		}
	}
}

func TestSemicolonInsertionAndComments(t *testing.T) {
	src := `poopaye("hi") // inline comment
poopaye("bye")`
	l := New("inline.🍌", src)
	got := collectTypes(nextTypes(l))
	want := []TokenType{IDENT, LPAREN, STRING, RPAREN, SEMICOLON, IDENT, LPAREN, STRING, RPAREN, EOF}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tokens mismatch: got=%v want=%v", got, want)
	}
}

func TestArrowAndRuneEdgeTokens(t *testing.T) {
	src := "tatata<- x\n<-tatata y\n\"hola 🥴\""
	l := New("edge.🍌", src)
	got := collectTypes(nextTypes(l))
	want := []TokenType{TATATA, ARROW, IDENT, SEMICOLON, ARROW, TATATA, IDENT, SEMICOLON, STRING, EOF}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tokens mismatch: got=%v want=%v", got, want)
	}
}

func nextTypes(l *Lexer) []Token {
	var out []Token
	for {
		tok := l.Next()
		out = append(out, tok)
		if tok.Type == EOF || tok.Type == ILLEGAL {
			break
		}
	}
	return out
}
