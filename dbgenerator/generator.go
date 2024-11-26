package dbgenerator

import "text/template"

type Generator interface {
	GetTemplates() (*template.Template, error)
	GetMapping() Mapping
}

type GoGenerator struct {
	path string
}

func NewGoGenerator(path string) *GoGenerator {
	return &GoGenerator{path: path}
}

func (g GoGenerator) GetTemplates() (*template.Template, error) {
	var err error
	templates := template.New("root")
	templates, err = templates.ParseGlob("templates" + "/*")
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (g GoGenerator) GetMapping() Mapping {
	return GetGoTypeMap()
}
