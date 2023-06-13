package main

import (
	"os"

	"github.com/chainguard-dev/bomshell/pkg/shell"
	"github.com/sirupsen/logrus"
)

func main() {
	sbomPath := os.Args[1]
	program := os.Args[2]

	bomshell, err := shell.New()
	if err != nil {
		logrus.Fatal("creating bomshell: %w", err)
	}

	if err := bomshell.OpenSBOM(sbomPath); err != nil {
		logrus.Fatalf("loading SBOM: %v", err)
	}

	if err := bomshell.Run(program); err != nil {
		logrus.Fatal(err)
	}
}
