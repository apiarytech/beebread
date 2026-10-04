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

// Command beedancegen writes the table of the beebread library that the
// beedance transpiler uses to call it: the packages, the functions with
// their parameter and result types, the function blocks with their fields,
// and the types.
//
//	go run ./tools/beedancegen > ../beedance/transpiler/beebread_library.go
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/apiarytech/beebread/tools/internal/oscat"
)

const module = "github.com/apiarytech/beebread"

// aliases are the names the generated code gives the beebread packages,
// which must not hide Go's or royaljelly's packages.
var aliases = map[string]string{
	"basic":             "oscat",
	"basic/buffer":      "oscatbuffer",
	"basic/engineering": "oscateng",
	"basic/list":        "oscatlist",
	"basic/logic":       "oscatlogic",
	"basic/math":        "oscatmath",
	"basic/other":       "oscatother",
	"basic/string":      "oscatstring",
	"basic/time_date":   "oscattime",

	"building":            "oscatbuild",
	"building/actuators":  "oscatactuators",
	"building/electrical": "oscatelectrical",
	"building/hlk":        "oscathlk",
	"building/jalousie":   "oscatjalousie",

	"network":          "oscatnet",
	"network/ads":      "oscatads",
	"network/crypto":   "oscatcrypto",
	"network/dlog":     "oscatdlog",
	"network/encoding": "oscatencoding",
	"network/file":     "oscatfile",
	"network/inet":     "oscatinet",
	"network/ip":       "oscatip",
	"network/irtrans":  "oscatirtrans",
	"network/logging":  "oscatlogging",
	"network/modbus":   "oscatmodbus",
	"network/netvar":   "oscatnetvar",
	"network/parser":   "oscatparser",
	"network/tcpip":    "oscattcpip",
	"network/telnet":   "oscattelnet",
	"network/weather":  "oscatweather",
}

// libraries are the OSCAT libraries beebread ports: the source their POUs
// are read from, and their global variables with the package that declares
// them and their types: a structured type, or the Go type of an IEC value.
var libraries = []struct {
	source  string
	globals map[string]global
}{
	{"doc/beedance_basic.st", map[string]global{
		"MATH": {"oscat", "CONSTANTS_MATH"}, "PHYS": {"oscat", "CONSTANTS_PHYS"},
		"LANGUAGE": {"oscat", "CONSTANTS_LANGUAGE"}, "SETUP": {"oscat", "CONSTANTS_SETUP"},
		"LOCATION": {"oscat", "CONSTANTS_LOCATION"}, "STRING_LENGTH": {"oscat", "iec.INT"}, "LIST_LENGTH": {"oscat", "iec.INT"},
	}},
	{"doc/beedance_building.st", nil},
	{"doc/beedance_network.st", map[string]global{
		"NETWORK_BUFFER_LONG_SIZE": {"oscatnet", "iec.UINT"}, "NETWORK_BUFFER_SHORT_SIZE": {"oscatnet", "iec.UINT"},
		"LOG_MAX": {"oscatnet", "iec.INT"}, "LOG_SIZE": {"oscatnet", "iec.INT"}, "ELEMENT_LENGTH": {"oscatnet", "iec.INT"},
		"TCP_SERVER_RESET": {"oscatnet", "iec.BYTE"}, "SSRVNETID": {"oscatnet", "iec.STRING"}, "SLOCALHOST": {"oscatnet", "iec.STRING"},
		"SYSLIBSOCKETS_OPTION": {"oscatnet", "iec.BYTE"}, "LOG_CL": {"oscatnet", "LOG_CONTROL"},
	}},
}

// global is an OSCAT global variable: the package of its port and its type.
type global struct{ pkg, typ string }

// external are the packages that port blocks OSCAT uses from outside its
// libraries, all of whose exported declarations go in the table: TwinCAT's
// TCP/IP blocks, which OSCAT NETWORK calls.
var external = map[string]bool{"network/tcpip": true}

// importNames maps the names a file imports the packages by to their
// aliases, for the file being read.
type file struct {
	pkg     string            // the package's alias
	imports map[string]string // import name to alias; "." for the dot import of basic
	own     map[string]bool   // the exported types the package declares
}

