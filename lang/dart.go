package lang

import (
	"github.com/ixre/gof/db/db"
)

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

type dart struct {
}

func (d dart) ParseType(typeId int) string {
	switch typeId {
	case db.TypeString:
		return "String"
	case db.TypeBoolean:
		return "bool"
	case db.TypeInt16, db.TypeInt32, db.TypeInt64:
		return "int"
	case db.TypeFloat32, db.TypeFloat64, db.TypeDecimal:
		return "num"
	case db.TypeDateTime:
		return "DateTime"
	}
	return "dynamic"
}

func (d dart) SqlMapType(typeId int, len int) string {
	panic("not support for dart")
}

func (d dart) DefaultValue(typeId int) string {
	switch typeId {
	case db.TypeString:
		return "\"\""
	case db.TypeBoolean:
		return "false"
	case db.TypeInt16, db.TypeInt32, db.TypeInt64:
		return "0"
	case db.TypeFloat32, db.TypeFloat64, db.TypeDecimal:
		return "0.0"
	case db.TypeDateTime:
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
