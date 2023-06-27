package shell

import (
	"fmt"

	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/chainguard-dev/bomshell/pkg/functions"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	//"github.com/google/cel-go/common/operators"
	//"github.com/google/cel-go/common/types/traits"
	//celfuncs "github.com/google/cel-go/interpreter/functions"
)

type shellLibrary struct{}

// createEnvironment creates the CEL execution environment that the runner will
// use to compile and evaluate programs on the SBOM
func (shellLibrary) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Types(&elements.Document{}),
		cel.Types(&elements.NodeList{}),

		cel.Variable("sbom",
			// cel.ObjectType(protoDocumentType),
			elements.DocumentType,
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
			"NodeByID",
			cel.MemberOverload(
				"sbom_elementbyid_binding", []*cel.Type{elements.DocumentType, cel.StringType}, elements.NodeType,
				cel.BinaryBinding(functions.NodeById),
			),
			cel.MemberOverload(
				"nodelist_elementbyid_binding", []*cel.Type{elements.NodeListType, cel.StringType}, elements.NodeType,
				cel.BinaryBinding(functions.NodeById),
			),
		),

		cel.Function(
			"ToDocument",
			cel.MemberOverload(
				"nodelist_todocument_binding", []*cel.Type{elements.NodeListType}, elements.DocumentType,
				cel.UnaryBinding(functions.ToDocument),
			),
		),

		cel.Function(
			"LoadSBOM",
			cel.Overload(
				"global_load_sbom_binding",
				[]*cel.Type{cel.StringType},
				elements.DocumentType,
				cel.UnaryBinding(functions.LoadSBOM),
				//cel.BinaryBinding(functions.LoadSBOM),
			),
		),
		/*
			cel.Macros(
				// cel.bind(var, <init>, <expr>)
				cel.NewReceiverMacro("LoadSBOM", 1, celBind),
			),
		*/
	}
}

func (shellLibrary) LibraryName() string {
	return "cel.chainguard.bomshell"
}

func (shellLibrary) ProgramOptions() []cel.ProgramOption {
	/*
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
	*/
	return []cel.ProgramOption{}
}

func Library() cel.EnvOption {
	return cel.Lib(shellLibrary{})
}

func createEnvironment(opts *Options) (*cel.Env, error) {
	env, err := cel.NewEnv(
		Library(),
		ext.Bindings(),
		ext.Strings(),
		ext.Encoders(),
	)
	if err != nil {
		return nil, (fmt.Errorf("creating CEL environment: %w", err))
	}

	return env, nil
}
