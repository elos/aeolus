aeolus [![GoDoc](https://godoc.org/github.com/elos/aeolus?status.svg)](https://godoc.org/github.com/elos/aeolus)
------

Package aeolus provides logical data structures to model web services and their individual endpoints. It provides a configuration which specifies both external interface, which routes are available and to which HTTP actions they response, and internal requirements, which middleware is used and which services are required.

#### Architecture
The architecture of aeolus is quite simple. There is a `Host` definition and `Endpoint` definitions. A `Host` has a name, a host (to listen on), a port, and lists of recognized middleware and services. It also has a list of endpoints. `Enpoint`s may be recursively defined through `Subpoints` on the `Endpoint`, and have a name, a path, a list of HTTP Actions responded to and the services and middleware required for each.

##### Endpoint
##### Host
