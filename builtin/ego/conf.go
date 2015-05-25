package ego

import (
	"log"
	"path/filepath"
	"text/template"

	"github.com/elos/aeolus"
	"github.com/elos/ehttp/templates"
)

var (
	importPath    = "github.com/elos/aeolus/builtin/ego"
	packagePath   = templates.PackagePath(importPath)
	templatesPath = filepath.Join(packagePath, "templates")
)

// Recognized templates
const (
	Core templates.Name = iota
	Router
	Routes
	RoutesContext
)

// aelous.Actions to equivalent go syntax strings
var actionsToLiterals = map[aeolus.Action]string{
	aeolus.POST:   "POST",
	aeolus.GET:    "GET",
	aeolus.DELETE: "DELETE",
}

var templateSet = &templates.TemplateSet{
	Core:          []string{"core.tmpl"},
	Router:        []string{"router.tmpl"},
	Routes:        []string{"routes.tmpl"},
	RoutesContext: []string{"routes_context.tmpl"},
}

var functionMap = template.FuncMap{
	"name":            name,
	"action":          ActionLiteral,
	"argsFor":         argsFor,
	"userAuth":        userAuth,
	"signatureFor":    signatureFor,
	"interpolatorFor": interpolatorFor,
	"staticPath":      staticPath,
	"downcase":        downcase,
	"reverse":         aeolus.Reverse,
}

var engine = templates.NewEngine(templatesPath, templateSet).WithFuncMap(functionMap)

func init() {
	if err := engine.Parse(); err != nil {
		log.Fatal(err)
	}
}
