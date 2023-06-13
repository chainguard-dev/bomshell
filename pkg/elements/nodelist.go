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
	NodeListObject    = decls.NewObjectType("bomsquad.protobom.NodeList")
	NodeListTypeValue = types.NewTypeValue("bomsquad.protobom.NodeList")
	NodeListType      = cel.ObjectType("bomsquad.protobom.NodeList")
)

type NodeList struct {
	*sbom.NodeList
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (nl NodeList) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	if reflect.TypeOf(nl).AssignableTo(typeDesc) {
		return nl, nil
	} else if reflect.TypeOf(nl.NodeList).AssignableTo(typeDesc) {
		return nl.NodeList, nil
	}
	//if reflect.TypeOf("").AssignableTo(typeDesc) {
	//		return d.URL.String(), nil
	//	}
	return nil, fmt.Errorf("type conversion error from 'NodeList' to '%v'", typeDesc)
}

// ConvertToType implements ref.Val.ConvertToType.
func (nl NodeList) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case NodeListTypeValue:
		return nl
	case types.TypeType:
		return NodeListTypeValue

	}
	return types.NewErr("type conversion error from '%s' to '%s'", NodeListTypeValue, typeVal)
}

// Equal implements ref.Val.Equal.
func (nl NodeList) Equal(other ref.Val) ref.Val {
	// otherDur, ok := other.(NodeList)
	_, ok := other.(NodeList)
	if !ok {
		return types.MaybeNoSuchOverloadErr(other)
	}

	// TODO: Moar tests like:
	// return types.Bool(d.URL.String() == otherDur.URL.String())
	return types.True
}

// Type implements ref.Val.Type.
func (nl NodeList) Type() ref.Type {
	return NodeListTypeValue
}

// Value implements ref.Val.Value.
func (nl NodeList) Value() interface{} {
	return nl
}
