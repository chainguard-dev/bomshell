package shell

import (
	"fmt"

	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
)

type Runner struct {
	Environment *cel.Env
	options     Options
	impl        RunnerImplementation
}

func NewRunner() (*Runner, error) {
	return NewRunnerWithOptions(&defaultOptions)
}

func NewRunnerWithOptions(opts *Options) (*Runner, error) {
	env, err := createEnvironment(&defaultOptions)
	if err != nil {
		return nil, err
	}
	runner := Runner{
		Environment: env,
		impl:        &defaultRunnerImplementation{},
	}

	return &runner, nil
}

// Compile reads CEL code from string, compiles it and
// returns the Abstract Syntax Tree (AST). The AST can then be evaluated
// in the environment. As compilation of the AST is expensive, it can
// be cached for better performance.
func (r *Runner) Compile(code string) (*cel.Ast, error) {
	// Run the compilation step
	ast, err := r.impl.Compile(r.Environment, code)
	if err != nil {
		return nil, fmt.Errorf("compiling program: %w", err)
	}
	return ast, nil
}

// EvaluateAST evaluates a CEL syntax tree on an SBOM. Returns the program
// evaluation result or an error.
func (r *Runner) EvaluateAST(ast *cel.Ast, doc *sbom.Document) (ref.Val, error) {
	program, err := r.Environment.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("generating program from AST: %w", err)
	}

	val := map[string]interface{}{
		"sbom": doc,
	}

	result, _, err := program.Eval(val)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %w", err)
	}

	return result, nil
}
