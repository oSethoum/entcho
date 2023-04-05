package main

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path"

	"entgo.io/ent/entc/gen"
)

//go:embed templates
var templates embed.FS

func main() {
	writeTemplae("entc", "entc.go")
	writeTemplae("generate", "generate.go")
}

func writeTemplae(name string, p string) {
	in, err := templates.ReadFile("templates/" + name + ".tmpl")
	catch(err)
	t, err := template.New(name).Funcs(gen.Funcs).Parse(string(in))
	catch(err)
	out := new(bytes.Buffer)
	err = t.Execute(out, nil)
	catch(err)
	p = path.Join("ent/generate/", p)
	err = os.MkdirAll(path.Dir(p), 0777)
	catch(err)
	err = os.WriteFile(p, out.Bytes(), fs.ModePerm)
	catch(err)
}

func catch(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
