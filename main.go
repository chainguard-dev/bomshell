package main

import (
	"os"

	"github.com/chainguard-dev/bomshell/pkg/shell"
	"github.com/sirupsen/logrus"
)

func main() {
	sbomPath := os.Args[1]
	program := os.Args[2]

	bomshell, err := shell.NewWithOptions(shell.Options{
		SBOM:   sbomPath,
		Format: shell.DefaultFormat,
	})
	if err != nil {
		logrus.Fatal("creating bomshell: %w", err)
	}

	if err := bomshell.RunFile(program); err != nil {
		logrus.Fatal(err)
	}
}