// typeString writes a type expression with the aliases of the generated
// code: a type of package basic, which the sources import with a dot or
// declare, is oscat.T.
func (f *file) typeString(e ast.Expr) string {
	switch e := e.(type) {
	case *ast.Ident:
		if f.own[e.Name] {
			return f.pkg + "." + e.Name // a type of the package itself
		}
		if ast.IsExported(e.Name) {
			// A type of basic, by its dot import.
			return "oscat." + e.Name
		}
		return e.Name
	case *ast.SelectorExpr:
		x := e.X.(*ast.Ident).Name
		if alias, ok := f.imports[x]; ok {
			return alias + "." + e.Sel.Name
		}
		return x + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + f.typeString(e.X)
	case *ast.ArrayType:
		if e.Len == nil {
			return "[]" + f.typeString(e.Elt)
		}
		n := f.constString(e.Len)
		if strings.Contains(n, "?") {
			return "?"
		}
		return "[" + n + "]" + f.typeString(e.Elt)
	}
	return "?"
}

// constString writes a constant expression, an array length such as
// NETWORK_BUFFER_LONG_SIZE + 1, with the names qualified by their packages.
func (f *file) constString(e ast.Expr) string {
	switch e := e.(type) {
	case *ast.BasicLit:
		return e.Value
	case *ast.Ident:
		if f.own[e.Name] {
			return f.pkg + "." + e.Name
		}
		if ast.IsExported(e.Name) {
			return "oscat." + e.Name
		}
	case *ast.SelectorExpr:
		x := e.X.(*ast.Ident).Name
		if alias, ok := f.imports[x]; ok {
			return alias + "." + e.Sel.Name
		}
	case *ast.ParenExpr:
		return "(" + f.constString(e.X) + ")"
	case *ast.BinaryExpr:
		return f.constString(e.X) + " " + e.Op.String() + " " + f.constString(e.Y)
	}
	return "?"
}

type function struct {
	pkg     string
	params  []string
	results []string
}

type field struct {
	name, typ string
}

type block struct {
	pkg     string
	fields  []field
	methods map[string]bool
}

