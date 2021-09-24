package internal

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/open-policy-agent/opa/loader"
	"github.com/snyk/snyk-iac-rules/util"

	"github.com/open-policy-agent/opa/util/test"

	"github.com/stretchr/testify/assert"
)

// Most of the test logic was taken from https://github.com/open-policy-agent/opa/blob/v0.31.0/cmd/build_test.go

func mockBuildParams() *BuildCommandParams {
	return &BuildCommandParams{
		Entrypoint: util.NewRepeatedStringFlag("test"),
		OutputFile: "",
		Ignore:     []string{},
		Target:     util.NewEnumFlag(TargetWasm, []string{TargetRego, TargetWasm}),
	}
}

func TestBuildProducesRegoBundle(t *testing.T) {
	files := map[string]string{
		"test.rego": `
			package test
			p = 1
		`,
	}

	test.WithTempFS(files, func(root string) {
		buildParams := mockBuildParams()
		buildParams.OutputFile = path.Join(root, "bundle.tar.gz")
		err := buildParams.Target.Set(TargetRego)
		assert.Nil(t, err)

		err = RunBuild([]string{root}, buildParams)
		assert.Nil(t, err)

		_, err = loader.NewFileLoader().AsBundle(buildParams.OutputFile)
		assert.Nil(t, err)

		// Check that manifest is not written given no input manifest and no other flags
		f, err := os.Open(buildParams.OutputFile)
		assert.Nil(t, err)
		defer f.Close()

		gr, err := gzip.NewReader(f)
		assert.Nil(t, err)

		tr := tar.NewReader(gr)

		for {
			f, err := tr.Next()
			if err == io.EOF {
				break
			}
			assert.Nil(t, err)

			if f.Name == "/data.json" || strings.HasSuffix(f.Name, "/test.rego") {
				continue
			}
			t.Fatal("unexpected file:", f.Name)
		}
	})
}

func TestBuildProducesWasmBundle(t *testing.T) {
	files := map[string]string{
		"test.rego": `
			package test
			p = 1
		`,
	}

	test.WithTempFS(files, func(root string) {
		buildParams := mockBuildParams()
		buildParams.OutputFile = path.Join(root, "bundle.tar.gz")

		err := RunBuild([]string{root}, buildParams)
		assert.Nil(t, err)

		_, err = loader.NewFileLoader().AsBundle(buildParams.OutputFile)
		assert.Nil(t, err)

		// Check that manifest is not written given no input manifest and no other flags
		f, err := os.Open(buildParams.OutputFile)
		assert.Nil(t, err)
		defer f.Close()

		gr, err := gzip.NewReader(f)
		assert.Nil(t, err)

		tr := tar.NewReader(gr)

		for {
			f, err := tr.Next()
			if err == io.EOF {
				break
			}
			assert.Nil(t, err)

			if f.Name == "/data.json" || strings.HasSuffix(f.Name, "/test.rego") || strings.HasSuffix(f.Name, "/policy.wasm") || strings.HasSuffix(f.Name, ".manifest") {
				continue
			}
			t.Fatal("unexpected file:", f.Name)
		}
	})
}

func TestBuildFilesystemModeIgnoresTarGz(t *testing.T) {
	files := map[string]string{
		"test.rego": `
			package test
			p = 1
		`,
	}

	test.WithTempFS(files, func(root string) {
		buildParams := mockBuildParams()
		buildParams.OutputFile = path.Join(root, "bundle.tar.gz")

		err := RunBuild([]string{root}, buildParams)
		assert.Nil(t, err)

		_, err = loader.NewFileLoader().AsBundle(buildParams.OutputFile)
		assert.Nil(t, err)

		// Just run the build again to simulate the user doing back-to-back builds.
		err = RunBuild([]string{root}, buildParams)
		assert.Nil(t, err)
	})
}

func TestBuildErrorDoesNotWriteFile(t *testing.T) {
	files := map[string]string{
		"test.rego": `
			package test
			p { p }
		`,
	}

	test.WithTempFS(files, func(root string) {
		buildParams := mockBuildParams()
		buildParams.OutputFile = path.Join(root, "bundle.tar.gz")

		err := RunBuild([]string{root}, buildParams)
		if err != nil {
			assert.Contains(t, err.Error(), "rule p is recursive")
		}

		_, err = os.Stat(buildParams.OutputFile)
		assert.NotNil(t, err)
	})
}
