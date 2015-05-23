package aeolus

import (
	"errors"
	"fmt"
	"path/filepath"
)

type Action int

const (
	POST Action = iota
	GET
	DELETE
	OPTIONS
)

var actionLiterals = map[string]Action{
	"GET":     GET,
	"POST":    POST,
	"DELETE":  DELETE,
	"OPTIONS": OPTIONS,
}

type (
	// An EndpointDef is a go struct form of the JSON encoded Endpoint definition
	// An endpoint has a name, a path (which may be relative if the endpoint is a
	// subpoint), middleware and services defined for each HTTP Actions, and
	// subpoints (which are endpoint definitions with paths rooted with this
	// endpoint's path
	// i.e.,
	//		{
	//			"name": "groups",
	//			"path": "/groups",
	//			"actions": ["GET"],
	//			"middleware": {
	//				"GET": ["log", "user-auth"],
	//			},
	//			"services": {
	//				"GET": ["db", "response"],
	//			}
	//		}
	EndpointDef struct {
		Name, Path           string
		Actions              []string
		Middleware, Services map[string][]string
		Subpoints            []*EndpointDef
	}

	// A HostDef is the rooted JSON encoded definition of an aeolus.Host.
	// It has a name, a host (for listening on), and a port. It also must
	// specify officially recognized middleware and services. It has a list
	// of endpoints which may be defined heirarchically using subpoints.
	// i.e.,
	//		{
	//			"name": "api",
	//			"host": "localhost",
	//			"port": 8000,
	//			"middleware": ["log", "user-auth"],
	//			"services": ["db", "response"],
	//			"endpoints": [
	//				{ ... },
	//				{ ... }
	//			]
	//		}
	HostDef struct {
		Name, Host           string
		Port                 int
		Middleware, Services []string
		Endpoints            []*EndpointDef
		Static               map[string]string
	}

	// An Endpoint is the logical construct representing a API route handler.
	// An enpoint must have a name and path defined. The Actions are which
	// actions it accepts (which may be none) and Middleware and Services
	// are the corresponding requirements for each of those actions.
	Endpoint struct {
		Name, Path           string
		Action               Action
		Middleware, Services []string
	}

	// A Host is the logical construct represent and collection of endpoints,
	// often thought of a single server, host, api, or application.
	// A host must have a name, a serving host, a port and optionally
	// explicitly recognied middleware and services. Any middleware/services
	// not explicity recognized will be considered invalid if referred to in
	// any endpoint of the host.
	Host struct {
		Name, Host           string
		Port                 uint
		Middleware, Services map[string]bool
		Routes               map[string]string
		Endpoints            map[string]*Endpoint
		Static               map[string]string
	}
)

func (ed *EndpointDef) Valid(h *HostDef) error {
	if ed.Name == "" {
		return errors.New("Endpoint must have a name")
	}

	if ed.Path == "" {
		return fmt.Errorf("Endpoint %s must have a path", ed.Name)
	}

	if string(ed.Path[0]) != "/" {
		return fmt.Errorf("Endpoint %s must have a path that starts with a /", ed.Name)
	}

	for _, a := range ed.Actions {
		if _, ok := actionLiterals[a]; !ok {
			return fmt.Errorf("Endpoint %s has invalid HTTP Action: %s", ed.Name, a)
		}
	}

	for a, values := range ed.Middleware {
		if _, ok := actionLiterals[a]; !ok {
			return fmt.Errorf("Endpoint %s has invalid HTTP Action: %s, defined in middleware", ed.Name, a)
		}

		for _, v := range values {
			if !includes(h.Middleware, v) {
				return fmt.Errorf("Endpoint %s has invalid middleware %s", ed.Name, v)
			}
		}
	}

	for a, values := range ed.Services {
		if _, ok := actionLiterals[a]; !ok {
			return fmt.Errorf("Endpoint %s has invalid HTTP Action: %s, defined in services", ed.Name, a)
		}

		for _, v := range values {
			if !includes(h.Services, v) {
				return fmt.Errorf("Endpoint %s has invalid service %s", ed.Name, v)
			}
		}
	}

	for _, e := range ed.Subpoints {
		if err := e.Valid(h); err != nil {
			return err
		}
	}

	return nil
}

func includes(ss []string, s string) bool {
	for i := range ss {
		if s == ss[i] {
			return true
		}
	}

	return false
}

func (ed *EndpointDef) Process(namespace string, path string) []*Endpoint {
	endpoints := make([]*Endpoint, len(ed.Actions))

	var name string
	path = filepath.Join(path, ed.Path)
	if namespace != "" {
		name = namespace + "_" + ed.Name
	} else {
		name = ed.Name
	}

	for i, aString := range ed.Actions {
		action := actionLiterals[aString]

		m := make([]string, len(ed.Middleware[aString]))
		s := make([]string, len(ed.Services[aString]))

		for i, middleware := range ed.Middleware[aString] {
			m[i] = ProcessName(middleware)
		}

		for i, service := range ed.Services[aString] {
			s[i] = ProcessName(service)
		}

		endpoints[i] = &Endpoint{
			Name: name, Path: path,
			Action:     action,
			Middleware: m, Services: s,
		}
	}

	for _, e := range ed.Subpoints {
		subpoints := e.Process(name, path)
		for _, s := range subpoints {
			endpoints = append(endpoints, s)
		}
	}

	return endpoints
}

func (hd *HostDef) Valid() error {
	if hd.Name == "" {
		return fmt.Errorf("Host must have a name")
	}

	if hd.Host == "" {
		return fmt.Errorf("Host %s must have a host", hd.Name)
	}

	if hd.Port == 0 {
		return fmt.Errorf("Host %s must have port (0 is invalid)", hd.Name)
	}

	if hd.Port < 0 {
		return fmt.Errorf("Host %s must have positive port", hd.Name)
	}

	if hd.Port > 65535 {
		return fmt.Errorf("Host %s must have port less that 65535", hd.Name)
	}

	for _, m := range hd.Middleware {
		if m == "" {
			return fmt.Errorf("Host %s has invalid middleware: \"\"", hd.Name)
		}
	}

	for _, s := range hd.Services {
		if s == "" {
			return fmt.Errorf("Host %s has invalid service: \"\"", hd.Name)
		}
	}

	for _, e := range hd.Endpoints {
		if err := e.Valid(hd); err != nil {
			return err
		}
	}

	return nil
}

func (hd *HostDef) Process() *Host {
	endpoints := make(map[string]*Endpoint, 0)
	routes := make(map[string]string)

	for _, e := range hd.Endpoints {
		subpoints := e.Process("", "")
		for _, s := range subpoints {
			endpoints[s.Name+string(s.Action)] = s
			routes[s.Name] = s.Path
		}
	}

	m := make(map[string]bool)
	for _, middleware := range hd.Middleware {
		m[ProcessName(middleware)] = true
	}

	s := make(map[string]bool)
	for _, service := range hd.Services {
		s[ProcessName(service)] = true
	}

	return &Host{
		Name:       hd.Name,
		Host:       hd.Host,
		Port:       uint(hd.Port),
		Middleware: m,
		Services:   s,
		Routes:     routes,
		Endpoints:  endpoints,
		Static:     hd.Static,
	}
}
