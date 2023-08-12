// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package cmd

import (
	"fmt"

	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/chainguard-dev/bomshell/pkg/shell"
	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/log"
)

type commandLineOptions struct {
	DocumentFormat string
	NodeListFormat string
	logLevel       string
	sboms          []string
}

type execOptions struct {
	ExecLine string
}

var execOpts = &execOptions{}

var commandLineOpts = &commandLineOptions{
	DocumentFormat: string(formats.SPDX23JSON),
	NodeListFormat: "application/json",
}

func (o *commandLineOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(
		&o.DocumentFormat,
		"document-format",
		string(shell.DefaultFormat),
		"format to output generated documents",
	)

	cmd.PersistentFlags().StringVar(
		&o.NodeListFormat,
		"nodelist-format",
		commandLineOpts.NodeListFormat,
		"format to output nodelsits (SBOM fragments)",
	)

	cmd.PersistentFlags().StringArrayVar(
		&o.sboms,
		"sbom",
		commandLineOpts.sboms,
		"path to one or more SBOMs to load into the bomshell environment",
	)

	cmd.PersistentFlags().StringVar(
		&o.logLevel,
		"log-level",
		"info",
		fmt.Sprintf("the logging verbosity, either %s", log.LevelNames()),
	)
}

func (eo *execOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&eo.ExecLine,
		"exec",
		"e",
		"",
		"CEL code to execute (overrides filename)",
	)
}
