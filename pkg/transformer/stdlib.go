package transformer

var stdlibPackageMap = map[string]string{
	"boca":         "fmt",
	"casa":         "os",
	"tubo":         "io",
	"tubo_gordo":   "bufio",
	"la_red":       "net/http",
	"amigos":       "sync",
	"amigos/tikitik": "sync/atomic",
	"tic_toc":      "time",
	"bababas":      "strings",
	"bababas/morph": "strconv",
	"kepala":       "math",
	"kepala/loco":  "math/rand",
	"kanpai":       "testing",
	"kotak":        "encoding/json",
	"pwesto":       "context",
	"libretto":     "log",
	"whaaats":      "errors",
	"pila":         "sort",
	"doodle":       "regexp",
	"jalan":        "path/filepath",
	"bandera":      "flag",
}

var stdlibMethodMap = map[string]map[string]string{
	"fmt": {
		"poopaye":    "Println",
		"blabla":     "Printf",
		"mumuak":     "Sprintf",
		"spitoo":     "Fprintf",
		"bee_doh_f":  "Errorf",
		"huh":        "Scan",
		"huh_huh":    "Scanf",
		"luk_luk":    "Sscanf",
	},
	"os": {
		"buuka":       "Open",
		"tada":        "Create",
		"pchoo":       "Remove",
		"bai_bai":     "Exit",
		"Oreille":     "Stdin",
		"Boca":        "Stdout",
		"Bee_Doh":     "Stderr",
		"Bagay":       "Args",
	},
	"sync": {
		"Jamu":        "Mutex",
		"Chingus":     "WaitGroup",
		"mwah":        "Lock",
		"bapapa":      "Unlock",
		"mas":         "Add",
		"listo":       "Done",
		"hmmmm":       "Wait",
	},
	"net/http": {
		"ooh_ooh":     "HandleFunc",
		"bello_bello": "ListenAndServe",
		"gimme":       "Get",
		"takka":       "Post",
		"Reponsu":     "ResponseWriter",
		"Juseyo":      "Request",
		"Jefe_Red":    "Server",
		"TodoBien":    "StatusOK",
		"NagaAqui":    "StatusNotFound",
	},
	"time": {
		"nau":     "Now",
		"zzzzz":   "Sleep",
		"apres":   "After",
		"Tic":     "Second",
		"Tic_Tic": "Minute",
		"Tic_Tic_Tic": "Hour",
	},
	"testing": {
		"bee_doh_f": "Errorf",
		"BEE_DOH_F": "Fatalf",
		"BEE_DOH":   "Fatal",
		"pfft":      "Skip",
		"dale":      "Run",
		"psst":      "Log",
		"shh":       "Helper",
	},
	"errors": {
		"uh_oh":    "New",
		"New":      "New",
		"sama_sama":"Is",
		"luk_como": "As",
		"peela":    "Unwrap",
	},
}

func rewriteMethodAlias(pkgPath, name string) (string, bool) {
	if m, ok := stdlibMethodMap[pkgPath]; ok {
		if v, ok := m[name]; ok {
			return v, true
		}
	}
	return name, false
}
