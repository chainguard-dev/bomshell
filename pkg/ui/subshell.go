package ui

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

type InteractiveSubshell struct{}

// func buildSubshell() InteractiveSubshell {}

func helpFunc(_ ...ref.Val) ref.Val {
	h := "Welcome to bomshell!\n"
	h += "--------------------\n"
	h += "This is the interactive mode of bomshell. In here you can run bomshell\n"
	h += "recipes on any loaded SBOMs. The interactive mode will render the\n"
	h += "evaluation results using a special human-readable renderer that outputs\n"
	h += "information about the returned types. For example, when a recipe \n"
	h += "evaluates to a NodeList, bomshell will show some data about it:\n\n"
	h += "   protobom NodeList\n"
	h += "   Root Elements: 1\n"
	h += "   Number of nodes: 69 (14 packages 55 files)\n"
	h += "   Package URL types: oci: 2 apk: 12 \n\n"
	h += "Each element (document, nodelist, node, etc) has its own information\n"
	h += "display. They are supposed to be read nby humans and we expect to\n"
	h += "modify them constantly with each release of bomshell. In other words,\n"
	h += "don't script based on the human (TTY) display\n\n"
	h += "Examples!\n"
	h += "If you have SBOMs loaded in the environment, try pasting these examples\n"
	h += "in the prompt below:\n\n"
	h += "// Print information about the first SBOM:\n"
	h += "sboms[0]\n\n"
	h += "// Query the SBOM and print data about its packages:\n"
	h += "sboms[0].packages()\n\n"
	h += "// Query a specific node in the document:\n"
	h += `sboms[0].NodeByID("my-node-identifier")` + "\n\n"

	return types.String(h)
}

func (subshell InteractiveSubshell) LibraryName() string {
	return "cel.chainguard.bomshell.interactive"
}

// CompileOptions
func (subshell InteractiveSubshell) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Function(
			"help",
			cel.Overload(
				"help_overload", nil, cel.StringType, cel.FunctionBinding(helpFunc),
			),
		),
	}
}

func (subshell InteractiveSubshell) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}
