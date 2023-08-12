// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/version"
)

var longHelp = `💣🐚 bomshell: An SBOM Programming Interface

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
		Version:       version.GetVersionInfo().GitVersion,
		Use:           "exec",
		SilenceUsage:  true,
		SilenceErrors: true,
		Hidden:        true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// List of SBOMs to prelaod
			sbomPaths := []string{}

			// If there is an SBOM piped through STDIN, it will always
			// be SBOM #0 in our list
			stdinSBOM, err := testStdin()
			if err != nil {
				return fmt.Errorf("checking STDIN for a piped SBOM")
			}

			if stdinSBOM != "" {
				sbomPaths = append(sbomPaths, stdinSBOM)
				defer os.Remove(stdinSBOM)
			}

			// If no file was piped and no args, then print help and exit
			if len(args) == 0 && execOpts.ExecLine == "" {
				return cmd.Help()
			}

			// Case 1: Run snippet from the --execute flag
			if execOpts.ExecLine != "" {
				sbomPaths = append(sbomPaths, args...)
				sbomPaths = append(sbomPaths, commandLineOpts.sboms...)
				if err := runCode(commandLineOpts, execOpts.ExecLine, sbomPaths); err != nil {
					return fmt.Errorf("running code snippet: %w", err)
				}
				return nil
			}

			// The next two cases take SBOMs from arg[1] on
			if len(args) > 1 {
				sbomPaths = append(sbomPaths, args[1:]...)
			}

			// Case 2: If the first argument is not a file, then we asume it is code
			if _, err := os.Stat(args[0]); errors.Is(err, os.ErrNotExist) {
				// TODO(puerco): Implemnent code to check if args[0] is code :D
				if err := runCode(commandLineOpts, args[0], append(sbomPaths, commandLineOpts.sboms...)); err != nil {
					return fmt.Errorf("running code snippet: %w", err)
				}
				return nil
			}

			// Case 3: First argument is the recipe file
			if err := runFile(commandLineOpts, args[0], append(sbomPaths, commandLineOpts.sboms...)); err != nil {
				return fmt.Errorf("executing recipe: %w", err)
			}
			return nil
		},
	}

	execOpts.AddFlags(execCmd)
	commandLineOpts.AddFlags(execCmd)
	return execCmd
}
