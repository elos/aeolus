package aeolus

// Action represents an HTTP Action
type Action int

// Enumeration of the 4 supported action types
const (
	GET Action = iota
	POST
	DELETE
	OPTIONS
)

// The Aeolus literals which map to these types
var actionLiterals = map[string]Action{
	"GET":     GET,
	"POST":    POST,
	"DELETE":  DELETE,
	"OPTIONS": OPTIONS,
}

// Each of the logical components of defining an aeolus.Host
type (
	// An Endpoint is the logical construct representing a API route handler.
	// An enpoint must have a name and path defined. The Actions are which
	// actions it accepts (which may be none) and Middleware and Services
	// are the corresponding requirements for each of those actions.
	Endpoint struct {
		Name, Path           string
		Action               Action
		Middleware, Services []string
	}

	// A Host is the logical construct represent by a collection of endpoints,
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
		Endpoints            []*Endpoint
		Static               map[string]string
	}
)
