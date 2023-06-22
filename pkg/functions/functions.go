package functions

import (
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// NodeToNodeList takes a node and returns a new NodeList
// with that nodelist with the node as the only member.
var NodeToNodeList = func(lhs ref.Val) ref.Val {
	var node elements.Node
	var ok bool
	if node, ok = lhs.(elements.Node); !ok {
		return types.NewErr("attemtp to convert a non node")
	}
	return node.ToNodeList()
}

var Addition = func(lhs, rhs ref.Val) ref.Val {
	return elements.NodeList{
		NodeList: &sbom.NodeList{},
	}
}

var AdditionOp = func(vals ...ref.Val) ref.Val {
	return elements.NodeList{
		NodeList: &sbom.NodeList{},
	}
}

// ElementById returns a Node matching with the specified ID
var ElementById = func(lhs, rawID ref.Val) ref.Val {
	queryID, ok := rawID.Value().(string)
	if !ok {
		return types.NewErr("argument to element by id has to be a string")
	}

	bom, ok := lhs.Value().(*sbom.Document)
	if !ok {
		return types.NewErr("unable to convert sbom to native (wrong type?)")
	}
	for _, n := range bom.Nodes {
		if n.Id == queryID {
			return elements.Node{
				Node: n,
			}
		}
	}
	return nil
}

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
	cleanEdges(&nl)
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
	cleanEdges(&nl)
	return nl
}
