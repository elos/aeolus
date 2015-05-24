package ego

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/elos/aeolus"
	"github.com/elos/ehttp/templates"
)

func Generate(hostFile, root string) error {
	h := aeolus.ParseHostFile(hostFile)

	var buf bytes.Buffer

	if err := engine.Execute(&buf, Routes, h); err != nil {
		return fmt.Errorf("error executing routes template: %s", err)
	}

	coreFile := filepath.Join(root, strings.ToLower(h.Name)+".go")
	routerFile := filepath.Join(root, "router.go")
	routesFile := filepath.Join(root, "routes/routes.go")
	routesContextFile := filepath.Join(root, "views/routes_context.go")

	if err := templates.ExecuteAndWriteGoFile(engine, Core, coreFile, h); err != nil {
		return err
	}

	if err := templates.ExecuteAndWriteGoFile(engine, Router, routerFile, h); err != nil {
		return err
	}

	if err := templates.ExecuteAndWriteGoFile(engine, Routes, routesFile, h); err != nil {
		return err
	}

	if err := templates.ExecuteAndWriteGoFile(engine, RoutesContext, routesContextFile, h); err != nil {
		return err
	}

	return nil
}
