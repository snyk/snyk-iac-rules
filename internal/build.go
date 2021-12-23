package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/compile"
	"github.com/snyk/snyk-iac-rules/util"
)

// Most of the logic was taken from https://github.com/open-policy-agent/opa/blob/v0.31.0/cmd/build.go

const (
	TargetRego = "rego"
	TargetWasm = "wasm"
)

type BuildCommandParams struct {
	Entrypoint   util.RepeatedStringFlag
	OutputFile   string
	Ignore       []string
	Target       util.EnumFlag
	Capabilities util.CapabilitiesFlag
}

func RunBuild(args []string, params *BuildCommandParams) error {
	buf := bytes.NewBuffer(nil)

	metadataFile, err := createManifest(args)
	if err == nil && metadataFile != "" {
		defer os.Remove(metadataFile)
	}

	var capabilities *ast.Capabilities
	// if capabilities are not provided as a cmd flag,
	// then ast.CapabilitiesForThisVersion must be called
	// within dobuild to ensure custom builtins are properly captured
	if params.Capabilities.C != nil {
		capabilities = params.Capabilities.C
	} else {
		capabilities = ast.CapabilitiesForThisVersion()
	}
	compiler := compile.New().
		WithCapabilities(capabilities).
		WithTarget(params.Target.String()).
		WithAsBundle(true).
		WithOptimizationLevel(0).
		WithOutput(buf).
		WithEntrypoints(params.Entrypoint.Strings()...).
		WithPaths(args...).
		WithFilter(buildCommandLoaderFilter(true, params.Ignore))

	err = compiler.Build(context.Background())
	if err != nil {
		return err
	}

	out, err := os.Create(params.OutputFile)
	defer out.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, buf)
	if err != nil {
		return err
	}

	return nil
}

func buildCommandLoaderFilter(bundleMode bool, ignore []string) func(string, os.FileInfo, int) bool {
	return func(abspath string, info os.FileInfo, depth int) bool {
		if !bundleMode {
			if !info.IsDir() && strings.HasSuffix(abspath, ".tar.gz") {
				return true
			}
		}
		return util.LoaderFilter{Ignore: ignore}.Apply(abspath, info, depth)
	}
}

func createManifest(inputPaths []string) (string, error) {
	rules, err := util.RetrieveRules(inputPaths)
	if err != nil {
		return "", err
	}

	// choose either one of the locations for the metadata.json
	// it will be included in the OPA generated data.json
	metadataFile := path.Join(inputPaths[0], ".manifest")
	err = os.WriteFile(metadataFile, []byte(fmt.Sprintf("{\"metadata\":{\"numberOfRules\":%d}}", len(rules))),
		0644)
	if err != nil {
		return "", err
	}
	return metadataFile, nil
}
