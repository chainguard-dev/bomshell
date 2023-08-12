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
	DocumentObject    = decls.NewObjectType("bomsquad.protobom.Document")
	DocumentTypeValue = types.NewTypeValue("bomsquad.protobom.Document")
	DocumentType      = cel.ObjectType("bomsquad.protobom.Document")
)

type Document struct {
	*sbom.Document
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (d Document) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	if reflect.TypeOf(d).AssignableTo(typeDesc) {
		return d, nil
	} else if reflect.TypeOf(d.Document).AssignableTo(typeDesc) {
		return d.Document, nil
	}

	return nil, fmt.Errorf("type conversion error from 'Document' to '%v'", typeDesc)
}

// ConvertToType implements ref.Val.ConvertToType.
func (d Document) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case DocumentTypeValue:
		return d
	case types.TypeType:
		return DocumentTypeValue

	}
	return types.NewErr("type conversion error from '%s' to '%s'", NodeListTypeValue, typeVal)
}

// Equal implements ref.Val.Equal.
func (d Document) Equal(other ref.Val) ref.Val {
	_, ok := other.(Document)
	if !ok {
		return types.MaybeNoSuchOverloadErr(other)
	}

	// TODO: Moar tests like:
	// return types.Bool(d.URL.String() == otherDur.URL.String())
	return types.True
}

func (d Document) Type() ref.Type {
	return DocumentTypeValue
}

// Value implements ref.Val.Value.
func (d Document) Value() interface{} {
	return d
}
