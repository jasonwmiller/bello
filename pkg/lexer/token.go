package lexer

import "fmt"

// Position represents a location in a Bello source file.
type Position struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}

// TokenType is a lexical token type.
type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL
	COMMENT

	// Identifiers and literals
	IDENT
	INT
	FLOAT
	IMAGINARY
	RUNE
	STRING

	// Keywords
	KAMPUNG
	MUAK
	BANANA
	BAPPLE
	POOKA
	GELATO
	LUK
	KAMPAI
	BUDDIES
	PAPOY
	PO
	KA
	TULALILOO
	TIKALI
	BUTTOM
	BAJO
	BEE
	DOH
	MEH
	UNDERPA
	TATATA
	CULO
	TANK_YU
	PATALAKI
	WAAAH
	DALA
	PWEDE

	// Operators
	ADD
	SUB
	MUL
	QUO
	REM
	AND
	OR
	XOR
	SHL
	SHR
	AND_NOT
	LAND
	LOR
	ARROW
	INC
	DEC
	EQL
	NEQ
	LSS
	LEQ
	GTR
	GEQ
	ASSIGN
	DEFINE
	ADD_ASSIGN
	SUB_ASSIGN
	MUL_ASSIGN
	QUO_ASSIGN
	REM_ASSIGN
	AND_ASSIGN
	OR_ASSIGN
	XOR_ASSIGN
	SHL_ASSIGN
	SHR_ASSIGN
	AND_NOT_ASSIGN
	NOT

	// Delimiters
	LPAREN
	RPAREN
	LBRACK
	RBRACK
	LBRACE
	RBRACE
	COMMA
	PERIOD
	SEMICOLON
	COLON
	ELLIPSIS
)

var tokenTypeNames = map[TokenType]string{
	EOF:            "EOF",
	ILLEGAL:        "ILLEGAL",
	COMMENT:        "COMMENT",
	IDENT:          "IDENT",
	INT:            "INT",
	FLOAT:          "FLOAT",
	IMAGINARY:      "IMAGINARY",
	RUNE:           "RUNE",
	STRING:         "STRING",
	KAMPUNG:        "kampung",
	MUAK:           "muak",
	BANANA:         "banana",
	BAPPLE:         "bapple",
	POOKA:          "pooka",
	GELATO:         "gelato",
	LUK:            "luk",
	KAMPAI:         "kampai",
	BUDDIES:        "buddies",
	PAPOY:          "papoy",
	PO:             "po",
	KA:             "ka",
	TULALILOO:      "tulaliloo",
	TIKALI:         "tikali",
	BUTTOM:         "buttom",
	BAJO:           "bajo",
	BEE:            "bee",
	DOH:            "doh",
	MEH:            "meh",
	UNDERPA:        "underpa",
	TATATA:         "tatata",
	CULO:           "culo",
	TANK_YU:        "tank_yu",
	PATALAKI:       "patalaki",
	WAAAH:          "waaah",
	DALA:           "dala",
	PWEDE:          "pwede",
	ADD:            "+",
	SUB:            "-",
	MUL:            "*",
	QUO:            "/",
	REM:            "%",
	AND:            "&",
	OR:             "|",
	XOR:            "^",
	SHL:            "<<",
	SHR:            ">>",
	AND_NOT:        "&^",
	LAND:           "&&",
	LOR:            "||",
	ARROW:          "<-",
	INC:            "++",
	DEC:            "--",
	EQL:            "==",
	NEQ:            "!=",
	LSS:            "<",
	LEQ:            "<=",
	GTR:            ">",
	GEQ:            ">=",
	ASSIGN:         "=",
	DEFINE:         ":=",
	ADD_ASSIGN:     "+=",
	SUB_ASSIGN:     "-=",
	MUL_ASSIGN:     "*=",
	QUO_ASSIGN:     "/=",
	REM_ASSIGN:     "%=",
	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",
	NOT:            "!",
	LPAREN:         "(",
	RPAREN:         ")",
	LBRACK:         "[",
	RBRACK:         "]",
	LBRACE:         "{",
	RBRACE:         "}",
	COMMA:          ",",
	PERIOD:         ".",
	SEMICOLON:      ";",
	COLON:          ":",
	ELLIPSIS:       "...",
}

func (tt TokenType) String() string {
	if s, ok := tokenTypeNames[tt]; ok {
		return s
	}
	return fmt.Sprintf("TokenType(%d)", tt)
}

// Token carries metadata and literal text.
type Token struct {
	Type TokenType
	Lit  string
	Pos  Position
	End  Position
}

var keywords = map[string]TokenType{
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

func KeywordLookup(lit string) (TokenType, bool) {
	t, ok := keywords[lit]
	return t, ok
}

func IsKeyword(tok TokenType) bool {
	switch tok {
	case KAMPUNG, MUAK, BANANA, BAPPLE, POOKA, GELATO, LUK, KAMPAI,
		BUDDIES, PAPOY, PO, KA, TULALILOO, TIKALI, BUTTOM, BAJO,
		BEE, DOH, MEH, UNDERPA, TATATA, CULO, TANK_YU, PATALAKI,
		WAAAH, DALA, PWEDE:
		return true
	default:
		return false
	}
}

func CanInsertSemicolon(tt TokenType) bool {
	switch tt {
	case IDENT, INT, FLOAT, IMAGINARY, RUNE, STRING,
		BAPPLE, BUTTOM, BAJO, PATALAKI,
		INC, DEC, RPAREN, RBRACE, RBRACK:
		return true
	default:
		return false
	}
}

