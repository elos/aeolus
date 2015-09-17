package aeolus

import "testing"

const (
	resttodo       = "./tests/todo-rest.json"
	imperativetodo = "./tests/todo-imperative.json"
)

func TestValidRestTodoExample(t *testing.T) {
	hd, err := ParseHostFile(resttodo)

	if err != nil {
		t.Fatal(err)
	}

	if err := hd.Valid(); err != nil {
		t.Errorf("rest todo host example should be valid, but got: %s", err)
	}
}

func TestValidImperativeTodoExample(t *testing.T) {
	hd, err := ParseHostFile(imperativetodo)

	if err != nil {
		t.Fatal(err)
	}

	if err := hd.Valid(); err != nil {
		t.Errorf("imperative todo host example should be valid, but got: %s", err)
	}
}
