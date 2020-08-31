package protocols

import (
	"fmt"
)

var serverReg = map[string]Factory{}

// Factory is used by output plugins to build an output instance
type Factory func() (Server, error)

// RegisterType registers a new output type.
func RegisterType(name string, f Factory) {
	if serverReg[name] != nil {
		panic(fmt.Errorf("server type  '%v' exists already", name))
	}
	serverReg[name] = f
}

// findFactory finds an output type its factory if available.
func findFactory(name string) Factory {
	return serverReg[name]
}

func Load(name string) (Server, error) {
	factory := findFactory(name)
	if factory == nil {
		return nil, fmt.Errorf("server type %v undefined", name)
	}

	return factory()
}
