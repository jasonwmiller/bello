package emitter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasonwmiller/bello/pkg/lexer"
	"github.com/jasonwmiller/bello/pkg/transformer"
)

// Emit writes the Go AST into the provided output directory.
func Emit(gf *ast.File, filename string, outDir string) (string, *transformer.PositionMap, error) {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, gf); err != nil {
		return "", nil, err
	}

	targetDir := outDir
	if targetDir == "" {
		var err error
		targetDir, err = os.MkdirTemp("", "bello-build-")
		if err != nil {
			return "", nil, err
		}
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", nil, err
	}

	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	if base == "" {
		base = "main"
	}
	outPath := filepath.Join(targetDir, base+".go")

	if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
		return "", nil, err
	}

	return outPath, &transformer.PositionMap{
		GoFile:     outPath,
		BelloFile:  filename,
		LineOffset: 0,
	}, nil
}

// EmitDefault writes the Go AST into a temporary directory.
func EmitDefault(gf *ast.File, filename string) (string, *transformer.PositionMap, error) {
	return Emit(gf, filename, "")
}

func RemapError(pm *transformer.PositionMap, goFile string, line, col int) lexer.Position {
	if pm == nil || goFile != pm.GoFile {
		return lexer.Position{Filename: goFile, Line: line, Column: col}
	}
	return pm.Remap(goFile, line, col)
}

func FormatBelloError(goFilename string, errLine, errCol int, pm *transformer.PositionMap, msg string) string {
	pos := RemapError(pm, goFilename, errLine, errCol)
	return fmt.Sprintf("BEE DOH! %s:%d:%d — %s", pos.Filename, pos.Line, pos.Column, msg)
}
