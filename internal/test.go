package internal

import (
	"context"
	"os"
	"time"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/compile"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/tester"
	"github.com/open-policy-agent/opa/topdown"
	"github.com/open-policy-agent/opa/topdown/lineage"
	"github.com/snyk/snyk-iac-custom-rules/util"
)

// Most of the logic was taken from https://github.com/open-policy-agent/opa/blob/v0.31.0/cmd/test.go

const (
	ExplainModeFull  = "full"
	ExplainModeNotes = "notes"
	ExplainModeFails = "fails"
)

type TestCommandParams struct {
	Verbose  bool
	Explain  util.EnumFlag
	Timeout  time.Duration
	Ignore   []string
	RunRegex string
}

func RunTest(args []string, params *TestCommandParams) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := util.LoaderFilter{
		Ignore: params.Ignore,
	}

	modules, store, err := tester.Load(args, filter.Apply)
	if err != nil {
		return err
	}

	txn, err := store.NewTransaction(ctx, storage.WriteParams)
	if err != nil {
		return err
	}
	defer store.Abort(ctx, txn)

	compiler := ast.NewCompiler().
		WithPathConflictsCheck(storage.NonEmpty(ctx, store, txn))

	info, err := util.Term()
	if err != nil {
		return err
	}

	runner := tester.NewRunner().
		SetCompiler(compiler).
		SetStore(store).
		EnableTracing(params.Verbose).
		EnableFailureLine(false).
		SetRuntime(info).
		SetModules(modules).
		SetTimeout(params.Timeout).
		Filter(params.RunRegex).
		Target(compile.TargetRego)

	reporter := tester.PrettyReporter{
		Verbose:     params.Verbose,
		FailureLine: false,
		Output:      os.Stdout,
	}

	return runTests(ctx, txn, runner, reporter, params)
}

func runTests(ctx context.Context, txn storage.Transaction, runner *tester.Runner, reporter tester.Reporter, params *TestCommandParams) error {
	ch, err := runner.RunTests(ctx, txn)
	if err != nil {
		return err
	}

	dup := make(chan *tester.Result)

	go func() {
		defer close(dup)
		for tr := range ch {
			tr.Trace = filterTrace(tr.Trace, params)
			dup <- tr
		}
	}()

	if err := reporter.Report(dup); err != nil {
		return err
	}

	return nil
}

func filterTrace(trace []*topdown.Event, params *TestCommandParams) []*topdown.Event {
	ops := map[topdown.Op]struct{}{}
	mode := params.Explain.String()

	if mode == ExplainModeFull {
		// Don't bother filtering anything
		return trace
	}

	// If an explain mode was specified, filter based
	// on the mode. If no explain mode was specified,
	// default to show both notes and fail events
	showDefault := !params.Explain.IsSet() && params.Verbose

	if mode == ExplainModeNotes || showDefault {
		ops[topdown.NoteOp] = struct{}{}
	}
	if mode == ExplainModeFails || showDefault {
		ops[topdown.FailOp] = struct{}{}
	}

	return lineage.Filter(trace, func(event *topdown.Event) bool {
		_, relevant := ops[event.Op]
		return relevant
	})
}
