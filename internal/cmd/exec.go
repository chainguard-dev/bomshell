package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/chainguard-dev/bomshell/pkg/shell"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func execCommand() *cobra.Command {
	type execOpts = struct {
		commandLineOptions
		sboms []string
	}
	opts := &execOpts{
		sboms: []string{},
	}
	execCmd := &cobra.Command{
		PersistentPreRunE: initLogging,
		Short:             "bomshell exec program.cel [sbom.spdx.json]",
		Long: `bomshell exec recipe.cel sbom.spdx.json → Execute a bomshell program

The exec subcommand executes a cell program in a file and outputs the result.
It can optionally load an SBOM into the environment and make it available to
the program statements.
`,
		Use:           "exec program.cel [sbom.spdx.json] ",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help() //nolint:errcheck
				return errors.New("no cel program specified")
			}

			bomshell, err := shell.NewWithOptions(shell.Options{
				SBOMs:  opts.sboms,
				Format: formats.Format(opts.DocumentFormat),
			})
			if err != nil {
				logrus.Fatal("creating bomshell: %w", err)
			}

			result, err := bomshell.RunFile(args[0])
			if err != nil {
				return fmt.Errorf("executing program: %w", err)
			}

			return bomshell.PrintResult(result, os.Stdout)
		},
	}

	commandLineOpts.AddFlags(execCmd)
	opts.commandLineOptions = *commandLineOpts

	execCmd.PersistentFlags().StringArrayVar(
		&opts.sboms,
		"sbom",
		opts.sboms,
		"path to one or more SBOMs to load into the bomshell environment",
	)

	return execCmd
}
