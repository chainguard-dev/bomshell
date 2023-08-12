// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2023 Chainguard Inc

package elements

import (
	"errors"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

var (
	BomshellObject = decls.NewObjectType("bomshell")
	BomshellType   = cel.ObjectType("bomshell", traits.ReceiverType)
)

type Bomshell struct{}

func (bs Bomshell) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	return bs, errors.New("bomshell cannot be converted to native")
}

// ConvertToType implements ref.Val.ConvertToType.
func (bs Bomshell) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case DocumentType:
		return bs
	case types.TypeType:
		return BomshellType
	}
	return types.NewErr("type conversion error not allowed in bomshell")
}

// Equal implements ref.Val.Equal.
func (bs Bomshell) Equal(other ref.Val) ref.Val {
	return types.NewErr("bomshell objects cannot be compared")
}

func (bs Bomshell) Type() ref.Type {
	return BomshellType
}

// Value implements ref.Val.Value.
func (bs Bomshell) Value() interface{} {
	return bs
}
