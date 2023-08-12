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

func runCommand() *cobra.Command {
	runCmd := &cobra.Command{
		PersistentPreRunE: initLogging,
		Short:             "Run bomshell recipe files",
		Example:           "bomshell run program.cel [sbom.spdx.json]...",
		Long: `bomshell exec recipe.cel sbom.spdx.json → Execute a bomshell program

The exec subcommand executes a cell program in a file and outputs the result.
It can optionally load an SBOM into the environment and make it available to
the program statements.
`,
		Use:           "run program.cel [sbom.spdx.json] ",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help() //nolint:errcheck
				return errors.New("no cel program specified")
			}

			return runFile(commandLineOpts, args[0])
		},
	}

	commandLineOpts.AddFlags(runCmd)

	return runCmd
}

func buildShell(opts *commandLineOptions) (*shell.Bomshell, error) {
	bomshell, err := shell.NewWithOptions(shell.Options{
		SBOMs:  opts.sboms,
		Format: formats.Format(opts.DocumentFormat),
	})
	if err != nil {
		logrus.Fatalf("creating bomshell: %v", err)
	}
	return bomshell, nil
}

func runFile(opts *commandLineOptions, recipePath string) error {
	bomshell, err := buildShell(opts)
	if err != nil {
		return err
	}

	result, err := bomshell.RunFile(recipePath)
	if err != nil {
		return fmt.Errorf("executing program: %w", err)
	}

	return bomshell.PrintResult(result, os.Stdout)
}
