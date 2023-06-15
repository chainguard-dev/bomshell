package functions

import (
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/chainguard-dev/bomshell/pkg/elements"
)

// cleanEdges removes all edges that have broken Froms and remove
// any destination IDs from elements no not in the NodeList.
func cleanEdges(nl *elements.NodeList) {
	// First copy the nodelist edges
	newEdges := []*sbom.Edge{}

	// Build a catalog of the elements ids
	idDict := map[string]struct{}{}
	for i := range nl.Nodes {
		idDict[nl.Nodes[i].Id] = struct{}{}
	}

	// Now list all edges and rebuild the list
	for _, edge := range nl.Edges {
		newTos := []string{}
		if _, ok := idDict[edge.From]; !ok {
			continue
		}

		for _, s := range edge.To {
			if _, ok := idDict[s]; ok {
				newTos = append(newTos, s)
			}
		}

		if len(newTos) == 0 {
			continue
		}

		edge.To = newTos
		newEdges = append(newEdges, edge)
	}

	nl.Edges = newEdges
}
