package lang

import (
	"github.com/ixre/gof/db/orm"
	"strconv"
)

type Thrift struct {
}

func (t Thrift) ParsePkg(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (t Thrift) ParseType(typeId int) string {
	return ThriftTypes(typeId)
}

func (t Thrift) DefaultValue(typeId int) string {
	panic("not support")
}

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : thrift.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-20 11:04
 * description :
 * history :
 */
var _ Lang = new(Thrift)

func ThriftTypes(typeId int) string {
	switch typeId {
	case orm.TypeString:
		return "string"
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeInt16:
		return "i16"
	case orm.TypeInt32:
		return "i32"
	case orm.TypeInt64:
		return "i64"
	case orm.TypeFloat32:
		return "f32"
	case orm.TypeFloat64:
		return "f64"
	case orm.TypeDecimal:
		return "f64"
	case orm.TypeDateTime:
		return "string"
	}
	return strconv.Itoa(typeId)
}
