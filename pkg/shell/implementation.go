// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bom-squad/protobom/pkg/reader"
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/bom-squad/protobom/pkg/writer"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
)

type BomshellImplementation interface {
	Compile(*Runner, string) (*cel.Ast, error)
	Evaluate(*Runner, *cel.Ast, map[string]interface{}) (ref.Val, error)
	LoadSBOM(io.ReadSeekCloser) (*elements.Document, error)
	OpenFile(path string) (*os.File, error)
	PrintDocumentResult(Options, ref.Val, io.WriteCloser) error
	ReadRecipeFile(io.Reader) (string, error)
}

type DefaultBomshellImplementation struct{}

func (di *DefaultBomshellImplementation) Compile(runner *Runner, code string) (*cel.Ast, error) {
	return runner.Compile(code)
}

func (di *DefaultBomshellImplementation) Evaluate(runner *Runner, ast *cel.Ast, variables map[string]interface{}) (ref.Val, error) {
	return runner.EvaluateAST(ast, variables)
}

func (di *DefaultBomshellImplementation) LoadSBOM(stream io.ReadSeekCloser) (*elements.Document, error) {
	r := reader.New()
	doc, err := r.ParseStream(stream)
	if err != nil {
		return nil, fmt.Errorf("parsing SBOM: %w", err)
	}

	return &elements.Document{Document: doc}, nil
}

func (di *DefaultBomshellImplementation) OpenFile(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening SBOM file: %w", err)
	}
	return f, nil
}

// PrintDocumentResult takes a document result from a bomshell query and outputs it
// as an SBOM in the format specified in the options
func (di *DefaultBomshellImplementation) PrintDocumentResult(opts Options, result ref.Val, w io.WriteCloser) error {
	protoWriter := writer.New(
		writer.WithFormat(opts.Format),
	)

	// Check to make sure the result is a document
	if result.Type() != elements.DocumentType {
		return errors.New("error printing result, value is not a document")
	}

	doc, ok := result.Value().(*sbom.Document)
	if !ok {
		return errors.New("error casting result into protobom document")
	}

	if err := protoWriter.WriteStream(doc, w); err != nil {
		return fmt.Errorf("writing document to stream: %w", err)
	}
	return nil
}

// ReadRecipeFile reads a bomshell recipe file and returns it as a string.
// This function will look for a pind-bag line at the start of the file and
// strip it if needed.
func (di *DefaultBomshellImplementation) ReadRecipeFile(f io.Reader) (string, error) {
	fileScanner := bufio.NewScanner(f)
	fileData := ""

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		if fileData == "" && strings.HasPrefix(fileScanner.Text(), "#!") {
			continue
		}
		fileData += fileScanner.Text() + "\n"
	}

	if fileData == "" {
		return fileData, errors.New("file read resulted in zero bytes")
	}

	return fileData, nil
}
