package lexer

import (
	"strings"
	"unicode"
)

type Lexer struct {
	input     []rune
	filename  string
	idx       int
	nextIdx   int
	ch        rune
	line      int
	column    int
	insertSemi bool
	prevType  TokenType
}

type State struct {
	idx, nextIdx int
	ch           rune
	line, column int
	insertSemi   bool
	prevType     TokenType
}

func New(filename, src string) *Lexer {
	l := &Lexer{input: []rune(src), filename: filename, line: 1, column: 0, prevType: ILLEGAL}
	l.read()
	return l
}

func (l *Lexer) Snapshot() State {
	return State{idx: l.idx, nextIdx: l.nextIdx, ch: l.ch, line: l.line, column: l.column, insertSemi: l.insertSemi, prevType: l.prevType}
}

func (l *Lexer) Restore(s State) {
	l.idx = s.idx
	l.nextIdx = s.nextIdx
	l.ch = s.ch
	l.line = s.line
	l.column = s.column
	l.insertSemi = s.insertSemi
	l.prevType = s.prevType
}

func (l *Lexer) Next() Token {
	// Consume trivia and emit implicit semicolons.
	if token := l.scanTrivia(); token != nil {
		l.prevType = SEMICOLON
		return *token
	}

	if l.ch == 0 {
		pos := l.position()
		l.prevType = EOF
		return Token{Type: EOF, Lit: "", Pos: pos, End: pos}
	}

	start := l.position()

	if isLetter(l.ch) {
		lit := l.readWhile(func(r rune) bool { return isLetter(r) || isDigit(r) || r == '_' })
		tokType, ok := KeywordLookup(lit)
		if !ok {
			tokType = IDENT
		}
		end := l.position()
		l.prevType = tokType
		l.insertSemi = CanInsertSemicolon(tokType)
		return Token{Type: tokType, Lit: lit, Pos: start, End: end}
	}

	if isDecimal(l.ch) {
		lit := l.readNumeric()
		tokType := tokenFromNumeric(lit)
		end := l.position()
		l.prevType = tokType
		l.insertSemi = CanInsertSemicolon(tokType)
		return Token{Type: tokType, Lit: lit, Pos: start, End: end}
	}

	switch l.ch {
	case '\'':
		lit := l.readRune()
		end := l.position()
		l.prevType = RUNE
		l.insertSemi = CanInsertSemicolon(RUNE)
		return Token{Type: RUNE, Lit: lit, Pos: start, End: end}
	case '`':
		lit := l.readRawString()
		end := l.position()
		l.prevType = STRING
		l.insertSemi = CanInsertSemicolon(STRING)
		return Token{Type: STRING, Lit: lit, Pos: start, End: end}
	case '"':
		lit := l.readInterpString()
		end := l.position()
		l.prevType = STRING
		l.insertSemi = CanInsertSemicolon(STRING)
		return Token{Type: STRING, Lit: lit, Pos: start, End: end}
	}

	tt, lit := l.readOperatorOrDelimiter()
	l.prevType = tt
	l.insertSemi = CanInsertSemicolon(tt)
	end := l.position()
	return Token{Type: tt, Lit: lit, Pos: start, End: end}
}

func (l *Lexer) position() Position {
	return Position{Filename: l.filename, Offset: l.idx, Line: l.line, Column: l.column}
}