func main() {
	// Only the ports of OSCAT's POUs and types go in the table; the helpers
	// of the packages, such as basic's IEC standard functions, do not.
	isPOU := map[string]bool{}
	inputs := map[string][]string{}
	bounds := map[string]map[string][]int64{}
	oscatName := map[string]string{}
	globals := map[string]global{}
	for _, lib := range libraries {
		pous, err := oscat.ReadCleanedPOUs(lib.source)
		check(err)
		for _, p := range pous {
			isPOU[p.GoName()] = true
			oscatName[p.GoName()] = p.Name
		}
		in, err := oscat.Inputs(lib.source)
		check(err)
		maps.Copy(inputs, in)
		b, err := oscat.ArrayBounds(lib.source)
		check(err)
		maps.Copy(bounds, b)
		maps.Copy(globals, lib.globals)
	}
	fset := token.NewFileSet()
	// The exported types each package declares, so that a type named
	// without a package is the package's own or else basic's.
	own := map[string]map[string]bool{}
	for dir, alias := range aliases {
		own[alias] = map[string]bool{}
		files, err := filepath.Glob(filepath.Join(dir, "*.go"))
		check(err)
		for _, path := range files {
			if strings.HasSuffix(path, "_test.go") {
				continue
			}
			src, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
			check(err)
			for _, d := range src.Decls {
				if g, ok := d.(*ast.GenDecl); ok {
					for _, s := range g.Specs {
						if vs, ok := s.(*ast.ValueSpec); ok && g.Tok == token.CONST {
							for _, n := range vs.Names {
								if n.IsExported() {
									own[alias][n.Name] = true // a constant array lengths may name
								}
							}
						}
						if ts, ok := s.(*ast.TypeSpec); ok && ts.Name.IsExported() {
							own[alias][ts.Name.Name] = true
							if external[dir] {
								isPOU[ts.Name.Name] = true
							}
						}
					}
				}
				if fd, ok := d.(*ast.FuncDecl); ok && external[dir] && fd.Recv == nil && fd.Name.IsExported() {
					isPOU[fd.Name.Name] = true
				}
			}
		}
	}
	if own["oscat"] != nil {
		own["oscat"] = nil // basic's types are oscat.T anyway
	}
	functions := map[string]function{}
	blocks := map[string]*block{}
	var order []string // blocks in the order they are declared
	for dir, alias := range aliases {
		files, err := filepath.Glob(filepath.Join(dir, "*.go"))
		check(err)
		for _, path := range files {
			if strings.HasSuffix(path, "_test.go") {
				continue
			}
			src, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
			check(err)
			f := &file{pkg: alias, imports: map[string]string{"iec": "iec"}, own: own[alias]}
			for _, imp := range src.Imports {
				p := strings.Trim(imp.Path.Value, `"`)
				rel := strings.TrimPrefix(p, module+"/")
				a, ok := aliases[rel]
				if !ok {
					continue
				}
				name := rel[strings.LastIndex(rel, "/")+1:]
				if imp.Name != nil {
					name = imp.Name.Name
				}
				f.imports[name] = a
			}
			for _, d := range src.Decls {
				switch d := d.(type) {
				case *ast.FuncDecl:
					if !d.Name.IsExported() {
						continue
					}
					if d.Recv != nil {
						recv := d.Recv.List[0].Type
						if star, ok := recv.(*ast.StarExpr); ok {
							recv = star.X
						}
						name := recv.(*ast.Ident).Name
						if blocks[name] == nil {
							blocks[name] = &block{pkg: alias, methods: map[string]bool{}}
						}
						blocks[name].methods[d.Name.Name] = true
						continue
					}
					if d.Type.TypeParams != nil {
						continue // a generic helper
					}
					fn := function{pkg: alias}
					for _, p := range d.Type.Params.List {
						n := len(p.Names)
						if n == 0 {
							n = 1
						}
						for i := 0; i < n; i++ {
							fn.params = append(fn.params, f.typeString(p.Type))
						}
					}
					if d.Type.Results != nil {
						for _, r := range d.Type.Results.List {
							fn.results = append(fn.results, f.typeString(r.Type))
						}
					}
					if old, ok := functions[d.Name.Name]; ok && old.pkg != fn.pkg && isPOU[d.Name.Name] {
						fmt.Fprintf(os.Stderr, "function %s in both %s and %s\n", d.Name.Name, old.pkg, fn.pkg)
					}
					functions[d.Name.Name] = fn
				case *ast.GenDecl:
					for _, s := range d.Specs {
						ts, ok := s.(*ast.TypeSpec)
						if !ok || !ts.Name.IsExported() {
							continue
						}
						st, ok := ts.Type.(*ast.StructType)
						if !ok {
							continue
						}
						name := ts.Name.Name
						if b := blocks[name]; b != nil && len(b.fields) > 0 && b.pkg != alias && isPOU[name] {
							fmt.Fprintf(os.Stderr, "type %s in both %s and %s\n", name, b.pkg, alias)
						}
						if blocks[name] == nil {
							blocks[name] = &block{pkg: alias, methods: map[string]bool{}}
						}
						blocks[name].pkg = alias
						order = append(order, name)
						for _, fl := range st.Fields.List {
							if len(fl.Names) == 0 {
								// An embedded struct of the package; its
								// fields are the block's.
								blocks[name].fields = append(blocks[name].fields, embedded(src, fl.Type.(*ast.Ident).Name, f)...)
								continue
							}
							for _, n := range fl.Names {
								if n.IsExported() {
									blocks[name].fields = append(blocks[name].fields, field{n.Name, f.typeString(fl.Type)})
								}
							}
						}
					}
				}
			}
		}
	}

	var out bytes.Buffer
	fmt.Fprintf(&out, `// Code generated by beebread's tools/beedancegen; DO NOT EDIT.

package transpiler

// beebreadModule is the module of beebread, the OSCAT libraries in Go.
const beebreadModule = %q

// beebreadPackages maps the names the generated code uses for the beebread
// packages to their import paths.
var beebreadPackages = map[string]string{
`, module)
	pkgs := make([]string, 0, len(aliases))
	for dir := range aliases {
		pkgs = append(pkgs, dir)
	}
	sort.Strings(pkgs)
	for _, dir := range pkgs {
		fmt.Fprintf(&out, "\t%q: beebreadModule + %q,\n", aliases[dir], "/"+dir)
	}
	fmt.Fprintf(&out, `}

// beebreadFunctions describes the functions of the OSCAT libraries: each parameter's
// Go type and the result's. A function with more than one result, which no
// OSCAT function has, is left out.
var beebreadFunctions = map[string]stdFunction{
`)
	names := make([]string, 0, len(functions))
	for n := range functions {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fn := functions[n]
		if !isPOU[n] || len(fn.results) != 1 || strings.Contains(strings.Join(fn.params, " ")+fn.results[0], "?") || strings.Contains(fn.results[0], "time.") {
			continue
		}
		fmt.Fprintf(&out, "\t%q: {pkg: %q, params: %#v, result: %q},\n", n, fn.pkg, quote(fn.params), fn.results[0])
	}
	fmt.Fprintf(&out, `}

// beebreadStandard describes the IEC 61131-3 standard functions of
// beebread's package basic, which OSCAT's sources use, with the behaviour
// OSCAT was written for; royaljelly's come first where it has them.
var beebreadStandard = map[string]stdFunction{
`)
	for _, n := range names {
		fn := functions[n]
		if isPOU[n] || fn.pkg != "oscat" || n == "CHARS" || n == "STR" || len(fn.results) != 1 || strings.Contains(strings.Join(fn.params, " ")+fn.results[0], "?") ||
			strings.Contains(strings.Join(fn.params, " ")+fn.results[0], "time.") {
			continue
		}
		fmt.Fprintf(&out, "\t%q: {pkg: %q, params: %#v, result: %q},\n", n, fn.pkg, quote(fn.params), fn.results[0])
	}
	fmt.Fprintf(&out, `}

// beebreadDefaults are the values of the parameters of the functions of
// the OSCAT libraries that a call leaves out: their initial values, or the zero
// value of their type.
var beebreadDefaults = map[string][]string{
`)
	for _, n := range names {
		fn := functions[n]
		if !isPOU[n] || len(fn.results) != 1 {
			continue
		}
		lits := inputs[oscatName[n]]
		if len(lits) != len(fn.params) {
			fmt.Fprintf(os.Stderr, "%s: %d inputs in OSCAT, %d parameters\n", n, len(lits), len(fn.params))
			lits = make([]string, len(fn.params))
		}
		vals := make([]string, len(fn.params))
		for i, typ := range fn.params {
			vals[i] = goValue(typ, lits[i])
		}
		fmt.Fprintf(&out, "\t%q: %#v,\n", strings.ToUpper(n), quote(vals))
	}
	fmt.Fprintf(&out, `}

// beebreadFunctionBlock is a function block or a structured type of an
// OSCAT library: its Go type and its fields' Go types. A function block has the
// methods INIT and Execute(now); a VAR_IN_OUT is a pointer field. lows are
// the lower bounds OSCAT declares for the dimensions of the fields that
// are arrays and do not start at 0; the Go arrays start at 0.
type beebreadFunctionBlock struct {
	goType string
	fields map[string]string
	isFB   bool
	lows   map[string][]int64
}

// beebreadFunctionBlocks describes the function blocks and structured types
// of the OSCAT libraries, by upper case name.
var beebreadFunctionBlocks = map[string]beebreadFunctionBlock{
`)
	sort.Strings(order)
	for _, n := range order {
		if !isPOU[n] {
			continue
		}
		b := blocks[n]
		fmt.Fprintf(&out, "\t%q: {goType: %q, isFB: %v, fields: map[string]string{", strings.ToUpper(n), b.pkg+"."+n, b.methods["INIT"] && b.methods["Execute"])
		for i, fl := range b.fields {
			if i > 0 {
				out.WriteString(", ")
			}
			fmt.Fprintf(&out, "%q: %q", fl.name, fl.typ)
		}
		out.WriteString("}")
		var lows []string
		for _, fl := range b.fields {
			l := bounds[oscatName[n]][fl.name]
			nonZero := false
			for _, v := range l {
				nonZero = nonZero || v != 0
			}
			if nonZero {
				lows = append(lows, fmt.Sprintf("%q: %#v", fl.name, l))
			}
		}
		if len(lows) > 0 {
			fmt.Fprintf(&out, ", lows: map[string][]int64{%s}", strings.Join(lows, ", "))
		}
		out.WriteString("},\n")
	}
	fmt.Fprintf(&out, `}

// beebreadGlobal is a global variable of an OSCAT library: the package of
// its port and its type: a structured type's name, or the Go type of an IEC value.
type beebreadGlobal struct{ pkg, typ string }

// beebreadGlobals are the global variables of the OSCAT libraries.
var beebreadGlobals = map[string]beebreadGlobal{
`)
	gnames := slices.Sorted(maps.Keys(globals))
	for _, n := range gnames {
		fmt.Fprintf(&out, "\t%q: {%q, %q},\n", n, globals[n].pkg, globals[n].typ)
	}
	out.WriteString("}\n")
	src, err := format.Source(out.Bytes())
	check(err)
	os.Stdout.Write(src)
}

