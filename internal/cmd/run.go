package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/chainguard-dev/bomshell/pkg/shell"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/version"
)

func runCommand() *cobra.Command {
	runCmd := &cobra.Command{
		PersistentPreRunE: initLogging,
		Short:             "Run bomshell recipe files",
		Example:           "bomshell run program.cel [sbom.spdx.json]...",
		Long: appName + ` run recipe.cel sbom.spdx.json → Execute a bomshell program

The exec subcommand executes a cell program in a file and outputs the result.
It can optionally load an SBOM into the environment and make it available to
the program statements.
`,
		Use:           "run",
		Version:       version.GetVersionInfo().GitVersion,
		SilenceUsage:  false,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help() //nolint:errcheck
				return errors.New("no cel program specified")
			}

			sbomPaths := []string{}
			if len(args) > 1 {
				sbomPaths = append(sbomPaths, args[1:]...)
			}

			return runFile(commandLineOpts, args[0], append(sbomPaths, commandLineOpts.sboms...))
		},
	}
	execOpts.AddFlags(runCmd)
	commandLineOpts.AddFlags(runCmd)

	return runCmd
}

func buildShell(opts *commandLineOptions, sbomList []string) (*shell.Bomshell, error) {
	bomshell, err := shell.NewWithOptions(shell.Options{
		SBOMs:  sbomList,
		Format: formats.Format(opts.DocumentFormat),
	})
	if err != nil {
		logrus.Fatalf("creating bomshell: %v", err)
	}
	return bomshell, nil
}

// runFile creates and configures a bomshell instance to run a recipe from a file
func runFile(opts *commandLineOptions, recipePath string, sbomList []string) error {
	bomshell, err := buildShell(opts, sbomList)
	if err != nil {
		return err
	}

	result, err := bomshell.RunFile(recipePath)
	if err != nil {
		return fmt.Errorf("executing program: %w", err)
	}

	return bomshell.PrintResult(result, os.Stdout)
}

// runCode creates and configures a bomshell instance to run a recipe from a string
func runCode(opts *commandLineOptions, celCode string, sbomList []string) error {
	bomshell, err := buildShell(opts, sbomList)
	if err != nil {
		return err
	}

	result, err := bomshell.Run(celCode)
	if err != nil {
		return fmt.Errorf("executing program: %w", err)
	}

	return bomshell.PrintResult(result, os.Stdout)
}
