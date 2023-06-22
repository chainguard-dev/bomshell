package shell

import (
	"fmt"

	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/chainguard-dev/bomshell/pkg/functions"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/traits"
	celfuncs "github.com/google/cel-go/interpreter/functions"
	//"github.com/google/cel-go/common/operators"
	//"github.com/google/cel-go/common/types/traits"
	//celfuncs "github.com/google/cel-go/interpreter/functions"
)

type shellLibrary struct{}

// createEnvironment creates the CEL execution environment that the runner will
// use to compile and evaluate programs on the SBOM
func (shellLibrary) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Types(&sbom.Document{}),
		cel.Types(&sbom.NodeList{}),

		cel.Variable("sbom",
			cel.ObjectType(protoDocumentType),
		),

		cel.Function(
			"files",
			cel.MemberOverload(
				"sbom_files_binding", []*cel.Type{cel.ObjectType(protoDocumentType)}, elements.NodeListType,
				cel.UnaryBinding(functions.Files),
			),
		),

		cel.Function(
			"packages",
			cel.MemberOverload(
				"sbom_packages_binding", []*cel.Type{cel.ObjectType(protoDocumentType)}, elements.NodeListType,
				cel.UnaryBinding(functions.Packages),
			),
		),

		cel.Function(
			"add",
			cel.MemberOverload(
				"add_nodelists",
				[]*cel.Type{elements.NodeListType, elements.NodeListType},
				elements.NodeListType,
				cel.BinaryBinding(functions.Addition),
				// cel.OverloadOperandTrait(traits.AdderType),
			),
		),

		cel.Function(
			"tonodelist",
			cel.MemberOverload(
				"node_tonodelist_binding",
				[]*cel.Type{elements.NodeType},
				elements.NodeListType,
				cel.UnaryBinding(functions.NodeToNodeList),
				// cel.OverloadOperandTrait(traits.AdderType),
			),
		),

		/*
			cel.Function(
				operators.Add,
				cel.MemberOverload(
					"add_nodelists",
					[]*cel.Type{elements.NodeListType, elements.NodeListType},
					elements.NodeListType,
					cel.BinaryBinding(functions.Addition),

					cel.OverloadOperandTrait(traits.AdderType),
				),
			),
		*/
		cel.Function(
			"elementbyid",
			cel.MemberOverload(
				"sbom_elementbyid_binding", []*cel.Type{cel.ObjectType(protoDocumentType), cel.StringType}, elements.NodeType,
				cel.BinaryBinding(functions.ElementById),
			),
		),
	}
}
func (shellLibrary) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{
		// cel.Functions(functions.StandardOverloads()...),

		cel.Functions(
			&celfuncs.Overload{
				Operator:     "++", /// Placegholder while I figure out how to overload operators.Add
				OperandTrait: traits.AdderType,
				Binary:       functions.Addition,
				// Function:     functions.AdditionOp,
				//NonStrict: false,
			},
		),
	}

	// return []cel.ProgramOption{}
}

func Library() cel.EnvOption {
	return cel.Lib(shellLibrary{})
}

func createEnvironment(opts *Options) (*cel.Env, error) {
	shlib := Library()
	env, err := cel.NewEnv(shlib)
	if err != nil {
		return nil, (fmt.Errorf("creating CEL environment: %w", err))
	}

	return env, nil
}
