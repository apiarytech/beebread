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

package other

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	gomath "math"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	"unsafe"
)

// initer is a function block: INIT sets it to its initial values.
type initer interface{ INIT() }

// initBlocks are the function blocks of the package.
var initBlocks = map[string]func() initer{
	"ESR_COLLECT": func() initer { return new(ESR_COLLECT) },
	"ESR_MON_B8":  func() initer { return new(ESR_MON_B8) },
	"ESR_MON_R4":  func() initer { return new(ESR_MON_R4) },
	"ESR_MON_X8":  func() initer { return new(ESR_MON_X8) },
}

// TestINIT checks that INIT resets every block: after its exported fields
// are changed, INIT brings it back to exactly the state of a fresh INIT.
func TestINIT(t *testing.T) {
	for name, newBlock := range initBlocks {
		want := newBlock()
		want.INIT()
		got := newBlock()
		got.INIT()
		if n := scramble(reflect.ValueOf(got).Elem()); n == 0 {
			t.Errorf("%s: no exported value to change", name)
			continue
		}
		if reflect.DeepEqual(got, want) {
			t.Errorf("%s: changing its fields left it unchanged", name)
			continue
		}
		got.INIT()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: INIT did not restore the initial values:\n got %+v\nwant %+v", name, got, want)
		}
	}
}

// TestINITDefaults checks INIT against the OSCAT initial values written in
// the source as "// default X" after a field: TestINIT alone cannot see a
// wrong default, since it compares INIT with itself.
func TestINITDefaults(t *testing.T) {
	defaults := sourceDefaults(t)
	checked := 0
	for name, newBlock := range initBlocks {
		b := newBlock()
		b.INIT()
		v := reflect.ValueOf(b).Elem()
		for _, d := range blockDefaults(v.Type(), defaults) {
			f := v.FieldByName(d.field)
			if !f.IsValid() {
				t.Errorf("%s.%s: no such field", name, d.field)
				continue
			}
			if err := matches(f, d.value); err != nil {
				t.Errorf("%s.%s after INIT: %v (OSCAT default %s)", name, d.field, err, d.value)
			}
			checked++
		}
	}
	switch {
	case len(defaults) == 0:
		t.Skip("no \"// default\" annotation in this package")
	case checked == 0:
		t.Fatal("annotated defaults, but none checked")
	}
	t.Logf("%d defaults checked", checked)
}

// fieldDefault is a field's documented initial value.
type fieldDefault struct{ field, value string }

// sourceDefaults reads the "// default" comments of the package's struct
// fields, by type name.
func sourceDefaults(t *testing.T) map[string][]fieldDefault {
	t.Helper()
	fset := token.NewFileSet()
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	out := map[string][]fieldDefault{}
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}
			for _, f := range st.Fields.List {
				if f.Comment == nil || len(f.Names) == 0 {
					continue
				}
				rest, ok := strings.CutPrefix(strings.TrimSpace(f.Comment.Text()), "default ")
				if !ok {
					continue
				}
				values := strings.Split(rest, ",")
				for i, name := range f.Names {
					v := values[0]
					if len(values) == len(f.Names) {
						v = values[i]
					} else if len(values) > 1 {
						continue // a list that does not match the names
					}
					out[ts.Name.Name] = append(out[ts.Name.Name], fieldDefault{name.Name, strings.TrimSpace(v)})
				}
			}
			return true
		})
	}
	return out
}

// blockDefaults returns the defaults of a block's own fields and of the
// structs it embeds.
func blockDefaults(t reflect.Type, defaults map[string][]fieldDefault) []fieldDefault {
	out := append([]fieldDefault(nil), defaults[t.Name()]...)
	for i := range t.NumField() {
		if sf := t.Field(i); sf.Anonymous && sf.Type.Kind() == reflect.Struct {
			out = append(out, blockDefaults(sf.Type, defaults)...)
		}
	}
	return out
}

