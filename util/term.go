package util

import (
	"os"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/version"
)

// Term was taken from https://github.com/open-policy-agent/opa/blob/v0.31.0/internal/runtime/runtime.go#L23

func Term() (*ast.Term, error) {
	obj := ast.NewObject()
	env := ast.NewObject()

	for _, s := range os.Environ() {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) == 1 {
			env.Insert(ast.StringTerm(parts[0]), ast.NullTerm())
		} else if len(parts) > 1 {
			env.Insert(ast.StringTerm(parts[0]), ast.StringTerm(parts[1]))
		}
	}

	obj.Insert(ast.StringTerm("env"), ast.NewTerm(env))
	obj.Insert(ast.StringTerm("version"), ast.StringTerm(version.Version))
	obj.Insert(ast.StringTerm("commit"), ast.StringTerm(version.Vcs))

	return ast.NewTerm(obj), nil
}
