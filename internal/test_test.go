package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/topdown"
	"github.com/open-policy-agent/opa/util/test"

	"github.com/snyk/snyk-iac-rules/util"

	"github.com/stretchr/testify/assert"
)

// Most of the test logic was taken from https://github.com/open-policy-agent/opa/blob/v0.31.0/cmd/test_test.go

func mockTestParams() *TestCommandParams {
	return &TestCommandParams{
		Verbose:  false,
		Explain:  util.NewEnumFlag(util.ExplainModeFails, []string{util.ExplainModeFails, util.ExplainModeFull, util.ExplainModeNotes}),
		Timeout:  5 * time.Second,
		Ignore:   []string{},
		RunRegex: "",
	}
}

func TestFilterTraceDefault(t *testing.T) {
	testParams := mockTestParams()
	testParams.Verbose = false
	expected := `Enter data.testing.test_p = _
| Enter data.testing.test_p
| | Enter data.testing.p
| | | Enter data.testing.q
| | | | Enter data.testing.r
| | | | | Fail x = data.x
| | | | Fail data.testing.r[x]
| | | Fail data.testing.q.foo
| | Fail data.testing.p with data.x as "bar"
| Fail data.testing.test_p = _
`
	verifyFilteredTrace(t, testParams, expected)
}

func TestFilterTraceVerbose(t *testing.T) {
	testParams := mockTestParams()
	testParams.Verbose = true
	expected := `Enter data.testing.test_p = _
| Enter data.testing.test_p
| | Enter data.testing.p
| | | Note "test test"
| | | Enter data.testing.q
| | | | Note "got this far"
| | | | Enter data.testing.r
| | | | | Note "got this far2"
| | | | | Fail x = data.x
| | | | Fail data.testing.r[x]
| | | Fail data.testing.q.foo
| | Fail data.testing.p with data.x as "bar"
| Fail data.testing.test_p = _
`
	verifyFilteredTrace(t, testParams, expected)
}

func TestFilterTraceExplainFails(t *testing.T) {
	testParams := mockTestParams()
	err := testParams.Explain.Set(util.ExplainModeFails)
	assert.Nil(t, err)

	expected := `Enter data.testing.test_p = _
| Enter data.testing.test_p
| | Enter data.testing.p
| | | Enter data.testing.q
| | | | Enter data.testing.r
| | | | | Fail x = data.x
| | | | Fail data.testing.r[x]
| | | Fail data.testing.q.foo
| | Fail data.testing.p with data.x as "bar"
| Fail data.testing.test_p = _
`
	verifyFilteredTrace(t, testParams, expected)
}

func TestFilterTraceExplainNotes(t *testing.T) {
	testParams := mockTestParams()
	err := testParams.Explain.Set(util.ExplainModeNotes)
	assert.Nil(t, err)

	expected := `Enter data.testing.test_p = _
| Enter data.testing.test_p
| | Enter data.testing.p
| | | Note "test test"
| | | Enter data.testing.q
| | | | Note "got this far"
| | | | Enter data.testing.r
| | | | | Note "got this far2"
`
	verifyFilteredTrace(t, testParams, expected)
}

func TestFilterTraceExplainFull(t *testing.T) {
	testParams := mockTestParams()
	err := testParams.Explain.Set(util.ExplainModeFull)
	assert.Nil(t, err)

	expected := `Enter data.testing.test_p = _
| Eval data.testing.test_p = _
| Unify data.testing.test_p = _
| Index data.testing.test_p (matched 1 rule, early exit)
| Enter data.testing.test_p
| | Eval data.testing.p with data.x as "bar"
| | Unify data.testing.p = _
| | Index data.testing.p (matched 1 rule, early exit)
| | Enter data.testing.p
| | | Eval data.testing.x
| | | Unify data.testing.x = _
| | | Index data.testing.x (matched 1 rule, early exit)
| | | Enter data.testing.x
| | | | Eval data.testing.y
| | | | Unify data.testing.y = _
| | | | Index data.testing.y (matched 1 rule, early exit)
| | | | Enter data.testing.y
| | | | | Eval true
| | | | | Unify true = _
| | | | | Exit data.testing.y early
| | | | Unify true = _
| | | | Exit data.testing.x early
| | | Unify true = _
| | | Eval trace("test test")
| | | Note "test test"
| | | Eval data.testing.q.foo
| | | Unify data.testing.q.foo = _
| | | Index data.testing.q (matched 1 rule)
| | | Enter data.testing.q
| | | | Unify x = "foo"
| | | | Eval trace("got this far")
| | | | Note "got this far"
| | | | Eval data.testing.r[x]
| | | | Unify data.testing.r[x] = _
| | | | Index data.testing.r (matched 1 rule)
| | | | Enter data.testing.r
| | | | | Unify x = "foo"
| | | | | Eval trace("got this far2")
| | | | | Note "got this far2"
| | | | | Eval x = data.x
| | | | | Unify "foo" = data.x
| | | | | Unify "foo" = "bar"
| | | | | Fail x = data.x
| | | | | Redo trace("got this far2")
| | | | Fail data.testing.r[x]
| | | | Redo trace("got this far")
| | | Fail data.testing.q.foo
| | | Redo trace("test test")
| | | Redo data.testing.x
| | | Redo data.testing.x
| | | | Redo data.testing.y
| | | | | Redo true
| | Fail data.testing.p with data.x as "bar"
| Fail data.testing.test_p = _
`
	verifyFilteredTrace(t, testParams, expected)
}

func verifyFilteredTrace(t *testing.T, params *TestCommandParams, expected string) {
	filtered := filterTrace(failTrace(t), params)

	var buff bytes.Buffer
	topdown.PrettyTrace(&buff, filtered)
	actual := buff.String()
	fmt.Println(actual)

	assert.Equal(t, expected, actual)
}

func failTrace(t *testing.T) []*topdown.Event {
	t.Helper()
	mod := `
	package testing
	
	p {
		x  # Always true
		trace("test test")
		q["foo"]
	}
	
	x {
		y
	}
	
	y {
		true
	}
	
	q[x] {
		some x
		trace("got this far")
		r[x]
		trace("got this far1")
	}
	
	r[x] {
		trace("got this far2")
		x := data.x
	}
	
	test_p {
		p with data.x as "bar"
	}
	`

	tracer := topdown.NewBufferTracer()

	_, err := rego.New(
		rego.Module("test.rego", mod),
		rego.Trace(true),
		rego.QueryTracer(tracer),
		rego.Query("data.testing.test_p"),
	).Eval(context.Background())
	assert.Nil(t, err)

	return *tracer
}

func TestInvalidRego(t *testing.T) {
	baseTestFiles := map[string]string{
		"test1.rego": `
			package test
			*
		`,
	}

	test.WithTempFS(baseTestFiles, func(root string) {
		testParams := mockTestParams()
		testParams.Verbose = false

		rescueStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		err := RunTest([]string{root}, testParams)
		assert.NotNil(t, err)

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stderr = rescueStderr

		assert.Contains(t, string(out), "unexpected mul token")
	})
}

func TestRegoUnsafeVariable(t *testing.T) {
	baseTestFiles := map[string]string{
		"test1.rego": `package rules

deny[msg] {
	unsafe
}
`,
	}

	test.WithTempFS(baseTestFiles, func(root string) {
		testParams := mockTestParams()
		testParams.Verbose = false

		rescueStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		err := RunTest([]string{root}, testParams)
		assert.NotNil(t, err)

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stderr = rescueStderr

		assert.Contains(t, string(out), "rego_unsafe_var_error")
	})
}