func (l *Lexer) read() {
	if l.nextIdx >= len(l.input) {
		l.idx = l.nextIdx
		l.ch = 0
		l.column++
		return
	}
	l.idx = l.nextIdx
	l.ch = l.input[l.nextIdx]
	l.nextIdx++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peek() rune {
	if l.nextIdx >= len(l.input) {
		return 0
	}
	return l.input[l.nextIdx]
}

// scanTrivia consumes spaces/comments. If a newline should create a semicolon, it returns that token.
func (l *Lexer) scanTrivia() *Token {
	for {
		if l.ch == 0 {
			return nil
		}
		switch l.ch {
		case ' ', '\t', '\r', '\f':
			l.read()
			continue
		case '\n':
			if CanInsertSemicolon(l.prevType) {
				pos := l.position()
				l.read()
				l.insertSemi = false
				return &Token{Type: SEMICOLON, Lit: ";", Pos: pos, End: pos}
			}
			l.read()
			continue
		case '/':
			if l.peek() == '/' {
				nl := l.readLineComment()
				if nl && CanInsertSemicolon(l.prevType) {
					pos := l.position()
					l.insertSemi = false
					return &Token{Type: SEMICOLON, Lit: ";", Pos: pos, End: pos}
				}
				continue
			}
			if l.peek() == '*' {
				nl := l.readBlockComment()
				if nl && CanInsertSemicolon(l.prevType) {
					pos := l.position()
					l.insertSemi = false
					return &Token{Type: SEMICOLON, Lit: ";", Pos: pos, End: pos}
				}
				continue
			}
		}
		return nil
	}
}

func (l *Lexer) readLineComment() bool {
	l.read() // '/'
	l.read() // '/'
	hadNewline := false
	for l.ch != 0 && l.ch != '\n' {
		l.read()
	}
	if l.ch == '\n' {
		hadNewline = true
		l.read()
	}
	return hadNewline
}

func (l *Lexer) readBlockComment() bool {
	l.read() // '/'
	l.read() // '*'
	hadNewline := false
	for l.ch != 0 {
		if l.ch == '*' && l.peek() == '/' {
			l.read()
			l.read()
			return hadNewline
		}
		if l.ch == '\n' {
			hadNewline = true
		}
		l.read()
	}
	return hadNewline
}

func (l *Lexer) readWhile(pred func(rune) bool) string {
	start := l.idx
	for pred(l.ch) {
		l.read()
	}
	return string(l.input[start:l.idx])
}

func (l *Lexer) readNumeric() string {
	start := l.idx
	if l.ch == '0' && (l.peek() == 'x' || l.peek() == 'X' || l.peek() == 'b' || l.peek() == 'B' || l.peek() == 'o' || l.peek() == 'O') {
		l.read()
		l.read()
		for isHexDigit(l.ch) || l.ch == '_' || l.ch == '.' || l.ch == 'p' || l.ch == 'P' || l.ch == '+' || l.ch == '-' || l.ch == 'e' || l.ch == 'E' {
			l.read()
		}
		if l.ch == 'i' {
			l.read()
		}
		return string(l.input[start:l.idx])
	}
	for isDecimal(l.ch) || l.ch == '_' || l.ch == '.' || l.ch == 'e' || l.ch == 'E' || l.ch == '+' || l.ch == '-' {
		prev := l.ch
		l.read()
		if (prev == '+' || prev == '-') && !(l.prevType == INT || l.prevType == FLOAT || l.prevType == IMAGINARY || l.prevType == RUNE || l.prevType == STRING || l.prevType == IDENT || l.prevType == RBRACK || l.prevType == RPAREN) {
			break
		}
	}
	if l.ch == 'i' {
		l.read()
	}
	return string(l.input[start:l.idx])
}

func tokenFromNumeric(lit string) TokenType {
	if strings.HasSuffix(lit, "i") {
		return IMAGINARY
	}
	if strings.ContainsAny(lit, ".eEpP") || strings.Contains(lit, "0x") || strings.Contains(lit, "0X") {
		if strings.ContainsAny(lit, ".eEpP") {
			return FLOAT
		}
	}
	return INT
}

func (l *Lexer) readRune() string {
	start := l.idx
	l.read()
	for l.ch != '\'' && l.ch != 0 {
		if l.ch == '\\' {
			l.read()
			if l.ch != 0 {
				l.read()
			}
			continue
		}
		l.read()
	}
	if l.ch == '\'' {
		l.read()
	}
	return string(l.input[start:l.idx])
}

func (l *Lexer) readRawString() string {
	start := l.idx
	l.read()
	for l.ch != '`' && l.ch != 0 {
		l.read()
	}
	if l.ch == '`' {
		l.read()
	}
	return string(l.input[start:l.idx])
}

func (l *Lexer) readInterpString() string {
	start := l.idx
	l.read()
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.read()
			if l.ch != 0 {
				l.read()
			}
			continue
		}
		l.read()
	}
	if l.ch == '"' {
		l.read()
	}
	return string(l.input[start:l.idx])
}

