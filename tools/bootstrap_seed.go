package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasonwmiller/bello/pkg/transformer"
)

func main() {
	sourceRoot := flag.String("source", ".", "Go source root for cmd/ and pkg/")
	outRoot := flag.String("out", "", "Seed output root (defaults to <source>/bootstrap/src)")
	flag.Parse()

	resolvedSource, err := filepath.Abs(*sourceRoot)
	if err != nil {
		panic(fmt.Sprintf("cannot resolve source root: %v", err))
	}

	out := *outRoot
	if out == "" {
		out = filepath.Join(resolvedSource, "bootstrap", "src")
	}
	resolvedOut, err := filepath.Abs(out)
	if err != nil {
		panic(fmt.Sprintf("cannot resolve output root: %v", err))
	}

	if err := os.RemoveAll(resolvedOut); err != nil {
		panic(fmt.Sprintf("clear output root: %v", err))
	}
	if err := os.MkdirAll(resolvedOut, 0o755); err != nil {
		panic(fmt.Sprintf("make output root: %v", err))
	}

	if err := filepath.WalkDir(resolvedSource, func(path string, d os.DirEntry, errIn error) error {
		if errIn != nil {
			return errIn
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if base == ".git" || base == ".jj" || base == ".tools" || strings.HasPrefix(base, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		rel, err := filepath.Rel(resolvedSource, path)
		if err != nil {
			return err
		}
		parts := strings.Split(filepath.ToSlash(rel), "/")
		if len(parts) < 2 || (parts[0] != "cmd" && parts[0] != "pkg") {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		outBello, err := transformer.RewriteGoToBelloSource(string(src))
		if err != nil {
			return fmt.Errorf("cannot translate %s: %w", path, err)
		}

		dst := filepath.Join(resolvedOut, strings.TrimSuffix(rel, ".go")+".🍌")
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return os.WriteFile(dst, outBello, 0o644)
	}); err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(resolvedOut, "go.mod"), []byte("module github.com/jasonwmiller/bello\n\ngo 1.23\n"), 0o644); err != nil {
		panic(fmt.Sprintf("write module: %v", err))
	}

	fmt.Println("bello bootstrap seed refreshed at", resolvedOut)
}
