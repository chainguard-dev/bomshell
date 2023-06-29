package cmd

import (
	"errors"
	"fmt"

	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/chainguard-dev/bomshell/pkg/shell"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func execCommand() *cobra.Command {
	type execOpts = struct {
		format string
		sbom   string
	}
	opts := &execOpts{
		format: string(formats.SPDX23JSON),
		sbom:   "",
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
				SBOM:   opts.sbom,
				Format: shell.DefaultFormat,
			})
			if err != nil {
				logrus.Fatal("creating bomshell: %w", err)
			}

			result, err := bomshell.RunFile(args[0])
			if err != nil {
				return fmt.Errorf("executing program: %w", err)
			}

			if result != nil {
				fmt.Printf("value: %v (%T)\n", result.Value(), result)
			} else {
				fmt.Printf("result is nil\n")
			}
			return nil
		},
	}

	execCmd.PersistentFlags().StringVar(
		&opts.format,
		"format",
		opts.format,
		"format to output generated SBOMs",
	)

	execCmd.PersistentFlags().StringVar(
		&opts.sbom,
		"sbom",
		opts.sbom,
		"path to sbom to ingest",
	)

	return execCmd
}
