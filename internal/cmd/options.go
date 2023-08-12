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

var commandLineOpts = &commandLineOptions{
	DocumentFormat: string(formats.SPDX23JSON),
	NodeListFormat: "application/json",
}

func (o *commandLineOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(
		&o.DocumentFormat,
		"document-format",
		string(shell.DefaultFormat),
		fmt.Sprintf("format to output generated documents"),
	)

	cmd.PersistentFlags().StringVar(
		&o.NodeListFormat,
		"nodelist-format",
		commandLineOpts.NodeListFormat,
		fmt.Sprintf("format to output nodelsits (SBOM fragments)"),
	)

	cmd.PersistentFlags().StringArrayVar(
		&o.sboms,
		"sbom",
		commandLineOpts.sboms,
		"path to one or more SBOMs to load into the bomshell environment",
	)

	cmd.PersistentFlags().StringVar(
		&commandLineOpts.logLevel,
		"log-level",
		"info",
		fmt.Sprintf("the logging verbosity, either %s", log.LevelNames()),
	)
}
