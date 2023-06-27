package loader

import (
	"fmt"
	"io"
	"os"

	"github.com/bom-squad/protobom/pkg/reader"
	"github.com/bom-squad/protobom/pkg/sbom"
)

func ReadSBOM(stream io.ReadSeekCloser) (*sbom.Document, error) {
	r := reader.New()
	doc, err := r.ParseReader(stream)
	if err != nil {
		return nil, fmt.Errorf("parsing SBOM: %w", err)
	}

	return doc, nil
}

func OpenFile(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening SBOM file: %w", err)
	}
	return f, nil

}
