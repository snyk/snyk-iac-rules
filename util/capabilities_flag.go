package util

import (
	"os"

	"github.com/open-policy-agent/opa/ast"
)

// CapabilitiesFlag was taken from https://github.com/open-policy-agent/opa/blob/v0.31.0/cmd/flags.go#L160

type CapabilitiesFlag struct {
	C    *ast.Capabilities
	path string
}

func NewCapabilitiesFlag() CapabilitiesFlag {
	return CapabilitiesFlag{
		// cannot call ast.CapabilitiesForThisVersion here because
		// custom builtins cannot be registered by this point in execution
		C: nil,
	}
}

func (f *CapabilitiesFlag) Type() string {
	return "string"
}

func (f *CapabilitiesFlag) String() string {
	return f.path
}

func (f *CapabilitiesFlag) Set(s string) error {
	f.path = s
	fd, err := os.Open(s)
	if err != nil {
		return err
	}
	defer fd.Close()
	f.C, err = ast.LoadCapabilitiesJSON(fd)
	return err
}
