package builtins

import (
	parsers "github.com/snyk/snyk-iac-parsers/pkg"
	"io/ioutil"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
)

func RegisterHCLBuiltin() {
	rego.RegisterBuiltin1(
		&rego.Function{
			Name:    "hcl2.unmarshal_file",
			Decl:    types.NewFunction(types.Args(types.S), types.A),
			Memoize: true,
		},
		func(bctx rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {

			var filePath string

			if err := ast.As(a.Value, &filePath); err != nil {
				return nil, err
			}

			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}

			var parsedInput interface{}
			if err := parsers.ParseHCL2(content, &parsedInput); err != nil {
				return nil, err
			}
			v, err := ast.InterfaceToValue(parsedInput)
			if err != nil {
				return nil, err
			}

			return ast.NewTerm(v), nil
		},
	)
}
