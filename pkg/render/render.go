// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package render

import (
	"fmt"
	"strings"

	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/chainguard-dev/bomshell/pkg/elements"
	"github.com/google/cel-go/common/types"
)

type RendererOptions struct {
	ListNodes bool
}

type Renderer interface {
	Display(any) string
}

func NewTTY() *TTY {
	return &TTY{
		Options: RendererOptions{
			ListNodes: false,
		},
	}
}

type TTY struct {
	Options RendererOptions
}

func (tty *TTY) Display(result any) string {
	switch v := result.(type) {
	case nil:
		return "<nil>"
	case types.String:
		return v.Value().(string)
	case *elements.Document:
		return tty.Display(v.Document)
	case *elements.NodeList:
		return tty.Display(v.NodeList)
	case elements.NodeList:
		return tty.Display(v.NodeList)
	case *elements.Node:
		return tty.Display(v.Node)
	case elements.Node:
		return tty.Display(v.Node)
	case *sbom.Document:
		ret := "protobom Document\n"
		ret += fmt.Sprintf("Document ID: %s", v.Metadata.Id)
		ret += "\n" + tty.Display(v.NodeList)
		return ret
	case *sbom.NodeList:
		ret := "protobom NodeList\n"
		ret += fmt.Sprintf("Root Elements: %d\n", len(v.GetRootElements()))
		ret += fmt.Sprintf("Number of nodes: %d (%d packages %d files)\n", len(v.Nodes), numPackages(v), numFiles(v))
		ptypes := purlTypes(v)
		ret += "Package URL types: "
		if len(ptypes) > 0 {
			for t, n := range ptypes {
				ret += fmt.Sprintf("%s: %d ", t, n)
			}
			ret += "\n"
		}

		if tty.Options.ListNodes {
			for _, n := range v.Nodes {
				ret += fmt.Sprintf(" Node %s %s\n", n.Id, n.Purl())
			}
		}
		return ret
	case *sbom.Node:
		ret := "Node Information:\n"
		ret += fmt.Sprintf("ID: %s (%s)\n", v.Id, []string{"file", "package"}[v.Type])
		ret += fmt.Sprintf("Package URL: %s\n", v.Purl())
		ret += fmt.Sprintf("Name: %s\n", v.Name)
		return ret
	default:
		ret := fmt.Sprintf("type %T renderer not implemented yet!\n", result)
		return ret
	}
}

func purlTypes(nl *sbom.NodeList) map[string]int {
	purls := map[string]int{}
	for _, n := range nl.Nodes {
		p := n.Purl()
		if p == "" {
			continue
		}
		ps := strings.TrimPrefix(string(p), "pkg:/")
		ps = strings.TrimPrefix(ps, "pkg:")

		parts := strings.Split(ps, "/")

		if _, ok := purls[parts[0]]; !ok {
			purls[parts[0]] = 0
		}
		purls[parts[0]]++
	}
	return purls
}

func numFiles(nl *sbom.NodeList) int {
	c := 0
	for _, n := range nl.Nodes {
		if n.Type == sbom.Node_FILE {
			c++
		}
	}
	return c
}

func numPackages(nl *sbom.NodeList) int {
	c := 0
	for _, n := range nl.Nodes {
		if n.Type == sbom.Node_PACKAGE {
			c++
		}
	}
	return c
}
