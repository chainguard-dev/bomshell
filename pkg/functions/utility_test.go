package functions

import (
	"testing"

	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/stretchr/testify/require"
)

func TestCleanEdges(t *testing.T) {
	for _, tc := range []struct {
		sut      *elements.NodeList
		expected *elements.NodeList
	}{
		// Edge does not need to be modified
		{
			sut: &elements.NodeList{
				NodeList: &sbom.NodeList{
					Nodes: []*sbom.Node{
						{Id: "node1"}, {Id: "node2"},
					},
					Edges: []*sbom.Edge{
						{
							Type: 0,
							From: "node1",
							To:   []string{"node2"},
						},
					},
					RootElements: []string{"node1"},
				},
			},
			expected: &elements.NodeList{
				NodeList: &sbom.NodeList{
					Nodes: []*sbom.Node{
						{Id: "node1"}, {Id: "node2"},
					},
					Edges: []*sbom.Edge{
						{
							Type: 0,
							From: "node1",
							To:   []string{"node2"},
						},
					},
					RootElements: []string{"node1"},
				},
			},
		},
		// Edge contains a broken To
		{
			sut: &elements.NodeList{
				NodeList: &sbom.NodeList{
					Nodes: []*sbom.Node{
						{Id: "node1"}, {Id: "node2"},
					},
					Edges: []*sbom.Edge{
						{
							Type: 0,
							From: "node1",
							To:   []string{"node2", "node3"},
						},
					},
					RootElements: []string{"node1"},
				},
			},
			expected: &elements.NodeList{
				NodeList: &sbom.NodeList{
					Nodes: []*sbom.Node{
						{Id: "node1"}, {Id: "node2"},
					},
					Edges: []*sbom.Edge{
						{
							Type: 0,
							From: "node1",
							To:   []string{"node2"},
						},
					},
					RootElements: []string{"node1"},
				},
			},
		},
		// Edge contains a broken From
		{
			sut: &elements.NodeList{
				NodeList: &sbom.NodeList{
					Nodes: []*sbom.Node{
						{Id: "node1"}, {Id: "node2"},
					},
					Edges: []*sbom.Edge{
						{
							Type: 0,
							From: "node3",
							To:   []string{"node1"},
						},
					},
					RootElements: []string{"node1"},
				},
			},
			expected: &elements.NodeList{
				NodeList: &sbom.NodeList{
					Nodes: []*sbom.Node{
						{Id: "node1"}, {Id: "node2"},
					},
					Edges:        []*sbom.Edge{},
					RootElements: []string{"node1"},
				},
			},
		},
	} {
		cleanEdges(tc.sut)
		require.Equal(t, tc.sut, tc.expected)
	}
}
