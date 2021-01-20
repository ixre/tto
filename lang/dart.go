package lang

import "github.com/ixre/gof/db/orm"

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : dart.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-12-16 12:40
 * description :
 * history :
 */

var _ Lang = new(dart)
type dart struct{

}

func (d dart) ParseType(typeId int) string {
	switch typeId {
	case orm.TypeString:
		return "String"
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeInt16:
		return "int"
	case orm.TypeInt32:
		return "int"
	case orm.TypeInt64:
		return "BigInt"
	case orm.TypeFloat32:
		return "num"
	case orm.TypeFloat64, orm.TypeDecimal:
		return "num"
	case orm.TypeDateTime:
		return "DateTime"
	}
	return "dynamic"
}

func (d dart) SqlMapType(typeId int, len int) string {
	panic("not support for dart")
}

func (d dart) DefaultValue(typeId int) string {
	switch typeId {
	case orm.TypeString:
		return "\"\""
	case orm.TypeBoolean:
		return "false"
	case orm.TypeInt16,orm.TypeInt32:
		return "0"
	case orm.TypeInt64:
		return "BigInt.from(0)"
	case orm.TypeFloat32,orm.TypeFloat64, orm.TypeDecimal:
		return "0.0"
	case orm.TypeDateTime:
		return "DateTime.now()"
	}
	return "{}"
}

func (d dart) PkgPath(pkg string) string {
	return pkg
}

func (d dart) PkgName(pkg string) string {
	return ""
}
