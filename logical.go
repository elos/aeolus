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
	EndpointDef struct {
		Name, Path, Auth  string
		Actions, Requires []string
		Endpoints         []*EndpointDef
	}

	HostDef struct {
		Name, Host string
		Endpoints  []*EndpointDef
	}

	Endpoint struct {
		Name, Path, Auth string
		Actions          []Action
		Requires         map[string]bool
	}

	Host struct {
		Name, Host string
		Endpoints  map[string]*Endpoint
	}
)

func (ed *EndpointDef) Valid() error {
	if ed.Name == "" {
		return errors.New("Endpoint must have a name")
	}

	if ed.Path == "" {
		return fmt.Errorf("Endpoint %s must have a path", ed.Name)
	}

	if string(ed.Path[0]) != "/" {
		return fmt.Errorf("Endpoint %s must have a path that starts with a /", ed.Name)
	}

	for _, e := range ed.Endpoints {
		if err := e.Valid(); err != nil {
			return err
		}
	}

	return nil
}

func (ed *EndpointDef) Process(namespace string, path string, auth string) []*Endpoint {
	endpoints := make([]*Endpoint, 1)

	actions := make([]Action, len(ed.Actions))
	for i, a := range ed.Actions {
		actions[i] = actionLiterals[a]
	}

	path = filepath.Join(path, ed.Path)
	var name string
	if namespace != "" {
		name = namespace + "_" + ed.Name
	} else {
		name = ed.Name
	}

	rs := make(map[string]bool)
	for _, a := range ed.Requires {
		rs[a] = true
	}

	if ed.Auth == "" {
		if auth == "" {
			auth = "none"
		}
	} else {
		auth = ed.Auth
	}

	endpoints[0] = &Endpoint{
		Name: name, Path: path, Auth: auth,
		Actions:  actions,
		Requires: rs,
	}

	for _, e := range ed.Endpoints {
		subpoints := e.Process(name, path, auth)
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

	for _, e := range hd.Endpoints {
		if err := e.Valid(); err != nil {
			return err
		}
	}

	return nil
}

func (hd *HostDef) Process() *Host {
	endpoints := make(map[string]*Endpoint, 0)

	for _, e := range hd.Endpoints {
		subpoints := e.Process("", "", "")
		for _, s := range subpoints {
			endpoints[s.Name] = s
		}
	}

	return &Host{
		Name:      hd.Name,
		Host:      hd.Host,
		Endpoints: endpoints,
	}
}
