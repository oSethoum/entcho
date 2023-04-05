package entcho

import (
	"fmt"
	"path"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"entgo.io/ent/schema/field"
)

var (
	snake      = gen.Funcs["snake"].(func(string) string)
	plural     = gen.Funcs["plural"].(func(string) string)
	buggyCamel = gen.Funcs["camel"].(func(string) string)
	camel      = func(s string) string { return buggyCamel(snake(s)) }
)

func init() {
	gen.Funcs["tag"] = tag
	gen.Funcs["imports"] = imports
	gen.Funcs["null_field_create"] = null_field_create
	gen.Funcs["null_field_update"] = null_field_update
	gen.Funcs["extract_type"] = extract_type
	gen.Funcs["edge_field"] = edge_field
	gen.Funcs["is_comparable"] = is_comparable
	gen.Funcs["enum_or_edge_filed"] = enum_or_edge_filed
	gen.Funcs["get_name"] = get_name
	gen.Funcs["get_type"] = get_type
	gen.Funcs["is_slice"] = is_slice
	gen.Funcs["id_type"] = id_type
	gen.Funcs["go_ts"] = go_ts
	gen.Funcs["order_fields"] = order_fields
	gen.Funcs["select_fields"] = select_fields
	gen.Funcs["dir"] = path.Dir
}

func tag(f *load.Field) string {
	if f.Tag == "" {
		name := camel(f.Name)
		if strings.HasSuffix(name, "ID") {
			name = strings.TrimSuffix(name, "ID")
			name += "Id"
		}
		return fmt.Sprintf("json:\"%s,omitempty\"", name)
	}
	return f.Tag
}

func imports(g *gen.Graph, isInput ...bool) []string {
	imps := []string{}

	for _, s := range g.Schemas {
		for _, f := range s.Fields {
			if len(f.Enums) > 0 && len(isInput) > 0 && isInput[0] {
				imps = append(imps, path.Join(g.Package, strings.Split(f.Info.Ident, ".")[0]))
			}
			if f.Info != nil && len(f.Info.PkgPath) != 0 {
				if !in(f.Info.PkgPath, imps) {
					imps = append(imps, f.Info.PkgPath)
				}
			}
		}
	}
	return imps
}

func null_field_create(f *load.Field) bool {
	return f.Optional || f.Default
}

func null_field_update(field *load.Field) bool {
	return !strings.HasPrefix(extract_type(field), "[]")
}

func extract_type(field *load.Field) string {
	if field.Info.Ident != "" {
		return field.Info.Ident
	}
	return field.Info.Type.String()
}

func edge_field(e *load.Edge) bool {
	return e.Field != ""
}

func is_comparable(f *load.Field) bool {
	return has_prefixes(extract_type(f), []string{
		"string",
		"int",
		"uint",
		"float",
		"time.Time",
	})
}

func enum_or_edge_filed(s *load.Schema, f *load.Field) bool {
	for _, e := range s.Edges {
		if e.Field == f.Name {
			return extract_type(f) == "enum"
		}
	}
	return false
}

func get_name(f *load.Field) string {
	n := camel(f.Name)
	if strings.HasSuffix(n, "ID") {
		n = strings.TrimSuffix(n, "ID") + "Id"
	}
	return n
}

func get_type(t *field.TypeInfo) string {
	return go_ts(t.Type.String())
}

func go_ts(s string) string {
	slice := false
	if strings.HasPrefix(s, "[]") {
		slice = true
		s = strings.TrimPrefix(s, "[]")
	}
	for k, v := range gots {
		if strings.HasPrefix(s, k) {
			if slice {
				return v + "[]"
			}
			return v
		}
	}
	if slice {
		return s + "[]"
	}
	return s
}

func is_slice(f *load.Field) bool {
	return strings.HasPrefix(get_type(f.Info), "[]")
}

func id_type(s *load.Schema) string {
	for _, f := range s.Fields {
		if strings.ToLower(f.Name) == "id" {
			return get_type(f.Info)
		}
	}
	return "number"
}

func order_fields(s *load.Schema) string {
	fields := []string{}
	for _, f := range s.Fields {
		if orderable(f) {
			fields = append(fields, get_name(f))
		}
	}
	return "\"" + strings.Join(fields, "\" | \"") + "\""
}

func select_fields(s *load.Schema) string {
	fields := []string{}
	for _, f := range s.Fields {
		fields = append(fields, get_name(f))
	}
	return "\"" + strings.Join(fields, "\" | \"") + "\""
}

func orderable(f *load.Field) bool {
	return has_prefixes(extract_type(f), []string{
		"string",
		"int",
		"uint",
		"float",
		"time.Time",
		"bool",
	})
}
