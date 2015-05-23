package ego

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/elos/aeolus"
	"github.com/elos/metis"
)

func ActionLiteral(a aeolus.Action) string {
	return actionsToLiterals[a]
}

func name(s string) string {
	if s == "db" {
		return "DB"
	}
	return strings.Title(metis.CamelCase(s))
}

func isDynamic(e *aeolus.Endpoint) bool {
	return strings.Contains(e.Path, ":")
}

func signatureFor(e *aeolus.Endpoint) string {
	var buf bytes.Buffer
	tokens := strings.Split(e.Path, "/")
	args := make([]string, 0)
	for i := range tokens {
		if strings.Contains(tokens[i], ":") {
			args = append(args, string(tokens[i][1:]))
		}
	}

	fmt.Fprint(&buf, "func (")
	for i, arg := range args {
		if i != 0 {
			fmt.Fprint(&buf, ",")
		}
		fmt.Fprintf(&buf, "%s", arg)
	}
	if len(args) > 0 {
		fmt.Fprint(&buf, " string) string")
	} else {
		fmt.Fprint(&buf, ") string")
	}
	return buf.String()
}

func interpolatorFor(n, path string) string {
	var buf bytes.Buffer
	tokens := strings.Split(path, "/")
	args := make([]string, 0)
	for i := range tokens {
		if strings.Contains(tokens[i], ":") {
			args = append(args, string(tokens[i][1:]))
		}
	}

	fmt.Fprintf(&buf, "func (r *RoutesContext) %s(", name(n))
	for i, arg := range args {
		if i != 0 {
			fmt.Fprint(&buf, " ,")
		}
		fmt.Fprintf(&buf, "%s", arg)
	}
	if len(args) > 0 {
		fmt.Fprint(&buf, " string) string {")
	} else {
		fmt.Fprint(&buf, ") string {")
	}
	for i, token := range tokens {
		if strings.Contains(token, ":") {
			tokens[i] = "%s"
		}
	}
	fmt.Fprintf(&buf, "return fmt.Sprintf(\"%s\"", strings.Join(tokens, "/"))
	for _, arg := range args {
		fmt.Fprint(&buf, ",")
		fmt.Fprintf(&buf, "%s", arg)
	}
	fmt.Fprint(&buf, ")\n}")

	return buf.String()
}

func argsFor(e *aeolus.Endpoint) string {
	a := "c"

	for _, service := range e.Services {
		a += ", s." + name(service)
	}

	return a
}

func userAuth(e *aeolus.Endpoint) bool {
	return includes(e.Services, "user-auth")
}

func staticPath(path string) string {
	return filepath.Join("/", path, "*filepath")
}

func downcase(s string) string {
	return strings.ToLower(s)
}

func includes(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}
