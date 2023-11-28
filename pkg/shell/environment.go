// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package shell

import (
	"fmt"

	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/chainguard-dev/bomshell/pkg/functions"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	// "github.com/google/cel-go/common/operators"
	// "github.com/google/cel-go/common/types/traits"
	// celfuncs "github.com/google/cel-go/interpreter/functions"
)

type shellLibrary struct{}

// createEnvironment creates the CEL execution environment that the runner will
// use to compile and evaluate programs on the SBOM
func (shellLibrary) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable("sboms", cel.MapType(cel.IntType, elements.DocumentType)),
		cel.Variable("sbom", elements.DocumentType),
		cel.Variable("bomshell", elements.BomshellType),

		cel.Function(
			"files",
			cel.MemberOverload(
				"sbom_files_binding", []*cel.Type{cel.ObjectType(protoDocumentType)}, elements.NodeListType,
				cel.UnaryBinding(functions.Files),
			),
			cel.MemberOverload(
				"nodelist_files_binding", []*cel.Type{elements.NodeListType}, elements.NodeListType,
				cel.UnaryBinding(functions.Files),
			),
			cel.MemberOverload(
				"node_files_binding", []*cel.Type{elements.NodeType}, elements.NodeListType,
				cel.UnaryBinding(functions.Files),
			),
		),

		cel.Function(
			"packages",
			cel.MemberOverload(
				"sbom_packages_binding", []*cel.Type{cel.ObjectType(protoDocumentType)}, elements.NodeListType,
				cel.UnaryBinding(functions.Packages),
			),
			cel.MemberOverload(
				"nodeslist_packages_binding", []*cel.Type{elements.NodeListType}, elements.NodeListType,
				cel.UnaryBinding(functions.Packages),
			),
			cel.MemberOverload(
				"node_packages_binding", []*cel.Type{elements.NodeType}, elements.NodeListType,
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
			"ToNodeList",
			cel.MemberOverload(
				"document_tonodelist_binding",
				[]*cel.Type{cel.ObjectType(protoDocumentType)}, elements.NodeListType,
				cel.UnaryBinding(functions.ToNodeList),
			),
			cel.MemberOverload(
				"nodelist_tonodelist_binding",
				[]*cel.Type{elements.NodeListType}, elements.NodeListType,
				cel.UnaryBinding(functions.ToNodeList),
			),
			cel.MemberOverload(
				"node_tonodelist_binding",
				[]*cel.Type{elements.NodeType}, elements.NodeListType,
				cel.UnaryBinding(functions.ToNodeList),
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
				"sbom_nodebyid_binding", []*cel.Type{elements.DocumentType, cel.StringType}, elements.NodeType,
				cel.BinaryBinding(functions.NodeByID),
			),
			cel.MemberOverload(
				"nodelist_nodebyid_binding", []*cel.Type{elements.NodeListType, cel.StringType}, elements.NodeType,
				cel.BinaryBinding(functions.NodeByID),
			),
		),

		cel.Function(
			"NodesByPurlType",
			cel.MemberOverload(
				"sbom_nodesbypurltype_binding", []*cel.Type{elements.DocumentType, cel.StringType}, elements.NodeListType,
				cel.BinaryBinding(functions.NodesByPurlType),
			),
			cel.MemberOverload(
				"nodelist_nodesbypurltype_binding", []*cel.Type{elements.NodeListType, cel.StringType}, elements.NodeListType,
				cel.BinaryBinding(functions.NodesByPurlType),
			),
		),

		cel.Function(
			"ToDocument",
			cel.MemberOverload(
				"document_todocument_binding",
				[]*cel.Type{elements.DocumentType}, elements.DocumentType,
				cel.UnaryBinding(functions.ToDocument),
			),
			cel.MemberOverload(
				"nodelist_todocument_binding",
				[]*cel.Type{elements.NodeListType}, elements.DocumentType,
				cel.UnaryBinding(functions.ToDocument),
			),
			cel.MemberOverload(
				"node_todocument_binding",
				[]*cel.Type{elements.NodeType}, elements.DocumentType,
				cel.UnaryBinding(functions.ToDocument),
			),
		),

		cel.Function(
			"LoadSBOM",
			cel.MemberOverload(
				"bomshell_loadsbom_binding",
				[]*cel.Type{elements.BomshellType, cel.StringType}, elements.DocumentType,
				cel.BinaryBinding(functions.LoadSBOM),
			),
		),

		cel.Function(
			"RelateNodeListAtID",
			cel.MemberOverload(
				"sbom_relatenodesatid_binding",
				[]*cel.Type{elements.DocumentType, elements.NodeListType, cel.StringType, cel.StringType},
				elements.DocumentType, // result
				cel.FunctionBinding(functions.RelateNodeListAtID),
			),
			cel.MemberOverload(
				"nodelist_relatenodesatid_binding",
				[]*cel.Type{elements.NodeListType, elements.NodeListType, cel.StringType, cel.StringType},
				elements.DocumentType, // result
				cel.FunctionBinding(functions.RelateNodeListAtID),
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

type customTypeAdapter struct{}

func (customTypeAdapter) NativeToValue(value interface{}) ref.Val {
	val, ok := value.(elements.Bomshell)
	if ok {
		return val
	} else {
		// let the default adapter handle other cases
		return types.DefaultTypeAdapter.NativeToValue(value)
	}
}

func createEnvironment(opts *Options) (*cel.Env, error) {
	envOpts := []cel.EnvOption{
		cel.CustomTypeAdapter(&customTypeAdapter{}),
		Library(),
		ext.Bindings(),
		ext.Strings(),
		ext.Encoders(),
	}

	// Add any additional environment options passed in the construcutor
	envOpts = append(envOpts, opts.EnvOptions...)
	env, err := cel.NewEnv(
		envOpts...,
	)
	if err != nil {
		return nil, (fmt.Errorf("creating CEL environment: %w", err))
	}

	return env, nil
}
