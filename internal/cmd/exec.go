package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var longHelp = `bomshell: An SBOM Programming Interface

bomshell is a programmable shell to work with SBOM (Software Bill of Materials) 
data using CEL expressions (Common Expression Language). bomshell can query and
remix SBOM data, split data into new documents and more.

The main bomshell command can run bomshell scripts (called recipes) from files,
or from the command line positional arguments. 

  # Execute recipe.cel, preloading SBOMs
  bomshell recipe.cel sbom1.json sbom2.json ....

  # Execute a bomshell recipe from the command line
  bomshell 'sbom.Packages()' sbom.json

  # Pipe an SBOM through bomshell, extract its file:
  cat sbom.json | bomshell 'sbom.Files()'

The root command tries to be smart when looking at arguments. It will look for
the CEL code to execute in the value of the -e|--execute flag. If it is not a file,
it will check the first positional argument:

 - If arg[0] is a file, bomshell will try to run it as a recipe.
 - If it is not a file, bomshell will treat the value as a CEL recipe and run it.
 
If a recipe is defined with -e|--execute, bomshell will treat all positional
arguments as sboms to be read and preloaded into the runtime environment:

  # Run a program with --execute:
  bomshell --execute="sbom.GetNodeByID("my-package")" sbom.json

For a more predictable runner, check the bomshell run subcommand.

`

func execCommand() *cobra.Command {
	var execCmd = &cobra.Command{
		Short:         "Default execution mode (hidden)",
		Long:          longHelp,
		Use:           "exec",
		SilenceUsage:  true,
		SilenceErrors: true,
		Hidden:        true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO(puerco): Detect open STDIN

			// If no file was piped and no args, then print help and exit
			if len(args) == 0 {
				return cmd.Help()
			}

			// If the first argument is a file, then we asume it is code
			if _, err := os.Stat(args[0]); errors.Is(err, os.ErrNotExist) {
				// TODO(puerco): Run code
			}

			// Run the recipe:
			if err := runFile(commandLineOpts, args[0]); err != nil {
				return fmt.Errorf("executing recipe: %w", err)
			}

			return nil
		},
	}

	return execCmd
}
