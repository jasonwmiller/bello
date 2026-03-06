package module

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Require struct {
	Module  string
	Version string
}

type Replace struct {
	OldMod  string
	OldVer  string
	NewMod  string
	NewVer  string
}

type ModuleFile struct {
	ModulePath string
	GoVersion  string
	Requires   []Require
	Replaces   []Replace
}

func Parse(path string) (*ModuleFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := &ModuleFile{}
	var mode string
	sc := bufio.NewScanner(f)
	lineNo := 0

	for sc.Scan() {
		lineNo++
		line := stripComment(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		if len(parts) == 1 && parts[0] == ")" {
			if mode == "" {
				return nil, fmt.Errorf("unexpected ) at line %d", lineNo)
			}
			mode = ""
			continue
		}

		switch mode {
		case "necesita":
			if req := parseRequire(parts); req != nil {
				m.Requires = append(m.Requires, *req)
				continue
			}
			if rep := parseReplace(parts); rep != nil {
				m.Replaces = append(m.Replaces, *rep)
				continue
			}
			return nil, fmt.Errorf("invalid necesita block entry at line %d", lineNo)
		case "cambio":
			if rep := parseReplace(parts); rep == nil {
				return nil, fmt.Errorf("invalid cambio block entry at line %d", lineNo)
			} else {
				m.Replaces = append(m.Replaces, *rep)
				continue
			}
		}

		switch parts[0] {
		case "modulo":
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid modulo directive at line %d", lineNo)
			}
			m.ModulePath = parts[1]
		case "bello":
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid bello directive at line %d", lineNo)
			}
			m.GoVersion = parts[1]
		case "necesita":
			// single-line require/replacement
			if len(parts) == 1 {
				mode = "necesita"
				continue
			}
			if len(parts) == 2 && parts[1] == "(" {
				mode = "necesita"
				continue
			}
			if len(parts) == 2 {
				// allow `necesita` grouped-open as separate tokenization fallback
				continue
			}
			if len(parts) == 3 {
				m.Requires = append(m.Requires, Require{Module: parts[1], Version: parts[2]})
				continue
			}
			if rep := parseReplace(parts[1:]); rep != nil {
				m.Replaces = append(m.Replaces, *rep)
				continue
			}
			if len(parts) > 1 && parts[1] == "(" {
				mode = "necesita"
				continue
			}
			return nil, fmt.Errorf("invalid necesita directive at line %d", lineNo)
	case "cambio":
		if len(parts) == 1 {
			mode = "cambio"
			continue
		}
		if len(parts) == 2 && parts[1] == "(" {
			mode = "cambio"
			continue
		}
		rep := parseReplace(parts[1:])
		if rep != nil {
			m.Replaces = append(m.Replaces, *rep)
			continue
		}
		return nil, fmt.Errorf("invalid cambio directive at line %d", lineNo)
	}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if m.ModulePath == "" {
		return nil, fmt.Errorf("missing modulo directive")
	}
	if m.GoVersion == "" {
		m.GoVersion = "1.23"
	}
	return m, nil
}

func parseRequire(parts []string) *Require {
	if len(parts) != 2 {
		return nil
	}
	return &Require{Module: parts[0], Version: parts[1]}
}

func parseReplace(parts []string) *Replace {
	if len(parts) == 3 && parts[1] == "=>" {
		return &Replace{OldMod: parts[0], NewMod: parts[2]}
	}
	if len(parts) == 4 && parts[2] == "=>" {
		return &Replace{OldMod: parts[0], OldVer: parts[1], NewMod: parts[3]}
	}
	if len(parts) == 5 && parts[2] == "=>" {
		return &Replace{OldMod: parts[0], OldVer: parts[1], NewMod: parts[3], NewVer: parts[4]}
	}
	return nil
}

func stripComment(line string) string {
	if i := strings.Index(line, "//"); i >= 0 {
		line = line[:i]
	}
	return strings.TrimSpace(strings.TrimSuffix(line, ";"))
}

func (m *ModuleFile) RenderGoMod() string {
	var b strings.Builder
	b.WriteString("module ")
	b.WriteString(m.ModulePath)
	b.WriteByte('\n')

	if m.GoVersion != "" {
		b.WriteString("go ")
		b.WriteString(m.GoVersion)
		b.WriteByte('\n')
	}

	if len(m.Requires) > 0 {
		b.WriteByte('\n')
		if len(m.Requires) == 1 {
			r := m.Requires[0]
			b.WriteString("require ")
			b.WriteString(r.Module)
			if r.Version != "" {
				b.WriteByte(' ')
				b.WriteString(r.Version)
			}
			b.WriteByte('\n')
		} else {
			b.WriteString("require (\n")
			for _, r := range m.Requires {
				b.WriteByte('\t')
				b.WriteString(r.Module)
				if r.Version != "" {
					b.WriteByte(' ')
					b.WriteString(r.Version)
				}
				b.WriteByte('\n')
			}
			b.WriteString(")\n")
		}
	}

	if len(m.Replaces) > 0 {
		b.WriteByte('\n')
		if len(m.Replaces) == 1 {
			r := m.Replaces[0]
			b.WriteString("replace ")
			b.WriteString(r.OldMod)
			if r.OldVer != "" {
				b.WriteByte(' ')
				b.WriteString(r.OldVer)
			}
			b.WriteByte(' ')
			b.WriteString("=> ")
			b.WriteString(r.NewMod)
			if r.NewVer != "" {
				b.WriteByte(' ')
				b.WriteString(r.NewVer)
			}
			b.WriteByte('\n')
		} else {
			b.WriteString("replace (\n")
			for _, r := range m.Replaces {
				b.WriteByte('\t')
				b.WriteString(r.OldMod)
				if r.OldVer != "" {
					b.WriteByte(' ')
					b.WriteString(r.OldVer)
				}
				b.WriteByte(' ')
				b.WriteString("=> ")
				b.WriteString(r.NewMod)
				if r.NewVer != "" {
					b.WriteByte(' ')
					b.WriteString(r.NewVer)
				}
				b.WriteByte('\n')
			}
			b.WriteString(")\n")
		}
	}
	return b.String()
}

func ModuleNameFromPath(path string) string {
	name := filepath.Base(path)
	name = strings.TrimSpace(strings.ReplaceAll(name, " ", "-"))
	if name == "" || name == "." {
		return "bello.local"
	}
	return "example.com/" + name
}
