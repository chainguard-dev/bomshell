package cmd

import (
	"fmt"

	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/chainguard-dev/bomshell/pkg/shell"
	"github.com/chainguard-dev/bomshell/pkg/ui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/log"
)

type commandLineOptions struct {
	DocumentFormat string
	NodeListFormat string
	logLevel       string
}

var commandLineOpts = &commandLineOptions{
	DocumentFormat: string(formats.SPDX23JSON),
	NodeListFormat: "application/json",
}

func rootCommand() *cobra.Command {
	var rootCmd = &cobra.Command{
		Short: "A programmable shell to work with SBOM data",
		Long: `bomshell

bomshell is a programmable shell that allows to work with SBOM data
using CEL (Common Expression Language) expressions. bomshell can query
and remix SBOM data, split data into new documents and more.
`,
		Use:               "bomshell",
		SilenceUsage:      false,
		PersistentPreRunE: initLogging,
		RunE: func(md *cobra.Command, args []string) error {
			i, err := ui.NewInteractive()
			if err != nil {
				return fmt.Errorf("creating interactive env: %w", err)
			}

			// Start the interactive shell
			if err := i.Start(); err != nil {
				return fmt.Errorf("executing interactive mode: %w", err)
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVar(
		&commandLineOpts.logLevel,
		"log-level",
		"info",
		fmt.Sprintf("the logging verbosity, either %s", log.LevelNames()),
	)

	rootCmd.AddCommand(execCommand())

	return rootCmd
}

func Execute() {
	if err := rootCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func initLogging(*cobra.Command, []string) error {
	return log.SetupGlobalLogger(commandLineOpts.logLevel)
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
}
