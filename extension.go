package entcho

import (
	"embed"

	"entgo.io/ent/entc/gen"
)

//go:embed templates
var templates embed.FS

func (e *extension) Hooks() []gen.Hook {
	return e.hooks
}

func NewExtension(opts ...option) *extension {
	e := new(extension)
	for _, opt := range opts {
		opt(e)
	}
	e.hooks = append(e.hooks, e.generate)
	return e
}
