package shell

import (
	"fmt"
	"io"
	"os"

	"github.com/bom-squad/protobom/pkg/reader"
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
)

type BomShellImplementation interface {
	Compile(*Runner, string) (*cel.Ast, error)
	Evaluate(*Runner, *cel.Ast, map[string]interface{}) (ref.Val, error)
	LoadSBOM(io.ReadSeekCloser) (*sbom.Document, error)
	OpenFile(path string) (*os.File, error)
}

type DefaultBomShellImplementation struct{}

func (di *DefaultBomShellImplementation) Compile(runner *Runner, code string) (*cel.Ast, error) {
	return runner.Compile(code)
}

func (di *DefaultBomShellImplementation) Evaluate(runner *Runner, ast *cel.Ast, variables map[string]interface{}) (ref.Val, error) {
	return runner.EvaluateAST(ast, variables)
}

func (di *DefaultBomShellImplementation) LoadSBOM(stream io.ReadSeekCloser) (*sbom.Document, error) {
	r := reader.New()
	doc, err := r.ParseReader(stream)
	if err != nil {
		return nil, fmt.Errorf("parsing SBOM: %w", err)
	}

	return doc, nil
}

func (di *DefaultBomShellImplementation) OpenFile(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening SBOM file: %w", err)
	}
	return f, nil

}

///func HandleDocumentResult()