// embedded returns the exported fields of the unexported struct name
// declared in src.
func embedded(src *ast.File, name string, f *file) []field {
	var out []field
	ast.Inspect(src, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != name {
			return true
		}
		for _, fl := range ts.Type.(*ast.StructType).Fields.List {
			for _, n := range fl.Names {
				if n.IsExported() {
					out = append(out, field{n.Name, f.typeString(fl.Type)})
				}
			}
		}
		return false
	})
	return out
}

func quote(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// goValue returns the Go expression of the ST literal lit for the Go type
// typ, or the type's zero value if lit is "".
func goValue(typ, lit string) string {
	lit = strings.TrimSpace(lit)
	switch {
	case strings.HasPrefix(typ, "[]"), strings.HasPrefix(typ, "*"):
		return "nil"
	case strings.HasPrefix(typ, "["), strings.HasPrefix(typ, "oscat."),
		typ == "iec.DT", typ == "iec.DATE", typ == "iec.TOD":
		return typ + "{}"
	case typ == "iec.BOOL":
		switch strings.ToUpper(lit) {
		case "TRUE", "1":
			return "iec.BOOL(true)"
		}
		return "iec.BOOL(false)"
	case typ == "iec.STRING":
		return fmt.Sprintf("iec.STRING(%q)", strings.Trim(lit, "'"))
	case typ == "iec.TIME":
		if lit == "" {
			return "iec.TIME(0)"
		}
		return fmt.Sprintf("iec.TIME(%d)", stTime(lit))
	}
	if lit == "" {
		return typ + "(0)"
	}
	lit = strings.ReplaceAll(lit, "_", "")
	if i := strings.Index(lit, "#"); i >= 0 {
		// A typed literal, such as INT#2, or a radix, such as 16#FF.
		base, digits := lit[:i], lit[i+1:]
		switch base {
		case "2":
			lit = "0b" + digits
		case "8":
			lit = "0o" + digits
		case "16":
			lit = "0x" + digits
		default:
			lit = digits
		}
	}
	return typ + "(" + lit + ")"
}

// stTime returns the nanoseconds of an ST TIME literal, such as T#1.2s or
// t#10d.
func stTime(lit string) int64 {
	s := strings.ToLower(lit[strings.Index(lit, "#")+1:])
	units := []struct {
		name string
		ns   float64
	}{{"ms", 1e6}, {"d", 864e11}, {"h", 36e11}, {"m", 6e10}, {"s", 1e9}}
	var total float64
	for s != "" {
		i := 0
		for i < len(s) && (s[i] >= '0' && s[i] <= '9' || s[i] == '.') {
			i++
		}
		var v float64
		fmt.Sscan(s[:i], &v)
		s = s[i:]
		for _, u := range units {
			if strings.HasPrefix(s, u.name) {
				total += v * u.ns
				s = s[len(u.name):]
				break
			}
		}
		s = strings.TrimLeft(s, "_")
	}
	return int64(total)
}