func (l *Lexer) readOperatorOrDelimiter() (TokenType, string) {
	p := l.ch
	n := l.peek()

	match3 := func(a, b, c rune) bool { return l.ch == a && n == b && l.inputAtNext(2) == c }
	if match3('.', '.', '.') {
		l.read(); l.read(); l.read()
		return ELLIPSIS, "..."
	}

	switch p {
	case '+':
		if n == '+' {
			l.read(); l.read(); return INC, "++"
		}
		if n == '=' {
			l.read(); l.read(); return ADD_ASSIGN, "+="
		}
		l.read(); return ADD, "+"
	case '-':
		if n == '-' {
			l.read(); l.read(); return DEC, "--"
		}
		if n == '=' {
			l.read(); l.read(); return SUB_ASSIGN, "-="
		}
		l.read(); return SUB, "-"
	case '*':
		if n == '=' {
			l.read(); l.read(); return MUL_ASSIGN, "*="
		}
		l.read(); return MUL, "*"
	case '/':
		if n == '=' {
			l.read(); l.read(); return QUO_ASSIGN, "/="
		}
		l.read(); return QUO, "/"
	case '%':
		if n == '=' { l.read(); l.read(); return REM_ASSIGN, "%=" }
		l.read(); return REM, "%"
	case '&':
		if n == '&' { l.read(); l.read(); return LAND, "&&" }
		if n == '^' {
			l.read()
			if l.ch == '=' { l.read(); return AND_NOT_ASSIGN, "&^=" }
			return AND_NOT, "&^"
		}
		if n == '=' { l.read(); l.read(); return AND_ASSIGN, "&=" }
		l.read(); return AND, "&"
	case '|':
		if n == '|' { l.read(); l.read(); return LOR, "||" }
		if n == '=' { l.read(); l.read(); return OR_ASSIGN, "|=" }
		l.read(); return OR, "|"
	case '^':
		if n == '=' { l.read(); l.read(); return XOR_ASSIGN, "^=" }
		l.read(); return XOR, "^"
	case '<':
		if n == '=' { l.read(); l.read(); return LEQ, "<=" }
		if n == '<' {
			l.read(); if l.ch == '=' { l.read(); return SHL_ASSIGN, "<<=" }
			return SHL, "<<"
		}
		if n == '-' { l.read(); l.read(); return ARROW, "<-" }
		l.read(); return LSS, "<"
	case '>':
		if n == '=' { l.read(); l.read(); return GEQ, ">=" }
		if n == '>' { l.read(); if l.ch == '=' { l.read(); return SHR_ASSIGN, ">>=" }; return SHR, ">>" }
		l.read(); return GTR, ">"
	case '=':
		if n == '=' { l.read(); l.read(); return EQL, "==" }
		l.read(); return ASSIGN, "="
	case '!':
		if n == '=' { l.read(); l.read(); return NEQ, "!=" }
		l.read(); return NOT, "!"
	case ':':
		if n == '=' { l.read(); l.read(); return DEFINE, ":=" }
		l.read(); return COLON, ":"
	case '.':
		l.read(); return PERIOD, "."
	case ',': l.read(); return COMMA, ","
	case ';': l.read(); return SEMICOLON, ";"
	case '(': l.read(); return LPAREN, "("
	case ')': l.read(); return RPAREN, ")"
	case '[': l.read(); return LBRACK, "["
	case ']': l.read(); return RBRACK, "]"
	case '{': l.read(); return LBRACE, "{"
	case '}': l.read(); return RBRACE, "}"
	}

	invalid := p
	l.read()
	return ILLEGAL, string(invalid)
}

func (l *Lexer) inputAtNext(offset int) rune {
	idx := l.idx + offset
	if idx >= len(l.input) {
		return 0
	}
	return l.input[idx]
}

func isLetter(r rune) bool { return r == '_' || unicode.IsLetter(r) }
func isDecimal(r rune) bool { return r >= '0' && r <= '9' }
func isDigit(r rune) bool { return isDecimal(r) }
func isHexDigit(r rune) bool { return isDecimal(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') }
