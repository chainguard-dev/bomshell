package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/log"
	"sigs.k8s.io/release-utils/version"
)

const (
	appName = "bomshell"
)

func rootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               appName + " [flags] [\"cel code\"| recipe.cel] [sbom.json]...",
		Short:             appName + " [flags] [\"cel code\"| recipe.cel] [sbom.json]...",
		Long:              longHelp,
		Example:           appName + " recipe.cel sbom.spdx.json sbom.cdx.json",
		Version:           version.GetVersionInfo().GitVersion,
		PersistentPreRunE: initLogging,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		// SilenceErrors:              false,
		SilenceUsage: false,
	}

	rootCmd.SetVersionTemplate(fmt.Sprintf("%s v{{.Version}}\n", appName))

	execOpts.AddFlags(rootCmd)
	commandLineOpts.AddFlags(rootCmd)

	rootCmd.AddCommand(execCommand())
	rootCmd.AddCommand(runCommand())
	rootCmd.AddCommand(interactiveCommand())
	rootCmd.AddCommand(version.WithFont("starwars"))

	return rootCmd
}

func Execute() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "completion", "exec", "version", "interactive", "run":
		default:
			os.Args = append([]string{os.Args[0], "exec"}, os.Args[1:]...)
		}
	}

	if err := rootCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func initLogging(*cobra.Command, []string) error {
	return log.SetupGlobalLogger(commandLineOpts.logLevel)
}
