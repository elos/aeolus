package aeolus

import (
	"errors"
	"fmt"
	"path/filepath"
)

// The structures representing a parsed aeolus definition. For the actual
// logical structures, see logical.go.
type (
	// An EndpointDef is a go struct form of the JSON encoded Endpoint definition
	// An endpoint has a name, a path (which may be relative if the endpoint is a
	// subpoint), middleware and services defined for each HTTP Action. It may also
	// have "subpoints," which are endpoint definitions with paths rooted with this
	// endpoint's path.
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
	//			},
	//			subpoints: [
	//				{ ... }
	//			]
	//		}
	EndpointDef struct {
		Name, Path           string
		Actions              []string
		Middleware, Services map[string][]string
		Subpoints            []*EndpointDef
	}

	// A HostDef is the rooted JSON encoded definition of an aeolus.Host.
	// It has a name, a host (for listening on), and a port. It also must
	// specify officially recognized (valid for endpoints to declare)
	// middleware and services. It has a list of endpoints which may be
	// defined hierarchically using subpoints. A host is a tree of endpoints.
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
)

// EndpointDef: Valid(h *HostDef) error, Process --- {{{

// Valid returns an error if the endpoint definition is invalid, otherwise nil
// An endpoint can be invalid for 6 reasons:
//  1. It lacks a name
//  2. It lacks a path
//  3. It lacks a path which begins "/"
//  4. It declares an invalid HTTP action either as supported,
//     or in as partition in middleware or services
//  5. It declares unrecognized middleware or services
//  6. Any of its subpoints are invalid for the aforementioned 5 reasons
func (ed *EndpointDef) Valid(h *HostDef) error {
	// (1): no name
	if ed.Name == "" {
		return errors.New("endpoint definition must have a name")
	}

	// (2): no path
	if ed.Path == "" {
		return fmt.Errorf("endpoint definition '%s' must have a path", ed.Name)
	}

	// (3): path doesn't start with /
	if string(ed.Path[0]) != "/" {
		return fmt.Errorf("endpoint definition '%s' must have a path that begins with a '/'", ed.Name)
	}

	// (4): unrecognized action
	for _, a := range ed.Actions {
		if _, ok := actionLiterals[a]; !ok {
			return fmt.Errorf("endpoint definition '%s' has an invalid http action: '%s'", ed.Name, a)
		}
	}

	for action, wares := range ed.Middleware {
		// (4): unrecognized action
		if _, ok := actionLiterals[action]; !ok {
			return fmt.Errorf("endpoint definition '%s' has middleware definition with an invalid http action: '%s'", ed.Name, action)
		}

		// (5): unrecognized middleware
		for _, v := range wares {
			if !includes(h.Middleware, v) {
				return fmt.Errorf("endpoint definition '%s' declares unrecognized middleware: '%s'", ed.Name, v)
			}
		}
	}

	for action, services := range ed.Services {
		// (4): unrecognized action
		if _, ok := actionLiterals[action]; !ok {
			return fmt.Errorf("endpoint definition '%s' has services definition with an invalid http action: '%s'", ed.Name, action)
		}

		// (5): unrecognized service
		for _, v := range services {
			if !includes(h.Services, v) {
				return fmt.Errorf("endpoint definition '%s' declares uncregonized service: '%s'", ed.Name, v)
			}
		}
	}

	// (6): recursively check the validity of each subpoint
	for _, e := range ed.Subpoints {
		if err := e.Valid(h); err != nil {
			return err
		}
	}

	return nil
}

func (ed *EndpointDef) Process(namespace string, path string) []*Endpoint {
	// each endpoint has unique action
	endpoints := make([]*Endpoint, len(ed.Actions))

	path = filepath.Join(path, ed.Path)

	var name string
	if namespace != "" {
		name = namespace + "_" + ed.Name
	} else {
		name = ed.Name
	}

	for i, aString := range ed.Actions {
		action := actionLiterals[aString]

		m := make([]string, len(ed.Middleware[aString]))
		s := make([]string, len(ed.Services[aString]))

		for j, middleware := range ed.Middleware[aString] {
			m[j] = ProcessName(middleware)
		}

		for j, service := range ed.Services[aString] {
			s[j] = ProcessName(service)
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

// --- }}}

// HostDef Valid() error, Process --- {{{

// Valid returns an error if the host definition is invalid, otherwise nil
// A host can be invalid for n reasons:
//  1. It lacks a name
//  2. It lacks a host address
//  3. It lacks a positive port
//  4. It has empty middleware or service string
//  5. Any of it's endpoints are invalid
func (hd *HostDef) Valid() error {
	// (1) no name
	if hd.Name == "" {
		return fmt.Errorf("host definition lacks  a name")
	}

	// (2) no host address
	if hd.Host == "" {
		return fmt.Errorf("host definition '%s' lacks a host address", hd.Name)
	}

	// (3) bad port
	if hd.Port <= 0 || hd.Port > 65535 {
		return fmt.Errorf("host definition %s lacks valid port have port (0 < p < 65535)", hd.Name)
	}

	// (4) empty middlware name
	for _, m := range hd.Middleware {
		if m == "" {
			return fmt.Errorf("host defintion '%s' has invalid middleware: \"\"", hd.Name)
		}
	}

	// (4) empty service name
	for _, s := range hd.Services {
		if s == "" {
			return fmt.Errorf("host definition '%s' has invalid service: \"\"", hd.Name)
		}
	}

	// (5) validity of endpoints
	for _, e := range hd.Endpoints {
		if err := e.Valid(hd); err != nil {
			return err
		}
	}

	return nil
}

func (hd *HostDef) Process() *Host {
	endpoints := make([]*Endpoint, 0)
	routes := make(map[string]string)

	for _, e := range hd.Endpoints {
		subpoints := e.Process("", "")
		for _, s := range subpoints {
			endpoints = append(endpoints, s)
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

// --- }}}
