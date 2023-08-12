package cmd

import (
	"fmt"

	"github.com/chainguard-dev/bomshell/pkg/ui"
	"github.com/spf13/cobra"
)

func interactiveCommand() *cobra.Command {
	type interactiveOpts = struct {
		commandLineOptions
		sboms []string
	}
	opts := &interactiveOpts{
		sboms: []string{},
	}
	execCmd := &cobra.Command{
		PersistentPreRunE: initLogging,
		Short:             "Launch bomshell interactive workbench",
		Long: `bomshell interactive sbom.spdx.json → Launch the bomshell interactive workbench

The interactive subcommand launches the bomshell interactive workbench
`,
		Use:           "interactive [sbom.spdx.json...] ",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return launchInteractive(commandLineOpts)
		},
	}

	commandLineOpts.AddFlags(execCmd)
	opts.commandLineOptions = *commandLineOpts

	return execCmd
}

func launchInteractive(opts *commandLineOptions) error {
	i, err := ui.NewInteractive()
	if err != nil {
		return fmt.Errorf("creating interactive env: %w", err)
	}

	// Start the interactive shell
	if err := i.Start(); err != nil {
		return fmt.Errorf("executing interactive mode: %w", err)
	}
	return nil
}
