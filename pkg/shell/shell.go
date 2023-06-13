package shell

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bom-squad/protobom/pkg/reader"
	"github.com/bom-squad/protobom/pkg/sbom"
)

const (
	protoDocumentType = "bomsquad.protobom.Document"
)

type Options struct {
	SBOM *sbom.Document
}

var defaultOptions = Options{}

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

func (bs *BomShell) LoadSBOM(stream io.ReadSeekCloser) error {
	r := reader.New()
	doc, err := r.ParseReader(stream)
	if err != nil {
		return fmt.Errorf("parsing SBOM: %w", err)
	}

	bs.Options.SBOM = doc
	return nil
}

func (bs *BomShell) OpenSBOM(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening SBOM file: %w", err)
	}
	defer f.Close()
	return bs.LoadSBOM(f)
}

func (bs *BomShell) RunFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading program data: %w", err)
	}
	return bs.Run(string(data))
}

func (bs *BomShell) Run(code string) error {
	if bs.Options.SBOM == nil {
		return errors.New("unable to run code, no SBOM has been loaded")
	}

	ast, err := bs.runner.Compile(code)
	if err != nil {
		return fmt.Errorf("compiling program: %w", err)
	}

	result, err := bs.runner.EvaluateAST(ast, bs.Options.SBOM)
	if err != nil {
		return err
	}
	fmt.Printf("value: %v (%T)\n", result.Value(), result)
	return nil
}
