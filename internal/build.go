package internal

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"

	"github.com/open-policy-agent/opa/compile"
	"github.com/snyk/snyk-iac-custom-rules/util"
)

// Most of the logic was taken from https://github.com/open-policy-agent/opa/blob/v0.31.0/cmd/build.go

const (
	TargetRego = "rego"
	TargetWasm = "wasm"
)

type BuildCommandParams struct {
	Entrypoint util.RepeatedStringFlag
	OutputFile string
	Ignore     []string
	Target     util.EnumFlag
}

func RunBuild(args []string, params *BuildCommandParams) error {
	buf := bytes.NewBuffer(nil)

	compiler := compile.New().
		WithTarget(params.Target.String()).
		WithAsBundle(false).
		WithOutput(buf).
		WithEntrypoints(params.Entrypoint.Strings()...).
		WithPaths(args...).
		WithFilter(buildCommandLoaderFilter(false, params.Ignore))

	err := compiler.Build(context.Background())
	if err != nil {
		return err
	}

	out, err := os.Create(params.OutputFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, buf)
	if err != nil {
		return err
	}

	return out.Close()
}

func buildCommandLoaderFilter(bundleMode bool, ignore []string) func(string, os.FileInfo, int) bool {
	return func(abspath string, info os.FileInfo, depth int) bool {
		if !info.IsDir() && strings.HasSuffix(abspath, ".tar.gz") {
			return true
		}
		return util.LoaderFilter{Ignore: ignore}.Apply(abspath, info, depth)
	}
}