// matches reports whether f holds the IEC literal lit.
func matches(f reflect.Value, lit string) error {
	switch f.Kind() {
	case reflect.Bool:
		want := strings.EqualFold(lit, "TRUE")
		if !want && !strings.EqualFold(lit, "FALSE") {
			return fmt.Errorf("cannot read %q as BOOL", lit)
		}
		if f.Bool() != want {
			return fmt.Errorf("got %v", f.Bool())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if strings.HasPrefix(lit, "T#") {
			want, err := iecTime(lit)
			if err != nil {
				return err
			}
			if time.Duration(f.Int()) != want {
				return fmt.Errorf("got %v", time.Duration(f.Int()))
			}
			return nil
		}
		want, err := iecInt(lit)
		if err != nil {
			return err
		}
		if f.Int() != want {
			return fmt.Errorf("got %d", f.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		want, err := iecInt(lit)
		if err != nil {
			return err
		}
		if f.Uint() != uint64(want) {
			return fmt.Errorf("got %d", f.Uint())
		}
	case reflect.Float32, reflect.Float64:
		want, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			return fmt.Errorf("cannot read %q as REAL", lit)
		}
		if got := f.Float(); gomath.Abs(got-want) > 1e-6*gomath.Max(1, gomath.Abs(want)) {
			return fmt.Errorf("got %v", got)
		}
	default:
		return fmt.Errorf("cannot compare a %s", f.Kind())
	}
	return nil
}

// iecInt reads an IEC integer literal: 42, -1, 16#FF, 2#1010_0000.
func iecInt(lit string) (int64, error) {
	base, digits := 10, strings.ReplaceAll(lit, "_", "")
	if b, d, ok := strings.Cut(digits, "#"); ok {
		n, err := strconv.Atoi(b)
		if err != nil {
			return 0, fmt.Errorf("cannot read %q as an integer", lit)
		}
		base, digits = n, d
	}
	u, err := strconv.ParseUint(strings.TrimPrefix(digits, "-"), base, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot read %q as an integer", lit)
	}
	if strings.HasPrefix(digits, "-") {
		return -int64(u), nil
	}
	return int64(u), nil
}

// iecTime reads an IEC TIME literal: T#500ms, T#1.2s, T#10d.
func iecTime(lit string) (time.Duration, error) {
	units := map[string]time.Duration{"d": 24 * time.Hour, "h": time.Hour, "m": time.Minute,
		"s": time.Second, "ms": time.Millisecond, "us": time.Microsecond}
	rest := strings.ToLower(strings.TrimPrefix(lit, "T#"))
	var total time.Duration
	for rest != "" {
		i := strings.IndexFunc(rest, func(r rune) bool { return (r < '0' || r > '9') && r != '.' })
		if i <= 0 {
			return 0, fmt.Errorf("cannot read %q as TIME", lit)
		}
		n, err := strconv.ParseFloat(rest[:i], 64)
		if err != nil {
			return 0, fmt.Errorf("cannot read %q as TIME", lit)
		}
		j := i
		for j < len(rest) && rest[j] >= 'a' && rest[j] <= 'z' {
			j++
		}
		u, ok := units[rest[i:j]]
		if !ok {
			return 0, fmt.Errorf("cannot read %q as TIME", lit)
		}
		total += time.Duration(n * float64(u))
		rest = rest[j:]
	}
	return total, nil
}

// scramble changes every exported field of v, in nested and embedded
// structs too, and returns how many it changed. Pointers are left alone: a
// block may keep one across INIT (INTEGRATE's output, for instance).
func scramble(v reflect.Value) int {
	n := 0
	for i := range v.NumField() {
		f := v.Field(i)
		if sf := v.Type().Field(i); sf.Anonymous && f.Kind() == reflect.Struct && !f.CanSet() {
			// An embedded unexported struct holds exported inputs and
			// outputs (FIFO_16 embeds fifo): reach them through its address.
			n += scramble(reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem())
			continue
		}
		if !f.CanSet() {
			continue
		}
		switch f.Kind() {
		case reflect.Bool:
			f.SetBool(!f.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.SetInt(f.Int() + 7)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			f.SetUint(f.Uint() + 7)
		case reflect.Float32, reflect.Float64:
			f.SetFloat(f.Float() + 7.5)
		case reflect.String:
			f.SetString(f.String() + "x")
		case reflect.Struct:
			n += scramble(f)
			continue
		case reflect.Array:
			for j := range f.Len() {
				if e := f.Index(j); e.Kind() == reflect.Struct {
					n += scramble(e)
				} else if e.Kind() == reflect.Bool {
					e.SetBool(!e.Bool())
					n++
				}
			}
			continue
		default:
			continue
		}
		n++
	}
	return n
}
