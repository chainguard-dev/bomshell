// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package functions

import (
	"fmt"
	"testing"

	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/stretchr/testify/require"
)

// ToNodeList takes a node and returns a new NodeList
// with that nodelist with the node as the only member.
func TestToNodeList(t *testing.T) {
	for _, tc := range []struct {
		name string
		sut  ref.Val
	}{
		{
			name: "doc",
			sut: &elements.Document{
				Document: &sbom.Document{
					Metadata: &sbom.Metadata{},
					NodeList: &sbom.NodeList{},
				},
			},
		},
		{
			name: "nodelist",
			sut: &elements.NodeList{
				NodeList: &sbom.NodeList{},
			},
		},
		{
			name: "node",
			sut: &elements.Node{
				Node: &sbom.Node{},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := ToNodeList(tc.sut)
			require.NotNil(t, res)
			require.Equal(t, "elements.NodeList", fmt.Sprintf("%T", res), res)
		})
	}
}

func TestNodeByID(t *testing.T) {
	node := &sbom.Node{
		Id: "mynode",
	}
	nl := &sbom.NodeList{
		Nodes:        []*sbom.Node{node},
		Edges:        []*sbom.Edge{},
		RootElements: []string{},
	}

	for _, tc := range []struct {
		name string
		sut  ref.Val
	}{
		{
			name: "doc",
			sut: &elements.Document{
				Document: &sbom.Document{
					NodeList: nl,
				},
			},
		},
		{
			name: "nodelist",
			sut: &elements.NodeList{
				NodeList: nl,
			},
		},
		{
			name: "node",
			sut: &elements.Node{
				Node: node,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := NodeByID(tc.sut, types.String("mynode"))
			require.NotNil(t, res)
			require.Equal(t, "elements.Node", fmt.Sprintf("%T", res), res)
			require.Equal(t, "mynode", res.Value().(*sbom.Node).Id)
		})
	}
}
