/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - EPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package building

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestCoverage checks that every POU of doc/beedance_building.st, those
// commented out for beedance too, has a port: a function or a type of the
// same name in one of the packages.
func TestCoverage(t *testing.T) {
	src, err := os.ReadFile("../doc/beedance_building.st")
	if err != nil {
		t.Fatal(err)
	}
	rePOU := regexp.MustCompile(`(?m)^(?:// )?(FUNCTION_BLOCK|FUNCTION)\s+([A-Za-z_0-9]+)`)
	pous := rePOU.FindAllStringSubmatch(string(src), -1)
	if len(pous) == 0 {
		t.Fatal("no POU found")
	}

	ported := map[string]string{}
	files, err := filepath.Glob("*/*.go")
	if err != nil {
		t.Fatal(err)
	}
	files = append(files, "building.go")
	fset := token.NewFileSet()
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil {
					ported[d.Name.Name] = "FUNCTION"
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						ported[ts.Name.Name] = "FUNCTION_BLOCK"
					}
				}
			}
		}
	}
	for _, p := range pous {
		kind, name := p[1], strings.ToUpper(p[2])
		if got, ok := ported[name]; !ok {
			t.Errorf("%s %s has no port", kind, name)
		} else if got != kind {
			t.Errorf("%s %s is ported as a %s", kind, name, got)
		}
	}
	t.Logf("%d POUs ported", len(pous))
}
