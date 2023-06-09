package shell

import (
	"fmt"
	"io"
	"os"

	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/google/cel-go/common/types/ref"
)

const (
	protoDocumentType = "bomsquad.protobom.Document"
	DefaultFormat     = formats.SPDX23JSON
)

type Options struct {
	SBOM   string
	Format formats.Format
}

var defaultOptions = Options{
	Format: DefaultFormat,
}

type BomShell struct {
	Options Options
	runner  *Runner
	impl    BomShellImplementation
}

func New() (*BomShell, error) {
	return NewWithOptions(defaultOptions)
}

func NewWithOptions(opts Options) (*BomShell, error) {
	runner, err := NewRunnerWithOptions(&opts)
	if err != nil {
		return nil, fmt.Errorf("creating runner: %w", err)
	}
	return &BomShell{
		Options: opts,
		runner:  runner,
		impl:    &DefaultBomShellImplementation{},
	}, nil
}

func (bs *BomShell) RunFile(path string) (ref.Val, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading program data: %w", err)
	}
	return bs.Run(string(data))
}

func (bs *BomShell) Run(code string) (ref.Val, error) {
	// Variables that wil be made available in the CEL env
	vars := map[string]interface{}{}

	// Load an SBNOM if defined
	if bs.Options.SBOM != "" {
		f, err := bs.impl.OpenFile(bs.Options.SBOM)
		if err != nil {
			return nil, fmt.Errorf("opening SBOM file: %w", err)
		}

		doc, err := bs.impl.LoadSBOM(f)
		if err != nil {
			return nil, fmt.Errorf("loading SBOM: %w", err)
		}

		vars["sbom"] = doc
	}

	ast, err := bs.impl.Compile(bs.runner, code)
	if err != nil {
		return nil, fmt.Errorf("compiling program: %w", err)
	}

	result, err := bs.impl.Evaluate(bs.runner, ast, vars)
	if err != nil {
		return nil, fmt.Errorf("evaluating: %w", err)
	}

	return result, nil
}

func (bs *BomShell) LoadSBOM(path string) (*sbom.Document, error) {
	f, err := bs.impl.OpenFile(bs.Options.SBOM)
	if err != nil {
		return nil, fmt.Errorf("opening SBOM file: %w", err)
	}

	doc, err := bs.impl.LoadSBOM(f)
	if err != nil {
		return nil, fmt.Errorf("loading SBOM: %w", err)
	}

	return doc, nil
}

// PrintResult writes result into writer w according to the format
// configured in the options
func (bs *BomShell) PrintResult(result ref.Val, w io.WriteCloser) error {
	// TODO(puerco): Check if result is an error
	if result == nil {
		fmt.Fprint(w, "<nil>")
	}

	switch result.Type() {
	case elements.DocumentTypeValue:
		return bs.impl.PrintDocumentResult(bs.Options, result, w)
	default:
		fmt.Printf("TMPRENDER:\nvalue: %v (%T)\n", result.Value(), result)
		return nil
	}
}
