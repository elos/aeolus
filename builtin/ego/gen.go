package ego

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
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

	if err := ExecuteAndWrite(Core, coreFile, h); err != nil {
		return err
	}

	if err := ExecuteAndWrite(Router, routerFile, h); err != nil {
		return err
	}

	if err := ExecuteAndWrite(Routes, routesFile, h); err != nil {
		return err
	}

	if err := ExecuteAndWrite(RoutesContext, routesContextFile, h); err != nil {
		return err
	}

	return nil
}

func ExecuteAndWrite(n templates.Name, file string, h *aeolus.Host) error {
	var buf bytes.Buffer

	if err := engine.Execute(&buf, n, h); err != nil {
		return fmt.Errorf("error executing template: %s", err)
	}

	if err := ioutil.WriteFile(file, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing file: %s", err)
	}

	if err := exec.Command("gofmt", "-w=true", file); err != nil {
		fmt.Errorf("error running gofmt  on file: %s", err)
	}

	if err := exec.Command("goimports", "-w=true", file); err != nil {
		fmt.Errorf("error running goimports on file: %s", err)
	}

	log.Printf("wrote %s", file)

	return nil
}
