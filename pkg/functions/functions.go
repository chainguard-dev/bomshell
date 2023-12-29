// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package functions

import (
	"fmt"
	"os"

	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/chainguard-dev/bomshell/pkg/loader"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sigs.k8s.io/release-utils/version"
)

// ToNodeList takes a node and returns a new NodeList
// with that nodelist with the node as the only member.
var ToNodeList = func(lhs ref.Val) ref.Val {
	switch v := lhs.Value().(type) {
	case *sbom.Document:
		return elements.NodeList{
			NodeList: v.NodeList,
		}
	case *elements.Document:
		return elements.NodeList{
			NodeList: v.Document.NodeList,
		}
	case *sbom.NodeList:
		return elements.NodeList{
			NodeList: v,
		}
	case *elements.NodeList:
		return v
	case *elements.Node:
		nl := v.ToNodeList()
		return *nl
	case *sbom.Node:
		nl := elements.Node{
			Node: v,
		}.ToNodeList()
		return *nl
	default:
		return types.NewErr("type does not support conversion to NodeList" + fmt.Sprintf("%T", v))
	}
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

// NodeByID returns a Node matching the specified ID
var NodeByID = func(lhs, rawID ref.Val) ref.Val {
	queryID, ok := rawID.Value().(string)
	if !ok {
		return types.NewErr("argument to element by id has to be a string")
	}
	var node *sbom.Node
	switch v := lhs.Value().(type) {
	case *sbom.Document:
		node = v.NodeList.GetNodeByID(queryID)
	case *sbom.NodeList:
		node = v.GetNodeByID(queryID)
	default:
		return types.NewErr("method unsupported on type %T", lhs.Value())
	}

	if node == nil {
		return nil
	}

	return elements.Node{
		Node: node,
	}
}

// Files returns all the Nodes marked as type file from an element. The function
// supports documents, nodelists, and nodes. If the node is a file, it will return
// a NodeList with it as the single node or empty if it is a package.
//
// If the passed type is not supported, the return value will be an error.
var Files = func(lhs ref.Val) ref.Val {
	nodeList, err := getTypedNodes(lhs, sbom.Node_FILE)
	if err != nil {
		return types.NewErr(err.Error())
	}
	return nodeList
}

// Packages returns a NodeList with any packages in the lhs element. It supports
// Documents, NodeLists and Nodes. If a node is provided it will return a NodeList
// with the single node it is a package, otherwise it will be empty.
//
// If lhs is an unsupprted type, Packages will return an error.
var Packages = func(lhs ref.Val) ref.Val {
	nodeList, err := getTypedNodes(lhs, sbom.Node_PACKAGE)
	if err != nil {
		return types.NewErr(err.Error())
	}
	return nodeList
}

// getTypedNodes takes an element and returns a nodelist containing all nodes
// of the specified type (package or file). If an unsupported types is provided,
// the function return an error
func getTypedNodes(element ref.Val, t sbom.Node_NodeType) (elements.NodeList, error) {
	var sourceNodeList *sbom.NodeList

	switch v := element.Value().(type) {
	case *sbom.Document:
		sourceNodeList = v.NodeList
	case *elements.Document:
		sourceNodeList = v.Document.NodeList
	case *sbom.NodeList:
		sourceNodeList = v
	case *elements.NodeList:
		sourceNodeList = v.NodeList
	case *elements.Node:
		sourceNodeList = &sbom.NodeList{
			RootElements: []string{},
		}

		if v.Node.Type == t {
			sourceNodeList.AddNode(v.Node)
			sourceNodeList.RootElements = append(sourceNodeList.RootElements, v.Id)
		}

		return elements.NodeList{
			NodeList: sourceNodeList,
		}, nil

	default:
		return elements.NodeList{}, fmt.Errorf("unable to list packages (unsupported type?) %T", element.Value())
	}
	resultNodeList := elements.NodeList{
		NodeList: &sbom.NodeList{
			RootElements: []string{},
			Edges:        sourceNodeList.Edges,
		},
	}

	for _, n := range sourceNodeList.Nodes {
		if n.Type == t {
			resultNodeList.AddNode(n)
		}
	}

	cleanEdges(&resultNodeList)
	reconnectOrphanNodes(&resultNodeList)
	return resultNodeList, nil
}

// ToDocument converts an element into a fill document. This is useful when
// bomshell needs to convert its results to a document to output them as an SBOM
var ToDocument = func(lhs ref.Val) ref.Val {
	var nodelist *elements.NodeList
	switch v := lhs.Value().(type) {
	case *sbom.NodeList:
		nodelist = &elements.NodeList{NodeList: v}
	case *elements.NodeList:
		nodelist = v
	case *elements.Node:
		nodelist = v.ToNodeList()
	case *sbom.Node:
		nodelist = elements.Node{Node: v}.ToNodeList()
	default:
		return types.NewErr("unable to convert element to document")
	}

	// Here we reconnect all orphaned nodelists to the root of the
	// nodelist. The produced document will describe all elements of
	// the nodelist except for those which are already related to other
	// nodes in the graph.
	reconnectOrphanNodes(nodelist)

	doc := elements.Document{
		Document: &sbom.Document{
			Metadata: &sbom.Metadata{
				Id:      "",
				Version: "1",
				Name:    "bomshell generated document",
				Date:    timestamppb.Now(),
				Tools: []*sbom.Tool{
					{
						Name:    "bomshell",
						Version: version.GetVersionInfo().GitVersion,
						Vendor:  "Chainguard Labs",
					},
				},
				Authors: []*sbom.Person{},
				Comment: "This document was generated by bomshell from a protobom nodelist",
			},
			NodeList: nodelist.NodeList,
		},
	}

	return doc
}

var LoadSBOM = func(_, pathVal ref.Val) ref.Val {
	path, ok := pathVal.Value().(string)
	if !ok {
		return types.NewErr("argument to element by id has to be a string")
	}

	f, err := os.Open(path)
	if err != nil {
		return types.NewErr("opening SBOM file: %w", err)
	}

	doc, err := loader.ReadSBOM(f)
	if err != nil {
		return types.NewErr("loading document: %w", err)
	}
	return elements.Document{
		Document: doc,
	}
}

var NodesByPurlType = func(lhs, rhs ref.Val) ref.Val {
	purlType, ok := rhs.Value().(string)
	if !ok {
		return types.NewErr("argument to GetNodesByPurlType must be a string")
	}

	var nl *sbom.NodeList
	switch v := lhs.Value().(type) {
	case *sbom.Document:
		nl = v.NodeList.GetNodesByPurlType(purlType)
	case *sbom.NodeList:
		nl = v.GetNodesByPurlType(purlType)
	default:
		return types.NewErr("method unsupported on type %T", lhs.Value())
	}

	return elements.NodeList{
		NodeList: nl,
	}
}

// RelateNodeListAtID relates a nodelist at the specified ID
var RelateNodeListAtID = func(vals ...ref.Val) ref.Val {
	if len(vals) != 4 {
		return types.NewErr("invalid number of arguments for RealteAtNodeListAtID")
	}
	id, ok := vals[2].Value().(string)
	if !ok {
		return types.NewErr("node id has to be a string")
	}
	// relType
	_, ok = vals[3].Value().(string)
	if !ok {
		return types.NewErr("relationship type has has to be a string")
	}

	nodelist, ok := vals[1].(elements.NodeList)
	if !ok {
		return types.NewErr("could not cast nodelist")
	}

	switch v := vals[0].Value().(type) {
	case *sbom.Document:
		// FIXME: Lookup reltype
		if err := v.NodeList.RelateNodeListAtID(nodelist.Value().(*sbom.NodeList), id, sbom.Edge_dependsOn); err != nil {
			return types.NewErr(err.Error())
		}
		return elements.Document{
			Document: v,
		}
	case *sbom.NodeList:
		if err := v.RelateNodeListAtID(nodelist.Value().(*sbom.NodeList), id, sbom.Edge_dependsOn); err != nil {
			return types.NewErr(err.Error())
		}
		return elements.NodeList{
			NodeList: v,
		}
	default:
		return types.NewErr("method unsupported on type %T", vals[0].Value())
	}
}
