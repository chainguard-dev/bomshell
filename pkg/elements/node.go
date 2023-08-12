// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package elements

import (
	"fmt"
	"reflect"

	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

var (
	NodeObject    = decls.NewObjectType("bomsquad.protobom.Node")
	NodeTypeValue = types.NewTypeValue("bomsquad.protobom.Node")
	NodeType      = cel.ObjectType("bomsquad.protobom.Node")
)

type Node struct {
	*sbom.Node
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (n Node) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	if reflect.TypeOf(n).AssignableTo(typeDesc) {
		return n, nil
	} else if reflect.TypeOf(n.Node).AssignableTo(typeDesc) {
		return n.Node, nil
	}

	return nil, fmt.Errorf("type conversion error from 'Node' to '%v'", typeDesc)
}

// ConvertToType implements ref.Val.ConvertToType.
func (n Node) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case NodeTypeValue:
		return n
	case types.TypeType:
		return NodeTypeValue
	}
	return types.NewErr("type conversion error from '%s' to '%s'", NodeListTypeValue, typeVal)
}

// Equal implements ref.Val.Equal.
func (n Node) Equal(other ref.Val) ref.Val {
	_, ok := other.(Node)
	if !ok {
		return types.MaybeNoSuchOverloadErr(other)
	}

	// TODO: Moar tests like:
	// return types.Bool(d.URL.String() == otherDur.URL.String())
	return types.True
}

func (n Node) Type() ref.Type {
	return NodeTypeValue
}

// Value implements ref.Val.Value.
func (n Node) Value() interface{} {
	return n
}

// ToNodeList returns a new NodeList with the node as the only member
func (n Node) ToNodeList() *NodeList {
	return &NodeList{
		NodeList: &sbom.NodeList{
			Nodes:        []*sbom.Node{n.Node},
			Edges:        []*sbom.Edge{},
			RootElements: []string{n.Id},
		},
	}
}
