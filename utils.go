package entcho

import (
	"bytes"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"entgo.io/ent/entc/gen"
)

func parseTemplate(name string, data any) string {
	in, err := templates.ReadFile("templates/" + name + ".tmpl")
	if err != nil {
		log.Fatalln(err)
	}
	t, err := template.New(name).Funcs(gen.Funcs).Parse(string(in))
	if err != nil {
		log.Fatalln(err)
	}
	out := new(bytes.Buffer)
	err = t.Execute(out, data)
	if err != nil {
		log.Fatalln(err)
	}
	return out.String()
}

func in[T comparable](v T, vs []T) bool {
	for _, v2 := range vs {
		if v == v2 {
			return true
		}
	}
	return false
}

func has_prefixes(s string, px []string) bool {
	for _, p := range px {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func writeFile(destination string, content string) {
	destination = path.Join(get_gomod_dir(), destination)
	os.MkdirAll(path.Dir(destination), 0777)
	err := os.WriteFile(destination, []byte(content), 07777)
	if err != nil {
		log.Fatalln(err)
	}
}

func catch(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func get_gomod_dir() string {
	current, err := os.Getwd()
	catch(err)

again:
	_, err = os.Stat(path.Join(current, "go.mod"))
	if err != nil {
		current = path.Join(current, "../")
		if current == "/" {
			log.Fatalln("go.mod not found")
		}
		goto again
	}
	return current
}
