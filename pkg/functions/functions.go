package functions

import (
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

var Files = func(lhs ref.Val) ref.Val {
	nl := elements.NodeList{
		NodeList: &sbom.NodeList{},
	}
	bom, ok := lhs.Value().(*sbom.Document)
	if !ok {
		return types.NewErr("unable to convert sbom to native (wrong type?)")
	}
	nl.Edges = bom.Edges
	for _, n := range bom.Nodes {
		if n.Type == sbom.Node_FILE {
			nl.NodeList.Nodes = append(nl.NodeList.Nodes, n)
		}
	}
	return nl
}

var Packages = func(lhs ref.Val) ref.Val {
	nl := elements.NodeList{
		NodeList: &sbom.NodeList{},
	}
	bom, ok := lhs.Value().(*sbom.Document)
	if !ok {
		return types.NewErr("unable to convert sbom to native (wrong type?)")
	}
	nl.Edges = bom.Edges
	for _, n := range bom.Nodes {
		if n.Type == sbom.Node_PACKAGE {
			nl.NodeList.Nodes = append(nl.NodeList.Nodes, n)
		}
	}
	return nl
}
