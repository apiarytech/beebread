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

// Command oscatcov reports which POUs and types of the OSCAT BASIC source
// have a port in the basic packages, by subject.
//
//	go run ./tools/oscatcov [-v]
//
// A FUNCTION must be a Go function and a FUNCTION_BLOCK a type with the
// methods INIT and Execute. Its test fails if anything is missing.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/apiarytech/beebread/tools/internal/oscat"
)

// POU is a POU of the OSCAT source.
type POU = oscat.POU

// ReadPOUs reads the POUs of an OSCAT source file.
var ReadPOUs = oscat.ReadPOUs

// Decl is a Go declaration of the basic packages.
type Decl struct {
	Package string
	Func    bool
	Methods map[string]bool
}

// ReadDecls reads the exported functions and types, with their methods, of
// the Go packages under dir.
func ReadDecls(dir string) (map[string][]*Decl, error) {
	decls := map[string][]*Decl{}
	types := map[string]*Decl{}
	fset := token.NewFileSet()
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return err
		}
		file, err := parser.ParseFile(fset, p, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		pkg := filepath.ToSlash(filepath.Dir(p))
		for _, d := range file.Decls {
			switch d := d.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil {
					if d.Name.IsExported() {
						decls[d.Name.Name] = append(decls[d.Name.Name], &Decl{Package: pkg, Func: true})
					}
					continue
				}
				recv := d.Recv.List[0].Type
				if star, ok := recv.(*ast.StarExpr); ok {
					recv = star.X
				}
				if id, ok := recv.(*ast.Ident); ok {
					key := pkg + "." + id.Name
					if types[key] == nil {
						types[key] = &Decl{Package: pkg, Methods: map[string]bool{}}
					}
					types[key].Methods[d.Name.Name] = true
				}
			case *ast.GenDecl:
				for _, s := range d.Specs {
					ts, ok := s.(*ast.TypeSpec)
					if !ok || !ts.Name.IsExported() {
						continue
					}
					key := pkg + "." + ts.Name.Name
					if types[key] == nil {
						types[key] = &Decl{Package: pkg, Methods: map[string]bool{}}
					}
					decls[ts.Name.Name] = append(decls[ts.Name.Name], types[key])
				}
			}
		}
		return nil
	})
	return decls, err
}

// Problem returns what is wrong with a POU's port, or "" if it is ported.
func Problem(p POU, decls map[string][]*Decl) string {
	d := decls[p.GoName()]
	switch {
	case len(d) == 0:
		return "missing"
	case len(d) > 1:
		pkgs := []string{}
		for _, x := range d {
			pkgs = append(pkgs, x.Package)
		}
		return "declared more than once: " + strings.Join(pkgs, ", ")
	case p.Kind == "FUNCTION" && !d[0].Func:
		return "not a function"
	case p.Kind == "FUNCTION_BLOCK" && (d[0].Func || !d[0].Methods["INIT"] || !d[0].Methods["Execute"]):
		return "not a type with INIT and Execute"
	case p.Kind == "TYPE" && d[0].Func:
		return "not a type"
	}
	return ""
}

func main() {
	verbose := flag.Bool("v", false, "list the ported POUs too")
	st := flag.String("st", "documents/oscat_basic_335.st", "the OSCAT source")
	dir := flag.String("go", "basic", "the directory of the port")
	flag.Parse()
	pous, err := ReadPOUs(*st)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	decls, err := ReadDecls(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	bySubject := map[string][]POU{}
	for _, p := range pous {
		bySubject[p.Subject] = append(bySubject[p.Subject], p)
	}
	subjects := make([]string, 0, len(bySubject))
	for s := range bySubject {
		subjects = append(subjects, s)
	}
	sort.Strings(subjects)
	total, ported := 0, 0
	for _, s := range subjects {
		var bad, good []string
		for _, p := range bySubject[s] {
			if prob := Problem(p, decls); prob != "" {
				bad = append(bad, fmt.Sprintf("%s (line %d): %s", p.Name, p.Line, prob))
			} else {
				good = append(good, p.Name)
			}
		}
		total += len(bySubject[s])
		ported += len(good)
		fmt.Printf("%-36s %3d/%3d\n", s, len(good), len(bySubject[s]))
		for _, b := range bad {
			fmt.Println("    " + b)
		}
		if *verbose {
			fmt.Println("    " + strings.Join(good, " "))
		}
	}
	fmt.Printf("%-36s %3d/%3d\n", "total", ported, total)
	if ported != total {
		os.Exit(1)
	}
}
