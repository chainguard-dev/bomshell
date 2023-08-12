package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/log"
)

func rootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "bomshell [flags] [\"cel code\"| recipe.cel] [sbom.json]...",
		Short:   "bomshell [flags] [\"cel code\"| recipe.cel] [sbom.json]...",
		Long:    longHelp,
		Example: "bomshell recipe.cel sbom.spdx.json sbom.cdx.json",
		//Deprecated:        "",
		//Version:           "",
		PersistentPreRunE: initLogging,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		// SilenceErrors:              false,
		SilenceUsage: false,
	}

	rootCmd.AddCommand(execCommand())
	rootCmd.AddCommand(runCommand())
	rootCmd.AddCommand(interactiveCommand())

	return rootCmd
}

func Execute() {
	root := rootCommand()

	if len(os.Args) > 1 {
		pcmd := os.Args[1]
		if pcmd == "completion" || pcmd == "--version" || pcmd == "exec" {
			return
		}
		for _, command := range root.Commands() {
			if command.Use == pcmd {
				return
			}
		}
		os.Args = append([]string{os.Args[0], "exec"}, os.Args[1:]...)
	}

	if err := root.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func initLogging(*cobra.Command, []string) error {
	return log.SetupGlobalLogger(commandLineOpts.logLevel)
}
