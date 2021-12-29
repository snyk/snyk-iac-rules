package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
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
	err := assertValidRuleIds(args)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)

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
		WithAsBundle(false).
		WithOptimizationLevel(0).
		WithOutput(buf).
		WithEntrypoints(params.Entrypoint.Strings()...).
		WithPaths(args...).
		WithFilter(buildCommandLoaderFilter(false, params.Ignore))

	err = compiler.Build(context.Background())
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

func assertUniqueRuleIds(rules []util.Rule) error {
	visitedRulePaths := make(map[string]string)

	for _, rule := range rules {
		if _, ok := visitedRulePaths[rule.PublicId]; ok {
			return fmt.Errorf(
				"We cannot create a bundle for your custom rules."+
					"\nThe bundle contains duplicate rules."+
					"\nPlease ensure all rules have unique public IDs."+
					"\n\nDuplicate rules are:"+
					"\n- %s"+
					"\n- %s",
				rule.Path,
				visitedRulePaths[rule.PublicId],
			)
		}

		visitedRulePaths[rule.PublicId] = rule.Path
	}

	return nil
}

func assertRuleIdsWithoutSnykPrefix(rules []util.Rule) error {
	invalidRulesPaths := []string{}

	for _, rule := range rules {
		if strings.HasPrefix(rule.PublicId, "SNYK-") {
			invalidRulesPaths = append(invalidRulesPaths, rule.Path)
		}
	}

	if len(invalidRulesPaths) > 0 {
		errMessage := "We cannot create a bundle for your custom rules." +
			"\nCustom rules cannot have a name that starts with \"SNYK-\"." +
			"\nPlease ensure your public ID does not start with \"SNYK-\"." +
			"\n\nRules that start with \"SNYK-\" are:"
		for _, invalidRulePath := range invalidRulesPaths {
			errMessage += fmt.Sprintf("\n- %s", invalidRulePath)
		}
		return errors.New(errMessage)
	}

	return nil
}

func assertValidRuleIds(paths []string) error {
	rules, err := util.RetrieveRules(paths)
	if err != nil {
		return err
	}

	err = assertUniqueRuleIds(rules)
	if err != nil {
		return err
	}

	err = assertRuleIdsWithoutSnykPrefix(rules)
	if err != nil {
		return err
	}

	return nil
}
