package lang

import (
	"fmt"
	"github.com/ixre/gof/db/orm"
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

func (p Protobuf) ParsePkg(pkg string) string {
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
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeFloat32:
		return "double"
	case orm.TypeFloat64:
		return "double"
	case orm.TypeInt16:
		return "int32"
	case orm.TypeInt32:
		return "int32"
	case orm.TypeInt64:
		return "int64"
	case orm.TypeString:
		return "string"
	case orm.TypeDecimal:
		return "decimal"
	case orm.TypeDateTime:
		return "string"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
