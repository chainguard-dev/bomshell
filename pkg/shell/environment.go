package shell

import (
	"fmt"

	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/chainguard-dev/bomshell/pkg/functions"
	"github.com/google/cel-go/cel"
)

// createEnvironment creates the CEL execution environment that the runner will
// use to compile and evaluate programs on the SBOM
func createEnvironment(opts *Options) (*cel.Env, error) {
	env, err := cel.NewEnv(
		cel.Types(&sbom.Document{}),
		cel.Types(&sbom.NodeList{}),

		cel.Variable("sbom",
			cel.ObjectType(protoDocumentType),
		),

		cel.Function(
			"files",
			cel.MemberOverload(
				"sbom_files_binding", []*cel.Type{cel.ObjectType(protoDocumentType)}, elements.NodeListType,
				cel.UnaryBinding(
					functions.Files,
				),
			),
		),
		cel.Function(
			"packages",
			cel.MemberOverload(
				"sbom_packages_binding", []*cel.Type{cel.ObjectType(protoDocumentType)}, elements.NodeListType,
				cel.UnaryBinding(
					functions.Packages,
				),
			),
		),
	)

	if err != nil {
		return nil, (fmt.Errorf("creating CEL environment: %w", err))
	}

	return env, nil
}
