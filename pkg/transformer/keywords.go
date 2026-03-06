package transformer

var keywordMap = map[string]string{
	"kampung":  "package",
	"muak":     "import",
	"banana":   "func",
	"bapple":   "return",
	"pooka":    "var",
	"gelato":   "const",
	"luk":      "type",
	"kampai":   "struct",
	"buddies":  "interface",
	"papoy":    "map",
	"po":       "if",
	"ka":       "else",
	"tulaliloo":"for",
	"tikali":   "range",
	"buttom":   "break",
	"bajo":     "continue",
	"bee":      "switch",
	"doh":      "case",
	"meh":      "default",
	"underpa":  "go",
	"tatata":   "chan",
	"culo":     "select",
	"tank_yu":  "defer",
	"patalaki": "fallthrough",
	"waaah":    "goto",
	"dala":     "make",
	"pwede":    "new",
}

var builtinMap = map[string]string{
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

var goKeywordInverse = func() map[string]string {
	out := make(map[string]string, len(keywordMap))
	for bello, goName := range keywordMap {
		out[goName] = bello
	}
	return out
}()

var goBuiltinInverse = func() map[string]string {
	out := make(map[string]string, len(builtinMap))
	for bello, goName := range builtinMap {
		out[goName] = bello
	}
	return out
}()
