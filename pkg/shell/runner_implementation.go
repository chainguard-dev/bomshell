// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package shell

import (
	"fmt"
	"io"

	"github.com/google/cel-go/cel"
)

type RunnerImplementation interface {
	ReadStream(io.Reader) (string, error)
	Compile(*cel.Env, string) (*cel.Ast, error)
}

type defaultRunnerImplementation struct{}

func (dri *defaultRunnerImplementation) ReadStream(reader io.Reader) (string, error) {
	// Read all the stream into a string
	contents, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("reading stram code: %w", err)
	}
	return string(contents), nil
}

func (dri *defaultRunnerImplementation) Compile(env *cel.Env, code string) (*cel.Ast, error) {
	// Run the compilation step
	ast, iss := env.Compile(code)
	if iss.Err() != nil {
		return nil, fmt.Errorf("compilation error: %w", iss.Err())
	}
	return ast, nil
}

// func (dri *defaultRunnerImplementation) Evaluate(env )
