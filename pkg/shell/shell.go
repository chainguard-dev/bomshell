package shell

import (
	"fmt"
	"os"

	"github.com/bom-squad/protobom/pkg/formats"
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

func (bs *BomShell) RunFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading program data: %w", err)
	}
	return bs.Run(string(data))
}

func (bs *BomShell) Run(code string) error {
	// Variables that wil be made available in the CEL env
	vars := map[string]interface{}{}

	// Load an SBNOM if defined
	if bs.Options.SBOM != "" {
		f, err := bs.impl.OpenFile(bs.Options.SBOM)
		if err != nil {
			return fmt.Errorf("opening SBOM file: %w", err)
		}

		doc, err := bs.impl.LoadSBOM(f)
		if err != nil {
			return fmt.Errorf("loading SBOM: %w", err)
		}

		vars["sbom"] = doc
	}

	ast, err := bs.impl.Compile(bs.runner, code)
	if err != nil {
		return fmt.Errorf("compiling program: %w", err)
	}

	result, err := bs.impl.Evaluate(bs.runner, ast, vars)
	if err != nil {
		return fmt.Errorf("evaluating: %w", err)
	}

	if result != nil {
		fmt.Printf("value: %v (%T)\n", result.Value(), result)
	} else {
		fmt.Printf("result is nil\n")
	}
	return nil
}
