// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package shell

import (
	"fmt"
	"io"

	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/google/cel-go/common/types/ref"
	"github.com/sirupsen/logrus"
)

const (
	protoDocumentType = "bomsquad.protobom.Document"
	DefaultFormat     = formats.SPDX23JSON
)

type Options struct {
	SBOMs  []string
	Format formats.Format
}

var defaultOptions = Options{
	Format: DefaultFormat,
}

type Bomshell struct {
	Options Options
	runner  *Runner
	impl    BomshellImplementation
}

func New() (*Bomshell, error) {
	return NewWithOptions(defaultOptions)
}

func NewWithOptions(opts Options) (*Bomshell, error) {
	runner, err := NewRunnerWithOptions(&opts)
	if err != nil {
		return nil, fmt.Errorf("creating runner: %w", err)
	}
	return &Bomshell{
		Options: opts,
		runner:  runner,
		impl:    &DefaultBomshellImplementation{},
	}, nil
}

// RunFile runs a bomshell recipe from a file
func (bs *Bomshell) RunFile(path string) (ref.Val, error) {
	f, err := bs.impl.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	defer f.Close()

	data, err := bs.impl.ReadRecipeFile(f)
	if err != nil {
		return nil, fmt.Errorf("reading program data: %w", err)
	}

	return bs.Run(string(data))
}

func (bs *Bomshell) Run(code string) (ref.Val, error) {
	// Variables that wil be made available in the CEL env
	vars := map[string]interface{}{}
	sbomList := []*elements.Document{}

	// Load defined SBOMs into the sboms array
	if len(bs.Options.SBOMs) > 0 {
		for _, sbomSpec := range bs.Options.SBOMs {
			// TODO(puerco): Split for varname
			f, err := bs.impl.OpenFile(sbomSpec)
			if err != nil {
				return nil, fmt.Errorf("opening SBOM file: %w", err)
			}

			doc, err := bs.impl.LoadSBOM(f)
			if err != nil {
				return nil, fmt.Errorf("loading SBOM: %w", err)
			}
			logrus.Debugf("Loaded %s", sbomSpec)

			sbomList = append(sbomList, doc)
		}
	}

	// Add the SBOM list to the runtim environment
	vars["sboms"] = sbomList
	if len(sbomList) > 0 {
		vars["sbom"] = sbomList[0]
	}
	// Add the default bomshell object
	vars["bomshell"] = elements.Bomshell{}

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

func (bs *Bomshell) LoadSBOM(path string) (*elements.Document, error) {
	f, err := bs.impl.OpenFile(path)
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
func (bs *Bomshell) PrintResult(result ref.Val, w io.WriteCloser) error {
	// TODO(puerco): Check if result is an error
	if result == nil {
		fmt.Fprint(w, "<nil>")
	}

	switch result.Type() {
	case elements.DocumentType:
		return bs.impl.PrintDocumentResult(bs.Options, result, w)
	default:
		fmt.Printf("TMPRENDER:\nvalue: %v (%T)\n", result.Value(), result)
		return nil
	}
}
