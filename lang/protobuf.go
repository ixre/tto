package lang

import (
	"fmt"

	"github.com/ixre/gof/db/db"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : protobuf
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-20 11:04
 * description :
 * history :
 */

type Protobuf struct {
}

func (p Protobuf) SqlMapType(typeId int, len int) string {
	return p.ParseType(typeId)
}

func (p Protobuf) PkgName(pkg string) string {
	return p.PkgPath(pkg)
}
func (p Protobuf) PkgPath(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (p Protobuf) ParseType(typeId int) string {
	return ProtobufTypes(typeId)
}

func (p Protobuf) DefaultValue(typeId int) string {
	panic("not support for protobuf")
}

var _ Lang = new(Protobuf)

func ProtobufTypes(typeId int) string {
	switch typeId {
	case db.TypeBoolean:
		return "bool"
	case db.TypeFloat32:
		return "double"
	case db.TypeFloat64:
		return "double"
	case db.TypeInt16:
		return "int32"
	case db.TypeInt32:
		return "int32"
	case db.TypeInt64:
		return "int64"
	case db.TypeString:
		return "string"
	case db.TypeDecimal:
		return "decimal"
	case db.TypeDateTime:
		return "string"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
